package config

import (
	"flag"
	"os"
	"strconv"
	"strings"
	"sync"

	"github.com/cwloo/gonet/logs"
	"github.com/cwloo/gonet/utils"

	"github.com/cwloo/uploader/src/global"
	"github.com/cwloo/uploader/src/global/tg_bot"
)

var (
	lock              = &sync.RWMutex{}
	ini    *utils.Ini = &utils.Ini{}
	Config *IniConfig
)

type IniConfig struct {
	Flag     int `json:"flag" form:"flag"`
	Interval int `json:"interval" form:"interval"`
	Path     struct {
		UpdateCfg string `json:"updatecfg" form:"updatecfg"`
		GetCfg    string `json:"getcfg" form:"getcfg"`
	} `json:"path" form:"path"`
	Log struct {
		Dir      string `json:"dir" form:"dir"`
		Level    int    `json:"level" form:"level"`
		Mode     int    `json:"mode" form:"mode"`
		Style    int    `json:"style" form:"style"`
		Timezone int    `json:"timezone" form:"timezone"`
	} `json:"log" form:"log"`
	Sub struct {
		Gate struct {
			Num  int    `json:"num" form:"num"`
			Dir  string `json:"dir" form:"dir"`
			Exec string `json:"exec" form:"exec"`
			Http struct {
				Num  int    `json:"num" form:"num"`
				Dir  string `json:"dir" form:"dir"`
				Exec string `json:"exec" form:"exec"`
			} `json:"http" form:"http"`
		} `json:"gate" form:"gate"`
		File struct {
			Num  int    `json:"num" form:"num"`
			Dir  string `json:"dir" form:"dir"`
			Exec string `json:"exec" form:"exec"`
		} `json:"file" form:"file"`
	} `json:"sub" form:"sub"`
	TgBot struct {
		Enable int    `json:"enable" form:"enable"`
		ChatId int64  `json:"chatId" form:"chatId"`
		Token  string `json:"token" form:"token"`
	} `json:"tg_bot" form:"tg_bot"`
	File struct {
		Ip          string `json:"ip" form:"ip"`
		Port        []int  `json:"port" form:"port"`
		MaxConn     int    `json:"maxConn" form:"maxConn"`
		IdleTimeout int    `json:"idleTimeout" form:"idleTimeout"`
		Upload      struct {
			Dir                string `json:"dir" form:"dir"`
			CheckMd5           int    `json:"checkMd5" form:"checkMd5"`
			WriteFile          int    `json:"writeFile" form:"writeFile"`
			MultiFile          int    `json:"multiFile" form:"multiFile"`
			UseAsync           int    `json:"useAsync" form:"useAsync"`
			MaxMemory          int64  `json:"maxMemory" form:"maxMemory"`
			MaxSegmentSize     int64  `json:"maxSegmentSize" form:"maxSegmentSize"`
			MaxSingleSize      int64  `json:"maxSingleSize" form:"maxSingleSize"`
			MaxTotalSize       int64  `json:"maxTotalSize" form:"maxTotalSize"`
			PendingTimeout     int    `json:"pendingTimeout" form:"pendingTimeout"`
			FileExpiredTimeout int    `json:"fileExpiredTimeout" form:"fileExpiredTimeout"`
			UseOriginFilename  int    `json:"useOriginFilename" form:"useOriginFilename"`
		} `json:"upload" form:"upload"`
		Path struct {
			Upload     string `json:"upload" form:"upload"`
			Get        string `json:"get" form:"get"`
			Del        string `json:"del" form:"del"`
			Fileinfo   string `json:"fileinfo" form:"fileinfo"`
			FileDetail string `json:"filedetail" form:"filedetail"`
			UuidList   string `json:"uuidlist" form:"uuidlist"`
			List       string `json:"list" form:"list"`
		} `json:"path" form:"path"`
	} `json:"file" form:"file"`
	Oss struct {
		Type   string `json:"type" form:"type"`
		Aliyun struct {
			BasePath        string `json:"basepath" form:"basepath"`
			BucketUrl       string `json:"bucketUrl" form:"bucketUrl"`
			BucketName      string `json:"bucketName" form:"bucketName"`
			EndPoint        string `json:"endpoint" form:"endpoint"`
			AccessKeyId     string `json:"accessKeyId" form:"accessKeyId"`
			AccessKeySecret string `json:"accessKeySecret" form:"accessKeySecret"`
			Routines        int    `json:"routines" form:"routines"`
		} `json:"aliyun" form:"aliyun"`
		Aws_s3 struct {
			Bucket           string `json:"bucket" form:"bucket"`
			Region           string `json:"region" form:"region"`
			EndPoint         string `json:"endpoint" form:"endpoint"`
			Force_path_style int    `json:"force_path_style" form:"force_path_style"`
			Disable_ssl      int    `json:"disable_ssl" form:"disable_ssl"`
			Secret_id        string `json:"secret_id" form:"secret_id"`
			Secret_key       string `json:"secret_key" form:"secret_key"`
			Base_url         string `json:"base_url" form:"base_url"`
			Path_prefix      string `json:"path_prefix" form:"path_prefix"`
		} `json:"aws_s3" form:"aws_s3"`
		Tencent_cos struct {
			Bucket      string `json:"bucket" form:"bucket"`
			Region      string `json:"region" form:"region"`
			Secret_id   string `json:"secret_id" form:"secret_id"`
			Secret_key  string `json:"secret_key" form:"secret_key"`
			Base_url    string `json:"base_url" form:"base_url"`
			Path_prefix string `json:"path_prefix" form:"path_prefix"`
		} `json:"tencent_cos" form:"tencent_cos"`
		Qiniu struct {
			Zone            string `json:"zone" form:"zone"`
			Bucket          string `json:"bucket" form:"bucket"`
			ImgPath         string `json:"imgPath" form:"imgPath"`
			UseHttps        string `json:"useHttps" form:"useHttps"`
			Access_key      string `json:"access_key" form:"access_key"`
			Secret_key      string `json:"secret_key" form:"secret_key"`
			Base_url        string `json:"base_url" form:"base_url"`
			Use_cdn_domains string `json:"use-cdn-domains" form:"use_cdn_domains"`
		} `json:"qniu" form:"qniu"`
		Huawei_obs struct {
			Path       string `json:"path" form:"path"`
			Bucket     string `json:"bucket" form:"bucket"`
			EndPoint   string `json:"endpoint" form:"endpoint"`
			Access_key string `json:"access_key" form:"access_key"`
			Secret_key string `json:"secret_key" form:"secret_key"`
			Base_url   string `json:"base_url" form:"base_url"`
		} `json:"huawei_obs" form:"huawei_obs"`
	} `json:"oss" form:"oss"`
	Etcd struct {
		Schema   string   `json:"schema" form:"schema"`
		Addr     []string `json:"addr" form:"addr"`
		UserName string   `json:"username" form:"username"`
		Password string   `json:"password" form:"password"`
	} `json:"etcd" form:"etcd"`
	Gate struct {
		Proto string `json:"proto" form:"proto"`
		Ip    string `json:"ip" form:"ip"`
		Port  []int  `json:"port" form:"port"`
		Path  struct {
			Handshake string `json:"handshake" form:"handshake"`
		} `json:"path" form:"path"`
		MaxConn          int `json:"maxConn" form:"maxConn"`
		UsePool          int `json:"usePool" form:"usePool"`
		HandshakeTimeout int `json:"handshakeTimeout" form:"handshakeTimeout"`
		IdleTimeout      int `json:"idleTimeout" form:"idleTimeout"`
		ReadBufferSize   int `json:"readBufferSize" form:"readBufferSize"`
		PrintInterval    int `json:"printInterval" form:"printInterval"`
		Http             struct {
			Ip          string `json:"ip" form:"ip"`
			Port        []int  `json:"port" form:"port"`
			MaxConn     int    `json:"maxConn" form:"maxConn"`
			IdleTimeout int    `json:"idleTimeout" form:"idleTimeout"`
			Path        struct {
				Fileserver string `json:"fileserver" form:"fileserver"`
			} `json:"path" form:"path"`
		} `json:"http" form:"http"`
	} `json:"gate" form:"gate"`
	Rpc struct {
		Ip   string `json:"ip" form:"ip"`
		Gate struct {
			Port []int  `json:"port" form:"port"`
			Node string `json:"node" form:"node"`
			Http struct {
				Port []int  `json:"port" form:"port"`
				Node string `json:"node" form:"node"`
			} `json:"http" form:"http"`
		} `json:"gate" form:"gate"`
		File struct {
			Port []int  `json:"port" form:"port"`
			Node string `json:"node" form:"node"`
		} `json:"file" form:"file"`
	} `json:"rpc" form:"rpc"`
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
	c.Log.Dir = ini.GetString("log", "gate.http.dir")
	c.Log.Level = ini.GetInt("log", "gate.http.level")
	c.Log.Mode = ini.GetInt("log", "gate.http.mode")
	c.Log.Style = ini.GetInt("log", "gate.http.style")
	c.Log.Timezone = ini.GetInt("log", "gate.http.timezone")
	// TgBot
	c.TgBot.Enable = ini.GetInt("tg_bot", "enable")
	c.TgBot.ChatId = ini.GetInt64("tg_bot", "chatId")
	c.TgBot.Token = ini.GetString("tg_bot", "token")
	// Path
	c.Path.UpdateCfg = ini.GetString("path", "updateconfig")
	c.Path.GetCfg = ini.GetString("path", "getconfig")
	c.File.Path.Upload = ini.GetString("path", "file.upload")
	c.File.Path.Get = ini.GetString("path", "file.get")
	c.File.Path.Del = ini.GetString("path", "file.del")
	c.File.Path.Fileinfo = ini.GetString("path", "file.fileinfo")
	c.File.Path.FileDetail = ini.GetString("path", "file.filedetail")
	c.File.Path.UuidList = ini.GetString("path", "file.uuidlist")
	c.File.Path.List = ini.GetString("path", "file.list")
	// File
	c.File.Ip = ini.GetString("file", "ip")
	ports := strings.Split(ini.GetString("file", "port"), ",")
	for _, port := range ports {
		switch port == "" {
		case false:
			c.File.Port = append(c.File.Port, utils.Atoi(port))
		}
	}
	c.File.MaxConn = ini.GetInt("file", "maxConn")
	c.File.IdleTimeout = ini.GetInt("file", "idleTimeout")
	c.File.Upload.Dir = ini.GetString("file", "upload.dir")
	c.File.Upload.CheckMd5 = ini.GetInt("file", "upload.checkMd5")
	c.File.Upload.WriteFile = ini.GetInt("file", "upload.writeFile")
	c.File.Upload.MultiFile = ini.GetInt("file", "upload.multiFile")
	c.File.Upload.UseAsync = ini.GetInt("file", "upload.useAsync")
	c.File.Upload.UseOriginFilename = ini.GetInt("file", "upload.useOriginFilename")
	str := ini.GetString("file", "upload.maxMemory")
	slice := strings.Split(str, "*")
	val := int64(1)
	for _, v := range slice {
		v = strings.ReplaceAll(v, " ", "")
		c, _ := strconv.ParseInt(v, 10, 0)
		val *= c
	}
	c.File.Upload.MaxMemory = val
	str = ini.GetString("file", "upload.maxSegmentSize")
	slice = strings.Split(str, "*")
	val = int64(1)
	for _, v := range slice {
		v = strings.ReplaceAll(v, " ", "")
		c, _ := strconv.ParseInt(v, 10, 0)
		val *= c
	}
	c.File.Upload.MaxSegmentSize = val
	str = ini.GetString("file", "upload.maxSingleSize")
	slice = strings.Split(str, "*")
	val = int64(1)
	for _, v := range slice {
		v = strings.ReplaceAll(v, " ", "")
		c, _ := strconv.ParseInt(v, 10, 0)
		val *= c
	}
	c.File.Upload.MaxSingleSize = val
	str = ini.GetString("file", "upload.maxTotalSize")
	slice = strings.Split(str, "*")
	val = int64(1)
	for _, v := range slice {
		v = strings.ReplaceAll(v, " ", "")
		c, _ := strconv.ParseInt(v, 10, 0)
		val *= c
	}
	c.File.Upload.MaxTotalSize = val
	str = ini.GetString("file", "upload.pendingTimeout")
	slice = strings.Split(str, "*")
	val1 := 1
	for _, v := range slice {
		v = strings.ReplaceAll(v, " ", "")
		c, _ := strconv.Atoi(v)
		val1 *= c
	}
	c.File.Upload.PendingTimeout = val1
	str = ini.GetString("file", "upload.fileExpiredTimeout")
	slice = strings.Split(str, "*")
	val1 = 1
	for _, v := range slice {
		v = strings.ReplaceAll(v, " ", "")
		c, _ := strconv.Atoi(v)
		val1 *= c
	}
	c.File.Upload.FileExpiredTimeout = val1
	// Oss
	c.Oss.Type = ini.GetString("oss", "type")
	c.Oss.Aliyun.BasePath = ini.GetString("aliyun", "basePath")
	c.Oss.Aliyun.BucketUrl = ini.GetString("aliyun", "bucketUrl")
	c.Oss.Aliyun.BucketName = ini.GetString("aliyun", "bucketName")
	c.Oss.Aliyun.EndPoint = ini.GetString("aliyun", "endpoint")
	c.Oss.Aliyun.AccessKeyId = ini.GetString("aliyun", "accessKeyId")
	c.Oss.Aliyun.AccessKeySecret = ini.GetString("aliyun", "accessKeySecret")
	c.Oss.Aliyun.Routines = ini.GetInt("aliyun", "routines")
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
	// Gate
	c.Gate.Proto = ini.GetString("gate", "proto")
	c.Gate.Ip = ini.GetString("gate", "ip")
	ports = strings.Split(ini.GetString("gate", "port"), ",")
	for _, port := range ports {
		switch port == "" {
		case false:
			c.Gate.Port = append(c.Gate.Port, utils.Atoi(port))
		}
	}
	c.Gate.MaxConn = ini.GetInt("gate", "maxConn")
	c.Gate.UsePool = ini.GetInt("gate", "usePool")
	c.Gate.HandshakeTimeout = ini.GetInt("gate", "handshakeTimeout")
	c.Gate.IdleTimeout = ini.GetInt("gate", "idleTimeout")
	c.Gate.ReadBufferSize = ini.GetInt("gate", "readBufferSize")
	c.Gate.PrintInterval = ini.GetInt("gate", "printInterval")
	c.Gate.Path.Handshake = ini.GetString("path", "gate.handshake")
	c.Gate.Http.Ip = ini.GetString("gate.http", "ip")
	ports = strings.Split(ini.GetString("gate.http", "port"), ",")
	for _, port := range ports {
		switch port == "" {
		case false:
			c.Gate.Http.Port = append(c.Gate.Http.Port, utils.Atoi(port))
		}
	}
	c.Gate.Http.MaxConn = ini.GetInt("gate.http", "maxConn")
	c.Gate.Http.IdleTimeout = ini.GetInt("gate.http", "idleTimeout")
	c.Gate.Http.Path.Fileserver = ini.GetString("path", "gate.http.fileserver")
	// Rpc
	c.Rpc.Ip = ini.GetString("rpc", "ip")
	c.Rpc.Gate.Node = ini.GetString("rpc", "gate.node")
	ports = strings.Split(ini.GetString("rpc", "gate.port"), ",")
	for _, port := range ports {
		switch port == "" {
		case false:
			c.Rpc.Gate.Port = append(c.Rpc.Gate.Port, utils.Atoi(port))
		}
	}
	c.Rpc.Gate.Http.Node = ini.GetString("rpc", "gate.http.node")
	ports = strings.Split(ini.GetString("rpc", "gate.http.port"), ",")
	for _, port := range ports {
		switch port == "" {
		case false:
			c.Rpc.Gate.Http.Port = append(c.Rpc.Gate.Http.Port, utils.Atoi(port))
		}
	}
	c.Rpc.File.Node = ini.GetString("rpc", "file.node")
	ports = strings.Split(ini.GetString("rpc", "file.port"), ",")
	for _, port := range ports {
		switch port == "" {
		case false:
			c.Rpc.File.Port = append(c.Rpc.File.Port, utils.Atoi(port))
		}
	}
	return
}

func check() {
	if Config.File.Upload.Dir == "" {
		Config.File.Upload.Dir = global.Dir_upload
	}
	if Config.File.Upload.WriteFile > 0 {
		_, err := os.Stat(Config.File.Upload.Dir)
		if err != nil && os.IsNotExist(err) {
			os.MkdirAll(Config.File.Upload.Dir, os.ModePerm)
		}
	}
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

func read(conf string) {
	Config = readIni(conf)
	if Config == nil {
		logs.Fatalf("error")
	}
}

func InitConfig(conf string) {
	read(conf)
	switch Config.Flag {
	case 1:
		flag.Parse()
	default:
	}
	check()
}

func readConfig(conf string) {
	read(conf)
	check()
}

func ReadConfig(conf string) {
	lock.RLock()
	readConfig(conf)
	lock.RUnlock()
}

func updateConfig(conf string, req *global.UpdateCfgReq) {
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
	if req.MaxMemory != "" {
		ini.SetString("file", "upload.maxMemory", req.MaxMemory)
	}
	if req.MaxSegmentSize != "" {
		ini.SetString("file", "upload.maxSegmentSize", req.MaxSegmentSize)
	}
	if req.MaxSingleSize != "" {
		ini.SetString("file", "upload.maxSingleSize", req.MaxSingleSize)
	}
	if req.MaxTotalSize != "" {
		ini.SetString("file", "upload.maxTotalSize", req.MaxTotalSize)
	}
	if req.PendingTimeout != "" {
		ini.SetString("file", "upload.pendingTimeout", req.PendingTimeout)
	}
	if req.FileExpiredTimeout != "" {
		ini.SetString("file", "upload.fileExpiredTimeout", req.FileExpiredTimeout)
	}
	if req.CheckMd5 != "" {
		ini.SetString("file", "upload.checkMd5", req.CheckMd5)
	}
	if req.WriteFile != "" {
		ini.SetString("file", "upload.writeFile", req.WriteFile)
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
	ini.SaveTo(conf)
}

func UpdateConfig(conf string, req *global.UpdateCfgReq) {
	lock.Lock()
	updateConfig(conf, req)
	readConfig(conf)
	lock.Unlock()
}

func GetConfig(req *global.GetCfgReq) (*global.GetCfgResp, bool) {
	lock.RLock()
	resp := &global.GetCfgResp{
		ErrCode: 0,
		ErrMsg:  "ok",
		Data: &global.CfgData{
			Interval:           Config.Interval,
			Log_dir:            Config.Log.Dir,
			Log_level:          Config.Log.Level,
			Log_mode:           Config.Log.Mode,
			Log_style:          Config.Log.Style,
			Log_timezone:       Config.Log.Timezone,
			HttpAddr:           strings.Join([]string{Config.File.Ip, strconv.Itoa(Config.File.Port[0])}, ":"),
			UploadPath:         Config.File.Path.Upload,
			GetPath:            Config.File.Path.Get,
			DelPath:            Config.File.Path.Del,
			FileinfoPath:       Config.File.Path.Fileinfo,
			UpdateCfgPath:      Config.Path.UpdateCfg,
			GetCfgPath:         Config.Path.GetCfg,
			CheckMd5:           Config.File.Upload.CheckMd5,
			WriteFile:          Config.File.Upload.WriteFile,
			MultiFile:          Config.File.Upload.MultiFile,
			UseAsync:           Config.File.Upload.UseAsync,
			MaxMemory:          Config.File.Upload.MaxMemory,
			MaxSegmentSize:     Config.File.Upload.MaxSegmentSize,
			MaxSingleSize:      Config.File.Upload.MaxSingleSize,
			MaxTotalSize:       Config.File.Upload.MaxTotalSize,
			PendingTimeout:     Config.File.Upload.PendingTimeout,
			FileExpiredTimeout: Config.File.Upload.FileExpiredTimeout,
			UploadDir:          Config.File.Upload.Dir,
			OssType:            Config.Oss.Type,
			UseTgBot:           Config.TgBot.Enable,
			TgBotChatId:        Config.TgBot.ChatId,
			TgBotToken:         Config.TgBot.Token,
		},
	}
	lock.RUnlock()
	return resp, true
}
