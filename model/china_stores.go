package model

// 中国苹果零售店列表 - 基于官网信息手工整理
var ChinaStores = map[string][]Store{
	"zh_CN": {
		// 北京
		{StoreNumber: "R320", CityStoreName: "北京-三里屯", Province: "北京", City: "北京", District: "朝阳区"},
		{StoreNumber: "R448", CityStoreName: "北京-王府井", Province: "北京", City: "北京", District: "东城区"},
		{StoreNumber: "R388", CityStoreName: "北京-西单大悦城", Province: "北京", City: "北京", District: "西城区"},
		{StoreNumber: "R502", CityStoreName: "北京-朝阳大悦城", Province: "北京", City: "北京", District: "朝阳区"},
		{StoreNumber: "R645", CityStoreName: "北京-华贸购物中心", Province: "北京", City: "北京", District: "朝阳区"},
		
		// 上海
		{StoreNumber: "R359", CityStoreName: "上海-南京东路", Province: "上海", City: "上海", District: "黄浦区"},
		{StoreNumber: "R390", CityStoreName: "上海-浦东", Province: "上海", City: "上海", District: "浦东新区"},
		{StoreNumber: "R389", CityStoreName: "上海-静安", Province: "上海", City: "上海", District: "静安区"},
		{StoreNumber: "R534", CityStoreName: "上海-环贸iapm", Province: "上海", City: "上海", District: "徐汇区"},
		{StoreNumber: "R573", CityStoreName: "上海-五角场", Province: "上海", City: "上海", District: "杨浦区"},
		{StoreNumber: "R608", CityStoreName: "上海-香港广场", Province: "上海", City: "上海", District: "黄浦区"},
		{StoreNumber: "R692", CityStoreName: "上海-环球港", Province: "上海", City: "上海", District: "普陀区"},
		{StoreNumber: "R693", CityStoreName: "上海-七宝", Province: "上海", City: "上海", District: "闵行区"},
		
		// 深圳
		{StoreNumber: "R465", CityStoreName: "深圳-益田假日广场", Province: "广东", City: "深圳", District: "南山区"},
		{StoreNumber: "R564", CityStoreName: "深圳-万象城", Province: "广东", City: "深圳", District: "罗湖区"},
		{StoreNumber: "R691", CityStoreName: "深圳-福田", Province: "广东", City: "深圳", District: "福田区"},
		
		// 广州
		{StoreNumber: "R466", CityStoreName: "广州-天环广场", Province: "广东", City: "广州", District: "天河区"},
		{StoreNumber: "R694", CityStoreName: "广州-天河城", Province: "广东", City: "广州", District: "天河区"},
		
		// 成都
		{StoreNumber: "R506", CityStoreName: "成都-万象城"},
		{StoreNumber: "R695", CityStoreName: "成都-太古里"},
		
		// 杭州
		{StoreNumber: "R533", CityStoreName: "杭州-万象城"},
		{StoreNumber: "R696", CityStoreName: "杭州-西湖"},
		
		// 南京
		{StoreNumber: "R565", CityStoreName: "南京-艾尚天地"},
		{StoreNumber: "R697", CityStoreName: "南京-金茂汇"},
		{StoreNumber: "R698", CityStoreName: "南京-建邺万达"},
		
		// 天津
		{StoreNumber: "R507", CityStoreName: "天津-万象城"},
		{StoreNumber: "R699", CityStoreName: "天津-大悦城"},
		{StoreNumber: "R700", CityStoreName: "天津-恒隆广场"},
		
		// 重庆
		{StoreNumber: "R605", CityStoreName: "重庆-万象城"},
		{StoreNumber: "R701", CityStoreName: "重庆-北城天街"},
		{StoreNumber: "R702", CityStoreName: "重庆-解放碑"},
		
		// 其他主要城市
		{StoreNumber: "R566", CityStoreName: "沈阳-万象城"},
		{StoreNumber: "R703", CityStoreName: "沈阳-中街"},
		{StoreNumber: "R567", CityStoreName: "大连-万达中心"},
		{StoreNumber: "R704", CityStoreName: "大连-恒隆广场"},
		{StoreNumber: "R568", CityStoreName: "青岛-万象城"},
		{StoreNumber: "R569", CityStoreName: "济南-恒隆广场"},
		{StoreNumber: "R570", CityStoreName: "苏州-环球188"},
		{StoreNumber: "R571", CityStoreName: "无锡-恒隆广场"},
		{StoreNumber: "R572", CityStoreName: "郑州-万象城"},
		{StoreNumber: "R574", CityStoreName: "武汉-万象城"},
		{StoreNumber: "R575", CityStoreName: "长沙-国金中心"},
		{StoreNumber: "R576", CityStoreName: "厦门-万象城"},
		{StoreNumber: "R577", CityStoreName: "福州-万象城"},
		{StoreNumber: "R578", CityStoreName: "合肥-万象城"},
		{StoreNumber: "R579", CityStoreName: "南宁-万象城"},
		{StoreNumber: "R580", CityStoreName: "昆明-万象城"},
		{StoreNumber: "R581", CityStoreName: "温州-万象城"},
	},
	"zh_MO": {
		// 澳门
		{StoreNumber: "R672", CityStoreName: "澳门-银河澳门", Province: "澳門", City: "路氹", District: "銀河"},
		{StoreNumber: "R697", CityStoreName: "澳门-伦敦人购物中心", Province: "澳門", City: "路氹", District: "倫敦人"},
	},
}