package main

import (
	"flag"
	"fmt"
	"math/rand"
	"os"
	"time"

	"github.com/FloatTech/ZeroBot-Plugin/kanban" // 在最前打印 banner

	_ "github.com/Suika99/AsoulBot/plugins/asoul"
	_ "github.com/Suika99/AsoulBot/plugins/bilibili"
	_ "github.com/Suika99/AsoulBot/plugins/covid"
	_ "github.com/Suika99/AsoulBot/plugins/emo"
	_ "github.com/Suika99/AsoulBot/plugins/emojimix"
	_ "github.com/Suika99/AsoulBot/plugins/ffxiv"
	_ "github.com/Suika99/AsoulBot/plugins/manager"
	_ "github.com/Suika99/AsoulBot/plugins/minecraft"
	_ "github.com/Suika99/AsoulBot/plugins/music"
	_ "github.com/Suika99/AsoulBot/plugins/nativesetu"
	_ "github.com/Suika99/AsoulBot/plugins/setu"

	// -----------------------以下为内置依赖，勿动------------------------ //
	"github.com/FloatTech/zbputils/process"
	"github.com/sirupsen/logrus"
	zero "github.com/wdvxdr1123/ZeroBot"
	"github.com/wdvxdr1123/ZeroBot/driver"
	"github.com/wdvxdr1123/ZeroBot/message"
)

var (
	nicks  = []string{"猫猫"}
	token  *string
	url    *string
	adana  *string
	prefix *string
)

func init() {
	// 解析命令行参数
	d := flag.Bool("d", false, "Enable debug level log and higher.")
	w := flag.Bool("w", false, "Enable warning level log and higher.")
	h := flag.Bool("h", false, "Display this help.")
	// 解析命令行参数，输入 `-g 监听地址:端口` 指定 gui 访问地址，默认 127.0.0.1:3000
	// g := flag.String("g", "127.0.0.1:3000", "Set web gui listening address.")

	// 直接写死 AccessToken 时，请更改下面第二个参数
	token = flag.String("t", "", "Set AccessToken of WSClient.")
	// 直接写死 URL 时，请更改下面第二个参数
	url = flag.String("u", "ws://go-cqhttp:6700", "Set Url of WSClient.")
	// 默认昵称
	adana = flag.String("n", "椛椛", "Set default nickname.")
	prefix = flag.String("p", "/", "Set command prefix.")

	flag.Parse()
	if *h {
		kanban.PrintBanner()
		fmt.Println("Usage:")
		flag.PrintDefaults()
		os.Exit(0)
	} else {
		if *d && !*w {
			logrus.SetLevel(logrus.DebugLevel)
		}
		if *w {
			logrus.SetLevel(logrus.WarnLevel)
		}
	}

	// 启用 gui
	// webctrl.InitGui(*g)
}

func main() {
	rand.Seed(time.Now().UnixNano()) // 全局 seed，其他插件无需再 seed
	// 帮助
	zero.OnFullMatchGroup([]string{"/help", ".help", "菜单"}).SetBlock(true).
		Handle(func(ctx *zero.Ctx) {
			ctx.SendChain(message.Text(kanban.Banner, "\n可发送\"/服务列表\"查看 bot 功能"))
		})
	zero.OnFullMatch("查看zbp公告", zero.OnlyToMe, zero.AdminPermission).SetBlock(true).
		Handle(func(ctx *zero.Ctx) {
			ctx.SendChain(message.Text(kanban.Kanban()))
		})
	zero.RunAndBlock(
		zero.Config{
			NickName:      append([]string{*adana}, nicks...),
			CommandPrefix: *prefix,
			// SuperUsers 某些功能需要主人权限，可通过以下两种方式修改
			// SuperUsers: []string{"12345678", "87654321"}, // 通过代码写死的方式添加主人账号
			SuperUsers: []int64{290760339}, // 通过命令行参数的方式添加主人账号
			Driver:     []zero.Driver{driver.NewWebSocketClient(*url, *token)},
		},
		process.GlobalInitMutex.Unlock,
	)
}
