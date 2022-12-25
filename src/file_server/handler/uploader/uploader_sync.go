package uploader

import (
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/cwloo/gonet/logs"
	"github.com/cwloo/uploader/src/config"
	"github.com/cwloo/uploader/src/file_server/handler"
	"github.com/cwloo/uploader/src/global"
	"github.com/cwloo/uploader/src/global/httpsrv"
	"github.com/cwloo/uploader/src/global/tg_bot"
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
	uuid  string
	state State
	tm    time.Time
	l_tm  *sync.RWMutex
}

func NewSyncUploader(uuid string) global.Uploader {
	s := syncUploaders.Get().(*SyncUploader)
	s.uuid = uuid
	s.state = NewUploaderState()
	s.tm = time.Now()
	s.l_tm = &sync.RWMutex{}
	return s
}

func (s *SyncUploader) reset() {
	s.state.Put()
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
	global.Uploaders.Remove(s.uuid).Put()
}

func (s *SyncUploader) NotifyClose() {
}

func (s *SyncUploader) Remove(md5 string) {
	if s.state.Remove(md5) && s.state.AllDone() {
		s.Clear()
		global.Uploaders.Remove(s.uuid).Put()
	}
}

func (s *SyncUploader) Clear() {
	msgs := []string{}
	s.state.Range(func(md5 string, ok bool) {
		if !ok {
			////// 任务退出，移除未决的文件
			if msg, ok := handler.RemovePendingFile(s.uuid, md5); ok {
				msgs = append(msgs, msg)
			}
		} else {
			////// 任务退出，移除校验失败的文件
			if msg, ok := handler.RemoveCheckErrFile(s.uuid, md5); ok {
				msgs = append(msgs, msg)
			}
		}
	})
	tg_bot.TgWarnMsg(msgs...)
}

func (s *SyncUploader) Upload(req *global.Req) {
	for _, key := range req.Keys {
		s.state.TryAdd(key.Md5)
	}
	s.uploading(req)
	exit := s.state.AllDone()
	if exit {
		logs.Tracef("--------------------- ****** 无待上传文件，结束任务 %v ...", s.uuid)
		global.Uploaders.Remove(s.uuid).Put()
	}
}

func (s *SyncUploader) uploading(req *global.Req) {
	logs.Tracef("%#x", *req)
	s.update()
	resp := req.Resp
	result := req.Result
	for _, k := range req.Keys {
		s.state.TryAdd(k.Md5)
		part, header, err := req.R.FormFile(k.Key)
		if err != nil {
			logs.Errorf(err.Error())
			size, _ := strconv.ParseInt(k.Total, 10, 0)
			result = append(result,
				global.Result{
					Uuid:    req.Uuid,
					File:    k.Filename,
					Md5:     k.Md5,
					Total:   size,
					ErrCode: global.ErrCheckReUpload.ErrCode,
					ErrMsg:  global.ErrCheckReUpload.ErrMsg,
					Message: strings.Join([]string{req.Uuid, " check reuploading ", k.Filename, " progress:", strconv.FormatInt(0, 10), "/", k.Total}, ""),
				})
			logs.Errorf("%v %v[%v] %v/%v offset:%v", req.Uuid, k.Filename, k.Md5, 0, k.Total, k.Offset)
			offset_n, _ := strconv.ParseInt(k.Offset, 10, 0)
			logs.Debugf("--------------------- ****** checking re-upload %v %v[%v] %v/%v offset:%v seg_size[%d]", req.Uuid, k.Filename, k.Md5, 0, k.Total, offset_n, k.Headersize)
			continue
		}
		info := global.FileInfos.Get(k.Md5)
		if info == nil {
			size, _ := strconv.ParseInt(k.Total, 10, 0)
			result = append(result,
				global.Result{
					Uuid:    req.Uuid,
					File:    k.Filename,
					Md5:     k.Md5,
					Total:   size,
					ErrCode: global.ErrCheckReUpload.ErrCode,
					ErrMsg:  global.ErrCheckReUpload.ErrMsg,
					Message: strings.Join([]string{req.Uuid, " check reuploading ", header.Filename, " progress:", strconv.FormatInt(0, 10), "/", k.Total}, ""),
				})
			logs.Errorf("%v %v[%v] %v/%v offset:%v", req.Uuid, header.Filename, k.Md5, 0, k.Total, k.Offset)
			offset_n, _ := strconv.ParseInt(k.Offset, 10, 0)
			logs.Debugf("--------------------- ****** checking re-upload %v %v[%v] %v/%v offset:%v seg_size[%d]", req.Uuid, header.Filename, k.Md5, 0, k.Total, offset_n, header.Size)
			continue
		}
		////// 还未接收完
		if info.Done(true) {
			logs.Fatalf("%v %v[%v] %v %v/%v finished\nurl[%v]", info.Uuid(), info.SrcName(), info.Md5(), info.DstName(), info.Now(true), info.Total(false), info.Url(false))
		}
		////// 校验uuid
		if req.Uuid != info.Uuid() {
			logs.Fatalf("%v %v(%v) %v", info.Uuid(), info.SrcName(), info.Md5(), req.Uuid)
		}
		////// 校验MD5
		if k.Md5 != info.Md5() {
			logs.Fatalf("%v %v(%v) md5:%v", info.Uuid(), info.SrcName(), info.Md5(), k.Md5)
		}
		////// 校验数据大小
		if k.Total != strconv.FormatInt(info.Total(false), 10) {
			logs.Fatalf("%v %v(%v) info.total:%v total:%v", info.Uuid(), info.SrcName(), info.Md5(), info.Total(false), k.Total)
		}
		////// 校验文件offset
		if k.Offset != strconv.FormatInt(info.Now(true), 10) {
			result = append(result,
				global.Result{
					Uuid:    info.Uuid(),
					File:    info.SrcName(),
					Md5:     info.Md5(),
					Now:     info.Now(true),
					Total:   info.Total(false),
					Expired: s.Get().Add(time.Duration(config.Config.File.Upload.PendingTimeout) * time.Second).Unix(),
					ErrCode: global.ErrCheckReUpload.ErrCode,
					ErrMsg:  global.ErrCheckReUpload.ErrMsg,
					Message: strings.Join([]string{info.Uuid(), " check reuploading ", info.DstName(), " progress:", strconv.FormatInt(info.Now(true), 10), "/", k.Total}, ""),
				})
			// logs.Errorf("%v %v(%v) %v/%v offset:%v", info.Uuid(), info.SrcName(), info.Md5(), info.Now(true), info.Total(false), k.Offset)
			offset_n, _ := strconv.ParseInt(k.Offset, 10, 0)
			logs.Infof("--------------------- checking re-upload %v %v[%v] %v/%v offset:%v seg_size[%d]", info.Uuid(), header.Filename, k.Md5, info.Now(true), k.Total, offset_n, header.Size)
			continue
		}
		f := config.Config.File.Upload.Dir + info.DstName()
		switch config.Config.File.Upload.WriteFile > 0 {
		case true:
			////// 检查上传目录
			_, err = os.Stat(config.Config.File.Upload.Dir)
			if err != nil && os.IsNotExist(err) {
				os.MkdirAll(config.Config.File.Upload.Dir, 0777)
			}
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
						Expired: s.Get().Add(time.Duration(config.Config.File.Upload.PendingTimeout) * time.Second).Unix(),
						ErrCode: global.ErrCheckReUpload.ErrCode,
						ErrMsg:  global.ErrCheckReUpload.ErrMsg,
						Message: strings.Join([]string{info.Uuid(), " check reuploading ", info.DstName(), " progress:", strconv.FormatInt(info.Now(true), 10), "/", k.Total}, ""),
					})
				logs.Errorf(err.Error())
				offset_n, _ := strconv.ParseInt(k.Offset, 10, 0)
				logs.Infof("--------------------- checking re-upload %v %v[%v] %v/%v offset:%v seg_size[%d]", info.Uuid(), header.Filename, k.Md5, info.Now(true), k.Total, offset_n, header.Size)
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
						Expired: s.Get().Add(time.Duration(config.Config.File.Upload.PendingTimeout) * time.Second).Unix(),
						ErrCode: global.ErrCheckReUpload.ErrCode,
						ErrMsg:  global.ErrCheckReUpload.ErrMsg,
						Message: strings.Join([]string{info.Uuid(), " check reuploading ", info.DstName(), " progress:", strconv.FormatInt(info.Now(true), 10), "/", k.Total}, ""),
					})
				logs.Errorf(err.Error())
				err = fd.Close()
				if err != nil {
					logs.Errorf(err.Error())
				}
				offset_n, _ := strconv.ParseInt(k.Offset, 10, 0)
				logs.Infof("--------------------- checking re-upload %v %v[%v] %v/%v offset:%v seg_size[%d]", info.Uuid(), header.Filename, k.Md5, info.Now(true), k.Total, offset_n, header.Size)
				continue
			}
			err = fd.Close()
			if err != nil {
				logs.Errorf(err.Error())
			}
			err = part.Close()
			if err != nil {
				logs.Errorf(err.Error())
			}
		default:
		}
		retry_c := 0
	retry:
		done, ok, url, errMsg, start := info.Update(header.Size,
			config.Config.Oss.Type,
			func(info global.FileInfo, oss global.Oss) (url string, err *global.ErrorMsg) {
				url, _, err = oss.UploadFile(info, header)
				if err != nil {
					logs.Errorf(err.Error())
				}
				return
			}, func(info global.FileInfo) (time.Time, bool) {
				start := time.Now()
				switch config.Config.File.Upload.WriteFile > 0 {
				case true:
					switch config.Config.File.Upload.CheckMd5 > 0 {
					case true:
						md5 := global.CalcFileMd5(f)
						ok := md5 == info.Md5()
						return start, ok
					default:
						return start, true
					}
				default:
					return start, true
				}
			})
		switch errMsg {
		case nil:
			if done {
				s.state.SetDone(info.Md5())
				logs.Debugf("%v %v[%v] %v ==>>> %v/%v +%v last_segment[finished] checking md5 ...", s.uuid, header.Filename, k.Md5, info.DstName(), info.Now(true), k.Total, header.Size)
				if ok {
					// global.FileInfos.Remove(info.Md5()).Put()
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
							Message: strings.Join([]string{info.Uuid(), " uploading ", info.DstName(), " progress:", strconv.FormatInt(info.Now(true), 10) + "/" + k.Total + " 上传成功!"}, "")})
					logs.Warnf("%v %v[%v] %v chkmd5 [ok] %v elapsed:%vms", req.Uuid, header.Filename, k.Md5, info.DstName(), url, time.Since(start).Milliseconds())
					tg_bot.TgSuccMsg(fmt.Sprintf("%v\n%v[%v]\n%v chkmd5 [ok]\n%v elapsed:%vms", req.Uuid, header.Filename, k.Md5, info.DstName(), url, time.Since(start).Milliseconds()))
				} else {
					global.FileInfos.Remove(info.Md5()).Put()
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
							Message: strings.Join([]string{info.Uuid(), " uploading ", info.DstName(), " progress:", strconv.FormatInt(info.Now(true), 10) + "/" + k.Total + " 上传完毕 MD5校验失败!"}, "")})
					logs.Errorf("%v %v[%v] %v chkmd5 [Err] elapsed:%vms", req.Uuid, header.Filename, k.Md5, info.DstName(), time.Since(start).Milliseconds())
					tg_bot.TgErrMsg(fmt.Sprintf("%v\n%v[%v]\n%v chkmd5 [Err] elapsed:%vms", req.Uuid, header.Filename, k.Md5, info.DstName(), time.Since(start).Milliseconds()))
				}
			} else {
				result = append(result,
					global.Result{
						Uuid:    req.Uuid,
						File:    header.Filename,
						Md5:     info.Md5(),
						Now:     info.Now(true),
						Total:   info.Total(false),
						Expired: s.Get().Add(time.Duration(config.Config.File.Upload.PendingTimeout) * time.Second).Unix(),
						ErrCode: global.ErrSegOk.ErrCode,
						ErrMsg:  global.ErrSegOk.ErrMsg,
						Message: strings.Join([]string{info.Uuid(), " uploading ", info.DstName(), " progress:", strconv.FormatInt(info.Now(true), 10) + "/" + k.Total}, "")})
				if info.Now(true) == header.Size {
					logs.Tracef("%v %v[%v] %v ==>>> %v/%v +%v first_segment", req.Uuid, header.Filename, k.Md5, info.DstName(), info.Now(true), k.Total, header.Size)
				} else {
					logs.Warnf("%v %v[%v] %v ==>>> %v/%v +%v continue_segment", req.Uuid, header.Filename, k.Md5, info.DstName(), info.Now(true), k.Total, header.Size)
				}
			}
		default:
			switch errMsg.ErrCode {
			case global.ErrCancel.ErrCode:
				global.FileInfos.Remove(info.Md5()).Put()
				os.Remove(f)
				size, _ := strconv.ParseInt(k.Total, 10, 0)
				result = append(result,
					global.Result{
						Uuid:    req.Uuid,
						File:    k.Filename,
						Md5:     k.Md5,
						Total:   size,
						ErrCode: global.ErrCheckReUpload.ErrCode,
						ErrMsg:  global.ErrCheckReUpload.ErrMsg,
						Message: strings.Join([]string{req.Uuid, " check reuploading ", header.Filename, " progress:", strconv.FormatInt(0, 10), "/", k.Total}, ""),
					})
				logs.Errorf("%v %v[%v] %v/%v offset:%v", req.Uuid, header.Filename, k.Md5, 0, k.Total, k.Offset)
				offset_n, _ := strconv.ParseInt(k.Offset, 10, 0)
				logs.Debugf("--------------------- ****** checking re-upload %v %v[%v] %v/%v offset:%v seg_size[%d]", req.Uuid, header.Filename, k.Md5, 0, k.Total, offset_n, header.Size)
			case global.ErrRetry.ErrCode:
				retry_c++
				switch retry_c <= 2 {
				case true:
					goto retry
				default:
					global.FileInfos.Remove(info.Md5()).Put()
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
							Message: strings.Join([]string{info.Uuid(), " uploading ", info.DstName(), " progress:", strconv.FormatInt(info.Now(true), 10) + "/" + k.Total + " 上传完毕 MD5校验失败!"}, "")})
					logs.Errorf("%v %v[%v] %v chkmd5 [Err] elapsed:%vms", req.Uuid, header.Filename, k.Md5, info.DstName(), time.Since(start).Milliseconds())
				}
			case global.ErrFatal.ErrCode:
				global.FileInfos.Remove(info.Md5()).Put()
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
						Message: strings.Join([]string{info.Uuid(), " uploading ", info.DstName(), " progress:", strconv.FormatInt(info.Now(true), 10) + "/" + k.Total + " 上传完毕 MD5校验失败!"}, "")})
				logs.Errorf("%v %v[%v] %v chkmd5 [Err] elapsed:%vms", req.Uuid, header.Filename, k.Md5, info.DstName(), time.Since(start).Milliseconds())
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
		logs.Tracef("httpsrv.WriteResponse")
		/// http.ResponseWriter 生命周期原因，不支持异步
		httpsrv.WriteResponse(req.W, req.R, resp)
		// logs.Errorf("%v %v", req.Uuid, string(j))
	} else {
		/// http.ResponseWriter 生命周期原因，不支持异步
		httpsrv.WriteResponse(req.W, req.R, &global.Resp{})
		logs.Fatalf("%v", req.Uuid)
	}
}
