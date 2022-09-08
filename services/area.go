package services

import (
	"apple-store-helper/model"
	"github.com/thoas/go-funk"
)

var Area = areaService{}

type areaService struct{}

func (s *areaService) ProductsByCode(local string) []model.Product {

	area := funk.Find(model.Areas, func(x model.Area) bool {
		return x.Locale == local
	}).(model.Area)

	var products []model.Product

	for _, result := range area.Products {
		products = append(products, model.Product{
			Color: result.Get("dimensionColor").String(),
			Title: result.Get("familyType").String(),
			Type:  result.Get("familyType").String(),
			Code:  result.Get("partNumber").String(),
		})
	}
	for i := range products {
		products[i].Title = products[i].Title + "-" + area.Color[products[i].Color].Get("value").String()
	}
	return products
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
