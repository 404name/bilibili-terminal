package video

import (
	"time"

	"github.com/404name/termui-demo/global"
	"github.com/404name/termui-demo/resource"
	"github.com/404name/termui-demo/utils"
	ui "github.com/gizak/termui/v3"
)

func VideoDownloadBarRender(current, total int) {
	if Player.Ready {
		global.Log.Errorln("Read未能正常关闭")
		return
	}
	global.Gauges[0].Percent = current * 100 / total
	if current == total {
		global.Log.Errorln("加载完毕")
		global.Gauges[0].Title = "loading~"
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
		global.Gauges[0].Title = ""
		ui.Render(global.Gauges[0])
	}

}

func VideoProgressBarRender(v *VideoDetail) {
	global.Gauges[0].Percent = v.CurrentPos * 100 / Player.Duration
	global.Gauges[0].Title = v.GetProgressTitle()
	ui.Render(global.Gauges[0])
}

// 视频播放器
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
					if v.CurrentPos >= v.Duration && v.FrameLeft == 0 {

						// 还有余下的几帧 播放完
						if len(v.FrameCache) != 0 {
							v.FrameLeft = len(v.FrameCache)
							continue
						}
						v.Audio.Close()
						v.Clear()
						// 添加播放完图片
						global.Img.Image = utils.LoadImg(resource.VideoCoverImg)
						ui.Render(global.Img)
						ticker.Stop()
						// 初始化播放资源

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
								// global.Log.Infoln("预加载====>", v.PreLoadPos+model.VideoPreLoadGap)
								go v.Load(true)
							}
						}
						v.FrameLeft--

						// global.Log.Infoln("当前秒剩余帧数====>", v.FrameLeft)

						global.Img.Image = <-v.FrameCache
						ui.Render(global.Img)

						// global.Log.Infoln("没卡住刷新界面====>", v.FrameLeft)
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
