package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"time"

	"github.com/cwloo/gonet/logs"
	"github.com/cwloo/gonet/utils"
	"github.com/cwloo/uploader/file_server/config"
	"github.com/cwloo/uploader/file_server/global"
)

func CalcFileMd5(f string) string {
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

func UpdateCfg(req *global.UpdateCfgReq) (*global.UpdateCfgResp, bool) {
	config.UpdateConfig(req)
	return &global.UpdateCfgResp{
		ErrCode: 0,
		ErrMsg:  "ok"}, true
}

func QueryFileinfoCache(md5 string) (*global.FileInfoResp, bool) {
	info := fileInfos.Get(md5)
	if info == nil {
		return &global.FileInfoResp{Md5: md5, ErrCode: 5, ErrMsg: "not found"}, false
	}
	return &global.FileInfoResp{
		Uuid:    info.Uuid(),
		File:    info.SrcName(),
		Md5:     md5,
		Now:     info.Now(false),
		Total:   info.Total(false),
		ErrCode: 0,
		ErrMsg:  "ok"}, true
}

func DelCacheFile(delType int, md5 string) {
	switch delType {
	case 1:
		// 1-取消文件上传(移除未决的文件)
		fileInfos.RemoveWithCond(md5, func(info FileInfo) bool {
			return !info.Done(false)
		}, func(info FileInfo) {
			os.Remove(config.Config.UploadlDir + info.DstName())
			uploaders.Get(info.Uuid()).Remove(md5)
			info.Put()
		})
	case 2:
		// 2-移除已上传的文件
		fileInfos.RemoveWithCond(md5, func(info FileInfo) bool {
			if ok, _ := info.Ok(false); ok {
				return true
			}
			return false
		}, func(info FileInfo) {
			os.Remove(config.Config.UploadlDir + info.DstName())
			info.Put()
		})
	}
}

func RemovePendingFile(uuid, md5 string) (msg string, ok bool) {
	fileInfos.RemoveWithCond(md5, func(info FileInfo) bool {
		if info.Uuid() != uuid {
			logs.LogFatal("error")
		}
		if info.Done(false) {
			logs.LogFatal("error")
		}
		return true
	}, func(info FileInfo) {
		msg = fmt.Sprintf("%v\n%v[%v]\n%v [Err]", info.Uuid(), info.SrcName(), md5, info.DstName())
		os.Remove(config.Config.UploadlDir + info.DstName())
		info.Put()
	})
	ok = msg != ""
	return
}

func RemoveCheckErrFile(uuid, md5 string) (msg string, ok bool) {
	fileInfos.RemoveWithCond(md5, func(info FileInfo) bool {
		if info.Uuid() != uuid {
			logs.LogFatal("error")
		}
		if !info.Done(false) {
			logs.LogFatal("error")
		}
		ok, _ := info.Ok(false)
		return !ok
	}, func(info FileInfo) {
		msg = fmt.Sprintf("%v\n%v[%v]\n%v chkmd5 [Err]", info.Uuid(), info.SrcName(), md5, info.DstName())
		os.Remove(config.Config.UploadlDir + info.DstName())
		info.Put()
	})
	ok = msg != ""
	return
}

func CheckExpiredFile() {
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

func CheckPendingUploader() {
	switch config.Config.UseAsync > 0 {
	case true:
		////// 异步
		uploaders.Range(func(_ string, uploader Uploader) {
			if time.Since(uploader.Get()) >= time.Duration(config.Config.PendingTimeout)*time.Second {
				uploader.NotifyClose()
			}
		})
	default:
		////// 同步
		uploaders.RangeRemoveWithCond(func(uploader Uploader) bool {
			return time.Since(uploader.Get()) >= time.Duration(config.Config.PendingTimeout)*time.Second
		}, func(uploader Uploader) {
			uploader.Clear()
			uploader.Put()
		})
	}
}
