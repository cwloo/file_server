package main

import (
	"io/ioutil"
	"os"
	"time"

	"github.com/cwloo/gonet/logs"
	"github.com/cwloo/gonet/utils"
)

// <summary>
// Uploader
// <summary>
type Uploader interface {
	Get() time.Time
	Upload(req *Req)
	Clear()
	Close()
	NotifyClose()
	Put()
}

func calFileMd5(f string) string {
	fd, err := os.OpenFile(f, os.O_RDONLY, 0)
	if err != nil {
		logs.LogFatal(err.Error())
	}
	b, err := ioutil.ReadAll(fd)
	if err != nil {
		logs.LogFatal(err.Error())
	}
	err = fd.Close()
	if err != nil {
		logs.LogFatal(err.Error())
	}
	return utils.MD5Byte(b, false)
}
