// Package asoul 相关功能
package asoul

import (
	"errors"
	"github.com/FloatTech/zbputils/control"
)

type littleEndian struct{}

type follows struct {
	TS   int `json:"ts"`
	Code int `json:"code"`
	Card struct {
		Mid        string  `json:"mid"`
		Name       string  `json:"name"`
		Sex        string  `json:"sex"`
		Face       string  `json:"face"`
		Regtime    int64   `json:"regtime"`
		Birthday   string  `json:"birthday"`
		Sign       string  `json:"sign"`
		Attentions []int64 `json:"attentions"`
		Fans       int     `json:"fans"`
		Friend     int     `json:"friend"`
		Attention  int     `json:"attention"`
		LevelInfo  struct {
			NextExp      int `json:"next_exp"`
			CurrentLevel int `json:"current_level"`
			CurrentMin   int `json:"current_min"`
			CurrentExp   int `json:"current_exp"`
		} `json:"level_info"`
		Pendant struct {
			Pid    int    `json:"pid"`
			Name   string `json:"name"`
			Image  string `json:"image"`
			Expire int    `json:"expire"`
		} `json:"pendant"`
	} `json:"card"`
}

type follower struct {
	Mid      int    `json:"mid"`
	Uname    string `json:"uname"`
	Video    int    `json:"video"`
	Roomid   int    `json:"roomid"`
	Rise     int    `json:"rise"`
	Follower int    `json:"follower"`
	GuardNum int    `json:"guardNum"`
	AreaRank int    `json:"areaRank"`
}

type vdInfo struct {
	Code int `json:"code"`
	Data struct {
		List struct {
			Vlist []struct {
				Pic     string `json:"pic"`
				Title   string `json:"title"`
				Created int    `json:"created"`
				Aid     int    `json:"aid"`
				Bvid    string `json:"bvid"`
			} `json:"vlist"`
		} `json:"list"`
		Page struct {
			Count int `json:"count"`
		} `json:"page"`
	} `json:"data"`
}

type medalInfo struct {
	Mid              int64  `json:"target_id"`
	MedalName        string `json:"medal_name"`
	Level            int64  `json:"level"`
	MedalColorStart  int64  `json:"medal_color_start"`
	MedalColorEnd    int64  `json:"medal_color_end"`
	MedalColorBorder int64  `json:"medal_color_border"`
}

type medal struct {
	Uname     string `json:"target_name"`
	medalInfo `json:"medal_info"`
}

type vup struct {
	Mid    int64  `gorm:"column:mid;primary_key"`
	Uname  string `gorm:"column:uname"`
	Roomid int64  `gorm:"column:roomid"`
}

const (
	diana    = 672328094
	ava      = 672346917
	kira     = 672353429
	queen    = 672342685
	datapath = "data/vtbs/"
	dbfile   = datapath + "vup.db"
)

var (
	LittleEndian  littleEndian
	asoul         = []int64{672328094, 672346917, 672353429, 672342685}
	cookie        = "添加自己的cookie"
 	errNeedCookie = errors.New("该api需要设置b站cookie，去plugin/asoul/config.go文件里加")
	engine        = control.Register("asoul", &control.Options{
		DisableOnDefault: false,
		Help: "=======asoul相关功能=======\n" +
			"- /查 [名字|uid] (查询bilibili用户关注vtb的情况)\n" +
			"- 日程表 (从asoul官号获取最新的日程表)\n" +
			"- 来点然/晚/牛/乃/狼能量 (随机推送一条对应账号的投稿)\n" +
			"- 粉丝信息 (发送bilibili平台asoul官号+5个小姐姐的粉丝数据)",
	})
)
