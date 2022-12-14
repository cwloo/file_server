package main

import (
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/cwloo/gonet/core/base/mq/lq"
	"github.com/cwloo/gonet/core/base/pipe"
	"github.com/cwloo/gonet/core/base/run"
	"github.com/cwloo/gonet/logs"
	"github.com/cwloo/uploader/file_server/config"
	"github.com/cwloo/uploader/file_server/global"
	"github.com/cwloo/uploader/file_server/tg_bot"
)

var (
	asyncUploaders = sync.Pool{
		New: func() any {
			return &AsyncUploader{}
		},
	}
)

// <summary>
// AsyncUploader 异步方式上传
// <summary>
type AsyncUploader struct {
	uuid     string
	pipe     pipe.Pipe
	state    State
	tm       time.Time
	l_tm     *sync.RWMutex
	signaled bool
	l_signal *sync.Mutex
	cond     *sync.Cond
}

func NewAsyncUploader(uuid string) Uploader {
	s := asyncUploaders.Get().(*AsyncUploader)
	s.signaled = false
	s.uuid = uuid
	s.state = NewUploaderState()
	s.tm = time.Now()
	s.l_tm = &sync.RWMutex{}
	s.l_signal = &sync.Mutex{}
	s.cond = sync.NewCond(s.l_signal)
	mq := lq.NewQueue(1000)
	runner := NewProcessor(s.handler)
	s.pipe = pipe.NewPipeWithQuit(global.I32.New(), "uploader.pipe", mq, runner, s.onQuit)
	return s
}

func (s *AsyncUploader) reset() {
	s.pipe = nil
	s.state.Put()
}

func (s *AsyncUploader) Put() {
	s.reset()
	asyncUploaders.Put(s)
}

func (s *AsyncUploader) notify() {
	s.l_signal.Lock()
	s.signaled = true
	// s.cond.Signal()
	s.cond.Broadcast()
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
	s.l_tm.Lock()
	s.tm = time.Now()
	s.l_tm.Unlock()
}

func (s *AsyncUploader) Get() time.Time {
	s.l_tm.RLock()
	tm := s.tm
	s.l_tm.RUnlock()
	return tm
}

func (s *AsyncUploader) Close() {
	s.pipe.Close()
}

func (s *AsyncUploader) NotifyClose() {
	s.pipe.NotifyClose()
}

func (s *AsyncUploader) Remove(md5 string) {
	if s.state.Remove(md5) && s.state.AllDone() {
		s.pipe.NotifyClose()
	}
}

func (s *AsyncUploader) Clear() {
	msgs := []string{}
	s.state.Range(func(md5 string, ok bool) {
		if !ok {
			////// 任务退出，移除未决的文件
			if msg, ok := RemovePendingFile(s.uuid, md5); ok {
				msgs = append(msgs, msg)
			}
		} else {
			////// 任务退出，移除校验失败的文件
			if msg, ok := RemoveCheckErrFile(s.uuid, md5); ok {
				msgs = append(msgs, msg)
			}
		}
	})
	tg_bot.TgWarnMsg(msgs...)
}

func (s *AsyncUploader) onQuit(slot run.Slot) {
	s.Clear()
	uploaders.Remove(s.uuid).Put()
}

func (s *AsyncUploader) handler(msg any, args ...any) (exit bool) {
	req := msg.(*global.Req)
	switch config.Config.MultiFile > 0 {
	case true:
		s.multi_uploading(req)
	default:
		s.uploading(req)
	}
	exit = s.state.AllDone()
	if exit {
		logs.LogTrace("--------------------- ****** 无待上传文件，结束任务 %v ...", s.uuid)
	}
	return
}

func (s *AsyncUploader) Upload(req *global.Req) {
	s.pipe.Do(req)
	/// http.ResponseWriter 生命周期原因，不支持异步，所以加了 wait
	s.wait()
}

func (s *AsyncUploader) uploading(req *global.Req) {
	s.update()
	resp := req.Resp
	result := req.Result
	for _, k := range req.Keys {
		s.state.TryAdd(req.Md5)
		part, header, err := req.R.FormFile(k)
		if err != nil {
			logs.LogError(err.Error())
			return
		}
		info := fileInfos.Get(req.Md5)
		if info == nil {
			logs.LogFatal("error")
			return
		}
		////// 还未接收完
		if info.Done(false) {
			logs.LogFatal("%v %v[%v] %v %v/%v finished\nurl[%v]", info.Uuid(), info.SrcName(), info.Md5(), info.DstName(), info.Now(false), info.Total(false), info.Url(false))
		}
		////// 校验uuid
		if req.Uuid != info.Uuid() {
			logs.LogFatal("%v %v[%v] %v", info.Uuid(), info.SrcName(), info.Md5(), req.Uuid)
		}
		////// 校验MD5
		if req.Md5 != info.Md5() {
			logs.LogFatal("%v %v[%v] md5:%v", info.Uuid(), info.SrcName(), info.Md5(), req.Md5)
		}
		////// 校验数据大小
		if req.Total != strconv.FormatInt(info.Total(false), 10) {
			logs.LogFatal("%v %v[%v] info.total:%v total:%v", info.Uuid, info.SrcName, info.Md5, info.Total, req.Total)
		}
		////// 校验文件offset
		if req.Offset != strconv.FormatInt(info.Now(false), 10) {
			result = append(result,
				global.Result{
					Uuid:    info.Uuid(),
					File:    info.SrcName(),
					Md5:     info.Md5(),
					Now:     info.Now(false),
					Total:   info.Total(false),
					Expired: s.Get().Add(time.Duration(config.Config.PendingTimeout) * time.Second).Unix(),
					ErrCode: global.ErrCheckReUpload.ErrCode,
					ErrMsg:  global.ErrCheckReUpload.ErrMsg,
					Message: strings.Join([]string{info.Uuid(), " check reuploading ", info.DstName(), " progress:", strconv.FormatInt(info.Now(false), 10), "/", req.Total}, ""),
				})
			logs.LogError("%v %v[%v] %v/%v offset:%v", info.Uuid(), info.SrcName(), info.Md5(), info.Now(false), info.Total(false), req.Offset)
			offset_n, _ := strconv.ParseInt(req.Offset, 10, 0)
			logs.LogDebug("--------------------- ****** checking re-upload %v %v[%v] %v/%v offset:%v seg_size[%d]", info.Uuid(), header.Filename, req.Md5, info.Now(false), req.Total, offset_n, header.Size)
			continue
		}
		////// 检查上传目录
		_, err = os.Stat(config.Config.UploadlDir)
		if err != nil && os.IsNotExist(err) {
			os.MkdirAll(config.Config.UploadlDir, 0777)
		}
		////// 检查上传文件
		f := config.Config.UploadlDir + info.DstName()
		switch config.Config.WriteFile > 0 {
		case true:
			_, err = os.Stat(f)
			if err != nil && os.IsNotExist(err) {
			} else {
				/// 第一次写如果文件存在则删除
				if info.Now(false) == int64(0) {
					os.Remove(f)
					logs.LogFatal("error")
				}
			}
			fd, err := os.OpenFile(f, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0777)
			if err != nil {
				result = append(result,
					global.Result{
						Uuid:    info.Uuid(),
						File:    info.SrcName(),
						Md5:     info.Md5(),
						Now:     info.Now(false),
						Total:   info.Total(false),
						Expired: s.Get().Add(time.Duration(config.Config.PendingTimeout) * time.Second).Unix(),
						ErrCode: global.ErrCheckReUpload.ErrCode,
						ErrMsg:  global.ErrCheckReUpload.ErrMsg,
						Message: strings.Join([]string{info.Uuid(), " check reuploading ", info.DstName(), " progress:", strconv.FormatInt(info.Now(false), 10), "/", req.Total}, ""),
					})
				logs.LogError(err.Error())
				offset_n, _ := strconv.ParseInt(req.Offset, 10, 0)
				logs.LogDebug("--------------------- ****** checking re-upload %v %v[%v] %v/%v offset:%v seg_size[%d]", info.Uuid(), header.Filename, req.Md5, info.Now(false), req.Total, offset_n, header.Size)
				continue
			}
			fd.Seek(0, io.SeekEnd)
			_, err = io.Copy(fd, part)
			if err != nil {
				result = append(result,
					global.Result{
						Uuid:    info.Uuid(),
						File:    info.SrcName(),
						Md5:     info.Md5(),
						Now:     info.Now(false),
						Total:   info.Total(false),
						Expired: s.Get().Add(time.Duration(config.Config.PendingTimeout) * time.Second).Unix(),
						ErrCode: global.ErrCheckReUpload.ErrCode,
						ErrMsg:  global.ErrCheckReUpload.ErrMsg,
						Message: strings.Join([]string{info.Uuid(), " check reuploading ", info.DstName(), " progress:", strconv.FormatInt(info.Now(false), 10), "/", req.Total}, ""),
					})
				logs.LogError(err.Error())
				err = fd.Close()
				if err != nil {
					logs.LogError(err.Error())
				}
				offset_n, _ := strconv.ParseInt(req.Offset, 10, 0)
				logs.LogDebug("--------------------- ****** checking re-upload %v %v[%v] %v/%v offset:%v seg_size[%d]", info.Uuid(), header.Filename, req.Md5, info.Now(false), req.Total, offset_n, header.Size)
				continue
			}
			err = fd.Close()
			if err != nil {
				logs.LogError(err.Error())
			}
			err = part.Close()
			if err != nil {
				logs.LogError(err.Error())
			}
		default:
		}
		done, ok, url, start := info.Update(header.Size, func(info FileInfo, oss OSS) (url string, err error) {
			url, _, err = oss.UploadFile(info, header)
			if err != nil {
				logs.LogError(err.Error())
			}
			return
		}, func(info FileInfo) (time.Time, bool) {
			start := time.Now()
			switch config.Config.WriteFile > 0 {
			case true:
				switch config.Config.CheckMd5 > 0 {
				case true:
					md5 := CalcFileMd5(f)
					ok := md5 == info.Md5()
					return start, ok
				default:
					return start, true
				}
			default:
				return start, true
			}
		})
		if done {
			s.state.SetDone(info.Md5())
			logs.LogDebug("%v %v[%v] %v ==>>> %v/%v +%v last_segment[finished] checking md5 ...", s.uuid, header.Filename, req.Md5, info.DstName(), info.Now(false), req.Total, header.Size)
			if ok {
				// fileInfos.Remove(info.Md5()).Put()
				result = append(result,
					global.Result{
						Uuid:    req.Uuid,
						File:    header.Filename,
						Md5:     info.Md5(),
						Now:     info.Now(false),
						Total:   info.Total(false),
						ErrCode: global.ErrOk.ErrCode,
						ErrMsg:  global.ErrOk.ErrMsg,
						Url:     url,
						Message: strings.Join([]string{info.Uuid(), " uploading ", info.DstName(), " progress:", strconv.FormatInt(info.Now(false), 10) + "/" + req.Total + " 上传成功!"}, "")})
				logs.LogWarn("%v %v[%v] %v chkmd5 [ok] %v elapsed:%vms", req.Uuid, header.Filename, req.Md5, info.DstName(), url, time.Since(start).Milliseconds())
				tg_bot.TgSuccMsg(fmt.Sprintf("%v\n%v[%v]\n%v chkmd5 [ok]\n%v elapsed:%vms", req.Uuid, header.Filename, req.Md5, info.DstName(), url, time.Since(start).Milliseconds()))
			} else {
				fileInfos.Remove(info.Md5()).Put()
				os.Remove(f)
				result = append(result,
					global.Result{
						Uuid:    req.Uuid,
						File:    header.Filename,
						Md5:     info.Md5(),
						Now:     info.Now(false),
						Total:   info.Total(false),
						ErrCode: global.ErrFileMd5.ErrCode,
						ErrMsg:  global.ErrFileMd5.ErrMsg,
						Message: strings.Join([]string{info.Uuid(), " uploading ", info.DstName(), " progress:", strconv.FormatInt(info.Now(false), 10) + "/" + req.Total + " 上传完毕 MD5校验失败!"}, "")})
				logs.LogError("%v %v[%v] %v chkmd5 [Err] elapsed:%vms", req.Uuid, header.Filename, req.Md5, info.DstName(), time.Since(start).Milliseconds())
				tg_bot.TgErrMsg(fmt.Sprintf("%v\n%v[%v]\n%v chkmd5 [Err] elapsed:%vms", req.Uuid, header.Filename, req.Md5, info.DstName(), time.Since(start).Milliseconds()))
			}
		} else {
			result = append(result,
				global.Result{
					Uuid:    req.Uuid,
					File:    header.Filename,
					Md5:     info.Md5(),
					Now:     info.Now(false),
					Total:   info.Total(false),
					Expired: s.Get().Add(time.Duration(config.Config.PendingTimeout) * time.Second).Unix(),
					ErrCode: global.ErrSegOk.ErrCode,
					ErrMsg:  global.ErrSegOk.ErrMsg,
					Message: strings.Join([]string{info.Uuid(), " uploading ", info.DstName(), " progress:", strconv.FormatInt(info.Now(false), 10) + "/" + req.Total}, "")})
			if info.Now(false) == header.Size {
				logs.LogTrace("%v %v[%v] %v ==>>> %v/%v +%v first_segment", req.Uuid, header.Filename, req.Md5, info.DstName(), info.Now(false), req.Total, header.Size)
			} else {
				logs.LogWarn("%v %v[%v] %v ==>>> %v/%v +%v continue_segment", req.Uuid, header.Filename, req.Md5, info.DstName(), info.Now(false), req.Total, header.Size)
			}
		}
	}
	if resp == nil {
		if len(result) > 0 {
			resp = &global.Resp{
				Data: result,
			}
		}
	} else {
		if len(result) > 0 {
			resp.Data = result
		}
	}
	if resp != nil {
		/// http.ResponseWriter 生命周期原因，不支持异步，所以加了 notify
		writeResponse(req.W, req.R, resp)
		s.notify()
		// logs.LogError("%v %v", req.Uuid, string(j))
	} else {
		/// http.ResponseWriter 生命周期原因，不支持异步，所以加了 notify
		writeResponse(req.W, req.R, &global.Resp{})
		s.notify()
		logs.LogFatal("%v", req.Uuid)
	}
}

func (s *AsyncUploader) multi_uploading(req *global.Req) {
	s.update()
	resp := req.Resp
	result := req.Result
	for _, k := range req.Keys {
		offset := req.R.FormValue(k + ".offset")
		total := req.R.FormValue(k + ".total")
		md5 := strings.ToLower(k)
		s.state.TryAdd(md5)
		part, header, err := req.R.FormFile(k)
		if err != nil {
			logs.LogError(err.Error())
			return
		}
		info := fileInfos.Get(md5)
		if info == nil {
			logs.LogFatal("error")
			return
		}
		////// 还未接收完
		if info.Done(false) {
			logs.LogFatal("%v %v[%v] %v %v/%v finished\nurl[%v]", info.Uuid(), info.SrcName(), info.Md5(), info.DstName(), info.Now(false), info.Total(false), info.Url(false))
		}
		////// 校验uuid
		if req.Uuid != info.Uuid() {
			logs.LogFatal("%v %v[%v] %v", info.Uuid(), info.SrcName(), info.Md5(), req.Uuid)
		}
		////// 校验MD5
		if md5 != info.Md5() {
			logs.LogFatal("%v %v[%v] md5:%v", info.Uuid(), info.SrcName(), info.Md5(), md5)
		}
		////// 校验数据大小
		if total != strconv.FormatInt(info.Total(false), 10) {
			logs.LogFatal("%v %v[%v] info.total:%v total:%v", info.Uuid, info.SrcName, info.Md5, info.Total, total)
		}
		////// 校验文件offset
		if offset != strconv.FormatInt(info.Now(false), 10) {
			result = append(result,
				global.Result{
					Uuid:    info.Uuid(),
					File:    info.SrcName(),
					Md5:     info.Md5(),
					Now:     info.Now(false),
					Total:   info.Total(false),
					Expired: s.Get().Add(time.Duration(config.Config.PendingTimeout) * time.Second).Unix(),
					ErrCode: global.ErrCheckReUpload.ErrCode,
					ErrMsg:  global.ErrCheckReUpload.ErrMsg,
					Message: strings.Join([]string{info.Uuid(), " check reuploading ", info.DstName(), " progress:", strconv.FormatInt(info.Now(false), 10), "/", total}, ""),
				})
			// logs.LogError("%v %v[%v] %v/%v offset:%v", info.Uuid(), info.SrcName(), info.Md5(), info.Now(false), info.Total(false), offset)
			offset_n, _ := strconv.ParseInt(offset, 10, 0)
			logs.LogDebug("--------------------- ****** checking re-upload %v %v[%v] %v/%v offset:%v seg_size[%d]", info.Uuid(), header.Filename, md5, info.Now(false), total, offset_n, header.Size)
			continue
		}
		////// 检查上传目录
		_, err = os.Stat(config.Config.UploadlDir)
		if err != nil && os.IsNotExist(err) {
			os.MkdirAll(config.Config.UploadlDir, 0777)
		}
		////// 检查上传文件
		f := config.Config.UploadlDir + info.DstName()
		switch config.Config.WriteFile > 0 {
		case true:
			_, err = os.Stat(f)
			if err != nil && os.IsNotExist(err) {
			} else {
				/// 第一次写如果文件存在则删除
				if info.Now(false) == int64(0) {
					os.Remove(f)
					logs.LogFatal("error")
				}
			}
			fd, err := os.OpenFile(f, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0777)
			if err != nil {
				result = append(result,
					global.Result{
						Uuid:    info.Uuid(),
						File:    info.SrcName(),
						Md5:     info.Md5(),
						Now:     info.Now(false),
						Total:   info.Total(false),
						Expired: s.Get().Add(time.Duration(config.Config.PendingTimeout) * time.Second).Unix(),
						ErrCode: global.ErrCheckReUpload.ErrCode,
						ErrMsg:  global.ErrCheckReUpload.ErrMsg,
						Message: strings.Join([]string{info.Uuid(), " check reuploading ", info.DstName(), " progress:", strconv.FormatInt(info.Now(false), 10), "/", total}, ""),
					})
				logs.LogError(err.Error())
				offset_n, _ := strconv.ParseInt(offset, 10, 0)
				logs.LogDebug("--------------------- ****** checking re-upload %v %v[%v] %v/%v offset:%v seg_size[%d]", info.Uuid(), header.Filename, md5, info.Now(false), total, offset_n, header.Size)
				continue
			}
			fd.Seek(0, io.SeekEnd)
			_, err = io.Copy(fd, part)
			if err != nil {
				result = append(result,
					global.Result{
						Uuid:    info.Uuid(),
						File:    info.SrcName(),
						Md5:     info.Md5(),
						Now:     info.Now(false),
						Total:   info.Total(false),
						Expired: s.Get().Add(time.Duration(config.Config.PendingTimeout) * time.Second).Unix(),
						ErrCode: global.ErrCheckReUpload.ErrCode,
						ErrMsg:  global.ErrCheckReUpload.ErrMsg,
						Message: strings.Join([]string{info.Uuid(), " check reuploading ", info.DstName(), " progress:", strconv.FormatInt(info.Now(false), 10), "/", total}, ""),
					})
				logs.LogError(err.Error())
				err = fd.Close()
				if err != nil {
					logs.LogError(err.Error())
				}
				offset_n, _ := strconv.ParseInt(offset, 10, 0)
				logs.LogDebug("--------------------- ****** checking re-upload %v %v[%v] %v/%v offset:%v seg_size[%d]", info.Uuid(), header.Filename, md5, info.Now(false), total, offset_n, header.Size)
				continue
			}
			err = fd.Close()
			if err != nil {
				logs.LogError(err.Error())
			}
			err = part.Close()
			if err != nil {
				logs.LogError(err.Error())
			}
		default:
		}
		done, ok, url, start := info.Update(header.Size, func(info FileInfo, oss OSS) (url string, err error) {
			url, _, err = oss.UploadFile(info, header)
			if err != nil {
				logs.LogError(err.Error())
			}
			return
		}, func(info FileInfo) (time.Time, bool) {
			start := time.Now()
			switch config.Config.WriteFile > 0 {
			case true:
				switch config.Config.CheckMd5 > 0 {
				case true:
					md5 := CalcFileMd5(f)
					ok := md5 == info.Md5()
					return start, ok
				default:
					return start, true
				}
			default:
				return start, true
			}
		})
		if done {
			s.state.SetDone(info.Md5())
			logs.LogDebug("%v %v[%v] %v ==>>> %v/%v +%v last_segment[finished] checking md5 ...", s.uuid, header.Filename, md5, info.DstName(), info.Now(false), total, header.Size)
			if ok {
				// fileInfos.Remove(info.Md5()).Put()
				result = append(result,
					global.Result{
						Uuid:    req.Uuid,
						File:    header.Filename,
						Md5:     info.Md5(),
						Now:     info.Now(false),
						Total:   info.Total(false),
						ErrCode: global.ErrOk.ErrCode,
						ErrMsg:  global.ErrOk.ErrMsg,
						Url:     url,
						Message: strings.Join([]string{info.Uuid(), " uploading ", info.DstName(), " progress:", strconv.FormatInt(info.Now(false), 10) + "/" + total + " 上传成功!"}, "")})
				logs.LogWarn("%v %v[%v] %v chkmd5 [ok] %v elapsed:%vms", req.Uuid, header.Filename, req.Md5, info.DstName(), url, time.Since(start).Milliseconds())
				tg_bot.TgSuccMsg(fmt.Sprintf("%v\n%v[%v]\n%v chkmd5 [ok]\n%v elapsed:%vms", req.Uuid, header.Filename, req.Md5, info.DstName(), url, time.Since(start).Milliseconds()))
			} else {
				fileInfos.Remove(info.Md5()).Put()
				os.Remove(f)
				result = append(result,
					global.Result{
						Uuid:    req.Uuid,
						File:    header.Filename,
						Md5:     info.Md5(),
						Now:     info.Now(false),
						Total:   info.Total(false),
						ErrCode: global.ErrFileMd5.ErrCode,
						ErrMsg:  global.ErrFileMd5.ErrMsg,
						Message: strings.Join([]string{info.Uuid(), " uploading ", info.DstName(), " progress:", strconv.FormatInt(info.Now(false), 10) + "/" + total + " 上传完毕 MD5校验失败!"}, "")})
				logs.LogError("%v %v[%v] %v chkmd5 [Err] elapsed:%vms", req.Uuid, header.Filename, md5, info.DstName(), time.Since(start).Milliseconds())
				tg_bot.TgErrMsg(fmt.Sprintf("%v\n%v[%v]\n%v chkmd5 [Err] elapsed:%vms", req.Uuid, header.Filename, md5, info.DstName(), time.Since(start).Milliseconds()))
			}
		} else {
			result = append(result,
				global.Result{
					Uuid:    req.Uuid,
					File:    header.Filename,
					Md5:     info.Md5(),
					Now:     info.Now(false),
					Total:   info.Total(false),
					Expired: s.Get().Add(time.Duration(config.Config.PendingTimeout) * time.Second).Unix(),
					ErrCode: global.ErrSegOk.ErrCode,
					ErrMsg:  global.ErrSegOk.ErrMsg,
					Message: strings.Join([]string{info.Uuid(), " uploading ", info.DstName(), " progress:", strconv.FormatInt(info.Now(false), 10) + "/" + total}, "")})
			if info.Now(false) == header.Size {
				logs.LogTrace("%v %v[%v] %v ==>>> %v/%v +%v first_segment", req.Uuid, header.Filename, md5, info.DstName(), info.Now(false), total, header.Size)
			} else {
				logs.LogWarn("%v %v[%v] %v ==>>> %v/%v +%v continue_segment", req.Uuid, header.Filename, md5, info.DstName(), info.Now(false), total, header.Size)
			}
		}
	}
	if resp == nil {
		if len(result) > 0 {
			resp = &global.Resp{
				Data: result,
			}
		}
	} else {
		if len(result) > 0 {
			resp.Data = result
		}
	}
	if resp != nil {
		/// http.ResponseWriter 生命周期原因，不支持异步，所以加了 notify
		writeResponse(req.W, req.R, resp)
		s.notify()
		// logs.LogError("%v %v", req.Uuid, string(j))
	} else {
		/// http.ResponseWriter 生命周期原因，不支持异步，所以加了 notify
		writeResponse(req.W, req.R, &global.Resp{})
		s.notify()
		logs.LogFatal("%v", req.Uuid)
	}
}
