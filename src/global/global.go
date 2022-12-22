package global

import (
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/cwloo/gonet/core/base/cc"
)

var (
	CheckMd5                 = true               //上传完毕是否校验文件完整性
	WriteFile                = false              //上传文件是否缓存服务器本地
	MultiFile                = false              //一次可以上传多个文件
	UseAsyncUploader         = true               //使用异步上传方式
	MaxMemory          int64 = 1024 * 1024 * 1024 //multipart缓存限制
	MaxSegmentSize     int64 = 1024 * 1024 * 20   //单个文件分片上传限制
	MaxSingleSize      int64 = 1024 * 1024 * 1024 //单个文件上传大小限制
	MaxTotalSize       int64 = 1024 * 1024 * 1024 //单次上传文件总大小限制
	PendingTimeout           = 30                 //定期清理未决的上传任务，即前端上传能暂停的最长时间
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
	ErrMultiFileNotSupport = ErrorMsg{14, "upload multifiles not supported"}      //MultiFile为false时，一次只能上传一个文件
	ErrRetry               = ErrorMsg{101, ""}                                    //
	ErrFatal               = ErrorMsg{102, ""}                                    //
	ErrCancel              = ErrorMsg{103, ""}                                    //
	path, _                = os.Executable()                                      //
	Dir, Exe               = filepath.Split(path)                                 //
	Dir_upload             = Dir + "upload/"                                      //上传服务端本地目录，末尾要加上'/'
	I32                    = cc.NewI32()
)

var (
	Name      string
	Cmd       CMD
	Server    Router
	RpcServer RPCServer
)

func init() {
	Cmd.Arg.ID = "id"
	Cmd.Arg.CONF = "config"
	Cmd.Arg.C = "c"
	Cmd.Arg.LOG = "log_dir"
}

// <summary>
// ErrorMsg
// <summary>
type ErrorMsg struct {
	ErrCode int    `json:"code" form:"code"`
	ErrMsg  string `json:"errmsg" form:"errmsg"`
}

func (s *ErrorMsg) Error() string {
	return strings.Join([]string{strconv.Itoa(s.ErrCode), s.ErrMsg}, ":")
}

// <summary>
// Req
// <summary>
type Req struct {
	Uuid   string
	Keys   []*File
	W      http.ResponseWriter
	R      *http.Request
	Resp   *Resp
	Result []Result
}

// <summary>
// File
// <summary>
type File struct {
	Md5        string
	Filename   string
	Headersize int64
	Offset     string
	Total      string
	Key        string
}

// <summary>
// DelReq
// <summary>
type DelReq struct {
	Type  int    `json:"type,omitempty"` // 1-取消文件上传(移除未决的文件) 2-移除已上传的文件
	Md5   string `json:"md5,omitempty"`
	Uuid  string `json:"uuid,omitempty"`
	Check bool   `json:"check,omitempty"` // 是否校验uuid
}

// <summary>
// DelResp
// <summary>
type DelResp struct {
	Type    int    `json:"type,omitempty"`
	Md5     string `json:"md5,omitempty"`
	ErrCode int    `json:"code" form:"code"`
	ErrMsg  string `json:"errmsg" form:"errmsg"`
}

// <summary>
// FileServerReq
// <summary>
type FileServerReq struct {
	Md5 string `json:"md5,omitempty"`
}

// <summary>
// FileServerResp
// <summary>
type FileServerResp struct {
	Md5     string `json:"md5,omitempty"`
	Dns     string `json:"dns,omitempty"`
	ErrCode int    `json:"code" form:"code"`
	ErrMsg  string `json:"errmsg" form:"errmsg"`
}

// <summary>
// FileInfoReq
// <summary>
type FileInfoReq struct {
	Md5 string `json:"md5,omitempty"`
}

// <summary>
// FileInfoResp
// <summary>
type FileInfoResp struct {
	Uuid    string `json:"uuid,omitempty"`
	File    string `json:"file,omitempty"`
	Md5     string `json:"md5,omitempty"`
	Now     int64  `json:"now,omitempty"`
	Total   int64  `json:"total,omitempty"`
	ErrCode int    `json:"code" form:"code"`
	ErrMsg  string `json:"errmsg" form:"errmsg"`
}

// <summary>
// UpdateCfgReq
// <summary>
type UpdateCfgReq struct {
	Interval           string `json:"interval,omitempty"`               //刷新配置间隔时间
	LogTimezone        string `json:"log_timezone" form:"log_timezone"` //
	LogMode            string `json:"log_mode" form:"log_mode"`         //
	LogStyle           string `json:"log_style" form:"log_style"`       //
	LogLevel           string `json:"log_level" form:"log_level"`       //
	MaxMemory          string `json:"maxMemory,omitempty"`              //multipart缓存限制
	MaxSegmentSize     string `json:"maxSegmentSize,omitempty"`         //单个文件分片上传限制
	MaxSingleSize      string `json:"maxSingleSize,omitempty"`          //单个文件上传大小限制
	MaxTotalSize       string `json:"maxTotalSize,omitempty"`           //单次上传文件总大小限制
	PendingTimeout     string `json:"pendingTimeout,omitempty"`         //定期清理未决的上传任务，即前端上传能暂停的最长时间
	FileExpiredTimeout string `json:"fileExpiredTimeout,omitempty"`     //定期清理长期未访问已上传文件记录
	CheckMd5           string `json:"checkMd5,omitempty"`               //上传完毕是否校验文件完整性
	WriteFile          string `json:"writeFile,omitempty"`              //上传文件是否缓存服务器本地
	UseTgBot           string `json:"useTgBot" form:"useTgBot"`
	TgBotChatId        string `json:"tg_chatId" form:"tg_chatId"`
	TgBotToken         string `json:"tg_token" form:"tg_token"`
}

// <summary>
// UpdateCfgResp
// <summary>
type UpdateCfgResp struct {
	ErrCode int    `json:"code" form:"code"`
	ErrMsg  string `json:"errmsg" form:"errmsg"`
}

// <summary>
// GetCfgReq
// <summary>
type GetCfgReq struct {
}

// <summary>
// GetCfgResp
// <summary>
type GetCfgResp struct {
	ErrCode int    `json:"code" form:"code"`
	ErrMsg  string `json:"errmsg" form:"errmsg"`
	Data    any    `json:"data" form:"data"`
}

// <summary>
// CfgData
// <summary>
type CfgData struct {
	Log_dir            string `json:"log_dir" form:"log_dir"`
	Log_level          int    `json:"log_level" form:"log_level"`
	Log_mode           int    `json:"log_mode" form:"log_mode"`
	Log_style          int    `json:"log_style" form:"log_style"`
	Log_timezone       int    `json:"log_timezone" form:"log_timezone"`
	HttpAddr           string `json:"http_addr" form:"http_addr"`
	UploadPath         string `json:"upload_path" form:"upload_path"`
	GetPath            string `json:"get_path" form:"get_path"`
	DelPath            string `json:"delfile_path" form:"delfile_path"`
	FileinfoPath       string `json:"fileinfo_path" form:"fileinfo_path"`
	UpdateCfgPath      string `json:"updatecfg_path" form:"updatecfg_path"`
	GetCfgPath         string `json:"getcfg_path" form:"getcfg_path"`
	CheckMd5           int    `json:"checkMd5" form:"checkMd5"`
	WriteFile          int    `json:"writeFile" form:"writeFile"`
	MultiFile          int    `json:"multiFile" form:"multiFile"`
	UseAsync           int    `json:"useAsync" form:"useAsync"`
	MaxMemory          int64  `json:"maxMemory" form:"maxMemory"`
	MaxSegmentSize     int64  `json:"maxSegmentSize" form:"maxSegmentSize"`
	MaxSingleSize      int64  `json:"maxSingleSize" form:"maxSingleSize"`
	MaxTotalSize       int64  `json:"maxTotalSize" form:"maxTotalSize"`
	PendingTimeout     int    `json:"pendingTimeout" form:"pendingTimeout"`
	FileExpiredTimeout int    `json:"fileExpiredTimeout" form:"fileExpiredTimeout"`
	UploadDir          string `json:"uploadDir" form:"uploadDir"`
	OssType            string `json:"ossType" form:"ossType"`
	UseTgBot           int    `json:"useTgBot" form:"useTgBot"`
	Interval           int    `json:"interval" form:"interval"`
	TgBotChatId        int64  `json:"tg_chatId,omitempty"`
	TgBotToken         string `json:"tg_token,omitempty"`
}

// <summary>
// Resp
// <summary>
type Resp struct {
	Uuid    string `json:"uuid,omitempty"`
	Data    any    `json:"data,omitempty"`
	ErrCode int    `json:"code" form:"code"`
	ErrMsg  string `json:"errmsg" form:"errmsg"`
}

// <summary>
// Result
// <summary>
type Result struct {
	Uuid    string `json:"uuid,omitempty"`
	File    string `json:"file,omitempty"`
	Md5     string `json:"md5,omitempty"`
	Now     int64  `json:"now" form:"now"`
	Total   int64  `json:"total" form:"total"`
	Expired int64  `json:"expired,omitempty"`
	ErrCode int    `json:"code" form:"code"`
	ErrMsg  string `json:"errmsg" form:"errmsg"`
	Message string `json:"message,omitempty"`
	Url     string `json:"url,omitempty"`
}

// <summary>
// FileDetailResp
// <summary>
type FileDetail struct {
	Uuid     string `json:"uuid" form:"uuid"`
	Md5      string `json:"md5" form:"md5"`
	FileName string `json:"filename" form:"filename"`
	DstName  string `json:"dstname" form:"dstname"`
	YunName  string `json:"yunname" form:"yunname"`
	Now      int64  `json:"now" form:"now"`
	Total    int64  `json:"total" form:"total"`
	Url      string `json:"url" form:"url"`
	Create   string `json:"create" form:"create"`
	Time     string `json:"time" form:"time"`
	Elapsed  string `json:"elapsed" form:"elapsed"`
	Percent  string `json:"progress" form:"progress"`
}

// <summary>
// FileDetailReq
// <summary>
type FileDetailReq struct {
	Md5 string `json:"md5,omitempty"`
}

// <summary>
// FileDetailResp
// <summary>
type FileDetailResp struct {
	File    *FileDetail `json:"file" form:"file"`
	ErrCode int         `json:"code" form:"code"`
	ErrMsg  string      `json:"errmsg" form:"errmsg"`
}

// <summary>
// UuidListReq
// <summary>
type UuidListReq struct {
}

// <summary>
// UuidListResp
// <summary>
type UuidListResp struct {
	Uuids   []string `json:"uuids" form:"uuids"`
	ErrCode int      `json:"code" form:"code"`
	ErrMsg  string   `json:"errmsg" form:"errmsg"`
}

// <summary>
// ListReq
// <summary>
type ListReq struct {
}

// <summary>
// ListResp
// <summary>
type ListResp struct {
	Uuids   []string      `json:"uuids" form:"uuids"`
	Files   []*FileDetail `json:"files" form:"files"`
	ErrCode int           `json:"code" form:"code"`
	ErrMsg  string        `json:"errmsg" form:"errmsg"`
}
