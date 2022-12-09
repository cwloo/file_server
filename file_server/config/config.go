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
	Flag               int
	Log_dir            string
	Log_level          int
	Log_mode           int
	Log_style          int
	Log_timezone       int64
	HttpAddr           string
	UploadPath         string
	GetPath            string
	MultiFile          int
	UseAsync           int
	MaxMemory          int64
	MaxSegmentSize     int64
	MaxSingleSize      int64
	MaxTotalSize       int64
	PendingTimeout     int
	FileExpiredTimeout int

	UploadLocalDir            string
	OssType                   string
	AliyunOSS_BasePath        string
	AliyunOSS_BucketUrl       string
	AliyunOSS_BucketName      string
	AliyunOSS_Endpoint        string
	AliyunOSS_AccessKeyId     string
	AliyunOSS_AccessKeySecret string
}

func readIni(filename string) (c *IniConfig) {
	ini := utils.Ini{}
	if err := ini.Load(filename); err != nil {
		logs.LogFatal(err.Error())
	}
	c = &IniConfig{}

	c.OssType = ini.GetString("upload", "ossType")
	c.AliyunOSS_BasePath = ini.GetString("aliyunOSS", "basePath")
	c.AliyunOSS_BucketUrl = ini.GetString("aliyunOSS", "bucketUrl")
	c.AliyunOSS_BucketName = ini.GetString("aliyunOSS", "bucketName")
	c.AliyunOSS_Endpoint = ini.GetString("aliyunOSS", "endpoint")
	c.AliyunOSS_AccessKeyId = ini.GetString("aliyunOSS", "accessKeyId")
	c.AliyunOSS_AccessKeySecret = ini.GetString("aliyunOSS", "accessKeySecret")

	c.Flag = ini.GetInt("flag", "flag")
	c.Log_dir = ini.GetString("log", "dir")
	c.Log_level = ini.GetInt("log", "level")
	c.Log_mode = ini.GetInt("log", "mode")
	c.Log_style = ini.GetInt("log", "style")
	c.Log_timezone = ini.GetInt64("log", "timezone")
	c.HttpAddr = ini.GetString("httpserver", "addr")
	c.UploadPath = ini.GetString("path", "upload")
	c.GetPath = ini.GetString("path", "get")
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
