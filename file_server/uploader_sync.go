package main

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/cwloo/gonet/logs"
)

// <summary>
// SyncUploader 同步方式上传
// <summary>
type SyncUploader struct {
	uuid string
	file map[string]bool
	l    *sync.RWMutex
	tm   time.Time
	l_tm *sync.RWMutex
}

func NewSyncUploader(uuid string) Uploader {
	s := &SyncUploader{
		uuid: uuid,
		tm:   time.Now(),
		file: map[string]bool{},
		l:    &sync.RWMutex{},
		l_tm: &sync.RWMutex{},
	}
	return s
}

func (s *SyncUploader) update() {
	s.l.Lock()
	s.tm = time.Now()
	s.l.Unlock()
}

func (s *SyncUploader) Get() time.Time {
	s.l.RLock()
	tm := s.tm
	s.l.RUnlock()
	return tm
}

func (s *SyncUploader) Close() {
	s.Clear()
	uploaders.Remove(s.uuid)
}

func (s *SyncUploader) NotifyClose() {
}

func (s *SyncUploader) Clear() {
	msgs := []string{}
	s.l.RLock()
	for md5, ok := range s.file {
		if !ok {
			////// 任务退出，移除未决的文件
			fileInfos.RemoveWithCond(md5, func(info FileInfo) bool {
				if info.Uuid() != s.uuid {
					logs.LogFatal("error")
				}
				if info.Done() {
					logs.LogFatal("error")
				}
				return true
			}, func(info FileInfo) {
				msgs = append(msgs, fmt.Sprintf("%v\n%v[%v]\n%v [Err]", info.Uuid(), info.SrcName(), md5, info.DstName()))
				os.Remove(dir_upload + info.DstName())
			})
		} else {
			////// 任务退出，移除校验失败的文件
			fileInfos.RemoveWithCond(md5, func(info FileInfo) bool {
				if info.Uuid() != s.uuid {
					logs.LogFatal("error")
				}
				if !info.Done() {
					logs.LogFatal("error")
				}
				ok, _ := info.Ok()
				return !ok
			}, func(info FileInfo) {
				msgs = append(msgs, fmt.Sprintf("%v\n%v[%v]\n%v chkmd5 [Err]", info.Uuid(), info.SrcName(), md5, info.DstName()))
				os.Remove(dir_upload + info.DstName())
			})
		}
	}
	s.l.RUnlock()
	TgWarnMsg(msgs...)
}

func (s *SyncUploader) tryAdd(md5 string) {
	s.l.Lock()
	if _, ok := s.file[md5]; !ok {
		s.file[md5] = false
	}
	s.l.Unlock()
}

func (s *SyncUploader) setDone(md5 string) {
	s.l.Lock()
	if _, ok := s.file[md5]; ok {
		s.file[md5] = true
	}
	s.l.Unlock()
}

func (s *SyncUploader) allDone() bool {
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

func (s *SyncUploader) Upload(req *Req) {
	switch MultiFile {
	case true:
		s.multi_uploading(req)
	default:
		s.uploading(req)
	}
	exit := s.allDone()
	if exit {
		logs.LogTrace("--------------------- ****** 无待上传文件，结束任务 uuid:%v ...", s.uuid)
		uploaders.Remove(s.uuid)
	}
}

func (s *SyncUploader) uploading(req *Req) {
	s.update()
	resp := req.resp
	result := req.result
	for _, k := range req.keys {
		s.tryAdd(req.md5)
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
		if info.Done() {
			logs.LogFatal("uuid:%v:%v(%v) finished", info.Uuid(), info.SrcName(), info.Md5())
		}
		////// 校验uuid
		if req.uuid != info.Uuid() {
			logs.LogFatal("uuid:%v:%v(%v) uuid:%v", info.Uuid(), info.SrcName(), info.Md5(), req.uuid)
		}
		////// 校验MD5
		if req.md5 != info.Md5() {
			logs.LogFatal("uuid:%v:%v(%v) md5:%v", info.Uuid(), info.SrcName(), info.Md5(), req.md5)
		}
		////// 校验数据大小
		if req.total != strconv.FormatInt(info.Total(), 10) {
			logs.LogFatal("uuid:%v:%v(%v) info.total:%v total:%v", info.Uuid(), info.SrcName(), info.Md5(), info.Total(), req.total)
		}
		////// 校验文件offset
		if req.offset != strconv.FormatInt(info.Now(), 10) {
			result = append(result,
				Result{
					Uuid:    info.Uuid(),
					File:    info.SrcName(),
					Md5:     info.Md5(),
					Now:     info.Now(),
					Total:   info.Total(),
					Expired: s.Get().Add(time.Duration(PendingTimeout) * time.Second).Unix(),
					ErrCode: ErrCheckReUpload.ErrCode,
					ErrMsg:  ErrCheckReUpload.ErrMsg,
					Message: strings.Join([]string{"uuid:", info.Uuid(), " check reuploading ", info.DstName(), " progress:", strconv.FormatInt(info.Now(), 10), "/", req.total}, ""),
				})
			// logs.LogError("uuid:%v:%v(%v) %v/%v offset:%v", info.Uuid(), info.SrcName(), info.Md5(), info.Now(), info.Total(), req.offset)
			offset_n, _ := strconv.ParseInt(req.offset, 10, 0)
			logs.LogInfo("--------------------- checking re-upload uuid:%v %v=%v[%v] %v/%v offset:%v seg_size[%d]", info.Uuid(), k, header.Filename, req.md5, info.Now(), req.total, offset_n, header.Size)
			continue
		}
		////// 检查上传目录
		_, err = os.Stat(dir_upload)
		if err != nil && os.IsNotExist(err) {
			os.MkdirAll(dir_upload, 0777)
		}
		////// 检查上传文件
		f := dir_upload + info.DstName()
		switch WriteDisk {
		case true:
			_, err = os.Stat(f)
			if err != nil && os.IsNotExist(err) {
			} else {
				/// 第一次写如果文件存在则删除
				if info.Now() == int64(0) {
					os.Remove(f)
				}
			}
			fd, err := os.OpenFile(f, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0777)
			if err != nil {
				result = append(result,
					Result{
						Uuid:    info.Uuid(),
						File:    info.SrcName(),
						Md5:     info.Md5(),
						Now:     info.Now(),
						Total:   info.Total(),
						Expired: s.Get().Add(time.Duration(PendingTimeout) * time.Second).Unix(),
						ErrCode: ErrCheckReUpload.ErrCode,
						ErrMsg:  ErrCheckReUpload.ErrMsg,
						Message: strings.Join([]string{"uuid:", info.Uuid(), " check reuploading ", info.DstName(), " progress:", strconv.FormatInt(info.Now(), 10), "/", req.total}, ""),
					})
				logs.LogError(err.Error())
				offset_n, _ := strconv.ParseInt(req.offset, 10, 0)
				logs.LogInfo("--------------------- checking re-upload uuid:%v %v=%v[%v] %v/%v offset:%v seg_size[%d]", info.Uuid(), k, header.Filename, req.md5, info.Now(), req.total, offset_n, header.Size)
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
						Now:     info.Now(),
						Total:   info.Total(),
						Expired: s.Get().Add(time.Duration(PendingTimeout) * time.Second).Unix(),
						ErrCode: ErrCheckReUpload.ErrCode,
						ErrMsg:  ErrCheckReUpload.ErrMsg,
						Message: strings.Join([]string{"uuid:", info.Uuid(), " check reuploading ", info.DstName(), " progress:", strconv.FormatInt(info.Now(), 10), "/", req.total}, ""),
					})
				logs.LogError(err.Error())
				err = fd.Close()
				if err != nil {
					logs.LogError(err.Error())
				}
				offset_n, _ := strconv.ParseInt(req.offset, 10, 0)
				logs.LogInfo("--------------------- checking re-upload uuid:%v %v=%v[%v] %v/%v offset:%v seg_size[%d]", info.Uuid(), k, header.Filename, req.md5, info.Now(), req.total, offset_n, header.Size)
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
		done, ok, start, url := info.Update(header.Size, func(info FileInfo, done bool) (url string, err error) {
			oss := NewOss()
			url, _, err = oss.UploadFile(info, header, done)
			if err != nil {
				logs.LogError(err.Error())
			}
			return
		}, func(info FileInfo) (time.Time, bool) {
			start := time.Now()
			switch WriteDisk {
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
			s.setDone(info.Md5())
			logs.LogDebug("uuid:%v %v=%v[%v] %v ==>>> %v/%v +%v last_segment[finished] checking md5 ...", s.uuid, k, header.Filename, req.md5, info.DstName(), info.Now(), req.total, header.Size)
			if ok {
				// fileInfos.Remove(info.Md5())
				result = append(result,
					Result{
						Uuid:    req.uuid,
						File:    header.Filename,
						Md5:     info.Md5(),
						Now:     info.Now(),
						Total:   info.Total(),
						ErrCode: ErrOk.ErrCode,
						ErrMsg:  ErrOk.ErrMsg,
						Url:     url,
						Message: strings.Join([]string{"uuid:", info.Uuid(), " uploading ", info.DstName(), " progress:", strconv.FormatInt(info.Now(), 10) + "/" + req.total + " 上传成功!"}, "")})
				logs.LogWarn("uuid:%v %v[%v] %v chkmd5 [ok] %v elapsed:%vms", req.uuid, header.Filename, req.md5, info.DstName(), url, time.Since(start).Milliseconds())
				TgSuccMsg(fmt.Sprintf("%v\n%v[%v]\n%v chkmd5 [ok]\n%v elapsed:%vms", req.uuid, header.Filename, req.md5, info.DstName(), url, time.Since(start).Milliseconds()))
			} else {
				fileInfos.Remove(info.Md5())
				os.Remove(f)
				result = append(result,
					Result{
						Uuid:    req.uuid,
						File:    header.Filename,
						Md5:     info.Md5(),
						Now:     info.Now(),
						Total:   info.Total(),
						ErrCode: ErrFileMd5.ErrCode,
						ErrMsg:  ErrFileMd5.ErrMsg,
						Message: strings.Join([]string{"uuid:", info.Uuid(), " uploading ", info.DstName(), " progress:", strconv.FormatInt(info.Now(), 10) + "/" + req.total + " 上传完毕 MD5校验失败!"}, "")})
				logs.LogError("uuid:%v %v[%v] %v chkmd5 [Err] elapsed:%vms", req.uuid, header.Filename, req.md5, info.DstName(), time.Since(start).Milliseconds())
				TgErrMsg(fmt.Sprintf("%v\n%v[%v]\n%v chkmd5 [Err] elapsed:%vms", req.uuid, header.Filename, req.md5, info.DstName(), time.Since(start).Milliseconds()))
			}
		} else {
			result = append(result,
				Result{
					Uuid:    req.uuid,
					File:    header.Filename,
					Md5:     info.Md5(),
					Now:     info.Now(),
					Total:   info.Total(),
					Expired: s.Get().Add(time.Duration(PendingTimeout) * time.Second).Unix(),
					ErrCode: ErrSegOk.ErrCode,
					ErrMsg:  ErrSegOk.ErrMsg,
					Message: strings.Join([]string{"uuid:", info.Uuid(), " uploading ", info.DstName(), " progress:", strconv.FormatInt(info.Now(), 10) + "/" + req.total}, "")})
			if info.Now() == header.Size {
				logs.LogTrace("uuid:%v %v[%v] %v ==>>> %v/%v +%v first_segment", req.uuid, header.Filename, req.md5, info.DstName(), info.Now(), req.total, header.Size)
			} else {
				logs.LogWarn("uuid:%v %v[%v] %v ==>>> %v/%v +%v continue_segment", req.uuid, header.Filename, req.md5, info.DstName(), info.Now(), req.total, header.Size)
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
		j, _ := json.Marshal(resp)
		setResponseHeader(req.w, req.r)
		req.w.Header().Set("Content-Type", "application/json")
		req.w.Header().Set("Content-Length", strconv.Itoa(len(j)))
		/// http.ResponseWriter 生命周期原因，不支持异步
		_, err := req.w.Write(j)
		if err != nil {
			logs.LogError(err.Error())
		}
		// logs.LogError("uuid:%v %v", req.uuid, string(j))
	} else {
		resp = &Resp{}
		j, _ := json.Marshal(resp)
		setResponseHeader(req.w, req.r)
		req.w.Header().Set("Content-Type", "application/json")
		req.w.Header().Set("Content-Length", strconv.Itoa(len(j)))
		/// http.ResponseWriter 生命周期原因，不支持异步
		_, err := req.w.Write(j)
		if err != nil {
			logs.LogError(err.Error())
		}
		logs.LogFatal("uuid:%v", req.uuid)
	}
}

func (s *SyncUploader) multi_uploading(req *Req) {
	s.update()
	resp := req.resp
	result := req.result
	for _, k := range req.keys {
		offset := req.r.FormValue(k + ".offset")
		total := req.r.FormValue(k + ".total")
		md5 := strings.ToLower(k)
		s.tryAdd(md5)
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
		if info.Done() {
			logs.LogFatal("uuid:%v:%v(%v) finished", info.Uuid(), info.SrcName(), info.Md5())
		}
		////// 校验uuid
		if req.uuid != info.Uuid() {
			logs.LogFatal("uuid:%v:%v(%v) uuid:%v", info.Uuid(), info.SrcName(), info.Md5(), req.uuid)
		}
		////// 校验MD5
		if md5 != info.Md5() {
			logs.LogFatal("uuid:%v:%v(%v) md5:%v", info.Uuid(), info.SrcName(), info.Md5(), md5)
		}
		////// 校验数据大小
		if total != strconv.FormatInt(info.Total(), 10) {
			logs.LogFatal("uuid:%v:%v(%v) info.total:%v total:%v", info.Uuid(), info.SrcName(), info.Md5(), info.Total(), total)
		}
		////// 校验文件offset
		if offset != strconv.FormatInt(info.Now(), 10) {
			result = append(result,
				Result{
					Uuid:    info.Uuid(),
					File:    info.SrcName(),
					Md5:     info.Md5(),
					Now:     info.Now(),
					Total:   info.Total(),
					Expired: s.Get().Add(time.Duration(PendingTimeout) * time.Second).Unix(),
					ErrCode: ErrCheckReUpload.ErrCode,
					ErrMsg:  ErrCheckReUpload.ErrMsg,
					Message: strings.Join([]string{"uuid:", info.Uuid(), " check reuploading ", info.DstName(), " progress:", strconv.FormatInt(info.Now(), 10), "/", total}, ""),
				})
			// logs.LogError("uuid:%v:%v(%v) %v/%v offset:%v", info.Uuid(), info.SrcName(), info.Md5(), info.Now(), info.Total(), offset)
			offset_n, _ := strconv.ParseInt(offset, 10, 0)
			logs.LogInfo("--------------------- checking re-upload uuid:%v %v=%v[%v] %v/%v offset:%v seg_size[%d]", info.Uuid(), k, header.Filename, md5, info.Now(), total, offset_n, header.Size)
			continue
		}
		////// 检查上传目录
		_, err = os.Stat(dir_upload)
		if err != nil && os.IsNotExist(err) {
			os.MkdirAll(dir_upload, 0777)
		}
		////// 检查上传文件
		f := dir_upload + info.DstName()
		switch WriteDisk {
		case true:
			_, err = os.Stat(f)
			if err != nil && os.IsNotExist(err) {
			} else {
				/// 第一次写如果文件存在则删除
				if info.Now() == int64(0) {
					os.Remove(f)
				}
			}
			fd, err := os.OpenFile(f, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0777)
			if err != nil {
				result = append(result,
					Result{
						Uuid:    info.Uuid(),
						File:    info.SrcName(),
						Md5:     info.Md5(),
						Now:     info.Now(),
						Total:   info.Total(),
						Expired: s.Get().Add(time.Duration(PendingTimeout) * time.Second).Unix(),
						ErrCode: ErrCheckReUpload.ErrCode,
						ErrMsg:  ErrCheckReUpload.ErrMsg,
						Message: strings.Join([]string{"uuid:", info.Uuid(), " check reuploading ", info.DstName(), " progress:", strconv.FormatInt(info.Now(), 10), "/", total}, ""),
					})
				logs.LogError(err.Error())
				offset_n, _ := strconv.ParseInt(offset, 10, 0)
				logs.LogInfo("--------------------- checking re-upload uuid:%v %v=%v[%v] %v/%v offset:%v seg_size[%d]", info.Uuid(), k, header.Filename, md5, info.Now(), total, offset_n, header.Size)
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
						Now:     info.Now(),
						Total:   info.Total(),
						Expired: s.Get().Add(time.Duration(PendingTimeout) * time.Second).Unix(),
						ErrCode: ErrCheckReUpload.ErrCode,
						ErrMsg:  ErrCheckReUpload.ErrMsg,
						Message: strings.Join([]string{"uuid:", info.Uuid(), " check reuploading ", info.DstName(), " progress:", strconv.FormatInt(info.Now(), 10), "/", total}, ""),
					})
				logs.LogError(err.Error())
				err = fd.Close()
				if err != nil {
					logs.LogError(err.Error())
				}
				offset_n, _ := strconv.ParseInt(offset, 10, 0)
				logs.LogInfo("--------------------- checking re-upload uuid:%v %v=%v[%v] %v/%v offset:%v seg_size[%d]", info.Uuid(), k, header.Filename, md5, info.Now(), total, offset_n, header.Size)
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
		done, ok, start, url := info.Update(header.Size, func(info FileInfo, done bool) (url string, err error) {
			oss := NewOss()
			url, _, err = oss.UploadFile(info, header, done)
			if err != nil {
				logs.LogError(err.Error())
			}
			return
		}, func(info FileInfo) (time.Time, bool) {
			start := time.Now()
			switch WriteDisk {
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
			s.setDone(info.Md5())
			logs.LogDebug("uuid:%v %v=%v[%v] %v ==>>> %v/%v +%v last_segment[finished] checking md5 ...", s.uuid, k, header.Filename, md5, info.DstName(), info.Now(), total, header.Size)
			if ok {
				// fileInfos.Remove(info.Md5())
				result = append(result,
					Result{
						Uuid:    req.uuid,
						File:    header.Filename,
						Md5:     info.Md5(),
						Now:     info.Now(),
						Total:   info.Total(),
						ErrCode: ErrOk.ErrCode,
						ErrMsg:  ErrOk.ErrMsg,
						Url:     url,
						Message: strings.Join([]string{"uuid:", info.Uuid(), " uploading ", info.DstName(), " progress:", strconv.FormatInt(info.Now(), 10) + "/" + total + " 上传成功!"}, "")})
				logs.LogWarn("uuid:%v %v[%v] %v chkmd5 [ok] %v elapsed:%vms", req.uuid, header.Filename, req.md5, info.DstName(), url, time.Since(start).Milliseconds())
				TgSuccMsg(fmt.Sprintf("%v\n%v[%v]\n%v chkmd5 [ok]\n%v elapsed:%vms", req.uuid, header.Filename, req.md5, info.DstName(), url, time.Since(start).Milliseconds()))
			} else {
				fileInfos.Remove(info.Md5())
				os.Remove(f)
				result = append(result,
					Result{
						Uuid:    req.uuid,
						File:    header.Filename,
						Md5:     info.Md5(),
						Now:     info.Now(),
						Total:   info.Total(),
						ErrCode: ErrFileMd5.ErrCode,
						ErrMsg:  ErrFileMd5.ErrMsg,
						Message: strings.Join([]string{"uuid:", info.Uuid(), " uploading ", info.DstName(), " progress:", strconv.FormatInt(info.Now(), 10) + "/" + total + " 上传完毕 MD5校验失败!"}, "")})
				logs.LogError("uuid:%v %v[%v] %v chkmd5 [Err] elapsed:%vms", req.uuid, header.Filename, md5, info.DstName(), time.Since(start).Milliseconds())
				TgErrMsg(fmt.Sprintf("%v\n%v[%v]\n%v chkmd5 [Err] elapsed:%vms", req.uuid, header.Filename, md5, info.DstName(), time.Since(start).Milliseconds()))
			}
		} else {
			result = append(result,
				Result{
					Uuid:    req.uuid,
					File:    header.Filename,
					Md5:     info.Md5(),
					Now:     info.Now(),
					Total:   info.Total(),
					Expired: s.Get().Add(time.Duration(PendingTimeout) * time.Second).Unix(),
					ErrCode: ErrSegOk.ErrCode,
					ErrMsg:  ErrSegOk.ErrMsg,
					Message: strings.Join([]string{"uuid:", info.Uuid(), " uploading ", info.DstName(), " progress:", strconv.FormatInt(info.Now(), 10) + "/" + total}, "")})
			if info.Now() == header.Size {
				logs.LogTrace("uuid:%v %v[%v] %v ==>>> %v/%v +%v first_segment", req.uuid, header.Filename, md5, info.DstName(), info.Now(), total, header.Size)
			} else {
				logs.LogWarn("uuid:%v %v[%v] %v ==>>> %v/%v +%v continue_segment", req.uuid, header.Filename, md5, info.DstName(), info.Now(), total, header.Size)
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
		j, _ := json.Marshal(resp)
		setResponseHeader(req.w, req.r)
		req.w.Header().Set("Content-Type", "application/json")
		req.w.Header().Set("Content-Length", strconv.Itoa(len(j)))
		/// http.ResponseWriter 生命周期原因，不支持异步
		_, err := req.w.Write(j)
		if err != nil {
			logs.LogError(err.Error())
		}
		// logs.LogError("uuid:%v %v", req.uuid, string(j))
	} else {
		resp = &Resp{}
		j, _ := json.Marshal(resp)
		setResponseHeader(req.w, req.r)
		req.w.Header().Set("Content-Type", "application/json")
		req.w.Header().Set("Content-Length", strconv.Itoa(len(j)))
		/// http.ResponseWriter 生命周期原因，不支持异步
		_, err := req.w.Write(j)
		if err != nil {
			logs.LogError(err.Error())
		}
		logs.LogFatal("uuid:%v", req.uuid)
	}
}
