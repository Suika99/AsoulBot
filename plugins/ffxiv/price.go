package ffxiv

import (
	"fmt"
	zero "github.com/wdvxdr1123/ZeroBot"
	"github.com/wdvxdr1123/ZeroBot/message"
	"time"
)

func init() {
	engine.OnRegex(`^p\s(.{1,25})$`).SetBlock(true).
		Handle(func(ctx *zero.Ctx) {
			keyword := ctx.State["regex_matched"].([]string)[1]
			cafeRes := getCafeSearch(keyword)
			itemID := cafeRes.Get("Results.#.ID").Array()[0].Int()
			itemCcon := fmt.Sprintf("https://cafemaker.wakingsands.com/%s", cafeRes.Get("Results.#.Icon").Array()[0].Str)
			priceData := getItemPrice(itemID).Get("listings").Array()
			var priceMessage string
			for _, v := range priceData {
				vPrice := fmt.Sprintf("%v     x%v=%v     %v:%v     %v\n",
					v.Get("pricePerUnit"),
					v.Get("quantity"),
					v.Get("total"),
					v.Get("worldName"),
					v.Get("retainerName"),
					time.Unix(v.Get("lastReviewTime").Int(), 0).Format("2006-01-02 15:04:05"),
				)
				priceMessage += vPrice
			}
			ctx.SendChain(message.Image(itemCcon), message.Text(priceMessage))
		})
}
