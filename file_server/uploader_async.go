package main

import (
	"encoding/json"
	"io"
	"io/ioutil"
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
// AsyncUploader 异步方式上传
// <summary>
type AsyncUploader struct {
	uuid     string
	pipe     pipe.Pipe
	file     map[string]bool
	l        *sync.RWMutex
	tm       time.Time
	l_tm     *sync.RWMutex
	signaled bool
	l_signal *sync.Mutex
	cond     *sync.Cond
}

func NewAsyncUploader(uuid string) Uploader {
	s := &AsyncUploader{
		signaled: false,
		uuid:     uuid,
		tm:       time.Now(),
		file:     map[string]bool{},
		l:        &sync.RWMutex{},
		l_tm:     &sync.RWMutex{},
		l_signal: &sync.Mutex{},
	}
	s.cond = sync.NewCond(s.l_signal)
	mq := lq.NewQueue(1000)
	runner := NewProcessor(s.handler)
	s.pipe = pipe.NewPipeWithQuit(i32.New(), "uploader.pipe", mq, runner, s.onQuit)
	return s
}

func (s *AsyncUploader) notify() {
	s.l_signal.Lock()
	s.signaled = true
	s.cond.Signal()
	s.l_signal.Unlock()
}

func (s *AsyncUploader) wait() {
	s.l_signal.Lock()
	for !s.signaled {
		s.cond.Wait()
	}
	s.signaled = false
	s.l_signal.Unlock()
}

func (s *AsyncUploader) update() {
	s.l.Lock()
	s.tm = time.Now()
	s.l.Unlock()
}

func (s *AsyncUploader) Get() time.Time {
	s.l.RLock()
	tm := s.tm
	s.l.RUnlock()
	return tm
}

func (s *AsyncUploader) Close() {
	s.pipe.Close()
}

func (s *AsyncUploader) NotifyClose() {
	s.pipe.NotifyClose()
}

func (s *AsyncUploader) Clear() {
	s.l.RLock()
	for md5 := range s.file {
		fileInfos.Remove(md5)
	}
	s.l.RUnlock()
}

func (s *AsyncUploader) onQuit(slot run.Slot) {
	// logs.LogError("uuid:%v", s.uuid)
	s.Clear()
	uploaders.Remove(s.uuid)
}

func (s *AsyncUploader) tryAdd(md5 string) {
	s.l.Lock()
	if _, ok := s.file[md5]; !ok {
		s.file[md5] = false
	}
	s.l.Unlock()
}

func (s *AsyncUploader) setFinished(md5 string) {
	s.l.Lock()
	if _, ok := s.file[md5]; ok {
		s.file[md5] = true
	}
	s.l.Unlock()
}

func (s *AsyncUploader) hasFinishedAll() bool {
	s.l.RLock()
	for _, v := range s.file {
		if !v {
			s.l.RUnlock()
			return false
		}
	}
	s.l.RUnlock()
	return true
}

func (s *AsyncUploader) handler(msg any, args ...any) (exit bool) {
	req := msg.(*Req)
	s.uploading(req)
	exit = s.hasFinishedAll()
	if exit {
		logs.LogTrace("--------------------- ****** 无待上传文件，结束任务 uuid:%v ...", s.uuid)
	}
	return
}

func (s *AsyncUploader) Upload(req *Req) {
	s.pipe.Do(req)
	/// http.ResponseWriter 生命周期原因，不支持异步，所以加了 wait
	s.wait()
}

func (s *AsyncUploader) uploading(req *Req) {
	s.update()
	resp := req.resp
	result := req.result
	for _, k := range req.keys {
		offset := req.r.FormValue(k + ".offset")
		total := req.r.FormValue(k + ".total")
		md5 := strings.ToLower(k)
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
		////// 还未接收完
		if info.Finished() {
			logs.LogFatal("uuid:%v:%v(%v) finished", info.Uuid, info.SrcName, info.Md5)
		}
		////// 校验uuid
		if req.uuid != info.Uuid {
			logs.LogFatal("uuid:%v:%v(%v) uuid:%v", info.Uuid, info.SrcName, info.Md5, req.uuid)
		}
		////// 校验MD5
		if md5 != info.Md5 {
			logs.LogFatal("uuid:%v:%v(%v) md5:%v", info.Uuid, info.SrcName, info.Md5, md5)
		}
		////// 校验数据大小
		if total != strconv.FormatInt(info.Total, 10) {
			logs.LogFatal("uuid:%v:%v(%v) info.total:%v total:%v", info.Uuid, info.SrcName, info.Md5, info.Total, total)
		}
		////// 校验文件offset
		if offset != strconv.FormatInt(info.Now, 10) {
			result = append(result,
				Result{
					Uuid:    info.Uuid,
					File:    info.SrcName,
					Md5:     info.Md5,
					Now:     info.Now,
					Total:   info.Total,
					Expired: s.Get().Add(time.Duration(PendingTimeout) * time.Second).Unix(),
					ErrCode: ErrCheckReUpload.ErrCode,
					ErrMsg:  ErrCheckReUpload.ErrMsg,
					Result:  strings.Join([]string{"uuid:", info.Uuid, " check reuploading ", info.DstName, " progress:", strconv.FormatInt(info.Now, 10), "/", total}, ""),
				})
			// logs.LogError("uuid:%v:%v(%v) %v/%v offset:%v", info.Uuid, info.SrcName, info.Md5, info.Now, info.Total, offset)
			offset_n, _ := strconv.ParseInt(offset, 10, 0)
			logs.LogDebug("--------------------- ****** checking re-upload uuid:%v %v=%v[%v] %v/%v offset:%v seg_size[%d]", info.Uuid, k, header.Filename, md5, info.Now, total, offset_n, header.Size)
			continue
		}
		////// 检查上传目录
		_, err = os.Stat(dir + "upload/")
		if err != nil && os.IsNotExist(err) {
			os.MkdirAll(dir+"upload/", 0777)
		}
		////// 检查上传文件
		f := dir + "upload/" + info.DstName
		_, err = os.Stat(f)
		if err != nil && os.IsNotExist(err) {
		} else {
			/// 第一次写如果文件存在则删除
			if info.Now == int64(0) {
				os.Remove(f)
				logs.LogFatal("error")
			}
		}
		fd, err := os.OpenFile(f, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0777)
		if err != nil {
			result = append(result,
				Result{
					Uuid:    info.Uuid,
					File:    info.SrcName,
					Md5:     info.Md5,
					Now:     info.Now,
					Total:   info.Total,
					Expired: s.Get().Add(time.Duration(PendingTimeout) * time.Second).Unix(),
					ErrCode: ErrCheckReUpload.ErrCode,
					ErrMsg:  ErrCheckReUpload.ErrMsg,
					Result:  strings.Join([]string{"uuid:", info.Uuid, " check reuploading ", info.DstName, " progress:", strconv.FormatInt(info.Now, 10), "/", total}, ""),
				})
			logs.LogError("%v", err.Error())
			offset_n, _ := strconv.ParseInt(offset, 10, 0)
			logs.LogDebug("--------------------- ****** checking re-upload uuid:%v %v=%v[%v] %v/%v offset:%v seg_size[%d]", info.Uuid, k, header.Filename, md5, info.Now, total, offset_n, header.Size)
			continue
		}
		fd.Seek(0, io.SeekEnd)
		_, err = io.Copy(fd, file)
		if err != nil {
			result = append(result,
				Result{
					Uuid:    info.Uuid,
					File:    info.SrcName,
					Md5:     info.Md5,
					Now:     info.Now,
					Total:   info.Total,
					Expired: s.Get().Add(time.Duration(PendingTimeout) * time.Second).Unix(),
					ErrCode: ErrCheckReUpload.ErrCode,
					ErrMsg:  ErrCheckReUpload.ErrMsg,
					Result:  strings.Join([]string{"uuid:", info.Uuid, " check reuploading ", info.DstName, " progress:", strconv.FormatInt(info.Now, 10), "/", total}, ""),
				})
			logs.LogError("%v", err.Error())
			err = fd.Close()
			if err != nil {
				logs.LogError("%v", err.Error())
			}
			offset_n, _ := strconv.ParseInt(offset, 10, 0)
			logs.LogDebug("--------------------- ****** checking re-upload uuid:%v %v=%v[%v] %v/%v offset:%v seg_size[%d]", info.Uuid, k, header.Filename, md5, info.Now, total, offset_n, header.Size)
			continue
		} else {
			info.Update(header.Size)
		}
		err = fd.Close()
		if err != nil {
			logs.LogError("%v", err.Error())
		}
		err = file.Close()
		if err != nil {
			logs.LogError("%v", err.Error())
		}
		if info.Finished() {
			s.setFinished(info.Md5)
			// logs.LogDebug("uuid:%v %v=%v[%v] %v ==>>> %v/%v +%v last_segment[finished] checking md5 ...", s.uuid, k, header.Filename, md5, info.DstName, info.Now, total, header.Size)
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
			info.Md5Ok = (md5_calc == info.Md5)
			if info.Md5Ok {
				now := time.Now()
				info.DoneTime = now.Unix()
				info.UpdateHitTime(now)
				// fileInfos.Remove(info.Md5)
				result = append(result,
					Result{
						Uuid:    req.uuid,
						File:    header.Filename,
						Md5:     info.Md5,
						Now:     info.Now,
						Total:   info.Total,
						ErrCode: ErrOk.ErrCode,
						ErrMsg:  ErrOk.ErrMsg,
						Result:  strings.Join([]string{"uuid:", info.Uuid, " uploading ", info.DstName, " progress:", strconv.FormatInt(info.Now, 10) + "/" + total + " 上传成功!"}, "")})
				logs.LogDebug("uuid:%v %v=%v[%v] %v chkmd5 [ok] elapsed:%vms", req.uuid, k, header.Filename, md5, info.DstName, time.Since(start).Milliseconds())
			} else {
				fileInfos.Remove(info.Md5)
				os.Remove(f)
				result = append(result,
					Result{
						Uuid:    req.uuid,
						File:    header.Filename,
						Md5:     info.Md5,
						Now:     info.Now,
						Total:   info.Total,
						ErrCode: ErrFileMd5.ErrCode,
						ErrMsg:  ErrFileMd5.ErrMsg,
						Result:  strings.Join([]string{"uuid:", info.Uuid, " uploading ", info.DstName, " progress:", strconv.FormatInt(info.Now, 10) + "/" + total + " 上传完毕 MD5校验失败!"}, "")})
				logs.LogError("uuid:%v %v=%v[%v] %v chkmd5 [Err] elapsed:%vms", req.uuid, k, header.Filename, md5, info.DstName, time.Since(start).Milliseconds())
			}
		} else {
			result = append(result,
				Result{
					Uuid:    req.uuid,
					File:    header.Filename,
					Md5:     info.Md5,
					Now:     info.Now,
					Total:   info.Total,
					Expired: s.Get().Add(time.Duration(PendingTimeout) * time.Second).Unix(),
					ErrCode: ErrSegOk.ErrCode,
					ErrMsg:  ErrSegOk.ErrMsg,
					Result:  strings.Join([]string{"uuid:", info.Uuid, " uploading ", info.DstName, " progress:", strconv.FormatInt(info.Now, 10) + "/" + total}, "")})
			// if info.Now == header.Size {
			// 	logs.LogTrace("uuid:%v %v=%v[%v] %v ==>>> %v/%v +%v first_segment", req.uuid, k, header.Filename, md5, info.DstName, info.Now, total, header.Size)
			// } else {
			// 	logs.LogWarn("uuid:%v %v=%v[%v] %v ==>>> %v/%v +%v continue_segment", req.uuid, k, header.Filename, md5, info.DstName, info.Now, total, header.Size)
			// }
		}
	}
	if resp == nil {
		if len(result) > 0 {
			resp = &Resp{
				Data: result,
			}
		}
	} else {
		if len(result) > 0 {
			resp.Data = result
		}
	}
	if resp != nil {
		j, _ := json.Marshal(resp)
		req.w.Header().Set("Content-Length", strconv.Itoa(len(j)))
		req.w.Header().Set("Content-Type", "application/json")
		/// http.ResponseWriter 生命周期原因，不支持异步，所以加了 notify
		_, err := req.w.Write(j)
		if err != nil {
			logs.LogError(err.Error())
		}
		s.notify()
		// logs.LogError("uuid:%v %v", req.uuid, string(j))
	} else {
		resp = &Resp{}
		j, _ := json.Marshal(resp)
		req.w.Header().Set("Content-Length", strconv.Itoa(len(j)))
		/// http.ResponseWriter 生命周期原因，不支持异步，所以加了 notify
		_, err := req.w.Write(j)
		if err != nil {
			logs.LogError(err.Error())
		}
		s.notify()
		logs.LogFatal("uuid:%v", req.uuid)
	}
}
