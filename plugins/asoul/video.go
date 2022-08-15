package asoul

import (
	zero "github.com/wdvxdr1123/ZeroBot"
	"github.com/wdvxdr1123/ZeroBot/message"
	"math/rand"
	"time"
)

// 烂写，refactor
func init() {
	engine.OnKeyword("来点然能量").
		Handle(func(ctx *zero.Ctx) {
			data := getVideo(diana)
			rand.Seed(time.Now().UnixNano())
			ranNub := rand.Intn(50)
			ctx.SendChain(message.Image(data.Data.List.Vlist[ranNub].Pic))
			ctx.SendChain(message.CustomMusic(
				"https://bilibili.com/video/"+data.Data.List.Vlist[ranNub].Bvid,
				"11111112355",
				data.Data.List.Vlist[ranNub].Title,
			))
		})
}

func init() {
	engine.OnKeyword("来点晚能量").
		Handle(func(ctx *zero.Ctx) {
			data := getVideo(ava)
			rand.Seed(time.Now().UnixNano())
			ranNub := rand.Intn(50)
			ctx.SendChain(message.Image(data.Data.List.Vlist[ranNub].Pic))
			ctx.SendChain(message.CustomMusic(
				"https://bilibili.com/video/"+data.Data.List.Vlist[ranNub].Bvid,
				"11111112355",
				data.Data.List.Vlist[ranNub].Title,
			))
		})
}

func init() {
	engine.OnKeyword("来点牛能量").
		Handle(func(ctx *zero.Ctx) {
			data := getVideo(bella)
			rand.Seed(time.Now().UnixNano())
			ranNub := rand.Intn(50)
			ctx.SendChain(message.Image(data.Data.List.Vlist[ranNub].Pic))
			ctx.SendChain(message.CustomMusic(
				"https://bilibili.com/video/"+data.Data.List.Vlist[ranNub].Bvid,
				"11111112355",
				data.Data.List.Vlist[ranNub].Title,
			))
		})
}

func init() {
	engine.OnKeyword("来点乃能量").
		Handle(func(ctx *zero.Ctx) {
			data := getVideo(eileen)
			rand.Seed(time.Now().UnixNano())
			ranNub := rand.Intn(50)
			ctx.SendChain(message.Image(data.Data.List.Vlist[ranNub].Pic))
			ctx.SendChain(message.CustomMusic(
				"https://bilibili.com/video/"+data.Data.List.Vlist[ranNub].Bvid,
				"11111112355",
				data.Data.List.Vlist[ranNub].Title,
			))
		})
}
