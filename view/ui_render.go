package view

import "github.com/gizak/termui/v3"

// 一个页面要提供
// @Refresh   刷新方法
// @Close     页面关闭方法
// @EventHander() 事件处理器(响应事件)
type Page interface {
	Refresh()
	Close()
	EventHander(e termui.Event)
}

var NowPage Page

func InitUI() {
	NowPage = NewVideoPage()
	NowPage.Refresh()
}
