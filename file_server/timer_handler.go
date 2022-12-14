package main

import (
	"time"

	"github.com/cwloo/gonet/core/base/task"
	"github.com/cwloo/gonet/core/cb"
	"github.com/cwloo/uploader/file_server/config"
)

func handlerReadConfig() {
	config.ReadConfig()
	task.After(time.Duration(config.Config.Interval)*time.Second, cb.NewFunctor00(func() {
		handlerReadConfig()
	}))
}

func handlerPendingUploader() {
	// 清理未决的任务，对应移除未决或校验失败的文件
	CheckPendingUploader()
	task.After(time.Duration(config.Config.PendingTimeout)*time.Second, cb.NewFunctor00(func() {
		handlerPendingUploader()
	}))
}

func handlerExpiredFile() {
	// 清理长期未访问的已上传文件记录
	CheckExpiredFile()
	task.After(time.Duration(config.Config.FileExpiredTimeout)*time.Second, cb.NewFunctor00(func() {
		handlerExpiredFile()
	}))
}
