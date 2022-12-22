package main

import (
	"github.com/cwloo/gonet/logs"
	"github.com/cwloo/gonet/utils"
	"github.com/cwloo/uploader/src/config"
	"github.com/cwloo/uploader/src/global"
	"github.com/cwloo/uploader/src/loader/handler"
	"github.com/cwloo/uploader/src/loader/handler/sub"
	loader "github.com/cwloo/uploader/src/loader/server"
)

func main() {
	global.Cmd.ParseArgs()
	config.InitMonitorConfig(global.Cmd.Conf_Dir)
	logs.SetTimezone(logs.TimeZone(config.Config.Log.Monitor.Timezone))
	logs.SetMode(logs.Mode(config.Config.Log.Monitor.Mode))
	logs.SetStyle(logs.Style(config.Config.Log.Monitor.Style))
	logs.Init(config.Config.Log.Monitor.Dir, logs.Level(config.Config.Log.Monitor.Level), global.Exe, 100000000)
	go func() {
		utils.ReadConsole(handler.OnInput)
	}()
	// dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	// if err != nil {
	// 	logs.Fatalf("%v", err)
	// }
	sub.Start()
	loader.Run(global.Cmd.ID, global.Name)
	sub.WaitAll()
	logs.Debugf("exit...")
	logs.Close()
}
