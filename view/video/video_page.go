package video

import (
	"github.com/404name/termui-demo/global"
	"github.com/404name/termui-demo/utils"

	ui "github.com/gizak/termui/v3"
	"github.com/gizak/termui/v3/widgets"
)

type VideoPage struct {
	VideoDetail
	Gauge   *widgets.Gauge
	List    *widgets.List
	P1      *widgets.Paragraph
	P2      *widgets.Paragraph
	Tab     *widgets.TabPane
	TabView []interface{}

	CurTabView interface{}
	Layout     *ui.Grid
}

func NewPage(layout *ui.Grid) *VideoPage {

	v := &VideoPage{}
	v.Layout = layout
	// 进度条
	v.ProgressBar = widgets.NewGauge()
	v.ProgressBar.BarColor = ui.ColorBlue
	v.Gauge = widgets.NewGauge()
	v.Gauge.BarColor = ui.ColorRed
	// list列表
	v.List = widgets.NewList()
	v.List.Rows = []string{
		"[1] Downloading File 1",
		"",
		"[2] Downloading File 2",
		"",
		"[3] Uploading File 3",
	}
	// 日志输出
	v.Log = widgets.NewParagraph()
	v.Log.Text = utils.VersionMessage()

	v.Log.Title = "系统日志"

	// 文字
	v.P1 = widgets.NewParagraph()
	v.P1.Text = `hello在这里。今生无悔入东方 来世愿生幻想乡
	红魔地灵夜神雪 永夜风神星莲船
	非想天则文花贴 萃梦神灵绯想天
	冥界地狱异变起 樱下华胥主谋现
	净罪无改渡黄泉 华鸟风月是非辨
	境界颠覆入迷途 幻想花开啸风弄
	二色花蝶双生缘 前缘未尽今生还 
	星屑洒落雨霖铃 虹彩彗光银尘耀
	无寿迷蝶彼岸归 幻真如画妖如月
	永劫夜宵哀伤起 幼社灵中幻似梦
	追忆往昔巫女缘 须弥之间冥梦现
	仁榀华诞井中天 歌雅风颂心无念
	此生无悔入东方 来世愿生幻想乡
	不求间隙一紫妹 但求回眸望幽香
	生死梦寄永远亭 孤魂永伴迷途林
	白玉楼前西行樱 花开无缘彼岸川
	何处觅得妖怪山 千年神恋绝不厌
	红魔地灵夜神雪 永夜风神星莲船
	非想天则文花贴 萃梦神灵[绯想天
	冥界地狱异变起 亡灵妖怪主谋现
	幻想乡内四季现 华鸟风月是非辨
	此生唯一是红白 黑白魔女串门来(^・ω・^ )`
	v.P1.Title = "中文测试"
	v.P2 = widgets.NewParagraph()
	v.P2.Title = "中文测试"
	v.P2.Text = "Press q to ther 中文测试nd things. Let's see how this thing renders out.\n    Here is a new paragraph and stuffs and things. There should be a tab indent at the beginning of the paragraph. Let's see if that worked as well."
	// P2.Text = utils.ChineseStringFormat("中文测试")
	v.P2.BorderStyle.Fg = ui.ColorBlue

	v.Img = widgets.NewImage(utils.LoadImg(global.CONFIG.BasePath.VideoCoverImg))

	// tab栏
	v.Tab = widgets.NewTabPane("首页", "视频", "动态", "我的")
	v.Tab.Border = true
	v.TabView = []interface{}{v.Img, v.List, v.List, v.P1}
	return v
}

func (v *VideoPage) Load() error {
	// 初始化加载
	// 初始化视频下载或者渲染
	return v.VideoDetail.Init()
}

func (v *VideoPage) IsReady() bool {
	return v.Ready
}

func (v *VideoPage) EventHander(e ui.Event) {
	switch e.ID {
	case "<Space>", "Enter":
		if v.Ready {
			var playControl interface{}
			v.PlayChan <- playControl
		}
	case "<Left>":
		v.Tab.FocusLeft()
		ui.Render(v.CurTabView.(ui.Drawable))
	case "<Right>":
		v.Tab.FocusRight()
		ui.Render(v.CurTabView.(ui.Drawable))
	}
}

func (v *VideoPage) Close() {
	// 关闭里面在执行的事件
	if v.Ready {
		var playControl interface{}
		v.CloseChan <- playControl
	}
	v.cancel()
}
func (v *VideoPage) Refresh() {
	ui.Clear()
	termWidth, termHeight := ui.TerminalDimensions()
	videoHeight := float64(termWidth * 4.0 / 16.0)
	videoHeightRate := videoHeight / float64(termHeight)
	// global.LOG.Infoln("系统窗口尺寸==>", termWidth, termHeight)
	// global.LOG.Infoln("视频窗口尺寸|高比率==>", termWidth, videoHeight, videoHeightRate)
	v.Layout.SetRect(0, 0, termWidth, termHeight)
	// 后期切换界面
	v.CurTabView = v.TabView[v.Tab.ActiveTabIndex]
	v.Layout.Set(
		ui.NewRow(1.0/10,
			ui.NewCol(1.0/4, v.List),
			ui.NewCol(2.0/4, v.Tab),
			ui.NewCol(1.0/4, v.List),
		),
		ui.NewRow(videoHeightRate, v.CurTabView),
		ui.NewRow(1.0/20, v.ProgressBar),
		// ui.NewRow(1.0/20, Gauges[1]),
		ui.NewRow(7.0/10-videoHeightRate-1.0/20,
			ui.NewCol(1.0/2,
				ui.NewRow(1.0/2, v.P2),
				ui.NewRow(1.0/2, v.Gauge),
			),
			ui.NewCol(1.0/2, v.P1),
		),
		ui.NewRow(2.0/10, v.Log),
	)
	ui.Render(v.Layout)
}
