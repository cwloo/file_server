package main

import (
	"time"

	"github.com/cwloo/gonet/core/base/sys/cmd"
	"github.com/cwloo/gonet/core/base/task"
	"github.com/cwloo/gonet/core/cb"
	"github.com/cwloo/gonet/logs"
	"github.com/cwloo/uploader/src/config"
	"github.com/cwloo/uploader/src/file_server/handler"
	file_server "github.com/cwloo/uploader/src/file_server/server"
	"github.com/cwloo/uploader/src/global"
)

func init() {
	cmd.InitArgs(func(arg *cmd.ARG) {
		arg.CONF.Dir = "config"
		arg.CONF.Name = "conf.ini"
		arg.Append("server", "server", "srv", "svr", "s")
		arg.Append("rpc", "rpc")
	})
}

func main() {
	cmd.ParseArgs()
	config.InitFileConfig(cmd.Conf())
	logs.SetTimezone(logs.TimeZone(config.Config.Log.File.Timezone))
	logs.SetMode(logs.Mode(config.Config.Log.File.Mode))
	logs.SetStyle(logs.Style(config.Config.Log.File.Style))
	logs.Init(config.Config.Log.File.Dir, logs.Level(config.Config.Log.File.Level), global.Exe, 100000000)

	task.After(time.Duration(config.Config.File.Upload.PendingTimeout)*time.Second, cb.NewFunctor00(func() {
		handler.PendingUploader()
	}))

	task.After(time.Duration(config.Config.File.Upload.FileExpiredTimeout)*time.Second, cb.NewFunctor00(func() {
		handler.ExpiredFile()
	}))

	task.After(time.Duration(config.Config.Interval)*time.Second, cb.NewFunctor00(func() {
		handler.ReadConfig()
	}))
	file_server.Run(cmd.Id(), global.Name)
	logs.Close()
}
