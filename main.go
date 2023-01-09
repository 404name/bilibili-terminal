// Copyright 2017 Zack Guo <zack.y.guo@gmail.com>. All rights reserved.
// Use of this source code is governed by a MIT license that can
// be found in the LICENSE file.

//go:build ignore
// +build ignore

package main

import (
	"context"
	"log"
	"os"
	"strings"
	"time"

	"github.com/404name/termui-demo/core"
	"github.com/404name/termui-demo/global"
	"github.com/404name/termui-demo/utils"
	"github.com/404name/termui-demo/view"
	ui "github.com/gizak/termui/v3"
	"github.com/go-resty/resty/v2"
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
	// ticker := time.NewTicker(time.Second).C
	for {
		select {
		case e := <-uiEvents:
			//println(e.ID)
			switch e.ID {

			case "<C-c>":
				return
			case "<Resize>":
				view.BiliUI.Refresh()
			case "q":
				{
					global.Command <- utils.CommondNavigateBack()
				}
			default:
				view.BiliUI.EventHander(e)
			}
		case cmd := <-global.Command:
			global.LOG.Infof("处理命令====>%s", cmd)
			// 预期用[type:action:param]组合区分指令
			s := strings.Split(cmd, ":")
			cmdType, cmdAction, param := s[0], s[1], s[2]
			switch cmdType {
			case "NavigateTo":
				view.BiliUI.NavigateTo(cmdAction, param)
			case "NavigateBack":
				view.BiliUI.NavigateBack()
			}
		}
	}
}

// 这边后面要抽出来，main文件里面只有main
func Init() {
	// 先初始化日志和系统服务
	initService()
	view.Init()
	// video.Init()
}

func initService() {
	// 读取本地配置优先
	global.VIPER = core.Viper()
	global.LOG = core.Zap()
	global.LOG.Infoln("读取本地配置====》", global.CONFIG)
	global.PATH = utils.GetCurrentDirectory()
	global.Request = resty.New().SetTimeout(time.Second * 10).SetLogger(global.LOG)
	global.Command = make(chan string, 10)
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
	utils.CallCommandRun(context.Background(), "chcp", []string{"65001"})
}
