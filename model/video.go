package model

import (
	"fmt"
	"image"

	"github.com/404name/termui-demo/ffmpeg"
	"github.com/404name/termui-demo/global"
	"github.com/404name/termui-demo/resource"
	"github.com/404name/termui-demo/utils"
	"github.com/hajimehoshi/oto/v2"
)

const (
	VideoPreLoadDuration = 3  // 提前几秒开始缓存
	VideoPreLoadGap      = 60 // 加载间隔/预加载秒数
	VideoFrameRate       = 12 // 视频帧率
)

type VideoDetail struct {
	URL        string           // 视频URL
	Duration   int              // 视频时长
	CurrentPos int              // 当前看到多少秒
	PreLoadPos int              // 预加载到了多少秒
	FrameLeft  int              // 当前一秒还剩下多少帧需要播放
	FrameCache chan image.Image // 截取的视频帧缓存器
	AudioCache chan oto.Player  // 截取的音频缓存器
	Audio      oto.Player
	PlayChan   chan interface{} // 让视频播放和暂停
}

func (v *VideoDetail) Init() error {
	if v.URL == "" {
		v.URL = resource.BaseVideoUrl
	}
	v.CurrentPos = 0
	v.FrameCache = make(chan image.Image, VideoPreLoadGap*VideoFrameRate)
	v.AudioCache = make(chan oto.Player, 1)
	v.PlayChan = make(chan interface{}, 1)
	v.Duration = ffmpeg.GetVideoDuration(v.URL)
	v.FrameLeft = VideoFrameRate
	go v.Load(false)
	return nil
}
func (v *VideoDetail) Clear() error {
	v.PreLoadPos = 0
	v.FrameLeft = VideoFrameRate
	v.Load(false)
	return nil
}

func (v *VideoDetail) GetProgressTitle() string {
	return utils.VideoDurationFormat(v.CurrentPos, false) + " / " + utils.VideoDurationFormat(v.Duration, false)
}

// 加载视频[音频&图片] preload为是否预加载
func (v *VideoDetail) Load(preload bool) {

	// 第一次不预加载
	if preload {
		v.PreLoadPos += VideoPreLoadGap
	}
	if v.PreLoadPos > v.Duration {
		v.PreLoadPos = v.Duration
	}

	toPos := v.PreLoadPos + VideoPreLoadGap
	if toPos > v.Duration {
		toPos = v.Duration
	}

	// 比如5秒加载一次,当前是第3秒,并且规定提前两秒去加载
	if err := ffmpeg.GetIpcScreenShot("ffmpeg", v.URL, resource.OutputImgPath, resource.OutputAudioPath, v.PreLoadPos, VideoFrameRate, toPos); err != nil {
		global.Log.Errorln("请求异常====>", err)
		return
	}

	// decode图片
	for i := 1; i <= VideoFrameRate*(toPos-v.PreLoadPos); i++ {
		v.FrameCache <- utils.LoadImg(fmt.Sprintf(resource.OutputImgPath, i))
		global.Log.Infoln("请求中:缓存池还剩下====>", len(v.FrameCache))
	}
	v.AudioCache <- utils.LoadAudio(resource.OutputAudioPath)
	global.Log.Infoln("获取%d-%ds内共%d张图片及音频", v.PreLoadPos, toPos, VideoFrameRate*(toPos-v.PreLoadPos))
}
