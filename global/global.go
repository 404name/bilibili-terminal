package global

import (
	"image"

	ui "github.com/gizak/termui/v3"
	"github.com/gizak/termui/v3/widgets"
	"go.uber.org/zap"
)

// 系统服务
var (
	Log *zap.SugaredLogger
)

// UI组件
var (
	SparklineGroup *widgets.SparklineGroup
	Plot           *widgets.Plot
	Gauges         []*widgets.Gauge
	List           *widgets.List
	P1             *widgets.Paragraph
	P2             *widgets.Paragraph
	Tab            *widgets.TabPane
	TabView        []interface{}
	Grid           *ui.Grid
	CurTabView     interface{}
	Img            *widgets.Image
	ImgList        []image.Image
)

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
		ui.NewRow(8.0/10-videoHeightRate-1.0/20,
			ui.NewCol(1.0/2,
				ui.NewRow(.9/3, List),
				ui.NewRow(.9/3, Gauges[1]),
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
