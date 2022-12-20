package main

import "github.com/cwloo/gonet/logs"

// cd loader
// ./loader 必须父进程启动
func main() {
	InitConfig()
	// logs.SetTimezone(logs.MY_CST)
	// logs.SetMode(logs.M_STDOUT_FILE)
	// logs.SetStyle(logs.F_DETAIL)
	// logs.Init(dir+"logs", logs.LVL_DEBUG, exe, 100000000)
	logs.SetTimezone(logs.TimeZone(Config.Log_timezone))
	logs.SetMode(logs.Mode(Config.Log_mode))
	logs.SetStyle(logs.Style(Config.Log_style))
	logs.Init(dir+"logs", logs.Level(Config.Log_level), exe, 100000000)

	switch MultiFile {
	case true:
		multiUpload()
	default:
		upload()
	}
}
