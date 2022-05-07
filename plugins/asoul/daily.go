package asoul

import (
	log "github.com/sirupsen/logrus"
	"github.com/tidwall/gjson"
	zero "github.com/wdvxdr1123/ZeroBot"
	"github.com/wdvxdr1123/ZeroBot/message"
	"io/ioutil"
	"net/http"
	"strings"
)

func init() {
	engine.OnKeyword("日程表").
		Handle(func(ctx *zero.Ctx) {
			url := getDynamic()
			ctx.SendChain(
				message.Text(url),
				message.Image(url))
		})
}

func getDynamic() string {
	api := "https://api.vc.bilibili.com/dynamic_svr/v1/dynamic_svr/space_history?host_uid=703007996"
	resp, err := http.Get(api)
	if err != nil {
		log.Error(err)
	}
	data, _ := ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	json := gjson.ParseBytes(data)

	dy := json.Get("data.cards.#.card").Array()
	for _, i := range dy {
		if strings.Index(i.String(), "日程表") >= 0 {
			org := (gjson.Parse(i.String()).Get("origin"))
			pic := (gjson.Parse(org.String()).Get("item.pictures").Array())
			return pic[0].Get("img_src").String()
		} else {
			return "Not found image"
		}
	}
	return "unknown error"
}
