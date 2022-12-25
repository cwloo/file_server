package config

import (
	"flag"
	"os"
	"strconv"
	"strings"
	"sync"

	"github.com/cwloo/gonet/core/base/sys/cmd"
	"github.com/cwloo/gonet/core/net/conn"
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
		Monitor struct {
			Dir      string `json:"dir" form:"dir"`
			Level    int    `json:"level" form:"level"`
			Mode     int    `json:"mode" form:"mode"`
			Style    int    `json:"style" form:"style"`
			Timezone int    `json:"timezone" form:"timezone"`
		} `json:"monitor" form:"monitor"`
		Client struct {
			Dir      string `json:"dir" form:"dir"`
			Level    int    `json:"level" form:"level"`
			Mode     int    `json:"mode" form:"mode"`
			Style    int    `json:"style" form:"style"`
			Timezone int    `json:"timezone" form:"timezone"`
		} `json:"client" form:"client"`
		Gate struct {
			Dir      string `json:"dir" form:"dir"`
			Level    int    `json:"level" form:"level"`
			Mode     int    `json:"mode" form:"mode"`
			Style    int    `json:"style" form:"style"`
			Timezone int    `json:"timezone" form:"timezone"`
			Http     struct {
				Dir      string `json:"dir" form:"dir"`
				Level    int    `json:"level" form:"level"`
				Mode     int    `json:"mode" form:"mode"`
				Style    int    `json:"style" form:"style"`
				Timezone int    `json:"timezone" form:"timezone"`
			} `json:"http" form:"http"`
		} `json:"gate" form:"gate"`
		File struct {
			Dir      string `json:"dir" form:"dir"`
			Level    int    `json:"level" form:"level"`
			Mode     int    `json:"mode" form:"mode"`
			Style    int    `json:"style" form:"style"`
			Timezone int    `json:"timezone" form:"timezone"`
		} `json:"file" form:"file"`
	} `json:"log" form:"log"`
	Sub struct {
		Client struct {
			Num  int    `json:"num" form:"num"`
			Dir  string `json:"dir" form:"dir"`
			Exec string `json:"exec" form:"exec"`
		} `json:"client" form:"client"`
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
	Etcd struct {
		Schema   string   `json:"schema" form:"schema"`
		Addr     []string `json:"addr" form:"addr"`
		UserName string   `json:"username" form:"username"`
		Password string   `json:"password" form:"password"`
	} `json:"etcd" form:"etcd"`
	Monitor struct {
		Name        string `json:"name" form:"name"`
		Ip          string `json:"ip" form:"ip"`
		Port        []int  `json:"port" form:"port"`
		MaxConn     int    `json:"maxConn" form:"maxConn"`
		IdleTimeout int    `json:"idleTimeout" form:"idleTimeout"`
		Path        struct {
			Start   string `json:"start" form:"start"`
			Kill    string `json:"kill" form:"kill"`
			KillAll string `json:"killall" form:"killall"`
			SubList string `json:"sublist" form:"sublist"`
		} `json:"path" form:"path"`
	} `json:"monitor" form:"monitor"`
	Client struct {
		Name string `json:"name" form:"name"`
		Addr []struct {
			Proto string `json:"proto" form:"proto"`
			Ip    string `json:"ip" form:"ip"`
			Port  int    `json:"port" form:"port"`
		} `json:"addr" form:"addr"`
		Path struct {
			Router   string `json:"router" form:"router"`
			Upload   string `json:"upload" form:"upload"`
			Get      string `json:"get" form:"get"`
			Fileinfo string `json:"fileinfo" form:"fileinfo"`
		} `json:"path" form:"path"`
		Upload struct {
			SegmentSize int64    `json:"segmentSize" form:"segmentSize"`
			MultiFile   int      `json:"multiFile" form:"multiFile"`
			Filelist    []string `json:"filelist" form:"filelist"`
		} `json:"upload" form:"upload"`
	} `json:"client" form:"client"`
	Gate struct {
		Name  string `json:"name" form:"name"`
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
			Name        string `json:"name" form:"name"`
			Ip          string `json:"ip" form:"ip"`
			Port        []int  `json:"port" form:"port"`
			MaxConn     int    `json:"maxConn" form:"maxConn"`
			IdleTimeout int    `json:"idleTimeout" form:"idleTimeout"`
			Path        struct {
				Router string `json:"router" form:"router"`
			} `json:"path" form:"path"`
		} `json:"http" form:"http"`
	} `json:"gate" form:"gate"`
	File struct {
		Name        string   `json:"name" form:"name"`
		Ip          string   `json:"ip" form:"ip"`
		Port        []int    `json:"port" form:"port"`
		MaxConn     int      `json:"maxConn" form:"maxConn"`
		IdleTimeout int      `json:"idleTimeout" form:"idleTimeout"`
		Domain      []string `json:"domain" form:"domain"`
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
	Rpc struct {
		Ip      string `json:"ip" form:"ip"`
		Monitor struct {
			Port []int  `json:"port" form:"port"`
			Node string `json:"node" form:"node"`
		} `json:"monitor" form:"monitor"`
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
	Mysql struct {
		Addr          []string `json:"addr" form:"addr"`
		Username      string   `json:"username" form:"username"`
		Password      string   `json:"password" form:"password"`
		Database      string   `json:"database" form:"database"`
		Tablename     string   `json:"tablename" form:"tablename"`
		MaxConn       int      `json:"maxConn" form:"maxConn"`
		MaxIdleConn   int      `json:"maxIdleConn" form:"maxIdleConn"`
		MaxLifeTime   int      `json:"maxLifeTime" form:"maxLifeTime"`
		SetLevel      int      `json:"SetLevel" form:"SetLevel"`
		SlowThreshold int      `json:"slowThreshold" form:"slowThreshold"`
	} `json:"mysql" form:"mysql"`
	Mongo struct {
		Url         string   `json:"url" form:"url"`
		Addr        []string `json:"addr" form:"addr"`
		Direct      bool     `json:"direct" form:"direct"`
		Timeout     int      `json:"timeout" form:"timeout"`
		Database    string   `json:"database" form:"database"`
		Source      string   `json:"source" form:"source"`
		Username    string   `json:"username" form:"username"`
		Password    string   `json:"password" form:"password"`
		MaxPoolSize int      `json:"maxPoolSize" form:"maxPoolSize"`
	} `json:"mongo" form:"mongo"`
	Redis struct {
		Addr        []string `json:"addr" form:"addr"`
		MaxIdle     int      `json:"maxIdle" form:"maxIdle"`
		MaxActive   int      `json:"maxActive" form:"maxActive"`
		IdleTimeout int      `json:"idleTimeout" form:"idleTimeout"`
		Username    string   `json:"username" form:"username"`
		Password    string   `json:"password" form:"password"`
		Cluster     bool     `json:"cluster" form:"cluster"`
	} `json:"redis" form:"redis"`
	Kafka struct {
		SASLUserName string `json:"SASLUserName" form:"SASLUserName"`
		SASLPassword string `json:"SASLPassword" form:"SASLPassword"`
		Ws2mschat    struct {
			Addr  []string `json:"addr" form:"addr"`
			Topic string   `json:"topic" form:"topic"`
		}
		MsgToMongo struct {
			Addr  []string `json:"addr" form:"addr"`
			Topic string   `json:"topic" form:"topic"`
		}
		Ms2pschat struct {
			Addr  []string `json:"addr" form:"addr"`
			Topic string   `json:"topic" form:"topic"`
		}
		ConsumerGroupID struct {
			MsgToRedis string `json:"msgToTransfer" form:"msgToTransfer"`
			MsgToMongo string `json:"msgToMongo" form:"msgToMongo"`
			MsgToMySql string `json:"msgToMySql" form:"msgToMySql"`
			MsgToPush  string `json:"msgToPush" form:"msgToPush"`
		}
	} `json:"kafka" form:"kafka"`
	Prometheus struct {
		Enable bool   `json:"enable" form:"enable"`
		Ip     string `json:"ip" form:"ip"`
		Gate   struct {
			Port []int `json:"port" form:"port"`
		} `json:"gate" form:"gate"`
		Msg struct {
			Port []int `json:"port" form:"port"`
		} `json:"msg" form:"msg"`
		Push struct {
			Port []int `json:"port" form:"port"`
		} `json:"push" form:"push"`
		Transfer struct {
			Port []int `json:"port" form:"port"`
		} `json:"transfer" form:"transfer"`
		User struct {
			Port []int `json:"port" form:"port"`
		} `json:"user" form:"user"`
		Friend struct {
			Port []int `json:"port" form:"port"`
		} `json:"friend" form:"friend"`
		Group struct {
			Port []int `json:"port" form:"port"`
		} `json:"group" form:"group"`
		Auth struct {
			Port []int `json:"port" form:"port"`
		} `json:"auth" form:"auth"`
		Cache struct {
			Port []int `json:"port" form:"port"`
		} `json:"cache" form:"cache"`
		Admin struct {
			Port []int `json:"port" form:"port"`
		} `json:"admin" form:"admin"`
		RealTime struct {
			Port []int `json:"port" form:"port"`
		} `json:"realtime" form:"realtime"`
	} `json:"prometheus" form:"prometheus"`
}

func readIni(filename string, cb func(*IniConfig) string) (c *IniConfig) {
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
	c.Monitor.Name = ini.GetString("monitor", "name")
	c.Client.Name = ini.GetString("client", "name")
	c.Gate.Name = ini.GetString("gate", "name")
	c.Gate.Http.Name = ini.GetString("gate.http", "name")
	c.File.Name = ini.GetString("file", "name")
	setServiceName(cb, c)
	// Log
	c.Log.Monitor.Dir = ini.GetString("log", "monitor.dir")
	c.Log.Monitor.Level = ini.GetInt("log", "monitor.level")
	c.Log.Monitor.Mode = ini.GetInt("log", "monitor.mode")
	c.Log.Monitor.Style = ini.GetInt("log", "monitor.style")
	c.Log.Monitor.Timezone = ini.GetInt("log", "monitor.timezone")
	c.Log.Client.Dir = ini.GetString("log", "client.dir")
	c.Log.Client.Level = ini.GetInt("log", "client.level")
	c.Log.Client.Mode = ini.GetInt("log", "client.mode")
	c.Log.Client.Style = ini.GetInt("log", "client.style")
	c.Log.Client.Timezone = ini.GetInt("log", "client.timezone")
	c.Log.Gate.Dir = ini.GetString("log", "gate.dir")
	c.Log.Gate.Level = ini.GetInt("log", "gate.level")
	c.Log.Gate.Mode = ini.GetInt("log", "gate.mode")
	c.Log.Gate.Style = ini.GetInt("log", "gate.style")
	c.Log.Gate.Timezone = ini.GetInt("log", "gate.timezone")
	c.Log.Gate.Http.Dir = ini.GetString("log", "gate.http.dir")
	c.Log.Gate.Http.Level = ini.GetInt("log", "gate.http.level")
	c.Log.Gate.Http.Mode = ini.GetInt("log", "gate.http.mode")
	c.Log.Gate.Http.Style = ini.GetInt("log", "gate.http.style")
	c.Log.Gate.Http.Timezone = ini.GetInt("log", "gate.http.timezone")
	c.Log.File.Dir = ini.GetString("log", "file.dir")
	c.Log.File.Level = ini.GetInt("log", "file.level")
	c.Log.File.Mode = ini.GetInt("log", "file.mode")
	c.Log.File.Style = ini.GetInt("log", "file.style")
	c.Log.File.Timezone = ini.GetInt("log", "file.timezone")
	// Sub
	c.Sub.Client.Num = ini.GetInt("sub", "client.num")
	c.Sub.Client.Dir = ini.GetString("sub", "client.dir")
	c.Sub.Client.Exec = ini.GetString("sub", "client.execname")
	c.Sub.Gate.Num = ini.GetInt("sub", "gate.num")
	c.Sub.Gate.Dir = ini.GetString("sub", "gate.dir")
	c.Sub.Gate.Exec = ini.GetString("sub", "gate.execname")
	c.Sub.Gate.Http.Num = ini.GetInt("sub", "gate.http.num")
	c.Sub.Gate.Http.Dir = ini.GetString("sub", "gate.http.dir")
	c.Sub.Gate.Http.Exec = ini.GetString("sub", "gate.http.execname")
	c.Sub.File.Num = ini.GetInt("sub", "file.num")
	c.Sub.File.Dir = ini.GetString("sub", "file.dir")
	c.Sub.File.Exec = ini.GetString("sub", "file.execname")
	// TgBot
	c.TgBot.Enable = ini.GetInt("tg_bot", "enable")
	c.TgBot.ChatId = ini.GetInt64("tg_bot", "chatId")
	c.TgBot.Token = ini.GetString("tg_bot", "token")
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
	// Monitor
	// c.Monitor.Name = ini.GetString("monitor", "name")
	c.Monitor.Ip = ini.GetString("monitor", "ip")
	ports := strings.Split(ini.GetString("monitor", "port"), ",")
	for _, port := range ports {
		switch port == "" {
		case false:
			c.Monitor.Port = append(c.Monitor.Port, utils.Atoi(port))
		}
	}
	c.Monitor.MaxConn = ini.GetInt("monitor", "maxConn")
	c.Monitor.IdleTimeout = ini.GetInt("monitor", "idleTimeout")
	// Client
	// c.Client.Name = ini.GetString("client", "name")
	addrs = strings.Split(ini.GetString("client", "addr"), ",")
	for _, addr := range addrs {
		switch addr == "" {
		case false:
			v := conn.ParseAddress(addr)
			switch v {
			case nil:
				logs.Fatalf("error")
			default:
				c.Client.Addr = append(c.Client.Addr, struct {
					Proto string `json:"proto" form:"proto"`
					Ip    string `json:"ip" form:"ip"`
					Port  int    `json:"port" form:"port"`
				}{
					Proto: v.Proto,
					Ip:    v.Ip,
					Port:  utils.Atoi(v.Port),
				})
			}
		}
	}
	c.Client.Path.Router = ini.GetString("client", "path.router")
	c.Client.Path.Upload = ini.GetString("client", "path.upload")
	c.Client.Path.Get = ini.GetString("client", "path.get")
	c.Client.Path.Fileinfo = ini.GetString("client", "path.fileinfo")
	str := ini.GetString("client", "upload.segmentSize")
	slice := strings.Split(str, "*")
	val := int64(1)
	for _, v := range slice {
		v = strings.ReplaceAll(v, " ", "")
		c, _ := strconv.ParseInt(v, 10, 0)
		val *= c
	}
	c.Client.Upload.SegmentSize = val
	c.Client.Upload.MultiFile = ini.GetInt("client", "upload.multiFile")
	num := ini.GetInt("client", "upload.num")
	for i := 0; i < num; i++ {
		c.Client.Upload.Filelist = append(
			c.Client.Upload.Filelist,
			ini.GetString("client", strings.Join([]string{"upload.file", strconv.Itoa(i)}, "")))
	}
	// Path
	c.Path.UpdateCfg = ini.GetString("path", "updateconfig")
	c.Path.GetCfg = ini.GetString("path", "getconfig")
	c.Monitor.Path.Start = ini.GetString("path", "monitor.start")
	c.Monitor.Path.Kill = ini.GetString("path", "monitor.kill")
	c.Monitor.Path.KillAll = ini.GetString("path", "monitor.killall")
	c.Monitor.Path.SubList = ini.GetString("path", "monitor.sublist")
	c.File.Path.Upload = ini.GetString("path", "file.upload")
	c.File.Path.Get = ini.GetString("path", "file.get")
	c.File.Path.Del = ini.GetString("path", "file.del")
	c.File.Path.Fileinfo = ini.GetString("path", "file.fileinfo")
	c.File.Path.FileDetail = ini.GetString("path", "file.filedetail")
	c.File.Path.UuidList = ini.GetString("path", "file.uuidlist")
	c.File.Path.List = ini.GetString("path", "file.list")
	// File
	// c.File.Name = ini.GetString("file", "name")
	c.File.Ip = ini.GetString("file", "ip")
	ports = strings.Split(ini.GetString("file", "port"), ",")
	for _, port := range ports {
		switch port == "" {
		case false:
			c.File.Port = append(c.File.Port, utils.Atoi(port))
		}
	}
	c.File.MaxConn = ini.GetInt("file", "maxConn")
	c.File.IdleTimeout = ini.GetInt("file", "idleTimeout")
	domains := strings.Split(ini.GetString("file", "domain"), ",")
	for _, domain := range domains {
		switch domain == "" {
		case false:
			c.File.Domain = append(c.File.Domain, domain)
		}
	}
	c.File.Upload.Dir = ini.GetString("file", "upload.dir")
	c.File.Upload.CheckMd5 = ini.GetInt("file", "upload.checkMd5")
	c.File.Upload.WriteFile = ini.GetInt("file", "upload.writeFile")
	c.File.Upload.MultiFile = ini.GetInt("file", "upload.multiFile")
	c.File.Upload.UseAsync = ini.GetInt("file", "upload.useAsync")
	c.File.Upload.UseOriginFilename = ini.GetInt("file", "upload.useOriginFilename")
	str = ini.GetString("file", "upload.maxMemory")
	slice = strings.Split(str, "*")
	val = int64(1)
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
	// Gate
	// c.Gate.Name = ini.GetString("gate", "name")
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
	// c.Gate.Http.Name = ini.GetString("gate.http", "name")
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
	c.Gate.Http.Path.Router = ini.GetString("path", "gate.http.router")
	// Rpc
	c.Rpc.Ip = ini.GetString("rpc", "ip")
	c.Rpc.Monitor.Node = ini.GetString("rpc", "monitor.node")
	ports = strings.Split(ini.GetString("rpc", "monitor.port"), ",")
	for _, port := range ports {
		switch port == "" {
		case false:
			c.Rpc.Monitor.Port = append(c.Rpc.Monitor.Port, utils.Atoi(port))
		}
	}
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
	// Mysql
	addrs = strings.Split(ini.GetString("mysql", "addr"), ",")
	for _, addr := range addrs {
		switch addr == "" {
		case false:
			c.Mysql.Addr = append(c.Mysql.Addr, addr)
		}
	}
	c.Mysql.Username = ini.GetString("mysql", "username")
	c.Mysql.Password = ini.GetString("mysql", "password")
	c.Mysql.Database = ini.GetString("mysql", "database")
	c.Mysql.Tablename = ini.GetString("mysql", "tablename")
	c.Mysql.MaxConn = ini.GetInt("mysql", "maxConn")
	c.Mysql.MaxIdleConn = ini.GetInt("mysql", "maxIdleConn")
	c.Mysql.MaxLifeTime = ini.GetInt("mysql", "maxLifeTime")
	c.Mysql.SetLevel = ini.GetInt("mysql", "SetLevel")
	c.Mysql.SlowThreshold = ini.GetInt("mysql", "slowThreshold")
	// Mongo
	c.Mongo.Url = ini.GetString("mongo", "url")
	addrs = strings.Split(ini.GetString("mongo", "addr"), ",")
	for _, addr := range addrs {
		switch addr == "" {
		case false:
			c.Mongo.Addr = append(c.Mongo.Addr, addr)
		}
	}
	c.Mongo.Direct = ini.GetInt("mongo", "direct") > 0
	c.Mongo.Timeout = ini.GetInt("mongo", "timeout")
	c.Mongo.Database = ini.GetString("mongo", "database")
	c.Mongo.Source = ini.GetString("mongo", "source")
	c.Mongo.Username = ini.GetString("mongo", "username")
	c.Mongo.Password = ini.GetString("mongo", "password")
	c.Mongo.MaxPoolSize = ini.GetInt("mongo", "maxPoolSize")
	// Redis
	addrs = strings.Split(ini.GetString("redis", "addr"), ",")
	for _, addr := range addrs {
		switch addr == "" {
		case false:
			c.Redis.Addr = append(c.Redis.Addr, addr)
		}
	}
	c.Redis.MaxIdle = ini.GetInt("redis", "maxIdle")
	c.Redis.MaxActive = ini.GetInt("redis", "maxActive")
	c.Redis.IdleTimeout = ini.GetInt("redis", "idleTimeout")
	c.Redis.Username = ini.GetString("redis", "username")
	c.Redis.Password = ini.GetString("redis", "password")
	c.Redis.Cluster = ini.GetInt("redis", "cluster") > 0
	// Prometheus
	c.Prometheus.Enable = ini.GetInt("prometheus", "enable") > 0
	c.Prometheus.Ip = ini.GetString("prometheus", "ip")
	ports = strings.Split(ini.GetString("prometheus", "gate.port"), ",")
	for _, port := range ports {
		switch port == "" {
		case false:
			c.Prometheus.Gate.Port = append(c.Prometheus.Gate.Port, utils.Atoi(port))
		}
	}
	ports = strings.Split(ini.GetString("prometheus", "msg.port"), ",")
	for _, port := range ports {
		switch port == "" {
		case false:
			c.Prometheus.Msg.Port = append(c.Prometheus.Msg.Port, utils.Atoi(port))
		}
	}
	ports = strings.Split(ini.GetString("prometheus", "push.port"), ",")
	for _, port := range ports {
		switch port == "" {
		case false:
			c.Prometheus.Push.Port = append(c.Prometheus.Push.Port, utils.Atoi(port))
		}
	}
	ports = strings.Split(ini.GetString("prometheus", "user.port"), ",")
	for _, port := range ports {
		switch port == "" {
		case false:
			c.Prometheus.User.Port = append(c.Prometheus.User.Port, utils.Atoi(port))
		}
	}
	ports = strings.Split(ini.GetString("prometheus", "friend.port"), ",")
	for _, port := range ports {
		switch port == "" {
		case false:
			c.Prometheus.Friend.Port = append(c.Prometheus.Friend.Port, utils.Atoi(port))
		}
	}
	ports = strings.Split(ini.GetString("prometheus", "group.port"), ",")
	for _, port := range ports {
		switch port == "" {
		case false:
			c.Prometheus.Group.Port = append(c.Prometheus.Group.Port, utils.Atoi(port))
		}
	}
	ports = strings.Split(ini.GetString("prometheus", "auth.port"), ",")
	for _, port := range ports {
		switch port == "" {
		case false:
			c.Prometheus.Auth.Port = append(c.Prometheus.Auth.Port, utils.Atoi(port))
		}
	}
	ports = strings.Split(ini.GetString("prometheus", "cache.port"), ",")
	for _, port := range ports {
		switch port == "" {
		case false:
			c.Prometheus.Cache.Port = append(c.Prometheus.Cache.Port, utils.Atoi(port))
		}
	}
	ports = strings.Split(ini.GetString("prometheus", "admin.port"), ",")
	for _, port := range ports {
		switch port == "" {
		case false:
			c.Prometheus.Admin.Port = append(c.Prometheus.Admin.Port, utils.Atoi(port))
		}
	}
	ports = strings.Split(ini.GetString("prometheus", "admin.port"), ",")
	for _, port := range ports {
		switch port == "" {
		case false:
			c.Prometheus.Admin.Port = append(c.Prometheus.Admin.Port, utils.Atoi(port))
		}
	}
	ports = strings.Split(ini.GetString("prometheus", "realTime.port"), ",")
	for _, port := range ports {
		switch port == "" {
		case false:
			c.Prometheus.RealTime.Port = append(c.Prometheus.RealTime.Port, utils.Atoi(port))
		}
	}
	ports = strings.Split(ini.GetString("prometheus", "transfer.port"), ",")
	for _, port := range ports {
		switch port == "" {
		case false:
			c.Prometheus.Transfer.Port = append(c.Prometheus.Transfer.Port, utils.Atoi(port))
		}
	}
	return
}

func check() {
	switch ServiceName() {
	case Config.Monitor.Name:
		switch cmd.Log() == "" {
		case true:
			switch Config.Log.Monitor.Dir == "" {
			case true:
				Config.Log.Monitor.Dir = global.Dir + "log"
			default:
			}
		default:
			Config.Log.Monitor.Dir = cmd.Log()
		}
		switch Config.Log.Monitor.Timezone != int(logs.GetTimezone()) {
		case true:
			logs.SetTimezone(logs.Timezone(Config.Log.Monitor.Timezone))
		}
		switch Config.Log.Monitor.Mode != int(logs.GetMode()) {
		case true:
			logs.SetMode(logs.Mode(Config.Log.Monitor.Mode))
		}
		switch Config.Log.Monitor.Style != int(logs.GetStyle()) {
		case true:
			logs.SetStyle(logs.Style(Config.Log.Monitor.Style))
		}
		switch Config.Log.Monitor.Level != int(logs.GetLevel()) {
		case true:
			logs.SetLevel(logs.Level(Config.Log.Monitor.Level))
		}
	case Config.Gate.Name:
		switch cmd.Log() == "" {
		case true:
			switch Config.Log.Gate.Dir == "" {
			case true:
				Config.Log.Gate.Dir = global.Dir + "log"
			default:
			}
		default:
			Config.Log.Gate.Dir = cmd.Log()
		}
		switch Config.Log.Gate.Timezone != int(logs.GetTimezone()) {
		case true:
			logs.SetTimezone(logs.Timezone(Config.Log.Gate.Timezone))
		}
		switch Config.Log.Gate.Mode != int(logs.GetMode()) {
		case true:
			logs.SetMode(logs.Mode(Config.Log.Gate.Mode))
		}
		switch Config.Log.Gate.Style != int(logs.GetStyle()) {
		case true:
			logs.SetStyle(logs.Style(Config.Log.Gate.Style))
		}
		switch Config.Log.Gate.Level != int(logs.GetLevel()) {
		case true:
			logs.SetLevel(logs.Level(Config.Log.Gate.Level))
		}
	case Config.Gate.Http.Name:
		switch cmd.Log() == "" {
		case true:
			switch Config.Log.Gate.Http.Dir == "" {
			case true:
				Config.Log.Gate.Http.Dir = global.Dir + "log"
			default:
			}
		default:
			Config.Log.Gate.Http.Dir = cmd.Log()
		}
		switch Config.Log.Gate.Http.Timezone != int(logs.GetTimezone()) {
		case true:
			logs.SetTimezone(logs.Timezone(Config.Log.Gate.Http.Timezone))
		}
		switch Config.Log.Gate.Http.Mode != int(logs.GetMode()) {
		case true:
			logs.SetMode(logs.Mode(Config.Log.Gate.Http.Mode))
		}
		switch Config.Log.Gate.Http.Style != int(logs.GetStyle()) {
		case true:
			logs.SetStyle(logs.Style(Config.Log.Gate.Http.Style))
		}
		switch Config.Log.Gate.Http.Level != int(logs.GetLevel()) {
		case true:
			logs.SetLevel(logs.Level(Config.Log.Gate.Http.Level))
		}
	case Config.File.Name:
		switch cmd.Log() == "" {
		case true:
			switch Config.Log.File.Dir == "" {
			case true:
				Config.Log.File.Dir = global.Dir + "log"
			default:
			}
		default:
			Config.Log.File.Dir = cmd.Log()
		}
		switch Config.Log.File.Timezone != int(logs.GetTimezone()) {
		case true:
			logs.SetTimezone(logs.Timezone(Config.Log.File.Timezone))
		}
		switch Config.Log.File.Mode != int(logs.GetMode()) {
		case true:
			logs.SetMode(logs.Mode(Config.Log.File.Mode))
		}
		switch Config.Log.File.Style != int(logs.GetStyle()) {
		case true:
			logs.SetStyle(logs.Style(Config.Log.File.Style))
		}
		switch Config.Log.File.Level != int(logs.GetLevel()) {
		case true:
			logs.SetLevel(logs.Level(Config.Log.File.Level))
		}
		switch Config.File.Upload.Dir == "" {
		case true:
			Config.File.Upload.Dir = global.Dir_upload
		}
		switch Config.File.Upload.WriteFile > 0 {
		case true:
			_, err := os.Stat(Config.File.Upload.Dir)
			if err != nil && os.IsNotExist(err) {
				os.MkdirAll(Config.File.Upload.Dir, os.ModePerm)
			}
		}
	}
	// 中国大陆这里可能因为被墙了卡住
	tg_bot.NewTgBot(Config.TgBot.Token, Config.TgBot.ChatId, Config.TgBot.Enable > 0)
}

func read(conf string, cb func(*IniConfig) string) {
	Config = readIni(conf, cb)
	if Config == nil {
		logs.Fatalf("error")
	}
}

func InitMonitorConfig(conf string) {
	read(conf, func(c *IniConfig) string {
		return c.Monitor.Name
	})
	switch Config.Flag {
	case 1:
		flag.Parse()
	default:
	}
	check()
}

func InitClientConfig(conf string) {
	read(conf, func(c *IniConfig) string {
		return c.Client.Name
	})
	switch Config.Flag {
	case 1:
		flag.Parse()
	default:
	}
	check()
}

func InitGateConfig(conf string) {
	read(conf, func(c *IniConfig) string {
		return c.Gate.Name
	})
	switch Config.Flag {
	case 1:
		flag.Parse()
	default:
	}
	check()
}

func InitHttpGateConfig(conf string) {
	read(conf, func(c *IniConfig) string {
		return c.Gate.Http.Name
	})
	switch Config.Flag {
	case 1:
		flag.Parse()
	default:
	}
	check()
}

func InitFileConfig(conf string) {
	read(conf, func(c *IniConfig) string {
		return c.File.Name
	})
	switch Config.Flag {
	case 1:
		flag.Parse()
	default:
	}
	check()
}

func readConfig(conf string) {
	read(conf, nil)
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
	dir, level, mode, style, timezone := "", 0, 0, 0, 0
	lock.RLock()
	switch ServiceName() {
	case Config.Monitor.Name:
		dir = Config.Log.Monitor.Dir
		level = Config.Log.Monitor.Level
		mode = Config.Log.Monitor.Mode
		style = Config.Log.Monitor.Style
		timezone = Config.Log.Monitor.Timezone
	case Config.Gate.Name:
		dir = Config.Log.Monitor.Dir
		level = Config.Log.Monitor.Level
		mode = Config.Log.Monitor.Mode
		style = Config.Log.Monitor.Style
		timezone = Config.Log.Monitor.Timezone
	case Config.Gate.Http.Name:
		dir = Config.Log.Monitor.Dir
		level = Config.Log.Monitor.Level
		mode = Config.Log.Monitor.Mode
		style = Config.Log.Monitor.Style
		timezone = Config.Log.Monitor.Timezone
	case Config.File.Name:
		dir = Config.Log.Monitor.Dir
		level = Config.Log.Monitor.Level
		mode = Config.Log.Monitor.Mode
		style = Config.Log.Monitor.Style
		timezone = Config.Log.Monitor.Timezone
	}
	resp := &global.GetCfgResp{
		ErrCode: 0,
		ErrMsg:  "ok",
		Data: &global.CfgData{
			Interval:           Config.Interval,
			Log_dir:            dir,
			Log_level:          level,
			Log_mode:           mode,
			Log_style:          style,
			Log_timezone:       timezone,
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
