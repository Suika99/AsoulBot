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
	Face     string `json:"face"`
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
	bella    = 672353429
	eileen   = 672342685
	datapath = "data/vtbs/"
	dbfile   = datapath + "vup.db"
)

var (
	LittleEndian  littleEndian
	Asoul         = []int64{672328094, 672346917, 672353429, 672342685}
	cookie        = "_uuid=9C533C58-92103-D238-C6F7-DB427E2492DD08771infoc; buvid3=D224246E-4D43-4A41-9371-789AAFCE8D59167638infoc; b_nut=1641285109; rpdid=|(kmJYRuumYR0J'uYRuJY)klY; buvid4=D0B36D56-FB05-4A1D-CA4A-6194E92BFD3518669-022020812-Hv0N1ZTLOoipgbP9gUsRyQ%3D%3D; CURRENT_FNVAL=4048; fingerprint=15865356dc2b14e5e1dab341b8af2ed2; buvid_fp_plain=undefined; buvid_fp=1ca8d1bb83d616039b6657f6ea74fd32; SESSDATA=359149b0%2C1660551200%2Ca6ade%2A21; bili_jct=9191c1bac641cc8cb640dc16757b42d3; DedeUserID=1306274; DedeUserID__ckMd5=1f2dd48dfa071812; sid=aubhtwi7; i-wanna-go-back=-1; b_ut=5; LIVE_BUVID=AUTO2016449992304388; PVID=1; CURRENT_BLACKGAP=0; blackside_state=0; bp_video_offset_1306274=676678028957319200; innersign=0; b_lsid=79357B12_181DCD03319; b_timer=%7B%22ffp%22%3A%7B%22333.1007.fp.risk_D224246E%22%3A%22181DCD03565%22%7D%7D"
	errNeedCookie = errors.New("该api需要设置b站cookie，请发送命令设置cookie，例如\"设置b站cookie SESSDATA=82da790d,1663822823,06ecf*31\"")
	engine        = control.Register("asoul", &control.Options{
		DisableOnDefault: false,
		Help: "=======asoul相关功能=======\n" +
			"- /查 [名字|uid] (查询bilibili用户关注vtb的情况)\n" +
			"- 日程表 (从asoul官号获取最新的日程表)\n" +
			"- 来点然/晚/牛/乃/狼能量 (随机推送一条对应账号的投稿)\n" +
			"- 粉丝信息 (发送bilibili平台asoul官号+5个小姐姐的粉丝数据)",
	})
)
