package asoul

import (
	"fmt"
	"github.com/FloatTech/zbputils/file"
	"github.com/FloatTech/zbputils/img"
	"github.com/FloatTech/zbputils/img/text"
	"github.com/FloatTech/zbputils/img/writer"
	"github.com/fogleman/gg"
	log "github.com/sirupsen/logrus"
	zero "github.com/wdvxdr1123/ZeroBot"
	"github.com/wdvxdr1123/ZeroBot/message"
	"image"
	"image/color"
	"os"
	"path"
	"strconv"
	"time"
)

// 指令触发查询 asoul 粉丝数据
func init() {
	engine.OnKeyword("粉丝信息").
		Handle(func(ctx *zero.Ctx) {
			asoulFansImg := rAsoulImage()
			ctx.SendChain(message.Image(asoulFansImg))
		})

	engine.OnRegex(`^查粉丝\s(.{1,25})$`).
		Handle(func(ctx *zero.Ctx) {
			keyword := ctx.State["regex_matched"].([]string)[1]
			rest, err := getMid(keyword)
			if err != nil {
				ctx.SendChain(message.Text("ERROR: ", err))
				return
			}
			mid := rest.Get("data.result.0.mid").String()
			fansImg := rFansImage1(mid)
			ctx.SendChain(message.Image(fansImg))
		})
}

func rFansImage1(mid string) (fansImg string) {
	var err error
	uid, _ := strconv.Atoi(mid)
	fansData := fansapi(uid)
	facePath := "data/cache/" + mid + "vupFace" + path.Ext(fansData.Face)
	fansFile := "data/cache/" + mid + "fans.png"
	var back image.Image
	if path.Ext(fansData.Face) != ".webp" {
		err = initFacePic(facePath, fansData.Face)
		if err != nil {
			log.Errorln("cha ERROR:", err)
			return
		}
		back, err = gg.LoadImage(facePath)
		if err != nil {
			log.Errorln("cha ERROR:", err)
			return
		}
		back = img.Size(back, 500, 500).Im
	}
	canvas := gg.NewContext(1200, 550)
	fontSize := 40.0
	canvas.SetColor(color.White)
	canvas.Clear()
	if back != nil {
		canvas.DrawImage(back, 0, 0)
	}
	canvas.SetColor(color.Black)
	_, err = file.GetLazyData(text.BoldFontFile, true)
	if err != nil {
		log.Errorln("cha ERROR:", err)
		return
	}
	if err = canvas.LoadFontFace(text.BoldFontFile, fontSize); err != nil {
		log.Errorln("cha ERROR:", err)
		return
	}

	length, h := canvas.MeasureString(strconv.Itoa(fansData.Mid))
	n, _ := canvas.MeasureString(fansData.Uname)
	canvas.DrawString(fansData.Uname, 550, 100-h)
	canvas.DrawRoundedRectangle(600+n-length*0.1, 100-h*2.5, length*1.2, h*2, fontSize*0.2)
	canvas.SetRGB255(221, 221, 221)
	canvas.Fill()
	canvas.SetColor(color.Black)
	canvas.DrawString(strconv.Itoa(fansData.Mid), 600+n, 100-h)
	canvas.DrawString(fmt.Sprintf("舰长：%d", fansData.GuardNum), 550, 180-h)
	canvas.DrawString(fmt.Sprintf("视频投稿：%d", fansData.Video), 550, 260-h)
	canvas.DrawString(fmt.Sprintf("当前粉丝：%d", fansData.Follower), 550, 340-h)
	canvas.DrawString("24H涨粉：", 550, 420-h)
	canvas.DrawString("查询日期："+time.Now().Format("2006-01-02"), 550, 500-h)
	if fansData.Rise <= 0 {
		canvas.SetRGB255(255, 0, 0)
		canvas.DrawString(strconv.Itoa(fansData.Rise), 750, 420-h)
	} else {
		canvas.SetRGB255(0, 255, 0)
		canvas.DrawString(fmt.Sprintf("+%d", fansData.Rise), 750, 420-h)
	}

	f, err := os.Create(fansFile)
	if err != nil {
		log.Errorln("cha ERROR:", err)
		data, cl := writer.ToBytes(canvas.Image())
		fmt.Println(data)
		cl()
	}
	_, err = writer.WriteTo(canvas.Image(), f)
	_ = f.Close()
	if err != nil {
		log.Errorln("cha ERROR:", err)
		return
	}
	fansImg = "file:///" + file.BOTPATH + "/" + fansFile
	return fansImg
}

func rAsoulImage() (asoulFansImg string) {
	var (
		diana  = fansapi(diana)
		ava    = fansapi(ava)
		eileen = fansapi(eileen)
		bella  = fansapi(bella)
		asoul  = fansapi(asoul)
	)
	var err error
	fontSize := 80.0
	id := strconv.FormatInt(time.Now().Unix(), 10)
	fansFile := "data/cache/" + id + "fans.png"
	canvas := gg.NewContext(2000, 950)
	canvas.SetColor(color.White)
	canvas.Clear()
	_, err = file.GetLazyData(text.BoldFontFile, true)
	if err != nil {
		log.Errorln("cha ERROR:", err)
		return
	}
	if err = canvas.LoadFontFace(text.BoldFontFile, fontSize); err != nil {
		log.Errorln("cha ERROR:", err)
		return
	}
	canvas.SetColor(color.Black)
	canvas.DrawString("名字", 100, 100)
	canvas.DrawString("当前粉丝", 850, 100)
	canvas.DrawString("24小时涨粉", 1400, 100)
	canvas.SetColor(color.Black)
	canvas.DrawString(ava.Uname, 100, 250)
	canvas.DrawString(strconv.Itoa(ava.Follower), 850, 250)
	if ava.Rise <= 0 {
		canvas.SetColor(color.RGBA{255, 0, 0, 255})
	} else {
		canvas.SetColor(color.RGBA{0, 255, 0, 255})
	}
	canvas.DrawString(strconv.Itoa(ava.Rise), 1400, 250)
	canvas.SetColor(color.Black)
	canvas.DrawString(bella.Uname, 100, 400)
	canvas.DrawString(strconv.Itoa(bella.Follower), 850, 400)
	if bella.Rise <= 0 {
		canvas.SetColor(color.RGBA{255, 0, 0, 255})
	} else {
		canvas.SetColor(color.RGBA{0, 255, 0, 255})
	}
	canvas.DrawString(strconv.Itoa(bella.Rise), 1400, 400)
	canvas.SetColor(color.Black)
	canvas.DrawString(diana.Uname, 100, 550)
	canvas.DrawString(strconv.Itoa(diana.Follower), 850, 550)
	if diana.Rise <= 0 {
		canvas.SetColor(color.RGBA{255, 0, 0, 255})
	} else {
		canvas.SetColor(color.RGBA{0, 255, 0, 255})
	}
	canvas.DrawString(strconv.Itoa(diana.Rise), 1400, 550)
	canvas.SetColor(color.Black)
	canvas.DrawString(eileen.Uname, 100, 700)
	canvas.DrawString(strconv.Itoa(eileen.Follower), 850, 700)
	if eileen.Rise <= 0 {
		canvas.SetColor(color.RGBA{255, 0, 0, 255})
	} else {
		canvas.SetColor(color.RGBA{0, 255, 0, 255})
	}
	canvas.DrawString(strconv.Itoa(eileen.Rise), 1400, 700)
	canvas.SetColor(color.Black)
	canvas.DrawString(asoul.Uname, 100, 850)
	canvas.DrawString(strconv.Itoa(asoul.Follower), 850, 850)
	if asoul.Rise <= 0 {
		canvas.SetColor(color.RGBA{255, 0, 0, 255})
	} else {
		canvas.SetColor(color.RGBA{0, 255, 0, 255})
	}
	canvas.DrawString(strconv.Itoa(asoul.Rise), 1400, 850)

	f, err := os.Create(fansFile)
	if err != nil {
		log.Errorln("cha ERROR:", err)
		data, cl := writer.ToBytes(canvas.Image())
		fmt.Println(data)
		cl()
	}
	_, err = writer.WriteTo(canvas.Image(), f)
	_ = f.Close()
	if err != nil {
		log.Errorln("cha ERROR:", err)
		return
	}
	asoulFansImg = "file:///" + file.BOTPATH + "/" + fansFile
	return asoulFansImg
}
