package main

import (
	"github.com/cwloo/gonet/logs"
	"github.com/cwloo/gonet/utils"
	"github.com/cwloo/uploader/src/config"
	"github.com/cwloo/uploader/src/global"
	"github.com/cwloo/uploader/src/loader/handler"
	loader "github.com/cwloo/uploader/src/loader/server"
	"github.com/cwloo/uploader/src/loader/sub"
)

func main() {
	global.Name = "monitor"
	global.Cmd.ID, global.Cmd.Dir, global.Cmd.Conf_Dir, global.Cmd.Log_Dir = global.ParseArgs()
	config.InitConfig(global.Name, global.Cmd.Conf_Dir)
	logs.Errorf("%v log_dir=%v", global.Name, config.Config.Log.Monitor.Dir)
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
	loader.Run(global.Cmd.ID)
	sub.WaitAll()
	logs.Debugf("exit...")
	logs.Close()
}
