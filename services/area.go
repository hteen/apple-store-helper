package services

import (
	"apple-store-helper/model"
	"fmt"
	"github.com/thoas/go-funk"
	"github.com/tidwall/gjson"
)

var Area = areaService{}

type areaService struct{}

func (s *areaService) ProductsByCode(local string) []model.Product {

	area := funk.Find(model.Areas, func(x model.Area) bool {
		return x.Locale == local
	}).(model.Area)

	var products []model.Product

	for _, pJson := range area.ProductsJson {
		json := gjson.Parse(pJson)
		for _, result := range json.Get("products").Array() {
			color := json.Get(fmt.Sprintf("displayValues.dimensionColor.%s.value", result.Get("dimensionColor")))
			products = append(products, model.Product{
				Title: fmt.Sprintf("%s - %s - %s", result.Get("familyType"), color, result.Get("dimensionCapacity")),
				Type:  result.Get("familyType").String(),
				Code:  result.Get("partNumber").String(),
			})
		}
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
