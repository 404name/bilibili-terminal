package model

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"time"

	"github.com/404name/termui-demo/global"
	"github.com/404name/termui-demo/utils"
	"github.com/pkg/errors"
)

type RcmdVideoList struct {
	Item                  []RcmdVideo `json:"item"`
	BusinessCard          interface{} `json:"business_card"`
	FloorInfo             interface{} `json:"floor_info"`
	UserFeature           interface{} `json:"user_feature"`
	PreloadExposePct      float64     `json:"preload_expose_pct"`
	PreloadFloorExposePct float64     `json:"preload_floor_expose_pct"`
	Mid                   int         `json:"mid"`
}

type RcmdVideo struct {
	ID       int    `json:"id"`
	Bvid     string `json:"bvid"`     // 稿件bvid
	Cid      int    `json:"cid"`      // 稿件cid
	Goto     string `json:"goto"`     // 跳转类别
	URI      string `json:"uri"`      // 稿件cid
	Pic      string `json:"pic"`      // 封面
	Title    string `json:"title"`    // 标题
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

func GetRcmdVideo() ([]RcmdVideo, error) {
	rand.Seed(time.Now().Unix())
	resp, err := global.Request.SetTimeout(time.Second*10).R().SetHeader("Content-Type", "application/x-www-form-urlencoded").Get(fmt.Sprintf("https://api.bilibili.com/x/web-interface/wbi/index/top/feed/rcmd?feed_version=V4&fresh_idx_1h=1&fetch_row=1&fresh_idx=%d&brush=1&homepage_ver=1&ps=%d", rand.Intn(10000), 10))
	if err != nil {
		return nil, errors.WithStack(err)
	}
	data, err := utils.GetRespData(resp, "获取推荐视频列表")
	if err != nil {
		return nil, err
	}
	var content RcmdVideoList
	err = json.Unmarshal(data, &content)
	global.LOG.Infoln(content.Item)
	return content.Item, errors.WithStack(err)
}
