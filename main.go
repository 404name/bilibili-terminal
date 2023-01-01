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

	"github.com/404name/termui-demo/core"
	"github.com/404name/termui-demo/global"
	"github.com/404name/termui-demo/utils"
	"github.com/404name/termui-demo/view"
	ui "github.com/gizak/termui/v3"
)

func main() {
	if err := ui.Init(); err != nil {
		log.Fatalf("failed to initialize termui: %v", err)
	}
	defer ui.Close()

	Init()

	// 初始化模型

	defer global.LOG.Sync()

	// 开始渲染

	uiEvents := ui.PollEvents()
	ticker := time.NewTicker(time.Second).C
	for {
		select {
		case e := <-uiEvents:
			//println(e.ID)
			switch e.ID {

			case "q", "<C-c>":
				return
			case "<Resize>":
				view.NowPage.Refresh()
			default:
				view.NowPage.EventHander(e)
			}
		case <-ticker:
			view.NowPage.(*view.VideoPage).Gauge.Percent = (view.NowPage.(*view.VideoPage).Gauge.Percent + 3) % 100
			ui.Render(view.NowPage.(*view.VideoPage).Gauge)
			// tickerCount++
		}
	}
}

// 这边后面要抽出来，main文件里面只有main
func Init() {
	// 先初始化日志和系统服务
	initService()
	view.InitUI()
	// video.Init()
}

func initService() {
	// 读取本地配置优先
	global.VIPER = core.Viper()
	global.LOG = core.Zap()
	global.LOG.Infoln("读取本地配置====》", global.CONFIG)
	global.PATH = utils.GetCurrentDirectory()
	global.LOG.Infoln("运行路径====》", global.PATH)
	// 初始化系统及设置命令行UTF-8格式
	initOS()
}

func initOS() {
	// 创建ffmpeg输出图片和音频的文件夹防止ffmpeg生成时候报错
	os.MkdirAll(global.CONFIG.Output.OutputAudioPath[:strings.LastIndex(global.CONFIG.Output.OutputAudioPath, "/")], os.ModePerm)
	os.MkdirAll(global.CONFIG.Output.OutputImgPath[:strings.LastIndex(global.CONFIG.Output.OutputImgPath, "/")], os.ModePerm)
	os.MkdirAll(global.CONFIG.Output.OutputVideoPath[:strings.LastIndex(global.CONFIG.Output.OutputVideoPath, "/")], os.ModePerm)

	// 添加UTF-8来支持中文
	utils.CallCommandRun("chcp", []string{"65001"})
}
