package config

import (
	"flag"
	"os"
	"strconv"
	"strings"
	"sync"

	"github.com/cwloo/gonet/logs"
	"github.com/cwloo/gonet/utils"
	"github.com/cwloo/uploader/file_server/global"
	"github.com/cwloo/uploader/file_server/tg_bot"
)

var (
	lock              = &sync.RWMutex{}
	ini    *utils.Ini = &utils.Ini{}
	Config *IniConfig
)

type IniConfig struct {
	Flag                   int
	Log_dir                string
	Log_level              int
	Log_mode               int
	Log_style              int
	Log_timezone           int64
	HttpAddr               string
	UploadPath             string
	GetPath                string
	DelPath                string
	FileinfoPath           string
	UpdateCfgPath          string
	GetCfgPath             string
	CheckMd5               int
	WriteFile              int
	MultiFile              int
	UseAsync               int
	MaxMemory              int64
	MaxSegmentSize         int64
	MaxSingleSize          int64
	MaxTotalSize           int64
	PendingTimeout         int
	FileExpiredTimeout     int
	UploadDir              string
	OssType                string
	Aliyun_BasePath        string
	Aliyun_BucketUrl       string
	Aliyun_BucketName      string
	Aliyun_Endpoint        string
	Aliyun_AccessKeyId     string
	Aliyun_AccessKeySecret string
	Aliyun_Routines        int

	TgBot_ChatId      int64
	TgBot_Token       string
	UseTgBot          int
	Interval          int
	UseOriginFilename int
}

func readIni(filename string) (c *IniConfig) {
	if err := ini.Load(filename); err != nil {
		logs.LogFatal(err.Error())
	}
	c = &IniConfig{}
	c.TgBot_ChatId = ini.GetInt64("tg_bot", "chatId")
	c.TgBot_Token = ini.GetString("tg_bot", "token")
	c.UploadDir = ini.GetString("upload", "dir")
	c.OssType = ini.GetString("upload", "ossType")
	c.Aliyun_BasePath = ini.GetString("aliyun", "basePath")
	c.Aliyun_BucketUrl = ini.GetString("aliyun", "bucketUrl")
	c.Aliyun_BucketName = ini.GetString("aliyun", "bucketName")
	c.Aliyun_Endpoint = ini.GetString("aliyun", "endpoint")
	c.Aliyun_AccessKeyId = ini.GetString("aliyun", "accessKeyId")
	c.Aliyun_AccessKeySecret = ini.GetString("aliyun", "accessKeySecret")
	c.Aliyun_Routines = ini.GetInt("aliyun", "routines")

	c.Flag = ini.GetInt("flag", "flag")
	c.Log_dir = ini.GetString("log", "dir")
	c.Log_level = ini.GetInt("log", "level")
	c.Log_mode = ini.GetInt("log", "mode")
	c.Log_style = ini.GetInt("log", "style")
	c.Log_timezone = ini.GetInt64("log", "timezone")
	c.HttpAddr = ini.GetString("httpserver", "addr")
	c.UploadPath = ini.GetString("path", "upload")
	c.GetPath = ini.GetString("path", "get")
	c.DelPath = ini.GetString("path", "del")
	c.FileinfoPath = ini.GetString("path", "fileinfo")
	c.UpdateCfgPath = ini.GetString("path", "updateconfig")
	c.GetCfgPath = ini.GetString("path", "getconfig")
	c.CheckMd5 = ini.GetInt("upload", "checkMd5")
	c.WriteFile = ini.GetInt("upload", "writeFile")
	c.MultiFile = ini.GetInt("upload", "multiFile")
	c.UseAsync = ini.GetInt("upload", "useAsync")
	c.UseTgBot = ini.GetInt("upload", "useTgBot")
	c.UseOriginFilename = ini.GetInt("upload", "useOriginFilename")
	str := ini.GetString("upload", "maxMemory")
	slice := strings.Split(str, "*")
	val := int64(1)
	for _, v := range slice {
		v = strings.ReplaceAll(v, " ", "")
		c, _ := strconv.ParseInt(v, 10, 0)
		val *= c
	}
	c.MaxMemory = val
	str = ini.GetString("upload", "maxSegmentSize")
	slice = strings.Split(str, "*")
	val = int64(1)
	for _, v := range slice {
		v = strings.ReplaceAll(v, " ", "")
		c, _ := strconv.ParseInt(v, 10, 0)
		val *= c
	}
	c.MaxSegmentSize = val
	str = ini.GetString("upload", "maxSingleSize")
	slice = strings.Split(str, "*")
	val = int64(1)
	for _, v := range slice {
		v = strings.ReplaceAll(v, " ", "")
		c, _ := strconv.ParseInt(v, 10, 0)
		val *= c
	}
	c.MaxSingleSize = val
	str = ini.GetString("upload", "maxTotalSize")
	slice = strings.Split(str, "*")
	val = int64(1)
	for _, v := range slice {
		v = strings.ReplaceAll(v, " ", "")
		c, _ := strconv.ParseInt(v, 10, 0)
		val *= c
	}
	c.MaxTotalSize = val
	str = ini.GetString("upload", "pendingTimeout")
	slice = strings.Split(str, "*")
	val1 := 1
	for _, v := range slice {
		v = strings.ReplaceAll(v, " ", "")
		c, _ := strconv.Atoi(v)
		val1 *= c
	}
	c.PendingTimeout = val1
	str = ini.GetString("upload", "fileExpiredTimeout")
	slice = strings.Split(str, "*")
	val1 = 1
	for _, v := range slice {
		v = strings.ReplaceAll(v, " ", "")
		c, _ := strconv.Atoi(v)
		val1 *= c
	}
	c.FileExpiredTimeout = val1
	str = ini.GetString("flag", "interval")
	slice = strings.Split(str, "*")
	val1 = 1
	for _, v := range slice {
		v = strings.ReplaceAll(v, " ", "")
		c, _ := strconv.Atoi(v)
		val1 *= c
	}
	c.Interval = val1
	return
}

func check() {
	if Config.UploadDir == "" {
		Config.UploadDir = global.Dir_upload
	}
	_, err := os.Stat(Config.UploadDir)
	if err != nil && os.IsNotExist(err) {
		os.MkdirAll(Config.UploadDir, os.ModePerm)
	}
	if Config.Log_dir == "" {
		Config.Log_dir = global.Dir + "logs"
	}
	// 中国大陆这里可能因为被墙了卡住
	tg_bot.NewTgBot(Config.TgBot_Token, Config.TgBot_ChatId, Config.UseTgBot > 0)
}

func InitConfig() {
	Config = readIni("conf.ini")
	if Config == nil {
		logs.LogFatal("error")
	}
	switch Config.Flag {
	case 1:
		flag.Parse()
	default:
	}
	check()
}

func ReadConfig() {
	lock.RLock()
	Config = readIni("conf.ini")
	if Config == nil {
		logs.LogFatal("error")
	}
	check()
	lock.RUnlock()
}

func UpdateConfig(req *global.UpdateCfgReq) {
	lock.Lock()
	if req.Interval != "" {
		ini.SetString("flag", "interval", req.Interval)
	}
	if req.MaxMemory != "" {
		ini.SetString("upload", "maxMemory", req.MaxMemory)
	}
	if req.MaxSegmentSize != "" {
		ini.SetString("upload", "maxSegmentSize", req.MaxSegmentSize)
	}
	if req.MaxSingleSize != "" {
		ini.SetString("upload", "maxSingleSize", req.MaxSingleSize)
	}
	if req.MaxTotalSize != "" {
		ini.SetString("upload", "maxTotalSize", req.MaxTotalSize)
	}
	if req.PendingTimeout != "" {
		ini.SetString("upload", "pendingTimeout", req.PendingTimeout)
	}
	if req.FileExpiredTimeout != "" {
		ini.SetString("upload", "fileExpiredTimeout", req.FileExpiredTimeout)
	}
	if req.CheckMd5 != "" {
		ini.SetString("upload", "checkMd5", req.CheckMd5)
	}
	if req.WriteFile != "" {
		ini.SetString("upload", "writeFile", req.WriteFile)
	}
	ini.SaveTo("conf.ini")
	lock.Unlock()
}

func GetConfig(req *global.GetCfgReq) (*global.GetCfgResp, bool) {
	lock.RLock()
	resp := &global.GetCfgResp{
		ErrCode: 0,
		ErrMsg:  "ok",
		Data: &global.CfgData{
			Log_dir:            Config.Log_dir,
			Log_level:          Config.Log_level,
			Log_mode:           Config.Log_mode,
			Log_style:          Config.Log_style,
			Log_timezone:       Config.Log_timezone,
			HttpAddr:           Config.HttpAddr,
			UploadPath:         Config.UploadPath,
			GetPath:            Config.GetPath,
			DelPath:            Config.DelPath,
			FileinfoPath:       Config.FileinfoPath,
			UpdateCfgPath:      Config.UpdateCfgPath,
			GetCfgPath:         Config.GetCfgPath,
			CheckMd5:           Config.CheckMd5,
			WriteFile:          Config.WriteFile,
			MultiFile:          Config.MultiFile,
			UseAsync:           Config.UseAsync,
			MaxMemory:          Config.MaxMemory,
			MaxSegmentSize:     Config.MaxSegmentSize,
			MaxSingleSize:      Config.MaxSingleSize,
			MaxTotalSize:       Config.MaxTotalSize,
			PendingTimeout:     Config.PendingTimeout,
			FileExpiredTimeout: Config.FileExpiredTimeout,
			UploadDir:          Config.UploadDir,
			OssType:            Config.OssType,
			UseTgBot:           Config.UseTgBot,
			Interval:           Config.Interval,
		},
	}
	lock.RUnlock()
	return resp, true
}
