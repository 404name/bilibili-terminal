package main

import (
	"log"

	ui "github.com/gizak/termui/v3"
	"github.com/gizak/termui/v3/widgets"
)

func main() {
	// 1. 初始化项目
	if err := ui.Init(); err != nil {
		log.Fatalf("failed to initialize termui: %v", err)
	}
	defer ui.Close()

	// 2. 通过widgets创建画布
	p := widgets.NewParagraph()
	p.Text = "Hello World!"

	// 3 通过定位放置画布 或者 通过grid响应式布局嵌入画布
	p.SetRect(0, 0, 25, 5) // 或者grid := ui.NewGrid()
	// 4. 渲染 画布 / 布局
	ui.Render(p)

	// 5. 监听事件响应内容或者更新界面
	uiEvents := ui.PollEvents()
	for {
		select {
		// UI事件触发动作或者resize页面
		case e := <-uiEvents:
			println("e", e.ID, e.Type, e.Payload)
			switch e.ID {
			case "q", "<C-c>":
				return
			}
		}
	}
}
