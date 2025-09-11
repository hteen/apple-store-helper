# 抢你妹 - Apple 产品库存监控工具

<div align="center">

![Version](https://img.shields.io/badge/version-1.6.2-blue.svg)
![Platform](https://img.shields.io/badge/platform-macOS-lightgrey.svg)
![Go](https://img.shields.io/badge/Go-1.17+-00ADD8.svg)
![License](https://img.shields.io/badge/license-MIT-green.svg)

**实时监控 Apple Store 库存，第一时间抢购心仪产品！**

[下载地址](#下载安装) | [使用教程](#使用教程) | [Bark 推送](#bark-推送设置) | [问题反馈](#问题反馈)

</div>

## ✨ 功能特色

### 🚀 核心功能
- **实时监控**：自动扫描 Apple Store 库存状态
- **多地区支持**：支持中国大陆、香港、日本、新加坡、美国、英国、澳大利亚
- **智能提醒**：Bark 推送 + 声音提醒，不错过任何机会
- **一键下单**：推送直达苹果官网购物车
- **独立运行**：无需外部数据文件，单文件即可运行

### 📱 支持的产品
- **iPhone 系列**：iPhone 16、iPhone 16 Plus、iPhone 16 Pro、iPhone 16 Pro Max
- **iPad 系列**：iPad Air、iPad Pro、iPad mini
- **Apple Watch**：Apple Watch Series 10、Apple Watch Ultra 3
- **Mac 系列**：MacBook Air、MacBook Pro、iMac、Mac Studio
- **其他产品**：AirPods、Apple Vision Pro 等

### 🌍 支持的地区
| 地区 | 产品数据 | 门店数据 | 状态 |
|------|----------|----------|------|
| 🇨🇳 中国大陆 | ✅ 4 个系列 | ✅ 49 个门店 | 完全支持 |
| 🇭🇰 香港 | ✅ 6 个系列 | ✅ 6 个门店 | 完全支持 |
| 🇯🇵 日本 | ✅ 6 个系列 | ✅ 11 个门店 | 完全支持 |
| 🇸🇬 新加坡 | ✅ 6 个系列 | ❌ 门店数据 | 部分支持 |
| 🇺🇸 美国 | ✅ 6 个系列 | ✅ 274 个门店 | 完全支持 |
| 🇬🇧 英国 | ✅ 6 个系列 | ✅ 40 个门店 | 完全支持 |
| 🇦🇺 澳大利亚 | ✅ 6 个系列 | ✅ 22 个门店 | 完全支持 |

## 📥 下载安装

### 系统要求
- **操作系统**：macOS 10.15+ (Catalina 或更高版本)
- **架构**：Intel x64 / Apple Silicon (M1/M2/M3)
- **内存**：至少 4GB RAM
- **网络**：稳定的互联网连接

### 下载方式

#### 方式一：直接下载（推荐）
1. 访问 [Releases 页面](https://github.com/your-repo/apple-store-helper/releases)
2. 下载最新版本的 `apple-store-helper-macOS.zip`
3. 解压后双击 `apple-store-helper` 运行

#### 方式二：从源码编译
```bash
# 克隆仓库
git clone https://github.com/your-repo/apple-store-helper.git
cd apple-store-helper

# 编译程序
go build -o apple-store-helper .

# 运行程序
./apple-store-helper
```

## 🚀 使用教程

### 首次使用

1. **启动程序**
   - 双击 `apple-store-helper` 启动
   - 首次启动会自动加载数据，请耐心等待

2. **选择地区**
   - 在「地区」下拉框中选择目标地区
   - 程序会自动加载该地区的产品和门店数据

3. **选择产品**
   - 在「型号」下拉框中选择产品型号
   - 在「容量/尺寸」下拉框中选择规格
   - 在「颜色」下拉框中选择颜色

4. **选择门店**
   - 中国大陆：选择省份 → 城市 → 门店
   - 其他地区：选择具体门店

5. **设置提醒**
   - 配置 Bark 推送（见下方教程）
   - 选择提醒方式（通知推送/持续响铃）

6. **开始监控**
   - 点击「开始监控」按钮
   - 程序将自动扫描库存状态

### 高级功能

#### 批量监控
- 可以同时监控多个产品
- 支持不同地区、不同门店的监控
- 每个监控任务独立运行

#### 自定义扫描间隔
- 默认扫描间隔：30 秒
- 可根据网络情况调整
- 建议不要设置过短，避免被限制

#### 数据更新
- 程序内置最新数据，无需手动更新
- 支持手动刷新产品数据
- 自动检测数据更新

## 📱 Bark 推送设置

Bark 是一款优秀的 iOS 推送工具，让你第一时间收到库存提醒！

### 步骤 1：下载 Bark

从 App Store 下载 Bark：

🔗 **下载链接**：https://apps.apple.com/my/app/bark-custom-notifications/id1403753865

![Bark App Store](https://ibeta-vue.oss-cn-hangzhou.aliyuncs.com/undefined20250910234709.png)

### 步骤 2：注册设备

1. 打开 Bark 应用
2. 点击「注册设备」按钮
3. 系统会自动生成你的推送地址

![注册设备](https://ibeta-vue.oss-cn-hangzhou.aliyuncs.com/undefined20250910234832.png)

### 步骤 3：配置推送

1. 复制 Bark 生成的推送地址
2. 在「抢你妹」程序中粘贴到「Bark 推送」输入框
3. 选择提醒方式：
   - **通知推送**：静默推送，适合办公环境
   - **持续响铃**：持续响铃提醒，确保不错过

![配置推送](https://ibeta-vue.oss-cn-hangzhou.aliyuncs.com/undefined3a01fa675551970c68a81f9858b97637.png)

### 步骤 4：开始监控

1. 在「抢你妹」中点击「开始监控」
2. 程序开始扫描库存
3. 有库存时会立即推送通知

### 步骤 5：快速下单

1. 点击 iPhone 上的推送通知
2. 自动跳转到苹果官网购物车
3. 完成下单购买

## ⚙️ 配置说明

### 程序设置

| 设置项 | 说明 | 默认值 |
|--------|------|--------|
| 扫描间隔 | 库存检查频率 | 30 秒 |
| 超时时间 | 网络请求超时 | 10 秒 |
| 重试次数 | 失败后重试次数 | 3 次 |
| 声音提醒 | 是否播放提醒音 | 开启 |

### 推送设置

| 设置项 | 说明 | 推荐值 |
|--------|------|--------|
| Bark 推送 | 推送服务地址 | 必填 |
| 提醒方式 | 通知/响铃 | 根据环境选择 |
| 推送频率 | 同一产品推送间隔 | 5 分钟 |

## 🔧 故障排除

### 常见问题

#### 1. 程序无法启动
**问题**：双击程序无反应或报错
**解决方案**：
- 检查系统版本是否满足要求
- 尝试在终端中运行：`./apple-store-helper`
- 检查文件权限：`chmod +x apple-store-helper`

#### 2. 无法加载数据
**问题**：程序启动后显示"数据加载失败"
**解决方案**：
- 检查网络连接
- 重启程序
- 检查防火墙设置

#### 3. 推送不工作
**问题**：设置了 Bark 推送但没有收到通知
**解决方案**：
- 检查 Bark 推送地址是否正确
- 确认 iPhone 网络连接正常
- 检查 Bark 应用是否正常运行

#### 4. 监控无结果
**问题**：程序运行但没有检测到库存
**解决方案**：
- 确认产品型号和规格选择正确
- 检查门店选择是否正确
- 尝试更换门店或地区

### 日志查看

程序运行时会在控制台输出日志信息：
- `INFO`：正常信息
- `WARN`：警告信息
- `ERROR`：错误信息

### 性能优化

1. **减少监控任务**：同时监控的产品越少，扫描速度越快
2. **选择合适的门店**：选择距离较近或库存较多的门店
3. **调整扫描间隔**：网络较慢时可适当增加间隔时间

## 📊 数据统计

### 内置数据
- **产品数据**：7 个地区，40+ 个产品系列
- **门店数据**：6 个地区，400+ 个门店
- **总数据量**：约 100KB
- **更新频率**：随程序版本更新

### 支持的产品系列
- iPhone 16 系列（4 款）
- iPhone 16 Plus 系列（4 款）
- iPhone 16 Pro 系列（4 款）
- iPhone 16 Pro Max 系列（4 款）
- iPad Air 系列（2 款）
- iPad Pro 系列（4 款）
- Apple Watch Series 10（3 款）
- Apple Watch Ultra 3（2 款）
- MacBook Air（2 款）
- MacBook Pro（4 款）
- 更多产品持续更新中...

## 🤝 贡献指南

### 报告问题
1. 在 [Issues](https://github.com/your-repo/apple-store-helper/issues) 中搜索是否已有相同问题
2. 如果没有，请创建新的 Issue
3. 详细描述问题现象和复现步骤

### 功能建议
1. 在 [Discussions](https://github.com/your-repo/apple-store-helper/discussions) 中提出建议
2. 详细描述功能需求和预期效果
3. 参与讨论和投票

### 代码贡献
1. Fork 本仓库
2. 创建功能分支
3. 提交代码并创建 Pull Request
4. 等待代码审查

## 📄 许可证

本项目采用 [MIT 许可证](LICENSE)。

## 🙏 致谢

- 感谢所有贡献者的支持
- 感谢 Bark 团队提供的推送服务
- 感谢 Apple 提供的产品数据接口

## 📞 联系我们

- **GitHub**：https://github.com/your-repo/apple-store-helper
- **邮箱**：your-email@example.com
- **QQ 群**：123456789
- **微信群**：扫描二维码加入

---

<div align="center">

**⭐ 如果这个项目对你有帮助，请给个 Star！**

Made with ❤️ by [Your Name]

</div>