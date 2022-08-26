package ffxiv

import (
	"fmt"
	zero "github.com/wdvxdr1123/ZeroBot"
	"github.com/wdvxdr1123/ZeroBot/message"
	"strings"
)

func init() {
	engine.OnRegex(`^s\s(.{1,25})$`).SetBlock(true).
		Handle(func(ctx *zero.Ctx) {
			keyword := ctx.State["regex_matched"].([]string)[1]
			results := getSearch(keyword)
			if len(results.Raw) == 2 {
				ctx.Send("未搜索到结果123")
				return
			}

			elId := results.Get("#.id").Array()[0].Int()
			resultsType := results.Get("#.type").Array()[0].Str
			cafeRes := getCafeSearch(keyword)
			itemData := getItemData(elId)

			switch resultsType {
			case "instance":
				ctx.SendChain(message.Image(fmt.Sprintf("https://cafemaker.wakingsands.com/%s", cafeRes.Get("Results.#.Icon").Array()[0].Str)),
					message.Text(
						"副本：", results.Get("#.obj.n").Array()[0].Str, "\n",
						fmt.Sprintf("https://garlandtools.cn/db/#instance/%d", elId),
					))

			case "quest":
				ctx.SendChain(message.Text(
					"任务：", results.Get("#.obj.n").Array()[0].Str, "\n",
					"开始位置：", results.Get("#.obj.l").Array()[0].Str, "\n",
					"！注意：该内容可能涉及剧透，请谨慎浏览", "\n",
					fmt.Sprintf("https://garlandtools.cn/db/#quest/%d", elId),
				))

			case "action":
				ctx.SendChain(message.Image(fmt.Sprintf("https://cafemaker.wakingsands.com/%s", cafeRes.Get("Results.#.Icon").Array()[0].Str)),
					message.Text(
						"技能：", results.Get("#.obj.n").Array()[0].Str, "\n",
						fmt.Sprintf("https://garlandtools.cn/db/#action/%d", elId),
					))

			case "item":
				ctx.SendChain(message.Image(fmt.Sprintf("https://cafemaker.wakingsands.com/%s", cafeRes.Get("Results.#.Icon").Array()[0].Str)),
					message.Text(
						itemData.Get("item.name"), "\n",
						"类别：", getItemType(itemData.Get("item.category").Int()), "\n",
						"版本：", itemData.Get("item.patch"), "\n",
						"物品等级：", itemData.Get("item.ilvl"), "\n",
						"物品描述：", strings.Replace(itemData.Get("item.description").Str, "<br>", "\n", -1), "\n",
						"物品来源：", fmt.Sprintf("https://garlandtools.cn/db/#item/%d", elId),
					))
			}
		})
}
