package global

import (
	"os"
	"strings"

	"github.com/404name/termui-demo/resource"
	"github.com/404name/termui-demo/utils"
	ui "github.com/gizak/termui/v3"
	"github.com/gizak/termui/v3/widgets"
)

func Init() {
	initService()
	initUI()
	RefreshGrid()
}

func initService() {
	// 创建ffmpeg输出图片和音频的文件夹防止ffmpeg生成时候报错
	os.MkdirAll(resource.OutputAudioPath[:strings.LastIndex(resource.OutputAudioPath, "/")], os.ModePerm)
	os.MkdirAll(resource.OutputImgPath[:strings.LastIndex(resource.OutputImgPath, "/")], os.ModePerm)
	os.MkdirAll(resource.OutputVideoPath[:strings.LastIndex(resource.OutputVideoPath, "/")], os.ModePerm)
	InitLogger()
}

func initUI() {
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
	P1.Text = "<> This row has 3 columns\n<- Widgets can \nbe stacked\n up like left side\n<- Stacked widgets are treated as a single widget"
	P1.Title = "Demonstration"

	P2 = widgets.NewParagraph()
	P2.Title = "Text Box with Wrapping"
	P2.Text = "Press q to QUIT THE DEMO. [There](fg:blue,mod:bold) are other things [that](fg:red) are going to fit in here I think. What do you think? Now is the time for all good [men to](bg:blue) come to the aid of their country. [This is going to be one really really really long line](fg:green) that is going to go together and stuffs and things. Let's see how this thing renders out.\n    Here is a new paragraph and stuffs and things. There should be a tab indent at the beginning of the paragraph. Let's see if that worked as well."
	P2.BorderStyle.Fg = ui.ColorBlue

	Img = widgets.NewImage(utils.LoadImg(resource.VideoCoverImg))

	// tab栏
	Tab = widgets.NewTabPane("index", "text", "moretext", "list")
	Tab.Border = true
	TabView = []interface{}{Img, Gauges[1], List, SparklineGroup}
}

func RefreshGrid() {
	Grid = ui.NewGrid()
	termWidth, termHeight := ui.TerminalDimensions()

	videoHeight := float64(termWidth * 4.0 / 16.0)
	videoHeightRate := videoHeight / float64(termHeight)
	Log.Infoln("系统窗口尺寸==>", termWidth, termWidth)
	Log.Infoln("视频窗口尺寸|高比率==>", termWidth, videoHeight, videoHeightRate)
	Grid.SetRect(0, 0, termWidth, termHeight)
	// 后期切换界面
	CurTabView = TabView[Tab.ActiveTabIndex]

	Grid.Set(
		ui.NewRow(1.0/10,
			ui.NewCol(1.0/4, List),
			ui.NewCol(2.0/4, Tab),
			ui.NewCol(1.0/4, List),
		),
		ui.NewRow(videoHeightRate, CurTabView),
		ui.NewRow(1.0/20, Gauges[0]),
		// ui.NewRow(1.0/20, Gauges[1]),
		ui.NewRow(8.0/10-videoHeightRate-1.0/20,
			ui.NewCol(1.0/2,
				ui.NewRow(.9/3, Gauges[1]),
				ui.NewRow(.9/3, List),
				ui.NewRow(1.2/3, Gauges[2]),
			),
			ui.NewCol(1.0/2,
				ui.NewRow(1.0/2, SparklineGroup),
				ui.NewRow(1.0/2, Plot),
			),
		),
		ui.NewRow(1.0/10,
			ui.NewCol(1.0/3, List),
			ui.NewCol(1.0/3, List),
			ui.NewCol(1.0/3, List),
		),
	)
}
