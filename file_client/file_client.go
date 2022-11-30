package main

import (
	"bytes"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"net/http/cookiejar"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/cwloo/gonet/logs"
	"github.com/cwloo/gonet/utils"
)

const (
	BUFSIZ = 1024 * 1024 * 10
)

func main() {
	path, _ := os.Executable()
	dir, exe := filepath.Split(path)

	logs.LogTimezone(logs.MY_CST)
	logs.LogInit(dir+"/logs", logs.LVL_DEBUG, exe, 100000000)
	logs.LogMode(logs.M_STDOUT_FILE)

	transport := &http.Transport{
		DisableKeepAlives:     false,
		TLSHandshakeTimeout:   time.Duration(3600) * time.Second,
		IdleConnTimeout:       time.Duration(3600) * time.Second,
		ResponseHeaderTimeout: time.Duration(3600) * time.Second,
		ExpectContinueTimeout: time.Duration(3600) * time.Second,
	}
	jar, _ := cookiejar.New(nil)
	client := &http.Client{
		Jar:       jar,
		Timeout:   time.Duration(3600) * time.Second,
		Transport: transport,
	}

	method := "POST"
	url := "http://192.168.1.113:8088/upload"

	userId := "1001"           //上传操作用户
	uuid := utils.CreateGUID() //本次上传标识
	filelist := []string{
		"/home/go1.19.3.linux-amd64.tar.gz",
		// "/home/OpenIMSetup1.1.2.exe",
	}
	offset := make([]int64, len(filelist))  //分段读取文件偏移
	finished := make([]bool, len(filelist)) //标识文件读取完毕
	total := make([]int64, len(filelist))   //文件总大小
	md5 := make([]string, len(filelist))    //文件md5值
	for i, filename := range filelist {
		offset[i] = int64(0)
		finished[i] = false
		sta, err := os.Stat(filename) //读取当前上传文件大小
		if err != nil {
			logs.LogFatal("%v", err.Error())
		}
		total[i] = sta.Size() //单个文件总大小
	}
	//计算文件md5值
	for i, filename := range filelist {
		fd, err := os.OpenFile(filename, os.O_RDONLY, 0)
		if err != nil {
			logs.LogError("%v", err.Error())
			return
		}
		b, err := ioutil.ReadAll(fd)
		if err != nil {
			logs.LogError("%v", err.Error())
			return
		}
		md5[i] = utils.MD5Byte(b, false)
		err = fd.Close()
		if err != nil {
			logs.LogFatal("%v", err.Error())
		}
	}
	finished_c := 0
	finished_all := false
	for {
		// 每次断点续传的payload数据
		payload := &bytes.Buffer{}
		writer := multipart.NewWriter(payload)
		_ = writer.WriteField("userId", userId) //上传操作用户
		_ = writer.WriteField("uuid", uuid)     //本次上传标识
		// 要上传的文件列表，各个文件都上传一点
		for i, filename := range filelist {
			_ = writer.WriteField("file"+strconv.Itoa(i)+".total", strconv.FormatInt(total[i], 10)) //文件总大小
			_ = writer.WriteField("file"+strconv.Itoa(i)+".md5", md5[i])                            //文件md5值
			// 每次断点续传上传 BUFSIZ 字节大小
			part, err := writer.CreateFormFile("file"+strconv.Itoa(i), filepath.Base(filename))
			if err != nil {
				logs.LogFatal("%v", err.Error())
			}
			if !finished[i] {
				//打开当前上传文件
				file, err := os.OpenFile(filename, os.O_RDONLY, 0)
				if err != nil {
					logs.LogFatal("%v", err.Error())
				}
				//读取分段数据
				file.Seek(offset[i], io.SeekStart)
				n, err := io.CopyN(part, file, int64(BUFSIZ))
				if err != nil && err != io.EOF {
					logs.LogFatal("%v", err.Error())
				}
				//关闭当前文件
				err = file.Close()
				if err != nil {
					logs.LogFatal("%v", err.Error())
				}
				if n == 0 {
					finished[i] = true
					finished_c++
					// logs.LogInfo("%v Content-Range: %v-%v/%v finished", "file"+strconv.Itoa(i), offset[i], offset[i]+n, total[i])
					continue
				} else {
					// logs.LogInfo("%v Content-Range: %v-%v/%v", "file"+strconv.Itoa(i), offset[i], offset[i]+n, total[i])
					offset[i] += n
				}
			} else if finished_c == len(filelist) {
				finished_all = true
				break
			}
		}
		//必须发起HTTP请求之前关闭 writer
		err := writer.Close()
		if err != nil {
			logs.LogFatal("%v", err.Error())
		}
		if !finished_all {
			req, err := http.NewRequest(method, url, payload)
			if err != nil {
				logs.LogFatal("%v", err.Error())
			}

			req.Header.Set("Connection", "keep-alive")
			req.Header.Set("Keep-Alive", strings.Join([]string{"timeout=", strconv.Itoa(120)}, ""))
			req.Header.Set("Content-Type", writer.FormDataContentType())
			// logs.LogInfo("user:%v:%v %v %v %v", userId, uuid, method, url, filelist)

			//request
			res, err := client.Do(req)
			if err != nil {
				logs.LogError("%v", err.Error())
				logs.LogClose()
				return
			}
			defer res.Body.Close()
			for {
				// response
				body, err := ioutil.ReadAll(res.Body)
				if err != nil {
					logs.LogError("%v", err.Error())
					break
				}
				if len(body) == 0 {
					break
				}
				logs.LogInfo(string(body))
			}
		} else {
			break
		}
	}
	logs.LogClose()
}
