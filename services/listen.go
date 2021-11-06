package services

import (
    "apple-store-helper/model"
    "apple-store-helper/theme"
    "apple-store-helper/view"
    "bytes"
    "errors"
    "fmt"
    "fyne.io/fyne/v2"
    "fyne.io/fyne/v2/data/binding"
    "fyne.io/fyne/v2/dialog"
    "fyne.io/fyne/v2/widget"
    "github.com/PuerkitoBio/goquery"
    "github.com/faiface/beep"
    "github.com/faiface/beep/mp3"
    "github.com/faiface/beep/speaker"
    "github.com/golang-module/carbon"
    "github.com/parnurzeal/gorequest"
    "github.com/tidwall/gjson"
    "io/ioutil"
    "log"
    "net/url"
    "strconv"
    "strings"
    "time"
)

const (
    StatusOutStock = "无货"
    StatusInStock  = "有货"
    StatusWait     = "等待"
    
    Pause   = "暂停"
    Running = "监听中"
)

var Listen = listenService{
    items:    map[string]ListenItem{},
    Status:   binding.NewString(),
    Area:     model.Areas[0],
    Logs:     widget.NewLabel(""),
}

type listenService struct {
    items    map[string]ListenItem
    Status   binding.String
    Area     model.Area
    Logs     *widget.Label
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
                skus, err := s.groupByStore()
                if err != nil {
                    dialog.ShowInformation("遇到错误", err.Error(), view.Window)
                    s.Status.Set(Pause)
                }
                
                for key, item := range s.items {
                    status := skus[item.Store.StoreNumber+"."+item.Product.Code]
                    
                    if status {
                        s.UpdateStatus(key, StatusInStock)
                        s.Status.Set(Pause)
                        
                        // 进入购物袋, 手动选择门店
                        s.openBrowser(fmt.Sprintf("https://www.apple.com/%s/shop/bag", s.Area.ShortCode))
                        msg := fmt.Sprintf("%s %s 有货", item.Store.CityStoreName, item.Product.Title)
                        dialog.ShowInformation("匹配成功", msg, view.Window)
                        view.App.SendNotification(&fyne.Notification{
                            Title:   "有货提醒",
                            Content: msg,
                        })
                        go s.AlertMp3()
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

func  (s *listenService) groupByStore() (map[string]bool,error) {
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
        return skus,nil
    }
    
    ch := make(chan map[string]bool, count)
    chErr := make(chan error)
    for _, link := range reqs {
        go s.getSkuByLink(ch,chErr, link)
    }
    
    for i := 0; i < count; i++ {
        select {
        case err := <-chErr:
            return skus, err
        case val := <-ch:
            for key, v := range val {
                skus[key] = v
            }
        }
    }
    
    return skus,nil
}

func (s *listenService) getSkuByLink(ch chan map[string]bool, chErr chan error,  skUrl string) {
    skus := map[string]bool{}
    
    resp, body, errs := gorequest.New().
        Set("referer", "https://www.apple.com/shop/buy-iphone").
        Set("user-agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/94.0.4606.71 Safari/537.36").
        Timeout(time.Second*3).Get(skUrl).End()
    if len(errs) > 0 {
        log.Println(errs)
        // chErr <- fmt.Errorf("请求网络遇到错误 %w", errs[0])
		// 提示信息太复杂，不方便显示，暂时不处理
		ch <- skus
        return
    }
    
    log.Println(resp.Status, skUrl)

    //开始尝试按照html解析，如果html解析成功，按照失败处理
    dom, err := goquery.NewDocumentFromReader(strings.NewReader(body))
    if err == nil {
        log.Println(body)
        msg := dom.Find("head").Find("title").Text()
        if len(msg) > 0 {
            chErr <- errors.New(msg)
            return
        }
    }
    for _, result := range gjson.Get(body, "body.content.pickupMessage.stores").Array() {
        for productCode, availability := range result.Get("partsAvailability").Map() {
            uniqKey := fmt.Sprintf("%s.%s", result.Get("storeNumber").String(), productCode)
            skus[uniqKey] = availability.Get("storeSelectionEnabled").Bool()
        }
    }
    
    ch <- skus
}

// 型号对应预约地址
func (s *listenService) model2Url(productType string, partNumber string) string {
    // https://www.apple.com.cn/shop/buy-iphone/iphone-13/MLE73CH/A
    var t string
    switch productType {
    case "iphone13promax", "iphone13pro":
        t = "iphone-13-pro"
    case "iphone13mini", "iphone13":
        t = "iphone-13"
    }
    
    return fmt.Sprintf(
        "https://www.apple.com/%s/shop/buy-iphone/%s/%s",
        s.Area.ShortCode,
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
