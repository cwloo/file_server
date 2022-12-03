package main

import (
	"os"
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
	task.After(time.Duration(PendingTimeout)*time.Second, cb.NewFunctor00(func() {
		handlerExpiredFile()
	}))
}

func checkExpiredFile() {
	////// 清理长期未访问的已上传文件记录
	fileInfos.RangeRemoveWithCond(func(info *FileInfo) bool {
		if info.Ok() && info.Md5Ok {
			return time.Since(info.HitTime()) >= time.Duration(FileExpiredTimeout)*time.Second
		}
		return false
	}, func(info *FileInfo) {
		os.Remove(dir_upload + info.DstName)
	})
}

func checkPendingUploader() {
	switch UseAsyncUploader {
	case true:
		////// 清理未决的任务，对应移除未决或校验失败的文件(异步)
		uploaders.Range(func(_ string, uploader Uploader) {
			if time.Since(uploader.Get()) > time.Duration(PendingTimeout)*time.Second {
				uploader.NotifyClose()
			}
		})
	default:
		////// 清理未决的任务，对应移除未决或校验失败的文件(同步)
		uploaders.RangeRemoveWithCond(func(uploader Uploader) bool {
			return time.Since(uploader.Get()) > time.Duration(PendingTimeout)*time.Second
		}, func(uploader Uploader) {
			uploader.Clear()
		})
	}
}
