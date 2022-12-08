package main

import (
	"net/http"
	"os"
	"path/filepath"

	"github.com/cwloo/gonet/core/base/cc"
)

var (
	UseAsyncUploader         = true               //使用异步上传方式
	MaxMemory          int64 = 1024 * 1024 * 1024 //multipart缓存限制
	MaxSegmentSize     int64 = 1024 * 1024 * 20   //单个文件分片上传限制
	MaxSingleSize      int64 = 1024 * 1024 * 1024 //单个文件上传大小限制
	MaxTotalSize       int64 = 1024 * 1024 * 1024 //单次上传文件总大小限制
	PendingTimeout           = 30                 //定期清理未决的上传任务
	FileExpiredTimeout       = 120                //定期清理长期未访问已上传文件记录
)

var (
	ErrOk                  = ErrorMsg{0, "Ok"}                                    //上传完成，并且成功
	ErrSegOk               = ErrorMsg{1, "upload file segment succ"}              //上传成功(分段续传)                       --需要继续分段上传剩余数据
	ErrFileMd5             = ErrorMsg{2, "upload file over, but md5 failed"}      //上传完成，校验出错                       --上传失败
	ErrRepeat              = ErrorMsg{3, "Repeat upload same file"}               //文件重复上传                             --别人上传了
	ErrCheckReUpload       = ErrorMsg{4, "check and re-upload file"}              //文件校正重传                             --需要继续 客户端拿到返回校正数据继续上传
	ErrParamsUUID          = ErrorMsg{5, "upload param error uuid"}               //上传参数错误 uuid                        --上传错误
	ErrParamsMD5           = ErrorMsg{6, "upload param error md5"}                //上传参数错误 文件md5                     --上传错误
	ErrParamsOffset        = ErrorMsg{7, "upload param error offset"}             //上传参数错误 文件已读大小偏移数           --上传错误
	ErrParamsTotalLimit    = ErrorMsg{8, "upload param error total size"}         //上传参数错误 单个上传文件字节数           --上传错误
	ErrParamsSegSizeLimit  = ErrorMsg{9, "upload per-segment size limited"}       //上传参数错误 单次上传字节数限制           --上传错误
	ErrParamsAllTotalLimit = ErrorMsg{10, "upload all total szie limited"}        //上传参数错误 单次上传文件总大小           --上传错误
	ErrParsePartData       = ErrorMsg{11, "parse multipart form-data err"}        //解析multipart form-data数据错误          --上传失败
	ErrParseFormFile       = ErrorMsg{12, "parse multipart form-file err"}        //解析multipart form-file文件错误          --上传失败
	ErrParamsSegSizeZero   = ErrorMsg{13, "upload multipart form-data size zero"} //上传form-data数据字节大小为0             --上传失败
	path, _                = os.Executable()
	dir, exe               = filepath.Split(path)
	dir_upload             = dir + "upload/"
	i32                    = cc.NewI32()
	fileInfos              = NewFileInfos()
	uploaders              = NewSessionToHandler()
)

// <summary>
// ErrorMsg
// <summary>
type ErrorMsg struct {
	ErrCode int    `json:"code,omitempty"`
	ErrMsg  string `json:"errmsg,omitempty"`
}

// <summary>
// Req
// <summary>
type Req struct {
	uuid   string
	keys   []string
	w      http.ResponseWriter
	r      *http.Request
	resp   *Resp
	result []Result
}

// <summary>
// Resp
// <summary>
type Resp struct {
	Uuid    string      `json:"uuid,omitempty"`
	ErrCode int         `json:"code,omitempty"`
	ErrMsg  string      `json:"errmsg,omitempty"`
	Data    interface{} `json:"data,omitempty"`
}

// <summary>
// Result
// <summary>
type Result struct {
	Uuid    string `json:"uuid,omitempty"`
	File    string `json:"file,omitempty"`
	Md5     string `json:"md5,omitempty"`
	Now     int64  `json:"now,omitempty"`
	Total   int64  `json:"total,omitempty"`
	Expired int64  `json:"expired,omitempty"`
	ErrCode int    `json:"code,omitempty"`
	ErrMsg  string `json:"errmsg,omitempty"`
	Message string `json:"message,omitempty"`
}

func Init() {
	UseAsyncUploader = Config.UseAsync > 0
	MaxMemory = Config.MaxMemory
	MaxSegmentSize = Config.MaxSegmentSize
	MaxSingleSize = Config.MaxSingleSize
	MaxTotalSize = Config.MaxTotalSize
	PendingTimeout = Config.PendingTimeout
	FileExpiredTimeout = Config.FileExpiredTimeout
}
