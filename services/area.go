package services

import (
	"apple-store-helper/model"
	"github.com/thoas/go-funk"
)

var Area = areaService{}

type areaService struct{}

// ProductsByCode - Deprecated: Use Product.GetDynamicProducts() instead
func (s *areaService) ProductsByCode(local string) []model.Product {
	// Return empty slice as we're using dynamic data only
	return []model.Product{}
}

func (s *areaService) ForOptions() []string {
	return funk.Get(model.Areas, "Title").([]string)
}

func (s *areaService) Title2Code(title string) string {
	area := funk.Find(model.Areas, func(x model.Area) bool {
		return x.Title == title
	}).(model.Area)

	return area.Locale
}

func (s *areaService) GetArea(title string) model.Area {
	area := funk.Find(model.Areas, func(x model.Area) bool {
		return x.Title == title
	}).(model.Area)

	return area
}
