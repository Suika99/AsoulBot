// Package emo 网易云音乐热评
package emo

import (
	"github.com/FloatTech/zbputils/control"
	"github.com/FloatTech/zbputils/ctxext"
	"github.com/FloatTech/zbputils/web"
	zero "github.com/wdvxdr1123/ZeroBot"
	"github.com/wdvxdr1123/ZeroBot/message"
	"github.com/wdvxdr1123/ZeroBot/utils/helper"
)

const (
	wangyiyunURL     = "https://api.gmit.vip/Api/HotComments?format=text"
	wangyiyunReferer = "https://api.gmit.vip/"
	ua               = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/96.0.4664.110 Safari/537.36"
)

func init() {
	engine := control.Register("emo", &control.Options{
		DisableOnDefault: false,
		Help:             "wangyiyun \n- 来份网易云热评",
	})
	engine.OnKeywordGroup([]string{"emo", "网抑云"}).SetBlock(true).Limit(ctxext.LimitByUser).
		Handle(func(ctx *zero.Ctx) {
			data, err := web.RequestDataWith(web.NewDefaultClient(), wangyiyunURL, "GET", wangyiyunReferer, ua)
			if err != nil {
				ctx.SendChain(message.Text("ERROR:", err))
				return
			}
			ctx.SendChain(message.Text(helper.BytesToString(data)))
		})
}
