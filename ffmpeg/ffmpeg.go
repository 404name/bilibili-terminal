package ffmpeg

import (
	"fmt"
	"os/exec"
	"strconv"
	"strings"

	"github.com/404name/termui-demo/global"
)

// func main() {
// 	fmt.Println(GetIpcScreenShot("ffmpeg", "https://404name.oss-cn-shanghai.aliyuncs.com/lightcloud/2022/3/20/%E5%8F%A3%E7%BD%A9%E9%AA%8C%E6%94%B6.mp4", "C:/Users/Admin/go/src/github.com/404name/termui-demo/ffmpeg/images/test.png"))
// }

// 根据URL调用ffmpeg 获取截图
func GetIpcScreenShot(ffmpegPath string, url string, screenShotPath string, time int, frameRate int, perLoad int) error {
	// ffmpeg -timeout 10000000 -y -r 10 -i https://404name.oss-cn-shanghai.aliyuncs.com/lightcloud/2022/3/20/3724723_da3-1-16.mp4 -ss 00:00:01 -vframes 1 -f image2 -vcodec png ./ffmpeg/images/img-%02d.png
	var params []string

	if url[:4] == "http" {
		params = append(params, "-timeout")
		// 10s
		params = append(params, "10000000")
	}
	// params = append(params, "-ss")
	// params = append(params, utils.VideoDurationFormat(time, true))
	params = append(params, "-i")
	params = append(params, url)
	// params = append(params, "-r")
	params = append(params, "-ss")
	params = append(params, fmt.Sprint(time))
	params = append(params, "-t")
	params = append(params, fmt.Sprint(perLoad))
	params = append(params, "-r")
	params = append(params, fmt.Sprint(frameRate))
	// params = append(params, fmt.Sprint("select=between(t\\,%d\\,%d),fps=%d", time, time+1, frameRate))
	params = append(params, "-f")
	params = append(params, "image2")
	params = append(params, "-vcodec")
	params = append(params, "png")
	params = append(params, screenShotPath)

	_, err := CallCommandRun(ffmpegPath, params)
	if err != nil {
		global.Log.Errorln("获取截图出错，url为--->", url, err)
		return err
	}
	return nil
}

func GetVideoDuration(url string) int {
	var params []string
	params = append(params, "-v")
	params = append(params, "error")
	params = append(params, "-show_entries")
	params = append(params, "format=duration")
	params = append(params, "-of")
	params = append(params, "default=noprint_wrappers=1:nokey=1")
	params = append(params, url)
	var cmdName = "ffprobe"
	// out, err := CallCommandRun(cmdName, params)
	// cmd := exec.Command("ffmpeg", "-v error -show_entries format=duration -of default=noprint_wrappers=1:nokey=1 ./ffmpeg/video.mp4")
	// logger.Zap.Debugln("CallCommand Run 执行命令=> ", cmd)
	out, err := CallCommandRun(cmdName, params)
	if err != nil {
		return 0
	}
	// 这里必须通过 . 去把浮点数转整数，不知道为什么直接转float再int会失效
	duration, _ := strconv.Atoi(out[:strings.Index(out, ".")])
	return int(duration)
}

func resolveTime(seconds int) string {
	var day = seconds / (24 * 3600)
	hour := (seconds - day*3600*24) / 3600
	minute := (seconds - day*24*3600 - hour*3600) / 60
	second := seconds - day*24*3600 - hour*3600 - minute*60
	return fmt.Sprintf("%02d:%02d:%02d", hour, minute, second)
}

func CallCommandRun(cmdName string, args []string) (string, error) {

	cmd := exec.Command(cmdName, args...)

	bytes, err := cmd.Output()
	if err != nil {
		global.Log.Errorln("CallCommand Run 出错了.....", string(bytes), err.Error())
		return "", err
	}
	resp := string(bytes)
	global.Log.Infoln("CallCommand Run 参数=> ", args)
	global.Log.Infoln("CallCommand Run 执行命令=> ", cmd)
	global.Log.Infoln("CallCommand Run 调用完成=> ", resp)
	return resp, nil
}
