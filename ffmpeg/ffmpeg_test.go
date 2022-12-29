package ffmpeg

import (
	"testing"

	"github.com/404name/termui-demo/global"
	"github.com/404name/termui-demo/resource"
)

func init() {
	global.InitLogger()
}

func TestGetIpcScreenShot(t *testing.T) {
	GetIpcScreenShot("ffmpeg", "."+resource.BaseVideoUrl, "."+resource.OutputImgPath, "."+resource.OutputAudioPath, 6, 10, 1)
}
