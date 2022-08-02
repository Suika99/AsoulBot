// Package manager 群管
package manager

import (
	"fmt"
	zero "github.com/wdvxdr1123/ZeroBot"
	"github.com/wdvxdr1123/ZeroBot/message"

	"github.com/FloatTech/ZeroBot-Plugin/plugin/manager/timer"
	sql "github.com/FloatTech/sqlite"
	"github.com/FloatTech/zbputils/control"
	"github.com/FloatTech/zbputils/math"
)

const (
	hint = "====群管====\n" +
		"- 禁言@QQ 1分钟\n" +
		"- 解除禁言 @QQ\n" +
		"- 我要自闭 1分钟\n" +
		"- 开启全员禁言\n" +
		"- 解除全员禁言\n" +
		"- 升为管理@QQ\n" +
		"- 取消管理@QQ\n" +
		"- 修改名片@QQ XXX\n" +
		"- 修改头衔@QQ XXX\n" +
		"- 申请头衔 XXX\n" +
		"- 踢出群聊@QQ\n" +
		"- 入群/退群事件推送"
)

var (
	db    = &sql.Sqlite{}
	clock timer.Clock
)

func init() { // 插件主体
	engine := control.Register("manager", &control.Options{
		DisableOnDefault:  false,
		Help:              hint,
		PrivateDataFolder: "manager",
	})

	// 升为管理
	engine.OnRegex(`^升为管理.*?(\d+)`, zero.OnlyGroup, zero.SuperUserPermission).SetBlock(true).
		Handle(func(ctx *zero.Ctx) {
			ctx.SetGroupAdmin(
				ctx.Event.GroupID,
				math.Str2Int64(ctx.State["regex_matched"].([]string)[1]), // 被升为管理的人的qq
				true,
			)
			nickname := ctx.GetGroupMemberInfo( // 被升为管理的人的昵称
				ctx.Event.GroupID,
				math.Str2Int64(ctx.State["regex_matched"].([]string)[1]), // 被升为管理的人的qq
				false,
			).Get("nickname").Str
			ctx.SendChain(message.Text(nickname + " 升为了管理~"))
		})
	// 取消管理
	engine.OnRegex(`^取消管理.*?(\d+)`, zero.OnlyGroup, zero.SuperUserPermission).SetBlock(true).
		Handle(func(ctx *zero.Ctx) {
			ctx.SetGroupAdmin(
				ctx.Event.GroupID,
				math.Str2Int64(ctx.State["regex_matched"].([]string)[1]), // 被取消管理的人的qq
				false,
			)
			nickname := ctx.GetGroupMemberInfo( // 被取消管理的人的昵称
				ctx.Event.GroupID,
				math.Str2Int64(ctx.State["regex_matched"].([]string)[1]), // 被取消管理的人的qq
				false,
			).Get("nickname").Str
			ctx.SendChain(message.Text("残念~ " + nickname + " 暂时失去了管理员的资格"))
		})

		//	engine.OnKeyword("haha").SetBlock(true).
		//		Handle(func(ctx *zero.Ctx) {
		//			ctx.DeleteMessage(message.NewMessageIDFromInteger(ctx.Event.MessageID.(int64)))

	// 踢出群聊
	engine.OnRegex(`^踢出群聊.*?(\d+)`, zero.OnlyGroup, zero.AdminPermission).SetBlock(true).
		Handle(func(ctx *zero.Ctx) {
			ctx.SetGroupKick(
				ctx.Event.GroupID,
				math.Str2Int64(ctx.State["regex_matched"].([]string)[1]), // 被踢出群聊的人的qq
				false,
			)
			nickname := ctx.GetGroupMemberInfo( // 被踢出群聊的人的昵称
				ctx.Event.GroupID,
				math.Str2Int64(ctx.State["regex_matched"].([]string)[1]), // 被踢出群聊的人的qq
				false,
			).Get("nickname").Str
			ctx.SendChain(message.Text("残念~ " + nickname + " 被放逐"))
		})
	// 退出群聊
	engine.OnRegex(`^退出群聊.*?(\d+)`, zero.OnlyToMe, zero.SuperUserPermission).SetBlock(true).
		Handle(func(ctx *zero.Ctx) {
			ctx.SetGroupLeave(
				math.Str2Int64(ctx.State["regex_matched"].([]string)[1]), // 要退出的群的群号
				true,
			)
		})
	// 开启全体禁言
	engine.OnRegex(`^开启全员禁言$`, zero.OnlyGroup, zero.AdminPermission).SetBlock(true).
		Handle(func(ctx *zero.Ctx) {
			ctx.SetGroupWholeBan(
				ctx.Event.GroupID,
				true,
			)
			ctx.SendChain(message.Text("全员自闭开始~"))
		})
	// 解除全员禁言
	engine.OnRegex(`^解除全员禁言$`, zero.OnlyGroup, zero.AdminPermission).SetBlock(true).
		Handle(func(ctx *zero.Ctx) {
			ctx.SetGroupWholeBan(
				ctx.Event.GroupID,
				false,
			)
			ctx.SendChain(message.Text("全员自闭结束~"))
		})
	// 禁言
	engine.OnRegex(`^禁言.*?(\d+).*?\s(\d+)(.*)`, zero.OnlyGroup, zero.AdminPermission).SetBlock(true).
		Handle(func(ctx *zero.Ctx) {
			duration := math.Str2Int64(ctx.State["regex_matched"].([]string)[2])
			switch ctx.State["regex_matched"].([]string)[3] {
			case "分钟":
				//
			case "小时":
				duration *= 60
			case "天":
				duration *= 60 * 24
			default:
				//
			}
			if duration >= 43200 {
				duration = 43199 // qq禁言最大时长为一个月
			}
			ctx.SetGroupBan(
				ctx.Event.GroupID,
				math.Str2Int64(ctx.State["regex_matched"].([]string)[1]), // 要禁言的人的qq
				duration*60, // 要禁言的时间（分钟）
			)
			ctx.SendChain(message.Text("小黑屋收留成功~"))
		})
	// 解除禁言
	engine.OnRegex(`^解除禁言.*?(\d+)`, zero.OnlyGroup, zero.AdminPermission).SetBlock(true).
		Handle(func(ctx *zero.Ctx) {
			ctx.SetGroupBan(
				ctx.Event.GroupID,
				math.Str2Int64(ctx.State["regex_matched"].([]string)[1]), // 要解除禁言的人的qq
				0,
			)
			ctx.SendChain(message.Text("小黑屋释放成功~"))
		})
	// 自闭禁言
	engine.OnRegex(`^(我要自闭|禅定).*?(\d+)(.*)`, zero.OnlyGroup).SetBlock(true).
		Handle(func(ctx *zero.Ctx) {
			duration := math.Str2Int64(ctx.State["regex_matched"].([]string)[2])
			switch ctx.State["regex_matched"].([]string)[3] {
			case "分钟", "min", "mins", "m":
				break
			case "小时", "hour", "hours", "h":
				duration *= 60
			case "天", "day", "days", "d":
				duration *= 60 * 24
			default:
				break
			}
			if duration >= 43200 {
				duration = 43199 // qq禁言最大时长为一个月
			}
			ctx.SetGroupBan(
				ctx.Event.GroupID,
				ctx.Event.UserID,
				duration*60, // 要自闭的时间（分钟）
			)
			ctx.SendChain(message.Text("那我就不手下留情了~"))
		})
	// 修改名片
	engine.OnRegex(`^修改名片.*?(\d+).*?\s(.*)`, zero.OnlyGroup, zero.AdminPermission).SetBlock(true).
		Handle(func(ctx *zero.Ctx) {
			if len(ctx.State["regex_matched"].([]string)[2]) > 60 {
				ctx.SendChain(message.Text("名字太长啦！"))
				return
			}
			ctx.SetGroupCard(
				ctx.Event.GroupID,
				math.Str2Int64(ctx.State["regex_matched"].([]string)[1]), // 被修改群名片的人
				ctx.State["regex_matched"].([]string)[2],                 // 修改成的群名片
			)
			ctx.SendChain(message.Text("嗯！已经修改了"))
		})
	// 修改头衔
	engine.OnRegex(`^修改头衔.*?(\d+).*?\s(.*)`, zero.OnlyGroup, zero.AdminPermission).SetBlock(true).
		Handle(func(ctx *zero.Ctx) {
			if len(ctx.State["regex_matched"].([]string)[1]) > 18 {
				ctx.SendChain(message.Text("头衔太长啦！"))
				return
			}
			ctx.SetGroupSpecialTitle(
				ctx.Event.GroupID,
				math.Str2Int64(ctx.State["regex_matched"].([]string)[1]), // 被修改群头衔的人
				ctx.State["regex_matched"].([]string)[2],                 // 修改成的群头衔
			)
			ctx.SendChain(message.Text("嗯！已经修改了"))
		})
	// 申请头衔
	engine.OnRegex(`^申请头衔(.*)`, zero.OnlyGroup).SetBlock(true).
		Handle(func(ctx *zero.Ctx) {
			if len(ctx.State["regex_matched"].([]string)[1]) > 18 {
				ctx.SendChain(message.Text("头衔太长啦！"))
				return
			}
			ctx.SetGroupSpecialTitle(
				ctx.Event.GroupID,
				ctx.Event.UserID,                         // 被修改群头衔的人
				ctx.State["regex_matched"].([]string)[1], // 修改成的群头衔
			)
			ctx.SendChain(message.Text("嗯！不错的头衔呢~"))
		})

	// 入群欢迎
	engine.OnNotice().SetBlock(false).
		Handle(func(ctx *zero.Ctx) {
			if ctx.Event.NoticeType == "group_increase" {
				ctx.SendChain(message.At(ctx.Event.UserID),
					message.Text("耶！是新的冒险家！！\n欢迎欢迎~~和派蒙一起开始新的冒险吧！"))
			}
		})

	// 退群提醒
	engine.OnNotice().SetBlock(false).
		Handle(func(ctx *zero.Ctx) {
			if ctx.Event.NoticeType == "group_decrease" {
				ctx.SendChain(message.Text(fmt.Sprintf("%v 背叛了璃月永远的离开了我们", ctx.CardOrNickName(ctx.Event.UserID))))
			}
		})
}
