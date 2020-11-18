package main

import (
	"image/color"

	"fyne.io/fyne"
	"fyne.io/fyne/theme"
)

type myTheme struct{}

// return bundled font resource
func (myTheme) TextFont() fyne.Resource     { return resourceFzhtkTtf }
func (myTheme) TextBoldFont() fyne.Resource { return resourceFzhtkTtf }

func (myTheme) BackgroundColor() color.Color      { return theme.DarkTheme().BackgroundColor() }
func (myTheme) ButtonColor() color.Color          { return theme.DarkTheme().ButtonColor() }
func (myTheme) DisabledButtonColor() color.Color  { return theme.DarkTheme().DisabledButtonColor() }
func (myTheme) IconColor() color.Color            { return theme.DarkTheme().IconColor() }
func (myTheme) DisabledIconColor() color.Color    { return theme.DarkTheme().DisabledIconColor() }
func (myTheme) HyperlinkColor() color.Color       { return theme.DarkTheme().HyperlinkColor() }
func (myTheme) TextColor() color.Color            { return theme.DarkTheme().TextColor() }
func (myTheme) DisabledTextColor() color.Color    { return theme.DarkTheme().DisabledTextColor() }
func (myTheme) HoverColor() color.Color           { return theme.DarkTheme().HoverColor() }
func (myTheme) PlaceHolderColor() color.Color     { return theme.DarkTheme().PlaceHolderColor() }
func (myTheme) PrimaryColor() color.Color         { return theme.DarkTheme().PrimaryColor() }
func (myTheme) FocusColor() color.Color           { return theme.DarkTheme().FocusColor() }
func (myTheme) ScrollBarColor() color.Color       { return theme.DarkTheme().ScrollBarColor() }
func (myTheme) ShadowColor() color.Color          { return theme.DarkTheme().ShadowColor() }
func (myTheme) TextSize() int                     { return theme.DarkTheme().TextSize() }
func (myTheme) TextItalicFont() fyne.Resource     { return theme.DarkTheme().TextItalicFont() }
func (myTheme) TextBoldItalicFont() fyne.Resource { return theme.DarkTheme().TextBoldItalicFont() }
func (myTheme) TextMonospaceFont() fyne.Resource  { return theme.DarkTheme().TextMonospaceFont() }
func (myTheme) Padding() int                      { return theme.DarkTheme().Padding() }
func (myTheme) IconInlineSize() int               { return theme.DarkTheme().IconInlineSize() }
func (myTheme) ScrollBarSize() int                { return theme.DarkTheme().ScrollBarSize() }
func (myTheme) ScrollBarSmallSize() int           { return theme.DarkTheme().ScrollBarSmallSize() }