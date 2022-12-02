package main

import (
	"bytes"
	"encoding/json"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/cwloo/gonet/logs"
	"github.com/cwloo/gonet/utils"
)

const (
	BUFSIZ = 1024 * 1024 * 10
)

var (
	path, _  = os.Executable()
	dir, exe = filepath.Split(path)
)
var (
	ErrOk            = ErrorMsg{0, "Ok"}
	ErrSegOk         = ErrorMsg{1, "upload file segment succ"}
	ErrFileMd5       = ErrorMsg{2, "upload file over, but md5 failed"}
	ErrCheckReUpload = ErrorMsg{4, "check and re-upload file"}
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
	ErrCode int      `json:"code,omitempty"`
	ErrMsg  string   `json:"errmsg,omitempty"`
	Data    []Result `json:"data,omitempty"`
}

// <summary>
// Result
// <summary>
type Result struct {
	Uuid    string `json:"uuid,omitempty"`
	Key     string `json:"key,omitempty"`
	File    string `json:"file,omitempty"`
	Md5     string `json:"md5,omitempty"`
	Now     int64  `json:"now,omitempty"`
	Total   int64  `json:"total,omitempty"`
	Expired int64  `json:"expired,omitempty"`
	Result  string `json:"result,omitempty"`
	ErrCode int    `json:"code,omitempty"`
	ErrMsg  string `json:"errmsg,omitempty"`
}

func main() {
	logs.LogTimezone(logs.MY_CST)
	logs.LogInit(dir+"logs", logs.LVL_DEBUG, exe, 100000000)
	logs.LogMode(logs.M_STDOUT_FILE)

	_, filelist := parseargs()
	if len(filelist) == 0 {
		return
	}

	tmp_dir := dir + "tmp" // + fmt.Sprintf(".%v", id)
	os.MkdirAll(tmp_dir, 0666)

	client := httpclient()
	method := "POST"
	url := "http://192.168.1.113:8088/upload"

	uuid := utils.CreateGUID()           //本次上传标识
	MD5 := calcFileMd5(filelist)         //文件md5值
	total, offset := calcFileSize(MD5)   //文件大小/偏移
	results := loadTmpFile(tmp_dir, MD5) //未决临时文件

	start_new := len(results) == 0
CHECKPOINT:
	//////////////////////////////////////////// 先上传未传完的文件 ////////////////////////////////////////////
	for f, md5 := range MD5 {
		if result, ok := results[md5]; ok {
			// 校验文件总字节大小
			if total[md5] != result.Total {
				logs.LogFatal("error")
			}
			// 已经过期，当前文件无法继续上传
			if time.Now().Unix() >= result.Expired {
				os.Remove(tmp_dir + "/" + md5 + ".tmp")
				continue
			}
			// 定位读取文件偏移(上传进度)，从断点处继续上传
			offset := result.Now
			for {
				payload := &bytes.Buffer{}
				writer := multipart.NewWriter(payload)
				_ = writer.WriteField("uuid", result.Uuid)
				// 当前文件没有读完继续
				if result.Total > 0 && offset < result.Total {
					_ = writer.WriteField(md5+".offset", strconv.FormatInt(offset, 10))      //文件偏移量
					_ = writer.WriteField(md5+".total", strconv.FormatInt(result.Total, 10)) //文件总大小
					// 每次断点续传上传 BUFSIZ 字节大小
					part, err := writer.CreateFormFile(md5, filepath.Base(f))
					if err != nil {
						logs.LogFatal("%v", err.Error())
					}
					fd, err := os.OpenFile(f, os.O_RDONLY, 0)
					if err != nil {
						logs.LogFatal("%v", err.Error())
					}
					fd.Seek(offset, io.SeekStart)
					n, err := io.CopyN(part, fd, int64(BUFSIZ))
					if err != nil && err != io.EOF {
						logs.LogFatal("%v", err.Error())
					}
					err = fd.Close()
					if err != nil {
						logs.LogFatal("%v", err.Error())
					}
					err = writer.Close()
					if err != nil {
						logs.LogFatal("%v", err.Error())
					}
					req, err := http.NewRequest(method, url, payload)
					if err != nil {
						logs.LogFatal("%v", err.Error())
					}
					req.Header.Set("Connection", "keep-alive")
					req.Header.Set("Keep-Alive", strings.Join([]string{"timeout=", strconv.Itoa(120)}, ""))
					req.Header.Set("Content-Type", writer.FormDataContentType())
					// logs.LogInfo("user:%v:%v %v %v %v", userId, uuid, method, url, filelist)
					/// request
					res, err := client.Do(req)
					if err != nil {
						logs.LogError("%v", err.Error())
						logs.LogClose()
						return
					}
					defer res.Body.Close()
					for {
						/// response
						body, err := ioutil.ReadAll(res.Body)
						if err != nil {
							logs.LogError("%v", err.Error())
							break
						}
						if len(body) == 0 {
							break
						}
						resp := Resp{}
						err = json.Unmarshal(body, &resp)
						if err != nil {
							logs.LogFatal("%v", err.Error())
							break
						}
						for _, result := range resp.Data {
							switch result.ErrCode {
							case ErrSegOk.ErrCode:
								if result.Now <= 0 {
									break
								}
								// 上传进度写入临时文件
								fd, err := os.OpenFile(tmp_dir+"/"+result.Md5+".tmp", os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0777)
								if err != nil {
									logs.LogError("%v", err.Error())
									return
								}
								b, err := json.Marshal(&result)
								if err != nil {
									logs.LogFatal("%v", err.Error())
									break
								}
								_, err = fd.Write(b)
								if err != nil {
									logs.LogFatal("%v", err.Error())
									break
								}
								err = fd.Close()
								if err != nil {
									logs.LogFatal("%v", err.Error())
								}
							case ErrCheckReUpload.ErrCode:
								//校正需要重传
								if results == nil {
									results = map[string]Result{}
								}
								results[result.Md5] = result
								logs.LogInfo("--- *** --- uuid:%v %v[%v] %v", result.Uuid, result.Md5, result.File, result.ErrMsg)
								goto CHECKPOINT
							case ErrOk.ErrCode, ErrFileMd5.ErrCode:
								//上传完成，删除临时文件
								os.Remove(tmp_dir + "/" + result.Md5 + ".tmp")
								delete(results, result.Md5)
								logs.LogInfo("--- *** --- uuid:%v %v[%v] %v", result.Uuid, result.Md5, result.File, result.ErrMsg)
							}
						}
						// logs.LogInfo("--- *** ---\n%v", string(body))
					}
					if n > 0 {
						offset += n
						if offset == result.Total {
							break
						}
					}
				} // else if offset == result.Total {
				//	break
				//}
			}
		}
	}
	if !start_new {
		return
	}
	//////////////////////////////////////////// 再上传其他文件 ////////////////////////////////////////////
	for {
		finished := true
		// 每次断点续传的payload数据
		payload := &bytes.Buffer{}
		writer := multipart.NewWriter(payload)
		_ = writer.WriteField("uuid", uuid)
		// 要上传的文件列表，各个文件都上传一点
		for f, md5 := range MD5 {
			// 当前文件没有读完继续
			if total[md5] > 0 && offset[md5] < total[md5] {
				finished = false
				_ = writer.WriteField(md5+".offset", strconv.FormatInt(offset[md5], 10)) //文件偏移量
				_ = writer.WriteField(md5+".total", strconv.FormatInt(total[md5], 10))   //文件总大小
				// 每次断点续传上传 BUFSIZ 字节大小
				part, err := writer.CreateFormFile(md5, filepath.Base(f))
				if err != nil {
					logs.LogFatal("%v", err.Error())
				}
				fd, err := os.OpenFile(f, os.O_RDONLY, 0)
				if err != nil {
					logs.LogFatal("%v", err.Error())
				}
				fd.Seek(offset[md5], io.SeekStart)
				n, err := io.CopyN(part, fd, int64(BUFSIZ))
				if err != nil && err != io.EOF {
					logs.LogFatal("%v", err.Error())
				}
				err = fd.Close()
				if err != nil {
					logs.LogFatal("%v", err.Error())
				}
				if n > 0 {
					offset[md5] += n
					continue
				}
			}
		}
		err := writer.Close()
		if err != nil {
			logs.LogFatal("%v", err.Error())
		}
		if !finished {
			req, err := http.NewRequest(method, url, payload)
			if err != nil {
				logs.LogFatal("%v", err.Error())
			}
			req.Header.Set("Connection", "keep-alive")
			req.Header.Set("Keep-Alive", strings.Join([]string{"timeout=", strconv.Itoa(120)}, ""))
			req.Header.Set("Content-Type", writer.FormDataContentType())
			// logs.LogInfo("user:%v:%v %v %v %v", userId, uuid, method, url, filelist)
			/// request
			res, err := client.Do(req)
			if err != nil {
				logs.LogError("%v", err.Error())
				logs.LogClose()
				return
			}
			defer res.Body.Close()
			for {
				/// response
				body, err := ioutil.ReadAll(res.Body)
				if err != nil {
					logs.LogError("%v", err.Error())
					break
				}
				if len(body) == 0 {
					break
				}
				resp := Resp{}
				err = json.Unmarshal(body, &resp)
				if err != nil {
					logs.LogFatal("%v", err.Error())
					break
				}
				for _, result := range resp.Data {
					switch result.ErrCode {
					case ErrSegOk.ErrCode:
						if result.Now <= 0 {
							break
						}
						// 上传进度写入临时文件
						fd, err := os.OpenFile(tmp_dir+"/"+result.Md5+".tmp", os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0777)
						if err != nil {
							logs.LogError("%v", err.Error())
							break
						}
						b, err := json.Marshal(&result)
						if err != nil {
							logs.LogFatal("%v", err.Error())
							break
						}
						_, err = fd.Write(b)
						if err != nil {
							logs.LogFatal("%v", err.Error())
							break
						}
						err = fd.Close()
						if err != nil {
							logs.LogFatal("%v", err.Error())
						}
					case ErrCheckReUpload.ErrCode:
						//校正需要重传
						if results == nil {
							results = map[string]Result{}
						}
						results[result.Md5] = result
						logs.LogInfo("--- --- --- uuid:%v %v[%v] %v", result.Uuid, result.Md5, result.File, result.ErrMsg)
						goto CHECKPOINT
					case ErrOk.ErrCode, ErrFileMd5.ErrCode:
						//上传完成，删除临时文件
						os.Remove(tmp_dir + "/" + result.Md5 + ".tmp")
						delete(results, result.Md5)
						logs.LogInfo("--- --- --- uuid:%v %v[%v] %v", result.Uuid, result.Md5, result.File, result.ErrMsg)
					}
				}
				// logs.LogInfo("--- --- ---\n%v", string(body))
			}
		} else {
			break
		}
	}
	logs.LogClose()
}
