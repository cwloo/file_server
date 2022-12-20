package main

import (
	"github.com/cwloo/gonet/logs"
	"github.com/cwloo/uploader/file_server/config"
	"github.com/cwloo/uploader/file_server/gate"
	"github.com/cwloo/uploader/file_server/global"
)

func main() {
	config.InitConfig()
	logs.SetTimezone(logs.TimeZone(config.Config.Log.Timezone))
	logs.SetMode(logs.Mode(config.Config.Log.Mode))
	logs.SetStyle(logs.Style(config.Config.Log.Style))
	logs.Init(config.Config.Log.Dir, logs.Level(config.Config.Log.Level), global.Exe, 100000000)
	gate.Run()
	logs.Close()
}
