package main

import (
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/cwloo/gonet/logs"
	"github.com/cwloo/uploader/file_server/config"
	"github.com/cwloo/uploader/file_server/global"
	"github.com/cwloo/uploader/file_server/tg_bot"
)

var (
	syncUploaders = sync.Pool{
		New: func() any {
			return &SyncUploader{}
		},
	}
)

// <summary>
// SyncUploader 同步方式上传
// <summary>
type SyncUploader struct {
	uuid string
	data Data
	tm   time.Time
	l_tm *sync.RWMutex
}

func NewSyncUploader(uuid string) Uploader {
	s := syncUploaders.Get().(*SyncUploader)
	s.uuid = uuid
	s.data = NewUploaderData()
	s.tm = time.Now()
	s.l_tm = &sync.RWMutex{}
	return s
}

func (s *SyncUploader) reset() {
	s.data.Put()
}

func (s *SyncUploader) Put() {
	s.reset()
	syncUploaders.Put(s)
}

func (s *SyncUploader) update() {
	s.l_tm.Lock()
	s.tm = time.Now()
	s.l_tm.Unlock()
}

func (s *SyncUploader) Get() time.Time {
	s.l_tm.RLock()
	tm := s.tm
	s.l_tm.RUnlock()
	return tm
}

func (s *SyncUploader) Close() {
	s.Clear()
	uploaders.Remove(s.uuid).Put()
}

func (s *SyncUploader) NotifyClose() {
}

func (s *SyncUploader) Remove(md5 string) {
	if s.data.Remove(md5) && s.data.AllDone() {
		uploaders.Remove(s.uuid).Put()
	}
}

func (s *SyncUploader) Clear() {
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
				os.Remove(config.Config.UploadlDir + info.DstName())
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
				os.Remove(config.Config.UploadlDir + info.DstName())
				info.Put()
			})
		}
	})
	tg_bot.TgWarnMsg(msgs...)
}

func (s *SyncUploader) Upload(req *global.Req) {
	switch config.Config.MultiFile {
	default:
		s.multi_uploading(req)
	case 0:
		s.uploading(req)
	}
	exit := s.data.AllDone()
	if exit {
		logs.LogTrace("--------------------- ****** 无待上传文件，结束任务 %v ...", s.uuid)
		uploaders.Remove(s.uuid).Put()
	}
}

func (s *SyncUploader) uploading(req *global.Req) {
	s.update()
	resp := req.Resp
	result := req.Result
	for _, k := range req.Keys {
		s.data.TryAdd(req.Md5)
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
		if info.Done(true) {
			logs.LogFatal("%v %v(%v) finished", info.Uuid(), info.SrcName(), info.Md5())
		}
		////// 校验uuid
		if req.Uuid != info.Uuid() {
			logs.LogFatal("%v %v(%v) %v", info.Uuid(), info.SrcName(), info.Md5(), req.Uuid)
		}
		////// 校验MD5
		if req.Md5 != info.Md5() {
			logs.LogFatal("%v %v(%v) md5:%v", info.Uuid(), info.SrcName(), info.Md5(), req.Md5)
		}
		////// 校验数据大小
		if req.Total != strconv.FormatInt(info.Total(false), 10) {
			logs.LogFatal("%v %v(%v) info.total:%v total:%v", info.Uuid(), info.SrcName(), info.Md5(), info.Total(false), req.Total)
		}
		////// 校验文件offset
		if req.Offset != strconv.FormatInt(info.Now(true), 10) {
			result = append(result,
				global.Result{
					Uuid:    info.Uuid(),
					File:    info.SrcName(),
					Md5:     info.Md5(),
					Now:     info.Now(true),
					Total:   info.Total(false),
					Expired: s.Get().Add(time.Duration(config.Config.PendingTimeout) * time.Second).Unix(),
					ErrCode: global.ErrCheckReUpload.ErrCode,
					ErrMsg:  global.ErrCheckReUpload.ErrMsg,
					Message: strings.Join([]string{info.Uuid(), " check reuploading ", info.DstName(), " progress:", strconv.FormatInt(info.Now(true), 10), "/", req.Total}, ""),
				})
			// logs.LogError("%v %v(%v) %v/%v offset:%v", info.Uuid(), info.SrcName(), info.Md5(), info.Now(true), info.Total(false), req.Offset)
			offset_n, _ := strconv.ParseInt(req.Offset, 10, 0)
			logs.LogInfo("--------------------- checking re-upload %v %v[%v] %v/%v offset:%v seg_size[%d]", info.Uuid(), header.Filename, req.Md5, info.Now(true), req.Total, offset_n, header.Size)
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
				if info.Now(true) == int64(0) {
					os.Remove(f)
				}
			}
			fd, err := os.OpenFile(f, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0777)
			if err != nil {
				result = append(result,
					global.Result{
						Uuid:    info.Uuid(),
						File:    info.SrcName(),
						Md5:     info.Md5(),
						Now:     info.Now(true),
						Total:   info.Total(false),
						Expired: s.Get().Add(time.Duration(config.Config.PendingTimeout) * time.Second).Unix(),
						ErrCode: global.ErrCheckReUpload.ErrCode,
						ErrMsg:  global.ErrCheckReUpload.ErrMsg,
						Message: strings.Join([]string{info.Uuid(), " check reuploading ", info.DstName(), " progress:", strconv.FormatInt(info.Now(true), 10), "/", req.Total}, ""),
					})
				logs.LogError(err.Error())
				offset_n, _ := strconv.ParseInt(req.Offset, 10, 0)
				logs.LogInfo("--------------------- checking re-upload %v %v[%v] %v/%v offset:%v seg_size[%d]", info.Uuid(), header.Filename, req.Md5, info.Now(true), req.Total, offset_n, header.Size)
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
						Now:     info.Now(true),
						Total:   info.Total(false),
						Expired: s.Get().Add(time.Duration(config.Config.PendingTimeout) * time.Second).Unix(),
						ErrCode: global.ErrCheckReUpload.ErrCode,
						ErrMsg:  global.ErrCheckReUpload.ErrMsg,
						Message: strings.Join([]string{info.Uuid(), " check reuploading ", info.DstName(), " progress:", strconv.FormatInt(info.Now(true), 10), "/", req.Total}, ""),
					})
				logs.LogError(err.Error())
				err = fd.Close()
				if err != nil {
					logs.LogError(err.Error())
				}
				offset_n, _ := strconv.ParseInt(req.Offset, 10, 0)
				logs.LogInfo("--------------------- checking re-upload %v %v[%v] %v/%v offset:%v seg_size[%d]", info.Uuid(), header.Filename, req.Md5, info.Now(true), req.Total, offset_n, header.Size)
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
			logs.LogDebug("%v %v[%v] %v ==>>> %v/%v +%v last_segment[finished] checking md5 ...", s.uuid, header.Filename, req.Md5, info.DstName(), info.Now(true), req.Total, header.Size)
			if ok {
				// fileInfos.Remove(info.Md5()).Put()
				result = append(result,
					global.Result{
						Uuid:    req.Uuid,
						File:    header.Filename,
						Md5:     info.Md5(),
						Now:     info.Now(true),
						Total:   info.Total(false),
						ErrCode: global.ErrOk.ErrCode,
						ErrMsg:  global.ErrOk.ErrMsg,
						Url:     url,
						Message: strings.Join([]string{info.Uuid(), " uploading ", info.DstName(), " progress:", strconv.FormatInt(info.Now(true), 10) + "/" + req.Total + " 上传成功!"}, "")})
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
						Now:     info.Now(true),
						Total:   info.Total(false),
						ErrCode: global.ErrFileMd5.ErrCode,
						ErrMsg:  global.ErrFileMd5.ErrMsg,
						Message: strings.Join([]string{info.Uuid(), " uploading ", info.DstName(), " progress:", strconv.FormatInt(info.Now(true), 10) + "/" + req.Total + " 上传完毕 MD5校验失败!"}, "")})
				logs.LogError("%v %v[%v] %v chkmd5 [Err] elapsed:%vms", req.Uuid, header.Filename, req.Md5, info.DstName(), time.Since(start).Milliseconds())
				tg_bot.TgErrMsg(fmt.Sprintf("%v\n%v[%v]\n%v chkmd5 [Err] elapsed:%vms", req.Uuid, header.Filename, req.Md5, info.DstName(), time.Since(start).Milliseconds()))
			}
		} else {
			result = append(result,
				global.Result{
					Uuid:    req.Uuid,
					File:    header.Filename,
					Md5:     info.Md5(),
					Now:     info.Now(true),
					Total:   info.Total(false),
					Expired: s.Get().Add(time.Duration(config.Config.PendingTimeout) * time.Second).Unix(),
					ErrCode: global.ErrSegOk.ErrCode,
					ErrMsg:  global.ErrSegOk.ErrMsg,
					Message: strings.Join([]string{info.Uuid(), " uploading ", info.DstName(), " progress:", strconv.FormatInt(info.Now(true), 10) + "/" + req.Total}, "")})
			if info.Now(true) == header.Size {
				logs.LogTrace("%v %v[%v] %v ==>>> %v/%v +%v first_segment", req.Uuid, header.Filename, req.Md5, info.DstName(), info.Now(true), req.Total, header.Size)
			} else {
				logs.LogWarn("%v %v[%v] %v ==>>> %v/%v +%v continue_segment", req.Uuid, header.Filename, req.Md5, info.DstName(), info.Now(true), req.Total, header.Size)
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
		/// http.ResponseWriter 生命周期原因，不支持异步
		writeResponse(req.W, req.R, resp)
		// logs.LogError("%v %v", req.Uuid, string(j))
	} else {
		/// http.ResponseWriter 生命周期原因，不支持异步
		writeResponse(req.W, req.R, &global.Resp{})
		logs.LogFatal("%v", req.Uuid)
	}
}

func (s *SyncUploader) multi_uploading(req *global.Req) {
	s.update()
	resp := req.Resp
	result := req.Result
	for _, k := range req.Keys {
		offset := req.R.FormValue(k + ".offset")
		total := req.R.FormValue(k + ".total")
		md5 := strings.ToLower(k)
		s.data.TryAdd(md5)
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
		if info.Done(true) {
			logs.LogFatal("%v %v(%v) finished", info.Uuid(), info.SrcName(), info.Md5())
		}
		////// 校验uuid
		if req.Uuid != info.Uuid() {
			logs.LogFatal("%v %v(%v) %v", info.Uuid(), info.SrcName(), info.Md5(), req.Uuid)
		}
		////// 校验MD5
		if md5 != info.Md5() {
			logs.LogFatal("%v %v(%v) md5:%v", info.Uuid(), info.SrcName(), info.Md5(), md5)
		}
		////// 校验数据大小
		if total != strconv.FormatInt(info.Total(false), 10) {
			logs.LogFatal("%v %v(%v) info.total:%v total:%v", info.Uuid(), info.SrcName(), info.Md5(), info.Total(false), total)
		}
		////// 校验文件offset
		if offset != strconv.FormatInt(info.Now(true), 10) {
			result = append(result,
				global.Result{
					Uuid:    info.Uuid(),
					File:    info.SrcName(),
					Md5:     info.Md5(),
					Now:     info.Now(true),
					Total:   info.Total(false),
					Expired: s.Get().Add(time.Duration(config.Config.PendingTimeout) * time.Second).Unix(),
					ErrCode: global.ErrCheckReUpload.ErrCode,
					ErrMsg:  global.ErrCheckReUpload.ErrMsg,
					Message: strings.Join([]string{info.Uuid(), " check reuploading ", info.DstName(), " progress:", strconv.FormatInt(info.Now(true), 10), "/", total}, ""),
				})
			// logs.LogError("%v %v(%v) %v/%v offset:%v", info.Uuid(), info.SrcName(), info.Md5(), info.Now(true), info.Total(false), offset)
			offset_n, _ := strconv.ParseInt(offset, 10, 0)
			logs.LogInfo("--------------------- checking re-upload %v %v[%v] %v/%v offset:%v seg_size[%d]", info.Uuid(), header.Filename, md5, info.Now(true), total, offset_n, header.Size)
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
				if info.Now(true) == int64(0) {
					os.Remove(f)
				}
			}
			fd, err := os.OpenFile(f, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0777)
			if err != nil {
				result = append(result,
					global.Result{
						Uuid:    info.Uuid(),
						File:    info.SrcName(),
						Md5:     info.Md5(),
						Now:     info.Now(true),
						Total:   info.Total(false),
						Expired: s.Get().Add(time.Duration(config.Config.PendingTimeout) * time.Second).Unix(),
						ErrCode: global.ErrCheckReUpload.ErrCode,
						ErrMsg:  global.ErrCheckReUpload.ErrMsg,
						Message: strings.Join([]string{info.Uuid(), " check reuploading ", info.DstName(), " progress:", strconv.FormatInt(info.Now(true), 10), "/", total}, ""),
					})
				logs.LogError(err.Error())
				offset_n, _ := strconv.ParseInt(offset, 10, 0)
				logs.LogInfo("--------------------- checking re-upload %v %v[%v] %v/%v offset:%v seg_size[%d]", info.Uuid(), header.Filename, md5, info.Now(true), total, offset_n, header.Size)
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
						Now:     info.Now(true),
						Total:   info.Total(false),
						Expired: s.Get().Add(time.Duration(config.Config.PendingTimeout) * time.Second).Unix(),
						ErrCode: global.ErrCheckReUpload.ErrCode,
						ErrMsg:  global.ErrCheckReUpload.ErrMsg,
						Message: strings.Join([]string{info.Uuid(), " check reuploading ", info.DstName(), " progress:", strconv.FormatInt(info.Now(true), 10), "/", total}, ""),
					})
				logs.LogError(err.Error())
				err = fd.Close()
				if err != nil {
					logs.LogError(err.Error())
				}
				offset_n, _ := strconv.ParseInt(offset, 10, 0)
				logs.LogInfo("--------------------- checking re-upload %v %v[%v] %v/%v offset:%v seg_size[%d]", info.Uuid(), header.Filename, md5, info.Now(true), total, offset_n, header.Size)
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
					global.Result{
						Uuid:    req.Uuid,
						File:    header.Filename,
						Md5:     info.Md5(),
						Now:     info.Now(true),
						Total:   info.Total(false),
						ErrCode: global.ErrOk.ErrCode,
						ErrMsg:  global.ErrOk.ErrMsg,
						Url:     url,
						Message: strings.Join([]string{info.Uuid(), " uploading ", info.DstName(), " progress:", strconv.FormatInt(info.Now(true), 10) + "/" + total + " 上传成功!"}, "")})
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
						Now:     info.Now(true),
						Total:   info.Total(false),
						ErrCode: global.ErrFileMd5.ErrCode,
						ErrMsg:  global.ErrFileMd5.ErrMsg,
						Message: strings.Join([]string{info.Uuid(), " uploading ", info.DstName(), " progress:", strconv.FormatInt(info.Now(true), 10) + "/" + total + " 上传完毕 MD5校验失败!"}, "")})
				logs.LogError("%v %v[%v] %v chkmd5 [Err] elapsed:%vms", req.Uuid, header.Filename, md5, info.DstName(), time.Since(start).Milliseconds())
				tg_bot.TgErrMsg(fmt.Sprintf("%v\n%v[%v]\n%v chkmd5 [Err] elapsed:%vms", req.Uuid, header.Filename, md5, info.DstName(), time.Since(start).Milliseconds()))
			}
		} else {
			result = append(result,
				global.Result{
					Uuid:    req.Uuid,
					File:    header.Filename,
					Md5:     info.Md5(),
					Now:     info.Now(true),
					Total:   info.Total(false),
					Expired: s.Get().Add(time.Duration(config.Config.PendingTimeout) * time.Second).Unix(),
					ErrCode: global.ErrSegOk.ErrCode,
					ErrMsg:  global.ErrSegOk.ErrMsg,
					Message: strings.Join([]string{info.Uuid(), " uploading ", info.DstName(), " progress:", strconv.FormatInt(info.Now(true), 10) + "/" + total}, "")})
			if info.Now(true) == header.Size {
				logs.LogTrace("%v %v[%v] %v ==>>> %v/%v +%v first_segment", req.Uuid, header.Filename, md5, info.DstName(), info.Now(true), total, header.Size)
			} else {
				logs.LogWarn("%v %v[%v] %v ==>>> %v/%v +%v continue_segment", req.Uuid, header.Filename, md5, info.DstName(), info.Now(true), total, header.Size)
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
		/// http.ResponseWriter 生命周期原因，不支持异步
		writeResponse(req.W, req.R, resp)
		// logs.LogError("%v %v", req.Uuid, string(j))
	} else {
		/// http.ResponseWriter 生命周期原因，不支持异步
		writeResponse(req.W, req.R, &global.Resp{})
		logs.LogFatal("%v", req.Uuid)
	}
}
