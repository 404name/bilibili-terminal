package video

import (
	"time"

	"github.com/404name/termui-demo/global"
	"github.com/404name/termui-demo/resource"
	"github.com/404name/termui-demo/utils"
	ui "github.com/gizak/termui/v3"
)

func VideoTitleRender(title string) {
	global.Img.Title = title
	ui.Render(global.Img)
}

func VideoDownloadBarRender(current, total int) {
	if Player.Ready {
		utils.Log.Errorln("Read未能正常关闭")
		return
	}
	global.Gauges[0].Percent = current * 100 / total
	if current == total {
		utils.Log.Errorln("加 载 完 毕 ")
		global.Gauges[0].Title = "加 载 完 毕 "
		ui.Render(global.Gauges[0])
		if !Player.Ready {
			Player.Load(false)
		}
	} else {
		if global.Gauges[0].Percent >= 75 {
			global.Gauges[0].BarColor = ui.ColorGreen
		} else if global.Gauges[0].Percent >= 50 {
			global.Gauges[0].BarColor = ui.ColorBlue
		} else if global.Gauges[0].Percent >= 25 {
			global.Gauges[0].BarColor = ui.ColorYellow
		}
		global.Gauges[0].Title = "加 载 中 "
		ui.Render(global.Gauges[0])
	}

}

func VideoProgressBarRender(v *VideoDetail) {
	global.Gauges[0].Percent = v.CurrentPos * 100 / Player.Duration
	global.Gauges[0].Title = v.GetProgressTitle()
	ui.Render(global.Gauges[0])
}

// 视频播放器

// 这里应该放到外面去
func VideoRender(v *VideoDetail) {
	for {
		select {
		case <-v.PlayChan:
			// 如果是已经播放完再次播放就会从头开始
			if v.CurrentPos == v.Duration {
				v.CurrentPos = 0
			}
			// 如果音频为空表示第一次播放，此时阻塞等待缓存到来
			if v.Audio == nil {
				v.Audio = <-v.AudioCache
			}
			if !v.Audio.IsPlaying() {
				v.Audio.Play()
			}
			// 每一秒获取n张图片
			ticker := time.NewTicker(time.Second / VideoFrameRate)
		Loop:
			for {
				select {
				case <-ticker.C:

					// 处理视频
					// 当时间达到视频时长，并且没有剩余帧，或者缓存为空
					if v.CurrentPos >= v.Duration && (v.FrameLeft == 0 || len(v.FrameCache) == 0) {

						// 还有余下的几帧 播放完
						if len(v.FrameCache) != 0 {
							v.FrameLeft = len(v.FrameCache)
							continue
						}
						ticker.Stop()
						// 添加播放完图片
						global.Img.Image = utils.LoadImg(resource.VideoCoverImg)
						ui.Render(global.Img)
						// 初始化播放资源
						v.Ready = false
						v.Clear()
						break Loop
					} else {
						// 正常加载视频
						if v.FrameLeft == 0 {
							v.FrameLeft = VideoFrameRate
							// 当前秒结束
							v.CurrentPos++
							VideoProgressBarRender(v)
							if v.CurrentPos == v.PreLoadPos {
								// 此时应该更换音频
								go v.RefreshAudio()
							}

							// 预先一秒去加载
							if (v.CurrentPos+VideoPreLoadDuration)%VideoPreLoadGap == 0 {
								// utils.Log.Infoln("预加载====>", v.PreLoadPos+model.VideoPreLoadGap)
								go v.Load(true)
							}
						}
						v.FrameLeft--

						// utils.Log.Infoln("当前秒剩余帧数====>", v.FrameLeft)

						global.Img.Image = <-v.FrameCache
						termWidth, termHeight := ui.TerminalDimensions()
						if global.Img.Dx() > termWidth {
							utils.Log.Infof("宽度异常%d %d:%d", global.Img.Dx(), termWidth, termHeight)
						} else {
							utils.Log.Infof("宽度%d %d:%d", global.Img.Dx(), termWidth, termHeight)
							ui.Render(global.Img)
						}

						// utils.Log.Infoln("没卡住刷新界面====>", v.FrameLeft)
					}
				case <-v.PlayChan:
					// 暂停音频
					v.Audio.Pause()
					// 添加暂停封面
					global.Img.Image = utils.LoadImg(resource.VideoStopImg)
					ui.Render(global.Img)
					ticker.Stop()
					break Loop
				}
			}

		}
	}
}
