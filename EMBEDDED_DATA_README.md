# 嵌入数据说明

## 概述

本项目已成功将 data 文件嵌入到 Go 程序中，使用 Go 1.16+ 的 `embed` 包实现。这样程序就不再依赖外部的 JSON 文件，可以独立运行。

## 实现方式

### 1. 嵌入包结构

```
embedded/
├── embedded_data.go          # 嵌入数据定义
└── data/                     # 复制的数据文件
    ├── product/              # 产品数据
    │   ├── product_data_cn.json
    │   ├── product_data_hk.json
    │   ├── product_data_jp.json
    │   ├── product_data_sg.json
    │   ├── product_data_us.json
    │   ├── product_data_uk.json
    │   └── product_data_au.json
    └── store/                # 门店数据
        ├── store_cn.json
        ├── store_jp.json
        ├── store_us.json
        ├── store_uk.json
        └── store_au.json
```

### 2. 嵌入的数据

#### 产品数据 (7 个地区)
- **中国大陆 (cn)**: 4 个产品系列
- **香港 (hk)**: 6 个产品系列  
- **日本 (jp)**: 6 个产品系列
- **新加坡 (sg)**: 6 个产品系列
- **美国 (us)**: 6 个产品系列
- **英国 (uk)**: 6 个产品系列
- **澳大利亚 (au)**: 6 个产品系列

#### 门店数据 (5 个地区)
- **中国大陆 (cn)**: 49 个门店
- **日本 (jp)**: 11 个门店
- **美国 (us)**: 274 个门店
- **英国 (uk)**: 40 个门店
- **澳大利亚 (au)**: 22 个门店

### 3. 数据加载逻辑

程序会优先从嵌入的数据加载，如果嵌入数据不存在，则回退到文件系统：

```go
// 产品数据加载
func LoadProductData(areaCode string) (*ProductData, error) {
    // 首先尝试从嵌入的数据加载
    if data, exists := embedded.GetProductData(areaCode); exists {
        // 解析嵌入的数据
        // ...
    }
    
    // 如果嵌入数据不存在，回退到文件系统
    // ...
}

// 门店数据加载
func LoadStoreData(areaCode string) (*StoreData, error) {
    // 首先尝试从嵌入的数据加载
    if data, exists := embedded.GetStoreData(areaCode); exists {
        // 解析嵌入的数据
        // ...
    }
    
    // 如果嵌入数据不存在，回退到文件系统
    // ...
}
```

## 优势

### 1. 独立运行
- 程序不再依赖外部的 data 目录
- 可以打包成单个可执行文件分发
- 减少了文件丢失的风险

### 2. 性能提升
- 数据在编译时嵌入，启动时直接加载
- 避免了文件 I/O 操作
- 提高了程序启动速度

### 3. 部署简化
- 只需要分发一个可执行文件
- 不需要担心数据文件的路径问题
- 简化了安装和部署流程

## 数据更新

### 更新嵌入数据

1. **修改源数据文件**：
   ```bash
   # 修改原始数据文件
   vim data/product/product_data_cn.json
   ```

2. **复制到嵌入目录**：
   ```bash
   # 复制更新的文件到嵌入目录
   cp data/product/*.json embedded/data/product/
   cp data/store/store_*.json embedded/data/store/
   ```

3. **重新编译**：
   ```bash
   go build -o apple-store-helper .
   ```

### 添加新地区数据

1. **添加新的 embed 声明**：
   ```go
   //go:embed data/product/product_data_xx.json
   var ProductDataXX []byte
   
   //go:embed data/store/store_xx.json
   var StoreDataXX []byte
   ```

2. **更新映射表**：
   ```go
   var ProductDataMap = map[string][]byte{
       // ... 现有映射
       "xx": ProductDataXX,
   }
   
   var StoreDataMap = map[string][]byte{
       // ... 现有映射
       "xx": StoreDataXX,
   }
   ```

3. **复制数据文件并重新编译**

## 注意事项

### 1. 文件大小限制
- 嵌入的数据会增加可执行文件的大小
- 当前总嵌入数据约 100KB，影响较小
- 如果数据过大，可以考虑压缩或分片

### 2. 编译时间
- 嵌入大量数据会增加编译时间
- 建议在 CI/CD 中预编译

### 3. 数据更新
- 更新数据需要重新编译程序
- 对于频繁更新的数据，可以考虑保持文件系统方式

## 测试验证

使用以下命令验证嵌入的数据：

```bash
# 编译程序
go build -o apple-store-helper .

# 运行程序测试数据加载
./apple-store-helper
```

程序启动时会显示类似以下日志：
```
Loaded embedded product data for cn
Loaded embedded store data for cn: 49 stores
```

## 总结

通过使用 Go embed 包，我们成功将 data 文件嵌入到程序中，实现了：

- ✅ 独立运行：不依赖外部数据文件
- ✅ 性能提升：启动时直接加载嵌入数据
- ✅ 部署简化：单个可执行文件即可运行
- ✅ 向后兼容：支持文件系统回退机制

这种实现方式特别适合需要独立分发的桌面应用程序。
