package main

import (
	"apple-store-helper/common"
	"apple-store-helper/services"
	"apple-store-helper/theme"
	"apple-store-helper/view"
	"errors"
	"fmt"
	"strconv"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
	"github.com/faiface/beep"
	"github.com/faiface/beep/speaker"
)

// main 主函数 (Main function)
func main() {
	initMP3Player()
	initFyneApp()

	// 默认地区 (Default Area)
	defaultArea := services.Listen.Area.Title

	// 门店选择器 (Store Selector)
	storeWidget := widget.NewSelect(services.Store.ByAreaTitleForOptions(defaultArea), nil)
	storeWidget.PlaceHolder = "请选择自提门店"

	// 型号选择器 (Product Selector)
	productWidget := widget.NewSelect(services.Product.ByAreaTitleForOptions(defaultArea), nil)
	productWidget.PlaceHolder = "请选择 iPhone 型号"

	// Bark 通知输入框
	barkWidget := widget.NewEntry()
	barkWidget.SetPlaceHolder("https://api.day.app/你的BarkKey")

	// 检测次数阈值设置
	detectThresholdWidget := widget.NewEntry()
	detectThresholdWidget.SetText("3")
	detectThresholdWidget.SetPlaceHolder("检测次数阈值")

	// 时间阈值设置（分钟）
	timeThresholdWidget := widget.NewEntry()
	timeThresholdWidget.SetText("1")
	timeThresholdWidget.SetPlaceHolder("时间阈值（分钟）")

	// 刷新频率设置（秒）
	refreshIntervalWidget := widget.NewEntry()
	refreshIntervalWidget.SetText("3")
	refreshIntervalWidget.SetPlaceHolder("刷新频率（秒）")

	// 批次间隔设置（毫秒）
	batchIntervalWidget := widget.NewEntry()
	batchIntervalWidget.SetText("500")
	batchIntervalWidget.SetPlaceHolder("批次间隔（毫秒）")

	// 地区选择器 (Area Selector)
	areaWidget := widget.NewRadioGroup(services.Area.ForOptions(), nil)

	areaWidget.Horizontal = true

	// 为配置项添加变化监听
	areaWidget.OnChanged = func(value string) {
		// 防止空值或无效值导致崩溃
		if value == "" {
			return
		}

		// 更新商店列表
		storeWidget.Options = services.Store.ByAreaTitleForOptions(value)
		storeWidget.ClearSelected()

		// 更新产品列表
		productWidget.Options = services.Product.ByAreaTitleForOptions(value)
		productWidget.ClearSelected()

		// 更新监听服务
		services.Listen.Area = services.Area.GetArea(value)
		services.Listen.Clean()

		// 保存配置
		go saveCurrentSettings(areaWidget, storeWidget, productWidget, barkWidget, detectThresholdWidget, timeThresholdWidget, refreshIntervalWidget, batchIntervalWidget)
	}
	storeWidget.OnChanged = func(value string) {
		go saveCurrentSettings(areaWidget, storeWidget, productWidget, barkWidget, detectThresholdWidget, timeThresholdWidget, refreshIntervalWidget, batchIntervalWidget)
	}
	productWidget.OnChanged = func(value string) {
		go saveCurrentSettings(areaWidget, storeWidget, productWidget, barkWidget, detectThresholdWidget, timeThresholdWidget, refreshIntervalWidget, batchIntervalWidget)
	}
	barkWidget.OnChanged = func(value string) {
		go saveCurrentSettings(areaWidget, storeWidget, productWidget, barkWidget, detectThresholdWidget, timeThresholdWidget, refreshIntervalWidget, batchIntervalWidget)
	}
	detectThresholdWidget.OnChanged = func(value string) {
		go saveCurrentSettings(areaWidget, storeWidget, productWidget, barkWidget, detectThresholdWidget, timeThresholdWidget, refreshIntervalWidget, batchIntervalWidget)
	}
	timeThresholdWidget.OnChanged = func(value string) {
		go saveCurrentSettings(areaWidget, storeWidget, productWidget, barkWidget, detectThresholdWidget, timeThresholdWidget, refreshIntervalWidget, batchIntervalWidget)
	}
	refreshIntervalWidget.OnChanged = func(value string) {
		go saveCurrentSettings(areaWidget, storeWidget, productWidget, barkWidget, detectThresholdWidget, timeThresholdWidget, refreshIntervalWidget, batchIntervalWidget)
	}
	batchIntervalWidget.OnChanged = func(value string) {
		go saveCurrentSettings(areaWidget, storeWidget, productWidget, barkWidget, detectThresholdWidget, timeThresholdWidget, refreshIntervalWidget, batchIntervalWidget)
	}

	help := `1. 在 Apple 官网将需要购买的型号加入购物车
2. 选择地区、门店和型号，点击"添加"按钮，将需要监听的型号添加到监听列表
3. 设置检测阈值：检测次数阈值（默认3次）和时间阈值（默认1分钟）
4. 设置刷新频率（默认3秒），可根据需要调整监听频率
5. 设置批次间隔（默认500毫秒），控制同一轮内门店请求的间隔，避免API访问过于频繁
6. 点击"开始"按钮开始监听，在指定时间内检测到指定次数有货时会自动打开购物车页面
`

	loadUserSettingsCache(areaWidget, storeWidget, productWidget, barkWidget, detectThresholdWidget, timeThresholdWidget, refreshIntervalWidget, batchIntervalWidget)

	// 自动保存当前配置，确保所有配置项都被缓存
	go func() {
		time.Sleep(100 * time.Millisecond) // 等待UI初始化完成
		saveCurrentSettings(areaWidget, storeWidget, productWidget, barkWidget, detectThresholdWidget, timeThresholdWidget, refreshIntervalWidget, batchIntervalWidget)
	}()

	// 初始化 GUI 窗口内容 (Initialize GUI)
	view.Window.SetContent(container.NewVBox(
		widget.NewLabel(help),
		container.New(layout.NewFormLayout(), widget.NewLabel("选择地区:"), areaWidget),
		container.New(layout.NewFormLayout(), widget.NewLabel("选择门店:"), storeWidget),
		container.New(layout.NewFormLayout(), widget.NewLabel("选择型号:"), productWidget),
		container.New(layout.NewFormLayout(), widget.NewLabel("Bark 通知地址"), barkWidget),
		container.New(layout.NewFormLayout(), widget.NewLabel("检测次数阈值:"), detectThresholdWidget),
		container.New(layout.NewFormLayout(), widget.NewLabel("时间阈值(分钟):"), timeThresholdWidget),
		container.New(layout.NewFormLayout(), widget.NewLabel("刷新频率(秒):"), refreshIntervalWidget),
		container.New(layout.NewFormLayout(), widget.NewLabel("批次间隔(毫秒):"), batchIntervalWidget),

		container.NewBorder(nil, nil,
			createActionButtons(areaWidget, storeWidget, productWidget, barkWidget, detectThresholdWidget, timeThresholdWidget, refreshIntervalWidget, batchIntervalWidget),
			createControlButtons(),
		),

		services.Listen.Logs,
		layout.NewSpacer(),
		createVersionLabel(),
	))

	view.Window.Resize(fyne.NewSize(1000, 800))
	view.Window.CenterOnScreen()
	services.Listen.Run()
	view.Window.ShowAndRun()
}

// initMP3Player 初始化 MP3 播放器 (Initialize MP3 player)
func initMP3Player() {
	SampleRate := beep.SampleRate(44100)
	speaker.Init(SampleRate, SampleRate.N(time.Second/10))
}

// initFyneApp 初始化 Fyne 应用 (Initialize Fyne App)
func initFyneApp() {
	view.App = app.NewWithID("apple-store-helper")
	view.App.Settings().SetTheme(&theme.MyTheme{})
	view.Window = view.App.NewWindow("Apple Store Helper")
}

// 加载用户设置缓存 (Load user settings cache)
func loadUserSettingsCache(areaWidget *widget.RadioGroup, storeWidget *widget.Select, productWidget *widget.Select, barkNotifyWidget *widget.Entry, detectThresholdWidget *widget.Entry, timeThresholdWidget *widget.Entry, refreshIntervalWidget *widget.Entry, batchIntervalWidget *widget.Entry) {
	settings, err := services.LoadSettings()
	if err == nil {
		areaWidget.SetSelected(settings.SelectedArea)
		storeWidget.SetSelected(settings.SelectedStore)
		productWidget.SetSelected(settings.SelectedProduct)
		services.Listen.SetListenItems(settings.ListenItems)
		barkNotifyWidget.SetText(settings.BarkNotifyUrl)

		// 加载阈值设置，如果为0则使用默认值
		if settings.DetectThreshold > 0 {
			detectThresholdWidget.SetText(fmt.Sprintf("%d", settings.DetectThreshold))
			services.Listen.DetectThreshold = settings.DetectThreshold
		}
		if settings.TimeThreshold > 0 {
			timeThresholdWidget.SetText(fmt.Sprintf("%d", settings.TimeThreshold))
			services.Listen.TimeThreshold = settings.TimeThreshold
		}
		if settings.RefreshInterval > 0 {
			refreshIntervalWidget.SetText(fmt.Sprintf("%d", settings.RefreshInterval))
			services.Listen.RefreshInterval = settings.RefreshInterval
		}
		if settings.BatchInterval > 0 {
			batchIntervalWidget.SetText(fmt.Sprintf("%d", settings.BatchInterval))
			services.Listen.BatchInterval = settings.BatchInterval
		}
	} else {
		areaWidget.SetSelected(services.Listen.Area.Title)
	}
}

// 保存当前配置 saveCurrentSettings saves current settings
func saveCurrentSettings(areaWidget *widget.RadioGroup, storeWidget *widget.Select, productWidget *widget.Select, barkNotifyWidget *widget.Entry, detectThresholdWidget *widget.Entry, timeThresholdWidget *widget.Entry, refreshIntervalWidget *widget.Entry, batchIntervalWidget *widget.Entry) {
	detectThreshold, _ := strconv.Atoi(detectThresholdWidget.Text)
	timeThreshold, _ := strconv.Atoi(timeThresholdWidget.Text)
	refreshInterval, _ := strconv.Atoi(refreshIntervalWidget.Text)
	batchInterval, _ := strconv.Atoi(batchIntervalWidget.Text)

	// 确保有合理的默认值
	if detectThreshold <= 0 {
		detectThreshold = 3
	}
	if timeThreshold <= 0 {
		timeThreshold = 1
	}
	if refreshInterval <= 0 {
		refreshInterval = 3
	}
	if batchInterval <= 0 {
		batchInterval = 500
	}

	services.SaveAllCurrentSettings(
		areaWidget.Selected,
		storeWidget.Selected,
		productWidget.Selected,
		barkNotifyWidget.Text,
		detectThreshold,
		timeThreshold,
		refreshInterval,
		batchInterval,
	)
}

// 创建动作按钮 (Create action buttons)
func createActionButtons(areaWidget *widget.RadioGroup, storeWidget *widget.Select, productWidget *widget.Select, barkNotifyWidget *widget.Entry, detectThresholdWidget *widget.Entry, timeThresholdWidget *widget.Entry, refreshIntervalWidget *widget.Entry, batchIntervalWidget *widget.Entry) *fyne.Container {
	return container.NewHBox(
		widget.NewButton("添加", func() {
			if storeWidget.Selected == "" || productWidget.Selected == "" {
				dialog.ShowError(errors.New("请选择门店和型号"), view.Window)
			} else {
				// 解析阈值设置
				detectThreshold, err1 := strconv.Atoi(detectThresholdWidget.Text)
				timeThreshold, err2 := strconv.Atoi(timeThresholdWidget.Text)
				refreshInterval, err3 := strconv.Atoi(refreshIntervalWidget.Text)
				batchInterval, err4 := strconv.Atoi(batchIntervalWidget.Text)

				if err1 != nil || detectThreshold <= 0 {
					dialog.ShowError(errors.New("检测次数阈值必须是正整数"), view.Window)
					return
				}
				if err2 != nil || timeThreshold <= 0 {
					dialog.ShowError(errors.New("时间阈值必须是正整数"), view.Window)
					return
				}
				if err3 != nil || refreshInterval <= 0 {
					dialog.ShowError(errors.New("刷新频率必须是正整数"), view.Window)
					return
				}
				if err4 != nil || batchInterval <= 0 {
					dialog.ShowError(errors.New("批次间隔必须是正整数"), view.Window)
					return
				}

				// 设置阈值、刷新频率和批次间隔
				services.Listen.SetThresholds(detectThreshold, timeThreshold)
				services.Listen.SetRefreshInterval(refreshInterval)
				services.Listen.SetBatchInterval(batchInterval)

				services.Listen.Add(areaWidget.Selected, storeWidget.Selected, productWidget.Selected, barkNotifyWidget.Text)

				// 保存所有配置
				services.SaveAllCurrentSettings(
					areaWidget.Selected,
					storeWidget.Selected,
					productWidget.Selected,
					barkNotifyWidget.Text,
					detectThreshold,
					timeThreshold,
					refreshInterval,
					batchInterval,
				)
			}
		}),
		widget.NewButton("清空", func() {
			services.Listen.Clean()
			services.ClearSettings()

			// 重置UI到默认值
			areaWidget.SetSelected(services.Listen.Area.Title)
			storeWidget.ClearSelected()
			productWidget.ClearSelected()
			barkNotifyWidget.SetText("")
			detectThresholdWidget.SetText("3")
			timeThresholdWidget.SetText("1")
			refreshIntervalWidget.SetText("3")
			batchIntervalWidget.SetText("500")

			// 重置服务配置到默认值
			services.Listen.DetectThreshold = 3
			services.Listen.TimeThreshold = 1
			services.Listen.RefreshInterval = 3
			services.Listen.BatchInterval = 500

			// 保存重置后的配置
			saveCurrentSettings(areaWidget, storeWidget, productWidget, barkNotifyWidget, detectThresholdWidget, timeThresholdWidget, refreshIntervalWidget, batchIntervalWidget)
		}),
		widget.NewButton("试听(有货提示音)", func() {
			go services.Listen.AlertMp3()
		}),
		widget.NewButton("测试 Bark 通知", func() {
			services.Listen.BarkNotifyUrl = barkNotifyWidget.Text
			services.Listen.SendPushNotificationByBark("有货提醒（测试）", "此为测试提醒，点击通知将跳转到相关链接", "https://www.apple.com.cn/shop/bag")
		}),
	)
}

// 创建控制按钮 (Create control buttons)
func createControlButtons() *fyne.Container {
	return container.NewHBox(
		widget.NewButton("开始", func() {
			services.Listen.Status.Set(services.Running)
		}),
		widget.NewButton("暂停", func() {
			services.Listen.Status.Set(services.Pause)
		}),
		container.NewCenter(widget.NewLabel("状态:")),
		container.NewCenter(widget.NewLabelWithData(services.Listen.Status)),
	)
}

// createVersionLabel 创建版本标签 (Create version label)
func createVersionLabel() *fyne.Container {
	return container.NewHBox(
		layout.NewSpacer(),
		widget.NewLabel("version: "+common.VERSION),
	)
}
