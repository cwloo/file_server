package global

import (
	"io/ioutil"
	"os"
	"time"

	"github.com/cwloo/gonet/logs"
	"github.com/cwloo/gonet/utils"
)

type NewUploader func(bool, string) Uploader

// <summary>
// Uploader
// <summary>
type Uploader interface {
	Get() time.Time
	Upload(req *Req)
	Remove(md5 string)
	Clear()
	Close()
	NotifyClose()
	Put()
}

func CalcFileMd5(f string) string {
	fd, err := os.OpenFile(f, os.O_RDONLY, 0)
	if err != nil {
		logs.Fatalf(err.Error())
	}
	b, err := ioutil.ReadAll(fd)
	if err != nil {
		logs.Fatalf(err.Error())
	}
	err = fd.Close()
	if err != nil {
		logs.Fatalf(err.Error())
	}
	return utils.MD5Byte(b, false)
}
