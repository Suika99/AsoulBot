package minecraft

import (
	"encoding/json"
	"fmt"
	log "github.com/sirupsen/logrus"
	"net/http"
)

func getServerInfo(addr *Mcserver) (result *ServerInfo) {
	url := fmt.Sprintf("https://mcapi.us/server/status?ip=%s&port=%s", addr.IP, addr.Port)
	resp, err := http.Get(url)
	if err != nil {
		log.Errorln(err)
	}
	defer resp.Body.Close()
	result = &ServerInfo{}
	if err := json.NewDecoder(resp.Body).Decode(result); err != nil {
		log.Errorln(err)
	}
	return result
}
