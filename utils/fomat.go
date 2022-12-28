package utils

import "fmt"

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
