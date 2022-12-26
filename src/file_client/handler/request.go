package handler

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"math/rand"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/cwloo/gonet/core/base/sys"
	"github.com/cwloo/gonet/logs"
	"github.com/cwloo/uploader/src/config"
	"github.com/cwloo/uploader/src/global"
	"github.com/cwloo/uploader/src/global/httpcli"
)

func GetRouter(client *http.Client, md5 string) string {
	method := "GET"
	r := rand.Int() % len(config.Config.Client.Addr)
	url := strings.Join([]string{
		config.Config.Client.Addr[r].Proto, "://",
		config.Config.Client.Addr[r].Ip, ":",
		strconv.Itoa(config.Config.Client.Addr[r].Port),
		config.Config.Client.Path.Router}, "")
	logs.Tracef("request =>> %v %v", method, url+"?md5="+md5)
	res, err := client.Get(url + "?md5=" + md5)
	if err != nil {
		logs.Errorf(err.Error())
		return ""
	}
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		logs.Errorf(err.Error())
		return ""
	}
	if len(body) == 0 {
		return ""
	}
	resp := global.RouterResp{}
	err = json.Unmarshal(body, &resp)
	if err != nil {
		logs.Errorf(err.Error())
		logs.Warnf("%v", string(body))
		return ""
	}
	switch resp.ErrCode {
	case 0:
		if resp.Node.Domain == "" {
			logs.Fatalf("error")
		}
		logs.Tracef("%v", string(body))
		return resp.Node.Domain
	}
	return ""
}

func Upload() {
	_, filelist := parseargs()
	if len(filelist) == 0 {
		return
	}
	client := httpcli.HttpClient()
	method := "POST"
	MD5 := calcFileMd5(filelist)              //文件md5值
	total, offset, uuids := calcFileSize(MD5) //文件大小/偏移
	router := map[string]string{}
	for _, md5 := range MD5 {
		router[md5] = GetRouter(client, md5)
	}
	for {
		finished := true
		for f, md5 := range MD5 {
			if offset[md5] < total[md5] {
				finished = false
				payload := &bytes.Buffer{}
				writer := multipart.NewWriter(payload)
				_ = writer.WriteField("uuid", uuids[md5])
				_ = writer.WriteField("md5", md5)
				_ = writer.WriteField("offset", strconv.FormatInt(offset[md5], 10)) //文件偏移量
				_ = writer.WriteField("total", strconv.FormatInt(total[md5], 10))   //文件总大小
				part, err := writer.CreateFormFile("file", filepath.Base(f))
				if err != nil {
					logs.Fatalf(err.Error())
				}
				fd, err := os.OpenFile(f, os.O_RDONLY, 0)
				if err != nil {
					logs.Fatalf(err.Error())
				}
				// 单个文件分片上传大小
				fd.Seek(offset[md5], io.SeekStart)
				_, err = io.CopyN(part, fd, int64(config.Config.Client.Upload.SegmentSize))
				if err != nil && err != io.EOF {
					logs.Fatalf(err.Error())
				}
				err = fd.Close()
				if err != nil {
					logs.Fatalf(err.Error())
				}
				err = writer.Close()
				if err != nil {
					logs.Fatalf(err.Error())
				}
				if router[md5] == "" {
					continue
				}
				url := sys.CorrectPath(strings.Join([]string{router[md5], config.Config.Client.Path.Upload}, ""))
				req, err := http.NewRequest(method, url, payload)
				if err != nil {
					logs.Fatalf(err.Error())
				}
				req.Header.Set("Connection", "keep-alive")
				req.Header.Set("Keep-Alive", strings.Join([]string{"timeout=", strconv.Itoa(120)}, ""))
				req.Header.Set("Content-Type", writer.FormDataContentType())
				logs.Infof("request =>> %v %v %v", method, url, uuids[md5])
				/// request
				res, err := client.Do(req)
				if err != nil {
					logs.Errorf(err.Error())
					continue
				}
				/// response
				body, err := ioutil.ReadAll(res.Body)
				if err != nil {
					logs.Errorf(err.Error())
					break
				}
				if len(body) == 0 {
					break
				}
				resp := global.Resp{}
				err = json.Unmarshal(body, &resp)
				if err != nil {
					logs.Errorf(err.Error())
					logs.Warnf("%v", string(body))
					continue
				}
				// 检查有无 resp 错误码
				switch resp.ErrCode {
				case global.ErrMultiFileNotSupport.ErrCode:
					fallthrough
				case global.ErrParamsUUID.ErrCode:
					fallthrough
				case global.ErrParsePartData.ErrCode:
					logs.Errorf("--- %v %v", resp.Uuid, resp.ErrMsg)
					continue
				}
				// 读取每个文件上传状态数据
				for _, result := range resp.Data {
					switch result.ErrCode {
					case global.ErrParseFormFile.ErrCode:
						fallthrough
					case global.ErrParamsSegSizeLimit.ErrCode:
						fallthrough
					case global.ErrParamsSegSizeZero.ErrCode:
						fallthrough
					case global.ErrParamsTotalLimit.ErrCode:
						fallthrough
					case global.ErrParamsOffset.ErrCode:
						fallthrough
					case global.ErrParamsMD5.ErrCode:
						fallthrough
					case global.ErrParamsAllTotalLimit.ErrCode:
						logs.Errorf("--- %v %v[%v] %v => %v", result.Uuid, result.Md5, result.File, result.ErrMsg, result.Message)
						continue

						// 别人正在上传该文件的话，你要拿到上传文件的uuid和now值并继续重试，因为别人有可能暂停上传，这样你就会接着上传该文件
					case global.ErrRepeat.ErrCode:
						url := sys.CorrectPath(strings.Join([]string{router[md5], config.Config.Client.Path.Fileinfo}, ""))
						logs.Warnf("--- %v %v[%v] %v => %v", result.Uuid, result.Md5, result.File, result.ErrMsg, result.Message)
						logs.Warnf("request =>> %v %v", method, url+"?md5="+result.Md5)
						res, err := client.Get(url + "?md5=" + result.Md5)
						if err != nil {
							logs.Errorf(err.Error())
							continue
						}
						body, err := ioutil.ReadAll(res.Body)
						if err != nil {
							logs.Errorf(err.Error())
							continue
						}
						if len(body) == 0 {
							break
						}
						resp := global.FileInfoResp{}
						err = json.Unmarshal(body, &resp)
						if err != nil {
							logs.Errorf(err.Error())
							logs.Warnf("%v", string(body))
							continue
						}
						if resp.Uuid == "" {
							logs.Fatalf("error")
						}
						uuids[md5] = resp.Uuid
						offset[resp.Md5] = resp.Now
						continue

						// 上传成功(分段续传)，继续读取文件剩余字节继续上传
					case global.ErrSegOk.ErrCode:
						if result.Now <= 0 {
							break
						}
						progress, _ := strconv.ParseFloat(fmt.Sprintf("%f", float64(result.Now)/float64(result.Total)), 64)
						logs.Debugf("--- %v %v[%v] %v %.2f%%", result.Uuid, result.Md5, result.File, result.ErrMsg, progress*100)
						// 上传进度写入临时文件
						// fd, err := os.OpenFile(tmp_dir+result.Md5+".tmp", os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0777)
						// if err != nil {
						// 	logs.Errorf(err.Error())
						// 	break
						// }
						// b, err := json.Marshal(&result)
						// if err != nil {
						// 	logs.Fatalf(err.Error())
						// 	break
						// }
						// _, err = fd.Write(b)
						// if err != nil {
						// 	logs.Fatalf(err.Error())
						// 	break
						// }
						// err = fd.Close()
						// if err != nil {
						// 	logs.Fatalf(err.Error())
						// }
						offset[result.Md5] = result.Now

						// 校正需要重传，有可能别人正在上传该文件，你会一直收到校正重传，所以只需显示进度即可并继续重试，如果上传用户暂停的话，你会接着上传该文件
					case global.ErrCheckReUpload.ErrCode:
						offset[result.Md5] = result.Now
						progress, _ := strconv.ParseFloat(fmt.Sprintf("%f", float64(result.Now)/float64(result.Total)), 64)
						logs.Errorf("--- %v %v[%v] %v %.2f%%", result.Uuid, result.Md5, result.File, result.ErrMsg, progress*100)

						// 上传完成，校验失败
					case global.ErrFileMd5.ErrCode:
						fallthrough

						// 上传完成，并且成功
					case global.ErrOk.ErrCode:
						offset[result.Md5] = total[result.Md5]
						removeMd5File(&MD5, result.Md5)
						logs.Tracef("--- %v %v[%v] %v => %v", result.Uuid, result.Md5, result.File, result.ErrMsg, result.Url)
					}
				}
				res.Body.Close()
			}
		}
		if finished {
			break
		}
	}
	logs.Close()
}
