package services

import (
	"apple-store-helper/model"
	"errors"
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
	"github.com/golang-module/carbon"
	"github.com/parnurzeal/gorequest"
	"github.com/thoas/go-funk"
	"github.com/tidwall/gjson"
	"log"
	"os/exec"
	"runtime"
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
	Quantity: widget.NewSelect([]string{"1", "2"}, nil),
}

type listenService struct {
	items    map[string]ListenItem
	Status   binding.String
	Area     model.Area
	Window   fyne.Window
	Logs     *widget.Label
	Quantity *widget.Select
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
				skus := s.getSkus()

				for key, item := range s.items {
					status := skus[item.Store.StoreNumber+"."+item.Product.Code]

					if status {
						s.UpdateStatus(key, StatusInStock)
						s.openBrowser(s.model2Url(item.Product.Type, item.Store.StoreNumber, item.Product.Code))
						dialog.ShowInformation("匹配成功", fmt.Sprintf("%s %s 有货", item.Store.CityStoreName, item.Product.Title), s.Window)
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

func (s *listenService) getSkuByCode(code string) gjson.Result {
	skUrl := fmt.Sprintf(
		"https://reserve-prime.apple.com/%s/reserve/%s/availability.json",
		s.Area.Code,
		code,
	)

	_, body, _ := gorequest.New().Get(skUrl).End()

	return gjson.Get(body, "stores")
}

func (s *listenService) getSkus() map[string]bool {
	skus := map[string]bool{}
	for _, code := range funk.UniqString(funk.Values(model.TypeCode).([]string)) {
		stores := s.getSkuByCode(code)
		for storeCode, result := range stores.Map() {
			for productCode, availability := range result.Map() {
				inStock := availability.Get("contract").Bool() && availability.Get("unlocked").Bool()
				skus[storeCode+"."+productCode] = inStock
			}
		}
	}

	return skus
}

// RegisterCode 帮助提前获取注册码
func (s *listenService) RegisterCode(productType string) {
	stores := s.getSkuByCode(model.TypeCode[productType])

	// 寻找任意一个有货门店
	for store, items := range stores.Map() {
		for k, v := range items.Map() {
			if v.Get("availability.contract").Bool() && v.Get("availability.unlocked").Bool() {
				s.openBrowser(s.model2Url(productType, store, k))
				return
			}
		}
	}

	dialog.ShowError(errors.New("所有门店无货，无法前往注册码页面"), s.Window)
}

// 型号对应预约地址
func (s *listenService) model2Url(productType string, store string, partNumber string) string {
	return fmt.Sprintf(
		"https://reserve-prime.apple.com/%s/reserve/%s?quantity=%s&anchor-store=%s&store=%s&partNumber=%s&plan=unlocked",
		s.Area.Code,
		model.TypeCode[productType],
		s.Quantity.Selected,
		store,
		store,
		partNumber,
	)
}

func (s *listenService) openBrowser(url string) {
	var err error
	switch runtime.GOOS {
	case "linux":
		err = exec.Command("xdg-open", url).Start()
	case "windows":
		err = exec.Command("rundll32", "url.dll,FileProtocolHandler", url).Start()
	case "darwin":
		err = exec.Command("open", url).Start()
	default:
		err = fmt.Errorf("unsupported platform")
	}
	if err != nil {
		dialog.ShowError(errors.New("打开网页失败，请自行手动操作\n"+url), s.Window)
	}
}
