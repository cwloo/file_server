package main

import (
	"github.com/cwloo/gonet/core/base/sys/cmd"
	"github.com/cwloo/gonet/logs"
	"github.com/cwloo/uploader/src/config"
	"github.com/cwloo/uploader/src/global"
)

func init() {
	cmd.InitArgs(func(arg *cmd.ARG) {
		arg.SetConf("config/conf.ini")
	})
}

func main() {
	cmd.ParseArgs()
	config.InitClientConfig(cmd.Conf())
	logs.SetTimezone(logs.TimeZone(config.Config.Log.Client.Timezone))
	logs.SetMode(logs.Mode(config.Config.Log.Client.Mode))
	logs.SetStyle(logs.Style(config.Config.Log.Client.Style))
	logs.Init(config.Config.Log.Client.Dir, logs.Level(config.Config.Log.Client.Level), global.Exe, 100000000)
	switch config.Config.Client.Upload.MultiFile > 0 {
	case true:
		multiUpload()
	default:
		upload()
	}
}
