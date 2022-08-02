package asoul

import (
	zero "github.com/wdvxdr1123/ZeroBot"
	"github.com/wdvxdr1123/ZeroBot/message"
)

func init() {
	engine.OnKeyword("日程表").
		Handle(func(ctx *zero.Ctx) {
			picUrl := getDynamic()
			if len(picUrl) == 0 {
				ctx.Send("image not found")
				return
			} else if len(picUrl) == 1 {
				ctx.SendChain(message.Image(picUrl[0]))
			} else if len(picUrl) == 2 {
				ctx.SendChain(message.Image(picUrl[0]), message.Image(picUrl[1]))
			} else {
				ctx.Send("unknown error")
				return
			}
		})
}
