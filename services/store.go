package services

import (
	"apple-store-helper/model"
	"github.com/parnurzeal/gorequest"
	"github.com/thoas/go-funk"
	"github.com/tidwall/gjson"
	"sort"
)

var Store = storeService{
	stores: map[string][]model.Store{},
	storeListData: "",
}

type storeService struct {
	stores map[string][]model.Store
	storeListData string
}

func (s *storeService) ByArea(area model.Area) []model.Store {

	if len(s.stores[area.Locale]) > 0 {
		return s.stores[area.Locale]
	}
	
	if s.storeListData == "" {
		availability := "https://www.apple.com/rsp-web/store-list?locale=zh_CN"
		_, bd, errs := gorequest.New().Get(availability).End()
		
		if len(errs) != 0 {
			panic(errs[0])
		}
		
		s.storeListData = bd
	}

	for _, list := range gjson.Get(s.storeListData, "storeListData").Array() {
		if list.Get("locale").String() == area.Locale {
			if list.Get("state").Exists() {
				for _, state := range list.Get("state").Array() {
					for _, store := range state.Get("store").Array() {
						s.stores[area.Locale] = append(s.stores[area.Locale], model.Store{
							StoreNumber:   store.Get("id").String(),
							CityStoreName: store.Get("address.stateName").String() + "-" + store.Get("name").String(),
						})
					}
				}
			} else {
				for _, store := range list.Get("store").Array() {
					s.stores[area.Locale] = append(s.stores[area.Locale], model.Store{
						StoreNumber:   store.Get("id").String(),
						CityStoreName: store.Get("address.city").String() + "-" + store.Get("name").String(),
					})
				}
			}
		}
	}

	return s.stores[area.Locale]
}

func (s *storeService) ByAreaTitleForOptions(areaTitle string) []string {
	area := Area.GetArea(areaTitle)
	areas := funk.Get(s.ByArea(area), "CityStoreName").([]string)
	sort.Strings(areas)
	return areas
}

func (s *storeService) GetStore(areaTitle string, storeTitle string) model.Store {
	code := Area.Title2Code(areaTitle)

	return funk.Find(s.stores[code], func(x model.Store) bool {
		return x.CityStoreName == storeTitle
	}).(model.Store)
}
