package minecraft

import (
	"fmt"
	"github.com/FloatTech/zbputils/file"
	"github.com/FloatTech/zbputils/img/text"
	"github.com/FloatTech/zbputils/img/writer"
	"github.com/fogleman/gg"
	log "github.com/sirupsen/logrus"
	zero "github.com/wdvxdr1123/ZeroBot"
	"github.com/wdvxdr1123/ZeroBot/message"
	"image/color"
	"os"
	"time"
)

func init() {
	engine.OnRegex(`^add (.+) (.+):(.+)$`, zero.SuperUserPermission).
		Handle(func(ctx *zero.Ctx) {
			server := &Mcserver{
				Name: ctx.State["regex_matched"].([]string)[1],
				IP:   ctx.State["regex_matched"].([]string)[2],
				Port: ctx.State["regex_matched"].([]string)[3],
			}
			initDB()
			err := insertData(server)
			if err != nil {
				ctx.SendChain(message.Text("插入数据错误: ", err))
				return
			}
			ctx.SendChain(message.Text(
				"Server Add Success!\n",
				"Name：", server.Name, "\n",
				"IP：", server.IP, "\n",
				"Port：", server.Port))
		})

	engine.OnRegex(`^del (.+)$`, zero.SuperUserPermission).
		Handle(func(ctx *zero.Ctx) {
			delName := ctx.State["regex_matched"].([]string)[1]
			err := deleteData(delName)
			if err != nil {
				ctx.SendChain(message.Text("删除数据错误: ", err))
				return
			}
			ctx.SendChain(message.Text(delName, " 删除成功"))
		})

	engine.OnRegex(`^/list (.+)$`).
		Handle(func(ctx *zero.Ctx) {
			addName := ctx.State["regex_matched"].([]string)[1]
			addr, err := selectData(addName)
			if err != nil {
				ctx.SendChain(message.Text("查询数据失败:", err))
				return
			}
			info := getServerInfo(&addr)
			imgStr, err := rInfoImg(addr.Name, info)
			ctx.SendChain(message.Image(imgStr))
		})
}

func rInfoImg(name string, info *ServerInfo) (imgStr string, err error) {
	fontSize := 40.0
	playArr := info.Players.Sample
	canvas := gg.NewContext(600, 340+len(playArr)*80)
	canvas.SetColor(color.White)
	canvas.Clear()
	canvas.SetColor(color.Black)
	if err = canvas.LoadFontFace(text.BoldFontFile, fontSize); err != nil {
		log.Errorln("[minecraft]:", err)
		return
	}
	canvas.DrawString(fmt.Sprintf("服务器名称：%s", name), 100, 100)
	canvas.DrawString(fmt.Sprintf("在线状态：%v", info.Online), 100, 180)
	canvas.DrawString(fmt.Sprintf("在线人数：%d/%d", info.Players.Now, info.Players.Max), 100, 260)
	for i, v := range playArr {
		canvas.SetColor(color.Black)
		canvas.DrawString(v.Name, 100, float64(i*80+340))
	}
	tmp := time.Now().Unix()
	imgPath := "data/minecraft/images/"
	imgFile := fmt.Sprintf("%s%d.png", imgPath, tmp)
	if _, err = os.Stat(imgPath); err != nil || os.IsNotExist(err) {
		_ = os.MkdirAll(imgPath, 0755)
	}
	f, err := os.Create(imgFile)
	if err != nil {
		log.Errorln("[minecraft]:", err)
		data, cl := writer.ToBytes(canvas.Image())
		fmt.Println(data)
		cl()
	}
	_, err = writer.WriteTo(canvas.Image(), f)
	_ = f.Close()
	if err != nil {
		log.Errorln("[minecraft]:", err)
		return
	}
	imgStr = "file:///" + file.BOTPATH + "/" + imgFile
	return
}
