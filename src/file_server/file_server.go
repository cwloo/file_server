package main

import (
	"time"

	"github.com/cwloo/gonet/core/base/task"
	"github.com/cwloo/gonet/core/cb"
	"github.com/cwloo/gonet/logs"
	"github.com/cwloo/uploader/src/config"
	"github.com/cwloo/uploader/src/file_server/handler"
	file_server "github.com/cwloo/uploader/src/file_server/server"
	"github.com/cwloo/uploader/src/global"
)

func main() {
	global.Name = "file"
	global.Cmd.ID, global.Cmd.Dir, global.Cmd.Conf_Dir, global.Cmd.Log_Dir = global.ParseArgs()
	config.InitConfig(global.Name, global.Cmd.Conf_Dir)
	logs.Errorf("%v log_dir=%v", global.Name, config.Config.Log.File.Dir)
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
	file_server.Run(global.Cmd.ID)
	logs.Close()
}
