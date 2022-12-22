package main

import (
	"flag"
	"strconv"
	"strings"

	"github.com/cwloo/gonet/logs"
	"github.com/cwloo/gonet/utils"
)

var Config *IniConfig

type IniConfig struct {
	Flag         int
	Log_dir      string
	Log_level    int
	Log_mode     int
	Log_style    int
	Log_timezone int64
	HttpProto    string
	HttpAddr     string
	UploadPath   string
	FileinfoPath string
	GetPath      string
	SegmentSize  int64
	MultiFile    int

	TgBot_ChatId int64
	TgBot_Token  string
}

func readIni(filename string) (c *IniConfig) {
	ini := utils.Ini{}
	if err := ini.Load(filename); err != nil {
		logs.Fatalf(err.Error())
	}
	c = &IniConfig{}
	c.TgBot_ChatId = ini.GetInt64("tg_bot", "chatId")
	c.TgBot_Token = ini.GetString("tg_bot", "token")
	c.Flag = ini.GetInt("flag", "flag")
	c.Log_dir = ini.GetString("log", "dir")
	c.Log_level = ini.GetInt("log", "level")
	c.Log_mode = ini.GetInt("log", "mode")
	c.Log_style = ini.GetInt("log", "style")
	c.Log_timezone = ini.GetInt64("log", "timezone")
	c.HttpProto = ini.GetString("httpserver", "proto")
	c.HttpAddr = ini.GetString("httpserver", "addr")
	c.UploadPath = ini.GetString("path", "upload")
	c.GetPath = ini.GetString("path", "get")
	c.FileinfoPath = ini.GetString("path", "fileinfo")
	c.MultiFile = ini.GetInt("upload", "multiFile")
	str := ini.GetString("upload", "segmentSize")
	slice := strings.Split(str, "*")
	val := int64(1)
	for _, v := range slice {
		v = strings.ReplaceAll(v, " ", "")
		c, _ := strconv.ParseInt(v, 10, 0)
		val *= c
	}
	c.SegmentSize = val
	return
}

func InitConfig() {
	Config = readIni("conf.ini")
	if Config == nil {
		logs.Fatalf("error")
	}
	switch Config.Flag {
	case 1:
		flag.Parse()
	default:
	}
	Init()
}