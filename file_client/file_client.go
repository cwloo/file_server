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

// 一次只能上传一个文件
func upload() {

	_, filelist := parseargs()
	if len(filelist) == 0 {
		return
	}

	tmp_dir := dir + "tmp/"
	os.MkdirAll(tmp_dir, 0666)

	client := httpclient()
	method := "POST"
	url := strings.Join([]string{Config.HttpProto, Config.HttpAddr, Config.UploadPath}, "")

	uuid := utils.CreateGUID()           //本次上传标识
	MD5 := calcFileMd5(filelist)         //文件md5值
	total, offset := calcFileSize(MD5)   //文件大小/偏移
	results := loadTmpFile(tmp_dir, MD5) //未决临时文件
	//////////////////////////////////////////// 先上传未传完的文件 ////////////////////////////////////////////
	for {
		if len(results) == 0 {
			break
		}
		for md5, result := range results {
			f := filePathBy(&MD5, md5)
			// 校验文件总字节大小
			if total[md5] != result.Total {
				logs.LogFatal("error")
			}
			// 已经过期，当前文件无法继续上传
			if time.Now().Unix() >= result.Expired {
				delete(results, md5)
				os.Remove(tmp_dir + md5 + ".tmp")
				continue
			}
			// 定位读取文件偏移(上传进度)，从断点处继续上传
			offset[result.Md5] = result.Now
			for {
				// 当前文件没有读完继续
				if offset[result.Md5] < result.Total {
					payload := &bytes.Buffer{}
					writer := multipart.NewWriter(payload)
					_ = writer.WriteField("uuid", result.Uuid)
					_ = writer.WriteField("md5", result.Md5)
					_ = writer.WriteField("offset", strconv.FormatInt(offset[result.Md5], 10)) //文件偏移量
					_ = writer.WriteField("total", strconv.FormatInt(result.Total, 10))        //文件总大小
					part, err := writer.CreateFormFile("file", filepath.Base(f))
					if err != nil {
						logs.LogFatal(err.Error())
					}
					fd, err := os.OpenFile(f, os.O_RDONLY, 0)
					if err != nil {
						logs.LogFatal(err.Error())
					}
					// 单个文件分片上传大小
					fd.Seek(offset[result.Md5], io.SeekStart)
					_, err = io.CopyN(part, fd, int64(SegmentSize))
					if err != nil && err != io.EOF {
						logs.LogFatal(err.Error())
					}
					err = fd.Close()
					if err != nil {
						logs.LogFatal(err.Error())
					}
					err = writer.Close()
					if err != nil {
						logs.LogFatal(err.Error())
					}
					req, err := http.NewRequest(method, url, payload)
					if err != nil {
						logs.LogFatal(err.Error())
					}
					req.Header.Set("Connection", "keep-alive")
					req.Header.Set("Keep-Alive", strings.Join([]string{"timeout=", strconv.Itoa(120)}, ""))
					req.Header.Set("Content-Type", writer.FormDataContentType())
					logs.LogInfo("request =>> %v %v[%v] %v %v", method, url, 1, result.File, uuid)
					/// request
					res, err := client.Do(req)
					if err != nil {
						logs.LogError(err.Error())
						logs.LogClose()
						return
					}
					for {
						/// response
						body, err := ioutil.ReadAll(res.Body)
						if err != nil {
							logs.LogError(err.Error())
							break
						}
						if len(body) == 0 {
							break
						}
						resp := Resp{}
						err = json.Unmarshal(body, &resp)
						if err != nil {
							logs.LogError(err.Error())
							logs.LogWarn("%v", string(body))
							continue
						}
						// 检查有无 resp 错误码
						switch resp.ErrCode {
						case ErrMultiFileNotSupport.ErrCode:
							fallthrough
						case ErrParamsUUID.ErrCode:
							fallthrough
						case ErrParsePartData.ErrCode:
							// 需要继续重试
							logs.LogError("*** %v %v", resp.Uuid, resp.ErrMsg)
							continue
						}
						// 读取每个文件上传状态数据
						for _, result := range resp.Data {
							switch result.ErrCode {
							case ErrParseFormFile.ErrCode:
								fallthrough
							case ErrParamsSegSizeLimit.ErrCode:
								fallthrough
							case ErrParamsSegSizeZero.ErrCode:
								fallthrough
							case ErrParamsTotalLimit.ErrCode:
								fallthrough
							case ErrParamsOffset.ErrCode:
								fallthrough
							case ErrParamsMD5.ErrCode:
								fallthrough
							case ErrParamsAllTotalLimit.ErrCode:
								fallthrough
							case ErrRepeat.ErrCode:
								logs.LogError("*** %v %v[%v] %v => %v", result.Uuid, result.Md5, result.File, result.ErrMsg, result.Message)
							case ErrSegOk.ErrCode:
								logs.LogError("*** %v %v[%v] %v", result.Uuid, result.Md5, result.File, result.ErrMsg)
								if result.Now <= 0 {
									break
								}
								// 上传进度写入临时文件
								fd, err := os.OpenFile(tmp_dir+result.Md5+".tmp", os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0777)
								if err != nil {
									logs.LogError(err.Error())
									return
								}
								b, err := json.Marshal(&result)
								if err != nil {
									logs.LogFatal(err.Error())
									break
								}
								_, err = fd.Write(b)
								if err != nil {
									logs.LogFatal(err.Error())
									break
								}
								err = fd.Close()
								if err != nil {
									logs.LogFatal(err.Error())
								}
								// 更新文件读取偏移
								offset[result.Md5] = result.Now
							case ErrCheckReUpload.ErrCode:
								// 校正需要重传
								if results == nil {
									results = map[string]Result{}
								}
								results[result.Md5] = result
								offset[result.Md5] = result.Now
								logs.LogError("*** %v %v[%v] %v", result.Uuid, result.Md5, result.File, result.ErrMsg)
							case ErrFileMd5.ErrCode:
								// 上传失败了
								fallthrough
							case ErrOk.ErrCode:
								delete(results, result.Md5)
								offset[result.Md5] = total[result.Md5]
								removeMd5File(&MD5, result.Md5)
								// 上传完成，删除临时文件
								os.Remove(tmp_dir + result.Md5 + ".tmp")
								logs.LogTrace("*** %v %v[%v] %v => %v", result.Uuid, result.Md5, result.File, result.ErrMsg, result.Url)
							}
						}
					}
					res.Body.Close()
				}
			}
		}
	}
	//////////////////////////////////////////// 再上传其他文件 ////////////////////////////////////////////
	for {
		finished := true
		Filelist := []string{}
		// 每次断点续传的payload数据
		payload := &bytes.Buffer{}
		writer := multipart.NewWriter(payload)
		_ = writer.WriteField("uuid", uuid)
		// 要上传的文件列表，各个文件都上传一点
		for f, md5 := range MD5 {
			// 当前文件没有读完继续
			if offset[md5] < total[md5] {
				finished = false
				Filelist = append(Filelist, filepath.Base(f))
				_ = writer.WriteField("md5", md5)
				_ = writer.WriteField("offset", strconv.FormatInt(offset[md5], 10)) //文件偏移量
				_ = writer.WriteField("total", strconv.FormatInt(total[md5], 10))   //文件总大小
				part, err := writer.CreateFormFile("file", filepath.Base(f))
				if err != nil {
					logs.LogFatal(err.Error())
				}
				fd, err := os.OpenFile(f, os.O_RDONLY, 0)
				if err != nil {
					logs.LogFatal(err.Error())
				}
				// 单个文件分片上传大小
				fd.Seek(offset[md5], io.SeekStart)
				_, err = io.CopyN(part, fd, int64(SegmentSize))
				if err != nil && err != io.EOF {
					logs.LogFatal(err.Error())
				}
				err = fd.Close()
				if err != nil {
					logs.LogFatal(err.Error())
				}
			}
			break
		}
		err := writer.Close()
		if err != nil {
			logs.LogFatal(err.Error())
		}
		if !finished {
		retry:
			req, err := http.NewRequest(method, url, payload)
			if err != nil {
				logs.LogFatal(err.Error())
			}
			req.Header.Set("Connection", "keep-alive")
			req.Header.Set("Keep-Alive", strings.Join([]string{"timeout=", strconv.Itoa(120)}, ""))
			req.Header.Set("Content-Type", writer.FormDataContentType())
			logs.LogInfo("request =>> %v %v[%v] %v %v", method, url, len(Filelist), Filelist, uuid)
			/// request
			res, err := client.Do(req)
			if err != nil {
				logs.LogError(err.Error())
				// logs.LogClose()
				goto retry
			}
			for {
				/// response
				body, err := ioutil.ReadAll(res.Body)
				if err != nil {
					logs.LogError(err.Error())
					break
				}
				if len(body) == 0 {
					break
				}
				resp := Resp{}
				err = json.Unmarshal(body, &resp)
				if err != nil {
					logs.LogError(err.Error())
					logs.LogWarn("%v", string(body))
					continue
				}
				// 检查有无 resp 错误码
				switch resp.ErrCode {
				case ErrMultiFileNotSupport.ErrCode:
					fallthrough
				case ErrParamsUUID.ErrCode:
					fallthrough
				case ErrParsePartData.ErrCode:
					// 需要继续重试
					logs.LogError("--- %v %v", resp.Uuid, resp.ErrMsg)
					continue
				}
				// 读取每个文件上传状态数据
				for _, result := range resp.Data {
					switch result.ErrCode {
					case ErrParseFormFile.ErrCode:
						fallthrough
					case ErrParamsSegSizeLimit.ErrCode:
						fallthrough
					case ErrParamsSegSizeZero.ErrCode:
						fallthrough
					case ErrParamsTotalLimit.ErrCode:
						fallthrough
					case ErrParamsOffset.ErrCode:
						fallthrough
					case ErrParamsMD5.ErrCode:
						fallthrough
					case ErrParamsAllTotalLimit.ErrCode:
						fallthrough
					case ErrRepeat.ErrCode:
						logs.LogError("--- %v %v[%v] %v => %v", result.Uuid, result.Md5, result.File, result.ErrMsg, result.Message)
					// 上传成功(分段续传)，继续读取文件剩余字节继续上传
					case ErrSegOk.ErrCode:
						logs.LogError("--- %v %v[%v] %v", result.Uuid, result.Md5, result.File, result.ErrMsg)
						if result.Now <= 0 {
							break
						}
						// 上传进度写入临时文件
						fd, err := os.OpenFile(tmp_dir+result.Md5+".tmp", os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0777)
						if err != nil {
							logs.LogError(err.Error())
							break
						}
						b, err := json.Marshal(&result)
						if err != nil {
							logs.LogFatal(err.Error())
							break
						}
						_, err = fd.Write(b)
						if err != nil {
							logs.LogFatal(err.Error())
							break
						}
						err = fd.Close()
						if err != nil {
							logs.LogFatal(err.Error())
						}
						// 更新文件读取偏移
						offset[result.Md5] = result.Now
					case ErrCheckReUpload.ErrCode:
						// 校正需要重传
						offset[result.Md5] = result.Now
						logs.LogError("--- %v %v[%v] %v => %v", result.Uuid, result.Md5, result.File, result.ErrMsg, result.Message)
					case ErrFileMd5.ErrCode:
						// 上传失败了
						fallthrough
					case ErrOk.ErrCode:
						offset[result.Md5] = total[result.Md5]
						removeMd5File(&MD5, result.Md5)
						// 上传完成，删除临时文件
						os.Remove(tmp_dir + result.Md5 + ".tmp")
						logs.LogTrace("--- %v %v[%v] %v => %v", result.Uuid, result.Md5, result.File, result.ErrMsg, result.Url)
					}
				}
			}
			res.Body.Close()
		} else {
			break
		}
	}
	logs.LogClose()
}

// 一次可以上传多个文件
func multiUpload() {

	_, filelist := parseargs()
	if len(filelist) == 0 {
		return
	}

	tmp_dir := dir + "tmp/"
	os.MkdirAll(tmp_dir, 0666)

	client := httpclient()
	method := "POST"
	url := strings.Join([]string{Config.HttpProto, Config.HttpAddr, Config.UploadPath}, "")

	uuid := utils.CreateGUID()           //本次上传标识
	MD5 := calcFileMd5(filelist)         //文件md5值
	total, offset := calcFileSize(MD5)   //文件大小/偏移
	results := loadTmpFile(tmp_dir, MD5) //未决临时文件
	//////////////////////////////////////////// 先上传未传完的文件 ////////////////////////////////////////////
	for {
		if len(results) == 0 {
			break
		}
		for md5, result := range results {
			f := filePathBy(&MD5, md5)
			// 校验文件总字节大小
			if total[md5] != result.Total {
				logs.LogFatal("error")
			}
			// 已经过期，当前文件无法继续上传
			if time.Now().Unix() >= result.Expired {
				delete(results, md5)
				os.Remove(tmp_dir + md5 + ".tmp")
				continue
			}
			// 定位读取文件偏移(上传进度)，从断点处继续上传
			offset[result.Md5] = result.Now
			for {
				// 当前文件没有读完继续
				if offset[result.Md5] < result.Total {
					payload := &bytes.Buffer{}
					writer := multipart.NewWriter(payload)
					_ = writer.WriteField("uuid", result.Uuid)
					_ = writer.WriteField(md5+".offset", strconv.FormatInt(offset[result.Md5], 10)) //文件偏移量
					_ = writer.WriteField(md5+".total", strconv.FormatInt(result.Total, 10))        //文件总大小
					part, err := writer.CreateFormFile(md5, filepath.Base(f))
					if err != nil {
						logs.LogFatal(err.Error())
					}
					fd, err := os.OpenFile(f, os.O_RDONLY, 0)
					if err != nil {
						logs.LogFatal(err.Error())
					}
					// 单个文件分片上传大小
					fd.Seek(offset[result.Md5], io.SeekStart)
					_, err = io.CopyN(part, fd, int64(SegmentSize))
					if err != nil && err != io.EOF {
						logs.LogFatal(err.Error())
					}
					err = fd.Close()
					if err != nil {
						logs.LogFatal(err.Error())
					}
					err = writer.Close()
					if err != nil {
						logs.LogFatal(err.Error())
					}
					req, err := http.NewRequest(method, url, payload)
					if err != nil {
						logs.LogFatal(err.Error())
					}
					req.Header.Set("Connection", "keep-alive")
					req.Header.Set("Keep-Alive", strings.Join([]string{"timeout=", strconv.Itoa(120)}, ""))
					req.Header.Set("Content-Type", writer.FormDataContentType())
					logs.LogInfo("request =>> %v %v[%v] %v %v", method, url, 1, result.File, uuid)
					/// request
					res, err := client.Do(req)
					if err != nil {
						logs.LogError(err.Error())
						logs.LogClose()
						return
					}
					for {
						/// response
						body, err := ioutil.ReadAll(res.Body)
						if err != nil {
							logs.LogError(err.Error())
							break
						}
						if len(body) == 0 {
							break
						}
						resp := Resp{}
						err = json.Unmarshal(body, &resp)
						if err != nil {
							logs.LogError(err.Error())
							logs.LogWarn("%v", string(body))
							continue
						}
						// 检查有无 resp 错误码
						switch resp.ErrCode {
						case ErrParamsUUID.ErrCode:
							fallthrough
						case ErrParsePartData.ErrCode:
							// 需要继续重试
							logs.LogError("*** %v %v", resp.Uuid, resp.ErrMsg)
							continue
						}
						// 读取每个文件上传状态数据
						for _, result := range resp.Data {
							switch result.ErrCode {
							case ErrParseFormFile.ErrCode:
								fallthrough
							case ErrParamsSegSizeLimit.ErrCode:
								fallthrough
							case ErrParamsSegSizeZero.ErrCode:
								fallthrough
							case ErrParamsTotalLimit.ErrCode:
								fallthrough
							case ErrParamsOffset.ErrCode:
								fallthrough
							case ErrParamsMD5.ErrCode:
								fallthrough
							case ErrParamsAllTotalLimit.ErrCode:
								fallthrough
							case ErrRepeat.ErrCode:
								logs.LogError("*** %v %v[%v] %v => %v", result.Uuid, result.Md5, result.File, result.ErrMsg, result.Message)
							case ErrSegOk.ErrCode:
								logs.LogError("*** %v %v[%v] %v", result.Uuid, result.Md5, result.File, result.ErrMsg)
								if result.Now <= 0 {
									break
								}
								// 上传进度写入临时文件
								fd, err := os.OpenFile(tmp_dir+result.Md5+".tmp", os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0777)
								if err != nil {
									logs.LogError(err.Error())
									return
								}
								b, err := json.Marshal(&result)
								if err != nil {
									logs.LogFatal(err.Error())
									break
								}
								_, err = fd.Write(b)
								if err != nil {
									logs.LogFatal(err.Error())
									break
								}
								err = fd.Close()
								if err != nil {
									logs.LogFatal(err.Error())
								}
								// 更新文件读取偏移
								offset[result.Md5] = result.Now
							case ErrCheckReUpload.ErrCode:
								// 校正需要重传
								if results == nil {
									results = map[string]Result{}
								}
								results[result.Md5] = result
								offset[result.Md5] = result.Now
								logs.LogError("*** %v %v[%v] %v => %v", result.Uuid, result.Md5, result.File, result.ErrMsg, result.Message)
							case ErrFileMd5.ErrCode:
								// 上传失败了
								fallthrough
							case ErrOk.ErrCode:
								delete(results, result.Md5)
								offset[result.Md5] = total[result.Md5]
								removeMd5File(&MD5, result.Md5)
								// 上传完成，删除临时文件
								os.Remove(tmp_dir + result.Md5 + ".tmp")
								logs.LogTrace("*** %v %v[%v] %v => %v", result.Uuid, result.Md5, result.File, result.ErrMsg, result.Url)
							}
						}
					}
					res.Body.Close()
				}
			}
		}
	}
	//////////////////////////////////////////// 再上传其他文件 ////////////////////////////////////////////
	for {
		finished := true
		Filelist := []string{}
		// 每次断点续传的payload数据
		payload := &bytes.Buffer{}
		writer := multipart.NewWriter(payload)
		_ = writer.WriteField("uuid", uuid)
		// 要上传的文件列表，各个文件都上传一点
		for f, md5 := range MD5 {
			// 当前文件没有读完继续
			if offset[md5] < total[md5] {
				finished = false
				Filelist = append(Filelist, filepath.Base(f))
				_ = writer.WriteField(md5+".offset", strconv.FormatInt(offset[md5], 10)) //文件偏移量
				_ = writer.WriteField(md5+".total", strconv.FormatInt(total[md5], 10))   //文件总大小
				part, err := writer.CreateFormFile(md5, filepath.Base(f))
				if err != nil {
					logs.LogFatal(err.Error())
				}
				fd, err := os.OpenFile(f, os.O_RDONLY, 0)
				if err != nil {
					logs.LogFatal(err.Error())
				}
				// 单个文件分片上传大小
				fd.Seek(offset[md5], io.SeekStart)
				_, err = io.CopyN(part, fd, int64(SegmentSize))
				if err != nil && err != io.EOF {
					logs.LogFatal(err.Error())
				}
				err = fd.Close()
				if err != nil {
					logs.LogFatal(err.Error())
				}
			}
		}
		err := writer.Close()
		if err != nil {
			logs.LogFatal(err.Error())
		}
		if !finished {
		retry:
			req, err := http.NewRequest(method, url, payload)
			if err != nil {
				logs.LogFatal(err.Error())
			}
			req.Header.Set("Connection", "keep-alive")
			req.Header.Set("Keep-Alive", strings.Join([]string{"timeout=", strconv.Itoa(120)}, ""))
			req.Header.Set("Content-Type", writer.FormDataContentType())
			logs.LogInfo("request =>> %v %v[%v] %v %v", method, url, len(Filelist), Filelist, uuid)
			/// request
			res, err := client.Do(req)
			if err != nil {
				logs.LogError(err.Error())
				// logs.LogClose()
				goto retry
			}
			for {
				/// response
				body, err := ioutil.ReadAll(res.Body)
				if err != nil {
					logs.LogError(err.Error())
					break
				}
				if len(body) == 0 {
					break
				}
				resp := Resp{}
				err = json.Unmarshal(body, &resp)
				if err != nil {
					logs.LogError(err.Error())
					logs.LogWarn("%v", string(body))
					continue
				}
				// 检查有无 resp 错误码
				switch resp.ErrCode {
				case ErrParamsUUID.ErrCode:
					fallthrough
				case ErrParsePartData.ErrCode:
					// 需要继续重试
					logs.LogError("--- %v %v", resp.Uuid, resp.ErrMsg)
					continue
				}
				// 读取每个文件上传状态数据
				for _, result := range resp.Data {
					switch result.ErrCode {
					case ErrParseFormFile.ErrCode:
						fallthrough
					case ErrParamsSegSizeLimit.ErrCode:
						fallthrough
					case ErrParamsSegSizeZero.ErrCode:
						fallthrough
					case ErrParamsTotalLimit.ErrCode:
						fallthrough
					case ErrParamsOffset.ErrCode:
						fallthrough
					case ErrParamsMD5.ErrCode:
						fallthrough
					case ErrParamsAllTotalLimit.ErrCode:
						fallthrough
					case ErrRepeat.ErrCode:
						logs.LogError("--- %v %v[%v] %v => %v", result.Uuid, result.Md5, result.File, result.ErrMsg, result.Message)
					// 上传成功(分段续传)，继续读取文件剩余字节继续上传
					case ErrSegOk.ErrCode:
						logs.LogError("--- %v %v[%v] %v", result.Uuid, result.Md5, result.File, result.ErrMsg)
						if result.Now <= 0 {
							break
						}
						// 上传进度写入临时文件
						fd, err := os.OpenFile(tmp_dir+result.Md5+".tmp", os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0777)
						if err != nil {
							logs.LogError(err.Error())
							break
						}
						b, err := json.Marshal(&result)
						if err != nil {
							logs.LogFatal(err.Error())
							break
						}
						_, err = fd.Write(b)
						if err != nil {
							logs.LogFatal(err.Error())
							break
						}
						err = fd.Close()
						if err != nil {
							logs.LogFatal(err.Error())
						}
						// 更新文件读取偏移
						offset[result.Md5] = result.Now
					case ErrCheckReUpload.ErrCode:
						// 校正需要重传
						offset[result.Md5] = result.Now
						logs.LogError("--- %v %v[%v] %v => %v", result.Uuid, result.Md5, result.File, result.ErrMsg, result.Message)
					case ErrFileMd5.ErrCode:
						// 上传失败了
						fallthrough
					case ErrOk.ErrCode:
						offset[result.Md5] = total[result.Md5]
						removeMd5File(&MD5, result.Md5)
						// 上传完成，删除临时文件
						os.Remove(tmp_dir + result.Md5 + ".tmp")
						logs.LogTrace("--- %v %v[%v] %v => %v", result.Uuid, result.Md5, result.File, result.ErrMsg, result.Url)
					}
				}
			}
			res.Body.Close()
		} else {
			break
		}
	}
	logs.LogClose()
}
