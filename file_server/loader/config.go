package main

import (
	"flag"
	"strconv"
	"strings"
	"sync"

	"github.com/cwloo/gonet/logs"
	"github.com/cwloo/gonet/server/gate/global"
	"github.com/cwloo/gonet/server/gate/tg_bot"
	"github.com/cwloo/gonet/utils"
)

var (
	lock              = &sync.RWMutex{}
	ini    *utils.Ini = &utils.Ini{}
	Config *IniConfig
)

type IniConfig struct {
	Flag     int `json:"flag" form:"flag"`
	Interval int `json:"interval" form:"interval"`
	Log      struct {
		Dir      string `json:"dir" form:"dir"`
		Level    int    `json:"level" form:"level"`
		Mode     int    `json:"mode" form:"mode"`
		Style    int    `json:"style" form:"style"`
		Timezone int    `json:"timezone" form:"timezone"`
	} `json:"log" form:"log"`
	Sub struct {
		Num  int    `json:"num" form:"num"`
		Exec string `json:"exec" form:"exec"`
	} `json:"sub" form:"sub"`
	TgBot struct {
		Enable int    `json:"enable" form:"enable"`
		ChatId int64  `json:"chatId" form:"chatId"`
		Token  string `json:"token" form:"token"`
	} `json:"tg_bot" form:"tg_bot"`
	Monitor struct {
		Ip   string `json:"ip" form:"ip"`
		Port []int  `json:"port" form:"port"`
		Path struct {
			Start   string `json:"start" form:"start"`
			Kill    string `json:"kill" form:"kill"`
			KillAll string `json:"killall" form:"killall"`
			SubList string `json:"sublist" form:"sublist"`
		} `json:"path" form:"path"`
	} `json:"monitor" form:"monitor"`
	Etcd struct {
		Schema   string   `json:"schema" form:"schema"`
		Addr     []string `json:"addr" form:"addr"`
		UserName string   `json:"username" form:"username"`
		Password string   `json:"password" form:"password"`
	} `json:"etcd" form:"etcd"`
}

func readIni(filename string) (c *IniConfig) {
	if err := ini.Load(filename); err != nil {
		logs.Fatalf(err.Error())
	}
	c = &IniConfig{}
	// Flag
	c.Flag = ini.GetInt("flag", "flag")
	s := ini.GetString("flag", "interval")
	sli := strings.Split(s, "*")
	va := 1
	for _, v := range sli {
		v = strings.ReplaceAll(v, " ", "")
		c, _ := strconv.Atoi(v)
		va *= c
	}
	c.Interval = va
	// Log
	c.Log.Dir = ini.GetString("log", "dir")
	c.Log.Level = ini.GetInt("log", "level")
	c.Log.Mode = ini.GetInt("log", "mode")
	c.Log.Style = ini.GetInt("log", "style")
	c.Log.Timezone = ini.GetInt("log", "timezone")
	// Sub
	c.Sub.Num = ini.GetInt("sub", "num")
	c.Sub.Exec = ini.GetString("sub", "execname")
	// TgBot
	c.TgBot.Enable = ini.GetInt("tg_bot", "enable")
	c.TgBot.ChatId = ini.GetInt64("tg_bot", "chatId")
	c.TgBot.Token = ini.GetString("tg_bot", "token")
	// Monitor
	c.Monitor.Ip = ini.GetString("monitor", "ip")
	ports := strings.Split(ini.GetString("monitor", "port"), ",")
	for _, port := range ports {
		switch port == "" {
		case false:
			c.Monitor.Port = append(c.Monitor.Port, utils.Atoi(port))
		}
	}
	// Path
	c.Monitor.Path.Start = ini.GetString("path", "start")
	c.Monitor.Path.Kill = ini.GetString("path", "kill")
	c.Monitor.Path.KillAll = ini.GetString("path", "killall")
	c.Monitor.Path.SubList = ini.GetString("path", "sublist")
	// Etcd
	c.Etcd.Schema = ini.GetString("etcd", "schema")
	addrs := strings.Split(ini.GetString("etcd", "addr"), ",")
	for _, addr := range addrs {
		switch addr == "" {
		case false:
			c.Etcd.Addr = append(c.Etcd.Addr, addr)
		}
	}
	c.Etcd.UserName = ini.GetString("etcd", "username")
	c.Etcd.Password = ini.GetString("etcd", "password")
	return
}

func check() {
	if Config.Log.Dir == "" {
		Config.Log.Dir = global.Dir + "logs"
	}
	if Config.Log.Timezone != int(logs.GetTimeZone()) {
		logs.SetTimezone(logs.TimeZone(Config.Log.Timezone))
	}
	if Config.Log.Mode != int(logs.GetMode()) {
		logs.SetMode(logs.Mode(Config.Log.Mode))
	}
	if Config.Log.Style != int(logs.GetStyle()) {
		logs.SetStyle(logs.Style(Config.Log.Style))
	}
	if Config.Log.Level != int(logs.GetLevel()) {
		logs.SetLevel(logs.Level(Config.Log.Level))
	}
	// 中国大陆这里可能因为被墙了卡住
	tg_bot.NewTgBot(Config.TgBot.Token, Config.TgBot.ChatId, Config.TgBot.Enable > 0)
}

func read() {
	Config = readIni("../config/conf.ini")
	if Config == nil {
		logs.Fatalf("error")
	}
}

func InitConfig() {
	read()
	switch Config.Flag {
	case 1:
		flag.Parse()
	default:
	}
	check()
}

func readConfig() {
	read()
	check()
}

func ReadConfig() {
	lock.RLock()
	readConfig()
	lock.RUnlock()
}

func updateConfig(req *global.UpdateCfgReq) {
	if req.Interval != "" {
		ini.SetString("flag", "interval", req.Interval)
	}
	if req.LogTimezone != "" {
		v, _ := strconv.Atoi(req.LogTimezone)
		ini.SetInt("log", "timezone", v)
	}
	if req.LogMode != "" {
		v, _ := strconv.Atoi(req.LogMode)
		ini.SetInt("log", "mode", v)
	}
	if req.LogStyle != "" {
		v, _ := strconv.Atoi(req.LogStyle)
		ini.SetInt("log", "style", v)
	}
	if req.LogLevel != "" {
		v, _ := strconv.Atoi(req.LogLevel)
		ini.SetInt("log", "level", v)
	}
	if req.UseTgBot != "" {
		v, _ := strconv.Atoi(req.UseTgBot)
		ini.SetInt("tg_bot", "enable", v)
	}
	if req.TgBotChatId != "" {
		v, _ := strconv.ParseInt(req.TgBotChatId, 10, 0)
		ini.SetInt64("tg_bot", "chatId", v)
	}
	if req.TgBotToken != "" {
		ini.SetString("tg_bot", "token", req.TgBotToken)
	}
	ini.SaveTo("../config/conf.ini")
}

func UpdateConfig(req *global.UpdateCfgReq) {
	lock.Lock()
	updateConfig(req)
	readConfig()
	lock.Unlock()
}

func GetConfig(req *global.GetCfgReq) (*global.GetCfgResp, bool) {
	lock.RLock()
	resp := &global.GetCfgResp{
		ErrCode: 0,
		ErrMsg:  "ok",
		Data: &global.CfgData{
			Interval:     Config.Interval,
			Log_dir:      Config.Log.Dir,
			Log_level:    Config.Log.Level,
			Log_mode:     Config.Log.Mode,
			Log_style:    Config.Log.Style,
			Log_timezone: Config.Log.Timezone,
			UseTgBot:     Config.TgBot.Enable,
			TgBotChatId:  Config.TgBot.ChatId,
			TgBotToken:   Config.TgBot.Token,
		},
	}
	lock.RUnlock()
	return resp, true
}
