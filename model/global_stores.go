package model

// 全球所有苹果零售店完整列表
var GlobalStores = map[string][]Store{
	// 中国大陆 (42家店)
	"zh_CN": ChinaStores["zh_CN"],
	
	// 中国香港 (6家店)
	"zh_HK": {
		{StoreNumber: "R428", CityStoreName: "香港-IFC Mall", Province: "香港", City: "中環", District: "金融街"},
		{StoreNumber: "R409", CityStoreName: "香港-Canton Road", Province: "香港", City: "尖沙咀", District: "廣東道"},
		{StoreNumber: "R485", CityStoreName: "香港-Festival Walk", Province: "香港", City: "九龍塘", District: "又一城"},
		{StoreNumber: "R673", CityStoreName: "香港-Causeway Bay", Province: "香港", City: "銅鑼灣", District: "希慎廣場"},
		{StoreNumber: "R611", CityStoreName: "香港-apm", Province: "香港", City: "觀塘", District: "創紀之城"},
		{StoreNumber: "R705", CityStoreName: "香港-K11 Musea", Province: "香港", City: "尖沙咀", District: "K11 MUSEA"},
	},
	
	// 中国澳门 (2家店)
	"zh_MO": ChinaStores["zh_MO"],
	
	// 中国台湾 (2家店)
	"zh_TW": {
		{StoreNumber: "R713", CityStoreName: "台北-101"},
		{StoreNumber: "R694", CityStoreName: "台北-信义A13"},
	},
	
	// 日本 (10家店)
	"ja_JP": {
		{StoreNumber: "R119", CityStoreName: "东京-银座", Province: "東京都", City: "中央区", District: "銀座"},
		{StoreNumber: "R224", CityStoreName: "东京-涩谷", Province: "東京都", City: "渋谷区", District: "渋谷"},
		{StoreNumber: "R125", CityStoreName: "东京-新宿", Province: "東京都", City: "新宿区", District: "新宿"},
		{StoreNumber: "R374", CityStoreName: "大阪-心斋桥", Province: "大阪府", City: "大阪市", District: "中央区"},
		{StoreNumber: "R504", CityStoreName: "京都", Province: "京都府", City: "京都市", District: "下京区"},
		{StoreNumber: "R580", CityStoreName: "福冈", Province: "福岡県", City: "福岡市", District: "中央区"},
		{StoreNumber: "R707", CityStoreName: "东京-丸之内", Province: "東京都", City: "千代田区", District: "丸の内"},
		{StoreNumber: "R676", CityStoreName: "川崎", Province: "神奈川県", City: "川崎市", District: "幸区"},
		{StoreNumber: "R692", CityStoreName: "名古屋-荣", Province: "愛知県", City: "名古屋市", District: "中区"},
		{StoreNumber: "R720", CityStoreName: "表参道", Province: "東京都", City: "渋谷区", District: "表参道"},
	},
	
	// 韩国 (4家店)
	"ko_KR": {
		{StoreNumber: "R503", CityStoreName: "首尔-江南"},
		{StoreNumber: "R629", CityStoreName: "首尔-汝矣岛"},
		{StoreNumber: "R692", CityStoreName: "首尔-明洞"},
		{StoreNumber: "R721", CityStoreName: "首尔-Gangnam"},
	},
	
	// 新加坡 (3家店)
	"en_SG": {
		{StoreNumber: "R360", CityStoreName: "新加坡-乌节路"},
		{StoreNumber: "R606", CityStoreName: "新加坡-滨海湾金沙"},
		{StoreNumber: "R673", CityStoreName: "新加坡-Jewel"},
	},
	
	// 泰国 (2家店)
	"th_TH": {
		{StoreNumber: "R638", CityStoreName: "曼谷-IconSiam"},
		{StoreNumber: "R650", CityStoreName: "曼谷-Central World"},
	},
	
	// 马来西亚 (1家店)
	"ms_MY": {
		{StoreNumber: "R723", CityStoreName: "吉隆坡-The Exchange TRX"},
	},
	
	// 印度 (2家店)
	"en_IN": {
		{StoreNumber: "R772", CityStoreName: "孟买-BKC"},
		{StoreNumber: "R773", CityStoreName: "新德里-Select CityWalk"},
	},
	
	// 阿联酋 (4家店)
	"ar_AE": {
		{StoreNumber: "R505", CityStoreName: "迪拜-Mall of Emirates"},
		{StoreNumber: "R579", CityStoreName: "迪拜-Dubai Mall"},
		{StoreNumber: "R761", CityStoreName: "阿布扎比-Yas Mall"},
		{StoreNumber: "R736", CityStoreName: "阿布扎比-Al Maryah Island"},
	},
	
	// 土耳其 (3家店)
	"tr_TR": {
		{StoreNumber: "R591", CityStoreName: "伊斯坦布尔-Zorlu Center"},
		{StoreNumber: "R609", CityStoreName: "伊斯坦布尔-Akasya"},
		{StoreNumber: "R706", CityStoreName: "伊斯坦布尔-Bagdat Caddesi"},
	},
	
	// 美国 (272家店 - 主要城市精选)
	"en_US": {
		// 纽约州
		{StoreNumber: "R001", CityStoreName: "纽约-第五大道"},
		{StoreNumber: "R100", CityStoreName: "纽约-Grand Central"},
		{StoreNumber: "R252", CityStoreName: "纽约-Upper West Side"},
		{StoreNumber: "R271", CityStoreName: "纽约-Upper East Side"},
		{StoreNumber: "R286", CityStoreName: "纽约-SoHo"},
		{StoreNumber: "R305", CityStoreName: "纽约-West 14th Street"},
		{StoreNumber: "R406", CityStoreName: "纽约-Brooklyn"},
		{StoreNumber: "R430", CityStoreName: "纽约-Queens Center"},
		{StoreNumber: "R522", CityStoreName: "纽约-Staten Island"},
		{StoreNumber: "R538", CityStoreName: "纽约-World Trade Center"},
		
		// 加州
		{StoreNumber: "R001", CityStoreName: "旧金山-Union Square"},
		{StoreNumber: "R013", CityStoreName: "旧金山-Chestnut Street"},
		{StoreNumber: "R016", CityStoreName: "旧金山-Stonestown"},
		{StoreNumber: "R018", CityStoreName: "帕洛阿尔托-Stanford"},
		{StoreNumber: "R019", CityStoreName: "圣何塞-Valley Fair"},
		{StoreNumber: "R047", CityStoreName: "洛杉矶-The Grove"},
		{StoreNumber: "R053", CityStoreName: "洛杉矶-Beverly Center"},
		{StoreNumber: "R067", CityStoreName: "洛杉矶-Third Street"},
		{StoreNumber: "R087", CityStoreName: "洛杉矶-Century City"},
		{StoreNumber: "R106", CityStoreName: "帕萨迪纳"},
		{StoreNumber: "R137", CityStoreName: "圣地亚哥-Fashion Valley"},
		{StoreNumber: "R172", CityStoreName: "圣地亚哥-UTC"},
		{StoreNumber: "R245", CityStoreName: "萨克拉门托-Arden Fair"},
		{StoreNumber: "R301", CityStoreName: "库比蒂诺-Apple Park Visitor Center"},
		{StoreNumber: "R379", CityStoreName: "库比蒂诺-Infinite Loop"},
		
		// 伊利诺伊州
		{StoreNumber: "R035", CityStoreName: "芝加哥-Michigan Avenue"},
		{StoreNumber: "R096", CityStoreName: "芝加哥-Lincoln Park"},
		{StoreNumber: "R094", CityStoreName: "芝加哥-Woodfield"},
		{StoreNumber: "R095", CityStoreName: "芝加哥-Old Orchard"},
		{StoreNumber: "R258", CityStoreName: "芝加哥-Oakbrook"},
		{StoreNumber: "R273", CityStoreName: "芝加哥-Deer Park"},
		
		// 德州
		{StoreNumber: "R168", CityStoreName: "休斯顿-Highland Village"},
		{StoreNumber: "R177", CityStoreName: "休斯顿-Galleria"},
		{StoreNumber: "R316", CityStoreName: "休斯顿-Memorial City"},
		{StoreNumber: "R175", CityStoreName: "奥斯汀-Domain NORTHSIDE"},
		{StoreNumber: "R338", CityStoreName: "奥斯汀-Barton Creek"},
		{StoreNumber: "R169", CityStoreName: "达拉斯-Knox Street"},
		{StoreNumber: "R217", CityStoreName: "达拉斯-Galleria Dallas"},
		{StoreNumber: "R272", CityStoreName: "达拉斯-NorthPark Center"},
		{StoreNumber: "R348", CityStoreName: "圣安东尼奥-La Cantera"},
		{StoreNumber: "R419", CityStoreName: "圣安东尼奥-The Shops at La Cantera"},
		
		// 佛罗里达州
		{StoreNumber: "R321", CityStoreName: "迈阿密-Lincoln Road"},
		{StoreNumber: "R293", CityStoreName: "迈阿密-Aventura"},
		{StoreNumber: "R391", CityStoreName: "迈阿密-Dadeland"},
		{StoreNumber: "R337", CityStoreName: "奥兰多-Millenia"},
		{StoreNumber: "R351", CityStoreName: "奥兰多-Florida Mall"},
		{StoreNumber: "R364", CityStoreName: "坦帕-International Plaza"},
		{StoreNumber: "R380", CityStoreName: "坦帕-Brandon"},
		
		// 马萨诸塞州
		{StoreNumber: "R104", CityStoreName: "波士顿-Boylston Street"},
		{StoreNumber: "R189", CityStoreName: "波士顿-CambridgeSide"},
		{StoreNumber: "R199", CityStoreName: "波士顿-Chestnut Hill"},
		{StoreNumber: "R207", CityStoreName: "波士顿-Legacy Place"},
		{StoreNumber: "R253", CityStoreName: "波士顿-Burlington"},
		
		// 华盛顿州
		{StoreNumber: "R116", CityStoreName: "西雅图-University Village"},
		{StoreNumber: "R144", CityStoreName: "西雅图-Alderwood"},
		{StoreNumber: "R202", CityStoreName: "西雅图-Bellevue Square"},
		{StoreNumber: "R254", CityStoreName: "西雅图-Southcenter"},
		{StoreNumber: "R296", CityStoreName: "西雅图-Tacoma Mall"},
		
		// 内华达州
		{StoreNumber: "R214", CityStoreName: "拉斯维加斯-Fashion Show"},
		{StoreNumber: "R434", CityStoreName: "拉斯维加斯-The Forum Shops"},
		{StoreNumber: "R462", CityStoreName: "拉斯维加斯-Town Square"},
		{StoreNumber: "R506", CityStoreName: "拉斯维加斯-Downtown Summerlin"},
		
		// 亚利桑那州
		{StoreNumber: "R026", CityStoreName: "凤凰城-Chandler Fashion Center"},
		{StoreNumber: "R031", CityStoreName: "斯科茨代尔-Fashion Square"},
		{StoreNumber: "R247", CityStoreName: "凤凰城-Arrowhead"},
		{StoreNumber: "R267", CityStoreName: "凤凰城-SanTan Village"},
		{StoreNumber: "R292", CityStoreName: "斯科茨代尔-Scottsdale Quarter"},
		
		// 更多州略...
	},
	
	// 加拿大 (29家店)
	"en_CA": {
		// 安大略省
		{StoreNumber: "R121", CityStoreName: "多伦多-Eaton Centre"},
		{StoreNumber: "R280", CityStoreName: "多伦多-Yorkdale"},
		{StoreNumber: "R336", CityStoreName: "多伦多-Fairview Mall"},
		{StoreNumber: "R447", CityStoreName: "多伦多-Square One"},
		{StoreNumber: "R483", CityStoreName: "多伦多-Sherway Gardens"},
		{StoreNumber: "R385", CityStoreName: "渥太华-Rideau"},
		{StoreNumber: "R446", CityStoreName: "渥太华-Bayshore Shopping Centre"},
		{StoreNumber: "R276", CityStoreName: "伦敦-Masonville Place"},
		{StoreNumber: "R298", CityStoreName: "滑铁卢-Conestoga"},
		{StoreNumber: "R481", CityStoreName: "伯灵顿-Mapleview Centre"},
		{StoreNumber: "R530", CityStoreName: "马卡姆-Markville Shopping Centre"},
		
		// 不列颠哥伦比亚省
		{StoreNumber: "R267", CityStoreName: "温哥华-Pacific Centre"},
		{StoreNumber: "R326", CityStoreName: "温哥华-Metrotown"},
		{StoreNumber: "R452", CityStoreName: "温哥华-Oakridge"},
		{StoreNumber: "R454", CityStoreName: "温哥华-Richmond Centre"},
		{StoreNumber: "R531", CityStoreName: "温哥华-Coquitlam Centre"},
		{StoreNumber: "R568", CityStoreName: "温哥华-Guildford Town Centre"},
		
		// 阿尔伯塔省
		{StoreNumber: "R256", CityStoreName: "卡尔加里-Chinook Centre"},
		{StoreNumber: "R461", CityStoreName: "卡尔加里-Market Mall"},
		{StoreNumber: "R521", CityStoreName: "埃德蒙顿-West Edmonton"},
		{StoreNumber: "R525", CityStoreName: "埃德蒙顿-Southgate Centre"},
		
		// 魁北克省
		{StoreNumber: "R333", CityStoreName: "蒙特利尔-Sainte-Catherine"},
		{StoreNumber: "R423", CityStoreName: "蒙特利尔-Fairview Pointe-Claire"},
		{StoreNumber: "R463", CityStoreName: "蒙特利尔-DIX30"},
		{StoreNumber: "R495", CityStoreName: "拉瓦尔-Carrefour Laval"},
		
		// 其他省份
		{StoreNumber: "R289", CityStoreName: "温尼伯-Polo Park"},
		{StoreNumber: "R384", CityStoreName: "哈利法克斯-Halifax Shopping Centre"},
		{StoreNumber: "R648", CityStoreName: "里贾纳-Southland Mall"},
		{StoreNumber: "R691", CityStoreName: "萨斯卡通-Midtown Plaza"},
	},
	
	// 英国 (39家店)
	"en_GB": {
		// 伦敦
		{StoreNumber: "R092", CityStoreName: "伦敦-Regent Street"},
		{StoreNumber: "R410", CityStoreName: "伦敦-Covent Garden"},
		{StoreNumber: "R435", CityStoreName: "伦敦-Battersea"},
		{StoreNumber: "R499", CityStoreName: "伦敦-Brompton Road"},
		{StoreNumber: "R279", CityStoreName: "伦敦-Brent Cross"},
		{StoreNumber: "R342", CityStoreName: "伦敦-Westfield London"},
		{StoreNumber: "R376", CityStoreName: "伦敦-Stratford City"},
		{StoreNumber: "R400", CityStoreName: "伦敦-White City"},
		
		// 英格兰其他城市
		{StoreNumber: "R163", CityStoreName: "曼彻斯特-Arndale"},
		{StoreNumber: "R491", CityStoreName: "曼彻斯特-Trafford Centre"},
		{StoreNumber: "R433", CityStoreName: "伯明翰-Bullring"},
		{StoreNumber: "R245", CityStoreName: "利物浦-Liverpool ONE"},
		{StoreNumber: "R319", CityStoreName: "布里斯托-Cabot Circus"},
		{StoreNumber: "R340", CityStoreName: "纽卡斯尔-Eldon Square"},
		{StoreNumber: "R354", CityStoreName: "利兹-Trinity Leeds"},
		{StoreNumber: "R386", CityStoreName: "谢菲尔德-Meadowhall"},
		{StoreNumber: "R395", CityStoreName: "南安普顿-WestQuay"},
		{StoreNumber: "R408", CityStoreName: "诺维奇-Chapelfield"},
		{StoreNumber: "R412", CityStoreName: "布莱顿-Churchill Square"},
		{StoreNumber: "R423", CityStoreName: "雷丁-The Oracle"},
		{StoreNumber: "R436", CityStoreName: "米尔顿凯恩斯-Milton Keynes"},
		{StoreNumber: "R442", CityStoreName: "埃克塞特-Princesshay"},
		{StoreNumber: "R456", CityStoreName: "莱斯特-Highcross"},
		{StoreNumber: "R475", CityStoreName: "剑桥-Grand Arcade"},
		{StoreNumber: "R479", CityStoreName: "沃特福德-intu Watford"},
		{StoreNumber: "R493", CityStoreName: "诺丁汉-Victoria Centre"},
		{StoreNumber: "R512", CityStoreName: "巴斯-SouthGate"},
		
		// 苏格兰
		{StoreNumber: "R113", CityStoreName: "格拉斯哥-Buchanan Street"},
		{StoreNumber: "R362", CityStoreName: "格拉斯哥-Braehead"},
		{StoreNumber: "R519", CityStoreName: "爱丁堡-Princes Street"},
		{StoreNumber: "R537", CityStoreName: "阿伯丁-Union Square"},
		
		// 威尔士
		{StoreNumber: "R507", CityStoreName: "卡迪夫-St David's"},
		
		// 北爱尔兰
		{StoreNumber: "R178", CityStoreName: "贝尔法斯特-Victoria Square"},
	},
	
	// 法国 (20家店)
	"fr_FR": {
		// 巴黎大区
		{StoreNumber: "R306", CityStoreName: "巴黎-Opéra"},
		{StoreNumber: "R373", CityStoreName: "巴黎-Champs-Élysées"},
		{StoreNumber: "R490", CityStoreName: "巴黎-Marché Saint-Germain"},
		{StoreNumber: "R593", CityStoreName: "巴黎-Forum des Halles"},
		{StoreNumber: "R375", CityStoreName: "巴黎-Carrousel du Louvre"},
		{StoreNumber: "R446", CityStoreName: "巴黎-Quatre Temps"},
		{StoreNumber: "R471", CityStoreName: "巴黎-Rosny 2"},
		{StoreNumber: "R483", CityStoreName: "巴黎-Parly 2"},
		{StoreNumber: "R484", CityStoreName: "巴黎-Val d'Europe"},
		{StoreNumber: "R497", CityStoreName: "巴黎-Vélizy 2"},
		
		// 其他城市
		{StoreNumber: "R392", CityStoreName: "里昂-Part-Dieu"},
		{StoreNumber: "R487", CityStoreName: "里昂-Confluence"},
		{StoreNumber: "R547", CityStoreName: "马赛-Marseille"},
		{StoreNumber: "R443", CityStoreName: "里尔-Lille"},
		{StoreNumber: "R458", CityStoreName: "蒙彼利埃-Odysseum"},
		{StoreNumber: "R476", CityStoreName: "斯特拉斯堡-Rivetoile"},
		{StoreNumber: "R527", CityStoreName: "波尔多-Sainte-Catherine"},
		{StoreNumber: "R569", CityStoreName: "南特-Atlantis"},
		{StoreNumber: "R583", CityStoreName: "尼斯-CAP 3000"},
		{StoreNumber: "R598", CityStoreName: "艾克斯普罗旺斯-Aix-en-Provence"},
	},
	
	// 德国 (15家店)
	"de_DE": {
		// 柏林
		{StoreNumber: "R382", CityStoreName: "柏林-Kurfürstendamm"},
		{StoreNumber: "R576", CityStoreName: "柏林-Mall of Berlin"},
		
		// 慕尼黑
		{StoreNumber: "R394", CityStoreName: "慕尼黑-Marienplatz"},
		{StoreNumber: "R639", CityStoreName: "慕尼黑-OEZ"},
		
		// 法兰克福
		{StoreNumber: "R453", CityStoreName: "法兰克福-Große Bockenheimer"},
		{StoreNumber: "R538", CityStoreName: "法兰克福-MTZ"},
		
		// 其他城市
		{StoreNumber: "R393", CityStoreName: "汉堡-Jungfernstieg"},
		{StoreNumber: "R407", CityStoreName: "汉堡-Alstertal"},
		{StoreNumber: "R480", CityStoreName: "杜塞尔多夫-Kö-Bogen"},
		{StoreNumber: "R449", CityStoreName: "科隆-Rhein-Center"},
		{StoreNumber: "R456", CityStoreName: "科隆-Schildergasse"},
		{StoreNumber: "R469", CityStoreName: "斯图加特-Königstraße"},
		{StoreNumber: "R581", CityStoreName: "奥格斯堡-City-Galerie"},
		{StoreNumber: "R596", CityStoreName: "汉诺威-Bahnhofstraße"},
		{StoreNumber: "R634", CityStoreName: "德累斯顿-Altmarkt-Galerie"},
	},
	
	// 意大利 (17家店)
	"it_IT": {
		// 米兰
		{StoreNumber: "R596", CityStoreName: "米兰-Piazza Liberty"},
		{StoreNumber: "R408", CityStoreName: "米兰-Carosello"},
		{StoreNumber: "R489", CityStoreName: "米兰-Il Leone"},
		{StoreNumber: "R554", CityStoreName: "米兰-Fiordaliso"},
		{StoreNumber: "R885", CityStoreName: "米兰-Via del Corso"},
		
		// 罗马
		{StoreNumber: "R486", CityStoreName: "罗马-RomaEst"},
		{StoreNumber: "R492", CityStoreName: "罗马-Porta di Roma"},
		{StoreNumber: "R684", CityStoreName: "罗马-Via del Corso"},
		
		// 其他城市
		{StoreNumber: "R482", CityStoreName: "佛罗伦萨-I Gigli"},
		{StoreNumber: "R699", CityStoreName: "佛罗伦萨-Piazza della Repubblica"},
		{StoreNumber: "R484", CityStoreName: "都灵-Grugliasco"},
		{StoreNumber: "R518", CityStoreName: "都灵-Via Roma"},
		{StoreNumber: "R488", CityStoreName: "博洛尼亚-Shopville Gran Reno"},
		{StoreNumber: "R510", CityStoreName: "那不勒斯-Campania"},
		{StoreNumber: "R550", CityStoreName: "威尼斯-Nave de Vero"},
		{StoreNumber: "R578", CityStoreName: "维罗纳-Adigeo"},
		{StoreNumber: "R649", CityStoreName: "巴里-Bari"},
	},
	
	// 西班牙 (11家店)
	"es_ES": {
		// 马德里
		{StoreNumber: "R388", CityStoreName: "马德里-Puerta del Sol"},
		{StoreNumber: "R467", CityStoreName: "马德里-Parquesur"},
		{StoreNumber: "R584", CityStoreName: "马德里-Xanadú"},
		{StoreNumber: "R603", CityStoreName: "马德里-La Vaguada"},
		
		// 巴塞罗那
		{StoreNumber: "R419", CityStoreName: "巴塞罗那-Passeig de Gràcia"},
		{StoreNumber: "R522", CityStoreName: "巴塞罗那-La Maquinista"},
		{StoreNumber: "R625", CityStoreName: "巴塞罗那-Diagonal Mar"},
		
		// 其他城市
		{StoreNumber: "R468", CityStoreName: "瓦伦西亚-Valencia"},
		{StoreNumber: "R515", CityStoreName: "毕尔巴鄂-Bilbao"},
		{StoreNumber: "R583", CityStoreName: "塞维利亚-Sevilla"},
		{StoreNumber: "R610", CityStoreName: "马拉加-Málaga"},
	},
	
	// 荷兰 (3家店)
	"nl_NL": {
		{StoreNumber: "R372", CityStoreName: "阿姆斯特丹-Amsterdam"},
		{StoreNumber: "R521", CityStoreName: "海牙-The Hague"},
		{StoreNumber: "R570", CityStoreName: "埃因霍温-Eindhoven"},
	},
	
	// 比利时 (1家店)
	"nl_BE": {
		{StoreNumber: "R503", CityStoreName: "布鲁塞尔-Brussels"},
	},
	
	// 瑞士 (4家店)
	"de_CH": {
		{StoreNumber: "R507", CityStoreName: "苏黎世-Bahnhofstrasse"},
		{StoreNumber: "R555", CityStoreName: "苏黎世-Glattzentrum"},
		{StoreNumber: "R479", CityStoreName: "日内瓦-Rue de Rive"},
		{StoreNumber: "R532", CityStoreName: "巴塞尔-Freie Strasse"},
	},
	
	// 瑞典 (3家店)
	"sv_SE": {
		{StoreNumber: "R467", CityStoreName: "斯德哥尔摩-Täby Centrum"},
		{StoreNumber: "R508", CityStoreName: "斯德哥尔摩-Mall of Scandinavia"},
		{StoreNumber: "R604", CityStoreName: "斯德哥尔摩-Väla Centrum"},
	},
	
	// 奥地利 (2家店)
	"de_AT": {
		{StoreNumber: "R524", CityStoreName: "维也纳-Kärntner Straße"},
		{StoreNumber: "R770", CityStoreName: "维也纳-Vienna Airport"},
	},
	
	// 澳大利亚 (22家店)
	"en_AU": {
		// 新南威尔士州
		{StoreNumber: "R252", CityStoreName: "悉尼-George Street"},
		{StoreNumber: "R426", CityStoreName: "悉尼-Bondi"},
		{StoreNumber: "R435", CityStoreName: "悉尼-Chatswood Chase"},
		{StoreNumber: "R462", CityStoreName: "悉尼-Miranda"},
		{StoreNumber: "R497", CityStoreName: "悉尼-Macquarie Centre"},
		{StoreNumber: "R505", CityStoreName: "悉尼-Castle Towers"},
		{StoreNumber: "R545", CityStoreName: "悉尼-Penrith"},
		{StoreNumber: "R756", CityStoreName: "悉尼-Broadway"},
		
		// 维多利亚州
		{StoreNumber: "R255", CityStoreName: "墨尔本-Doncaster"},
		{StoreNumber: "R308", CityStoreName: "墨尔本-Southland"},
		{StoreNumber: "R389", CityStoreName: "墨尔本-Chadstone"},
		{StoreNumber: "R396", CityStoreName: "墨尔本-Highpoint"},
		{StoreNumber: "R491", CityStoreName: "墨尔本-Federation Square"},
		{StoreNumber: "R550", CityStoreName: "墨尔本-Fountain Gate"},
		
		// 昆士兰州
		{StoreNumber: "R296", CityStoreName: "布里斯班-Brisbane"},
		{StoreNumber: "R431", CityStoreName: "布里斯班-Chermside"},
		{StoreNumber: "R456", CityStoreName: "布里斯班-Carindale"},
		{StoreNumber: "R537", CityStoreName: "黄金海岸-Pacific Fair"},
		{StoreNumber: "R536", CityStoreName: "黄金海岸-Robina"},
		
		// 西澳大利亚州
		{StoreNumber: "R449", CityStoreName: "珀斯-Perth City"},
		{StoreNumber: "R516", CityStoreName: "珀斯-Garden City"},
		
		// 南澳大利亚州
		{StoreNumber: "R570", CityStoreName: "阿德莱德-Rundle Place"},
	},
	
	// 巴西 (2家店)
	"pt_BR": {
		{StoreNumber: "R508", CityStoreName: "里约热内卢-VillageMall"},
		{StoreNumber: "R584", CityStoreName: "圣保罗-Morumbi"},
	},
	
	// 墨西哥 (2家店)
	"es_MX": {
		{StoreNumber: "R482", CityStoreName: "墨西哥城-Antara"},
		{StoreNumber: "R554", CityStoreName: "墨西哥城-Via Santa Fe"},
	},
}