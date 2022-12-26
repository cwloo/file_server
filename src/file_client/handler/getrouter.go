package handler

import (
	"encoding/json"
	"io/ioutil"
	"math/rand"
	"net/http"
	"strconv"
	"strings"

	"github.com/cwloo/gonet/logs"
	"github.com/cwloo/uploader/src/config"
	"github.com/cwloo/uploader/src/global"
)

func GetRouter(client *http.Client, md5 string) string {
	method := "GET"
	r := rand.Int() % len(config.Config.Client.Addr)
	url := strings.Join([]string{
		config.Config.Client.Addr[r].Proto, "://",
		config.Config.Client.Addr[r].Ip, ":",
		strconv.Itoa(config.Config.Client.Addr[r].Port),
		config.Config.Client.Path.Router}, "")
	logs.Tracef("request =>> %v %v", method, url+"?md5="+md5)
	res, err := client.Get(url + "?md5=" + md5)
	if err != nil {
		logs.Errorf(err.Error())
		return ""
	}
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		logs.Errorf(err.Error())
		return ""
	}
	if len(body) == 0 {
		return ""
	}
	resp := global.RouterResp{}
	err = json.Unmarshal(body, &resp)
	if err != nil {
		logs.Errorf(err.Error())
		logs.Warnf("%v", string(body))
		return ""
	}
	switch resp.ErrCode {
	case 0:
		if resp.Node.Domain == "" {
			logs.Fatalf("error")
		}
		logs.Tracef("%v", string(body))
		return resp.Node.Domain
	}
	return ""
}
