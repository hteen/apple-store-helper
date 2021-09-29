package services

import (
    "apple-store-helper/model"
    "fmt"
    "github.com/parnurzeal/gorequest"
    "github.com/thoas/go-funk"
    "github.com/tidwall/gjson"
    "log"
    "net/url"
    "strconv"
)

var Product = productService{}

type productService struct{}

func (s *productService) ByAreaTitleForOptions(areaTitle string) []string {
    area := Area.GetArea(areaTitle)
    return funk.Get(Area.ProductsByCode(area.Locale), "Title").([]string)
}

func (s *productService) GetProduct(areaTitle string, productTitle string) model.Product {
    area := Area.GetArea(areaTitle)
    
    return funk.Find(Area.ProductsByCode(area.Locale), func(x model.Product) bool {
        return x.Title == productTitle
    }).(model.Product)
}

func (s *productService) GetProductReserve(area string, location string, stores []string, parts []string) {
    // "https://www.apple.com/cn/shop/fulfillment-messages?pl=true&parts.0=MLDR3CH/A&location=浙江 杭州&parts.1=MLDG3CH/A"
    var uri url.URL
    q := uri.Query()
    
    q.Set("pl", "true")
    q.Set("location", location)
    for i, part := range parts {
        q.Set("parts."+strconv.FormatInt(int64(i), 10), part)
    }
    queryStr := q.Encode()
    
    link := fmt.Sprintf(
        "https://www.apple.com/%s/shop/fulfillment-messages?%s",
        area,
        queryStr,
    )
    
    log.Println(link)
    _, bd, errs := gorequest.New().
        Proxy("http://127.0.0.1:1087").
        Get(link).End()
    if len(errs) != 0 {
        panic(errs[0])
    }
    
    for _, result := range gjson.Get(bd, "body.content.pickupMessage.stores").Array() {
        for _, r := range result.Get("partsAvailability").Map() {
            if funk.ContainsString(stores, result.Get("storeNumber").String()) {
                log.Println(result.Get("storeNumber"), result.Get("storeName"), r.Get("storePickupProductTitle"), r.Get("pickupSearchQuote"))
            }
        }
    }
}
