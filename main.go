package main

import (
	"apple-store-helper/common"
	"apple-store-helper/services"
	"apple-store-helper/theme"
	"apple-store-helper/view"
	"errors"
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

	// 地区选择器 (Area Selector)
	areaWidget := widget.NewRadioGroup(services.Area.ForOptions(), func(value string) {
		storeWidget.Options = services.Store.ByAreaTitleForOptions(value)
		storeWidget.ClearSelected()

		productWidget.Options = services.Product.ByAreaTitleForOptions(value)
		productWidget.ClearSelected()

		services.Listen.Area = services.Area.GetArea(value)
		services.Listen.Clean()
	})

	areaWidget.Horizontal = true

	help := `1. 在 Apple 官网将需要购买的型号加入购物车
2. 选择地区、门店和型号，点击“添加”按钮，将需要监听的型号添加到监听列表
3. 点击“开始”按钮开始监听，检测到有货时会自动打开购物车页面
`

	loadUserSettingsCache(areaWidget, storeWidget, productWidget, barkWidget)

	// 初始化 GUI 窗口内容 (Initialize GUI)
	view.Window.SetContent(container.NewVBox(
		widget.NewLabel(help),
		container.New(layout.NewFormLayout(), widget.NewLabel("选择地区:"), areaWidget),
		container.New(layout.NewFormLayout(), widget.NewLabel("选择门店:"), storeWidget),
		container.New(layout.NewFormLayout(), widget.NewLabel("选择型号:"), productWidget),
		container.New(layout.NewFormLayout(), widget.NewLabel("Bark 通知地址"), barkWidget),

		container.NewBorder(nil, nil,
			createActionButtons(areaWidget, storeWidget, productWidget, barkWidget),
			createControlButtons(),
		),

		services.Listen.Logs,
		layout.NewSpacer(),
		createVersionLabel(),
	))

	view.Window.Resize(fyne.NewSize(1000, 800))
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
func loadUserSettingsCache(areaWidget *widget.RadioGroup, storeWidget *widget.Select, productWidget *widget.Select, barkNotifyWidget *widget.Entry) {
	settings, err := services.LoadSettings()
	if err == nil {
		areaWidget.SetSelected(settings.SelectedArea)
		storeWidget.SetSelected(settings.SelectedStore)
		productWidget.SetSelected(settings.SelectedProduct)
		services.Listen.SetListenItems(settings.ListenItems)
		barkNotifyWidget.Text = settings.BarkNotifyUrl
		services.Listen.BarkNotifyUrl = settings.BarkNotifyUrl
	} else {
		areaWidget.SetSelected(services.Listen.Area.Title)
	}
}

// 创建动作按钮 (Create action buttons)
func createActionButtons(areaWidget *widget.RadioGroup, storeWidget *widget.Select, productWidget *widget.Select, barkNotifyWidget *widget.Entry) *fyne.Container {
	return container.NewHBox(
		widget.NewButton("添加", func() {
			if storeWidget.Selected == "" || productWidget.Selected == "" {
				dialog.ShowError(errors.New("请选择门店和型号"), view.Window)
			} else {
				services.Listen.Add(areaWidget.Selected, storeWidget.Selected, productWidget.Selected, barkNotifyWidget.Text)
				services.SaveSettings(services.UserSettings{
					SelectedArea:    areaWidget.Selected,
					SelectedStore:   storeWidget.Selected,
					SelectedProduct: productWidget.Selected,
					BarkNotifyUrl:   barkNotifyWidget.Text,
					ListenItems:     services.Listen.GetListenItems(),
				})
			}
		}),
		widget.NewButton("清空", func() {
			services.Listen.Clean()
			services.ClearSettings()
		}),
		widget.NewButton("试听(有货提示音)", func() {
			go services.Listen.AlertMp3()
		}),
		widget.NewButton("测试 Bark 通知", func() {
			services.Listen.BarkNotifyUrl = barkNotifyWidget.Text
			services.Listen.SendPushNotificationByBark("有货提醒（测试）", "此为测试提醒，点击通知将跳转到相关链接", "https://www.apple.com.cn/")
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
