package asoul

import (
	"encoding/binary"
	"fmt"
	"image"
	"image/color"
	"os"
	"path"
	"sort"
	"strconv"
	"time"

	"github.com/FloatTech/zbputils/file"
	"github.com/FloatTech/zbputils/img"
	"github.com/FloatTech/zbputils/img/text"
	"github.com/FloatTech/zbputils/img/writer"
	"github.com/fogleman/gg"
	"github.com/jinzhu/gorm"
	log "github.com/sirupsen/logrus"
	zero "github.com/wdvxdr1123/ZeroBot"
	"github.com/wdvxdr1123/ZeroBot/message"
)

// 查成分的
func init() {
	// 插件主体,匹配用户名字
	engine.OnRegex(`^查成分\s(.{1,25})$`).SetBlock(true).
		Handle(func(ctx *zero.Ctx) {
			keyword := ctx.State["regex_matched"].([]string)[1]
			rest, err := getMid(keyword)
			if err != nil {
				ctx.SendChain(message.Text("ERROR: ", err))
				return
			}
			mid := rest.Get("data.result.0.mid").String()
			info := getinfo(mid)
			attentions := info.Card.Attentions
			picPath := renderImg(attentions, info)
			ctx.SendChain(message.Image(picPath))
			//msgid := ctx.Event.MessageID
			ctx.DeleteMessage(ctx.Event.MessageID.(message.MessageID))

		})
}

// 渲染图片
func renderImg(attentions []int64, info *follows) (imgPath string) {
	var err error
	id := info.Card.Mid
	facePath := "data/cache/" + id + "vupFace" + path.Ext(info.Card.Face)
	drawedFile := "data/cache/" + id + "vupLike.png"
	vups, err := filterVup(attentions)
	var as []vup
	for _, v := range Asoul {
		for k, w := range vups {
			if v == w.Mid {
				as = append(as, w)
				vups[k] = vups[len(vups)-1]
				vups = vups[:len(vups)-1]
			}
		}
	}

	medals, err := medalwall(info.Card.Mid, cookie)
	sort.Sort(medalSlice(medals))
	if err != nil {
		log.Errorln(err)
	}
	frontVups := make([]vup, 0)
	medalMap := make(map[int64]medal)
	for _, v := range medals {
		up := vup{
			Mid:   v.Mid,
			Uname: v.Uname,
		}
		frontVups = append(frontVups, up)
		medalMap[v.Mid] = v
	}

	backX := 500
	backY := 400
	var back image.Image
	if path.Ext(info.Card.Face) != ".webp" {
		err = initFacePic(facePath, info.Card.Face)
		if err != nil {
			log.Errorln("cha ERROR:", err)
			return
		}
		back, err = gg.LoadImage(facePath)
		if err != nil {
			log.Errorln("cha ERROR:", err)
			return
		}
		back = img.Size(back, backX, 500).Im
	}
	canvas := gg.NewContext(1300, int(450*(1.1+float64(len(vups)+len(as))/3)))
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

	sl, _ := canvas.MeasureString("好")
	length, h := canvas.MeasureString(info.Card.Mid)
	n, _ := canvas.MeasureString(info.Card.Name)
	canvas.DrawString(info.Card.Name, 550, 100-h)
	canvas.DrawRoundedRectangle(600+n-length*0.1, 100-h*2.5, length*1.2, h*2, fontSize*0.2)
	canvas.SetRGB255(221, 221, 221)
	canvas.Fill()
	canvas.SetColor(color.Black)
	canvas.DrawString(info.Card.Mid, 600+n, 100-h)
	canvas.DrawString(fmt.Sprintf("关注：%d", info.Card.Attention), 550, 180-h)
	canvas.DrawString(fmt.Sprintf("粉丝：%d", info.Card.Fans), 650+n, 180-h)
	canvas.DrawString(fmt.Sprintf("关注vup：%d", len(vups)+len(as)), 550, 260-h)
	canvas.DrawString(fmt.Sprintf("使用装扮：%s", info.Card.Pendant.Name), 550, 340-h)
	canvas.DrawString("注册日期："+time.Unix(info.Card.Regtime, 0).Format("2006-01-02 15:04:05"), 550, 420-h)
	canvas.DrawString("查询日期："+time.Now().Format("2006-01-02"), 550, 500-h)
	for i, v := range as {
		canvas.SetColor(color.RGBA{231, 121, 176, 255})
		nl, _ := canvas.MeasureString(v.Uname)
		canvas.DrawString(v.Uname, float64(backX)*0.1, float64(backY)*1.4+float64(i+1)*float64(backY)/3-2*h)
		ml, _ := canvas.MeasureString(strconv.FormatInt(v.Mid, 10))
		canvas.DrawRoundedRectangle(nl-0.1*ml+float64(backX)*0.2, float64(backY)*1.4+float64(i+1)*float64(backY)/3-h*3.5, ml*1.2, h*2, fontSize*0.2)
		canvas.SetRGB255(221, 221, 221)
		canvas.Fill()
		canvas.SetColor(color.RGBA{231, 121, 176, 255})
		canvas.DrawString(strconv.FormatInt(v.Mid, 10), nl+float64(backX)*0.2, float64(backY)*1.4+float64(i+1)*float64(backY)/3-2*h)
		if m, ok := medalMap[v.Mid]; ok {
			mnl, _ := canvas.MeasureString(m.MedalName)
			grad := gg.NewLinearGradient(nl+ml-sl/2+float64(backX)*0.4, float64(backY)*1.4+float64(i+1)*float64(backY)/3-3.5*h, nl+ml+mnl+sl/2+float64(backX)*0.4, float64(backY)*1.4+float64(i+1)*float64(backY)/3-1.5*h)
			r, g, b := int2rbg(m.MedalColorStart)
			grad.AddColorStop(0, color.RGBA{R: uint8(r), G: uint8(g), B: uint8(b), A: 255})
			r, g, b = int2rbg(m.MedalColorEnd)
			grad.AddColorStop(1, color.RGBA{R: uint8(r), G: uint8(g), B: uint8(b), A: 255})
			canvas.SetFillStyle(grad)
			canvas.SetLineWidth(4)
			canvas.MoveTo(nl+ml-sl/2+float64(backX)*0.4, float64(backY)*1.4+float64(i+1)*float64(backY)/3-3.5*h)
			canvas.LineTo(nl+ml+mnl+sl/2+float64(backX)*0.4, float64(backY)*1.4+float64(i+1)*float64(backY)/3-3.5*h)
			canvas.LineTo(nl+ml+mnl+sl/2+float64(backX)*0.4, float64(backY)*1.4+float64(i+1)*float64(backY)/3-1.5*h)
			canvas.LineTo(nl+ml-sl/2+float64(backX)*0.4, float64(backY)*1.4+float64(i+1)*float64(backY)/3-1.5*h)
			canvas.ClosePath()
			canvas.Fill()
			canvas.SetColor(color.White)
			canvas.DrawString(m.MedalName, nl+ml+float64(backX)*0.4, float64(backY)*1.4+float64(i+1)*float64(backY)/3-2*h)
			r, g, b = int2rbg(m.MedalColorBorder)
			canvas.SetRGB255(int(r), int(g), int(b))
			canvas.DrawString(strconv.FormatInt(m.Level, 10), nl+ml+mnl+sl+float64(backX)*0.4, float64(backY)*1.4+float64(i+1)*float64(backY)/3-2*h)
			mll, _ := canvas.MeasureString(strconv.FormatInt(m.Level, 10))
			canvas.SetLineWidth(4)
			canvas.MoveTo(nl+ml-sl/2+float64(backX)*0.4, float64(backY)*1.4+float64(i+1)*float64(backY)/3-3.5*h)
			canvas.LineTo(nl+ml+mnl+mll+sl/2+float64(backX)*0.5, float64(backY)*1.4+float64(i+1)*float64(backY)/3-3.5*h)
			canvas.LineTo(nl+ml+mnl+mll+sl/2+float64(backX)*0.5, float64(backY)*1.4+float64(i+1)*float64(backY)/3-1.5*h)
			canvas.LineTo(nl+ml-sl/2+float64(backX)*0.4, float64(backY)*1.4+float64(i+1)*float64(backY)/3-1.5*h)
			canvas.ClosePath()
			canvas.Stroke()
		}
	}

	for i, v := range vups {
		canvas.SetColor(color.Black)
		nl, _ := canvas.MeasureString(v.Uname)
		canvas.DrawString(v.Uname, float64(backX)*0.1, float64(backY)*1.4+float64(i+1+len(as))*float64(backY)/3-2*h)
		ml, _ := canvas.MeasureString(strconv.FormatInt(v.Mid, 10))
		canvas.DrawRoundedRectangle(nl-0.1*ml+float64(backX)*0.2, float64(backY)*1.4+float64(i+1+len(as))*float64(backY)/3-h*3.5, ml*1.2, h*2, fontSize*0.2)
		canvas.SetRGB255(221, 221, 221)
		canvas.Fill()
		canvas.SetColor(color.Black)
		canvas.DrawString(strconv.FormatInt(v.Mid, 10), nl+float64(backX)*0.2, float64(backY)*1.4+float64(i+1+len(as))*float64(backY)/3-2*h)
		if m, ok := medalMap[v.Mid]; ok {
			mnl, _ := canvas.MeasureString(m.MedalName)
			grad := gg.NewLinearGradient(nl+ml-sl/2+float64(backX)*0.4, float64(backY)*1.4+float64(i+1+len(as))*float64(backY)/3-3.5*h, nl+ml+mnl+sl/2+float64(backX)*0.4, float64(backY)*1.4+float64(i+1+len(as))*float64(backY)/3-1.5*h)
			r, g, b := int2rbg(m.MedalColorStart)
			grad.AddColorStop(0, color.RGBA{R: uint8(r), G: uint8(g), B: uint8(b), A: 255})
			r, g, b = int2rbg(m.MedalColorEnd)
			grad.AddColorStop(1, color.RGBA{R: uint8(r), G: uint8(g), B: uint8(b), A: 255})
			canvas.SetFillStyle(grad)
			canvas.SetLineWidth(4)
			canvas.MoveTo(nl+ml-sl/2+float64(backX)*0.4, float64(backY)*1.4+float64(i+1+len(as))*float64(backY)/3-3.5*h)
			canvas.LineTo(nl+ml+mnl+sl/2+float64(backX)*0.4, float64(backY)*1.4+float64(i+1+len(as))*float64(backY)/3-3.5*h)
			canvas.LineTo(nl+ml+mnl+sl/2+float64(backX)*0.4, float64(backY)*1.4+float64(i+1+len(as))*float64(backY)/3-1.5*h)
			canvas.LineTo(nl+ml-sl/2+float64(backX)*0.4, float64(backY)*1.4+float64(i+1+len(as))*float64(backY)/3-1.5*h)
			canvas.ClosePath()
			canvas.Fill()
			canvas.SetColor(color.White)
			canvas.DrawString(m.MedalName, nl+ml+float64(backX)*0.4, float64(backY)*1.4+float64(i+1+len(as))*float64(backY)/3-2*h)
			r, g, b = int2rbg(m.MedalColorBorder)
			canvas.SetRGB255(int(r), int(g), int(b))
			canvas.DrawString(strconv.FormatInt(m.Level, 10), nl+ml+mnl+sl+float64(backX)*0.4, float64(backY)*1.4+float64(i+1+len(as))*float64(backY)/3-2*h)
			mll, _ := canvas.MeasureString(strconv.FormatInt(m.Level, 10))
			canvas.SetLineWidth(4)
			canvas.MoveTo(nl+ml-sl/2+float64(backX)*0.4, float64(backY)*1.4+float64(i+1+len(as))*float64(backY)/3-3.5*h)
			canvas.LineTo(nl+ml+mnl+mll+sl/2+float64(backX)*0.5, float64(backY)*1.4+float64(i+1+len(as))*float64(backY)/3-3.5*h)
			canvas.LineTo(nl+ml+mnl+mll+sl/2+float64(backX)*0.5, float64(backY)*1.4+float64(i+1+len(as))*float64(backY)/3-1.5*h)
			canvas.LineTo(nl+ml-sl/2+float64(backX)*0.4, float64(backY)*1.4+float64(i+1+len(as))*float64(backY)/3-1.5*h)
			canvas.ClosePath()
			canvas.Stroke()
		}
	}

	f, err := os.Create(drawedFile)
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

	imgPath = "file:///" + file.BOTPATH + "/" + drawedFile
	return imgPath
}

// 牌子数据处理
type medalSlice []medal

func (m medalSlice) Len() int {
	return len(m)
}
func (m medalSlice) Swap(i, j int) {
	m[i], m[j] = m[j], m[i]
}
func (m medalSlice) Less(i, j int) bool {
	return m[i].Level > m[j].Level
}

// filterVup 筛选vup
func filterVup(ids []int64) (vups []vup, err error) {
	db, err := gorm.Open("sqlite3", dbfile)
	if err != nil {
		log.Errorln("[Element]打开数据库失败：", err)
	}
	if err = db.Model(&vup{}).Find(&vups, "mid in (?)", ids).Error; err != nil {
		return vups, err
	}

	defer db.Close()
	return
}

// 转换颜色rbg
func int2rbg(t int64) (int64, int64, int64) {
	var buf [8]byte
	binary.LittleEndian.PutUint64(buf[:], uint64(t))
	b, g, r := int64(buf[0]), int64(buf[1]), int64(buf[2])
	return r, g, b
}

// 获取vups数据入库
func init() {
	go func() {
		var err error
		if _, err = os.Stat(dbfile); err != nil || os.IsNotExist(err) {
			// 生成文件
			_ = os.MkdirAll(datapath, 0755)
			f, err := os.Create(dbfile)
			if err != nil {
				log.Error("[Element]", err)
			}
			log.Infof("[Element]数据库文件(%v)创建成功", dbfile)
			time.Sleep(1 * time.Second)
			defer f.Close()
			// 打开数据库制表
			db, err := gorm.Open("sqlite3", dbfile)
			if err != nil {
				log.Errorln("[Element]打开数据库失败：", err)
			}
			db.AutoMigrate(vup{})
			time.Sleep(1 * time.Second)
			// 插入数据
			vupsData := getVupsData()
			for _, i := range vupsData.Array() {
				db.Create(&vup{
					Mid:    i.Get("mid").Int(),
					Uname:  i.Get("uname").Str,
					Roomid: i.Get("roomid").Int(),
				})
			}
			log.Infof("[Element]vtbs更新完成，插入（%v）条数据", len(vupsData.Array()))
			defer db.Close()
		}
	}()
}
