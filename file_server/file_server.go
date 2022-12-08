package main

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/cwloo/gonet/core/base/task"
	"github.com/cwloo/gonet/core/cb"
	"github.com/cwloo/gonet/logs"
)

func Upload(w http.ResponseWriter, r *http.Request) {
	switch MultiFile {
	case true:
		handlerMultiUpload(w, r)
	default:
		handlerUpload(w, r)
	}
}

func Get(w http.ResponseWriter, r *http.Request) {
	resp := &Resp{
		ErrCode: 0,
		ErrMsg:  "OK",
	}
	j, _ := json.Marshal(resp)
	setResponseHeader(w, r)
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Content-Length", strconv.Itoa(len(j)))
	_, err := w.Write(j)
	if err != nil {
		logs.LogError(err.Error())
		return
	}
	logs.LogDebug("%v", string(j))
}

func main() {
	InitConfig()
	// logs.LogTimezone(logs.MY_CST)
	// logs.LogMode(logs.M_STDOUT_FILE)
	// logs.LogStyle(logs.F_DETAIL)
	// logs.LogInit(dir+"logs", logs.LVL_DEBUG, exe, 100000000)
	logs.LogTimezone(logs.TimeZone(Config.Log_timezone))
	logs.LogMode(logs.Mode(Config.Log_mode))
	logs.LogStyle(logs.Style(Config.Log_style))
	logs.LogInit(dir+"logs", int32(Config.Log_level), exe, 100000000)

	task.After(time.Duration(PendingTimeout)*time.Second, cb.NewFunctor00(func() {
		handlerPendingUploader()
	}))

	task.After(time.Duration(FileExpiredTimeout)*time.Second, cb.NewFunctor00(func() {
		handlerExpiredFile()
	}))

	mux := http.NewServeMux()
	mux.HandleFunc(Config.UploadPath, Upload)
	mux.HandleFunc(Config.GetPath, Get)

	server := &http.Server{
		Addr:              Config.HttpAddr,
		Handler:           mux,
		ReadTimeout:       time.Duration(PendingTimeout) * time.Second,
		ReadHeaderTimeout: time.Duration(PendingTimeout) * time.Second,
		WriteTimeout:      time.Duration(PendingTimeout) * time.Second,
		IdleTimeout:       time.Duration(PendingTimeout) * time.Second,
	}

	logs.LogInfo(server.Addr)

	server.SetKeepAlivesEnabled(true)
	err := server.ListenAndServe()
	if err != nil {
		logs.LogFatal(err.Error())
	}
}
