package embedded

import (
	_ "embed"
)

// 嵌入产品数据文件
//
//go:embed data/product/product_data_cn.json
var ProductDataCN []byte

//go:embed data/product/product_data_hk.json
var ProductDataHK []byte

//go:embed data/product/product_data_jp.json
var ProductDataJP []byte

//go:embed data/product/product_data_sg.json
var ProductDataSG []byte

//go:embed data/product/product_data_us.json
var ProductDataUS []byte

//go:embed data/product/product_data_uk.json
var ProductDataUK []byte

//go:embed data/product/product_data_au.json
var ProductDataAU []byte

// 嵌入门店数据文件（只嵌入存在的文件）
//
//go:embed data/store/store_cn.json
var StoreDataCN []byte

//go:embed data/store/store_hk.json
var StoreDataHK []byte

//go:embed data/store/store_jp.json
var StoreDataJP []byte

//go:embed data/store/store_us.json
var StoreDataUS []byte

//go:embed data/store/store_uk.json
var StoreDataUK []byte

//go:embed data/store/store_au.json
var StoreDataAU []byte

// 嵌入配置文件
// 注意：apple_urls.json 文件较大，暂时不嵌入

// 获取产品数据的映射
var ProductDataMap = map[string][]byte{
	"cn": ProductDataCN,
	"hk": ProductDataHK,
	"jp": ProductDataJP,
	"sg": ProductDataSG,
	"us": ProductDataUS,
	"uk": ProductDataUK,
	"au": ProductDataAU,
}

// 获取门店数据的映射（只包含存在的文件）
var StoreDataMap = map[string][]byte{
	"cn": StoreDataCN,
	"hk": StoreDataHK,
	"jp": StoreDataJP,
	"us": StoreDataUS,
	"uk": StoreDataUK,
	"au": StoreDataAU,
}

// GetProductData 获取指定地区的产品数据
func GetProductData(areaCode string) ([]byte, bool) {
	data, exists := ProductDataMap[areaCode]
	return data, exists
}

// GetStoreData 获取指定地区的门店数据
func GetStoreData(areaCode string) ([]byte, bool) {
	data, exists := StoreDataMap[areaCode]
	return data, exists
}

// GetAppleURLs 获取 Apple URLs 配置
func GetAppleURLs() []byte {
	// 暂时返回空，因为文件较大未嵌入
	return []byte{}
}
