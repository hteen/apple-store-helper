package main

import (
	"apple-store-helper/common"
	"apple-store-helper/services"
	"apple-store-helper/theme"
	"errors"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
)

func main() {

	a := app.NewWithID("ip13")
	a.Settings().SetTheme(&theme.MyTheme{})
	w := a.NewWindow("iPhone13|Mini|Pro|ProMax")

	defaultArea := services.Listen.Area.Title

	// 门店 selector
	storeWidget := widget.NewSelect(services.Store.ByAreaTitleForOptions(defaultArea), nil)
	storeWidget.PlaceHolder = "请选择自提门店"

	// 型号 selector
	productWidget := widget.NewSelect(services.Product.ByAreaTitleForOptions(defaultArea), nil)
	productWidget.PlaceHolder = "请选择 iPhone 型号"

	// 地区 selector
	areaWidget := widget.NewRadioGroup(services.Area.ForOptions(), func(value string) {
		storeWidget.Options = services.Store.ByAreaTitleForOptions(value)
		storeWidget.ClearSelected()

		productWidget.Options = services.Product.ByAreaTitleForOptions(value)
		productWidget.ClearSelected()

		services.Listen.Area = services.Area.GetArea(value)
	})
	areaWidget.SetSelected(defaultArea)
	areaWidget.Horizontal = true

	w.SetContent(container.NewVBox(
		widget.NewLabel("1.选择门店和型号，点击添加按钮\n"+
			"2.点击开始\n"+
			"3.匹配到之后会直接进入产品预购页面，选择预约门店和时间",
		),
		container.New(layout.NewFormLayout(), widget.NewLabel("选择地区:"), areaWidget),
		container.New(layout.NewFormLayout(), widget.NewLabel("选择门店:"), storeWidget),
		container.New(layout.NewFormLayout(), widget.NewLabel("选择型号:"), productWidget),

		container.NewBorder(nil, nil,
			container.NewHBox(
				widget.NewButton("添加", func() {
					if storeWidget.Selected == "" || productWidget.Selected == "" {
						dialog.ShowError(errors.New("请选择门店和型号"), w)
					} else {
						services.Listen.Add(areaWidget.Selected, storeWidget.Selected, productWidget.Selected)
					}
				}),
				widget.NewButton("清空", func() {
					services.Listen.Clean()
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

	w.Resize(fyne.NewSize(1000, 800))
	services.Listen.Run()
	services.Listen.Window = w
	w.ShowAndRun()
}
