package ffxiv

import (
	"github.com/FloatTech/zbputils/ctxext"
	zero "github.com/wdvxdr1123/ZeroBot"
)

func init() {
	engine.OnRegex(`^m\s(.{1,25})$`).SetBlock(true).Limit(ctxext.LimitByUser).
		Handle(func(ctx *zero.Ctx) {
			ctx.SendChain(cloud163(ctx.State["regex_matched"].([]string)[1] + " 植松伸夫"))
		})
}
