package main

import (
	"net/http"
	"os"
	"path/filepath"

	"github.com/cwloo/gonet/core/base/cc"
)

const (
	MaxMemory       int64 = 1024 * 1024 * 1024
	MaxTotalSize    int64 = 1024 * 1024 * 1024 //单个文件上传不超过1G
	MaxSegmentSize  int64 = 1024 * 1024 * 20   //单个文件断点续传字节数限制
	MaxAllTotalSize int64 = 1024 * 1024 * 1024 //单次上传文件总大小字节数限制
	PendingTimeout        = 10                 //间隔秒数检查未决的上传任务
)

var (
	ErrOk                  = ErrorMsg{0, "Ok"}                                    //上传完成，并且成功
	ErrSegOk               = ErrorMsg{1, "upload file segment succ"}              //上传分段成功
	ErrFileMd5             = ErrorMsg{2, "upload file over, but md5 failed"}      //上传完成，文件出错
	ErrRepeat              = ErrorMsg{3, "Repeat upload same file"}               //文件重复上传
	ErrParamsUUID          = ErrorMsg{4, "upload param error uuid"}               //上传参数错误 uuid
	ErrParamsMD5           = ErrorMsg{5, "upload param error md5"}                //上传参数错误 文件md5
	ErrParamsTotalLimit    = ErrorMsg{6, "upload param error total size"}         //上传参数错误 单个上传文件字节数
	ErrParamsSegSizeLimit  = ErrorMsg{7, "upload per-segment size limited"}       //上传参数错误 单次上传字节数限制
	ErrParamsAllTotalLimit = ErrorMsg{8, "upload all total szie limited"}         //上传参数错误 单次上传文件总大小
	ErrParsePartData       = ErrorMsg{9, "parse multipart form-data err"}         //解析multipart form-data数据错误
	ErrParseFormFile       = ErrorMsg{9, "parse multipart form-file err"}         //解析multipart form-file文件错误
	ErrParamsSegSizeZero   = ErrorMsg{10, "upload multipart form-data size zero"} //上传form-data数据字节大小为0
	path, _                = os.Executable()
	dir, _                 = filepath.Split(path)
	i32                    = cc.NewI32()
	fileInfos              = NewFileInfos()
	uploaders              = NewSessionToHandler()
)

// <summary>
// ErrorMsg
// <summary>
type ErrorMsg struct {
	ErrCode int
	ErrMsg  string
}

// <summary>
// Req
// <summary>
type Req struct {
	uuid   string
	keys   []string
	ignore []*FileInfo
	w      http.ResponseWriter
	r      *http.Request
}

// <summary>
// Resp
// <summary>
type Resp struct {
	ErrCode int         `json:"code,omitempty"`
	ErrMsg  string      `json:"errmsg,omitempty"`
	Data    interface{} `json:"data,omitempty"`
}

// <summary>
// Result
// <summary>
type Result struct {
	Uuid    string
	File    string
	Md5     string
	Result  string
	ErrCode int    `json:"code,omitempty"`
	ErrMsg  string `json:"errmsg,omitempty"`
}
