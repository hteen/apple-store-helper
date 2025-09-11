package model

// 中国省市区数据结构
type Location struct {
	Province string
	Cities   map[string][]string // city -> districts
}

// 中国省市区数据（包含Apple Store的地区）
var ChinaLocations = []Location{
	{
		Province: "北京",
		Cities: map[string][]string{
			"北京": {"朝阳区", "东城区", "西城区", "海淀区", "丰台区", "石景山区", "通州区", "昌平区", "大兴区", "顺义区"},
		},
	},
	{
		Province: "上海",
		Cities: map[string][]string{
			"上海": {"黄浦区", "浦东新区", "静安区", "徐汇区", "杨浦区", "普陀区", "闵行区", "虹口区", "长宁区", "宝山区"},
		},
	},
	{
		Province: "广东",
		Cities: map[string][]string{
			"深圳": {"南山区", "罗湖区", "福田区", "宝安区", "龙岗区", "龙华区", "盐田区", "光明区"},
			"广州": {"天河区", "越秀区", "海珠区", "荔湾区", "白云区", "黄埔区", "番禺区", "花都区"},
		},
	},
	{
		Province: "江苏",
		Cities: map[string][]string{
			"南京": {"玄武区", "秦淮区", "建邺区", "鼓楼区", "栖霞区", "雨花台区", "江宁区", "浦口区"},
			"苏州": {"姑苏区", "工业园区", "吴中区", "相城区", "吴江区", "虎丘区"},
			"无锡": {"梁溪区", "锡山区", "惠山区", "滨湖区", "新吴区"},
		},
	},
	{
		Province: "浙江",
		Cities: map[string][]string{
			"杭州": {"上城区", "下城区", "江干区", "拱墅区", "西湖区", "滨江区", "萧山区", "余杭区"},
			"宁波": {"海曙区", "江北区", "北仑区", "镇海区", "鄞州区", "奉化区"},
		},
	},
	{
		Province: "天津",
		Cities: map[string][]string{
			"天津": {"和平区", "河东区", "河西区", "南开区", "河北区", "红桥区", "滨海新区", "东丽区", "西青区", "津南区"},
		},
	},
	{
		Province: "重庆",
		Cities: map[string][]string{
			"重庆": {"渝中区", "江北区", "南岸区", "沙坪坝区", "九龙坡区", "大渡口区", "渝北区", "巴南区", "北碚区"},
		},
	},
	{
		Province: "四川",
		Cities: map[string][]string{
			"成都": {"锦江区", "青羊区", "金牛区", "武侯区", "成华区", "龙泉驿区", "青白江区", "新都区", "温江区"},
		},
	},
	{
		Province: "湖南",
		Cities: map[string][]string{
			"长沙": {"芙蓉区", "天心区", "岳麓区", "开福区", "雨花区", "望城区", "长沙县"},
		},
	},
	{
		Province: "福建",
		Cities: map[string][]string{
			"厦门": {"思明区", "海沧区", "湖里区", "集美区", "同安区", "翔安区"},
			"福州": {"鼓楼区", "台江区", "仓山区", "马尾区", "晋安区", "长乐区"},
		},
	},
	{
		Province: "山东",
		Cities: map[string][]string{
			"青岛": {"市南区", "市北区", "黄岛区", "崂山区", "李沧区", "城阳区", "即墨区"},
			"济南": {"历下区", "市中区", "槐荫区", "天桥区", "历城区", "长清区"},
		},
	},
	{
		Province: "辽宁",
		Cities: map[string][]string{
			"沈阳": {"和平区", "沈河区", "大东区", "皇姑区", "铁西区", "苏家屯区", "浑南区", "沈北新区"},
			"大连": {"中山区", "西岗区", "沙河口区", "甘井子区", "旅顺口区", "金州区"},
		},
	},
	{
		Province: "安徽",
		Cities: map[string][]string{
			"合肥": {"瑶海区", "庐阳区", "蜀山区", "包河区", "长丰县", "肥东县", "肥西县", "庐江县", "巢湖市"},
		},
	},
	{
		Province: "河南",
		Cities: map[string][]string{
			"郑州": {"中原区", "二七区", "管城回族区", "金水区", "上街区", "惠济区", "中牟县", "巩义市", "荥阳市", "新密市", "新郑市", "登封市"},
		},
	},
	{
		Province: "湖北",
		Cities: map[string][]string{
			"武汉": {"江岸区", "江汉区", "硚口区", "汉阳区", "武昌区", "青山区", "洪山区", "东西湖区", "汉南区", "蔡甸区", "江夏区", "黄陂区", "新洲区"},
		},
	},
	{
		Province: "云南",
		Cities: map[string][]string{
			"昆明": {"五华区", "盘龙区", "官渡区", "西山区", "东川区", "呈贡区", "晋宁区", "富民县", "宜良县", "石林彝族自治县", "嵩明县", "禄劝彝族苗族自治县", "寻甸回族彝族自治县", "安宁市"},
		},
	},
	{
		Province: "广西壮族自治区",
		Cities: map[string][]string{
			"南宁": {"兴宁区", "青秀区", "江南区", "西乡塘区", "良庆区", "邕宁区", "武鸣区", "隆安县", "马山县", "上林县", "宾阳县", "横县"},
		},
	},
}

// 获取所有省份
func GetProvinces() []string {
	provinces := make([]string, len(ChinaLocations))
	for i, loc := range ChinaLocations {
		provinces[i] = loc.Province
	}
	return provinces
}

// 根据省份获取城市列表
func GetCitiesByProvince(province string) []string {
	for _, loc := range ChinaLocations {
		if loc.Province == province {
			cities := make([]string, 0, len(loc.Cities))
			for city := range loc.Cities {
				cities = append(cities, city)
			}
			return cities
		}
	}
	return []string{}
}

// 根据省份和城市获取区域列表
func GetDistrictsByProvinceAndCity(province, city string) []string {
	for _, loc := range ChinaLocations {
		if loc.Province == province {
			if districts, exists := loc.Cities[city]; exists {
				return districts
			}
		}
	}
	return []string{}
}
