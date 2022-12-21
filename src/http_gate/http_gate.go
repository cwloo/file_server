package main

import (
	"github.com/cwloo/gonet/logs"
	"github.com/cwloo/uploader/src/global"
	"github.com/cwloo/uploader/src/http_gate/config"
	"github.com/cwloo/uploader/src/http_gate/gate"
)

func main() {
	id, _, conf := global.ParseArgs()
	config.InitConfig(conf)
	logs.SetTimezone(logs.TimeZone(config.Config.Log.Timezone))
	logs.SetMode(logs.Mode(config.Config.Log.Mode))
	logs.SetStyle(logs.Style(config.Config.Log.Style))
	logs.Init(config.Config.Log.Dir, logs.Level(config.Config.Log.Level), global.Exe, 100000000)
	gate.Run(id)
	logs.Close()
}
