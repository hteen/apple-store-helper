package services

import (
    "apple-store-helper/config"
    "apple-store-helper/model"
    "fmt"
    "os"
    
    "github.com/thoas/go-funk"
    "github.com/tidwall/gjson"
)

var Product = productService{}
type productService struct{}

// ProductsByCode 根据地区代码获取产品列表（合并内置和缓存数据）
func (s *productService) ProductsByCode(locale string) []model.Product {
    areaInterface := funk.Find(model.Areas, func(x model.Area) bool {
        return x.Locale == locale
    })
    if areaInterface == nil {
        return []model.Product{}
    }
    area := areaInterface.(model.Area)

    var products []model.Product
    productMap := make(map[string]model.Product) // 用于去重
    
    // 先加载内置数据
    filename := fmt.Sprintf("products_%s.json", area.Locale)
    if builtinData, err := config.ReadConfigFile(filename); err == nil {
        builtinProducts := gjson.ParseBytes(builtinData)
        for _, json := range builtinProducts.Array() {
            for _, result := range json.Get("products").Array() {
                color := json.Get(fmt.Sprintf("displayValues.dimensionColor.%s.value", result.Get("dimensionColor")))
                product := model.Product{
                    Title: fmt.Sprintf("%s - %s - %s", result.Get("familyType"), color, result.Get("dimensionCapacity")),
                    Type:  result.Get("familyType").String(),
                    Code:  result.Get("partNumber").String(),
                }
                productMap[product.Code] = product
            }
        }
    }
    
    // 再加载缓存数据，合并到productMap中（自动去重）
    if cachedData, err := os.ReadFile("user_config/" + filename); err == nil {
        cachedProducts := gjson.ParseBytes(cachedData)
        for _, json := range cachedProducts.Array() {
            for _, result := range json.Get("products").Array() {
                color := json.Get(fmt.Sprintf("displayValues.dimensionColor.%s.value", result.Get("dimensionColor")))
                product := model.Product{
                    Title: fmt.Sprintf("%s - %s - %s", result.Get("familyType"), color, result.Get("dimensionCapacity")),
                    Type:  result.Get("familyType").String(),
                    Code:  result.Get("partNumber").String(),
                }
                productMap[product.Code] = product // 覆盖或添加
            }
        }
    }
    
    // 转换map为slice
    for _, product := range productMap {
        products = append(products, product)
    }

    return products
}

func (s *productService) ByAreaTitleForOptions(areaTitle string) []string {
    code := Area.Title2Code(areaTitle)
    return funk.Get(s.ProductsByCode(code), "Title").([]string)
}

func (s *productService) GetProduct(areaTitle string, productTitle string) model.Product {
    code := Area.Title2Code(areaTitle)
    
    return funk.Find(s.ProductsByCode(code), func(x model.Product) bool {
        return x.Title == productTitle
    }).(model.Product)
}