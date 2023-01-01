package video

import (
	"fmt"
	"time"

	"github.com/404name/termui-demo/global"
	"github.com/404name/termui-demo/utils"
	ui "github.com/gizak/termui/v3"
)

func (v *VideoDetail) AddLog(str string) {
	// 刷新log
	v.Log.Text = "bilicli: " + str + "\n" + v.Log.Text
	ui.Render(v.Log)
	// 大于清空
	if len(v.Log.Text) > 1000 {
		v.Log.Text = ""
	}
}
func (v *VideoDetail) VideoTitleRender() {
	v.Img.Title = v.Title
	ui.Render(v.Img)
}

func (v *VideoDetail) VideoDownloadBarRender(current, total int) {
	if v.Ready {
		global.LOG.Errorln("Read未能正常关闭")
		return
	}
	v.ProgressBar.Percent = current * 100 / total
	if current == total {
		global.LOG.Errorln("加 载 完 毕 ")
		v.AddLog(fmt.Sprintf("视频%v加 载 完 毕 ", v.bilibiliCid))
		v.AddLog(utils.VersionMessage())
		v.ProgressBar.Title = "加 载 完 毕 "
		ui.Render(v.ProgressBar)
		if !v.Ready {
			v.Load(false)
		}
	} else {
		if v.ProgressBar.Percent >= 75 {
			v.ProgressBar.BarColor = ui.ColorGreen
		} else if v.ProgressBar.Percent >= 50 {
			v.ProgressBar.BarColor = ui.ColorBlue
		} else if v.ProgressBar.Percent >= 25 {
			v.ProgressBar.BarColor = ui.ColorYellow
		}
		v.AddLog(fmt.Sprintf("视频%s[%s]加载中[%d/%]====>", v.Title, v.Bvid, v.ProgressBar.Percent))
		v.ProgressBar.Title = "加 载 中 "
		ui.Render(v.ProgressBar)
	}

}

func VideoProgressBarRender(v *VideoDetail) {
	v.ProgressBar.Percent = v.CurrentPos * 100 / v.Duration
	v.ProgressBar.Title = v.GetProgressTitle()
	ui.Render(v.ProgressBar)
}

// 视频播放器

// 这里应该放到外面去
func (v *VideoDetail) VideoRender() {
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
						v.Img.Image = utils.LoadImg(global.CONFIG.BasePath.VideoCoverImg)
						ui.Render(v.Img)
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
								// global.LOG.Infoln("预加载====>", v.PreLoadPos+model.VideoPreLoadGap)
								go v.Load(true)
							}
						}
						v.FrameLeft--

						// global.LOG.Infoln("当前秒剩余帧数====>", v.FrameLeft)

						v.Img.Image = <-v.FrameCache
						termWidth, termHeight := ui.TerminalDimensions()
						if v.Img.Dx() > termWidth {
							global.LOG.Infof("宽度异常%d %d:%d", v.Img.Dx(), termWidth, termHeight)
						} else {
							global.LOG.Infof("宽度%d %d:%d", v.Img.Dx(), termWidth, termHeight)
							ui.Render(v.Img)
						}

						// global.LOG.Infoln("没卡住刷新界面====>", v.FrameLeft)
					}
				case <-v.PlayChan:
					// 暂停音频
					v.Audio.Pause()
					// 添加暂停封面
					v.Img.Image = utils.LoadImg(global.CONFIG.BasePath.VideoStopImg)
					ui.Render(v.Img)
					ticker.Stop()
					break Loop
				}
			}

		}
	}
}
