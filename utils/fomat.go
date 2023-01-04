package utils

import (
	"fmt"
	"strings"

	"github.com/404name/termui-demo/global"
)

// 将时间格式化为 00:01 / 03:03
// hourFix为即使小时为0也会填充小时的前缀
func VideoDurationFormat(seconds int, hourFix bool) string {
	var day = seconds / (24 * 3600)
	hour := (seconds - day*3600*24) / 3600
	minute := (seconds - day*24*3600 - hour*3600) / 60
	second := seconds - day*24*3600 - hour*3600 - minute*60
	if !hourFix && hour == 0 {
		return fmt.Sprintf("%02d:%02d", minute, second)
	}
	return fmt.Sprintf("%02d:%02d:%02d", hour, minute, second)
}

// func ChineseStringFormat(str string) string {
// 	r := regexp.MustCompile(`([\p{Han}])`)
// 	return r.ReplaceAllString(str, "$1 ")
// }

func VersionMessage() string {
	return fmt.Sprintf(`
	欢迎使用 哔哩哔哩 for cmd
	当前版本:%s
		操作方法：键盘上下左右切换、空格点击、回车刷新
		项目地址: https://github.com/404name/bilibili-terminal
		项目视频: https://www.bilibili.com/video/BV1844y1d7Eg
`, global.VERSION)
}

func CommondNavigateTo(path string, param []string) string {
	if len(param) == 0 {
		return fmt.Sprintf("%s:%s: ", "NavigateTo", path)
	}
	return fmt.Sprintf("%s:%s:%s&%s", "NavigateTo", path, strings.Join(param, "&"))
}

func CommondNavigateBack() string {
	return fmt.Sprintf("%s: : ", "NavigateBack")
}
