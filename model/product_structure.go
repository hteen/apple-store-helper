package model

// ProductInfo 包含完整的产品信息
type ProductInfo struct {
	Model    string // 型号: iPhone Air, iPhone 17, iPhone 17 Pro, iPhone 17 Pro Max
	Capacity string // 容量: 256GB, 512GB, 1TB, 2TB
	Color    string // 颜色
	Code     string // 产品代码
	Type     string // 产品类型标识
}

// ModelConfig 定义每个型号的配置
type ModelConfig struct {
	Name       string
	Type       string
	Capacities []string
	Colors     map[string][]string // capacity -> colors
	Products   map[string]ProductInfo // "capacity_color" -> ProductInfo
}

// iPhone16Models 定义iPhone 16系列型号配置
var iPhone16Models = map[string]*ModelConfig{
	"iPhone 16": {
		Name: "iPhone 16",
		Type: "iphone16",
		Capacities: []string{"128GB", "256GB", "512GB"},
		Colors: map[string][]string{
			"128GB": {"黑色", "白色", "粉色", "深青色", "群青色"},
			"256GB": {"黑色", "白色", "粉色", "深青色", "群青色"},
			"512GB": {"黑色", "白色", "粉色", "深青色", "群青色"},
		},
		Products: map[string]ProductInfo{
			// 128GB
			"128GB_黑色":  {Model: "iPhone 16", Capacity: "128GB", Color: "黑色", Code: "MYEA3CH/A", Type: "iphone16"},
			"128GB_白色":  {Model: "iPhone 16", Capacity: "128GB", Color: "白色", Code: "MYEW3CH/A", Type: "iphone16"},
			"128GB_粉色":  {Model: "iPhone 16", Capacity: "128GB", Color: "粉色", Code: "MYEX3CH/A", Type: "iphone16"},
			"128GB_深青色": {Model: "iPhone 16", Capacity: "128GB", Color: "深青色", Code: "MYF03CH/A", Type: "iphone16"},
			"128GB_群青色": {Model: "iPhone 16", Capacity: "128GB", Color: "群青色", Code: "MYEK3CH/A", Type: "iphone16"},
			// 256GB
			"256GB_黑色":  {Model: "iPhone 16", Capacity: "256GB", Color: "黑色", Code: "MYF23CH/A", Type: "iphone16"},
			"256GB_白色":  {Model: "iPhone 16", Capacity: "256GB", Color: "白色", Code: "MYF53CH/A", Type: "iphone16"},
			"256GB_粉色":  {Model: "iPhone 16", Capacity: "256GB", Color: "粉色", Code: "MYFA3CH/A", Type: "iphone16"},
			"256GB_深青色": {Model: "iPhone 16", Capacity: "256GB", Color: "深青色", Code: "MYFE3CH/A", Type: "iphone16"},
			"256GB_群青色": {Model: "iPhone 16", Capacity: "256GB", Color: "群青色", Code: "MYF83CH/A", Type: "iphone16"},
			// 512GB
			"512GB_黑色":  {Model: "iPhone 16", Capacity: "512GB", Color: "黑色", Code: "MYFG3CH/A", Type: "iphone16"},
			"512GB_白色":  {Model: "iPhone 16", Capacity: "512GB", Color: "白色", Code: "MYFK3CH/A", Type: "iphone16"},
			"512GB_粉色":  {Model: "iPhone 16", Capacity: "512GB", Color: "粉色", Code: "MYFQ3CH/A", Type: "iphone16"},
			"512GB_深青色": {Model: "iPhone 16", Capacity: "512GB", Color: "深青色", Code: "MYFU3CH/A", Type: "iphone16"},
			"512GB_群青色": {Model: "iPhone 16", Capacity: "512GB", Color: "群青色", Code: "MYFN3CH/A", Type: "iphone16"},
		},
	},
	"iPhone 16 Plus": {
		Name: "iPhone 16 Plus",
		Type: "iphone16plus",
		Capacities: []string{"128GB", "256GB", "512GB"},
		Colors: map[string][]string{
			"128GB": {"黑色", "白色", "粉色", "深青色", "群青色"},
			"256GB": {"黑色", "白色", "粉色", "深青色", "群青色"},
			"512GB": {"黑色", "白色", "粉色", "深青色", "群青色"},
		},
		Products: map[string]ProductInfo{
			// 128GB
			"128GB_黑色":  {Model: "iPhone 16 Plus", Capacity: "128GB", Color: "黑色", Code: "MXU43CH/A", Type: "iphone16plus"},
			"128GB_白色":  {Model: "iPhone 16 Plus", Capacity: "128GB", Color: "白色", Code: "MXU73CH/A", Type: "iphone16plus"},
			"128GB_粉色":  {Model: "iPhone 16 Plus", Capacity: "128GB", Color: "粉色", Code: "MXUA3CH/A", Type: "iphone16plus"},
			"128GB_深青色": {Model: "iPhone 16 Plus", Capacity: "128GB", Color: "深青色", Code: "MXUE3CH/A", Type: "iphone16plus"},
			"128GB_群青色": {Model: "iPhone 16 Plus", Capacity: "128GB", Color: "群青色", Code: "MXUD3CH/A", Type: "iphone16plus"},
			// 256GB
			"256GB_黑色":  {Model: "iPhone 16 Plus", Capacity: "256GB", Color: "黑色", Code: "MXUG3CH/A", Type: "iphone16plus"},
			"256GB_白色":  {Model: "iPhone 16 Plus", Capacity: "256GB", Color: "白色", Code: "MXUH3CH/A", Type: "iphone16plus"},
			"256GB_粉色":  {Model: "iPhone 16 Plus", Capacity: "256GB", Color: "粉色", Code: "MXUL3CH/A", Type: "iphone16plus"},
			"256GB_深青色": {Model: "iPhone 16 Plus", Capacity: "256GB", Color: "深青色", Code: "MXUP3CH/A", Type: "iphone16plus"},
			"256GB_群青色": {Model: "iPhone 16 Plus", Capacity: "256GB", Color: "群青色", Code: "MXUJ3CH/A", Type: "iphone16plus"},
			// 512GB
			"512GB_黑色":  {Model: "iPhone 16 Plus", Capacity: "512GB", Color: "黑色", Code: "MXUR3CH/A", Type: "iphone16plus"},
			"512GB_白色":  {Model: "iPhone 16 Plus", Capacity: "512GB", Color: "白色", Code: "MXUT3CH/A", Type: "iphone16plus"},
			"512GB_粉色":  {Model: "iPhone 16 Plus", Capacity: "512GB", Color: "粉色", Code: "MXUX3CH/A", Type: "iphone16plus"},
			"512GB_深青色": {Model: "iPhone 16 Plus", Capacity: "512GB", Color: "深青色", Code: "MXV13CH/A", Type: "iphone16plus"},
			"512GB_群青色": {Model: "iPhone 16 Plus", Capacity: "512GB", Color: "群青色", Code: "MXUW3CH/A", Type: "iphone16plus"},
		},
	},
	"iPhone 16 Pro": {
		Name: "iPhone 16 Pro",
		Type: "iphone16pro",
		Capacities: []string{"128GB", "256GB", "512GB", "1TB"},
		Colors: map[string][]string{
			"128GB": {"黑色钛金属", "白色钛金属", "原色钛金属", "沙漠色钛金属"},
			"256GB": {"黑色钛金属", "白色钛金属", "原色钛金属", "沙漠色钛金属"},
			"512GB": {"黑色钛金属", "白色钛金属", "原色钛金属", "沙漠色钛金属"},
			"1TB":   {"黑色钛金属", "白色钛金属", "原色钛金属", "沙漠色钛金属"},
		},
		Products: map[string]ProductInfo{
			// 128GB
			"128GB_黑色钛金属": {Model: "iPhone 16 Pro", Capacity: "128GB", Color: "黑色钛金属", Code: "MYLP3CH/A", Type: "iphone16pro"},
			"128GB_白色钛金属": {Model: "iPhone 16 Pro", Capacity: "128GB", Color: "白色钛金属", Code: "MYLQ3CH/A", Type: "iphone16pro"},
			"128GB_原色钛金属": {Model: "iPhone 16 Pro", Capacity: "128GB", Color: "原色钛金属", Code: "MYLT3CH/A", Type: "iphone16pro"},
			"128GB_沙漠色钛金属": {Model: "iPhone 16 Pro", Capacity: "128GB", Color: "沙漠色钛金属", Code: "MYLR3CH/A", Type: "iphone16pro"},
			// 256GB
			"256GB_黑色钛金属": {Model: "iPhone 16 Pro", Capacity: "256GB", Color: "黑色钛金属", Code: "MYLW3CH/A", Type: "iphone16pro"},
			"256GB_白色钛金属": {Model: "iPhone 16 Pro", Capacity: "256GB", Color: "白色钛金属", Code: "MYLX3CH/A", Type: "iphone16pro"},
			"256GB_原色钛金属": {Model: "iPhone 16 Pro", Capacity: "256GB", Color: "原色钛金属", Code: "MYM03CH/A", Type: "iphone16pro"},
			"256GB_沙漠色钛金属": {Model: "iPhone 16 Pro", Capacity: "256GB", Color: "沙漠色钛金属", Code: "MYLY3CH/A", Type: "iphone16pro"},
			// 512GB
			"512GB_黑色钛金属": {Model: "iPhone 16 Pro", Capacity: "512GB", Color: "黑色钛金属", Code: "MYM13CH/A", Type: "iphone16pro"},
			"512GB_白色钛金属": {Model: "iPhone 16 Pro", Capacity: "512GB", Color: "白色钛金属", Code: "MYM23CH/A", Type: "iphone16pro"},
			"512GB_原色钛金属": {Model: "iPhone 16 Pro", Capacity: "512GB", Color: "原色钛金属", Code: "MYM53CH/A", Type: "iphone16pro"},
			"512GB_沙漠色钛金属": {Model: "iPhone 16 Pro", Capacity: "512GB", Color: "沙漠色钛金属", Code: "MYM33CH/A", Type: "iphone16pro"},
			// 1TB
			"1TB_黑色钛金属": {Model: "iPhone 16 Pro", Capacity: "1TB", Color: "黑色钛金属", Code: "MYM63CH/A", Type: "iphone16pro"},
			"1TB_白色钛金属": {Model: "iPhone 16 Pro", Capacity: "1TB", Color: "白色钛金属", Code: "MYM73CH/A", Type: "iphone16pro"},
			"1TB_原色钛金属": {Model: "iPhone 16 Pro", Capacity: "1TB", Color: "原色钛金属", Code: "MYMA3CH/A", Type: "iphone16pro"},
			"1TB_沙漠色钛金属": {Model: "iPhone 16 Pro", Capacity: "1TB", Color: "沙漠色钛金属", Code: "MYM83CH/A", Type: "iphone16pro"},
		},
	},
	"iPhone 16 Pro Max": {
		Name: "iPhone 16 Pro Max", 
		Type: "iphone16promax",
		Capacities: []string{"256GB", "512GB", "1TB"},
		Colors: map[string][]string{
			"256GB": {"黑色钛金属", "白色钛金属", "原色钛金属", "沙漠色钛金属"},
			"512GB": {"黑色钛金属", "白色钛金属", "原色钛金属", "沙漠色钛金属"},
			"1TB":   {"黑色钛金属", "白色钛金属", "原色钛金属", "沙漠色钛金属"},
		},
		Products: map[string]ProductInfo{
			// 256GB
			"256GB_黑色钛金属": {Model: "iPhone 16 Pro Max", Capacity: "256GB", Color: "黑色钛金属", Code: "MYTM3CH/A", Type: "iphone16promax"},
			"256GB_白色钛金属": {Model: "iPhone 16 Pro Max", Capacity: "256GB", Color: "白色钛金属", Code: "MYTN3CH/A", Type: "iphone16promax"},
			"256GB_原色钛金属": {Model: "iPhone 16 Pro Max", Capacity: "256GB", Color: "原色钛金属", Code: "MYTQ3CH/A", Type: "iphone16promax"},
			"256GB_沙漠色钛金属": {Model: "iPhone 16 Pro Max", Capacity: "256GB", Color: "沙漠色钛金属", Code: "MYTP3CH/A", Type: "iphone16promax"},
			// 512GB
			"512GB_黑色钛金属": {Model: "iPhone 16 Pro Max", Capacity: "512GB", Color: "黑色钛金属", Code: "MYTR3CH/A", Type: "iphone16promax"},
			"512GB_白色钛金属": {Model: "iPhone 16 Pro Max", Capacity: "512GB", Color: "白色钛金属", Code: "MYTT3CH/A", Type: "iphone16promax"},
			"512GB_原色钛金属": {Model: "iPhone 16 Pro Max", Capacity: "512GB", Color: "原色钛金属", Code: "MYTX3CH/A", Type: "iphone16promax"},
			"512GB_沙漠色钛金属": {Model: "iPhone 16 Pro Max", Capacity: "512GB", Color: "沙漠色钛金属", Code: "MYTW3CH/A", Type: "iphone16promax"},
			// 1TB
			"1TB_黑色钛金属": {Model: "iPhone 16 Pro Max", Capacity: "1TB", Color: "黑色钛金属", Code: "MYTY3CH/A", Type: "iphone16promax"},
			"1TB_白色钛金属": {Model: "iPhone 16 Pro Max", Capacity: "1TB", Color: "白色钛金属", Code: "MYU03CH/A", Type: "iphone16promax"},
			"1TB_原色钛金属": {Model: "iPhone 16 Pro Max", Capacity: "1TB", Color: "原色钛金属", Code: "MYU33CH/A", Type: "iphone16promax"},
			"1TB_沙漠色钛金属": {Model: "iPhone 16 Pro Max", Capacity: "1TB", Color: "沙漠色钛金属", Code: "MYU23CH/A", Type: "iphone16promax"},
		},
	},
}

// iPhone17Models 定义所有iPhone 17系列型号配置
var iPhone17Models = map[string]*ModelConfig{
	"iPhone Air": {
		Name: "iPhone Air",
		Type: "iphoneair",
		Capacities: []string{"256GB", "512GB", "1TB"},
		Colors: map[string][]string{
			"256GB": {"云白色", "天蓝色", "深空黑色", "浅金色"},
			"512GB": {"云白色", "天蓝色", "深空黑色", "浅金色"},
			"1TB":   {"云白色", "天蓝色", "深空黑色", "浅金色"},
		},
		Products: map[string]ProductInfo{
			"256GB_云白色":  {Model: "iPhone Air", Capacity: "256GB", Color: "云白色", Code: "MG334CH/A", Type: "iphoneair"},
			"256GB_天蓝色":  {Model: "iPhone Air", Capacity: "256GB", Color: "天蓝色", Code: "MG364CH/A", Type: "iphoneair"},
			"256GB_深空黑色": {Model: "iPhone Air", Capacity: "256GB", Color: "深空黑色", Code: "MG314CH/A", Type: "iphoneair"},
			"256GB_浅金色":  {Model: "iPhone Air", Capacity: "256GB", Color: "浅金色", Code: "MG344CH/A", Type: "iphoneair"},
			"512GB_云白色":  {Model: "iPhone Air", Capacity: "512GB", Color: "云白色", Code: "MG394CH/A", Type: "iphoneair"},
			"512GB_天蓝色":  {Model: "iPhone Air", Capacity: "512GB", Color: "天蓝色", Code: "MG3C4CH/A", Type: "iphoneair"},
			"512GB_深空黑色": {Model: "iPhone Air", Capacity: "512GB", Color: "深空黑色", Code: "MG374CH/A", Type: "iphoneair"},
			"512GB_浅金色":  {Model: "iPhone Air", Capacity: "512GB", Color: "浅金色", Code: "MG3A4CH/A", Type: "iphoneair"},
			"1TB_云白色":   {Model: "iPhone Air", Capacity: "1TB", Color: "云白色", Code: "MG3E4CH/A", Type: "iphoneair"},
			"1TB_天蓝色":   {Model: "iPhone Air", Capacity: "1TB", Color: "天蓝色", Code: "MG3G4CH/A", Type: "iphoneair"},
			"1TB_深空黑色":  {Model: "iPhone Air", Capacity: "1TB", Color: "深空黑色", Code: "MG3D4CH/A", Type: "iphoneair"},
			"1TB_浅金色":   {Model: "iPhone Air", Capacity: "1TB", Color: "浅金色", Code: "MG3F4CH/A", Type: "iphoneair"},
		},
	},
	"iPhone 17": {
		Name: "iPhone 17",
		Type: "iphone17",
		Capacities: []string{"256GB", "512GB"},
		Colors: map[string][]string{
			"256GB": {"黑色", "鼠尾草色", "白色", "薰衣草色", "薄雾蓝色"},
			"512GB": {"黑色", "鼠尾草色", "白色", "薰衣草色", "薄雾蓝色"},
		},
		Products: map[string]ProductInfo{
			"256GB_黑色":    {Model: "iPhone 17", Capacity: "256GB", Color: "黑色", Code: "MG6W4CH/A", Type: "iphone17"},
			"256GB_鼠尾草色":  {Model: "iPhone 17", Capacity: "256GB", Color: "鼠尾草色", Code: "MG714CH/A", Type: "iphone17"},
			"256GB_白色":    {Model: "iPhone 17", Capacity: "256GB", Color: "白色", Code: "MG6X4CH/A", Type: "iphone17"},
			"256GB_薰衣草色":  {Model: "iPhone 17", Capacity: "256GB", Color: "薰衣草色", Code: "MG704CH/A", Type: "iphone17"},
			"256GB_薄雾蓝色":  {Model: "iPhone 17", Capacity: "256GB", Color: "薄雾蓝色", Code: "MG6Y4CH/A", Type: "iphone17"},
			"512GB_黑色":    {Model: "iPhone 17", Capacity: "512GB", Color: "黑色", Code: "MG724CH/A", Type: "iphone17"},
			"512GB_鼠尾草色":  {Model: "iPhone 17", Capacity: "512GB", Color: "鼠尾草色", Code: "MG764CH/A", Type: "iphone17"},
			"512GB_白色":    {Model: "iPhone 17", Capacity: "512GB", Color: "白色", Code: "MG734CH/A", Type: "iphone17"},
			"512GB_薰衣草色":  {Model: "iPhone 17", Capacity: "512GB", Color: "薰衣草色", Code: "MG754CH/A", Type: "iphone17"},
			"512GB_薄雾蓝色":  {Model: "iPhone 17", Capacity: "512GB", Color: "薄雾蓝色", Code: "MG744CH/A", Type: "iphone17"},
		},
	},
	"iPhone 17 Pro": {
		Name: "iPhone 17 Pro",
		Type: "iphone17pro",
		Capacities: []string{"256GB", "512GB", "1TB", "2TB"},
		Colors: map[string][]string{
			"256GB": {"宇宙橙色", "深蓝色", "银色"},
			"512GB": {"宇宙橙色", "深蓝色", "银色"},
			"1TB":   {"宇宙橙色", "深蓝色", "银色"},
			"2TB":   {"银色"},
		},
		Products: map[string]ProductInfo{
			"256GB_宇宙橙色": {Model: "iPhone 17 Pro", Capacity: "256GB", Color: "宇宙橙色", Code: "MG074CH/A", Type: "iphone17pro"},
			"256GB_深蓝色":  {Model: "iPhone 17 Pro", Capacity: "256GB", Color: "深蓝色", Code: "MG084CH/A", Type: "iphone17pro"},
			"256GB_银色":   {Model: "iPhone 17 Pro", Capacity: "256GB", Color: "银色", Code: "MG094CH/A", Type: "iphone17pro"},
			"512GB_宇宙橙色": {Model: "iPhone 17 Pro", Capacity: "512GB", Color: "宇宙橙色", Code: "MG0C4CH/A", Type: "iphone17pro"},
			"512GB_深蓝色":  {Model: "iPhone 17 Pro", Capacity: "512GB", Color: "深蓝色", Code: "MG0D4CH/A", Type: "iphone17pro"},
			"512GB_银色":   {Model: "iPhone 17 Pro", Capacity: "512GB", Color: "银色", Code: "MG0G4CH/A", Type: "iphone17pro"},
			"1TB_宇宙橙色":  {Model: "iPhone 17 Pro", Capacity: "1TB", Color: "宇宙橙色", Code: "MG0H4CH/A", Type: "iphone17pro"},
			"1TB_深蓝色":   {Model: "iPhone 17 Pro", Capacity: "1TB", Color: "深蓝色", Code: "MG0J4CH/A", Type: "iphone17pro"},
			"1TB_银色":    {Model: "iPhone 17 Pro", Capacity: "1TB", Color: "银色", Code: "MG0K4CH/A", Type: "iphone17pro"},
			"2TB_银色":    {Model: "iPhone 17 Pro", Capacity: "2TB", Color: "银色", Code: "MG0L4CH/A", Type: "iphone17pro"},
		},
	},
	"iPhone 17 Pro Max": {
		Name: "iPhone 17 Pro Max",
		Type: "iphone17promax",
		Capacities: []string{"256GB", "512GB", "1TB", "2TB"},
		Colors: map[string][]string{
			"256GB": {"宇宙橙色", "深蓝色", "银色"},
			"512GB": {"宇宙橙色", "深蓝色", "银色"},
			"1TB":   {"宇宙橙色", "深蓝色", "银色"},
			"2TB":   {"银色", "深蓝色", "宇宙橙色"},
		},
		Products: map[string]ProductInfo{
			"256GB_宇宙橙色": {Model: "iPhone 17 Pro Max", Capacity: "256GB", Color: "宇宙橙色", Code: "MG0M4CH/A", Type: "iphone17promax"},
			"256GB_深蓝色":  {Model: "iPhone 17 Pro Max", Capacity: "256GB", Color: "深蓝色", Code: "MG0N4CH/A", Type: "iphone17promax"},
			"256GB_银色":   {Model: "iPhone 17 Pro Max", Capacity: "256GB", Color: "银色", Code: "MG0P4CH/A", Type: "iphone17promax"},
			"512GB_宇宙橙色": {Model: "iPhone 17 Pro Max", Capacity: "512GB", Color: "宇宙橙色", Code: "MG0R4CH/A", Type: "iphone17promax"},
			"512GB_深蓝色":  {Model: "iPhone 17 Pro Max", Capacity: "512GB", Color: "深蓝色", Code: "MG0T4CH/A", Type: "iphone17promax"},
			"512GB_银色":   {Model: "iPhone 17 Pro Max", Capacity: "512GB", Color: "银色", Code: "MG0U4CH/A", Type: "iphone17promax"},
			"1TB_宇宙橙色":  {Model: "iPhone 17 Pro Max", Capacity: "1TB", Color: "宇宙橙色", Code: "MG0V4CH/A", Type: "iphone17promax"},
			"1TB_深蓝色":   {Model: "iPhone 17 Pro Max", Capacity: "1TB", Color: "深蓝色", Code: "MG0W4CH/A", Type: "iphone17promax"},
			"1TB_银色":    {Model: "iPhone 17 Pro Max", Capacity: "1TB", Color: "银色", Code: "MG0X4CH/A", Type: "iphone17promax"},
			"2TB_银色":    {Model: "iPhone 17 Pro Max", Capacity: "2TB", Color: "银色", Code: "MG0Y4CH/A", Type: "iphone17promax"},
			"2TB_深蓝色":   {Model: "iPhone 17 Pro Max", Capacity: "2TB", Color: "深蓝色", Code: "MG104CH/A", Type: "iphone17promax"},
			"2TB_宇宙橙色":  {Model: "iPhone 17 Pro Max", Capacity: "2TB", Color: "宇宙橙色", Code: "MG114CH/A", Type: "iphone17promax"},
		},
	},
}

// AppleWatchModels 定义Apple Watch系列型号配置
var AppleWatchModels = map[string]*ModelConfig{
	"Apple Watch SE3": {
		Name: "Apple Watch SE3",
		Type: "watchse3",
		Capacities: []string{"40mm", "44mm"},
		Colors: map[string][]string{
			"40mm": {"银色", "星光色", "午夜色"},
			"44mm": {"银色", "星光色", "午夜色"},
		},
		Products: map[string]ProductInfo{
			// 40mm
			"40mm_银色":  {Model: "Apple Watch SE3", Capacity: "40mm", Color: "银色", Code: "MRE43CH/A", Type: "watchse3"},
			"40mm_星光色": {Model: "Apple Watch SE3", Capacity: "40mm", Color: "星光色", Code: "MRE53CH/A", Type: "watchse3"},
			"40mm_午夜色": {Model: "Apple Watch SE3", Capacity: "40mm", Color: "午夜色", Code: "MRE63CH/A", Type: "watchse3"},
			// 44mm
			"44mm_银色":  {Model: "Apple Watch SE3", Capacity: "44mm", Color: "银色", Code: "MRE73CH/A", Type: "watchse3"},
			"44mm_星光色": {Model: "Apple Watch SE3", Capacity: "44mm", Color: "星光色", Code: "MRE83CH/A", Type: "watchse3"},
			"44mm_午夜色": {Model: "Apple Watch SE3", Capacity: "44mm", Color: "午夜色", Code: "MRE93CH/A", Type: "watchse3"},
		},
	},
	"Apple Watch S11": {
		Name: "Apple Watch S11",
		Type: "watchs11",
		Capacities: []string{"42mm", "46mm"},
		Colors: map[string][]string{
			"42mm": {"银色", "星光色", "午夜色", "玫瑰金色", "喷射黑色"},
			"46mm": {"银色", "星光色", "午夜色", "玫瑰金色", "喷射黑色"},
		},
		Products: map[string]ProductInfo{
			// 42mm
			"42mm_银色":   {Model: "Apple Watch S11", Capacity: "42mm", Color: "银色", Code: "MREA3CH/A", Type: "watchs11"},
			"42mm_星光色":  {Model: "Apple Watch S11", Capacity: "42mm", Color: "星光色", Code: "MREB3CH/A", Type: "watchs11"},
			"42mm_午夜色":  {Model: "Apple Watch S11", Capacity: "42mm", Color: "午夜色", Code: "MREC3CH/A", Type: "watchs11"},
			"42mm_玫瑰金色": {Model: "Apple Watch S11", Capacity: "42mm", Color: "玫瑰金色", Code: "MRED3CH/A", Type: "watchs11"},
			"42mm_喷射黑色": {Model: "Apple Watch S11", Capacity: "42mm", Color: "喷射黑色", Code: "MREE3CH/A", Type: "watchs11"},
			// 46mm
			"46mm_银色":   {Model: "Apple Watch S11", Capacity: "46mm", Color: "银色", Code: "MREF3CH/A", Type: "watchs11"},
			"46mm_星光色":  {Model: "Apple Watch S11", Capacity: "46mm", Color: "星光色", Code: "MREG3CH/A", Type: "watchs11"},
			"46mm_午夜色":  {Model: "Apple Watch S11", Capacity: "46mm", Color: "午夜色", Code: "MREH3CH/A", Type: "watchs11"},
			"46mm_玫瑰金色": {Model: "Apple Watch S11", Capacity: "46mm", Color: "玫瑰金色", Code: "MREJ3CH/A", Type: "watchs11"},
			"46mm_喷射黑色": {Model: "Apple Watch S11", Capacity: "46mm", Color: "喷射黑色", Code: "MREK3CH/A", Type: "watchs11"},
		},
	},
	"Apple Watch Ultra3": {
		Name: "Apple Watch Ultra3",
		Type: "watchultra3",
		Capacities: []string{"49mm"},
		Colors: map[string][]string{
			"49mm": {"原色钛金属", "喷射黑色钛金属"},
		},
		Products: map[string]ProductInfo{
			"49mm_原色钛金属":   {Model: "Apple Watch Ultra3", Capacity: "49mm", Color: "原色钛金属", Code: "MREL3CH/A", Type: "watchultra3"},
			"49mm_喷射黑色钛金属": {Model: "Apple Watch Ultra3", Capacity: "49mm", Color: "喷射黑色钛金属", Code: "MREM3CH/A", Type: "watchultra3"},
		},
	},
}

// AllModels 合并所有型号配置
var AllModels = map[string]*ModelConfig{}

// InitModels 初始化所有型号
func init() {
	// 合并iPhone 16系列
	for k, v := range iPhone16Models {
		AllModels[k] = v
	}
	// 合并iPhone 17系列
	for k, v := range iPhone17Models {
		AllModels[k] = v
	}
	// 合并Apple Watch系列
	for k, v := range AppleWatchModels {
		AllModels[k] = v
	}
}

// GetModelNames 获取所有型号名称
func GetModelNames() []string {
	return []string{
		// iPhone 16系列
		"iPhone 16", "iPhone 16 Plus", "iPhone 16 Pro", "iPhone 16 Pro Max",
		// iPhone 17系列
		"iPhone Air", "iPhone 17", "iPhone 17 Pro", "iPhone 17 Pro Max",
		// Apple Watch系列
		"Apple Watch SE3", "Apple Watch S11", "Apple Watch Ultra3",
	}
}

// GetCapacitiesByModel 根据型号获取可用容量
func GetCapacitiesByModel(model string) []string {
	if config, exists := AllModels[model]; exists {
		return config.Capacities
	}
	return []string{}
}

// GetColorsByModelAndCapacity 根据型号和容量获取可用颜色
func GetColorsByModelAndCapacity(model, capacity string) []string {
	if config, exists := AllModels[model]; exists {
		if colors, exists := config.Colors[capacity]; exists {
			return colors
		}
	}
	return []string{}
}

// GetProductInfo 根据型号、容量和颜色获取产品信息
func GetProductInfo(model, capacity, color string) *ProductInfo {
	if config, exists := AllModels[model]; exists {
		key := capacity + "_" + color
		if product, exists := config.Products[key]; exists {
			return &product
		}
	}
	return nil
}