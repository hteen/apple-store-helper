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
}

type storeService struct {
	stores map[string][]model.Store
}

func (s *storeService) ByAreaCode(areaCode string) []model.Store {

	if len(s.stores[areaCode]) > 0 {
		return s.stores[areaCode]
	}

	availability := "https://reserve-prime.apple.com/" + areaCode + "/reserve/A/stores.json"
	_, bd, errs := gorequest.New().Get(availability).End()

	if len(errs) != 0 {
		panic(errs[0])
	}

	for _, store := range gjson.Get(bd, "stores").Array() {
		s.stores[areaCode] = append(s.stores[areaCode], model.Store{
			StoreNumber:   store.Get("storeNumber").String(),
			CityStoreName: store.Get("city").String() + "-" + store.Get("storeName").String(),
		})
	}

	return s.stores[areaCode]
}

func (s *storeService) ByAreaTitleForOptions(areaTitle string) []string {
	code := Area.Title2Code(areaTitle)
	areas := funk.Get(s.ByAreaCode(code), "CityStoreName").([]string)
	sort.Strings(areas)
	return areas
}

func (s *storeService) GetStore(areaTitle string, storeTitle string) model.Store {
	code := Area.Title2Code(areaTitle)

	return funk.Find(s.stores[code], func(x model.Store) bool {
		return x.CityStoreName == storeTitle
	}).(model.Store)
}
