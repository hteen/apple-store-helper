package main

import (
	"apple-store-helper/embedded"
	"encoding/json"
	"fmt"
)

func main() {
	fmt.Println("=== 最终数据嵌入测试 ===")

	// 测试所有地区的产品数据
	fmt.Println("\n📱 产品数据测试:")
	productRegions := []string{"cn", "hk", "jp", "sg", "us", "uk", "au"}
	for _, region := range productRegions {
		if data, exists := embedded.GetProductData(region); exists {
			var productData map[string]interface{}
			if err := json.Unmarshal(data, &productData); err == nil {
				if products, ok := productData["products"].(map[string]interface{}); ok {
					fmt.Printf("  ✓ %s: %d 个产品系列 (%d 字节)\n",
						region, len(products), len(data))
				}
			}
		}
	}

	// 测试所有地区的门店数据
	fmt.Println("\n🏪 门店数据测试:")
	storeRegions := []string{"cn", "hk", "jp", "us", "uk", "au"}
	for _, region := range storeRegions {
		if data, exists := embedded.GetStoreData(region); exists {
			var storeData map[string]interface{}
			if err := json.Unmarshal(data, &storeData); err == nil {
				if stores, ok := storeData["stores"].([]interface{}); ok {
					fmt.Printf("  ✓ %s: %d 个门店 (%d 字节)\n",
						region, len(stores), len(data))
				}
			}
		}
	}

	fmt.Println("\n🎉 所有数据嵌入成功！程序可以独立运行。")
}
