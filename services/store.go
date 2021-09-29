package services

import (
    "apple-store-helper/model"
    "github.com/parnurzeal/gorequest"
    "github.com/thoas/go-funk"
    "github.com/tidwall/gjson"
    "log"
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
    area := Area.GetArea(areaTitle)
    areas := funk.Get(s.ByLocale(area.Locale), "CityStoreName").([]string)
    sort.Strings(areas)
    return areas
}

func (s *storeService) GetStore(areaTitle string, storeTitle string) model.Store {
    area := Area.GetArea(areaTitle)
    
    return funk.Find(s.stores[area.Locale], func(x model.Store) bool {
        return x.CityStoreName == storeTitle
    }).(model.Store)
}

func (s *storeService) ByLocale(locale string) []model.Store {
    if len(s.stores[locale]) > 0 {
        return s.stores[locale]
    }
    
    link := "https://www.apple.com/rsp-web/store-list?locale=" + locale
    log.Println(link)
    _, bd, errs := gorequest.New().
        Proxy("http://127.0.0.1:1087").
        Get(link).End()
    if len(errs) != 0 {
        panic(errs[0])
    }

    for _, state := range gjson.Get(bd, "storeListData.0.state").Array() {
        for _, store := range state.Get("store").Array() {
            s.stores[locale] = append(s.stores[locale], model.Store{
				StoreNumber:   store.Get("id").String(),
				CityStoreName: store.Get("address.city").String() + " " + store.Get("name").String(),
                Location: store.Get("address.stateName").String() + " " + store.Get("address.city").String(),
			})
        }
    }
	
	for _, store := range s.stores[locale] {
		log.Println(store.StoreNumber, store.CityStoreName)
	}

    return s.stores[locale]
}
