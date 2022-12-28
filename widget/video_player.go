package widget

import (
	"time"

	"github.com/404name/termui-demo/global"
	"github.com/404name/termui-demo/model"
	"github.com/404name/termui-demo/resource"
	"github.com/404name/termui-demo/utils"
	ui "github.com/gizak/termui/v3"
)

// 视频播放器
func VideoRender(v *model.VideoDetail) {
	for {
		select {
		case <-v.PlayChan:
			// 如果是已经播放完再次播放就会从头开始
			if v.CurrentPos == v.Duration {
				v.CurrentPos = 0
			}
			// 每一秒获取n张图片
			ticker := time.NewTicker(time.Second / model.VideoFrameRate)
		Loop:
			for {
				select {
				case <-ticker.C:
					global.Log.Infoln("当前秒====>", v.CurrentPos)
					if v.CurrentPos >= v.Duration && v.FrameLeft == 0 {

						// 还有余下的几帧 播放完
						if len(v.FrameChan) != 0 {
							v.FrameLeft = len(v.FrameChan)
							continue
						}

						// 添加播放完图片
						global.Img.Image = utils.LoadImg(resource.VideoCoverImg)
						ui.Render(global.Img)
						ticker.Stop()
						// 初始化播放资源
						v.Clear()
						break Loop
					} else {
						// 正常加载视频
						if v.FrameLeft == 0 {
							v.FrameLeft = model.VideoFrameRate
							// 当前秒结束
							v.CurrentPos++
							// 预先一秒去加载
							if (v.CurrentPos+model.VideoPreLoadDuration)%model.VideoPreLoadGap == 0 {
								global.Log.Infoln("预加载====>", v.PreLoadPos+model.VideoPreLoadGap)
								go v.GetImgWithPreload(true)
							}
						}
						v.FrameLeft--

						global.Log.Infoln("当前秒剩余帧数====>", v.FrameLeft)

						global.Img.Image = <-v.FrameChan

						ui.Render(global.Img)
						global.Log.Infoln("没卡住刷新界面====>", v.FrameLeft)
					}
				case <-v.PlayChan:
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
