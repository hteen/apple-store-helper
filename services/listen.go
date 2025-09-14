package services

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
	"github.com/faiface/beep"
	"github.com/faiface/beep/mp3"
	"github.com/faiface/beep/speaker"
	"github.com/golang-module/carbon"
	"github.com/parnurzeal/gorequest"
	"github.com/tidwall/gjson"

	"apple-store-helper/model"
	"apple-store-helper/theme"
	"apple-store-helper/view"
)

const (
	StatusOutStock = "无货"
	StatusInStock  = "有货"
	StatusWait     = "等待"

	Pause   = "暂停"
	Running = "监听中"
)

var Listen = listenService{
	items:           map[string]ListenItem{},
	Status:          binding.NewString(),
	Area:            model.Areas[0],
	Logs:            widget.NewLabel(""),
	DetectThreshold: 3,   // 默认检测3次
	TimeThreshold:   1,   // 默认1分钟
	RefreshInterval: 3,   // 默认3秒刷新一次
	BatchInterval:   500, // 默认500毫秒门店请求间隔
}

type listenService struct {
	items           map[string]ListenItem
	Status          binding.String
	Area            model.Area
	Logs            *widget.Label
	BarkNotifyUrl   string
	detectCounts    map[string]int             // 检测次数计数器
	firstDetect     map[string]carbon.DateTime // 第一次检测时间
	DetectThreshold int                        // 检测次数阈值（默认3次）
	TimeThreshold   int                        // 时间阈值（默认1分钟）
	RefreshInterval int                        // 刷新间隔（秒，默认3秒）
	BatchInterval   int                        // 同一轮内门店请求间隔（毫秒，默认500毫秒）
}

type ListenItem struct {
	Store   model.Store
	Product model.Product
	Status  string
	Time    carbon.DateTime
}

func (s *listenService) Add(areaTitle string, storeTitle string, productTitle string, barkNotifyUrl string) {

	store := Store.GetStore(areaTitle, storeTitle)
	product := Product.GetProduct(areaTitle, productTitle)

	uniqKey := store.StoreNumber + "." + product.Code

	if s.items[uniqKey].Store.StoreNumber == "" {
		s.items[uniqKey] = ListenItem{
			Store:   store,
			Product: product,
			Status:  StatusWait,
		}
	}

	s.BarkNotifyUrl = barkNotifyUrl
	s.UpdateLogStr()
}

func (s *listenService) Clean() {
	s.items = map[string]ListenItem{}
	s.detectCounts = map[string]int{}
	s.firstDetect = map[string]carbon.DateTime{}
	s.UpdateLogStr()
}

func (s *listenService) SetListenItems(items map[string]ListenItem) {
	s.items = items
	s.UpdateLogStr()
}

func (s *listenService) GetListenItems() map[string]ListenItem {
	return s.items
}

func (s *listenService) SetThresholds(detectThreshold, timeThreshold int) {
	s.DetectThreshold = detectThreshold
	s.TimeThreshold = timeThreshold
}

func (s *listenService) GetThresholds() (int, int) {
	return s.DetectThreshold, s.TimeThreshold
}

func (s *listenService) SetRefreshInterval(interval int) {
	s.RefreshInterval = interval
}

func (s *listenService) GetRefreshInterval() int {
	return s.RefreshInterval
}

func (s *listenService) SetBatchInterval(interval int) {
	s.BatchInterval = interval
}

func (s *listenService) GetBatchInterval() int {
	return s.BatchInterval
}

func (s *listenService) UpdateLogStr() {
	var str string

	for key, item := range s.items {
		detectInfo := ""
		if count, exists := s.detectCounts[key]; exists && count > 0 {
			detectInfo = fmt.Sprintf(" (检测%d次)", count)
		}

		str += fmt.Sprintf(
			"[%s] %s %s %s%s %s",
			item.Status,
			item.Time,
			item.Store.CityStoreName,
			item.Product.Title,
			detectInfo,
			"\n",
		)
	}

	s.Logs.SetText(str)
}

func (s *listenService) UpdateStatus(uniqKey string, status string) {
	item := s.items[uniqKey]
	item.Time = carbon.DateTime{Carbon: carbon.Now(carbon.Shanghai)}
	item.Status = status
	s.items[uniqKey] = item
}

func (s *listenService) Run() {
	s.Status.Set(Pause)

	go func() {
		for {
			if stats, ok := s.Status.Get(); ok == nil && stats == Running && len(s.items) > 0 {
				skus := s.groupByStore()

				for key, item := range s.items {
					status := skus[item.Store.StoreNumber+"."+item.Product.Code]

					if status {
						// 更新检测次数和时间
						s.detectCounts[key]++
						now := carbon.Now(carbon.Shanghai)

						// 如果是第一次检测到有货，记录时间
						if s.detectCounts[key] == 1 {
							s.firstDetect[key] = carbon.DateTime{Carbon: now}
						}

						// 检查是否在指定时间内检测到指定次数
						if s.detectCounts[key] >= s.DetectThreshold {
							// 检查时间是否在指定阈值内
							if s.firstDetect[key].Carbon.DiffInMinutes(now) <= int64(s.TimeThreshold) {
								s.UpdateStatus(key, StatusInStock)
								s.Status.Set(Pause)

								var bagUrl = fmt.Sprintf("https://www.apple.com/%s/shop/bag", s.Area.ShortCode)
								// 进入购物袋
								s.openBrowser(bagUrl)
								msg := fmt.Sprintf("%s %s 有货 (检测%d次)", item.Store.CityStoreName, item.Product.Title, s.detectCounts[key])
								dialog.ShowInformation("匹配成功", msg, view.Window)
								view.App.SendNotification(&fyne.Notification{
									Title:   "有货提醒",
									Content: msg,
								})
								go s.AlertMp3()
								go s.SendPushNotificationByBark("有货提醒", msg, bagUrl)
								break
							} else {
								// 超过一分钟，重置计数
								s.detectCounts[key] = 1
								s.firstDetect[key] = carbon.DateTime{Carbon: now}
							}
						}

						// 更新状态为有货（但可能还没达到3次）
						s.UpdateStatus(key, StatusInStock)
					} else {
						s.UpdateStatus(key, StatusOutStock)
					}
				}

				s.UpdateLogStr()
			}

			time.Sleep(time.Duration(s.RefreshInterval) * time.Second)
		}
	}()
}

func (s *listenService) groupByStore() map[string]bool {
	skus := map[string]bool{}

	defer func() {
		if r := recover(); r != nil {
			log.Println(r)
		}
	}()

	group := map[string][]ListenItem{}
	reqs := map[string]string{}

	for _, item := range s.items {
		group[item.Store.StoreNumber] = append(group[item.Store.StoreNumber], item)
	}

	for storeNumber, items := range group {

		var uri url.URL
		q := uri.Query()
		q.Set("little", "true")
		q.Set("mt", "regular")
		q.Set("store", storeNumber)

		for index, item := range items {
			q.Set("parts."+strconv.FormatInt(int64(index), 10), item.Product.Code)
		}

		queryStr := q.Encode()

		link := fmt.Sprintf(
			"https://www.apple.com/%s/shop/fulfillment-messages?%s",
			s.Area.ShortCode,
			queryStr,
		)

		reqs[storeNumber] = link
	}

	count := len(reqs)
	if count < 1 {
		return skus
	}

	ch := make(chan map[string]bool, count)

	// 按顺序发送请求，每个请求之间添加间隔
	requestCount := 0
	for _, link := range reqs {
		go s.getSkuByLink(ch, link)
		requestCount++
		// 如果不是最后一个请求，添加间隔
		if requestCount < len(reqs) {
			time.Sleep(time.Duration(s.BatchInterval) * time.Millisecond)
		}
	}

	for i := 0; i < count; i++ {
		for key, v := range <-ch {
			skus[key] = v
		}
	}

	return skus
}

func (s *listenService) getSkuByLink(ch chan map[string]bool, skUrl string) {
	skus := map[string]bool{}

	resp, body, errs := gorequest.New().
		Set("referer", "https://www.apple.com/shop/buy-iphone").
		Set("user-agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/94.0.4606.71 Safari/537.36").
		Timeout(time.Second * 3).Get(skUrl).End()
	if len(errs) > 0 {
		log.Println(errs)
		ch <- skus
		return
	}

	log.Println(resp.Status, skUrl)
	for _, result := range gjson.Get(body, "body.content.pickupMessage.stores").Array() {
		for productCode, availability := range result.Get("partsAvailability").Map() {
			uniqKey := fmt.Sprintf("%s.%s", result.Get("storeNumber").String(), productCode)
			skus[uniqKey] = availability.Get("messageTypes.compact.storeSelectionEnabled").Bool()
		}
	}

	ch <- skus
}

// 型号对应预约地址
//func (s *listenService) model2Url(productType string) string {
//	// https://www.apple.com.cn/shop/buy-iphone/iphone-16
//	// https://www.apple.com.cn/shop/buy-iphone/iphone-16-pro
//
//	var t string
//	switch productType {
//	case "iphone16promax", "iphone16pro":
//		t = "iphone-16-pro"
//	case "iphone16":
//		t = "iphone-16"
//	}
//
//	return fmt.Sprintf(
//		"https://www.apple.com/%s/shop/buy-iphone/%s",
//		s.Area.ShortCode,
//		t,
//	)
//}

func (s *listenService) openBrowser(link string) {
	parse, err := url.Parse(link)
	if err != nil {
		dialog.ShowError(err, view.Window)
		return
	}

	err = view.App.OpenURL(parse)
	if err != nil {
		dialog.ShowError(err, view.Window)
		return
	}
}

func (s *listenService) AlertMp3() {
	reader := bytes.NewReader(theme.Mp3().Content())
	streamer, _, err := mp3.Decode(ioutil.NopCloser(reader))
	if err != nil {
		panic(err)
	}
	defer streamer.Close()

	done := make(chan bool)
	speaker.Play(beep.Seq(streamer, beep.Callback(func() {
		done <- true
	})))
	<-done
}

func (s *listenService) SendPushNotificationByBark(title string, content string, bagUrl string) {

	if len(s.BarkNotifyUrl) <= 0 {
		return
	}

	apiUrl := fmt.Sprintf("%s/%s/%s?url=%s", strings.TrimRight(s.BarkNotifyUrl, "/"), title, content, bagUrl)

	response, err := http.Get(apiUrl)
	if err != nil {
		panic(err)
	}
	defer response.Body.Close()
}
