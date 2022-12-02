package main

import (
	"time"

	"github.com/cwloo/gonet/core/base/task"
	"github.com/cwloo/gonet/core/cb"
)

func handlerUploadFileOnTimer() {
	checkOnTimer()
	task.After(time.Duration(PendingTimeout)*time.Second, cb.NewFunctor00(func() {
		handlerUploadFileOnTimer()
	}))
}

func checkOnTimer() {
	switch UseAsyncUploader {
	case true:
		uploaders.Range(func(_ string, uploader Uploader) {
			if time.Since(uploader.Get()) > time.Duration(PendingTimeout)*time.Second {
				uploader.NotifyClose()
			}
		})
	default:
		uploaders.RangeRemove(func(uploader Uploader) bool {
			return time.Since(uploader.Get()) > time.Duration(PendingTimeout)*time.Second
		}, func(uploader Uploader) {
			uploader.Clear()
		})
	}
}
