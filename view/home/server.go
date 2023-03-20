package home

import (
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"

	tea "github.com/charmbracelet/bubbletea"
)

// var baseUrl string = `https://api.iwyu.com/API/baiduresou/`
var baseUrl string = fmt.Sprintf("https://api.bilibili.com/x/web-interface/index/top/feed/rcmd?feed_version=V4&fresh_idx_1h=1&fetch_row=1&fresh_idx=%d&brush=1&homepage_ver=1&ps=%d", rand.Intn(10000), 20)

type Result struct {
	Code int           `json:"code"`
	Mgs  string        `json:"mgs"`
	Ttl  int           `json:"ttl"`
	Data RcmdVideoList `json:"data"`
}

type RcmdVideoList struct {
	Item                  []*RcmdVideo `json:"item"`
	BusinessCard          interface{}  `json:"business_card"`
	FloorInfo             interface{}  `json:"floor_info"`
	UserFeature           interface{}  `json:"user_feature"`
	PreloadExposePct      float64      `json:"preload_expose_pct"`
	PreloadFloorExposePct float64      `json:"preload_floor_expose_pct"`
	Mid                   int          `json:"mid"`
}

type RcmdVideo struct {
	ID       int    `json:"id"`
	Bvid     string `json:"bvid"`     // 稿件bvid
	Cid      int    `json:"cid"`      // 稿件cid
	Goto     string `json:"goto"`     // 跳转类别
	URI      string `json:"uri"`      // 稿件cid
	Pic      string `json:"pic"`      // 封面
	Headline string `json:"title"`    // 标题
	Duration int    `json:"duration"` // 时长
	Pubdate  int    `json:"pubdate"`  // 发布日期
	Owner    struct {
		Mid  int    `json:"mid"`  // 发布人id
		Name string `json:"name"` // 发布人名字
		Face string `json:"face"` // 发布人头像
	} `json:"owner"`
	Stat struct {
		View    int `json:"view"`    // 观看数量
		Like    int `json:"like"`    // 喜欢数
		Danmaku int `json:"danmaku"` // 弹幕数量
	} `json:"stat"`
	AvFeature  interface{} `json:"av_feature"`
	IsFollowed int         `json:"is_followed"`
	RcmdReason struct {
		Content    string `json:"content"`     // 推荐理由
		ReasonType int    `json:"reason_type"` // 推荐种类
	} `json:"rcmd_reason"`
	ShowInfo     int         `json:"show_info"`
	TrackID      string      `json:"track_id"`
	Pos          int         `json:"pos"`
	RoomInfo     interface{} `json:"room_info"`
	OgvInfo      interface{} `json:"ogv_info"`
	BusinessInfo interface{} `json:"business_info"`
	IsStock      int         `json:"is_stock"`
}

// 下载远程图片到当前目录下cache文件夹
func (m *Tui) getRemoteImg(url string) tea.Cmd {

	return func() tea.Msg {
		// 获取图片文件名
		filename := filepath.Base(url)
		// 拼接本地缓存路径
		cachePath := filepath.Join("cache", filename)
		if _, err := os.Stat("cache"); os.IsNotExist(err) {
			err = os.Mkdir("cache", 0755)
			if err != nil {
				return err
			}
		}
		// 如果本地已经存在该文件，直接返回文件路径
		if _, err := os.Stat(cachePath); err == nil {
			return LoadImgMsg(cachePath)
		}

		// 否则，下载图片并保存到本地缓存
		resp, err := http.Get(url)
		if err != nil {
			return err
		}
		defer resp.Body.Close()

		out, err := os.Create(cachePath)
		if err != nil {
			return err
		}
		defer out.Close()

		_, err = io.Copy(out, resp.Body)
		if err != nil {
			return err
		}
		return LoadImgMsg(cachePath)
	}
}

func GetTrendingList() ([]*RcmdVideo, error) {
	// 用http从baseUrl获取搜索结果映射到result
	resp, err := http.Get(baseUrl)
	defer resp.Body.Close()
	// 得到请求的status code
	if err != nil {
		return nil, err
	}

	var res Result
	if err = json.NewDecoder(resp.Body).Decode(&res); err != nil {
		return nil, err
	}

	return res.Data.Item, nil
}

// 写个go写个终端打开链接的函数
func open(url string) tea.Cmd {
	return func() tea.Msg {
		var cmd string
		var args []string

		switch runtime.GOOS {
		case "windows":
			cmd = "cmd"
			args = []string{"/c", "start"}
		case "darwin":
			cmd = "open"
		case "linux":
			cmd = "xdg-open"
		default:
			return fmt.Errorf("unsupported platform")
		}

		args = append(args, url)
		return exec.Command(cmd, args...).Start()
	}

}
