package main

import (
	"errors"
	"fmt"
	"fyne.io/fyne"
	"fyne.io/fyne/app"
	"fyne.io/fyne/container"
	"fyne.io/fyne/dialog"
	"fyne.io/fyne/layout"
	"fyne.io/fyne/widget"
	"github.com/parnurzeal/gorequest"
	"github.com/tidwall/gjson"
	"log"
	"net/url"
	"os/exec"
	"runtime"
	"strings"
	"time"
)

const VERSION = "1.0.5"

var isListen = false
var body *widget.Label
var tip *widget.Label
var status *widget.Label
var versionWgt *widget.Hyperlink
var modelCode = map[string]string{
	"iphone12mini": "H",
	"iphone12": "F",
	"iphone12pro": "A",
	"iphone12promax": "G",
}

var models = map[string][]string{
	"CN/zh_CN": {
		"iphone12mini 64gb 黑色-MG7Y3CH/A",
		"iphone12mini 64gb 白色-MG803CH/A",
		"iphone12mini 64gb 蓝色-MG823CH/A",
		"iphone12mini 64gb 绿色-MG833CH/A",
		"iphone12mini 64gb 红色-MG813CH/A",
		"iphone12mini 128gb 黑色-MG843CH/A",
		"iphone12mini 128gb 白色-MG853CH/A",
		"iphone12mini 128gb 蓝色-MG873CH/A",
		"iphone12mini 128gb 绿色-MG883CH/A",
		"iphone12mini 128gb 红色-MG863CH/A",
		"iphone12mini 256gb 黑色-MG893CH/A",
		"iphone12mini 256gb 白色-MG8A3CH/A",
		"iphone12mini 256gb 蓝色-MG8D3CH/A",
		"iphone12mini 256gb 绿色-MG8E3CH/A",
		"iphone12mini 256gb 红色-MG8C3CH/A",
		"iphone12 64gb 黑色-MGGM3CH/A",
		"iphone12 64gb 白色-MGGN3CH/A",
		"iphone12 64gb 蓝色-MGGQ3CH/A",
		"iphone12 64gb 绿色-MGGT3CH/A",
		"iphone12 64gb 红色-MGGP3CH/A",
		"iphone12 128gb 黑色-MGGU3CH/A",
		"iphone12 128gb 白色-MGGV3CH/A",
		"iphone12 128gb 蓝色-MGGX3CH/A",
		"iphone12 128gb 绿色-MGGY3CH/A",
		"iphone12 128gb 红色-MGGW3CH/A",
		"iphone12 256gb 黑色-MGH13CH/A",
		"iphone12 256gb 白色-MGH23CH/A",
		"iphone12 256gb 蓝色-MGH43CH/A",
		"iphone12 256gb 绿色-MGH53CH/A",
		"iphone12 256gb 红色-MGH33CH/A",
		"iphone12pro 128gb 石墨色-MGL93CH/A",
		"iphone12pro 128gb 银色-MGLA3CH/A",
		"iphone12pro 128gb 金色-MGLC3CH/A",
		"iphone12pro 128gb 海蓝色-MGLD3CH/A",
		"iphone12pro 256gb 石墨色-MGLE3CH/A",
		"iphone12pro 256gb 银色-MGLF3CH/A",
		"iphone12pro 256gb 金色-MGLG3CH/A",
		"iphone12pro 256gb 海蓝色-MGLH3CH/A",
		"iphone12pro 512gb 石墨色-MGLJ3CH/A",
		"iphone12pro 512gb 银色-MGLK3CH/A",
		"iphone12pro 512gb 金色-MGLL3CH/A",
		"iphone12pro 512gb 海蓝色-MGLM3CH/A",
		"iphone12promax 128gb 石墨色-MGC03CH/A",
		"iphone12promax 128gb 银色-MGC13CH/A",
		"iphone12promax 128gb 金色-MGC23CH/A",
		"iphone12promax 128gb 海蓝色-MGC33CH/A",
		"iphone12promax 256gb 石墨色-MGC43CH/A",
		"iphone12promax 256gb 银色-MGC53CH/A",
		"iphone12promax 256gb 金色-MGC63CH/A",
		"iphone12promax 256gb 海蓝色-MGC73CH/A",
		"iphone12promax 512gb 石墨色-MGC93CH/A",
		"iphone12promax 512gb 银色-MGCA3CH/A",
		"iphone12promax 512gb 金色-MGCC3CH/A",
		"iphone12promax 512gb 海蓝色-MGCE3CH/A",
	},
	"MO/zh_MO": {
		"iphone12mini 64gb 黑色-MG7Y3ZA/A",
		"iphone12mini 64gb 白色-MG803ZA/A",
		"iphone12mini 64gb 藍色-MG823ZA/A",
		"iphone12mini 64gb 綠色-MG833ZA/A",
		"iphone12mini 64gb (PRODUCT)RED-MG813ZA/A",
		"iphone12mini 128gb 黑色-MG843ZA/A",
		"iphone12mini 128gb 白色-MG853ZA/A",
		"iphone12mini 128gb 藍色-MG873ZA/A",
		"iphone12mini 128gb 綠色-MG883ZA/A",
		"iphone12mini 128gb (PRODUCT)RED-MG863ZA/A",
		"iphone12mini 256gb 黑色-MG893ZA/A",
		"iphone12mini 256gb 白色-MG8A3ZA/A",
		"iphone12mini 256gb 藍色-MG8D3ZA/A",
		"iphone12mini 256gb 綠色-MG8E3ZA/A",
		"iphone12mini 256gb (PRODUCT)RED-MG8C3ZA/A",
		"iphone12 64gb 黑色-MGGM3ZA/A",
		"iphone12 64gb 白色-MGGN3ZA/A",
		"iphone12 64gb 藍色-MGGQ3ZA/A",
		"iphone12 64gb 綠色-MGGT3ZA/A",
		"iphone12 64gb (PRODUCT)RED-MGGP3ZA/A",
		"iphone12 128gb 黑色-MGGU3ZA/A",
		"iphone12 128gb 白色-MGGV3ZA/A",
		"iphone12 128gb 藍色-MGGX3ZA/A",
		"iphone12 128gb 綠色-MGGY3ZA/A",
		"iphone12 128gb (PRODUCT)RED-MGGW3ZA/A",
		"iphone12 256gb 黑色-MGH13ZA/A",
		"iphone12 256gb 白色-MGH23ZA/A",
		"iphone12 256gb 藍色-MGH43ZA/A",
		"iphone12 256gb 綠色-MGH53ZA/A",
		"iphone12 256gb (PRODUCT)RED-MGH33ZA/A",
		"iphone12pro 128gb 石墨色-MGL93ZA/A",
		"iphone12pro 128gb 銀色-MGLA3ZA/A",
		"iphone12pro 128gb 金色-MGLC3ZA/A",
		"iphone12pro 128gb 太平洋藍色-MGLD3ZA/A",
		"iphone12pro 256gb 石墨色-MGLE3ZA/A",
		"iphone12pro 256gb 銀色-MGLF3ZA/A",
		"iphone12pro 256gb 金色-MGLG3ZA/A",
		"iphone12pro 256gb 太平洋藍色-MGLH3ZA/A",
		"iphone12pro 512gb 石墨色-MGLJ3ZA/A",
		"iphone12pro 512gb 銀色-MGLK3ZA/A",
		"iphone12pro 512gb 金色-MGLL3ZA/A",
		"iphone12pro 512gb 太平洋藍色-MGLM3ZA/A",
		"iphone12promax 128gb 石墨色-MGC03ZA/A",
		"iphone12promax 128gb 銀色-MGC13ZA/A",
		"iphone12promax 128gb 金色-MGC23ZA/A",
		"iphone12promax 128gb 太平洋藍色-MGC33ZA/A",
		"iphone12promax 256gb 石墨色-MGC43ZA/A",
		"iphone12promax 256gb 銀色-MGC53ZA/A",
		"iphone12promax 256gb 金色-MGC63ZA/A",
		"iphone12promax 256gb 太平洋藍色-MGC73ZA/A",
		"iphone12promax 512gb 石墨色-MGC93ZA/A",
		"iphone12promax 512gb 銀色-MGCA3ZA/A",
		"iphone12promax 512gb 金色-MGCC3ZA/A",
		"iphone12promax 512gb 太平洋藍色-MGCE3ZA/A",
	},
}
var selectQuantity = "1" // 默认一台
var selectStore string
var selectModel string
var listenStores map[string]string
var w fyne.Window

// 地区，中国大陆: CN/zh_CN, 中国澳门: MO/zh_MO
var area = "CN/zh_CN"

func main() {
	// 调试模式
	//os.Setenv("FYNE_FONT", "./fzhtk.ttf")
	a := app.NewWithID("ip12")
	// 打包时自动加载字体
	a.Settings().SetTheme(&myTheme{})
	w = a.NewWindow("iPhone12|Mini|Pro|ProMax")
	w.Resize(fyne.NewSize(800, 600))

	body = widget.NewLabel("")
	tip = widget.NewLabel("请选择门店和型号")
	status = widget.NewLabel("暂停")
	listenStores = make(map[string]string)
	stores := getStores()
	// 单次抢购数量，最多2
	quantityWgt := widget.NewSelect([]string{"1", "2"}, func(b string) {
		selectQuantity = b
	})
	quantityWgt.PlaceHolder ="预约台数"
	// 门店选择组件
	storesWgt := widget.NewSelect(stores, func(b string) {
		selectStore = b
	})
	//
	modelsWgt := widget.NewSelect(models[area], func(b string) {
		selectModel = b
	})
	// 地区选择组件
	areaWgt := widget.NewSelect([]string{"中国大陆", "中国澳门"}, func(b string) {
		area = "CN/zh_CN"

		if b == "中国澳门" {
			area = "MO/zh_MO"
		}
		// 重置门店
		stores = getStores()
		storesWgt.Options = stores
		storesWgt.ClearSelected()
		// 重置型号
		modelsWgt.Options = models[area]
		modelsWgt.ClearSelected()
		quantityWgt.ClearSelected()

		// 重置已有变量
		listenStores = map[string]string{}
		selectStore = ""
		selectModel = ""
		selectQuantity = "1"
		isListen = false
		body.SetText("")
		status.SetText("暂停")
	})
	areaWgt.PlaceHolder ="地区选择，默认中国大陆"

	releaseUrl, _ := url.Parse("https://github.com/hteen/apple-store-helper/releases")
	versionWgt = widget.NewHyperlink("", releaseUrl)
	go getLatestVersion()

	w.SetContent(container.NewVBox(
		container.NewHBox(
			widget.NewLabel("1.首次运行请先获取Apple注册码，确保能正确打开网页\n" +
				"2.选择门店和型号，点击添加按钮\n" +
				"3.点击开始\n" +
				"4.匹配到之后会直接进入门店预购页面，输入注册码选择预约时间即可",
			),
			layout.NewSpacer(),
			widget.NewLabel("当前版本:"+VERSION),
			versionWgt,
		),
		areaWgt,
		container.NewHBox(
			storesWgt,
			modelsWgt,
			widget.NewButton("添加", func() {
				if selectModel != "" && selectStore != "" {
					md := strings.Split(selectStore, "-")[0]+"."+strings.Split(selectModel, "-")[1]
					mdText := strings.Split(selectStore, "-")[1]+" "+strings.Split(selectModel, "-")[0]

					if !inArray(listenStores, md) {
						listenStores[md] = mdText
					}

					body.SetText(strings.Join(getValues(listenStores), "\n"))
				}
			}),
			widget.NewButton("清空", func() {
				listenStores = map[string]string{}
				body.SetText("")
				isListen = false
				status.SetText("暂停")
			}),
		),
		tip,
		body,
		layout.NewSpacer(),
		container.NewHBox(
			quantityWgt,
			widget.NewButton("开始", func() {
				if len(listenStores) < 1 {
					dialog.NewError(errors.New("请添加要监听的门店和型号"), w)
					return
				}

				isListen = true
				status.SetText("监听中")
			}),
			widget.NewButton("暂停", func() {
				isListen = false
				status.SetText("暂停")
			}),
			widget.NewButton("12mini注册码", func() {
				go registerCode("iphone12mini")
			}),
			widget.NewButton("12注册码", func() {
				go registerCode("iphone12")
			}),
			widget.NewButton("12Pro注册码", func() {
				go registerCode("iphone12pro")
			}),
			widget.NewButton("ProMax注册码", func() {
				go registerCode("iphone12promax")
			}),
			widget.NewButton("退出", func() {
				a.Quit()
			}),
			layout.NewSpacer(),
			widget.NewLabel("状态: "),
			status,
		),
	))
	go listen()
	w.ShowAndRun()
}

func listen() {
	for  {
		time.Sleep(time.Second*1)

		if !isListen {
			continue
		}

		sku := map[string]string{}
		str := ""
		t := time.Now().Format("2006-01-02 15:04:05")
		for model, title := range listenStores {
			md := title2model(title)
			if sku[md] == "" {
				skuUrl := "https://reserve-prime.apple.com/"+area+"/reserve/"+modelCode[md]+"/availability.json"
				_, bd, _ := gorequest.New().Get(skuUrl).End()
				sku[md] = bd
			}

			value := gjson.Get(sku[md], "stores."+model+".availability")
			if value.Map()["contract"].Bool() && value.Map()["unlocked"].Bool() {
				openBrowser(caleURL(model, title))

				status.SetText("暂停")
				isListen = false
				dialog.NewInformation("匹配成功", "已匹配到: " + title+ ", 暂停监听", w)
			} else {
				str += t+" "+title+"无货\n"
			}
		}

		body.SetText(str)
	}
}

// 帮助提前获取注册码
func registerCode(model string){
	_, bd, errs := gorequest.New().Get("https://reserve-prime.apple.com/"+area+"/reserve/"+modelCode[model]+"/availability.json").End()
	if len(errs) != 0 {
		dialog.NewError(errs[0], w)
		return
	}

	// 寻找任意一个有货门店
	for store, items := range gjson.Get(bd, "stores").Map() {
		for k,v := range items.Map(){
			if v.Get("availability.contract").Bool() && v.Get("availability.unlocked").Bool() {
				openBrowser(model2Url(model, store, k))
				return
			}
		}
	}

	dialog.NewError(errors.New("所有门店无货，无法前往注册码页面"), w)
}

// 型号对应预约地址
func model2Url(model string, store string, partNumber string) string {
	return "https://reserve-prime.apple.com/"+area+"/reserve/"+modelCode[model]+"?quantity="+selectQuantity+"&anchor-store="+store+
		"&store="+store+"&partNumber="+partNumber+"&plan=unlocked"
}

func caleURL(model string, title string)  string {
	// e.g: [R389.MGL93CH/A] -> [R389 MGL93CH/A]
	m := strings.Split(model, ".")
	return model2Url(title2model(title), m[0], m[1])
}

func title2model(title string) string {
	t := strings.Split(title, " ")
	t = t[len(t) - 3:]
	return t[0]
}

func openBrowser(url string) {
	var err error
	switch runtime.GOOS {
	case "linux":
		err = exec.Command("xdg-open", url).Start()
	case "windows":
		err = exec.Command("rundll32", "url.dll,FileProtocolHandler", url).Start()
	case "darwin":
		err = exec.Command("open", url).Start()
	default:
		err = fmt.Errorf("unsupported platform")
	}
	if err != nil {
		dialog.NewError(errors.New("打开网页失败，请自行手动操作\n"+url), w)
	}
}

func getStores() []string {
	// 门店列表
	var stores []string

	availability := "https://reserve-prime.apple.com/"+area+"/reserve/A/stores.json"
	_, bd, errs := gorequest.New().Get(availability).End()

	if len(errs) != 0 {
		log.Fatalln(errs)
	}

	for _, store := range gjson.Get(bd, "stores").Array() {
		str := store.Get("storeNumber").String()+
			"-"+store.Get("city").String()+
			" "+store.Get("storeName").String()

		stores = append(stores, str)
	}

	return stores
}

func inArray(slice map[string]string, s string) bool {
	for key := range slice {
		if key == s {
			return true
		}
	}
	return false
}

func getValues(slice map[string]string) []string {
	var values []string
	for _, value := range slice {
		values = append(values, value)
	}

	return values
}

// 版本查询
func getLatestVersion() {
	_, bd, _ := gorequest.New().Get("https://api.github.com/repos/hteen/apple-store-helper/releases").End()
	latest := gjson.Get(bd, "0.tag_name").String()
	if latest != "" {
		versionWgt.SetText("最新:"+latest)
	} else {
		versionWgt.SetText("最新版查询失败")
	}
}