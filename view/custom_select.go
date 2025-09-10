package view

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/widget"
)

// CustomSelect 自定义选择框，支持更宽的下拉列表
type CustomSelect struct {
	widget.Select
}

// NewCustomSelect 创建自定义选择框
func NewCustomSelect(options []string, changed func(string)) *CustomSelect {
	s := &CustomSelect{}
	s.Options = options
	s.OnChanged = changed
	s.ExtendBaseWidget(s)
	return s
}

// MinSize 返回最小尺寸，确保有足够的宽度
func (s *CustomSelect) MinSize() fyne.Size {
	// 计算最长选项的宽度
	maxWidth := float32(200) // 默认最小宽度
	
	for _, option := range s.Options {
		// 估算文本宽度（每个中文字符约14像素，英文字符约7像素）
		width := float32(0)
		for _, ch := range option {
			if ch > 127 {
				width += 14 // 中文字符
			} else {
				width += 7 // 英文字符
			}
		}
		if width > maxWidth {
			maxWidth = width
		}
	}
	
	// 加上一些边距和下拉箭头的空间
	maxWidth += 60
	
	return fyne.NewSize(maxWidth, s.Select.MinSize().Height)
}