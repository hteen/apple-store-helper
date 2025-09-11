package services

import (
	"apple-store-helper/model"
	"log"
	"sort"

	"github.com/thoas/go-funk"
)

var Store = storeService{
	stores: map[string][]model.Store{},
}

type storeService struct {
	stores map[string][]model.Store
}

func (s *storeService) ByArea(area model.Area) []model.Store {
	// 检查内存缓存
	if len(s.stores[area.Locale]) > 0 {
		return s.stores[area.Locale]
	}

	// 从 store_data 目录加载门店数据
	storeData, err := LoadStoreData(area.ShortCode)
	if err == nil && len(storeData.Stores) > 0 {
		s.stores[area.Locale] = storeData.Stores
		return storeData.Stores
	}

	log.Printf("Failed to load store data for %s: %v", area.ShortCode, err)

	// 回退到静态门店列表
	if stores, exists := model.GlobalStores[area.Locale]; exists {
		s.stores[area.Locale] = stores
		return stores
	}

	// 如果没有匹配的地区，返回空列表
	return []model.Store{}
}

func (s *storeService) ByAreaTitleForOptions(areaTitle string) []string {
	area := Area.GetArea(areaTitle)
	areas := funk.Get(s.ByArea(area), "CityStoreName").([]string)
	sort.Strings(areas)
	return areas
}

// ByAreaAndProvinceForOptions 根据地区和省份筛选门店选项
func (s *storeService) ByAreaAndProvinceForOptions(areaTitle string, province string) []string {
	area := Area.GetArea(areaTitle)
	stores := s.ByArea(area)

	// 如果省份为空，返回所有门店
	if province == "" {
		areas := funk.Get(stores, "CityStoreName").([]string)
		sort.Strings(areas)
		return areas
	}

	// 根据省份筛选门店
	var filteredStores []model.Store
	for _, store := range stores {
		// 检查门店是否属于该省份
		if s.isStoreInProvince(store, province) {
			filteredStores = append(filteredStores, store)
		}
	}

	areas := funk.Get(filteredStores, "CityStoreName").([]string)
	sort.Strings(areas)
	return areas
}

// isStoreInProvince 检查门店是否属于指定省份
func (s *storeService) isStoreInProvince(store model.Store, province string) bool {
	// 根据门店的Province字段判断是否属于该省份
	return store.Province == province
}

// GetStatesForArea 获取指定地区的州/省份列表
func (s *storeService) GetStatesForArea(areaTitle string) []string {
	area := Area.GetArea(areaTitle)
	stores := s.ByArea(area)

	// 获取所有唯一的州/省份
	stateMap := make(map[string]bool)
	for _, store := range stores {
		if store.Province != "" {
			stateMap[store.Province] = true
		}
	}

	// 转换为切片并排序
	var states []string
	for state := range stateMap {
		states = append(states, state)
	}
	sort.Strings(states)

	return states
}

func (s *storeService) GetStore(areaTitle string, storeTitle string) model.Store {
	// 确保门店数据已加载
	area := Area.GetArea(areaTitle)
	stores := s.ByArea(area)

	result := funk.Find(stores, func(x model.Store) bool {
		return x.CityStoreName == storeTitle
	})

	// 检查是否找到门店，如果没找到返回空门店
	if result == nil {
		log.Printf("Warning: Store not found for area '%s', store '%s'", areaTitle, storeTitle)
		return model.Store{}
	}

	return result.(model.Store)
}

// LoadForArea 加载指定地区的门店数据
func (s *storeService) LoadForArea(areaCode string) error {
	storeData, err := LoadStoreData(areaCode)
	if err != nil {
		return err
	}

	// 根据areaCode找到对应的locale
	var locale string
	for _, area := range model.Areas {
		if area.ShortCode == areaCode {
			locale = area.Locale
			break
		}
	}

	if locale != "" {
		s.stores[locale] = storeData.Stores
		log.Printf("Loaded %d stores for %s", len(storeData.Stores), areaCode)
	}

	return nil
}
