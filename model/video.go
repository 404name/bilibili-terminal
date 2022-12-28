package model

import (
	"fmt"
	"image"

	"github.com/404name/termui-demo/ffmpeg"
	"github.com/404name/termui-demo/global"
	"github.com/404name/termui-demo/resource"
	"github.com/404name/termui-demo/utils"
)

const (
	VideoPreLoadDuration = 5  // 提前几秒开始缓存
	VideoPreLoadGap      = 60 // 加载间隔/预加载秒数
	VideoFrameRate       = 12 // 视频帧率
)

type VideoDetail struct {
	URL        string           // 视频URL
	Duration   int              // 视频时长
	CurrentPos int              // 当前看到多少秒
	PreLoadPos int              // 预加载到了多少秒
	FrameLeft  int              // 当前一秒还剩下多少帧需要播放
	FrameChan  chan image.Image // 截取的视频帧
	PlayChan   chan interface{} // 让视频播放和暂停
}

func (v *VideoDetail) Init() error {
	if v.URL == "" {
		v.URL = resource.BaseVideoUrl
	}
	v.CurrentPos = 0
	v.FrameChan = make(chan image.Image, VideoPreLoadGap*VideoFrameRate)
	v.PlayChan = make(chan interface{}, 1)
	v.Duration = ffmpeg.GetVideoDuration(v.URL)
	v.FrameLeft = VideoFrameRate
	v.GetImgWithPreload(false)
	return nil
}
func (v *VideoDetail) Clear() error {
	v.PreLoadPos = 0
	v.FrameLeft = VideoFrameRate
	v.GetImgWithPreload(false)
	return nil
}

func (v *VideoDetail) GetProgressTitle() string {
	return utils.VideoDurationFormat(v.CurrentPos, false) + " / " + utils.VideoDurationFormat(v.Duration, false)
}

func (v *VideoDetail) GetImgWithPreload(preload bool) {

	// 第一次不预加载

	if preload {
		v.PreLoadPos += VideoPreLoadGap
	}
	if v.PreLoadPos > v.Duration {
		v.PreLoadPos = v.Duration
	}

	gap := VideoPreLoadGap
	if v.PreLoadPos+gap > v.Duration {
		gap = v.Duration - v.PreLoadPos
	}

	// 比如5秒加载一次,当前是第3秒,并且规定提前两秒去加载
	if err := ffmpeg.GetIpcScreenShot("ffmpeg", v.URL, resource.BaseImg, v.PreLoadPos, VideoFrameRate, gap); err != nil {
		global.Log.Errorln("请求异常====>", err)
		return
	}

	// decode图片
	for i := 1; i <= VideoFrameRate*gap; i++ {
		v.FrameChan <- utils.LoadImg(fmt.Sprintf(resource.BaseImg, i))
		global.Log.Infoln("请求中:缓存池还剩下====>", len(v.FrameChan))
	}
	global.Log.Infoln("获取%d-%ds内共%d张图片", v.PreLoadPos, v.PreLoadPos+gap, VideoFrameRate*gap)
}
