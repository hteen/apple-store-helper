package main

import (
	"apple-store-helper/common"
	"apple-store-helper/model"
	"apple-store-helper/services"
	appTheme "apple-store-helper/theme"
	"apple-store-helper/view"
	"errors"
	"fmt"
	"log"
	"math/rand"
	"net/url"
	"sort"
	"strings"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
	"github.com/faiface/beep"
	"github.com/faiface/beep/speaker"
)

func main() {
	// 初始 mp3 播放器
	SampleRate := beep.SampleRate(44100)
	speaker.Init(SampleRate, SampleRate.N(time.Second/10))

	view.App = app.NewWithID("apple-store-helper")
	view.App.Settings().SetTheme(&appTheme.MyTheme{})
	view.Window = view.App.NewWindow("抢你妹 - Apple 产品库存监控工具")

	// 定义全局变量（后面会初始化）
	var modelWidget *widget.Select
	var areaWidget *widget.Select

	// 加载动态产品数据（必需）
	go func() {
		hasData := false
		// 默认加载中国大陆数据
		if productData, err := services.LoadProductData("cn"); err == nil {
			services.Product.UpdateFromDynamicData(productData)
			hasData = true
			log.Println("Loaded cached product data for cn")
		}

		// 加载中国大陆门店数据
		if err := services.Store.LoadForArea("cn"); err != nil {
			log.Printf("Failed to load store data for cn: %v", err)
		}

		// 如果没有缓存数据，必须从网络获取
		if !hasData {
			log.Println("No cached product data, fetching from Apple...")
			time.Sleep(time.Second) // 等待窗口初始化

			if err := services.UpdateProductDatabase("cn"); err != nil {
				dialog.ShowError(errors.New("无法获取产品数据，请检查网络连接后点击「更新数据」按钮重试"), view.Window)
			}
		}

		// 数据加载后更新型号选择器
		time.Sleep(time.Millisecond * 500)
		if services.Product.GetDynamicProducts() != nil && modelWidget != nil {
			modelSet := make(map[string]bool)
			for _, products := range services.Product.GetDynamicProducts() {
				for _, p := range products {
					// 从Title解析型号
					parts := strings.Split(p.Title, " ")
					if len(parts) >= 2 {
						for i, part := range parts {
							if strings.HasSuffix(part, "GB") || strings.HasSuffix(part, "TB") || strings.HasSuffix(part, "mm") {
								model := strings.Join(parts[:i], " ")
								if model != "" {
									modelSet[model] = true
								}
								break
							}
						}
					}
				}
			}

			var models []string
			for model := range modelSet {
				models = append(models, model)
			}

			modelWidget.Options = models
			modelWidget.Refresh()
		}
	}()

	defaultArea := services.Listen.Area.Title

	// 省市区选择器（中国大陆专用）
	var provinceWidget, cityWidget, districtWidget *widget.Select
	var storeWidget *widget.Select
	var currentProvince, currentCity, currentDistrict string

	// 门店 selector
	storeWidget = widget.NewSelect([]string{}, nil)
	storeWidget.PlaceHolder = "请先选择地区"
	storeWidget.Disable()

	// 更新门店列表的函数
	updateStoreList := func() {
		if currentProvince != "" && currentCity != "" && currentDistrict != "" {
			// 确保门店数据已加载
			selectedArea := services.Area.GetArea("中国大陆")
			if err := services.Store.LoadForArea(selectedArea.ShortCode); err != nil {
				log.Printf("Failed to load store data: %v", err)
			}

			// 根据省市区筛选门店
			var stores []string
			allStores := services.Store.ByArea(selectedArea)
			for _, store := range allStores {
				if store.Province == currentProvince &&
					store.City == currentCity &&
					store.District == currentDistrict {
					stores = append(stores, store.CityStoreName)
				}
			}

			// 如果该区域没有门店，显示附近所有门店
			if len(stores) == 0 {
				for _, store := range allStores {
					if store.Province == currentProvince && store.City == currentCity {
						stores = append(stores, store.CityStoreName)
					}
				}
			}

			storeWidget.Options = stores
			storeWidget.Enable()
			if len(stores) > 0 {
				storeWidget.PlaceHolder = "选择门店"
			} else {
				storeWidget.PlaceHolder = "该地区暂无门店"
			}
		}
	}

	// 区域选择器
	districtWidget = widget.NewSelect([]string{}, func(district string) {
		currentDistrict = district
		// 区域选择后，更新门店列表
		updateStoreList()
	})
	districtWidget.PlaceHolder = "请先选择城市"
	districtWidget.Disable()

	// 城市选择器
	cityWidget = widget.NewSelect([]string{}, func(city string) {
		currentCity = city
		// 更新区域列表
		districts := model.GetDistrictsByProvinceAndCity(currentProvince, city)
		districtWidget.Options = districts
		districtWidget.ClearSelected()
		districtWidget.Enable()
		districtWidget.PlaceHolder = "选择区域"
	})
	cityWidget.PlaceHolder = "请先选择省份"
	cityWidget.Disable()

	// 省份选择器 - 先创建，回调在地区选择器定义后设置
	provinceWidget = widget.NewSelect([]string{}, func(province string) {
		currentProvince = province
		// 更新门店列表 - 根据省份筛选
		if areaWidget != nil {
			// 确保门店数据已加载
			selectedArea := services.Area.GetArea(areaWidget.Selected)
			if err := services.Store.LoadForArea(selectedArea.ShortCode); err != nil {
				log.Printf("Failed to load store data: %v", err)
			}

			stores := services.Store.ByAreaAndProvinceForOptions(areaWidget.Selected, province)
			storeWidget.Options = stores
			storeWidget.ClearSelected()
			storeWidget.Enable()
			storeWidget.PlaceHolder = "选择门店"
		}
	})
	provinceWidget.PlaceHolder = "选择省份"

	// 型号 selector - iPhone型号选择
	var capacityWidget, colorWidget *widget.Select
	var capacityLabel, colorLabel *widget.Label

	// 创建动态标签
	capacityLabel = widget.NewLabelWithStyle("容量", fyne.TextAlignLeading, fyne.TextStyle{Bold: true})
	colorLabel = widget.NewLabelWithStyle("颜色", fyne.TextAlignLeading, fyne.TextStyle{Bold: true})

	modelWidget = widget.NewSelect([]string{}, func(selectedModel string) {
		// 根据型号类型更新标签
		if strings.Contains(selectedModel, "Watch") {
			capacityLabel.Text = "尺寸"
		} else {
			capacityLabel.Text = "容量"
		}
		capacityLabel.Refresh()

		// 从本地数据获取该型号的所有容量选项
		capacitySet := make(map[string]bool)

		if services.Product.GetDynamicProducts() != nil {
			for _, products := range services.Product.GetDynamicProducts() {
				for _, p := range products {
					// 精确匹配型号（Title应该是 "型号 容量 颜色" 格式）
					parts := strings.Split(p.Title, " ")
					if len(parts) >= 2 {
						// 获取型号部分（容量之前的所有部分）
						modelParts := []string{}
						capacityIdx := -1
						for i, part := range parts {
							if strings.HasSuffix(part, "GB") || strings.HasSuffix(part, "TB") || strings.HasSuffix(part, "mm") {
								capacityIdx = i
								break
							}
							modelParts = append(modelParts, part)
						}

						currentModel := strings.Join(modelParts, " ")
						if currentModel == selectedModel && capacityIdx >= 0 && capacityIdx < len(parts) {
							capacitySet[parts[capacityIdx]] = true
						}
					}
				}
			}
		}

		// 转换为排序的列表
		var capacities []string
		for capacity := range capacitySet {
			capacities = append(capacities, capacity)
		}
		sort.Strings(capacities)

		capacityWidget.Options = capacities
		capacityWidget.ClearSelected()
		capacityWidget.Enable()
		if strings.Contains(selectedModel, "Watch") {
			capacityWidget.PlaceHolder = "选择尺寸"
		} else {
			capacityWidget.PlaceHolder = "选择容量"
		}

		// 清空颜色选择
		colorWidget.Options = []string{}
		colorWidget.ClearSelected()
		colorWidget.Disable()
		if strings.Contains(selectedModel, "Watch") {
			colorWidget.PlaceHolder = "请先选择尺寸"
		} else {
			colorWidget.PlaceHolder = "请先选择容量"
		}
	})
	modelWidget.PlaceHolder = "选择型号"

	// 容量/尺寸 selector
	capacityWidget = widget.NewSelect([]string{}, func(selectedCapacity string) {
		if modelWidget.Selected != "" {
			// 从本地数据获取该型号+容量的所有颜色选项
			colorSet := make(map[string]bool)

			if services.Product.GetDynamicProducts() != nil {
				for _, products := range services.Product.GetDynamicProducts() {
					for _, p := range products {
						// 精确匹配型号和容量
						parts := strings.Split(p.Title, " ")
						if len(parts) >= 3 {
							// 获取型号、容量和颜色
							modelParts := []string{}
							capacityIdx := -1
							for i, part := range parts {
								if strings.HasSuffix(part, "GB") || strings.HasSuffix(part, "TB") || strings.HasSuffix(part, "mm") {
									capacityIdx = i
									break
								}
								modelParts = append(modelParts, part)
							}

							if capacityIdx >= 0 && capacityIdx < len(parts) {
								currentModel := strings.Join(modelParts, " ")
								currentCapacity := parts[capacityIdx]

								// 匹配型号和容量
								if currentModel == modelWidget.Selected && currentCapacity == selectedCapacity {
									// 获取颜色（容量后面的所有部分）
									if capacityIdx+1 < len(parts) {
										color := strings.Join(parts[capacityIdx+1:], " ")
										if color != "" {
											colorSet[color] = true
										}
									}
								}
							}
						}
					}
				}
			}

			// 转换为排序的列表
			var colors []string
			for color := range colorSet {
				colors = append(colors, color)
			}
			sort.Strings(colors)

			colorWidget.Options = colors
			colorWidget.ClearSelected()
			colorWidget.Enable()
			colorWidget.PlaceHolder = "选择颜色"
		}
	})
	capacityWidget.PlaceHolder = "请先选择型号"
	capacityWidget.Disable()

	// 颜色 selector
	colorWidget = widget.NewSelect([]string{}, nil)
	colorWidget.PlaceHolder = "请先选择容量/尺寸"
	colorWidget.Disable()

	// 创建位置选择容器（根据地区动态切换）
	// 日本邮编输入框
	zipCodeEntry := widget.NewEntry()
	zipCodeEntry.PlaceHolder = "输入邮编（如：110-0006）"
	zipCodeEntry.Disable()

	locationContainer := container.NewVBox()

	// 地区 selector - 使用下拉框选择
	areaWidget = widget.NewSelect(services.Area.ForOptions(), func(value string) {
		services.Listen.Area = services.Area.GetArea(value)

		// 尝试加载该地区的产品数据
		selectedArea := services.Area.GetArea(value)
		areaCode := selectedArea.ShortCode
		if err := services.Product.LoadForArea(areaCode); err != nil {
			// 如果没有本地数据，提示用户更新
			log.Printf("No local data for area %s, need to update: %v", areaCode, err)
			// 清空产品选择器
			modelWidget.Options = []string{}
			modelWidget.ClearSelected()
			modelWidget.PlaceHolder = "请先更新数据"
			capacityWidget.Options = []string{}
			capacityWidget.ClearSelected()
			capacityWidget.Disable()
			colorWidget.Options = []string{}
			colorWidget.ClearSelected()
			colorWidget.Disable()
		} else {
			// 成功加载数据，更新型号选择器
			modelSet := make(map[string]bool)
			for _, products := range services.Product.GetDynamicProducts() {
				for _, p := range products {
					// 从Title解析型号
					parts := strings.Split(p.Title, " ")
					if len(parts) >= 2 {
						for i, part := range parts {
							if strings.HasSuffix(part, "GB") || strings.HasSuffix(part, "TB") || strings.HasSuffix(part, "mm") {
								// 容量/尺寸之前的部分是型号
								model := strings.Join(parts[:i], " ")
								if model != "" {
									modelSet[model] = true
								}
								break
							}
						}
					}
				}
			}

			models := []string{}
			for model := range modelSet {
				models = append(models, model)
			}

			modelWidget.Options = models
			modelWidget.PlaceHolder = "选择型号"
			modelWidget.Enable()
		}

		// 根据地区显示不同的选择器
		if value == "中国大陆" {
			// 确保门店数据已加载
			selectedArea := services.Area.GetArea("中国大陆")
			if err := services.Store.LoadForArea(selectedArea.ShortCode); err != nil {
				log.Printf("Failed to load store data for %s: %v", selectedArea.ShortCode, err)
			}

			// 显示省份选择器和门店选择器
			provinceWidget.Options = model.GetProvinces()
			storeWidget.Options = []string{} // 初始为空，选择省份后更新
			storeWidget.Disable()
			storeWidget.PlaceHolder = "请先选择省份"
			locationContainer.Objects = []fyne.CanvasObject{
				container.NewGridWithColumns(2,
					container.NewVBox(
						widget.NewLabelWithStyle("省份", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
						provinceWidget,
					),
					container.NewVBox(
						widget.NewLabelWithStyle("门店", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
						storeWidget,
					),
				),
			}
		} else if value == "香港" {
			// 香港：直接显示所有门店
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
		} else if value == "日本" {
			// 日本：直接显示所有门店
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
		} else if value == "美国" {
			// 美国：先选择州，再选择门店
			zipCodeEntry.Disable()
			zipCodeEntry.Text = ""
			provinceWidget.Options = services.Store.GetStatesForArea(value)
			storeWidget.Options = []string{} // 初始为空，选择州后更新
			storeWidget.Disable()
			storeWidget.PlaceHolder = "请先选择州"
			locationContainer.Objects = []fyne.CanvasObject{
				container.NewGridWithColumns(2,
					container.NewVBox(
						widget.NewLabelWithStyle("州", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
						provinceWidget,
					),
					container.NewVBox(
						widget.NewLabelWithStyle("门店", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
						storeWidget,
					),
				),
			}
		} else if value == "英国" {
			// 英国：直接显示所有门店
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
		} else if value == "澳大利亚" {
			// 澳大利亚：直接显示所有门店
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
		locationContainer.Refresh()

		// 清空型号、容量、颜色选择
		modelWidget.ClearSelected()
		capacityWidget.Options = []string{}
		capacityWidget.ClearSelected()
		capacityWidget.Disable()
		colorWidget.Options = []string{}
		colorWidget.ClearSelected()
		colorWidget.Disable()
	})
	areaWidget.SetSelected(defaultArea)
	areaWidget.PlaceHolder = "选择地区"

	// 创建精简的标题栏
	titleLabel := widget.NewLabelWithStyle("抢你妹 - iPhone 库存监控 (支持 iPhone 17/17 Pro/17 Pro Max/Air)",
		fyne.TextAlignLeading, fyne.TextStyle{Bold: true})
	versionLabel := widget.NewLabel("v" + common.VERSION)
	authorLink := widget.NewHyperlink("@Sunbelife", parseURL("https://weibo.com/x1nyang"))
	ibetaLink := widget.NewHyperlink("@iBeta", parseURL("https://ibeta.me"))

	// Bark 通知 URL 输入框
	barkEntry := widget.NewEntry()
	barkEntry.SetPlaceHolder("输入Bark推送地址（如：https://api.day.app/your_device_key）")

	// 加载保存的设置
	if settings, err := services.LoadSettings(); err == nil {
		if settings.BarkNotifyUrl != "" {
			barkEntry.SetText(settings.BarkNotifyUrl)
			services.Listen.SetBarkUrl(settings.BarkNotifyUrl)
		}
		if settings.SelectedArea != "" {
			areaWidget.SetSelected(settings.SelectedArea)
		}
		// 恢复监听项
		if len(settings.ListenItems) > 0 {
			services.Listen.SetListenItems(settings.ListenItems)
		}
	}

	// 创建操作按钮组
	addButton := widget.NewButton("添加到监控列表", func() {
		// 检查必填项
		var errorMsg string
		if areaWidget.Selected == "中国大陆" || areaWidget.Selected == "美国" {
			if currentProvince == "" || storeWidget.Selected == "" {
				errorMsg = "请选择州/省份和门店"
			}
		} else if areaWidget.Selected == "香港" || areaWidget.Selected == "日本" || areaWidget.Selected == "新加坡" || areaWidget.Selected == "英国" || areaWidget.Selected == "澳大利亚" {
			if storeWidget.Selected == "" {
				errorMsg = "请选择门店"
			}
		}

		if errorMsg == "" && (modelWidget.Selected == "" || capacityWidget.Selected == "" || colorWidget.Selected == "") {
			errorMsg = "请完整选择型号、容量/尺寸和颜色"
		}

		if errorMsg != "" {
			dialog.ShowError(errors.New(errorMsg), view.Window)
			return
		}

		// 更新 Bark URL
		barkUrl := strings.TrimSpace(barkEntry.Text)
		if barkUrl != "" {
			// 简单的Bark URL验证
			if !strings.HasPrefix(barkUrl, "https://") {
				dialog.ShowError(errors.New("Bark URL必须以https://开头"), view.Window)
				return
			}
			if !strings.Contains(barkUrl, "api.day.app") && !strings.Contains(barkUrl, "bark") {
				dialog.ShowError(errors.New("请输入有效的Bark推送地址"), view.Window)
				return
			}
		}
		services.Listen.SetBarkUrl(barkUrl)

		// 只使用动态数据
		var productCode, productType string
		productTitle := ""

		// 从动态数据获取产品代码
		if services.Product.GetDynamicProducts() == nil || len(services.Product.GetDynamicProducts()) == 0 {
			dialog.ShowError(errors.New("产品数据未加载，请点击「更新数据」按钮获取最新产品信息"), view.Window)
			return
		}

		// 从动态数据中查找匹配的产品
		found := false
		for _, products := range services.Product.GetDynamicProducts() {
			for _, p := range products {
				// 尝试匹配产品标题
				expectedTitle := modelWidget.Selected + " " + capacityWidget.Selected + " " + colorWidget.Selected
				if p.Title == expectedTitle {
					productCode = p.Code
					productType = p.Type
					productTitle = p.Title
					found = true
					break
				}
			}
			if found {
				break
			}
		}

		if !found {
			dialog.ShowError(errors.New("未找到对应的产品代码，请点击「更新数据」按钮更新产品信息"), view.Window)
			return
		}

		if productCode != "" {
			// 根据地区类型添加监控项
			if areaWidget.Selected == "中国大陆" {
				// 中国大陆：从动态数据中获取门店信息
				selectedStore := services.Store.GetStore(areaWidget.Selected, storeWidget.Selected)
				services.Listen.AddWithStoreInfo(selectedStore, productTitle, productCode, productType)
			} else if areaWidget.Selected == "香港" || areaWidget.Selected == "日本" || areaWidget.Selected == "新加坡" ||
				areaWidget.Selected == "美国" || areaWidget.Selected == "英国" || areaWidget.Selected == "澳大利亚" {
				// 确保门店数据已加载
				selectedArea := services.Area.GetArea(areaWidget.Selected)
				if err := services.Store.LoadForArea(selectedArea.ShortCode); err != nil {
					log.Printf("Failed to load store data for %s: %v", areaWidget.Selected, err)
					dialog.ShowError(fmt.Errorf("加载门店数据失败: %v", err), view.Window)
					return
				}

				// 从动态数据中获取门店信息
				selectedStore := services.Store.GetStore(areaWidget.Selected, storeWidget.Selected)
				if selectedStore.StoreNumber == "" {
					log.Printf("Store not found: %s in %s", storeWidget.Selected, areaWidget.Selected)
					dialog.ShowError(fmt.Errorf("未找到门店: %s", storeWidget.Selected), view.Window)
					return
				}

				services.Listen.AddWithStoreInfo(selectedStore, productTitle, productCode, productType)
			}

			// 保存设置
			settings := services.UserSettings{
				SelectedArea:    areaWidget.Selected,
				SelectedStore:   storeWidget.Selected,
				SelectedProduct: productTitle,
				BarkNotifyUrl:   barkEntry.Text,
				ListenItems:     services.Listen.GetListenItems(),
			}
			services.SaveSettings(settings)
		}
	})
	addButton.Importance = widget.HighImportance

	clearButton := widget.NewButton("清空列表", func() {
		services.Listen.Clean()
		services.ClearSettings()
	})

	testSoundButton := widget.NewButton("测试提示音", func() {
		go services.Listen.AlertMp3()
	})

	// 更新数据按钮
	updateDataButton := widget.NewButton("更新数据", func() {
		progressDialog := dialog.NewCustom("更新中", "关闭",
			container.NewVBox(
				widget.NewProgressBarInfinite(),
				widget.NewLabel("正在从Apple官网获取最新数据..."),
			), view.Window)
		progressDialog.Show()

		go func() {
			var productErr, storeErr error

			// 更新所有地区的产品数据
			areaCodes := []string{"cn", "hk", "jp", "sg", "us", "uk", "au"}
			for _, areaCode := range areaCodes {
				log.Printf("Updating product data for %s...", areaCode)
				if err := services.UpdateProductDatabase(areaCode); err != nil {
					log.Printf("Failed to update product data for %s: %v", areaCode, err)
					productErr = err
				}
				// 添加延迟避免频繁请求
				time.Sleep(time.Duration(2+rand.Intn(2)) * time.Second)
			}

			// 更新所有地区的门店数据
			log.Println("Updating store data for all areas...")
			if err := services.UpdateStoresForAllAreas(); err != nil {
				log.Printf("Failed to update store data: %v", err)
				storeErr = err
			}

			progressDialog.Hide()

			if productErr != nil || storeErr != nil {
				dialog.ShowError(errors.New("部分数据更新失败，请检查网络连接后重试"), view.Window)
			} else {
				dialog.ShowInformation("成功", "所有数据已更新完成", view.Window)

				// 重新加载当前地区的门店数据到内存
				if areaWidget.Selected != "" {
					selectedArea := services.Area.GetArea(areaWidget.Selected)
					services.Store.LoadForArea(selectedArea.ShortCode)
				}

				// 更新型号选择器选项
				if services.Product.GetDynamicProducts() != nil {
					modelSet := make(map[string]bool)
					for _, products := range services.Product.GetDynamicProducts() {
						for _, p := range products {
							// 从Title解析型号
							parts := strings.Split(p.Title, " ")
							if len(parts) >= 2 {
								for i, part := range parts {
									if strings.HasSuffix(part, "GB") || strings.HasSuffix(part, "TB") || strings.HasSuffix(part, "mm") {
										model := strings.Join(parts[:i], " ")
										if model != "" {
											modelSet[model] = true
										}
										break
									}
								}
							}
						}
					}

					var models []string
					for model := range modelSet {
						models = append(models, model)
					}

					modelWidget.Options = models
					modelWidget.Refresh()
				}
			}
		}()
	})
	updateDataButton.Importance = widget.WarningImportance

	// 控制按钮组
	startButton := widget.NewButton("开始监控", func() {
		services.Listen.Status.Set(services.Running)
	})
	startButton.Importance = widget.HighImportance

	pauseButton := widget.NewButton("暂停监控", func() {
		services.Listen.Status.Set(services.Pause)
	})

	// 创建日志滚动容器（精简高度）
	logScroll := container.NewScroll(services.Listen.Logs)
	logScroll.SetMinSize(fyne.NewSize(200, 200))

	// 第一列：地区和门店选择
	column1 := widget.NewCard("配置选择", "", container.NewVBox(
		container.NewVBox(
			widget.NewLabelWithStyle("地区", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
			areaWidget,
		),
		widget.NewSeparator(),
		locationContainer, // 动态位置选择容器
	))

	// 第二列：产品选择和控制
	column2 := widget.NewCard("产品与控制", "", container.NewVBox(
		container.NewGridWithColumns(2,
			container.NewVBox(
				widget.NewLabelWithStyle("型号", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
				modelWidget,
			),
			container.NewVBox(
				capacityLabel,
				capacityWidget,
			),
		),
		container.NewVBox(
			colorLabel,
			colorWidget,
		),
		widget.NewSeparator(),
		container.NewVBox(
			widget.NewLabelWithStyle("Bark推送（可选）", fyne.TextAlignLeading, fyne.TextStyle{}),
			barkEntry,
		),
		widget.NewSeparator(),
		container.NewGridWithColumns(2,
			addButton,
			clearButton,
		),
		container.NewGridWithColumns(2,
			testSoundButton,
			updateDataButton,
		),
		widget.NewSeparator(),
		container.NewGridWithColumns(2,
			startButton,
			pauseButton,
		),
		container.NewHBox(
			widget.NewLabelWithStyle("状态：", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
			widget.NewLabelWithData(services.Listen.Status),
		),
	))

	// 第三列：监控日志
	column3 := widget.NewCard("监控日志", "", logScroll)

	// 创建三列布局
	threeColumns := container.NewGridWithColumns(3,
		column1,
		column2,
		column3,
	)

	// 顶部栏
	topBar := container.NewBorder(
		nil, widget.NewSeparator(), nil,
		container.NewHBox(
			versionLabel,
			widget.NewLabel("·"),
			authorLink,
			widget.NewLabel("·"),
			ibetaLink,
		),
		titleLabel,
	)

	// 主布局（更紧凑）
	content := container.NewBorder(
		topBar,
		nil, nil, nil,
		threeColumns,
	)

	// 添加少量内边距
	paddedContent := container.NewPadded(content)

	view.Window.SetContent(paddedContent)
	view.Window.Resize(fyne.NewSize(1400, 800))
	view.Window.CenterOnScreen()

	services.Listen.Run()
	view.Window.ShowAndRun()
}

func parseURL(urlStr string) *url.URL {
	link, _ := url.Parse(urlStr)
	return link
}
