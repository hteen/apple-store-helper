package services

import (
	"apple-store-helper/model"
	"fmt"
	"log"
)

var Product = &productService{
	dynamicProducts: make(map[string][]model.Product),
	currentAreaCode: "cn", // 默认中国大陆
	useDynamic:      true, // 默认使用动态数据
}

type productService struct {
	dynamicProducts map[string][]model.Product
	currentAreaCode string
	useDynamic      bool
}

func (s *productService) ByAreaTitleForOptions(areaTitle string) []string {
	// 只使用动态数据
	var titles []string
	for _, products := range s.dynamicProducts {
		for _, p := range products {
			titles = append(titles, p.Title)
		}
	}
	return titles
}

func (s *productService) GetProduct(areaTitle string, productTitle string) model.Product {
	// 只使用动态数据
	for _, products := range s.dynamicProducts {
		for _, p := range products {
			if p.Title == productTitle {
				return p
			}
		}
	}

	// 如果没找到，返回空的Product
	return model.Product{}
}

// UpdateFromDynamicData 从动态数据更新产品列表
func (s *productService) UpdateFromDynamicData(data *ProductData) {
	s.dynamicProducts = make(map[string][]model.Product)
	s.currentAreaCode = data.AreaCode

	for series, products := range data.Products {
		var modelProducts []model.Product
		for _, p := range products {
			// 构建标题 - 只有当所有字段都存在时才生成有效的Title
			if p.Model != "" && p.Capacity != "" && p.Color != "" {
				title := fmt.Sprintf("%s %s %s", p.Model, p.Capacity, p.Color)
				modelProducts = append(modelProducts, model.Product{
					Title: title,
					Code:  p.Code,
					Type:  p.Type,
				})
			} else if p.Code != "" {
				// 如果解析失败但有Code，记录警告
				log.Printf("Warning: Product with code %s has incomplete data (Model: %s, Capacity: %s, Color: %s)",
					p.Code, p.Model, p.Capacity, p.Color)
			}
		}
		if len(modelProducts) > 0 {
			s.dynamicProducts[series] = modelProducts
		}
	}

	s.useDynamic = true
	log.Printf("Updated product database with %d series", len(s.dynamicProducts))
}

// GetDynamicProducts 获取动态产品列表
func (s *productService) GetDynamicProducts() map[string][]model.Product {
	return s.dynamicProducts
}

// SetUseDynamic 设置是否使用动态数据
func (s *productService) SetUseDynamic(use bool) {
	s.useDynamic = use
}

// GetCurrentAreaCode 获取当前地区代码
func (s *productService) GetCurrentAreaCode() string {
	return s.currentAreaCode
}

// LoadForArea 加载指定地区的产品数据
func (s *productService) LoadForArea(areaCode string) error {
	data, err := LoadProductData(areaCode)
	if err != nil {
		return err
	}
	s.UpdateFromDynamicData(data)
	return nil
}
