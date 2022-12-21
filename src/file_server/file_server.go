package main

import (
	"time"

	"github.com/cwloo/gonet/core/base/task"
	"github.com/cwloo/gonet/core/cb"
	"github.com/cwloo/gonet/logs"
	"github.com/cwloo/uploader/src/file_server/config"
	file_server "github.com/cwloo/uploader/src/file_server/server"
	"github.com/cwloo/uploader/src/global"
	"github.com/cwloo/uploader/src/global/handler"
)

func main() {
	id, _, conf := global.ParseArgs()
	config.InitConfig(conf)
	logs.SetTimezone(logs.TimeZone(config.Config.Log.Timezone))
	logs.SetMode(logs.Mode(config.Config.Log.Mode))
	logs.SetStyle(logs.Style(config.Config.Log.Style))
	logs.Init(config.Config.Log.Dir, logs.Level(config.Config.Log.Level), global.Exe, 100000000)

	task.After(time.Duration(config.Config.File.Upload.PendingTimeout)*time.Second, cb.NewFunctor00(func() {
		handler.PendingUploader()
	}))

	task.After(time.Duration(config.Config.File.Upload.FileExpiredTimeout)*time.Second, cb.NewFunctor00(func() {
		handler.ExpiredFile()
	}))

	task.After(time.Duration(config.Config.Interval)*time.Second, cb.NewFunctor00(func() {
		handler.ReadConfig()
	}))
	file_server.Run(id)
	logs.Close()
}
