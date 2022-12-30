package ffmpeg

import (
	"testing"

	"github.com/404name/termui-demo/resource"
	"github.com/404name/termui-demo/utils"
)

func init() {
	utils.InitLogger()
}

func TestGetIpcScreenShot(t *testing.T) {
	GetIpcScreenShot("ffmpeg", "."+resource.BaseVideoUrl, "."+resource.OutputImgPath, "."+resource.OutputAudioPath, 6, 10, 1)
}
