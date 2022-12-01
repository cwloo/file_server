package main

import (
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/cwloo/gonet/core/base/task"
	"github.com/cwloo/gonet/core/cb"
	"github.com/cwloo/gonet/logs"
)

func uploadFile(w http.ResponseWriter, r *http.Request) {
	handlerUploadFile(w, r)
}

func main() {
	path, _ := os.Executable()
	dir, exe := filepath.Split(path)

	logs.LogTimezone(logs.MY_CST)
	logs.LogInit(dir+"logs", logs.LVL_DEBUG, exe, 100000000)
	logs.LogMode(logs.M_STDOUT_FILE)

	task.After(time.Duration(PendingTimeout)*time.Second, cb.NewFunctor00(func() {
		handlerUploadFileOnTimer()
	}))

	mux := http.NewServeMux()
	mux.HandleFunc("/upload", uploadFile)

	server := &http.Server{
		Addr:              "192.168.1.113:8088",
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
