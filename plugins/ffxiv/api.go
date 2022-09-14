package ffxiv

import (
	"fmt"
	"github.com/FloatTech/zbputils/web"
	log "github.com/sirupsen/logrus"
	"github.com/tidwall/gjson"
	"github.com/wdvxdr1123/ZeroBot/message"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
)

func getItemSearch(itemName string) gjson.Result {
	api := fmt.Sprintf("https://cafemaker.wakingsands.com/search?indexes=Item&string=%s", itemName)
	client := &http.Client{}
	req, err := http.NewRequest("GET", api, nil)
	if err != nil {
		log.Errorln("[ff14]", err)
	}

	res, err := client.Do(req)
	if err != nil {
		log.Errorln("[ff14]", err)
	}

	defer res.Body.Close()
	data, err := io.ReadAll(res.Body)
	if err != nil {
		log.Errorln("[ff14]", err)
	}

	json := gjson.ParseBytes(data)

	return json
}

func getCafeSearch(keyword string) gjson.Result {
	api := fmt.Sprintf("https://cafemaker.wakingsands.com/search?string=%s", keyword)
	client := &http.Client{}
	req, err := http.NewRequest("GET", api, nil)
	if err != nil {
		log.Errorln("[ff14]", err)
	}

	res, err := client.Do(req)
	if err != nil {
		log.Errorln("[ff14]", err)
	}

	defer res.Body.Close()
	data, err := io.ReadAll(res.Body)
	if err != nil {
		log.Errorln("[ff14]", err)
	}

	json := gjson.ParseBytes(data)

	return json
}

func getItemData(itemId int64) gjson.Result {
	api := fmt.Sprintf("https://garlandtools.cn/db/doc/item/chs/3/%d.json", itemId)
	client := &http.Client{}
	req, err := http.NewRequest("GET", api, nil)
	if err != nil {
		log.Errorln("[ff14]", err)
	}

	res, err := client.Do(req)
	if err != nil {
		log.Errorln("[ff14]", err)
	}

	defer res.Body.Close()
	data, err := io.ReadAll(res.Body)
	if err != nil {
		log.Errorln("[ff14]", err)
	}

	json := gjson.ParseBytes(data)
	return json
}

func getItemPrice(itemId int64) gjson.Result {
	api := fmt.Sprintf("https://universalis.app/api/v2/莫古力/%d?listings=10", itemId)
	client := &http.Client{}
	req, err := http.NewRequest("GET", api, nil)
	if err != nil {
		log.Errorln("[ff14]", err)
	}

	res, err := client.Do(req)
	if err != nil {
		log.Errorln("[ff14]", err)
	}

	defer res.Body.Close()
	data, err := io.ReadAll(res.Body)
	if err != nil {
		log.Errorln("[ff14]", err)
	}

	json := gjson.ParseBytes(data)
	return json
}

func getItemType(categoryId int64) gjson.Result {
	jsonFile, err := os.Open("data/ffxiv/category.json")
	if err != nil {
		log.Errorln("[ff14]", err)
	}
	defer jsonFile.Close()

	data, err := ioutil.ReadAll(jsonFile)
	if err != nil {
		log.Errorln("[ff14]", err)
	}

	json := gjson.ParseBytes(data)

	categoryIndex := fmt.Sprintf("item.categoryIndex.%d.name", categoryId)
	itemType := json.Get(categoryIndex)
	return itemType
}

// cloud163 返回网易云音乐卡片
func cloud163(keyword string) (msg message.MessageSegment) {
	requestURL := "https://music.cyrilstudio.top/search?keywords=" + url.QueryEscape(keyword)
	data, err := web.GetData(requestURL)
	if err != nil {
		msg = message.Text("ERROR:", err)
		return
	}
	msg = message.Music("163", gjson.ParseBytes(data).Get("result.songs.0.id").Int())
	return
}
