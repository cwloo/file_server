package main

import (
	"flag"
	"fmt"

	"github.com/cwloo/gonet/utils"
)

var Config *IniConfig

type IniConfig struct {
	Flag     int
	Sub      int
	FileList []string
}

func readini(filename string) (c *IniConfig) {
	ini := utils.Ini{}
	if err := ini.Load("conf.ini"); err != nil {
		fmt.Printf("load %s err: [%s]\n", filename, err.Error())
		return
	}
	c = &IniConfig{}
	c.Flag = ini.GetInt("flag", "flag")
	c.Sub = ini.GetInt("sub", "num")
	num := ini.GetInt("file", "num")
	for i := 0; i < num; i++ {
		c.FileList = append(c.FileList, ini.GetString("file", fmt.Sprintf("file%v", i)))
	}
	return
}

func InitConfig() {
	Config = readini("conf.ini")
	if Config == nil {
		panic(utils.Stack())
	}
	switch Config.Flag {
	case 1:
		//解析命令行解析
		flag.Parse()
		//.\loader -sub= -c=2 -file0= -file1=
		Config.Sub = *flag.Int("sub", 1, "")
		num := *flag.Int("c", 0, "")
		for i := 0; i < num; i++ {
			Config.FileList = append(Config.FileList, *flag.String(fmt.Sprintf("file%v", i), "", ""))
		}
	default:
	}
}
