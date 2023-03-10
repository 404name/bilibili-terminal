package video

import (
	"context"
	"fmt"
	"image"

	"github.com/404name/termui-demo/ffmpeg"
	"github.com/404name/termui-demo/global"
	"github.com/404name/termui-demo/utils"
	"github.com/gizak/termui/v3/widgets"
	"github.com/hajimehoshi/oto/v2"
)

// 这边必须是常量
const (
	VideoPreLoadDuration = 5  // 提前几秒开始缓存
	VideoPreLoadGap      = 30 // 加载间隔/预加载秒数
	VideoFrameRate       = 12 // 视频帧率
)

type VideoDetail struct {
	bilibiliCid // bilicard

	ProgressBar *widgets.Gauge
	Img         *widgets.Image
	Log         *widgets.Paragraph

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
	CloseChan  chan interface{} // 退出视频

	ctx    context.Context    // ctx本身应该在参数中传递，但这里有些地方不用一直携带很侵入原有的代码
	cancel context.CancelFunc // 控制退出
}

func (v *VideoDetail) Init() error {
	v.ctx, v.cancel = context.WithCancel(context.Background())
	if v.URL == "" {
		// 这里也可以读取网络视频。把url改成网络的即可，前提是公开访问
		v.URL = global.CONFIG.Output.OutputVideoPath
	}
	v.Img.Image = utils.LoadImg(global.CONFIG.BasePath.VideoCoverImg)
	v.Ready = false
	v.CurrentPos = 0
	v.PreLoadPos = v.CurrentPos
	v.FrameCache = make(chan image.Image, VideoPreLoadGap*VideoFrameRate)
	v.AudioCache = make(chan oto.Player, 1)
	v.PlayChan = make(chan interface{}, 1)
	v.CloseChan = make(chan interface{}, 1)

	v.FrameLeft = VideoFrameRate

	// 加载视频
	videos := getCidList(v.Bvid, qn)
	global.LOG.Infoln(videos)
	v.bilibiliCid = videos[0]
	// 暂时只读取第一个视频
	v.bilibiliCid.PlayURLs = v.bilibiliCid.PlayURLs[:1]

	// 下面二选一，第一个是下载 第二个直接载入
	go v.download(global.CONFIG.Output.OutputVideoPath)
	//go v.Load(false)

	// 开启渲染
	go v.VideoRender(v.ctx)
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
	//global.LOG.Errorln("音频播放耗时====>", time.Since(startT))
	if err := v.Audio.Close(); err != nil {
		v.AddLog(fmt.Sprintf("音频释放失败====>", err))
		global.LOG.Errorln("音频释放失败====>", err)
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
		v.Duration = ffmpeg.GetVideoDuration(v.ctx, v.URL)
		select {
		case <-v.ctx.Done():
			global.LOG.Infoln("ctx中断任务")
			return
		default: // 默认退出
		}
		v.VideoTitleRender()
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
	if err := ffmpeg.GetIpcScreenShot(v.ctx, "ffmpeg", v.URL, global.CONFIG.Output.OutputImgPath, global.CONFIG.Output.OutputAudioPath, v.PreLoadPos, VideoFrameRate, toPos); err != nil {
		v.AddLog(fmt.Sprintf("请求异常====>%v", err))
		global.LOG.Errorln("请求异常====>%v", err)
		// 也可能是ctx Done了
		return
	}

	// decode图片
	for i := 1; i <= VideoFrameRate*(toPos-v.PreLoadPos); i++ {
		v.FrameCache <- utils.LoadImg(fmt.Sprintf(global.CONFIG.Output.OutputImgPath, i))

	}

	v.AudioCache <- utils.LoadAudio(global.CONFIG.Output.OutputAudioPath)
	v.AddLog(fmt.Sprintf("请求中:缓存池还剩下====>%d", len(v.FrameCache)/VideoPreLoadGap*VideoFrameRate))
	global.LOG.Infoln("请求中:缓存池还剩下====>", len(v.FrameCache)/VideoPreLoadGap*VideoFrameRate)

	v.AddLog(fmt.Sprintf("获取%d-%ds内共%d张图片及音频", v.PreLoadPos, toPos, VideoFrameRate*(toPos-v.PreLoadPos)))
	global.LOG.Infof("获取%d-%ds内共%d张图片及音频", v.PreLoadPos, toPos, VideoFrameRate*(toPos-v.PreLoadPos))
}
