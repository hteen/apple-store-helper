package services

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"github.com/parnurzeal/gorequest"

	"apple-store-helper/embedded"
	"apple-store-helper/model"
)

// ProductData 存储从Apple官网获取的产品数据
type ProductData struct {
	UpdateTime string                         `json:"update_time"`
	AreaCode   string                         `json:"area_code"`
	Products   map[string][]model.ProductInfo `json:"products"`
}

// FetchProductData 从Apple官网获取产品数据
func FetchProductData(areaCode string) (*ProductData, error) {
	productData := &ProductData{
		UpdateTime: time.Now().Format("2006-01-02 15:04:05"),
		AreaCode:   areaCode,
		Products:   make(map[string][]model.ProductInfo),
	}

	// 根据地区代码构建基础URL - 支持中国大陆、香港、日本、新加坡、美国、英国和澳大利亚
	baseURL := ""
	switch areaCode {
	case "cn":
		baseURL = "https://www.apple.com.cn"
	case "hk":
		baseURL = "https://www.apple.com/hk"
	case "jp":
		baseURL = "https://www.apple.com/jp"
	case "sg":
		baseURL = "https://www.apple.com/sg"
	case "us":
		baseURL = "https://www.apple.com/us"
	case "uk":
		baseURL = "https://www.apple.com/uk"
	case "au":
		baseURL = "https://www.apple.com/au"
	default:
		return nil, fmt.Errorf("unsupported area code: %s", areaCode)
	}

	// 构建所有产品系列的URL
	// 只获取当前存在的产品系列
	series := []struct {
		name      string
		url       string
		modelType string
	}{
		{"iPhone 16", fmt.Sprintf("%s/shop/buy-iphone/iphone-16", baseURL), "iphone16"},
		{"iPhone 16 Pro", fmt.Sprintf("%s/shop/buy-iphone/iphone-16-pro", baseURL), "iphone16pro"},
		// iPhone 17 系列
		{"iPhone 17", fmt.Sprintf("%s/shop/buy-iphone/iphone-17", baseURL), "iphone17"},
		{"iPhone 17 Pro", fmt.Sprintf("%s/shop/buy-iphone/iphone-17-pro", baseURL), "iphone17pro"},
		{"iPhone Air", fmt.Sprintf("%s/shop/buy-iphone/iphone-air", baseURL), "iphoneair"},
	}

	for _, s := range series {
		products, err := fetchSeriesProducts(s.url, s.modelType)
		if err != nil {
			log.Printf("Failed to fetch %s: %v", s.name, err)
			continue
		}
		if len(products) > 0 {
			productData.Products[s.name] = products
		}
	}

	return productData, nil
}

// fetchSeriesProducts 获取特定系列的产品
func fetchSeriesProducts(url string, modelType string) ([]model.ProductInfo, error) {
	resp, body, errs := gorequest.New().
		Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36").
		Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8").
		Set("Accept-Language", "zh-CN,zh;q=0.9").
		Timeout(time.Second * 10).
		Get(url).
		End()

	if len(errs) > 0 {
		return nil, fmt.Errorf("request failed: %v", errs[0])
	}

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("status code: %d", resp.StatusCode)
	}

	// 查找 metrics script 标签 - 使用更灵活的正则表达式
	re := regexp.MustCompile(`<script[^>]*id=["']metrics["'][^>]*>(.*?)</script>`)
	matches := re.FindStringSubmatch(body)
	if len(matches) < 2 {
		// 尝试另一种格式
		re = regexp.MustCompile(`<script\s+type=["']application/json["']\s+id=["']metrics["']>(.*?)</script>`)
		matches = re.FindStringSubmatch(body)
		if len(matches) < 2 {
			return nil, fmt.Errorf("metrics data not found in HTML")
		}
	}

	// 解析JSON数据
	var metricsData map[string]interface{}
	if err := json.Unmarshal([]byte(matches[1]), &metricsData); err != nil {
		return nil, fmt.Errorf("failed to parse metrics data: %v", err)
	}

	// 提取products数组
	products := []model.ProductInfo{}
	if data, ok := metricsData["data"].(map[string]interface{}); ok {
		if productsArray, ok := data["products"].([]interface{}); ok {
			for _, p := range productsArray {
				if product, ok := p.(map[string]interface{}); ok {
					partNumber, _ := product["partNumber"].(string)
					name, _ := product["name"].(string)
					sku, _ := product["sku"].(string)

					// 记录提取的原始数据
					log.Printf("Found product: SKU=%s, PartNumber=%s, Name=%s", sku, partNumber, name)

					// 解析产品信息
					info := parseProductInfo(name, partNumber, modelType)
					if info.Code != "" {
						products = append(products, info)
					}
				}
			}
		}
	}

	return products, nil
}

// normalizeSpaces 规范化字符串中的各种空格字符
func normalizeSpaces(s string) string {
	// 替换各种Unicode空格字符为普通空格
	// \u00A0 = 不间断空格 (NBSP)
	// \u2002 = En空格
	// \u2003 = Em空格
	// \u3000 = 全角空格
	s = strings.ReplaceAll(s, "\u00A0", " ")
	s = strings.ReplaceAll(s, "\u2002", " ")
	s = strings.ReplaceAll(s, "\u2003", " ")
	s = strings.ReplaceAll(s, "\u3000", " ")

	// 将多个连续空格替换为单个空格
	re := regexp.MustCompile(`\s+`)
	s = re.ReplaceAllString(s, " ")

	// 去除首尾空格
	return strings.TrimSpace(s)
}

// parseProductInfo 解析产品信息
func parseProductInfo(name string, partNumber string, modelType string) model.ProductInfo {
	info := model.ProductInfo{
		Code: partNumber,
		Type: modelType,
	}

	// 如果name为空，直接返回
	if name == "" {
		return info
	}

	// 规范化名称中的空格
	name = normalizeSpaces(name)

	// 解析产品名称
	parts := strings.Split(name, " ")
	if len(parts) >= 3 {
		// 格式: "iPhone 16 Plus 128GB Black" 或 "iPhone 16 128GB Black"
		modelParts := []string{}
		capacityIdx := -1

		for i, part := range parts {
			part = strings.TrimSpace(part)
			if strings.HasSuffix(part, "GB") || strings.HasSuffix(part, "TB") {
				capacityIdx = i
				break
			}
			if part != "" {
				modelParts = append(modelParts, part)
			}
		}

		if capacityIdx > 0 && len(modelParts) > 0 {
			info.Model = strings.Join(modelParts, " ")
			info.Capacity = strings.TrimSpace(parts[capacityIdx])
			if capacityIdx+1 < len(parts) {
				var colorParts []string
				for j := capacityIdx + 1; j < len(parts); j++ {
					part := strings.TrimSpace(parts[j])
					if part != "" {
						colorParts = append(colorParts, part)
					}
				}
				if len(colorParts) > 0 {
					info.Color = translateColor(strings.Join(colorParts, " "))
				}
			}
		}
	}

	return info
}

// translateColor 翻译颜色名称到官方中文颜色
func translateColor(color string) string {
	// 如果颜色为空，直接返回
	if color == "" {
		return color
	}

	// Apple 官方颜色翻译映射表
	colorMap := map[string]string{
		// iPhone 基础颜色
		"Black":       "黑色",
		"White":       "白色",
		"Pink":        "粉色",
		"Teal":        "深青色",
		"Ultramarine": "群青色",

		// iPhone Pro 系列颜色
		"Black Titanium":   "黑色钛金属",
		"White Titanium":   "白色钛金属",
		"Natural Titanium": "原色钛金属",
		"Desert Titanium":  "沙漠色钛金属",

		// iPhone 17 系列颜色
		"Sage":          "鼠尾草色",
		"Lavender":      "薰衣草色",
		"Mist Blue":     "薄雾蓝色",
		"Space Orange":  "宇宙橙色",
		"Deep Blue":     "深蓝色",
		"Silver":        "银色",
		"Cosmic Orange": "宇宙橙色",

		// iPhone Air 系列颜色
		"Cloud White":      "云白色",
		"Sky Blue":         "天蓝色",
		"Deep Space Black": "深空黑色",
		"Light Gold":       "浅金色",

		// Apple Watch 系列颜色
		"Space Black":   "深空黑色",
		"Gold":          "金色",
		"Rose Gold":     "玫瑰金色",
		"Midnight":      "午夜色",
		"Starlight":     "星光色",
		"Blue":          "蓝色",
		"Red":           "红色",
		"Green":         "绿色",
		"Yellow":        "黄色",
		"Orange":        "橙色",
		"Purple":        "紫色",
		"Product Red":   "红色",
		"Forest Green":  "森林绿色",
		"Ocean Blue":    "海洋蓝色",
		"Sunset Orange": "日落橙色",

		// 其他常见颜色
		"Graphite":          "石墨色",
		"Titanium":          "钛金属色",
		"Aluminum":          "铝金属色",
		"Stainless Steel":   "不锈钢色",
		"Ceramic":           "陶瓷色",
		"Leather":           "皮革色",
		"Fabric":            "织物色",
		"Sport Band":        "运动表带",
		"Sport Loop":        "运动表环",
		"Braided Solo Loop": "编织单圈表带",
		"Solo Loop":         "单圈表带",
		"Link Bracelet":     "链式表带",
		"Milanese Loop":     "米兰尼斯表带",
		"Leather Link":      "皮革链式表带",
		"Leather Loop":      "皮革表环",
		"Modern Buckle":     "现代扣式表带",
		"Classic Buckle":    "经典扣式表带",
	}

	if cn, ok := colorMap[color]; ok {
		return cn
	}

	// 如果没有找到翻译，尝试部分匹配
	for english, chinese := range colorMap {
		if strings.Contains(strings.ToLower(color), strings.ToLower(english)) {
			return chinese
		}
	}

	// 如果都没有找到，返回原始颜色
	return color
}

// SaveProductData 保存产品数据到本地文件
func SaveProductData(data *ProductData) error {
	// 获取可执行文件所在目录
	execDir, err := os.Executable()
	if err != nil {
		// 如果获取可执行文件路径失败，使用当前工作目录
		execDir, err = os.Getwd()
		if err != nil {
			return fmt.Errorf("failed to get current directory: %v", err)
		}
	} else {
		execDir = filepath.Dir(execDir)
	}

	// 在程序目录下创建data/product子目录
	dataDir := filepath.Join(execDir, "data", "product")
	if err := os.MkdirAll(dataDir, 0755); err != nil {
		return err
	}

	// 使用地区代码作为文件名的一部分
	fileName := fmt.Sprintf("product_data_%s.json", data.AreaCode)
	filePath := filepath.Join(dataDir, fileName)

	jsonData, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return err
	}

	return ioutil.WriteFile(filePath, jsonData, 0644)
}

// LoadProductData 从本地文件加载产品数据
func LoadProductData(areaCode string) (*ProductData, error) {
	// 首先尝试从嵌入的数据加载
	if data, exists := embedded.GetProductData(areaCode); exists {
		var productData ProductData
		if err := json.Unmarshal(data, &productData); err != nil {
			return nil, fmt.Errorf("failed to unmarshal embedded product data for %s: %v", areaCode, err)
		}
		log.Printf("Loaded embedded product data for %s", areaCode)
		return &productData, nil
	}

	// 如果嵌入数据不存在，回退到文件系统
	workDir, err := os.Getwd()
	if err != nil {
		return nil, fmt.Errorf("failed to get current directory: %v", err)
	}

	// 从 data/product 目录加载
	fileName := fmt.Sprintf("product_data_%s.json", areaCode)
	filePath := filepath.Join(workDir, "data", "product", fileName)

	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return nil, fmt.Errorf("product data file not found: %s", filePath)
	}

	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	var productData ProductData
	if err := json.Unmarshal(data, &productData); err != nil {
		return nil, err
	}

	log.Printf("Loaded product data for %s from %s", areaCode, filePath)
	return &productData, nil
}

// UpdateProductDatabase 更新产品数据库
func UpdateProductDatabase(areaCode string) error {
	log.Println("Fetching latest product data from Apple...")

	data, err := FetchProductData(areaCode)
	if err != nil {
		return fmt.Errorf("failed to fetch product data: %v", err)
	}

	if err := SaveProductData(data); err != nil {
		return fmt.Errorf("failed to save product data: %v", err)
	}

	log.Printf("Product data updated successfully. Total series: %d", len(data.Products))

	// 更新Product服务的产品列表
	Product.UpdateFromDynamicData(data)

	return nil
}
