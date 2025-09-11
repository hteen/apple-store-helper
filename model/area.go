package model

type Area struct {
	Title     string
	Locale    string
	ShortCode string
}

/**
 * Areas 地区
 * 中国大陆: CN/zh_CN
 * 中国澳门: MO/zh_MO
 *
 * ProductsJson 是一个数组，第一个值是基础款的 Data, 第二个值是 Pro 款的 Data
 * 打开购买页面，在开发者工具 Elements 面板中搜索 productSelectionData 提取值
 * Example URLs:
 * - [iPhone 15](https://www.apple.com.cn/shop/buy-iphone/iphone-15)
 * - [iPhone 15 Pro](https://www.apple.com.cn/shop/buy-iphone/iphone-15-pro)
 */
var Areas = []Area{
	{
		Title:     "中国大陆",
		Locale:    "zh_CN",
		ShortCode: "cn",
	},
	{
		Title:     "中国香港",
		ShortCode: "hk-zh", // hk
		Locale:    "zh_HK",
	},
	{
		Title:     "中国台湾",
		ShortCode: "tw",
		Locale:    "zh_TW",
	},
	{
		Title:     "Singapore",
		ShortCode: "sg",
		Locale:    "en_SG",
	},
	{
		Title:     "日本",
		ShortCode: "jp",
		Locale:    "ja_JP",
	},
	{
		Title:     "Australia",
		ShortCode: "au",
		Locale:    "en_AU",
	},
	{
		Title:     "Malaysia",
		ShortCode: "my",
		Locale:    "en_MY",
	},
}
