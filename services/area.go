package services

import (
	"apple-store-helper/model"
	"github.com/thoas/go-funk"
	"github.com/tidwall/gjson"
)

var Area = areaService{}

type areaService struct{}

func (s *areaService) ProductsByCode(code string) []model.Product {

	area := funk.Find(model.Areas, func(x model.Area) bool {
		return x.Code == code
	}).(model.Area)

	var products []model.Product

	for _, result := range gjson.Parse(area.ProductsJson).Array() {
		products = append(products, model.Product{
			Title: result.Get("familyType").String() + "-" + result.Get("seoUrlToken").String(),
			Type:  result.Get("familyType").String(),
			Code:  result.Get("partNumber").String(),
		})
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

	return area.Code
}

func (s *areaService) GetArea(title string) model.Area {
	area := funk.Find(model.Areas, func(x model.Area) bool {
		return x.Title == title
	}).(model.Area)

	return area
}
