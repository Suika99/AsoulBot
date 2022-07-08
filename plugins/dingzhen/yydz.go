package dingzhen

import (
	"encoding/json"
	"github.com/FloatTech/zbputils/control"
	log "github.com/sirupsen/logrus"
	zero "github.com/wdvxdr1123/ZeroBot"
	"github.com/wdvxdr1123/ZeroBot/message"
	"net/http"
)

type yiyandingzhen []struct {
	Picpath Picpath `json:"picpath"`
}

type Picpath struct {
	Num0    string `json:"0"`
	PicPath string `json:"pic_path"`
}

func init() {
	engine := control.Register("yiyandingzhen", &control.Options{
		DisableOnDefault: false,
		Help: "- 一眼丁真\n" +
			"- 随机发送一张丁真图",
	})
	engine.OnKeywordGroup([]string{"yydz", "dz", "一眼丁真", "丁真"}).
		Handle(func(ctx *zero.Ctx) {
			ctx.SendChain(message.Image(getImage()))
		})
}

func getImage() string {
	url := "https://www.yiyandingzhen.top/"
	api := "getpic.php"
	resp, err := http.Get(url + api)
	if err != nil {
		log.Errorln(err)
	}
	defer resp.Body.Close()
	result := yiyandingzhen{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		log.Errorln(err)
	}
	image := url + result[0].Picpath.PicPath
	return image
}
