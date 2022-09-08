package model

import (
	"embed"
	"fmt"
	"github.com/tidwall/gjson"
)

type Area struct {
	Title     string
	Locale    string
	ShortCode string
	Color     map[string]gjson.Result
	Products  []gjson.Result
}

//go:embed config/*
var config embed.FS

func readArray(area string, path string) []gjson.Result {
	content, err := config.ReadFile("config/" + area + "/iphone14.json")
	if err != nil {
		fmt.Printf("文件打开失败 [Err:%s]\n", err.Error())
		return nil
	}
	con := gjson.Parse(string(content)).Get(path).Array()
	//products := gjson.Get(string(content), "products").Array()
	content, err = config.ReadFile("config/" + area + "/iphone14pro.json")
	if err != nil {
		fmt.Printf("文件打开失败 [Err:%s]\n", err.Error())
		return nil
	}
	con2 := gjson.Parse(string(content)).Get(path).Array()
	//products2 := gjson.Get(string(content), "products").Array()
	return append(con, con2...)
}

func readMap(area string, path string) map[string]gjson.Result {
	content, err := config.ReadFile("config/" + area + "/iphone14.json")
	if err != nil {
		fmt.Printf("文件打开失败 [Err:%s]\n", err.Error())
		return nil
	}
	con := gjson.Parse(string(content)).Get(path).Map()
	//products := gjson.Get(string(content), "products").Array()
	content, err = config.ReadFile("config/" + area + "/iphone14pro.json")
	if err != nil {
		fmt.Printf("文件打开失败 [Err:%s]\n", err.Error())
		return nil
	}
	con2 := gjson.Parse(string(content)).Get(path).Map()
	//products2 := gjson.Get(string(content), "products").Array()
	m := map[string]gjson.Result{}
	for s, result := range con {
		m[s] = result
	}
	for s, result := range con2 {
		m[s] = result
	}
	return m
}

// Areas 地区，中国大陆: CN/zh_CN, 中国澳门: MO/zh_MO
var Areas = []Area{
	{
		Title:     "中国大陆",
		Locale:    "zh_CN",
		ShortCode: "cn",
		// ProductsJson view-source:https://www.apple.com.cn/shop/buy-iphone/iphone-13 productSelectionData.products
		Color:    readMap("zh-cn", "displayValues.dimensionColor"),
		Products: readArray("zh-cn", "products"),
	},
	{
		Title:     "中国香港",
		ShortCode: "hk", // hk-zh
		Locale:    "zh_HK",
		Color:     readMap("zh-hk", "displayValues.dimensionColor"),
		Products:  readArray("zh-hk", "products"),
	},
	{
		Title:     "中国台湾",
		ShortCode: "tw",
		Locale:    "zh_TW",
		Color:     readMap("zh-tw", "displayValues.dimensionColor"),
		Products:  readArray("zh-tw", "products"),
	},
	{
		Title:     "日本",
		ShortCode: "jp",
		Locale:    "ja_JP",
		Color:     readMap("jp", "displayValues.dimensionColor"),
		Products:  readArray("jp", "products"),
	},
}
