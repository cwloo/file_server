package main

import (
	"flag"
	"fmt"

	"github.com/cwloo/gonet/utils"
)

var Config *IniConfig

type IniConfig struct {
	Flag         int
	Sub          int
	Exec         string
	FileList     []string
	Log_dir      string
	Log_level    int
	Log_mode     int
	Log_style    int
	Log_timezone int64
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
	c.Log_style = ini.GetInt("log", "style")
	c.Log_timezone = ini.GetInt64("log", "timezone")
	c.Sub = ini.GetInt("sub", "num")
	c.Exec = ini.GetString("sub", "execname")
	num := ini.GetInt("file", "num")
	for i := 0; i < num; i++ {
		c.FileList = append(c.FileList, ini.GetString("file", fmt.Sprintf("file%v", i)))
	}
	return
}

func InitConfig() {
	Config = readIni("conf.ini")
	if Config == nil {
		panic(utils.Stack())
	}
	switch Config.Flag {
	case 1:
		//解析命令行解析
		//.\loader -sub=5 -c=2 -file0= -file1=
		Config.Sub = *flag.Int("sub", 1, "")
		num := *flag.Int("c", 0, "")
		for i := 0; i < num; i++ {
			Config.FileList = append(Config.FileList, *flag.String(fmt.Sprintf("file%v", i), "", ""))
		}
		flag.Parse()
	default:
	}
}
