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
 * ProductsJson 是一个数组，每个元素对应一个购买页的 productSelectionData
 * 打开购买页面，在开发者工具 Elements 面板中搜索 productSelectionData 提取值
 * Example URLs:
 * - [iPhone 17](https://www.apple.com.cn/shop/buy-iphone/iphone-17)
 * - [iPhone 17 Pro](https://www.apple.com.cn/shop/buy-iphone/iphone-17-pro)
 * - [iPhone Air](https://www.apple.com.cn/shop/buy-iphone/iphone-air)
 * - [iPhone Duo](https://www.apple.com.cn/shop/buy-iphone/iphone-duo)
 * - [iPhone 18 Pro](https://www.apple.com.cn/shop/buy-iphone/iphone-18-pro)
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

// ShopOrigin 返回对应地区 Apple Store 的站点根地址。
// 中国大陆必须使用 apple.com.cn，apple.com/cn 的库存接口会被拦截或返回空数据。
func (a Area) ShopOrigin() string {
	if a.ShortCode == "cn" {
		return "https://www.apple.com.cn"
	}
	return "https://www.apple.com/" + a.ShortCode
}

func (a Area) AcceptLanguage() string {
	switch a.Locale {
	case "zh_CN":
		return "zh-CN,zh;q=0.9"
	case "zh_HK":
		return "zh-HK,zh;q=0.9,en;q=0.8"
	case "zh_TW":
		return "zh-TW,zh;q=0.9,en;q=0.8"
	case "ja_JP":
		return "ja-JP,ja;q=0.9,en;q=0.8"
	default:
		return "en-US,en;q=0.9"
	}
}
