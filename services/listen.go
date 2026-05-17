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

	DefaultIntervalSeconds = 3
)

// IntervalOptions 用户可选的轮询间隔档位(秒)。
// 1s 偏激进,容易触发 Apple 限流;>=3s 是常见安全档位;30/60s 适合长期挂机。
var IntervalOptions = []int{1, 3, 5, 10, 30, 60}

// IntervalLabel 把秒数转成下拉框显示文案。
func IntervalLabel(seconds int) string {
	if seconds == DefaultIntervalSeconds {
		return fmt.Sprintf("%d 秒 (默认)", seconds)
	}
	return fmt.Sprintf("%d 秒", seconds)
}

// IntervalLabels 生成下拉框完整选项列表。
func IntervalLabels() []string {
	labels := make([]string, 0, len(IntervalOptions))
	for _, s := range IntervalOptions {
		labels = append(labels, IntervalLabel(s))
	}
	return labels
}

// ParseIntervalLabel 把下拉框选中的文案反解成秒数;无法识别时回退到默认值。
func ParseIntervalLabel(label string) int {
	for _, s := range IntervalOptions {
		if IntervalLabel(s) == label {
			return s
		}
	}
	return DefaultIntervalSeconds
}

var Listen = listenService{
	items:           map[string]ListenItem{},
	Status:          binding.NewString(),
	Area:            model.Areas[0],
	Logs:            widget.NewLabel(""),
	IntervalSeconds: binding.NewInt(),
}

type listenService struct {
	items           map[string]ListenItem
	Status          binding.String
	Area            model.Area
	Logs            *widget.Label
	BarkNotifyUrl   string
	IntervalSeconds binding.Int
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
	s.UpdateLogStr()
}

func (s *listenService) SetListenItems(items map[string]ListenItem) {
	s.items = items
	s.UpdateLogStr()
}

func (s *listenService) GetListenItems() map[string]ListenItem {
	return s.items
}

func (s *listenService) UpdateLogStr() {
	var str string

	for _, item := range s.items {

		str += fmt.Sprintf(
			"[%s] %s %s %s %s",
			item.Status,
			item.Time,
			item.Store.CityStoreName,
			item.Product.Title,
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

// IntervalDuration 返回当前生效的轮询间隔;未配置或异常值时回退到默认值。
func (s *listenService) IntervalDuration() time.Duration {
	seconds, err := s.IntervalSeconds.Get()
	if err != nil || seconds <= 0 {
		seconds = DefaultIntervalSeconds
	}
	return time.Duration(seconds) * time.Second
}

func (s *listenService) Run() {
	s.Status.Set(Pause)

	go func() {
		// 空闲时用短间隔轮询状态变化,避免用户按下"开始"后还要等一整个监听周期才生效。
		const idleSleep = 200 * time.Millisecond

		for {
			stats, ok := s.Status.Get()
			running := ok == nil && stats == Running && len(s.items) > 0

			if !running {
				time.Sleep(idleSleep)
				continue
			}

			skus := s.groupByStore()

			for key, item := range s.items {
				status := skus[item.Store.StoreNumber+"."+item.Product.Code]

				if status {
					s.UpdateStatus(key, StatusInStock)
					s.Status.Set(Pause)

					var bagUrl = fmt.Sprintf("https://www.apple.com/%s/shop/bag", s.Area.ShortCode)
					// 进入购物袋
					s.openBrowser(bagUrl)
					msg := fmt.Sprintf("%s %s 有货", item.Store.CityStoreName, item.Product.Title)
					dialog.ShowInformation("匹配成功", msg, view.Window)
					view.App.SendNotification(&fyne.Notification{
						Title:   "有货提醒",
						Content: msg,
					})
					go s.AlertMp3()
					go s.SendPushNotificationByBark("有货提醒", msg, bagUrl)
					break
				} else {
					s.UpdateStatus(key, StatusOutStock)
				}
			}

			s.UpdateLogStr()

			time.Sleep(s.IntervalDuration())
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

	for _, link := range reqs {
		go s.getSkuByLink(ch, link)
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
		Get(skUrl).
		Set("referer", "https://www.apple.com/shop/buy-iphone").
		Set("user-agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/94.0.4606.71 Safari/537.36").
		Timeout(time.Second * 3).End()
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
