package main

import "github.com/cwloo/gonet/logs"

// cd loader
// ./loader 必须父进程启动
func main() {
	InitConfig()
	// logs.LogTimezone(logs.MY_CST)
	// logs.LogMode(logs.M_STDOUT_FILE)
	// logs.LogStyle(logs.F_DETAIL)
	// logs.LogInit(dir+"logs", logs.LVL_DEBUG, exe, 100000000)
	logs.LogTimezone(logs.TimeZone(Config.Log_timezone))
	logs.LogMode(logs.Mode(Config.Log_mode))
	logs.LogStyle(logs.Style(Config.Log_style))
	logs.LogInit(dir+"logs", int32(Config.Log_level), exe, 100000000)

	switch MultiFile {
	case true:
		multiUpload()
	default:
		upload()
	}
}
