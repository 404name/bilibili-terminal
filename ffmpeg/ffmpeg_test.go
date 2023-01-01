package ffmpeg

import (
	"testing"

	"github.com/404name/termui-demo/core"
	"github.com/404name/termui-demo/global"
)

func init() {
	global.VIPER = core.Viper()
	global.LOG = core.Zap()
}

func TestGetIpcScreenShot(t *testing.T) {
	// GetIpcScreenShot("ffmpeg", "."+resource.BaseVideoUrl, "."+resource.OutputImgPath, "."+resource.OutputAudioPath, 6, 10, 1)
}
