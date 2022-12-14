package main

import (
	"time"

	"github.com/cwloo/uploader/file_server/global"
)

// <summary>
// Uploader
// <summary>
type Uploader interface {
	Get() time.Time
	Upload(req *global.Req)
	Remove(md5 string)
	Clear()
	Close()
	NotifyClose()
	Put()
}
