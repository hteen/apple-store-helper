# Apple Store Helper - 新地区和新产品适配指南

本文档详细说明如何为 Apple Store Helper 添加新地区支持和新产品系列支持，确保任何人都可以按照此流程进行适配。

## 📋 目录

1. [概述](#概述)
2. [新地区适配流程](#新地区适配流程)
3. [新产品系列适配流程](#新产品系列适配流程)
4. [技术细节说明](#技术细节说明)
5. [测试验证](#测试验证)
6. [常见问题](#常见问题)

## 概述

Apple Store Helper 支持多个地区的 Apple Store 库存监控，包括：
- 中国大陆 (cn)
- 香港 (hk) 
- 日本 (jp)
- 新加坡 (sg)

每个地区都有独立的产品数据和门店数据，通过 Apple 官方 API 获取。

## 新地区适配流程

### 步骤 1: 添加地区定义

**文件**: `model/area.go`

```go
// 在 Areas 数组中添加新地区
var Areas = []Area{
    {Title: "中国大陆", Locale: "zh_CN", ShortCode: "cn", ProductsJson: iPhone17ProductsJson},
    {Title: "香港", Locale: "zh_HK", ShortCode: "hk", ProductsJson: iPhone17ProductsJson},
    {Title: "日本", Locale: "ja_JP", ShortCode: "jp", ProductsJson: iPhone17ProductsJson},
    {Title: "新加坡", Locale: "en_SG", ShortCode: "sg", ProductsJson: iPhone17ProductsJson},
    // 添加新地区，例如：
    {Title: "美国", Locale: "en_US", ShortCode: "us", ProductsJson: iPhone17ProductsJson},
}
```

**参数说明**:
- `Title`: 在UI中显示的地区名称
- `Locale`: 地区语言代码 (如: en_US, zh_CN, ja_JP)
- `ShortCode`: 地区代码，用于API调用和文件命名
- `ProductsJson`: 产品JSON配置 (通常使用 iPhone17ProductsJson)

### 步骤 2: 更新产品获取逻辑

**文件**: `services/product_fetcher.go`

在 `FetchProductData` 函数的 switch 语句中添加新地区：

```go
// 根据地区代码构建基础URL
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
case "us":  // 新地区
    baseURL = "https://www.apple.com/us"
default:
    return nil, fmt.Errorf("unsupported area code: %s", areaCode)
}
```

### 步骤 3: 更新门店获取逻辑

**文件**: `services/store_fetcher.go`

#### 3.1 添加API调用逻辑

在 `FetchStoresForArea` 函数的 switch 语句中添加：

```go
switch areaCode {
case "cn":
    sampleProduct = "MYEW3CH/A" // iPhone 16 白色 128GB
    postalCode := getPostalCodeForLocation(location)
    apiURL = fmt.Sprintf("https://www.apple.com.cn/shop/fulfillment-messages?fae=true&pl=true&mts.0=regular&parts.0=%s&location=%s",
        sampleProduct, postalCode)
case "hk":
    sampleProduct = "MYEW3ZA/A"
    apiURL = fmt.Sprintf("https://www.apple.com/hk/shop/fulfillment-messages?fae=true&pl=true&mts.0=regular&parts.0=%s&location=Central",
        sampleProduct)
case "jp":
    sampleProduct = "MYDR3J/A"
    apiURL = fmt.Sprintf("https://www.apple.com/jp/shop/fulfillment-messages?fae=true&pl=true&mts.0=regular&parts.0=%s&location=%s&cppart=UNLOCKED_JP",
        sampleProduct, url.QueryEscape(location))
case "sg":
    sampleProduct = "MXY23ZP/A"
    apiURL = fmt.Sprintf("https://www.apple.com/sg/shop/fulfillment-messages?fae=true&pl=true&mts.0=regular&parts.0=%s&location=Singapore",
        sampleProduct)
case "us":  // 新地区
    sampleProduct = "MYEW3LL/A" // 美国iPhone 16 白色 128GB
    apiURL = fmt.Sprintf("https://www.apple.com/us/shop/fulfillment-messages?fae=true&pl=true&mts.0=regular&parts.0=%s&location=%s",
        sampleProduct, url.QueryEscape(location))
default:
    return nil, fmt.Errorf("unsupported area code: %s", areaCode)
}
```

**重要**: 需要找到该地区有效的产品代码，可以通过以下方式获取：

1. 访问该地区的 Apple 官网产品页面
2. 查看页面源码中的 `<script type="application/json" id="metrics">` 标签
3. 找到产品代码 (如: MYEW3LL/A)

#### 3.2 添加地区配置

在 `UpdateStoresForAllAreas` 函数中添加新地区：

```go
areas := map[string][]string{
    "cn": {"北京", "上海", "深圳", "广州", "成都", "杭州", "南京", "武汉", "西安", "重庆", "天津", "苏州", "青岛", "长沙", "大连", "厦门", "无锡", "福州", "济南", "宁波", "温州", "郑州", "沈阳", "哈尔滨", "石家庄", "太原", "呼和浩特", "长春", "合肥", "南昌", "南宁", "海口", "贵阳", "昆明", "兰州", "西宁", "银川", "乌鲁木齐"},
    "hk": {"Central"},
    "jp": {"100-0001", "150-0001", "460-0008", "530-0001", "650-0001", "700-0001", "800-0001", "900-0001"},
    "sg": {"Singapore"},
    "us": {"10001", "90210", "60601", "33101", "75201"}, // 美国主要城市邮编
}
```

#### 3.3 添加语言和Referer配置

在 `getAcceptLanguage` 函数中添加：

```go
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
    case "us":  // 新地区
        return "en-US,en;q=0.9"
    default:
        return "en-US,en;q=0.9"
    }
}
```

在 `getReferer` 函数中添加：

```go
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
    case "us":  // 新地区
        return "https://www.apple.com/us/shop/buy-iphone/iphone-16"
    default:
        return "https://www.apple.com/shop/buy-iphone/iphone-16"
    }
}
```

### 步骤 4: 更新UI逻辑

**文件**: `main.go`

#### 4.1 添加地区到更新列表

在更新数据按钮的逻辑中添加新地区：

```go
// 更新所有地区的产品数据
areaCodes := []string{"cn", "hk", "jp", "sg", "us"} // 添加新地区
```

#### 4.2 添加UI显示逻辑

在地区选择器的回调函数中添加新地区的UI逻辑：

```go
} else if value == "新加坡" {
    // 新加坡：直接显示所有门店
    zipCodeEntry.Disable()
    zipCodeEntry.Text = ""
    storeWidget.Options = services.Store.ByAreaTitleForOptions(value)
    storeWidget.Enable()
    storeWidget.PlaceHolder = "选择门店"
    locationContainer.Objects = []fyne.CanvasObject{
        container.NewVBox(
            widget.NewLabelWithStyle("门店", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
            storeWidget,
        ),
    }
} else if value == "美国" {  // 新地区
    // 美国：直接显示所有门店
    zipCodeEntry.Disable()
    zipCodeEntry.Text = ""
    storeWidget.Options = services.Store.ByAreaTitleForOptions(value)
    storeWidget.Enable()
    storeWidget.PlaceHolder = "选择门店"
    locationContainer.Objects = []fyne.CanvasObject{
        container.NewVBox(
            widget.NewLabelWithStyle("门店", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
            storeWidget,
        ),
    }
}
```

## 新产品系列适配流程

### 步骤 1: 添加产品系列定义

**文件**: `services/product_fetcher.go`

在 `FetchProductData` 函数的 series 数组中添加新产品系列：

```go
series := []struct {
    name string
    url  string
}{
    {"iPhone 16", baseURL + "/shop/buy-iphone/iphone-16"},
    {"iPhone 16 Plus", baseURL + "/shop/buy-iphone/iphone-16-plus"},
    {"iPhone 16 Pro", baseURL + "/shop/buy-iphone/iphone-16-pro"},
    {"iPhone 16 Pro Max", baseURL + "/shop/buy-iphone/iphone-16-pro-max"},
    {"Apple Watch Series 10", baseURL + "/shop/buy-watch/apple-watch-series-10"}, // 新产品系列
    {"Apple Watch Ultra 2", baseURL + "/shop/buy-watch/apple-watch-ultra-2"},     // 新产品系列
}
```

### 步骤 2: 更新产品解析逻辑

如果新产品系列有特殊的解析需求，需要在 `parseProductFromMetrics` 函数中添加相应的处理逻辑。

### 步骤 3: 更新产品类型映射

**文件**: `model/product.go`

如果需要新的产品类型，在 `TypeCode` 映射中添加：

```go
var TypeCode = map[string]string{
    "iPhone": "iPhone",
    "Apple Watch": "Apple Watch",
    "iPad": "iPad",        // 新产品类型
    "Mac": "Mac",          // 新产品类型
}
```

## 技术细节说明

### API 调用机制

#### 产品数据获取

1. **URL 格式**: `https://www.apple.com/{region}/shop/buy-{product}/{product-series}`
2. **数据源**: 页面中的 `<script type="application/json" id="metrics">` 标签
3. **解析方式**: 使用 `gjson` 库解析 JSON 数据

#### 门店数据获取

1. **API 端点**: `https://www.apple.com/{region}/shop/fulfillment-messages`
2. **必需参数**:
   - `fae=true`: 启用门店查询
   - `pl=true`: 启用位置查询
   - `mts.0=regular`: 消息类型
   - `parts.0={product_code}`: 产品代码
   - `location={location}`: 位置信息

3. **位置参数格式**:
   - **中国大陆**: 邮编 (如: 100000)
   - **香港**: Central (固定)
   - **日本**: 邮编 (如: 100-0001)
   - **新加坡**: Singapore (固定)
   - **美国**: 邮编 (如: 10001)

### 反爬虫策略

1. **随机 User-Agent**: 使用真实的浏览器 User-Agent
2. **随机 Referer**: 使用对应的产品页面作为 Referer
3. **请求延迟**: 在请求间添加随机延迟 (1-3秒)
4. **Accept-Language**: 根据地区设置合适的语言头

### 数据存储结构

#### 产品数据文件
- **路径**: `data/product_data_{area_code}.json`
- **结构**: 
```json
{
  "UpdateTime": "2025-01-01 12:00:00",
  "AreaCode": "cn",
  "Products": {
    "iPhone 16": [
      {
        "SKU": "MYEW3CH",
        "PartNumber": "MYEW3CH/A",
        "Name": "iPhone 16 128GB 白色"
      }
    ]
  }
}
```

#### 门店数据文件
- **路径**: `data/stores_{area_code}.json`
- **结构**:
```json
{
  "UpdateTime": "2025-01-01 12:00:00",
  "AreaCode": "cn",
  "Stores": [
    {
      "StoreNumber": "R448",
      "CityStoreName": "北京-王府井",
      "Province": "北京",
      "City": "北京",
      "District": ""
    }
  ]
}
```

## 测试验证

### 步骤 1: 创建测试脚本

```go
// test_new_region.go
package main

import (
    "apple-store-helper/services"
    "fmt"
    "log"
)

func main() {
    areaCode := "us" // 新地区代码
    
    fmt.Printf("=== 测试新地区: %s ===\n", areaCode)
    
    // 测试产品数据获取
    fmt.Println("\n1. 测试产品数据获取...")
    if err := services.UpdateProductDatabase(areaCode); err != nil {
        log.Printf("产品数据获取失败: %v", err)
    } else {
        fmt.Println("✓ 产品数据获取成功")
    }
    
    // 测试门店数据获取
    fmt.Println("\n2. 测试门店数据获取...")
    stores, err := services.FetchStoresForArea(areaCode, "10001") // 使用该地区的测试位置
    if err != nil {
        log.Printf("门店数据获取失败: %v", err)
    } else {
        fmt.Printf("✓ 成功获取 %d 个门店\n", len(stores))
        for i, store := range stores {
            if i < 3 {
                fmt.Printf("  %d. %s (%s)\n", i+1, store.CityStoreName, store.StoreNumber)
            }
        }
    }
    
    fmt.Println("\n=== 测试完成 ===")
}
```

### 步骤 2: 运行测试

```bash
go run test_new_region.go
```

### 步骤 3: 验证数据文件

检查生成的数据文件：

```bash
ls -la data/
cat data/product_data_us.json | jq '.Products | keys'
cat data/stores_us.json | jq '.Stores | length'
```

### 步骤 4: 测试UI功能

1. 编译并运行程序
2. 在地区选择器中选择新地区
3. 验证门店列表是否正确显示
4. 测试产品选择功能

## 常见问题

### Q1: 产品数据获取失败

**可能原因**:
- 产品页面URL不正确
- 页面结构发生变化
- 网络连接问题

**解决方案**:
1. 检查产品页面URL是否正确
2. 查看页面源码确认 `<script id="metrics">` 标签存在
3. 检查网络连接和代理设置

### Q2: 门店数据获取失败

**可能原因**:
- 产品代码无效
- 位置参数格式不正确
- API 端点不支持该地区

**解决方案**:
1. 使用 `curl` 测试 API 调用
2. 确认产品代码在该地区有效
3. 检查位置参数格式

### Q3: UI 中不显示新地区

**可能原因**:
- 地区定义不正确
- 地区选择器未更新

**解决方案**:
1. 检查 `model/area.go` 中的地区定义
2. 确认 `services.Area.ForOptions()` 包含新地区
3. 重新编译程序

### Q4: 门店筛选不工作

**可能原因**:
- 门店数据结构不正确
- 筛选逻辑有误

**解决方案**:
1. 检查门店数据的 `Province` 字段
2. 验证 `ByAreaAndProvinceForOptions` 函数逻辑
3. 查看控制台日志输出

## 调试技巧

### 1. 启用详细日志

在代码中添加日志输出：

```go
log.Printf("Fetching stores for area %s, location %s", areaCode, location)
log.Printf("API URL: %s", apiURL)
log.Printf("Response: %s", string(body))
```

### 2. 使用 curl 测试 API

```bash
curl -s "https://www.apple.com/us/shop/fulfillment-messages?fae=true&pl=true&mts.0=regular&parts.0=MYEW3LL/A&location=10001" \
  -H "User-Agent: Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36" \
  -H "Accept: application/json" | jq '.body.content.pickupMessage.stores'
```

### 3. 检查数据文件

```bash
# 查看产品数据
cat data/product_data_us.json | jq '.Products | keys'

# 查看门店数据
cat data/stores_us.json | jq '.Stores[0:3]'

# 统计门店数量
cat data/stores_us.json | jq '.Stores | length'
```

## 总结

按照本指南的步骤，任何人都可以成功适配新地区和新产品系列。关键是要：

1. **仔细检查 API 调用参数**
2. **验证产品代码的有效性**
3. **测试数据获取和解析**
4. **验证 UI 功能**

如果遇到问题，请参考常见问题部分或查看调试技巧。记住，每个地区的 Apple 官网可能有细微差异，需要根据实际情况调整参数。
