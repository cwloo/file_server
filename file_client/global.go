package main

import (
	"os"
	"path/filepath"
)

var (
	MultiFile         = false            //一次可以上传多个文件
	SegmentSize int64 = 1024 * 1024 * 10 //单个文件分片上传大小
)

var (
	path, _  = os.Executable()
	dir, exe = filepath.Split(path)
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
	ErrMultiFileNotSupport = ErrorMsg{14, "upload multifiles not supported"}      //MultiFile为false时，一次只能上传一个文件
)

// <summary>
// ErrorMsg
// <summary>
type ErrorMsg struct {
	ErrCode int    `json:"code,omitempty"`
	ErrMsg  string `json:"errmsg,omitempty"`
}

// <summary>
// Resp
// <summary>
type Resp struct {
	Uuid    string   `json:"uuid,omitempty"`
	ErrCode int      `json:"code,omitempty"`
	ErrMsg  string   `json:"errmsg,omitempty"`
	Data    []Result `json:"data,omitempty"`
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
	Url     string `json:"url,omitempty"`
}

func Init() {
	SegmentSize = Config.SegmentSize
	MultiFile = Config.MultiFile > 0
}
