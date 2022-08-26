package ffxiv

import "github.com/FloatTech/zbputils/control"

var (
	engine = control.Register("ff14", &control.Options{
		DisableOnDefault: false,
		Help: "=======ff14相关功能=======\n" +
			"- s 关键字 (搜索功能：支持副本、任务、技能、物品)\n" +
			"- p 物品名 (查询该物品的目前板子价格，默认数据是莫古力大区下)\n",
	})
)
