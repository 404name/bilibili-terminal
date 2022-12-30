// Copyright 2017 Zack Guo <zack.y.guo@gmail.com>. All rights reserved.
// Use of this source code is governed by a MIT license that can
// be found in the LICENSE file.

//go:build ignore
// +build ignore

package main

import (
	"log"
	"os"
	"strings"
	"time"

	"github.com/404name/termui-demo/global"
	"github.com/404name/termui-demo/model/video"
	"github.com/404name/termui-demo/resource"
	"github.com/404name/termui-demo/utils"
	ui "github.com/gizak/termui/v3"
)

func main() {
	if err := ui.Init(); err != nil {
		log.Fatalf("failed to initialize termui: %v", err)
	}
	defer ui.Close()

	Init()

	// 初始化模型

	defer utils.Log.Sync()

	// 开始渲染
	ui.Clear()
	ui.Render(global.Grid)

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
				if video.Player.Ready {
					var playControl interface{}
					video.Player.PlayChan <- playControl
				}

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

			// if tickerCount == 1000 {
			// 	return
			// }

			global.Gauges[2].Percent = (global.Gauges[2].Percent + 3) % 100

			// // global.SparklineGroup.Sparklines[0].Data = utils.SinFloat64[tickerCount : tickerCount+100]
			// // global.Plot.Data[0] = utils.SinFloat64[2*tickerCount:]
			ui.Render(global.Gauges[2])
			// // ui.Render(global.SparklineGroup)
			// // ui.Render(global.Plot)
			// tickerCount++
		}
	}
}

// 这边后面要抽出来，main文件里面只有main
func Init() {
	// 先初始化日志和系统服务
	initService()
	global.InitUI()
	global.RefreshGrid()
	video.Init()
}

func initService() {
	utils.InitLogger()
	// 初始化系统及设置命令行UTF-8格式
	initOS()
}

func initOS() {
	// 创建ffmpeg输出图片和音频的文件夹防止ffmpeg生成时候报错
	os.MkdirAll(resource.OutputAudioPath[:strings.LastIndex(resource.OutputAudioPath, "/")], os.ModePerm)
	os.MkdirAll(resource.OutputImgPath[:strings.LastIndex(resource.OutputImgPath, "/")], os.ModePerm)
	os.MkdirAll(resource.OutputVideoPath[:strings.LastIndex(resource.OutputVideoPath, "/")], os.ModePerm)

	// 添加UTF-8来支持中文
	utils.CallCommandRun("chcp", []string{"65001"})
}
