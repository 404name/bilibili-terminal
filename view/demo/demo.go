package demo

import (
	"fmt"

	"github.com/404name/termui-demo/model"
	"github.com/404name/termui-demo/utils"
	ui "github.com/gizak/termui/v3"
	"github.com/gizak/termui/v3/widgets"
)

type DemoPage struct {
	Layout          *ui.Grid
	Ready           bool
	Top             *widgets.Paragraph
	Tab             *widgets.TabPane
	VideoList       *widgets.List
	Bottom          *widgets.Paragraph
	TabView         []interface{}
	VideoListDetail []model.RcmdVideo
}

func NewPage(layout *ui.Grid) *DemoPage {

	v := &DemoPage{}
	v.Layout = layout
	v.Ready = false
	v.Top = widgets.NewParagraph()
	v.Top.Title = "我是顶部"
	v.Top.Text = "我是顶部"
	// list列表
	v.VideoList = widgets.NewList()
	v.VideoList.TextStyle = ui.NewStyle(ui.ColorYellow)
	v.VideoList.WrapText = false

	// tab栏
	v.Tab = widgets.NewTabPane("推荐", "热门", "动态", "我的")
	v.Tab.Border = true
	for _, name := range v.Tab.TabNames {
		p := widgets.NewParagraph()
		p.Text = "我是" + name + "页面"
		v.TabView = append(v.TabView, p)
	}

	// 日志输出
	v.Bottom = widgets.NewParagraph()
	v.Bottom.Text = utils.VersionMessage()
	v.Bottom.Title = "系统日志"
	v.Bottom.BorderStyle.Fg = ui.ColorBlue
	return v
}

func (v *DemoPage) Load() error {
	v.Ready = false
	defer func() {
		v.Ready = true
	}()
	// 初始化加载
	rcmdVideoList, err := model.GetRcmdVideo()
	v.VideoListDetail = rcmdVideoList
	if err != nil {
		return err
	}
	for _, video := range rcmdVideoList {
		v.VideoList.Rows = append(v.VideoList.Rows, video.Title)
	}

	ui.Render(v.VideoList)
	return nil
}

func (v *DemoPage) IsReady() bool {
	return true
}

func (v *DemoPage) EventHander(e ui.Event) {
	switch e.ID {
	case "<Space>", "Enter":
		if v.IsReady() {
			v.VideoList.Rows = []string{}
			v.Bottom.Text = fmt.Sprintf("清空列表")
			ui.Render(v.Bottom)
			v.Load()
			v.Bottom.Text = fmt.Sprintf("加载列表%v", v.VideoList.Rows)
			ui.Render(v.Bottom)
		}
	case "<Left>":
		v.Tab.FocusLeft()
		ui.Render(v.Tab)
		ui.Render(v.TabView[v.Tab.ActiveTabIndex].(ui.Drawable))
	case "<Right>":
		v.Tab.FocusRight()
		ui.Render(v.Tab)
		ui.Render(v.TabView[v.Tab.ActiveTabIndex].(ui.Drawable))
	case "<Up>":
		v.VideoList.ScrollUp()
		if v.VideoList.SelectedRow != -1 {
			v.Bottom.Text = fmt.Sprintf("%v", v.VideoListDetail[v.VideoList.SelectedRow])
		}
		ui.Render(v.VideoList)
		ui.Render(v.Bottom)
	case "<Down>":
		v.VideoList.ScrollDown()
		if v.VideoList.SelectedRow != -1 {
			v.Bottom.Text = fmt.Sprintf("%v", v.VideoListDetail[v.VideoList.SelectedRow])
		}

		ui.Render(v.VideoList)
		ui.Render(v.Bottom)
	}
}

func (v *DemoPage) Close() {
	// 关闭里面在执行的事件
}
func (v *DemoPage) Refresh() {
	termWidth, termHeight := ui.TerminalDimensions()

	v.Layout.SetRect(0, 0, termWidth, termHeight)

	v.Layout.Set(
		ui.NewRow(1.0/10, v.Top),
		ui.NewRow(1.0/10, v.Tab),
		ui.NewRow(1.0/10, v.TabView[v.Tab.ActiveTabIndex]),
		ui.NewRow(5.0/10, v.VideoList),
		ui.NewRow(2.0/10, v.Bottom),
	)
	ui.Render(v.Layout)
}
