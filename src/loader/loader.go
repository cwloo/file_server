package main

import (
	"time"

	"github.com/cwloo/gonet/core/base/sys/cmd"
	"github.com/cwloo/gonet/core/base/task"
	"github.com/cwloo/gonet/core/cb"
	"github.com/cwloo/gonet/logs"
	"github.com/cwloo/gonet/utils"
	"github.com/cwloo/uploader/src/config"
	"github.com/cwloo/uploader/src/global"
	"github.com/cwloo/uploader/src/loader/handler"
	"github.com/cwloo/uploader/src/loader/handler/sub"
	loader "github.com/cwloo/uploader/src/loader/server"
)

func init() {
	cmd.InitArgs(func(arg *cmd.ARG) {
		arg.SetConf("config/conf.ini")
		arg.AppendPattern("server", "server", "srv", "svr", "s")
		arg.AppendPattern("rpc", "rpc", "r")
	})
}

func main() {
	cmd.ParseArgs()
	config.InitMonitorConfig(cmd.Conf())
	logs.SetTimezone(logs.Timezone(config.Config.Log.Monitor.Timezone))
	logs.SetMode(logs.Mode(config.Config.Log.Monitor.Mode))
	logs.SetStyle(logs.Style(config.Config.Log.Monitor.Style))
	logs.SetLevel(logs.Level(config.Config.Log.Monitor.Level))
	logs.Init(config.Config.Log.Monitor.Dir, global.Exe, 100000000)
	go func() {
		utils.ReadConsole(handler.OnInput)
	}()
	// dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	// if err != nil {
	// 	logs.Fatalf(err.Error())
	// }

	task.After(time.Duration(config.Config.Interval)*time.Second, cb.NewFunctor00(func() {
		handler.ReadConfig()
	}))
	loader.Run(cmd.Id(), config.ServiceName())
	sub.Start()
	sub.WaitAll()
	logs.Debugf("exit...")
	logs.Close()
}
