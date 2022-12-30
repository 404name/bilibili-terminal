package video

import (
	"fmt"
	"image"

	"github.com/404name/termui-demo/ffmpeg"
	"github.com/404name/termui-demo/resource"
	"github.com/404name/termui-demo/utils"
	"github.com/hajimehoshi/oto/v2"
)

const (
	VideoPreLoadDuration = 5  // 提前几秒开始缓存
	VideoPreLoadGap      = 30 // 加载间隔/预加载秒数
	VideoFrameRate       = 12 // 视频帧率
)

var (
	Player *VideoDetail
)

func Init() {
	Player = &VideoDetail{}
	Player.Init()
}

type VideoDetail struct {
	bilibiliCid // bilicard

	Ready      bool             // 是否加载完成
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
		v.URL = resource.OutputVideoPath
	}

	v.Ready = false
	v.CurrentPos = 215
	v.PreLoadPos = v.CurrentPos
	v.FrameCache = make(chan image.Image, VideoPreLoadGap*VideoFrameRate)
	v.AudioCache = make(chan oto.Player, 1)
	v.PlayChan = make(chan interface{}, 1)

	v.FrameLeft = VideoFrameRate

	// 加载视频
	videos := getCidList(bvid, qn)
	v.bilibiliCid = videos[0]
	// 暂时只读取第一个视频
	v.bilibiliCid.PlayURLs = v.bilibiliCid.PlayURLs[:1]

	// 下面二选一，第一个是下载 第二个直接载入
	//go v.bilibiliCid.download(resource.OutputVideoPath)
	go v.Load(false)

	// 开启渲染
	go VideoRender(v)
	return nil
}
func (v *VideoDetail) Clear() error {
	v.CurrentPos = 0
	v.PreLoadPos = 0
	v.FrameLeft = VideoFrameRate
	select {
	case <-v.FrameCache:
	default:
	}
	v.Audio.Close()
	v.Audio = nil
	go v.Load(false)
	return nil
}
func (v *VideoDetail) RefreshAudio() {
	// 音频处理
	audio := <-v.AudioCache
	//startT := time.Now() //计算当前时间
	audio.Play()
	//utils.Log.Errorln("音频播放耗时====>", time.Since(startT))
	if err := v.Audio.Close(); err != nil {
		utils.Log.Errorln("音频释放失败====>", err)
	}
	// 这里同步
	v.Audio = audio
}
func (v *VideoDetail) GetProgressTitle() string {
	return utils.VideoDurationFormat(v.CurrentPos, false) + " / " + utils.VideoDurationFormat(v.Duration, false)
}

// 加载视频[音频&图片] preload为是否预加载
func (v *VideoDetail) Load(preload bool) {

	// 第一次不预加载
	if !preload {
		v.Ready = true
		// 加载时长
		v.Duration = ffmpeg.GetVideoDuration(v.URL)
		VideoTitleRender(v.Title)
		// 渲染进度条
		VideoProgressBarRender(v)
		// 渲染视频等待激活

	} else {
		v.PreLoadPos += VideoPreLoadGap
	}
	if v.PreLoadPos > v.Duration {
		v.PreLoadPos = v.Duration
	}

	toPos := v.PreLoadPos + VideoPreLoadGap
	if toPos > v.Duration {
		// 超出1s加载，比如时长219,219到219其实是0s，应该往后加载一秒
		toPos = v.Duration + 1
	}

	// 比如5秒加载一次,当前是第3秒,并且规定提前两秒去加载
	if err := ffmpeg.GetIpcScreenShot("ffmpeg", v.URL, resource.OutputImgPath, resource.OutputAudioPath, v.PreLoadPos, VideoFrameRate, toPos); err != nil {
		utils.Log.Errorln("请求异常====>", err)
		return
	}

	// decode图片
	for i := 1; i <= VideoFrameRate*(toPos-v.PreLoadPos); i++ {
		v.FrameCache <- utils.LoadImg(fmt.Sprintf(resource.OutputImgPath, i))

	}

	v.AudioCache <- utils.LoadAudio(resource.OutputAudioPath)
	utils.Log.Infoln("请求中:缓存池还剩下====>", len(v.FrameCache)/VideoPreLoadGap*VideoFrameRate)
	utils.Log.Infof("获取%d-%ds内共%d张图片及音频", v.PreLoadPos, toPos, VideoFrameRate*(toPos-v.PreLoadPos))
}
