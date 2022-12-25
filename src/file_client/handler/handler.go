package handler

import (
	"io/ioutil"
	"os"
	"strconv"
	"strings"

	"github.com/cwloo/gonet/core/base/sys/cmd"
	"github.com/cwloo/gonet/logs"
	"github.com/cwloo/gonet/utils"
)

func parseargs() (id int, filelist []string) {
	num := utils.Atoi(cmd.Arg("n"))
	for i := 0; i < num; i++ {
		filelist = append(filelist, cmd.Arg(strings.Join([]string{"file", strconv.Itoa(i)}, "")))
	}
	logs.Warnf("%v", os.Args)
	return
}

func calcFileSize(MD5 map[string]string) (total map[string]int64, offset map[string]int64, uuids map[string]string) {
	total = map[string]int64{}
	offset = map[string]int64{}
	uuids = map[string]string{}
	for f, md5 := range MD5 {
		sta, err := os.Stat(f)
		if err != nil && os.IsNotExist(err) {
			logs.Fatalf(err.Error())
		}
		if sta.Size() > 0 {
			offset[md5] = int64(0)
			total[md5] = sta.Size()
			uuids[md5] = utils.CreateGUID()
		}
	}
	return
}

func calcFileMd5(filelist []string) (md5 map[string]string) {
	md5 = map[string]string{}
	for _, f := range filelist {
		_, err := os.Stat(f)
		if err != nil && os.IsNotExist(err) {
			continue
		}
		fd, err := os.OpenFile(f, os.O_RDONLY, 0)
		if err != nil {
			logs.Errorf(err.Error())
			return nil
		}
		b, err := ioutil.ReadAll(fd)
		if err != nil {
			logs.Fatalf(err.Error())
			return nil
		}
		md5[f] = utils.MD5Byte(b, false)
		err = fd.Close()
		if err != nil {
			logs.Fatalf(err.Error())
		}
	}
	return
}

// func filePathBy(MD5 *map[string]string, md5 string) string {
// 	for f, v := range *MD5 {
// 		if v == md5 {
// 			return f
// 		}
// 	}
// 	return ""
// }

func removeMd5File(MD5 *map[string]string, md5 string) {
	for f, v := range *MD5 {
		if v == md5 {
			delete(*MD5, f)
			break
		}
	}
}

// func loadTmpFile(dir string, MD5 map[string]string) (results map[string]Result) {
// 	results = map[string]Result{}
// 	for _, md5 := range MD5 {
// 		f := dir + md5 + ".tmp"
// 		_, err := os.Stat(f)
// 		if err != nil && os.IsNotExist(err) {
// 			continue
// 		}
// 		fd, err := os.OpenFile(f, os.O_RDONLY, 0)
// 		if err != nil {
// 			logs.Fatalf(err.Error())
// 			return
// 		}
// 		data, err := ioutil.ReadAll(fd)
// 		if err != nil {
// 			logs.Fatalf(err.Error())
// 			return
// 		}
// 		var result Result
// 		err = json.Unmarshal(data, &result)
// 		if err != nil {
// 			logs.Fatalf(err.Error())
// 			return
// 		}
// 		results[md5] = result
// 		err = fd.Close()
// 		if err != nil {
// 			logs.Fatalf(err.Error())
// 			return
// 		}
// 	}
// 	return
// }
