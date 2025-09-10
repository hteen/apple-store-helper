package services

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"net/url"
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

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

var Listen = listenService{
	items:  map[string]ListenItem{},
	Status: binding.NewString(),
	Area:   model.Areas[0],
	Logs:   widget.NewLabel(""),
}

type listenService struct {
	items         map[string]ListenItem
	Status        binding.String
	Area          model.Area
	Logs          *widget.Label
	BarkNotifyUrl string
	requestCount  int64 // 请求计数器
}

type ListenItem struct {
	Store   model.Store
	Product model.Product
	Status  string
	Time    carbon.DateTime
}

func (s *listenService) Add(areaTitle string, storeTitle string, productTitle string) {

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

	s.UpdateLogStr()
}

func (s *listenService) AddWithProductInfo(areaTitle string, storeTitle string, productTitle string, productCode string, productType string) {
	store := Store.GetStore(areaTitle, storeTitle)

	product := model.Product{
		Title: productTitle,
		Code:  productCode,
		Type:  productType,
	}

	uniqKey := store.StoreNumber + "." + product.Code

	if s.items[uniqKey].Store.StoreNumber == "" {
		s.items[uniqKey] = ListenItem{
			Store:   store,
			Product: product,
			Status:  StatusWait,
		}
	}

	s.UpdateLogStr()
}

func (s *listenService) AddWithStoreInfo(store model.Store, productTitle string, productCode string, productType string) {
	// 验证门店信息是否有效
	if store.StoreNumber == "" {
		log.Printf("Error: Invalid store information - StoreNumber is empty")
		return
	}

	product := model.Product{
		Title: productTitle,
		Code:  productCode,
		Type:  productType,
	}

	uniqKey := store.StoreNumber + "." + product.Code

	log.Printf("Adding item with key: %s (Store: %s, Product: %s, Code: %s)",
		uniqKey, store.CityStoreName, productTitle, productCode)

	if s.items[uniqKey].Store.StoreNumber == "" {
		s.items[uniqKey] = ListenItem{
			Store:   store,
			Product: product,
			Status:  StatusWait,
		}
	}

	s.UpdateLogStr()
}

func (s *listenService) SetBarkUrl(barkUrl string) {
	s.BarkNotifyUrl = barkUrl
}

func (s *listenService) SetListenItems(items map[string]ListenItem) {
	s.items = items
	s.UpdateLogStr()
}

func (s *listenService) GetListenItems() map[string]ListenItem {
	return s.items
}

func (s *listenService) Clean() {
	s.items = map[string]ListenItem{}
	s.UpdateLogStr()
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
	item, exists := s.items[uniqKey]
	if !exists {
		log.Printf("WARNING: UpdateStatus called with non-existent key: %s", uniqKey)
		return
	}
	oldStatus := item.Status
	item.Time = carbon.DateTime{Carbon: carbon.Now(carbon.Shanghai)}
	item.Status = status
	s.items[uniqKey] = item
	log.Printf("Updated status for %s: %s -> %s (Store: %s, Product: %s)",
		uniqKey, oldStatus, status, item.Store.CityStoreName, item.Product.Title)
}

func (s *listenService) Run() {
	s.Status.Set(Pause)

	go func() {
		for {
			if stats, ok := s.Status.Get(); ok == nil && stats == Running && len(s.items) > 0 {
				skus := s.groupByStore()

				// 首先检查是否有任何店铺有货（用于location查询）
				availableStores := make(map[string][]string) // productCode -> []storeNumbers
				for skuKey, available := range skus {
					if available {
						parts := strings.Split(skuKey, ".")
						if len(parts) == 2 {
							storeNumber := parts[0]
							productCode := parts[1]
							// 检查是否是我们监控的产品
							for _, item := range s.items {
								if item.Product.Code == productCode {
									availableStores[productCode] = append(availableStores[productCode], storeNumber)
									break
								}
							}
						}
					}
				}

				// 创建店铺编号到名称的映射
				storeNameMap := make(map[string]string)
				for _, stores := range model.ChinaStores {
					for _, store := range stores {
						storeNameMap[store.StoreNumber] = store.CityStoreName
					}
				}

				// 打印有货店铺信息
				for productCode, stores := range availableStores {
					var storeNames []string
					for _, storeNum := range stores {
						if name, ok := storeNameMap[storeNum]; ok {
							storeNames = append(storeNames, name)
						} else {
							storeNames = append(storeNames, storeNum)
						}
					}
					log.Printf("Product %s available at %d stores: %v", productCode, len(stores), storeNames)
				}

				for key, item := range s.items {
					// 检查指定店铺是否有库存
					hasStock := false
					storeKey := fmt.Sprintf("%s.%s", item.Store.StoreNumber, item.Product.Code)

					if available, exists := skus[storeKey]; exists && available {
						hasStock = true
						log.Printf("Specific store query: Store %s (%s) has stock for product %s",
							item.Store.StoreNumber, item.Store.CityStoreName, item.Product.Code)
					}

					log.Printf("Checking item: key=%s, store=%s (%s), product=%s, hasStock=%v",
						key, item.Store.CityStoreName, item.Store.StoreNumber, item.Product.Title, hasStock)

					if hasStock {
						// 有货（指定店铺有库存）
						s.UpdateStatus(key, StatusInStock)
						s.Status.Set(Pause)

						// 构建提醒消息
						msg := fmt.Sprintf("%s %s 有货", item.Store.CityStoreName, item.Product.Title)

						// 进入购物袋, 手动选择门店
						bagUrl := fmt.Sprintf("https://www.apple.com/%s/shop/bag", s.Area.ShortCode)
						s.openBrowser(bagUrl)
						dialog.ShowInformation("有货提醒", msg, view.Window)
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
			}

			time.Sleep(time.Millisecond * 500)
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

	// 增加请求计数
	s.requestCount++

	// 每10次请求后添加额外延迟，模拟正常用户行为
	if s.requestCount%10 == 0 {
		time.Sleep(time.Duration(2000+rand.Intn(3000)) * time.Millisecond)
		log.Printf("Request count: %d, adding extra delay", s.requestCount)
	}

	group := map[string][]ListenItem{}
	reqs := map[string]string{}

	for _, item := range s.items {
		group[item.Store.StoreNumber] = append(group[item.Store.StoreNumber], item)
	}

	for storeNumber, items := range group {

		var link string

		// 检查是否为中国大陆、香港、日本、新加坡、美国、英国、澳大利亚
		if s.Area.ShortCode == "cn" || s.Area.ShortCode == "hk" || s.Area.ShortCode == "jp" || s.Area.ShortCode == "sg" || s.Area.ShortCode == "us" || s.Area.ShortCode == "uk" || s.Area.ShortCode == "au" {
			// 使用特定店铺查询API格式
			var params []string

			// 1. 固定参数 fae=true
			params = append(params, "fae=true")

			// 2. pl=true（只出现一次）
			params = append(params, "pl=true")

			// 3. 所有产品的mts.i=regular和parts.i=SKU（索引递增）
			for i, item := range items {
				params = append(params, fmt.Sprintf("mts.%d=regular", i))
				// 注意：SKU中的/不编码，保持原样
				params = append(params, fmt.Sprintf("parts.%d=%s", i, item.Product.Code))
			}

			// 4. 添加查询参数
			if len(items) > 0 {
				store := items[0].Store

				// 各地区特殊处理
				if s.Area.ShortCode == "jp" {
					// 日本使用location参数（邮编），cppart参数放在最后
					params = append(params, "location="+store.District)
					log.Printf("Japan location query for store %s (%s) with postal code %s",
						store.StoreNumber, store.CityStoreName, store.District)
				} else {
					// 其他地区使用店铺编号参数
					params = append(params, fmt.Sprintf("store=%s", store.StoreNumber))
					log.Printf("Store-specific query for store %s (%s)", store.StoreNumber, store.CityStoreName)
				}
			}

			// 5. 日本特殊参数cppart放在最后
			if s.Area.ShortCode == "jp" {
				params = append(params, "cppart=UNLOCKED_JP")
			}

			// 构建完整URL - 根据地区使用不同域名
			queryStr := strings.Join(params, "&")
			switch s.Area.ShortCode {
			case "hk":
				link = fmt.Sprintf("https://www.apple.com/hk/shop/fulfillment-messages?%s", queryStr)
			case "jp":
				link = fmt.Sprintf("https://www.apple.com/jp/shop/fulfillment-messages?%s", queryStr)
			case "sg":
				link = fmt.Sprintf("https://www.apple.com/sg/shop/fulfillment-messages?%s", queryStr)
			case "us":
				link = fmt.Sprintf("https://www.apple.com/us/shop/fulfillment-messages?%s", queryStr)
			case "uk":
				link = fmt.Sprintf("https://www.apple.com/uk/shop/fulfillment-messages?%s", queryStr)
			case "au":
				link = fmt.Sprintf("https://www.apple.com/au/shop/fulfillment-messages?%s", queryStr)
			default:
				link = fmt.Sprintf("https://www.apple.com.cn/shop/fulfillment-messages?%s", queryStr)
			}
		}

		reqs[storeNumber] = link
		log.Printf("Store %s URL: %s", storeNumber, link)
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
		result := <-ch
		for key, v := range result {
			skus[key] = v
			log.Printf("Received from channel: key=%s, available=%v", key, v)
		}
	}

	log.Printf("Total SKUs collected: %d", len(skus))
	for k, v := range skus {
		if v {
			log.Printf("Available SKU: %s", k)
		}
	}

	return skus
}

func (s *listenService) getSkuByLink(ch chan map[string]bool, skUrl string) {
	skus := map[string]bool{}

	// 生成随机User-Agent
	userAgents := []string{
		"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36",
		"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/119.0.0.0 Safari/537.36",
		"Mozilla/5.0 (Macintosh; Intel Mac OS X 13_0) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36",
		"Mozilla/5.0 (Macintosh; Intel Mac OS X 14_0) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/17.0 Safari/605.1.15",
		"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/16.6 Safari/605.1.15",
		"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36",
		"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/119.0.0.0 Safari/537.36",
	}

	// 生成随机referer
	referers := []string{
		"https://www.apple.com/shop/buy-iphone",
		"https://www.apple.com/shop/buy-iphone/iphone-17",
		"https://www.apple.com/shop/buy-iphone/iphone-16",
		"https://www.apple.com.cn/shop/buy-iphone",
		"https://www.apple.com.cn/shop/buy-iphone/iphone-17",
		"https://www.apple.com.cn/shop/buy-iphone/iphone-16",
		"https://www.apple.com/shop/",
		"https://www.apple.com.cn/shop/",
	}

	// 随机选择
	rand.Seed(time.Now().UnixNano())
	userAgent := userAgents[rand.Intn(len(userAgents))]
	referer := referers[rand.Intn(len(referers))]

	// 生成随机session ID
	sessionID := fmt.Sprintf("s_%d%d", time.Now().Unix(), rand.Intn(100000))

	// 生成随机cookie值
	cookieValues := []string{
		fmt.Sprintf("s_vi=[CS]v1|%X[CE]", rand.Int63()),
		fmt.Sprintf("s_fid=%X-%X", rand.Int63(), rand.Int63()),
		fmt.Sprintf("s_cc=true"),
		fmt.Sprintf("as_dc=nc"),
	}
	cookieString := strings.Join(cookieValues, "; ")

	// 创建新的请求客户端，禁用缓存
	req := gorequest.New().
		Set("referer", referer).
		Set("user-agent", userAgent).
		Set("accept", "application/json, text/javascript, */*; q=0.01").
		Set("accept-language", "zh-CN,zh;q=0.9,en;q=0.8").
		Set("accept-encoding", "gzip, deflate, br").
		Set("cache-control", "no-cache, no-store, must-revalidate").
		Set("pragma", "no-cache").
		Set("x-requested-with", "XMLHttpRequest").
		Set("sec-fetch-dest", "empty").
		Set("sec-fetch-mode", "cors").
		Set("sec-fetch-site", "same-origin").
		Set("sec-ch-ua", `"Not_A Brand";v="8", "Chromium";v="120", "Google Chrome";v="120"`).
		Set("sec-ch-ua-mobile", "?0").
		Set("sec-ch-ua-platform", `"macOS"`).
		Set("dnt", "1").
		Set("cookie", cookieString).
		Set("x-aos-model-page", "shop").
		Set("x-aos-stk", sessionID).
		Timeout(time.Second * 5)

	// 添加随机延迟（100-500ms）
	delay := time.Duration(100+rand.Intn(400)) * time.Millisecond
	time.Sleep(delay)

	resp, body, errs := req.Get(skUrl).End()
	if len(errs) > 0 {
		log.Println(errs)
		ch <- skus
		return
	}

	log.Println(resp.Status, skUrl)

	// 解析响应JSON
	// 中国地区的响应格式: body.stores[].partsAvailability
	// 其他地区的响应格式: body.content.pickupMessage.stores 或 body.content.pickupMessage.pickupEligibility
	// 库存判断字段: pickupDisplay ("available"=有货, "ineligible"=无货)

	// 先尝试中国地区格式 (body.stores)
	stores := gjson.Get(body, "body.stores").Array()
	if len(stores) > 0 {
		// 中国地区格式
		for _, store := range stores {
			storeNumber := store.Get("storeNumber").String()
			// 检查每个产品在该店铺的可用性
			partsAvailability := store.Get("partsAvailability").Map()
			for productCode, availability := range partsAvailability {
				uniqKey := fmt.Sprintf("%s.%s", storeNumber, productCode)
				// 检查 pickupDisplay 字段，"available" 表示有货，"ineligible" 表示无货
				pickupDisplay := availability.Get("pickupDisplay").String()
				isAvailable := (pickupDisplay == "available")

				// pickupDisplay 为 "available" 表示有货
				skus[uniqKey] = isAvailable

				if isAvailable {
					log.Printf("✅ Store %s Product %s: IN STOCK (pickupDisplay=%s)", storeNumber, productCode, pickupDisplay)
				} else {
					log.Printf("❌ Store %s Product %s: OUT OF STOCK (pickupDisplay=%s)", storeNumber, productCode, pickupDisplay)
				}
			}
		}
	} else {
		// 尝试其他地区格式 (body.content.pickupMessage.stores)
		stores = gjson.Get(body, "body.content.pickupMessage.stores").Array()
		if len(stores) > 0 {
			for _, result := range stores {
				storeNumber := result.Get("storeNumber").String()
				// 检查每个产品在该店铺的可用性
				for productCode, availability := range result.Get("partsAvailability").Map() {
					uniqKey := fmt.Sprintf("%s.%s", storeNumber, productCode)
					// 检查 pickupDisplay 字段，"available" 表示有货，"ineligible" 表示无货
					pickupDisplay := availability.Get("pickupDisplay").String()
					isAvailable := (pickupDisplay == "available")

					// pickupDisplay 为 "available" 表示有货
					skus[uniqKey] = isAvailable

					if isAvailable {
						log.Printf("✅ Store %s Product %s: IN STOCK (pickupDisplay=%s)", storeNumber, productCode, pickupDisplay)
					} else {
						log.Printf("❌ Store %s Product %s: OUT OF STOCK (pickupDisplay=%s)", storeNumber, productCode, pickupDisplay)
					}
				}
			}
		} else {
			// 单店铺查询逻辑（新格式）- 国外地区
			pickupEligibility := gjson.Get(body, "body.content.pickupMessage.pickupEligibility")
			if pickupEligibility.Exists() {
				// 从URL中提取店铺编号
				u, _ := url.Parse(skUrl)
				storeNumber := u.Query().Get("store")

				// 遍历所有产品
				pickupEligibility.ForEach(func(productCode, productInfo gjson.Result) bool {
					uniqKey := fmt.Sprintf("%s.%s", storeNumber, productCode)

					// 检查 pickupDisplay 字段，"available" 表示有货，"ineligible" 表示无货
					// 在这种格式下，pickupDisplay 可能在不同的路径
					pickupDisplay := productInfo.Get("messageTypes.regular.pickupDisplay").String()
					if pickupDisplay == "" {
						// 尝试其他可能的路径
						pickupDisplay = productInfo.Get("pickupDisplay").String()
					}
					if pickupDisplay == "" {
						// 再尝试另一个路径
						pickupDisplay = productInfo.Get("messageTypes.compact.pickupDisplay").String()
					}

					isAvailable := (pickupDisplay == "available")

					skus[uniqKey] = isAvailable

					if isAvailable {
						log.Printf("✅ Store %s Product %s: IN STOCK (pickupDisplay=%s)", storeNumber, productCode, pickupDisplay)
					} else {
						log.Printf("❌ Store %s Product %s: OUT OF STOCK (pickupDisplay=%s)", storeNumber, productCode, pickupDisplay)
					}

					return true
				})
			} else {
				// 如果没有找到任何已知格式，打印响应预览帮助调试
				log.Printf("Unknown response format, body preview: %s", body[:min(200, len(body))])
			}
		}
	}

	ch <- skus
}

// 型号对应预约地址
func (s *listenService) model2Url(productType string, partNumber string) string {
	// https://www.apple.com.cn/shop/buy-iphone/iphone-13/MLE73CH/A
	var t string
	switch productType {
	// iPhone 17系列
	case "iphone17promax", "iphone17pro":
		t = "iphone-17-pro"
	case "iphoneair":
		t = "iphone-air"
	case "iphone17":
		t = "iphone-17"
	// iPhone 16系列
	case "iphone16promax", "iphone16pro":
		t = "iphone-16-pro"
	case "iphone16plus":
		t = "iphone-16"
	case "iphone16":
		t = "iphone-16"
	// Apple Watch系列
	case "watchultra3":
		t = "apple-watch-ultra"
	case "watchs11":
		t = "apple-watch"
	case "watchse3":
		t = "apple-watch-se"
	}

	// 根据产品类型选择正确的URL路径
	urlPath := "buy-iphone"
	if strings.Contains(productType, "watch") {
		urlPath = "buy-watch"
	}

	return fmt.Sprintf(
		"https://www.apple.com/%s/shop/%s/%s/%s",
		s.Area.ShortCode,
		urlPath,
		t,
		partNumber,
	)
}

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
		log.Println("Bark URL is empty, skipping notification")
		return
	}

	// 对标题和内容进行URL编码
	encodedTitle := url.QueryEscape(title)
	encodedContent := url.QueryEscape(content)
	encodedUrl := url.QueryEscape(bagUrl)

	// 构建正确的Bark API URL格式
	// 格式: https://api.day.app/{device_key}?title={title}&body={body}&url={url}
	apiUrl := fmt.Sprintf("%s?title=%s&body=%s&url=%s",
		strings.TrimRight(s.BarkNotifyUrl, "/"),
		encodedTitle,
		encodedContent,
		encodedUrl)

	log.Printf("Sending Bark notification to: %s", apiUrl)

	response, err := http.Get(apiUrl)
	if err != nil {
		log.Printf("Bark notification error: %v", err)
		return
	}
	defer response.Body.Close()

	// 读取响应内容
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Printf("Failed to read Bark response: %v", err)
	}

	// 检查响应状态
	if response.StatusCode != 200 {
		log.Printf("Bark notification failed with status: %d, response: %s", response.StatusCode, string(body))
		return
	}

	log.Printf("Bark notification sent successfully, response: %s", string(body))
}
