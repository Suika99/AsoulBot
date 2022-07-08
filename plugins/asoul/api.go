package asoul

import (
	"encoding/json"
	"errors"
	"github.com/FloatTech/zbputils/binary"
	"github.com/FloatTech/zbputils/file"
	"github.com/FloatTech/zbputils/web"
	"github.com/jinzhu/gorm"
	log "github.com/sirupsen/logrus"
	"github.com/tidwall/gjson"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"time"
)

// 通过触发指令的名字获取用户的uid
func getMid(keyword string) (gjson.Result, error) {
	api := "http://api.bilibili.com/x/web-interface/search/type?search_type=bili_user&user_type=0&keyword=" + keyword
	resp, err := http.Get(api)
	if err != nil {
		return gjson.Result{}, err
	}
	if resp.StatusCode != http.StatusOK {
		return gjson.Result{}, errors.New("code not 200")
	}
	data, _ := ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	json := gjson.ParseBytes(data)
	if json.Get("data.numResults").Int() == 0 {
		return gjson.Result{}, errors.New("查无此人")
	}
	return json, nil
}

// 获取被查用户信息
func getinfo(mid string) *follows {
	url := "https://account.bilibili.com/api/member/getCardByMid?mid=" + mid
	resp, err := http.Get(url)
	if err != nil {
		log.Errorln(err)
	}
	defer resp.Body.Close()
	result := &follows{}
	if err := json.NewDecoder(resp.Body).Decode(result); err != nil {
		log.Errorln(err)
	}
	return result
}

// 获取牌子信息
func medalwall(uid string, cookie string) (result []medal, err error) {
	medalwallURL := "https://api.live.bilibili.com/xlive/web-ucenter/user/MedalWall?target_id=" + uid
	method := "GET"
	client := &http.Client{}
	req, err := http.NewRequest(method, medalwallURL, nil)
	if err != nil {
		return
	}
	//c := vdb.getBilibiliCookie()
	req.Header.Add("cookie", cookie)
	res, err := client.Do(req)
	if err != nil {
		return
	}
	defer res.Body.Close()
	data, err := io.ReadAll(res.Body)
	if err != nil {
		return
	}
	j := gjson.ParseBytes(data)
	if j.Get("code").Int() == -101 {
		err = errNeedCookie
		return
	}
	if j.Get("code").Int() != 0 {
		err = errors.New(j.Get("message").String())
	}
	_ = json.Unmarshal(binary.StringToBytes(j.Get("data.list").Raw), &result)
	return
}

// 获取头像图片
func initFacePic(filename, faceURL string) error {
	if file.IsNotExist(filename) {
		data, err := web.GetData(faceURL)
		if err != nil {
			return err
		}
		err = os.WriteFile(filename, data, 0666)
		if err != nil {
			return err
		}
	}
	return nil
}

// 首次启动初始化插件, 异步处理！！
// 获取vtbs数据返回
func init() {
	go func() {
		url := "https://api.vtbs.moe/v1/short"
		method := "GET"
		client := &http.Client{}
		req, err := http.NewRequest(method, url, nil)
		if err != nil {
			log.Error("[Element]请求api失败", err)
		}
		res, err := client.Do(req)
		if err != nil {
			log.Error("[Element]请求api失败")
			return
		}
		defer res.Body.Close()
		body, err := ioutil.ReadAll(res.Body)
		if err != nil {
			log.Errorln(err)
		}
		json := gjson.ParseBytes(body)

		if _, err = os.Stat(dbfile); err != nil || os.IsNotExist(err) {
			// 生成文件
			var err error
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
			for _, i := range json.Array() {
				db.Create(&vup{
					Mid:    i.Get("mid").Int(),
					Uname:  i.Get("uname").Str,
					Roomid: i.Get("roomid").Int(),
				})
			}
			log.Infof("[Element]vtbs更新完成，插入（%v）条数据", len(json.Array()))
			defer db.Close()
		}
	}()
}

// 对比数据库获取关注用户的名字
//func compared(follows []int64) []string {
//	var db *sqlx.DB
//	db, _ = sqlx.Open("sqlite3", dbfile)
//	defer db.Close()
//	query1, args, err := sqlx.In("select uname from vtbs where mid in (?)", follows)
//	if err != nil {
//		log.Errorln("[element]查找失败", err)
//	}
//	res := []string{}
//	err = db.Select(&res, query1, args...)
//	if err != nil {
//		log.Errorln("[element]查找失败", err)
//	}
//	return res
//}
