package main

import (
	"net/http"
	"time"

	"github.com/cwloo/gonet/core/base/task"
	"github.com/cwloo/gonet/core/cb"
	"github.com/cwloo/gonet/logs"
	"github.com/cwloo/uploader/file_server/config"
	"github.com/cwloo/uploader/file_server/global"
)

func Upload(w http.ResponseWriter, r *http.Request) {
	switch config.Config.MultiFile {
	default:
		handlerMultiUpload(w, r)
	case 0:
		handlerUpload(w, r)
	}
}

func Get(w http.ResponseWriter, r *http.Request) {
	// resp := &Resp{
	// 	ErrCode: 0,
	// 	ErrMsg:  "OK",
	// }
	// writeResponse(w, r, resp)
	handlerFileinfo(w, r)
}

func Del(w http.ResponseWriter, r *http.Request) {
	handlerDelCache(w, r)
}

func GetFileinfo(w http.ResponseWriter, r *http.Request) {
	handlerFileinfo(w, r)
}

func UpdateConfig(w http.ResponseWriter, r *http.Request) {
	handlerUpdateCfg(w, r)
}

func main() {
	config.InitConfig()
	// logs.LogTimezone(logs.MY_CST)
	// logs.LogMode(logs.M_STDOUT_FILE)
	// logs.LogStyle(logs.F_DETAIL)
	// logs.LogInit(global.Dir+"logs", logs.LVL_DEBUG, global.Exe, 100000000)
	logs.LogTimezone(logs.TimeZone(config.Config.Log_timezone))
	logs.LogMode(logs.Mode(config.Config.Log_mode))
	logs.LogStyle(logs.Style(config.Config.Log_style))
	logs.LogInit(config.Config.Log_dir, int32(config.Config.Log_level), global.Exe, 100000000)

	task.After(time.Duration(config.Config.PendingTimeout)*time.Second, cb.NewFunctor00(func() {
		handlerPendingUploader()
	}))

	task.After(time.Duration(config.Config.FileExpiredTimeout)*time.Second, cb.NewFunctor00(func() {
		handlerExpiredFile()
	}))

	task.After(time.Duration(config.Config.Interval)*time.Second, cb.NewFunctor00(func() {
		handlerReadConfig()
	}))

	mux := http.NewServeMux()
	mux.HandleFunc(config.Config.UploadPath, Upload)
	mux.HandleFunc(config.Config.GetPath, Get)
	mux.HandleFunc(config.Config.DelPath, Del)
	mux.HandleFunc(config.Config.FileinfoPath, GetFileinfo)
	mux.HandleFunc(config.Config.UpdateCfgPath, UpdateConfig)

	server := &http.Server{
		Addr:              config.Config.HttpAddr,
		Handler:           mux,
		ReadTimeout:       time.Duration(config.Config.PendingTimeout) * time.Second,
		ReadHeaderTimeout: time.Duration(config.Config.PendingTimeout) * time.Second,
		WriteTimeout:      time.Duration(config.Config.PendingTimeout) * time.Second,
		IdleTimeout:       time.Duration(config.Config.PendingTimeout) * time.Second,
	}

	logs.LogInfo(server.Addr)

	server.SetKeepAlivesEnabled(true)
	err := server.ListenAndServe()
	if err != nil {
		logs.LogFatal(err.Error())
	}
}
