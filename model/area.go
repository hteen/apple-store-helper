package model

type Area struct {
	Title        string
	Locale       string
	ShortCode    string
	ProductsJson []string
}

// Areas 地区列表 - 中国大陆、香港、日本、新加坡、美国、英国和澳大利亚
var Areas = []Area{
	{Title: "中国大陆", Locale: "zh_CN", ShortCode: "cn", ProductsJson: iPhone17ProductsJson},
	{Title: "香港", Locale: "zh_HK", ShortCode: "hk", ProductsJson: iPhone17ProductsJson},
	{Title: "日本", Locale: "ja_JP", ShortCode: "jp", ProductsJson: iPhone17ProductsJson},
	{Title: "新加坡", Locale: "en_SG", ShortCode: "sg", ProductsJson: iPhone17ProductsJson},
	{Title: "美国", Locale: "en_US", ShortCode: "us", ProductsJson: iPhone17ProductsJson},
	{Title: "英国", Locale: "en_GB", ShortCode: "uk", ProductsJson: iPhone17ProductsJson},
	{Title: "澳大利亚", Locale: "en_AU", ShortCode: "au", ProductsJson: iPhone17ProductsJson},
}
