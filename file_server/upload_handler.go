package main

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/cwloo/gonet/core/base/mq/lq"
	"github.com/cwloo/gonet/core/base/pipe"
	"github.com/cwloo/gonet/core/base/run"
	"github.com/cwloo/gonet/logs"
	"github.com/cwloo/gonet/utils"
)

// <summary>
// Uploader
// <summary>
type Uploader struct {
	uuid  string
	pipe  pipe.Pipe
	l     *sync.RWMutex
	flags map[string]bool
	l_tm  *sync.RWMutex
	tm    time.Time
}

func NewUploader(uuid string) *Uploader {
	s := &Uploader{
		uuid:  uuid,
		tm:    time.Now(),
		flags: map[string]bool{},
		l:     &sync.RWMutex{},
		l_tm:  &sync.RWMutex{}}
	mq := lq.NewQueue(1000)
	runner := NewProcessor(s.handler)
	s.pipe = pipe.NewPipeWithQuit(i32.New(), "uploader.pipe", mq, runner, s.onQuit)
	return s
}

func (s *Uploader) update() {
	s.l.Lock()
	s.tm = time.Now()
	s.l.Unlock()
}

func (s *Uploader) Get() time.Time {
	s.l.RLock()
	tm := s.tm
	s.l.RUnlock()
	return tm
}

func (s *Uploader) Close() {
	s.pipe.Close()
}

func (s *Uploader) NotifyClose() {
	s.pipe.NotifyClose()
}

func (s *Uploader) clear() {
	s.l.RLock()
	for md5 := range s.flags {
		fileInfos.Remove(md5)
	}
	s.l.RUnlock()
}

func (s *Uploader) Do(data any) {
	s.pipe.Do(data)
}

func (s *Uploader) onQuit(slot run.Slot) {
	logs.LogError("uuid:%v", s.uuid)
	s.clear()
	uploaders.Remove(s.uuid)
}

func (s *Uploader) tryAdd(md5 string) {
	s.l.Lock()
	if _, ok := s.flags[md5]; !ok {
		s.flags[md5] = false
	}
	s.l.Unlock()
}

func (s *Uploader) setFinished(md5 string) {
	s.l.Lock()
	if _, ok := s.flags[md5]; ok {
		s.flags[md5] = true
	}
	s.l.Unlock()
}

func (s *Uploader) hasFinishedAll() bool {
	s.l.RLock()
	for _, v := range s.flags {
		if !v {
			s.l.RUnlock()
			return false
		}
	}
	s.l.RUnlock()
	return true
}

func (s *Uploader) handler(msg any, args ...any) (exit bool) {
	s.update()
	req := msg.(*Req)
	result := []Result{}
	for _, info := range req.ignore {
		now := strconv.FormatInt(info.Now, 10)
		total := strconv.FormatInt(info.Total, 10)
		result = append(result,
			Result{
				Uuid:    req.uuid,
				File:    info.SrcName,
				Md5:     info.Md5,
				ErrCode: ErrRepeat.ErrCode,
				ErrMsg:  ErrRepeat.ErrMsg,
				Result:  info.Uuid + " 正在上传 " + info.SrcName + " 进度 " + now + "/" + total})
	}
	for _, k := range req.keys {
		total := req.r.FormValue(k + ".total")
		md5 := strings.ToLower(req.r.FormValue(k + ".md5"))
		s.tryAdd(md5)
		file, header, err := req.r.FormFile(k)
		if err != nil {
			logs.LogError("%v", err.Error())
			return
		}
		info := fileInfos.Get(md5)
		if info == nil {
			logs.LogFatal("error")
			return
		}
		////// 校验MD5
		if md5 != info.Md5 {
			logs.LogFatal("uuid:%v:%v(%v) conflict md5:%v", info.Uuid, info.SrcName, info.Md5, md5)
		}
		////// 校验数据大小
		if total != strconv.FormatInt(info.Total, 10) {
			logs.LogFatal("uuid:%v:%v(%v) conflict %v:%v", info.Uuid, info.SrcName, info.Md5, total, info.Total)
		}
		////// 校验uuid
		if req.uuid != info.Uuid {
			logs.LogFatal("uuid:%v:%v(%v) conflict uuid:%v", info.Uuid, info.SrcName, info.Md5, req.uuid)
		}
		////// 校验filename
		if header.Filename != info.SrcName {
			logs.LogFatal("uuid:%v:%v(%v) conflict %v", info.Uuid, info.SrcName, info.Md5, header.Filename)
		}
		////// 还未接收完
		if info.Finished() {
			logs.LogFatal("uuid:%v:%v(%v) finished", info.Uuid, info.SrcName, info.Md5)
		}
		//检查上传目录
		_, err = os.Stat(dir + "upload/")
		if err != nil && os.IsNotExist(err) {
			os.MkdirAll(dir+"upload/", 0666)
		}
		//检查上传文件
		f := dir + "upload/" + info.DstName
		_, err = os.Stat(f)
		if err != nil && os.IsNotExist(err) {
		} else {
			//第一次写如果文件存在则删除
			if info.Now == int64(0) {
				os.Remove(f)
			}
		}
		fd, err := os.OpenFile(f, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0777)
		if err != nil {
			logs.LogFatal("%v", err.Error())
		}
		fd.Seek(0, io.SeekEnd)
		_, err = io.Copy(fd, file)
		if err != nil {
			logs.LogFatal("%v", err.Error())
		}
		info.Now += header.Size
		now := strconv.FormatInt(info.Now, 10)
		err = fd.Close()
		if err != nil {
			logs.LogFatal("%v", err.Error())
		}
		err = file.Close()
		if err != nil {
			logs.LogFatal("%v", err.Error())
		}
		if info.Finished() {
			s.setFinished(info.Md5)
			logs.LogError("uuid:%v %v=%v[%v] %v ==>>> %v/%v +%v [ok] checking md5 ...", s.uuid, k, header.Filename, md5, info.DstName, info.Now, total, header.Size)
			start := time.Now()
			fd, err := os.OpenFile(f, os.O_RDONLY, 0)
			if err != nil {
				logs.LogFatal("%v", err.Error())
			}
			b, err := ioutil.ReadAll(fd)
			if err != nil {
				logs.LogFatal("%v", err.Error())
			}
			md5_calc := utils.MD5Byte(b, false)
			err = fd.Close()
			if err != nil {
				logs.LogFatal("%v", err.Error())
			}
			if md5_calc == info.Md5 {
				logs.LogDebug("uuid:%v %v=%v[%v] %v chkmd5 [ok] elapsed:%vms", req.uuid, k, header.Filename, md5, info.DstName, time.Since(start).Milliseconds())
				fileInfos.Remove(info.Md5)
				result = append(result,
					Result{
						Uuid:    req.uuid,
						File:    header.Filename,
						Md5:     info.Md5,
						ErrCode: ErrOk.ErrCode,
						ErrMsg:  ErrOk.ErrMsg,
						Result:  info.Uuid + " 正在上传 " + info.SrcName + " 进度 " + now + "/" + total + " 上传完毕!"})
			} else {
				logs.LogError("uuid:%v %v=%v[%v] %v chkmd5 [Err] elapsed:%vms", req.uuid, k, header.Filename, md5, info.DstName, time.Since(start).Milliseconds())
				fileInfos.Remove(info.Md5)
				//文件不完整
				result = append(result,
					Result{
						Uuid:    req.uuid,
						File:    header.Filename,
						Md5:     info.Md5,
						ErrCode: ErrFileMd5.ErrCode,
						ErrMsg:  ErrFileMd5.ErrMsg,
						Result:  info.Uuid + " 正在上传 " + info.SrcName + " 进度 " + now + "/" + total + " 上传失败，文件不完整!"})
			}
		} else {
			if info.Now == header.Size {
				logs.LogTrace("uuid:%v %v=%v[%v] %v ==>>> %v/%v +%v [first]", req.uuid, k, header.Filename, md5, info.DstName, info.Now, total, header.Size)
			} else {
				logs.LogWarn("uuid:%v %v=%v[%v] %v ==>>> %v/%v +%v [segment]", req.uuid, k, header.Filename, md5, info.DstName, info.Now, total, header.Size)
			}
			result = append(result,
				Result{
					Uuid:    req.uuid,
					File:    header.Filename,
					Md5:     info.Md5,
					ErrCode: ErrSegOk.ErrCode,
					ErrMsg:  ErrSegOk.ErrMsg,
					Result:  info.Uuid + " 正在上传 " + info.SrcName + " 进度 " + now + "/" + total})
		}
	}
	if len(result) > 0 {
		req.w.WriteHeader(http.StatusOK)
		obj := &Resp{
			Data: result,
		}
		j, _ := json.Marshal(obj)
		req.w.Write(j)
	} else {
		logs.LogFatal("error")
	}
	exit = s.hasFinishedAll()
	if exit {
		logs.LogTrace("--------------------- ****** 无待上传文件，结束任务 uuid:%v ...", s.uuid)
	}
	return
}
