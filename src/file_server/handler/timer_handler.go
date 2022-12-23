package handler

import (
	"time"

	"github.com/cwloo/gonet/core/base/task"
	"github.com/cwloo/gonet/core/cb"
	"github.com/cwloo/uploader/src/config"
	"github.com/cwloo/uploader/src/global"
)

func ReadConfig() {
	config.ReadConfig(global.Cmd.Conf)
	task.After(time.Duration(config.Config.Interval)*time.Second, cb.NewFunctor00(func() {
		ReadConfig()
	}))
}

func PendingUploader() {
	// 清理未决的任务，对应移除未决或校验失败的文件
	CheckPendingUploader()
	task.After(time.Duration(config.Config.File.Upload.PendingTimeout)*time.Second, cb.NewFunctor00(func() {
		PendingUploader()
	}))
}

func ExpiredFile() {
	// 清理长期未访问的已上传文件记录
	CheckExpiredFile()
	task.After(time.Duration(config.Config.File.Upload.FileExpiredTimeout)*time.Second, cb.NewFunctor00(func() {
		ExpiredFile()
	}))
}
