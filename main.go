// Copyright 2017 Zack Guo <zack.y.guo@gmail.com>. All rights reserved.
// Use of this source code is governed by a MIT license that can
// be found in the LICENSE file.

//go:build ignore
// +build ignore

package main

import (
	"log"
	"time"

	"github.com/404name/termui-demo/global"
	"github.com/404name/termui-demo/model"
	"github.com/404name/termui-demo/utils"
	"github.com/404name/termui-demo/widget"
	ui "github.com/gizak/termui/v3"
)

func main() {
	if err := ui.Init(); err != nil {
		log.Fatalf("failed to initialize termui: %v", err)
	}
	defer ui.Close()
	//初始化系统服务和UI布局
	global.Init()
	// 初始化模型
	model.Init()
	widget.Init()
	defer global.Log.Sync()

	// 开始渲染
	ui.Clear()
	ui.Render(global.Grid)
	tickerCount := 1
	uiEvents := ui.PollEvents()
	ticker := time.NewTicker(time.Second).C
	for {
		select {
		case e := <-uiEvents:
			//println(e.ID)
			switch e.ID {

			case "q", "<C-c>":
				return
			case "<Space>", "Enter":
				var playControl interface{}
				model.Video.PlayChan <- playControl
			case "<Left>":
				global.Tab.FocusLeft()
				global.RefreshGrid()
				ui.Render(global.Grid)
			case "<Right>":
				global.Tab.FocusRight()
				global.RefreshGrid()
				// global.CurTabView = global.TabView[global.Tab.ActiveTabIndex]
				// println(global.CurTabView)
				ui.Render(global.Grid)
			case "<Resize>":
				global.RefreshGrid()
				ui.Render(global.Grid)
			}
		case <-ticker:
			if tickerCount == 1000 {
				return
			}
			// global.Log.Debugln(model.Video)
			for _, g := range global.Gauges {
				g.Percent = (g.Percent + 3) % 100
			}
			global.Gauges[0].Percent = model.Video.CurrentPos * 100 / model.Video.Duration
			global.Gauges[0].Title = model.Video.GetProgressTitle()
			global.SparklineGroup.Sparklines[0].Data = utils.SinFloat64[tickerCount : tickerCount+100]
			global.Plot.Data[0] = utils.SinFloat64[2*tickerCount:]
			ui.Render(global.Gauges[0])
			ui.Render(global.Gauges[1])
			ui.Render(global.Gauges[2])
			ui.Render(global.SparklineGroup)
			ui.Render(global.Plot)
			tickerCount++
		}
	}
}
