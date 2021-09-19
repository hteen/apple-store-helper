package services

import (
    "apple-store-helper/model"
    "github.com/thoas/go-funk"
)

var Product = productService{}
type productService struct{}

func (s *productService) ByAreaTitleForOptions(areaTitle string) []string {
    code := Area.Title2Code(areaTitle)
    return funk.Get(Area.ProductsByCode(code), "Title").([]string)
}

func (s *productService) GetProduct(areaTitle string, productTitle string) model.Product {
    code := Area.Title2Code(areaTitle)
    
    return funk.Find(Area.ProductsByCode(code), func(x model.Product) bool {
        return x.Title == productTitle
    }).(model.Product)
}