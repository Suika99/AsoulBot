package ffxiv

import (
	"fmt"
	"github.com/FloatTech/zbputils/binary"
	"github.com/FloatTech/zbputils/img/text"
	zero "github.com/wdvxdr1123/ZeroBot"
	"github.com/wdvxdr1123/ZeroBot/message"
	"strings"
	"time"
)

func init() {
	engine.OnRegex(`^p\s(.{1,25})$`).SetBlock(true).
		Handle(func(ctx *zero.Ctx) {
			keyword := ctx.State["regex_matched"].([]string)[1]
			cafeRes := getItemSearch(keyword)
			itemID := cafeRes.Get("Results.#.ID").Array()[0].Int()
			itemCon := fmt.Sprintf("https://cafemaker.wakingsands.com/%s", cafeRes.Get("Results.#.Icon").Array()[0].Str)
			priceData := getItemPrice(itemID).Get("listings").Array()
			itemData := getItemData(itemID)
			var priceMessage string
			for _, v := range priceData {
				vPrice := fmt.Sprintf("%v     x%v=%v     %v     %v:%v\n",
					v.Get("pricePerUnit"),
					v.Get("quantity"),
					v.Get("total"),
					time.Unix(v.Get("lastReviewTime").Int(), 0).Format("2006-01-02 15:04:05"),
					v.Get("worldName"),
					v.Get("retainerName"),
				)
				priceMessage += vPrice
			}

			priceText := priceMessage[:len(priceMessage)-1]
			b, err := text.RenderToBase64(priceText, text.BoldFontFile, 1200, 35)
			if err != nil {
				ctx.SendChain(message.Text("ERROR: ", err))
				return
			}

			ctx.SendChain(message.Image(itemCon),
				message.Text(
					itemData.Get("item.name"), "\n",
					"类别：", getItemType(itemData.Get("item.category").Int()), "\n",
					"版本：", itemData.Get("item.patch"), "\n",
					"物品等级：", itemData.Get("item.ilvl"), "\n",
					"物品描述：", strings.Replace(itemData.Get("item.description").Str, "<br>", "\n", -1), "\n",
					"物品来源：", fmt.Sprintf("https://garlandtools.cn/db/#item/%d", itemID),
				),
				message.Image("base64://"+binary.BytesToString(b)))
		})
}
