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
	data     Data
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
	s.data = NewUploaderData()
	s.tm = time.Now()
	s.l_tm = &sync.RWMutex{}
	s.l_signal = &sync.Mutex{}
	s.cond = sync.NewCond(s.l_signal)
	mq := lq.NewQueue(1000)
	runner := NewProcessor(s.handler)
	s.pipe = pipe.NewPipeWithQuit(i32.New(), "uploader.pipe", mq, runner, s.onQuit)
	return s
}

func (s *AsyncUploader) reset() {
	s.pipe = nil
	s.data.Put()
}

func (s *AsyncUploader) Put() {
	s.reset()
	asyncUploaders.Put(s)
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
	if s.data.Remove(md5) && s.data.AllDone() {
		s.pipe.NotifyClose()
	}
}

func (s *AsyncUploader) Clear() {
	msgs := []string{}
	s.data.Range(func(md5 string, ok bool) {
		if !ok {
			////// 任务退出，移除未决的文件
			fileInfos.RemoveWithCond(md5, func(info FileInfo) bool {
				if info.Uuid() != s.uuid {
					logs.LogFatal("error")
				}
				if info.Done(false) {
					logs.LogFatal("error")
				}
				return true
			}, func(info FileInfo) {
				msgs = append(msgs, fmt.Sprintf("%v\n%v[%v]\n%v [Err]", info.Uuid(), info.SrcName(), md5, info.DstName()))
				os.Remove(dir_upload + info.DstName())
				info.Put()
			})
		} else {
			////// 任务退出，移除校验失败的文件
			fileInfos.RemoveWithCond(md5, func(info FileInfo) bool {
				if info.Uuid() != s.uuid {
					logs.LogFatal("error")
				}
				if !info.Done(false) {
					logs.LogFatal("error")
				}
				ok, _ := info.Ok(false)
				return !ok
			}, func(info FileInfo) {
				msgs = append(msgs, fmt.Sprintf("%v\n%v[%v]\n%v chkmd5 [Err]", info.Uuid(), info.SrcName(), md5, info.DstName()))
				os.Remove(dir_upload + info.DstName())
				info.Put()
			})
		}
	})
	TgWarnMsg(msgs...)
}

func (s *AsyncUploader) onQuit(slot run.Slot) {
	s.Clear()
	uploaders.Remove(s.uuid).Put()
}

func (s *AsyncUploader) handler(msg any, args ...any) (exit bool) {
	req := msg.(*Req)
	switch MultiFile {
	case true:
		s.multi_uploading(req)
	default:
		s.uploading(req)
	}
	exit = s.data.AllDone()
	if exit {
		logs.LogTrace("--------------------- ****** 无待上传文件，结束任务 %v ...", s.uuid)
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
		s.data.TryAdd(req.md5)
		part, header, err := req.r.FormFile(k)
		if err != nil {
			logs.LogError(err.Error())
			return
		}
		info := fileInfos.Get(req.md5)
		if info == nil {
			logs.LogFatal("error")
			return
		}
		////// 还未接收完
		if info.Done(true) {
			logs.LogFatal("%v %v(%v) finished", info.Uuid(), info.SrcName(), info.Md5())
		}
		////// 校验uuid
		if req.uuid != info.Uuid() {
			logs.LogFatal("%v %v(%v) %v", info.Uuid(), info.SrcName(), info.Md5(), req.uuid)
		}
		////// 校验MD5
		if req.md5 != info.Md5() {
			logs.LogFatal("%v %v(%v) md5:%v", info.Uuid(), info.SrcName(), info.Md5(), req.md5)
		}
		////// 校验数据大小
		if req.total != strconv.FormatInt(info.Total(true), 10) {
			logs.LogFatal("%v %v(%v) info.total:%v total:%v", info.Uuid, info.SrcName, info.Md5, info.Total, req.total)
		}
		////// 校验文件offset
		if req.offset != strconv.FormatInt(info.Now(true), 10) {
			result = append(result,
				Result{
					Uuid:    info.Uuid(),
					File:    info.SrcName(),
					Md5:     info.Md5(),
					Now:     info.Now(true),
					Total:   info.Total(true),
					Expired: s.Get().Add(time.Duration(PendingTimeout) * time.Second).Unix(),
					ErrCode: ErrCheckReUpload.ErrCode,
					ErrMsg:  ErrCheckReUpload.ErrMsg,
					Message: strings.Join([]string{info.Uuid(), " check reuploading ", info.DstName(), " progress:", strconv.FormatInt(info.Now(true), 10), "/", req.total}, ""),
				})
			// logs.LogError("%v %v(%v) %v/%v offset:%v", info.Uuid(), info.SrcName(), info.Md5(), info.Now(true), info.Total(true), req.offset)
			offset_n, _ := strconv.ParseInt(req.offset, 10, 0)
			logs.LogDebug("--------------------- ****** checking re-upload %v %v[%v] %v/%v offset:%v seg_size[%d]", info.Uuid(), header.Filename, req.md5, info.Now(true), req.total, offset_n, header.Size)
			continue
		}
		////// 检查上传目录
		_, err = os.Stat(dir_upload)
		if err != nil && os.IsNotExist(err) {
			os.MkdirAll(dir_upload, 0777)
		}
		////// 检查上传文件
		f := dir_upload + info.DstName()
		switch WriteFile {
		case true:
			_, err = os.Stat(f)
			if err != nil && os.IsNotExist(err) {
			} else {
				/// 第一次写如果文件存在则删除
				if info.Now(true) == int64(0) {
					os.Remove(f)
					logs.LogFatal("error")
				}
			}
			fd, err := os.OpenFile(f, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0777)
			if err != nil {
				result = append(result,
					Result{
						Uuid:    info.Uuid(),
						File:    info.SrcName(),
						Md5:     info.Md5(),
						Now:     info.Now(true),
						Total:   info.Total(true),
						Expired: s.Get().Add(time.Duration(PendingTimeout) * time.Second).Unix(),
						ErrCode: ErrCheckReUpload.ErrCode,
						ErrMsg:  ErrCheckReUpload.ErrMsg,
						Message: strings.Join([]string{info.Uuid(), " check reuploading ", info.DstName(), " progress:", strconv.FormatInt(info.Now(true), 10), "/", req.total}, ""),
					})
				logs.LogError(err.Error())
				offset_n, _ := strconv.ParseInt(req.offset, 10, 0)
				logs.LogDebug("--------------------- ****** checking re-upload %v %v[%v] %v/%v offset:%v seg_size[%d]", info.Uuid(), header.Filename, req.md5, info.Now(true), req.total, offset_n, header.Size)
				continue
			}
			fd.Seek(0, io.SeekEnd)
			_, err = io.Copy(fd, part)
			if err != nil {
				result = append(result,
					Result{
						Uuid:    info.Uuid(),
						File:    info.SrcName(),
						Md5:     info.Md5(),
						Now:     info.Now(true),
						Total:   info.Total(true),
						Expired: s.Get().Add(time.Duration(PendingTimeout) * time.Second).Unix(),
						ErrCode: ErrCheckReUpload.ErrCode,
						ErrMsg:  ErrCheckReUpload.ErrMsg,
						Message: strings.Join([]string{info.Uuid(), " check reuploading ", info.DstName(), " progress:", strconv.FormatInt(info.Now(true), 10), "/", req.total}, ""),
					})
				logs.LogError(err.Error())
				err = fd.Close()
				if err != nil {
					logs.LogError(err.Error())
				}
				offset_n, _ := strconv.ParseInt(req.offset, 10, 0)
				logs.LogDebug("--------------------- ****** checking re-upload %v %v[%v] %v/%v offset:%v seg_size[%d]", info.Uuid(), header.Filename, req.md5, info.Now(true), req.total, offset_n, header.Size)
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
			switch WriteFile {
			case true:
				switch CheckMd5 {
				case true:
					md5 := calFileMd5(f)
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
			s.data.SetDone(info.Md5())
			logs.LogDebug("%v %v[%v] %v ==>>> %v/%v +%v last_segment[finished] checking md5 ...", s.uuid, header.Filename, req.md5, info.DstName(), info.Now(true), req.total, header.Size)
			if ok {
				// fileInfos.Remove(info.Md5()).Put()
				result = append(result,
					Result{
						Uuid:    req.uuid,
						File:    header.Filename,
						Md5:     info.Md5(),
						Now:     info.Now(true),
						Total:   info.Total(true),
						ErrCode: ErrOk.ErrCode,
						ErrMsg:  ErrOk.ErrMsg,
						Url:     url,
						Message: strings.Join([]string{info.Uuid(), " uploading ", info.DstName(), " progress:", strconv.FormatInt(info.Now(true), 10) + "/" + req.total + " 上传成功!"}, "")})
				logs.LogWarn("%v %v[%v] %v chkmd5 [ok] %v elapsed:%vms", req.uuid, header.Filename, req.md5, info.DstName(), url, time.Since(start).Milliseconds())
				TgSuccMsg(fmt.Sprintf("%v\n%v[%v]\n%v chkmd5 [ok]\n%v elapsed:%vms", req.uuid, header.Filename, req.md5, info.DstName(), url, time.Since(start).Milliseconds()))
			} else {
				fileInfos.Remove(info.Md5()).Put()
				os.Remove(f)
				result = append(result,
					Result{
						Uuid:    req.uuid,
						File:    header.Filename,
						Md5:     info.Md5(),
						Now:     info.Now(true),
						Total:   info.Total(true),
						ErrCode: ErrFileMd5.ErrCode,
						ErrMsg:  ErrFileMd5.ErrMsg,
						Message: strings.Join([]string{info.Uuid(), " uploading ", info.DstName(), " progress:", strconv.FormatInt(info.Now(true), 10) + "/" + req.total + " 上传完毕 MD5校验失败!"}, "")})
				logs.LogError("%v %v[%v] %v chkmd5 [Err] elapsed:%vms", req.uuid, header.Filename, req.md5, info.DstName(), time.Since(start).Milliseconds())
				TgErrMsg(fmt.Sprintf("%v\n%v[%v]\n%v chkmd5 [Err] elapsed:%vms", req.uuid, header.Filename, req.md5, info.DstName(), time.Since(start).Milliseconds()))
			}
		} else {
			result = append(result,
				Result{
					Uuid:    req.uuid,
					File:    header.Filename,
					Md5:     info.Md5(),
					Now:     info.Now(true),
					Total:   info.Total(true),
					Expired: s.Get().Add(time.Duration(PendingTimeout) * time.Second).Unix(),
					ErrCode: ErrSegOk.ErrCode,
					ErrMsg:  ErrSegOk.ErrMsg,
					Message: strings.Join([]string{info.Uuid(), " uploading ", info.DstName(), " progress:", strconv.FormatInt(info.Now(true), 10) + "/" + req.total}, "")})
			if info.Now(true) == header.Size {
				logs.LogTrace("%v %v[%v] %v ==>>> %v/%v +%v first_segment", req.uuid, header.Filename, req.md5, info.DstName(), info.Now(true), req.total, header.Size)
			} else {
				logs.LogWarn("%v %v[%v] %v ==>>> %v/%v +%v continue_segment", req.uuid, header.Filename, req.md5, info.DstName(), info.Now(true), req.total, header.Size)
			}
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
		/// http.ResponseWriter 生命周期原因，不支持异步，所以加了 notify
		writeResponse(req.w, req.r, resp)
		s.notify()
		// logs.LogError("%v %v", req.uuid, string(j))
	} else {
		/// http.ResponseWriter 生命周期原因，不支持异步，所以加了 notify
		writeResponse(req.w, req.r, &Resp{})
		s.notify()
		logs.LogFatal("%v", req.uuid)
	}
}

func (s *AsyncUploader) multi_uploading(req *Req) {
	s.update()
	resp := req.resp
	result := req.result
	for _, k := range req.keys {
		offset := req.r.FormValue(k + ".offset")
		total := req.r.FormValue(k + ".total")
		md5 := strings.ToLower(k)
		s.data.TryAdd(md5)
		part, header, err := req.r.FormFile(k)
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
		if info.Done(true) {
			logs.LogFatal("%v %v(%v) finished", info.Uuid(), info.SrcName(), info.Md5())
		}
		////// 校验uuid
		if req.uuid != info.Uuid() {
			logs.LogFatal("%v %v(%v) %v", info.Uuid(), info.SrcName(), info.Md5(), req.uuid)
		}
		////// 校验MD5
		if md5 != info.Md5() {
			logs.LogFatal("%v %v(%v) md5:%v", info.Uuid(), info.SrcName(), info.Md5(), md5)
		}
		////// 校验数据大小
		if total != strconv.FormatInt(info.Total(true), 10) {
			logs.LogFatal("%v %v(%v) info.total:%v total:%v", info.Uuid, info.SrcName, info.Md5, info.Total, total)
		}
		////// 校验文件offset
		if offset != strconv.FormatInt(info.Now(true), 10) {
			result = append(result,
				Result{
					Uuid:    info.Uuid(),
					File:    info.SrcName(),
					Md5:     info.Md5(),
					Now:     info.Now(true),
					Total:   info.Total(true),
					Expired: s.Get().Add(time.Duration(PendingTimeout) * time.Second).Unix(),
					ErrCode: ErrCheckReUpload.ErrCode,
					ErrMsg:  ErrCheckReUpload.ErrMsg,
					Message: strings.Join([]string{info.Uuid(), " check reuploading ", info.DstName(), " progress:", strconv.FormatInt(info.Now(true), 10), "/", total}, ""),
				})
			// logs.LogError("%v %v(%v) %v/%v offset:%v", info.Uuid(), info.SrcName(), info.Md5(), info.Now(true), info.Total(true), offset)
			offset_n, _ := strconv.ParseInt(offset, 10, 0)
			logs.LogDebug("--------------------- ****** checking re-upload %v %v[%v] %v/%v offset:%v seg_size[%d]", info.Uuid(), header.Filename, md5, info.Now(true), total, offset_n, header.Size)
			continue
		}
		////// 检查上传目录
		_, err = os.Stat(dir_upload)
		if err != nil && os.IsNotExist(err) {
			os.MkdirAll(dir_upload, 0777)
		}
		////// 检查上传文件
		f := dir_upload + info.DstName()
		switch WriteFile {
		case true:
			_, err = os.Stat(f)
			if err != nil && os.IsNotExist(err) {
			} else {
				/// 第一次写如果文件存在则删除
				if info.Now(true) == int64(0) {
					os.Remove(f)
					logs.LogFatal("error")
				}
			}
			fd, err := os.OpenFile(f, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0777)
			if err != nil {
				result = append(result,
					Result{
						Uuid:    info.Uuid(),
						File:    info.SrcName(),
						Md5:     info.Md5(),
						Now:     info.Now(true),
						Total:   info.Total(true),
						Expired: s.Get().Add(time.Duration(PendingTimeout) * time.Second).Unix(),
						ErrCode: ErrCheckReUpload.ErrCode,
						ErrMsg:  ErrCheckReUpload.ErrMsg,
						Message: strings.Join([]string{info.Uuid(), " check reuploading ", info.DstName(), " progress:", strconv.FormatInt(info.Now(true), 10), "/", total}, ""),
					})
				logs.LogError(err.Error())
				offset_n, _ := strconv.ParseInt(offset, 10, 0)
				logs.LogDebug("--------------------- ****** checking re-upload %v %v[%v] %v/%v offset:%v seg_size[%d]", info.Uuid(), header.Filename, md5, info.Now(true), total, offset_n, header.Size)
				continue
			}
			fd.Seek(0, io.SeekEnd)
			_, err = io.Copy(fd, part)
			if err != nil {
				result = append(result,
					Result{
						Uuid:    info.Uuid(),
						File:    info.SrcName(),
						Md5:     info.Md5(),
						Now:     info.Now(true),
						Total:   info.Total(true),
						Expired: s.Get().Add(time.Duration(PendingTimeout) * time.Second).Unix(),
						ErrCode: ErrCheckReUpload.ErrCode,
						ErrMsg:  ErrCheckReUpload.ErrMsg,
						Message: strings.Join([]string{info.Uuid(), " check reuploading ", info.DstName(), " progress:", strconv.FormatInt(info.Now(true), 10), "/", total}, ""),
					})
				logs.LogError(err.Error())
				err = fd.Close()
				if err != nil {
					logs.LogError(err.Error())
				}
				offset_n, _ := strconv.ParseInt(offset, 10, 0)
				logs.LogDebug("--------------------- ****** checking re-upload %v %v[%v] %v/%v offset:%v seg_size[%d]", info.Uuid(), header.Filename, md5, info.Now(true), total, offset_n, header.Size)
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
			switch WriteFile {
			case true:
				switch CheckMd5 {
				case true:
					md5 := calFileMd5(f)
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
			s.data.SetDone(info.Md5())
			logs.LogDebug("%v %v[%v] %v ==>>> %v/%v +%v last_segment[finished] checking md5 ...", s.uuid, header.Filename, md5, info.DstName(), info.Now(true), total, header.Size)
			if ok {
				// fileInfos.Remove(info.Md5()).Put()
				result = append(result,
					Result{
						Uuid:    req.uuid,
						File:    header.Filename,
						Md5:     info.Md5(),
						Now:     info.Now(true),
						Total:   info.Total(true),
						ErrCode: ErrOk.ErrCode,
						ErrMsg:  ErrOk.ErrMsg,
						Url:     url,
						Message: strings.Join([]string{info.Uuid(), " uploading ", info.DstName(), " progress:", strconv.FormatInt(info.Now(true), 10) + "/" + total + " 上传成功!"}, "")})
				logs.LogWarn("%v %v[%v] %v chkmd5 [ok] %v elapsed:%vms", req.uuid, header.Filename, req.md5, info.DstName(), url, time.Since(start).Milliseconds())
				TgSuccMsg(fmt.Sprintf("%v\n%v[%v]\n%v chkmd5 [ok]\n%v elapsed:%vms", req.uuid, header.Filename, req.md5, info.DstName(), url, time.Since(start).Milliseconds()))
			} else {
				fileInfos.Remove(info.Md5()).Put()
				os.Remove(f)
				result = append(result,
					Result{
						Uuid:    req.uuid,
						File:    header.Filename,
						Md5:     info.Md5(),
						Now:     info.Now(true),
						Total:   info.Total(true),
						ErrCode: ErrFileMd5.ErrCode,
						ErrMsg:  ErrFileMd5.ErrMsg,
						Message: strings.Join([]string{info.Uuid(), " uploading ", info.DstName(), " progress:", strconv.FormatInt(info.Now(true), 10) + "/" + total + " 上传完毕 MD5校验失败!"}, "")})
				logs.LogError("%v %v[%v] %v chkmd5 [Err] elapsed:%vms", req.uuid, header.Filename, md5, info.DstName(), time.Since(start).Milliseconds())
				TgErrMsg(fmt.Sprintf("%v\n%v[%v]\n%v chkmd5 [Err] elapsed:%vms", req.uuid, header.Filename, md5, info.DstName(), time.Since(start).Milliseconds()))
			}
		} else {
			result = append(result,
				Result{
					Uuid:    req.uuid,
					File:    header.Filename,
					Md5:     info.Md5(),
					Now:     info.Now(true),
					Total:   info.Total(true),
					Expired: s.Get().Add(time.Duration(PendingTimeout) * time.Second).Unix(),
					ErrCode: ErrSegOk.ErrCode,
					ErrMsg:  ErrSegOk.ErrMsg,
					Message: strings.Join([]string{info.Uuid(), " uploading ", info.DstName(), " progress:", strconv.FormatInt(info.Now(true), 10) + "/" + total}, "")})
			if info.Now(true) == header.Size {
				logs.LogTrace("%v %v[%v] %v ==>>> %v/%v +%v first_segment", req.uuid, header.Filename, md5, info.DstName(), info.Now(true), total, header.Size)
			} else {
				logs.LogWarn("%v %v[%v] %v ==>>> %v/%v +%v continue_segment", req.uuid, header.Filename, md5, info.DstName(), info.Now(true), total, header.Size)
			}
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
		/// http.ResponseWriter 生命周期原因，不支持异步，所以加了 notify
		writeResponse(req.w, req.r, resp)
		s.notify()
		// logs.LogError("%v %v", req.uuid, string(j))
	} else {
		/// http.ResponseWriter 生命周期原因，不支持异步，所以加了 notify
		writeResponse(req.w, req.r, &Resp{})
		s.notify()
		logs.LogFatal("%v", req.uuid)
	}
}
