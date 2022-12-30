package global

import (
	"image"

	ui "github.com/gizak/termui/v3"
	"github.com/gizak/termui/v3/widgets"
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
