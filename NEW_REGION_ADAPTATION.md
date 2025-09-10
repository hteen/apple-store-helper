# 新地区适配指南

## 概述
本文档描述了如何为 Apple Store Helper 添加新地区支持的完整流程。

## 适配步骤

### 1. 添加地区定义

在 `model/area.go` 中添加新地区：

```go
var Areas = []Area{
    // ... 现有地区
    {
        Title:        "地区名称",      // 显示名称
        Locale:       "语言_地区",     // 如 en_US
        ShortCode:    "地区代码",      // 如 us
        ProductsJson: iPhone17ProductsJson,
    },
}
```

### 2. 获取产品数据

#### 2.1 配置产品数据获取器

在 `services/product_fetcher.go` 中添加地区支持：

```go
func FetchProductsFromApple(areaCode string) (*ProductData, error) {
    // 添加新地区的URL
    urls := map[string]string{
        "cn": "https://www.apple.com.cn/shop/buy-iphone/iphone-16",
        "hk": "https://www.apple.com/hk/shop/buy-iphone/iphone-16",
        "jp": "https://www.apple.com/jp/shop/buy-iphone/iphone-16",
        // 新地区：
        "us": "https://www.apple.com/shop/buy-iphone/iphone-16",
    }
}
```

#### 2.2 运行产品数据获取

产品数据会自动保存到 `data/product_data_[地区代码].json`

### 3. 配置门店数据获取

#### 3.1 更新门店获取服务

在 `services/store_fetcher.go` 中添加地区支持：

```go
func FetchStoresForArea(areaCode string, location string) ([]model.Store, error) {
    switch areaCode {
    case "新地区代码":
        // 配置API URL
        baseURL := "https://www.apple.com/[地区]/shop/address-lookup"
        // 设置查询参数
    }
}
```

#### 3.2 配置请求头

```go
func getAcceptLanguage(areaCode string) string {
    switch areaCode {
    case "新地区代码":
        return "语言代码"
    }
}

func getReferer(areaCode string) string {
    switch areaCode {
    case "新地区代码":
        return "https://www.apple.com/[地区]/shop/buy-iphone/iphone-16"
    }
}
```

### 4. 配置库存查询

在 `services/listen.go` 中更新 `monitorInventory` 函数：

#### 4.1 配置查询参数

```go
// 根据地区设置location参数
if s.Area.ShortCode == "新地区代码" {
    // 设置location格式
    // 例如：邮编、城市名或其他地区特定格式
}
```

#### 4.2 配置API域名

```go
// 根据地区选择域名
domain := "www.apple.com"
if s.Area.ShortCode == "cn" {
    domain = "www.apple.com.cn"
} else if s.Area.ShortCode == "新地区代码" {
    domain = "www.apple.com/[地区路径]"
}
```

### 5. 静态门店数据（可选）

如果需要提供静态门店数据作为备份，创建 `model/[地区]_stores.go`：

```go
var [地区]Stores = map[string][]Store{
    "语言_地区": {
        {
            StoreNumber:   "R###",
            CityStoreName: "城市-门店名",
            Province:      "省/州",
            City:          "城市",
            District:      "区/邮编",
        },
    },
}
```

然后在 `model/global_stores.go` 中注册：

```go
var GlobalStores = map[string][]Store{
    // ... 现有地区
    "语言_地区": [地区]Stores["语言_地区"],
}
```

### 6. UI适配

根据地区特性，可能需要在 `main.go` 中进行UI调整：

- **需要邮编的地区**（如日本）：添加邮编输入框
- **基于位置的地区**（如中国）：使用省市区选择器
- **简单地区**（如香港）：使用固定location值

### 7. 测试步骤

1. **获取数据**：
   ```bash
   # 程序会自动获取，或点击"更新数据"按钮
   ```

2. **验证产品数据**：
   检查 `data/product_data_[地区代码].json` 文件内容

3. **验证门店数据**：
   检查 `data/stores_[地区代码].json` 文件内容

4. **测试库存查询**：
   运行程序并验证库存监控功能

## 数据文件结构

### 产品数据 (product_data_[地区].json)

```json
{
  "update_time": "2025-09-10 18:20:57",
  "area_code": "地区代码",
  "products": {
    "产品系列": [
      {
        "Model": "型号",
        "Capacity": "容量",
        "Color": "颜色",
        "Code": "产品代码",
        "Type": "产品类型"
      }
    ]
  }
}
```

### 门店数据 (stores_[地区].json)

```json
{
  "update_time": "2025-09-10 18:20:57",
  "area_code": "地区代码",
  "stores": [
    {
      "StoreNumber": "R###",
      "CityStoreName": "门店名称",
      "Province": "省/州",
      "City": "城市",
      "District": "区/邮编"
    }
  ]
}
```

## 注意事项

1. **Unicode空格处理**：产品数据解析时需要规范化各种空格字符
2. **动态请求头**：使用随机User-Agent和动态headers避免被识别
3. **缓存管理**：定期清理HTTP客户端缓存
4. **库存字段**：使用 `pickupDisplay` 字段判断库存状态
   - `"available"` = 有货
   - `"ineligible"` = 无货

## 常见问题

### Q: 如何找到正确的API端点？
A: 打开浏览器开发者工具，访问Apple官网选择产品和门店，观察网络请求

### Q: 产品代码格式不正确？
A: 检查HTML中的metrics JSON，确保正确解析产品信息

### Q: 门店列表为空？
A: 验证location参数格式是否符合该地区要求