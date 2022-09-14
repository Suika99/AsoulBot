package asoul

import (
	"github.com/FloatTech/zbputils/ctxext"
	zero "github.com/wdvxdr1123/ZeroBot"
)

func init() {
	engine.OnRegex(`^向晚\s(.{1,25})$`).SetBlock(true).Limit(ctxext.LimitByUser).
		Handle(func(ctx *zero.Ctx) {
			ctx.SendChain(cloud163(ctx.State["regex_matched"].([]string)[1] + " 向晚"))
		})

	engine.OnRegex(`^嘉然\s(.{1,25})$`).SetBlock(true).Limit(ctxext.LimitByUser).
		Handle(func(ctx *zero.Ctx) {
			ctx.SendChain(cloud163(ctx.State["regex_matched"].([]string)[1] + " 嘉然"))
		})
}
