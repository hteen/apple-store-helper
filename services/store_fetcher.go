package services

import (
	"apple-store-helper/embedded"
	"apple-store-helper/model"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"math/rand"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"time"
)

// StoreData 门店数据结构
type StoreData struct {
	UpdateTime string        `json:"update_time"`
	AreaCode   string        `json:"area_code"`
	Stores     []model.Store `json:"stores"`
}

// FetchStoresForArea 获取指定地区的门店列表
func FetchStoresForArea(areaCode string, location string) ([]model.Store, error) {
	// 使用fulfillment API获取门店列表
	// 需要一个有效的产品代码
	var apiURL string
	var sampleProduct string

	switch areaCode {
	case "cn":
		sampleProduct = "MYEW3CH/A" // iPhone 16 白色 128GB
		// 中国大陆使用邮编格式，如果location不是邮编则转换为邮编
		postalCode := getPostalCodeForLocation(location)
		apiURL = fmt.Sprintf("https://www.apple.com.cn/shop/fulfillment-messages?fae=true&pl=true&mts.0=regular&parts.0=%s&location=%s",
			sampleProduct, postalCode)
	case "hk":
		sampleProduct = "MYEW3ZA/A" // 香港的iPhone 16 白色 128GB
		apiURL = fmt.Sprintf("https://www.apple.com/hk/shop/fulfillment-messages?fae=true&pl=true&mts.0=regular&parts.0=%s&location=Central",
			sampleProduct)
	case "jp":
		sampleProduct = "MYDR3J/A" // 日本的iPhone 16
		apiURL = fmt.Sprintf("https://www.apple.com/jp/shop/fulfillment-messages?fae=true&pl=true&mts.0=regular&parts.0=%s&location=%s&cppart=UNLOCKED_JP",
			sampleProduct, url.QueryEscape(location))
	case "sg":
		sampleProduct = "MXY23ZP/A" // 新加坡的iPhone 16 Plus 256GB Ultramarine
		apiURL = fmt.Sprintf("https://www.apple.com/sg/shop/fulfillment-messages?fae=true&pl=true&mts.0=regular&parts.0=%s&location=Singapore",
			sampleProduct)
	case "us":
		sampleProduct = "MYAR3LL/A" // 美国的iPhone 16 128GB Black
		apiURL = fmt.Sprintf("https://www.apple.com/us/shop/fulfillment-messages?fae=true&pl=true&mts.0=regular&parts.0=%s&location=%s",
			sampleProduct, url.QueryEscape(location))
	case "uk":
		sampleProduct = "MYE93QN/A" // 英国的iPhone 16 128GB White
		apiURL = fmt.Sprintf("https://www.apple.com/uk/shop/fulfillment-messages?fae=true&pl=true&mts.0=regular&parts.0=%s&location=%s",
			sampleProduct, url.QueryEscape(location))
	case "au":
		sampleProduct = "MYE93X/A" // 澳大利亚的iPhone 16 128GB White
		apiURL = fmt.Sprintf("https://www.apple.com/au/shop/fulfillment-messages?fae=true&pl=true&mts.0=regular&parts.0=%s&location=%s",
			sampleProduct, url.QueryEscape(location))
	default:
		return nil, fmt.Errorf("unsupported area code: %s", areaCode)
	}

	// 创建HTTP请求
	req, err := http.NewRequest("GET", apiURL, nil)
	if err != nil {
		return nil, err
	}

	// 设置动态请求头
	userAgents := []string{
		"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36",
		"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/17.1 Safari/605.1.15",
		"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36",
	}
	req.Header.Set("User-Agent", userAgents[rand.Intn(len(userAgents))])
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Accept-Language", getAcceptLanguage(areaCode))
	req.Header.Set("Referer", getReferer(areaCode))

	// 发送请求
	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// 读取响应
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	// 解析fulfillment API响应
	var fulfillmentResponse struct {
		Body struct {
			Content struct {
				PickupMessage struct {
					Stores []struct {
						StoreNumber       string `json:"storeNumber"`
						StoreName         string `json:"storeName"`
						State             string `json:"state"`
						City              string `json:"city"`
						StoreEmail        string `json:"storeEmail"`
						PartsAvailability map[string]struct {
							PickupDisplay         string `json:"pickupDisplay"`
							StoreSelectionEnabled bool   `json:"storeSelectionEnabled"`
						} `json:"partsAvailability"`
					} `json:"stores"`
				} `json:"pickupMessage"`
			} `json:"content"`
		} `json:"body"`
	}

	if err := json.Unmarshal(body, &fulfillmentResponse); err != nil {
		return nil, err
	}

	// 转换为model.Store格式
	var stores []model.Store
	storeMap := make(map[string]bool) // 去重

	for _, s := range fulfillmentResponse.Body.Content.PickupMessage.Stores {
		// 构建门店名称
		storeName := s.StoreName
		if s.State != "" && s.State != s.StoreName {
			storeName = fmt.Sprintf("%s-%s", s.State, s.StoreName)
		}

		// 去重
		if storeMap[storeName] {
			continue
		}
		storeMap[storeName] = true

		store := model.Store{
			StoreNumber:   s.StoreNumber,
			CityStoreName: storeName,
			City:          s.City,
			Province:      s.State,
			District:      "", // fulfillment API不提供详细地址
		}
		stores = append(stores, store)
	}

	return stores, nil
}

// LoadStoreData 从嵌入数据或文件系统加载门店数据
func LoadStoreData(areaCode string) (*StoreData, error) {
	// 首先尝试从嵌入的数据加载
	if data, exists := embedded.GetStoreData(areaCode); exists {
		var storeData StoreData
		if err := json.Unmarshal(data, &storeData); err != nil {
			return nil, fmt.Errorf("failed to unmarshal embedded store data for %s: %v", areaCode, err)
		}
		log.Printf("Loaded embedded store data for %s: %d stores", areaCode, len(storeData.Stores))
		return &storeData, nil
	}

	// 如果嵌入数据不存在，回退到文件系统
	workDir, err := os.Getwd()
	if err != nil {
		return nil, fmt.Errorf("failed to get current directory: %v", err)
	}

	// 使用 data/store 目录
	storeDataDir := filepath.Join(workDir, "data", "store")
	filename := fmt.Sprintf("store_%s.json", areaCode)
	filepath := filepath.Join(storeDataDir, filename)

	// 检查文件是否存在
	if _, err := os.Stat(filepath); os.IsNotExist(err) {
		return nil, fmt.Errorf("store data file not found: %s", filepath)
	}

	// 读取文件
	data, err := os.ReadFile(filepath)
	if err != nil {
		return nil, err
	}

	// 解析JSON
	var storeData StoreData
	if err := json.Unmarshal(data, &storeData); err != nil {
		return nil, err
	}

	// 记录加载成功
	log.Printf("Loaded %d stores for %s from store_data directory", len(storeData.Stores), areaCode)

	return &storeData, nil
}

// SaveStoreData 保存门店数据到本地
func SaveStoreData(areaCode string, stores []model.Store) error {
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

	// 在程序目录下创建data/store子目录
	dataDir := filepath.Join(execDir, "data", "store")
	if err := os.MkdirAll(dataDir, 0755); err != nil {
		return err
	}

	// 构建数据结构
	storeData := StoreData{
		UpdateTime: time.Now().Format("2006-01-02 15:04:05"),
		AreaCode:   areaCode,
		Stores:     stores,
	}

	// 序列化为JSON
	jsonData, err := json.MarshalIndent(storeData, "", "  ")
	if err != nil {
		return err
	}

	// 写入文件
	filename := fmt.Sprintf("stores_%s.json", areaCode)
	filepath := filepath.Join(dataDir, filename)
	if err := os.WriteFile(filepath, jsonData, 0644); err != nil {
		return err
	}

	log.Printf("Saved %d stores for area %s to %s", len(stores), areaCode, filepath)
	return nil
}

// UpdateStoresForAllAreas 更新所有地区的门店数据
func UpdateStoresForAllAreas() error {
	// 定义所有地区的主要城市/位置，用于获取完整的门店列表
	areas := map[string][]string{
		"cn": {"北京", "上海", "深圳", "广州", "成都", "杭州", "南京", "武汉", "重庆", "天津", "苏州", "青岛", "长沙", "大连", "厦门", "无锡", "福州", "济南", "宁波", "温州", "郑州", "沈阳", "合肥", "南宁", "昆明"},
		"hk": {"Central"},
		"jp": {"100-0001", "150-0001", "460-0008", "530-0001", "650-0001", "700-0001", "800-0001", "900-0001"},
		"sg": {"Singapore"},
		"us": {"10001", "90210", "60601", "33101", "75201", "98101", "85001", "30309", "02108", "02116"},
		"uk": {"London", "Manchester", "Birmingham", "Glasgow", "Edinburgh", "Liverpool", "Leeds", "Bristol", "Newcastle", "Cardiff"},
		"au": {"Sydney", "Melbourne", "Brisbane", "Perth", "Adelaide", "Gold Coast", "Newcastle", "Wollongong", "Geelong", "Hobart", "Darwin", "Canberra"},
	}

	totalStores := 0
	for areaCode, locations := range areas {
		log.Printf("Fetching stores for %s...", areaCode)

		var allStores []model.Store
		storeMap := make(map[string]bool) // 用于去重

		for _, location := range locations {
			stores, err := FetchStoresForArea(areaCode, location)
			if err != nil {
				log.Printf("Failed to fetch stores for %s %s: %v", areaCode, location, err)
				continue
			}

			// 合并门店数据并去重
			for _, store := range stores {
				key := store.StoreNumber
				if !storeMap[key] {
					storeMap[key] = true
					allStores = append(allStores, store)
				}
			}

			// 添加随机延迟，避免频繁请求
			time.Sleep(time.Duration(1+rand.Intn(2)) * time.Second)
		}

		if len(allStores) > 0 {
			if err := SaveStoreData(areaCode, allStores); err != nil {
				log.Printf("Failed to save stores for %s: %v", areaCode, err)
				continue
			}
			log.Printf("Successfully saved %d stores for %s", len(allStores), areaCode)
			totalStores += len(allStores)
		}

		// 地区间添加较长延迟
		time.Sleep(time.Duration(3+rand.Intn(3)) * time.Second)
	}

	log.Printf("Total stores updated: %d", totalStores)
	return nil
}

// getAcceptLanguage 根据地区返回Accept-Language
func getAcceptLanguage(areaCode string) string {
	switch areaCode {
	case "cn":
		return "zh-CN,zh;q=0.9,en;q=0.8"
	case "hk":
		return "zh-HK,zh;q=0.9,en;q=0.8"
	case "jp":
		return "ja-JP,ja;q=0.9,en;q=0.8"
	case "sg":
		return "en-SG,en;q=0.9,zh;q=0.8"
	case "us":
		return "en-US,en;q=0.9"
	case "uk":
		return "en-GB,en;q=0.9"
	case "au":
		return "en-AU,en;q=0.9"
	default:
		return "en-US,en;q=0.9"
	}
}

// getReferer 根据地区返回Referer
func getReferer(areaCode string) string {
	switch areaCode {
	case "cn":
		return "https://www.apple.com.cn/shop/buy-iphone/iphone-16"
	case "hk":
		return "https://www.apple.com/hk/shop/buy-iphone/iphone-16"
	case "jp":
		return "https://www.apple.com/jp/shop/buy-iphone/iphone-16"
	case "sg":
		return "https://www.apple.com/sg/shop/buy-iphone/iphone-16"
	case "us":
		return "https://www.apple.com/us/shop/buy-iphone/iphone-16"
	case "uk":
		return "https://www.apple.com/uk/shop/buy-iphone/iphone-16"
	case "au":
		return "https://www.apple.com/au/shop/buy-iphone/iphone-16"
	default:
		return "https://www.apple.com/shop/buy-iphone/iphone-16"
	}
}

// getPostalCodeForLocation 将城市名称转换为邮编
func getPostalCodeForLocation(location string) string {
	// 主要城市的邮编映射
	cityPostalCodes := map[string]string{
		"北京":   "100000",
		"上海":   "200000",
		"深圳":   "518000",
		"广州":   "510000",
		"成都":   "610000",
		"杭州":   "310000",
		"南京":   "210000",
		"武汉":   "430000",
		"西安":   "710000",
		"重庆":   "400000",
		"天津":   "300000",
		"苏州":   "215000",
		"青岛":   "266000",
		"长沙":   "410000",
		"大连":   "116000",
		"厦门":   "361000",
		"无锡":   "214000",
		"福州":   "350000",
		"济南":   "250000",
		"宁波":   "315000",
		"温州":   "325000",
		"郑州":   "450000",
		"沈阳":   "110000",
		"哈尔滨":  "150000",
		"石家庄":  "050000",
		"太原":   "030000",
		"呼和浩特": "010000",
		"长春":   "130000",
		"合肥":   "230000",
		"南昌":   "330000",
		"南宁":   "530000",
		"海口":   "570000",
		"贵阳":   "550000",
		"昆明":   "650000",
		"兰州":   "730000",
		"西宁":   "810000",
		"银川":   "750000",
		"乌鲁木齐": "830000",
	}

	// 如果已经是邮编格式（6位数字），直接返回
	if len(location) == 6 {
		isNumeric := true
		for _, char := range location {
			if char < '0' || char > '9' {
				isNumeric = false
				break
			}
		}
		if isNumeric {
			return location
		}
	}

	// 查找城市对应的邮编
	if postalCode, exists := cityPostalCodes[location]; exists {
		return postalCode
	}

	// 默认返回北京邮编
	return "100000"
}
