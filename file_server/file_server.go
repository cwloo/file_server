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
	logs.LogInit(dir+"/logs", logs.LVL_DEBUG, exe, 100000000)
	logs.LogMode(logs.M_STDOUT_FILE)

	http.HandleFunc("/upload", uploadFile)

	task.After(time.Duration(PendingTimeout)*time.Second, cb.NewFunctor00(func() {
		handlerUploadFileOnTimer()
	}))

	err := http.ListenAndServe("192.168.1.113:8088", nil)
	if err != nil {
		logs.LogFatal("%v", err.Error())
	}
	logs.LogInfo("listen:8080")
}
