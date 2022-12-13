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
	checkPendingUploader()
	task.After(time.Duration(config.Config.PendingTimeout)*time.Second, cb.NewFunctor00(func() {
		handlerPendingUploader()
	}))
}

func handlerExpiredFile() {
	checkExpiredFile()
	task.After(time.Duration(config.Config.FileExpiredTimeout)*time.Second, cb.NewFunctor00(func() {
		handlerExpiredFile()
	}))
}

// 清理长期未访问的已上传文件记录
func checkExpiredFile() {
	fileInfos.RangeRemoveWithCond(func(info FileInfo) bool {
		if ok, _ := info.Ok(false); ok {
			return time.Since(info.HitTime(false)) >= time.Duration(config.Config.FileExpiredTimeout)*time.Second
		}
		return false
	}, func(info FileInfo) {
		// os.Remove(dir_upload + info.DstName())
		info.Put()
	})
}

// 清理未决的任务，对应移除未决或校验失败的文件
func checkPendingUploader() {
	switch config.Config.UseAsync {
	default:
		////// 异步
		uploaders.Range(func(_ string, uploader Uploader) {
			if time.Since(uploader.Get()) >= time.Duration(config.Config.PendingTimeout)*time.Second {
				uploader.NotifyClose()
			}
		})
	case 0:
		////// 同步
		uploaders.RangeRemoveWithCond(func(uploader Uploader) bool {
			return time.Since(uploader.Get()) >= time.Duration(config.Config.PendingTimeout)*time.Second
		}, func(uploader Uploader) {
			uploader.Clear()
			uploader.Put()
		})
	}
}
