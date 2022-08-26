package asoul

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/FloatTech/zbputils/binary"
	"github.com/FloatTech/zbputils/file"
	"github.com/FloatTech/zbputils/web"
	log "github.com/sirupsen/logrus"
	"github.com/tidwall/gjson"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
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

// 获取粉丝数据信息
func fansapi(uid int) *follower {
	url := fmt.Sprintf("https://api.vtbs.moe/v1/detail/%d", uid)
	resp, err := http.Get(url)
	if err != nil {
		log.Errorln(err)
	}
	defer resp.Body.Close()
	result := &follower{}
	if err := json.NewDecoder(resp.Body).Decode(result); err != nil {
		log.Errorln(err)
	}
	return result
}

// 获取asoul视频投稿
func getVideo(uid int) *vdInfo {
	url := fmt.Sprintf("https://api.bilibili.com/x/space/arc/search?&ps=50&pn=1&order=pubdate&mid=%d", uid)
	method := "GET"
	client := &http.Client{}
	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		log.Errorln("[video]", err)
	}
	res, err := client.Do(req)
	if err != nil {
		log.Errorln("[video]", err)
	}
	defer res.Body.Close()
	result := &vdInfo{}
	if err := json.NewDecoder(res.Body).Decode(result); err != nil {
		log.Error(err)
	}
	return result
}

// 获取日程表图片
func getDynamic() (dynamicPic []string) {
	api := "https://api.vc.bilibili.com/dynamic_svr/v1/dynamic_svr/space_history?host_uid=703007996"
	resp, err := http.Get(api)
	if err != nil {
		log.Error(err)
	}
	data, _ := ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	json := gjson.ParseBytes(data)

	dy := json.Get("data.cards.#.card").Array()
	for _, i := range dy {
		if strings.Index(i.String(), "日程表") >= 0 {
			org := (gjson.Parse(i.String()).Get("origin"))
			pic := (gjson.Parse(org.String()).Get("item.pictures").Array())
			for _, u := range pic {
				dynamicPic = append(dynamicPic, u.Get("img_src").String())
			}
			return dynamicPic
		}
	}
	return dynamicPic
}

// 获取vups数据
func getVupsData() gjson.Result {
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
	}
	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Errorln(err)
	}
	vupsData := gjson.ParseBytes(body)
	return vupsData
}
