package main

import (
	"apple-store-helper/common"
	"apple-store-helper/services"
	"apple-store-helper/theme"
	"apple-store-helper/view"
	"errors"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
	"github.com/faiface/beep"
	"github.com/faiface/beep/speaker"
	"time"
)

func main() {
	// 初始 mp3 播放器
	SampleRate := beep.SampleRate(44100)
	speaker.Init(SampleRate, SampleRate.N(time.Second/10))

	view.App = app.NewWithID("apple-store-helper")
	view.App.Settings().SetTheme(&theme.MyTheme{})
	view.Window = view.App.NewWindow("Apple Store 抢你妹")

	defaultArea := services.Listen.Area.Title

	// 门店 selector
	storeWidget := widget.NewSelect(services.Store.ByAreaTitleForOptions(defaultArea), nil)
	storeWidget.PlaceHolder = "请选择自提门店"

	// 型号 selector
	productWidget := widget.NewSelect(services.Product.ByAreaTitleForOptions(defaultArea), nil)
	productWidget.PlaceHolder = "请选择手机"

	// 地区 selector

	areaWidget := widget.NewRadioGroup(services.Area.ForOptions(), func(value string) {
		storeWidget.Options = services.Store.ByAreaTitleForOptions(value)
		storeWidget.ClearSelected()

		productWidget.Options = services.Product.ByAreaTitleForOptions(value)
		productWidget.ClearSelected()

		services.Listen.Area = services.Area.GetArea(value)
		// services.Listen.Clean()
	})
	areaWidget.SetSelected(defaultArea)
	areaWidget.Horizontal = true
	help := `1. 提前将需要购买的型号加入购物车，检测有货会打开购物车页面，需要在购物车页面手动选择门店
2. 选择门店和型号，点击 添加 按钮
3. 点击 开始
4. 抢到了的话记得去微博关注 @Sunbelife，因为这个工具是他适配的
`

	view.Window.SetContent(container.NewVBox(
		widget.NewLabel(help),
		container.New(layout.NewFormLayout(), widget.NewLabel("选择地区:"), areaWidget),
		container.New(layout.NewFormLayout(), widget.NewLabel("选择门店:"), storeWidget),
		container.New(layout.NewFormLayout(), widget.NewLabel("选择型号:"), productWidget),

		container.NewBorder(nil, nil,
			container.NewHBox(
				widget.NewButton("添加", func() {
					if storeWidget.Selected == "" || productWidget.Selected == "" {
						dialog.ShowError(errors.New("请选择门店和型号"), view.Window)
					} else {
						services.Listen.Add(areaWidget.Selected, storeWidget.Selected, productWidget.Selected)
					}
				}),
				widget.NewButton("清空", func() {
					services.Listen.Clean()
				}),
				widget.NewButton("试听(有货提示音)", func() {
					go services.Listen.AlertMp3()
				}),
			),
			container.NewHBox(
				widget.NewButton("开始", func() {
					services.Listen.Status.Set(services.Running)
				}),
				widget.NewButton("暂停", func() {
					services.Listen.Status.Set(services.Pause)
				}),
				container.NewCenter(widget.NewLabel("状态:")),
				container.NewCenter(widget.NewLabelWithData(services.Listen.Status)),
			),
		),
		services.Listen.Logs,
		layout.NewSpacer(),
		container.NewHBox(
			layout.NewSpacer(),
			widget.NewLabel("version: "+common.VERSION),
		),
	))

	view.Window.Resize(fyne.NewSize(1000, 800))
	services.Listen.Run()
	view.Window.ShowAndRun()
}
