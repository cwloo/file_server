package config

import (
	"flag"
	"strconv"
	"strings"

	"github.com/cwloo/gonet/logs"
	"github.com/cwloo/gonet/utils"
)

var Config *IniConfig

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
	UploadlDir             string
	OssType                string
	Aliyun_BasePath        string
	Aliyun_BucketUrl       string
	Aliyun_BucketName      string
	Aliyun_Endpoint        string
	Aliyun_AccessKeyId     string
	Aliyun_AccessKeySecret string

	TgBot_ChatId int64
	TgBot_Token  string
}

func readIni(filename string) (c *IniConfig) {
	ini := utils.Ini{}
	if err := ini.Load(filename); err != nil {
		logs.LogFatal(err.Error())
	}
	c = &IniConfig{}
	c.TgBot_ChatId = ini.GetInt64("tg_bot", "chatId")
	c.TgBot_Token = ini.GetString("tg_bot", "token")
	c.UploadlDir = ini.GetString("upload", "dir")
	c.OssType = ini.GetString("upload", "ossType")
	c.Aliyun_BasePath = ini.GetString("aliyun", "basePath")
	c.Aliyun_BucketUrl = ini.GetString("aliyun", "bucketUrl")
	c.Aliyun_BucketName = ini.GetString("aliyun", "bucketName")
	c.Aliyun_Endpoint = ini.GetString("aliyun", "endpoint")
	c.Aliyun_AccessKeyId = ini.GetString("aliyun", "accessKeyId")
	c.Aliyun_AccessKeySecret = ini.GetString("aliyun", "accessKeySecret")

	c.Flag = ini.GetInt("flag", "flag")
	c.Log_dir = ini.GetString("log", "dir")
	c.Log_level = ini.GetInt("log", "level")
	c.Log_mode = ini.GetInt("log", "mode")
	c.Log_style = ini.GetInt("log", "style")
	c.Log_timezone = ini.GetInt64("log", "timezone")
	c.HttpAddr = ini.GetString("httpserver", "addr")
	c.UploadPath = ini.GetString("path", "upload")
	c.GetPath = ini.GetString("path", "get")
	c.CheckMd5 = ini.GetInt("upload", "checkMd5")
	c.WriteFile = ini.GetInt("upload", "writeFile")
	c.MultiFile = ini.GetInt("upload", "multiFile")
	c.UseAsync = ini.GetInt("upload", "useAsync")
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
	return
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
}
