package main

import (
	"flag"

	"github.com/cwloo/gonet/utils"
)

var Config *IniConfig

type IniConfig struct {
	Flag         int
	Log_dir      string
	Log_level    int
	Log_mode     int
	Log_timezone int64
	HttpAddr     string
}

func readIni(filename string) (c *IniConfig) {
	ini := utils.Ini{}
	if err := ini.Load(filename); err != nil {
		panic(err.Error())
	}
	c = &IniConfig{}
	c.Flag = ini.GetInt("flag", "flag")
	c.Log_dir = ini.GetString("log", "dir")
	c.Log_level = ini.GetInt("log", "level")
	c.Log_mode = ini.GetInt("log", "mode")
	c.Log_timezone = ini.GetInt64("log", "timezone")
	c.HttpAddr = ini.GetString("httpserver", "addr")
	return
}

func InitConfig() {
	Config = readIni("loader/conf.ini")
	if Config == nil {
		panic(utils.Stack())
	}
	switch Config.Flag {
	case 1:
		flag.Parse()
	default:
	}
}
