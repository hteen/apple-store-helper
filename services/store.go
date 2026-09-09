package services

import (
	"apple-store-helper/config"
	"apple-store-helper/model"
	"fmt"
	"sort"
	"strings"

	"github.com/thoas/go-funk"
	"github.com/tidwall/gjson"
)

var Store = storeService{
	stores: map[string][]model.Store{},
}

type storeService struct {
	stores map[string][]model.Store
}

func (s *storeService) ByArea(area model.Area) []model.Store {
	stores, err := config.ReadConfigFile("stores.json")
	if err != nil {
		panic(err)
	}

	for _, v := range gjson.ParseBytes(stores).Array() {
		locale := v.Get("locale").String()
		hasStates := v.Get("hasStates").Bool()

		localeStores := []model.Store{}

		if hasStates {
			for _, state := range v.Get("state").Array() {
				for _, store := range state.Get("store").Array() {
					stateName := store.Get("address.stateName").String()
					city := store.Get("address.city").String()
					localeStores = append(localeStores, model.Store{
						StoreNumber:   store.Get("id").String(),
						CityStoreName: fmt.Sprintf("%s-%s", stateName, store.Get("name").String()),
						Location:      strings.TrimSpace(stateName + " " + city),
					})
				}
			}
		} else {
			for _, store := range v.Get("store").Array() {
				city := store.Get("address.city").String()
				localeStores = append(localeStores, model.Store{
					StoreNumber:   store.Get("id").String(),
					CityStoreName: fmt.Sprintf("%s-%s", city, store.Get("name").String()),
					Location:      city,
				})
			}
		}

		// 去重
		localeStores = funk.UniqBy(localeStores, func(x model.Store) string {
			return x.StoreNumber
		}).([]model.Store)

		s.stores[locale] = localeStores
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
	s.ensureLoaded(code)

	return funk.Find(s.stores[code], func(x model.Store) bool {
		return x.CityStoreName == storeTitle
	}).(model.Store)
}

func (s *storeService) GetByNumber(locale string, storeNumber string) model.Store {
	s.ensureLoaded(locale)
	found := funk.Find(s.stores[locale], func(x model.Store) bool {
		return x.StoreNumber == storeNumber
	})
	if found == nil {
		return model.Store{StoreNumber: storeNumber}
	}
	return found.(model.Store)
}

func (s *storeService) ensureLoaded(locale string) {
	if len(s.stores[locale]) > 0 {
		return
	}
	for _, area := range model.Areas {
		if area.Locale == locale {
			s.ByArea(area)
			return
		}
	}
}
