package services

import (
    "apple-store-helper/model"
    
    "github.com/thoas/go-funk"
)

var Area = areaService{}

type areaService struct{}

func (s *areaService) ForOptions() []string {
	// 仅返回内置地区
	return funk.Get(model.Areas, "Title").([]string)
}

func (s *areaService) Title2Code(title string) string {
	areaInterface := funk.Find(model.Areas, func(x model.Area) bool {
		return x.Title == title
	})
	if areaInterface == nil {
		return ""
	}
	area := areaInterface.(model.Area)

	return area.Locale
}

// TitleByCode 将 locale 转换为地区 Title（找不到则返回原 locale）
func (s *areaService) TitleByCode(locale string) string {
    for _, area := range model.Areas {
        if area.Locale == locale {
            return area.Title
        }
    }
    return locale
}

func (s *areaService) GetArea(title string) model.Area {
	areaInterface := funk.Find(model.Areas, func(x model.Area) bool {
		return x.Title == title
	})
	if areaInterface == nil {
		// 返回默认地区或空地区
		if len(model.Areas) > 0 {
			return model.Areas[0]
		}
		return model.Area{}
	}
	area := areaInterface.(model.Area)

	return area
}
