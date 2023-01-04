package view

import (
	"strings"

	"github.com/404name/termui-demo/view/demo"
	"github.com/404name/termui-demo/view/index"
	"github.com/404name/termui-demo/view/video"
	ui "github.com/gizak/termui/v3"
)

// 这边应该一个全局的界面
// 接收里面的err等内容

// 一个页面要提供
// @Refresh   刷新方法
// @Close     页面关闭方法
// @EventHander() 事件处理器(响应事件)
type Page interface {
	Load() error //初始化
	Refresh()
	Close()
	EventHander(e ui.Event)
}

var BiliUI *BilibiliUI

func Init() {
	BiliUI = &BilibiliUI{}
	BiliUI.Layout = ui.NewGrid()
	BiliUI.NavigateTo("IndexPage", "")
}

type PageList struct {
	IndexPage *index.IndexPage
	VideoPage *video.VideoPage
	DemoPage  *demo.DemoPage
}

type BilibiliUI struct {
	NowPage Page
	Layout  *ui.Grid
	PageList
}

func (app *BilibiliUI) NavigateTo(toPage string, param string) {
	if app.NowPage != nil {
		// 释放当前页面资源
		app.NowPage.Close()
	}
	switch toPage {
	case "VideoPage":
		if app.VideoPage == nil {
			app.VideoPage = video.NewPage(app.Layout)
		}
		if params := strings.Split(param, "&"); len(params) >= 2 {
			app.VideoPage.Bvid = params[0]
			app.VideoPage.Cid = params[1]
		}
		app.loadPage(app.VideoPage)
	case "IndexPage":
		if app.IndexPage == nil {
			app.IndexPage = index.NewPage(app.Layout)
		}
		app.loadPage(app.IndexPage)
	}
}

func (app *BilibiliUI) loadPage(p Page) {
	app.NowPage = p
	app.NowPage.Load() // 后面这里要捕获错误
	app.NowPage.Refresh()
}

func (app *BilibiliUI) NavigateBack() {
	// 暂时返回首页
	if app.NowPage != app.IndexPage {
		app.NavigateTo("IndexPage", "")
	}
}

func (app *BilibiliUI) Refresh() {
	app.NowPage.Refresh()
}

func (app *BilibiliUI) EventHander(e ui.Event) {
	app.NowPage.EventHander(e)
}
