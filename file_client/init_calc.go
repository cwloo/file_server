package main

import (
	"encoding/json"
	"io/ioutil"
	"os"

	"github.com/cwloo/gonet/logs"
	"github.com/cwloo/gonet/utils"
)

func calcFileSize(MD5 map[string]string) (total map[string]int64, offset map[string]int64) {
	total = map[string]int64{}
	offset = map[string]int64{}
	for f, md5 := range MD5 {
		offset[md5] = int64(0)
		sta, err := os.Stat(f)
		if err != nil && os.IsNotExist(err) {
			logs.LogFatal("%v", err.Error())
		}
		total[md5] = sta.Size()
	}
	return
}

func calcFileMd5(filelist []string) (md5 map[string]string) {
	md5 = map[string]string{}
	for _, filename := range filelist {
		_, err := os.Stat(filename)
		if err != nil && os.IsNotExist(err) {
			continue
		}
		fd, err := os.OpenFile(filename, os.O_RDONLY, 0)
		if err != nil {
			logs.LogError("%v", err.Error())
			return nil
		}
		b, err := ioutil.ReadAll(fd)
		if err != nil {
			logs.LogFatal("%v", err.Error())
			return nil
		}
		md5[filename] = utils.MD5Byte(b, false)
		err = fd.Close()
		if err != nil {
			logs.LogFatal("%v", err.Error())
		}
	}
	return
}

func loadTmpFile(dir string, MD5 map[string]string) (results map[string]Result) {
	results = map[string]Result{}
	for _, md5 := range MD5 {
		f := dir + "/" + md5 + ".tmp"
		_, err := os.Stat(f)
		if err != nil && os.IsNotExist(err) {
			continue
		}
		fd, err := os.OpenFile(f, os.O_RDONLY, 0)
		if err != nil {
			logs.LogFatal("%v", err.Error())
			return
		}
		data, err := ioutil.ReadAll(fd)
		if err != nil {
			logs.LogFatal("%v", err.Error())
			return
		}
		var result Result
		err = json.Unmarshal(data, &result)
		if err != nil {
			logs.LogFatal("%v", err.Error())
			return
		}
		results[md5] = result
		err = fd.Close()
		if err != nil {
			logs.LogFatal("%v", err.Error())
			return
		}
	}
	return
}
