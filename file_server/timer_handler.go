package main

import (
	"time"

	"github.com/cwloo/gonet/core/base/task"
	"github.com/cwloo/gonet/core/cb"
)

func handlerPendingUploader() {
	checkPendingUploader()
	task.After(time.Duration(PendingTimeout)*time.Second, cb.NewFunctor00(func() {
		handlerPendingUploader()
	}))
}

func handlerExpiredFile() {
	checkExpiredFile()
	task.After(time.Duration(FileExpiredTimeout)*time.Second, cb.NewFunctor00(func() {
		handlerExpiredFile()
	}))
}

// 清理长期未访问的已上传文件记录
func checkExpiredFile() {
	fileInfos.RangeRemoveWithCond(func(info FileInfo) bool {
		if ok, _ := info.Ok(); ok {
			return time.Since(info.HitTime()) >= time.Duration(FileExpiredTimeout)*time.Second
		}
		return false
	}, func(info FileInfo) {
		// os.Remove(dir_upload + info.DstName())
		info.Put()
	})
}

// 清理未决的任务，对应移除未决或校验失败的文件
func checkPendingUploader() {
	switch UseAsyncUploader {
	case true:
		////// 异步
		uploaders.Range(func(_ string, uploader Uploader) {
			if time.Since(uploader.Get()) >= time.Duration(PendingTimeout)*time.Second {
				uploader.NotifyClose()
			}
		})
	default:
		////// 同步
		uploaders.RangeRemoveWithCond(func(uploader Uploader) bool {
			return time.Since(uploader.Get()) >= time.Duration(PendingTimeout)*time.Second
		}, func(uploader Uploader) {
			uploader.Clear()
			uploader.Put()
		})
	}
}
