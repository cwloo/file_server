package main

import (
	"net/http"
	"time"

	"github.com/cwloo/gonet/core/base/task"
	"github.com/cwloo/gonet/core/cb"
	"github.com/cwloo/gonet/logs"
)

func uploadFile(w http.ResponseWriter, r *http.Request) {
	handlerUploadFile(w, r)
}

func main() {
	InitConfig()
	// logs.LogTimezone(logs.MY_CST)
	// logs.LogMode(logs.M_STDOUT_FILE)
	// logs.LogInit(dir+"logs", logs.LVL_DEBUG, exe, 100000000)
	logs.LogTimezone(logs.TimeZone(Config.Log_timezone))
	logs.LogMode(logs.Mode(Config.Log_mode))
	logs.LogInit(dir+"logs", int32(Config.Log_level), exe, 100000000)

	task.After(time.Duration(PendingTimeout)*time.Second, cb.NewFunctor00(func() {
		handlerPendingUploader()
	}))

	task.After(time.Duration(FileExpiredTimeout)*time.Second, cb.NewFunctor00(func() {
		handlerExpiredFile()
	}))

	mux := http.NewServeMux()
	mux.HandleFunc("/upload", uploadFile)

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
		logs.LogFatal("%v", err.Error())
	}
}
