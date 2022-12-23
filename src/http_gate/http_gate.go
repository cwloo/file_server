package main

import (
	"github.com/cwloo/gonet/logs"
	"github.com/cwloo/uploader/src/config"
	"github.com/cwloo/uploader/src/global"
	"github.com/cwloo/uploader/src/global/cmd"
	http_gate "github.com/cwloo/uploader/src/http_gate/server"
)

func main() {
	cmd.ParseArgs()
	config.InitHttpGateConfig(cmd.Conf())
	logs.SetTimezone(logs.TimeZone(config.Config.Log.Gate.Http.Timezone))
	logs.SetMode(logs.Mode(config.Config.Log.Gate.Http.Mode))
	logs.SetStyle(logs.Style(config.Config.Log.Gate.Http.Style))
	logs.Init(config.Config.Log.Gate.Http.Dir, logs.Level(config.Config.Log.Gate.Http.Level), global.Exe, 100000000)
	http_gate.Run(cmd.Id(), global.Name)
	logs.Close()
}
