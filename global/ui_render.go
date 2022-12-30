package global

import (
	"github.com/404name/termui-demo/resource"
	"github.com/404name/termui-demo/utils"

	ui "github.com/gizak/termui/v3"
	"github.com/gizak/termui/v3/widgets"
)

func InitUI() {
	// 实心折线图
	sl := widgets.NewSparkline()
	sl.Data = utils.SinFloat64[:100]
	sl.LineColor = ui.ColorCyan
	sl.TitleStyle.Fg = ui.ColorWhite

	SparklineGroup = widgets.NewSparklineGroup(sl)
	SparklineGroup.Title = "Sparkline"

	// 打点折线图
	Plot = widgets.NewPlot()
	Plot.Title = "braille-mode Line Chart"
	Plot.Data = append(Plot.Data, utils.SinFloat64)
	Plot.AxesColor = ui.ColorWhite
	Plot.LineColors[0] = ui.ColorYellow

	// 进度条
	Gauges = make([]*widgets.Gauge, 3)
	for i := range Gauges {
		Gauges[i] = widgets.NewGauge()
		Gauges[i].Percent = 0
		Gauges[i].BarColor = ui.ColorRed
	}
	Gauges[0].BarColor = ui.ColorBlue
	// list列表
	List = widgets.NewList()
	List.Rows = []string{
		"[1] Downloading File 1",
		"",
		"[2] Downloading File 2",
		"",
		"[3] Uploading File 3",
	}

	// 文字
	P1 = widgets.NewParagraph()
	P1.Text = `hello在这里。今生无悔入东方 来世愿生幻想乡
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
	P1.Title = "中文测试"
	P2 = widgets.NewParagraph()
	P2.Title = "中文测试"
	P2.Text = "Press q to ther 中文测试nd things. Let's see how this thing renders out.\n    Here is a new paragraph and stuffs and things. There should be a tab indent at the beginning of the paragraph. Let's see if that worked as well."
	// P2.Text = utils.ChineseStringFormat("中文测试")
	P2.BorderStyle.Fg = ui.ColorBlue

	Img = widgets.NewImage(utils.LoadImg(resource.VideoCoverImg))

	// tab栏
	Tab = widgets.NewTabPane("index", "text", "moretext", "list")
	Tab.Border = true
	TabView = []interface{}{Img, Gauges[1], List, SparklineGroup}
	Grid = ui.NewGrid()
}

func RefreshGrid() {

	termWidth, termHeight := ui.TerminalDimensions()

	videoHeight := float64(termWidth * 4.0 / 16.0)
	videoHeightRate := videoHeight / float64(termHeight)
	// utils.Log.Infoln("系统窗口尺寸==>", termWidth, termHeight)
	// utils.Log.Infoln("视频窗口尺寸|高比率==>", termWidth, videoHeight, videoHeightRate)
	Grid.SetRect(0, 0, termWidth, termHeight)
	// 后期切换界面
	CurTabView = TabView[Tab.ActiveTabIndex]
	Grid.Set(
		ui.NewRow(1.0/10,
			ui.NewCol(1.0/4, List),
			ui.NewCol(2.0/4, Tab),
			ui.NewCol(1.0/4, List),
		),
		ui.NewRow(videoHeightRate, Img),
		ui.NewRow(1.0/20, Gauges[0]),
		// ui.NewRow(1.0/20, Gauges[1]),
		ui.NewRow(8.0/10-videoHeightRate-1.0/20,
			ui.NewCol(1.0/2,
				ui.NewRow(.9/3, Gauges[1]),
				ui.NewRow(.9/3, P2),
				ui.NewRow(1.2/3, Gauges[2]),
			),
			ui.NewCol(1.0/2, P1),
		),
		ui.NewRow(1.0/10,
			ui.NewCol(1.0/3, List),
			ui.NewCol(1.0/3, SparklineGroup),
			ui.NewCol(1.0/3, Plot),
		),
	)
}
