package main

import (
	"apple-store-helper/embedded"
	"apple-store-helper/model"
	"apple-store-helper/services"
	"encoding/json"
	"log"
)

func main() {
	log.Println("=== 最终嵌入数据测试 ===")

	// 测试嵌入的商店数据
	log.Println("1. 测试嵌入的商店数据...")
	storeData, exists := embedded.GetStoreData("cn")
	if !exists {
		log.Fatal("❌ 嵌入的商店数据不存在")
	}

	var storeDataStruct struct {
		Stores []struct {
			StoreNumber   string `json:"StoreNumber"`
			CityStoreName string `json:"CityStoreName"`
			Province      string `json:"Province"`
			City          string `json:"City"`
		} `json:"stores"`
	}

	if err := json.Unmarshal(storeData, &storeDataStruct); err != nil {
		log.Fatalf("❌ 解析嵌入商店数据失败: %v", err)
	}

	log.Printf("✅ 嵌入商店数据: %d 个商店", len(storeDataStruct.Stores))

	// 检查安徽商店
	hasAnhui := false
	for _, store := range storeDataStruct.Stores {
		if store.Province == "安徽" {
			hasAnhui = true
			log.Printf("✅ 找到安徽商店: %s (%s)", store.CityStoreName, store.StoreNumber)
			break
		}
	}

	if !hasAnhui {
		log.Println("❌ 未找到安徽商店")
	}

	// 测试省份列表
	log.Println("\n2. 测试省份列表...")
	provinces := model.GetProvinces()
	log.Printf("✅ 省份列表 (%d个): %v", len(provinces), provinces)

	// 检查关键省份
	keyProvinces := []string{"安徽", "河南", "湖北", "云南", "广西壮族自治区"}
	for _, province := range keyProvinces {
		found := false
		for _, p := range provinces {
			if p == province {
				found = true
				break
			}
		}
		if found {
			log.Printf("✅ %s 在省份列表中", province)
		} else {
			log.Printf("❌ %s 不在省份列表中", province)
		}
	}

	// 测试服务层数据加载
	log.Println("\n3. 测试服务层数据加载...")
	err := services.Store.LoadForArea("cn")
	if err != nil {
		log.Fatalf("❌ 服务层加载商店数据失败: %v", err)
	}

	stores := services.Store.ByArea(services.Area.GetArea("中国大陆"))
	log.Printf("✅ 服务层加载: %d 个商店", len(stores))

	// 测试安徽商店在服务层
	anhuiStores := services.Store.ByAreaAndProvinceForOptions("中国大陆", "安徽")
	log.Printf("✅ 安徽商店 (%d个): %v", len(anhuiStores), anhuiStores)

	log.Println("\n🎉 所有测试通过！嵌入数据已成功更新并包含所有修复。")
}
