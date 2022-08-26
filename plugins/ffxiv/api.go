package ffxiv

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"github.com/tidwall/gjson"
	"io"
	"io/ioutil"
	"net/http"
	"os"
)

func getSearch(itemName string) gjson.Result {
	api := fmt.Sprintf("https://garlandtools.cn/api/search.php?lang=chs&text=%s", itemName)
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
