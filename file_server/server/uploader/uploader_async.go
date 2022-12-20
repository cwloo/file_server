package uploader

import (
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/cwloo/gonet/core/base/cc"
	"github.com/cwloo/gonet/core/base/mq/lq"
	"github.com/cwloo/gonet/core/base/pipe"
	"github.com/cwloo/gonet/core/base/run"
	"github.com/cwloo/gonet/logs"
	"github.com/cwloo/uploader/file_server/config"
	"github.com/cwloo/uploader/file_server/global"
	"github.com/cwloo/uploader/file_server/httpsrv"
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
	flag     cc.AtomFlag
}

func NewAsyncUploader(uuid string) global.Uploader {
	s := asyncUploaders.Get().(*AsyncUploader)
	s.signaled = false
	s.uuid = uuid
	s.state = NewUploaderState()
	s.tm = time.Now()
	s.l_tm = &sync.RWMutex{}
	s.l_signal = &sync.Mutex{}
	s.cond = sync.NewCond(s.l_signal)
	s.flag = cc.NewAtomFlag()
	s.start()
	return s
}

func (s *AsyncUploader) start() {
	if s.pipe == nil && s.flag.TestSet() {
		mq := lq.NewQueue(1000)
		runner := NewProcessor(s.handler)
		s.pipe = pipe.NewPipeWithQuit(global.I32.New(), "uploader.pipe", mq, runner, s.onQuit)
		s.flag.Reset()
	}
	s.wait_started()
}

func (s *AsyncUploader) wait_started() {
	for {
		if s.pipe != nil {
			break
		}
	}
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
	if s.pipe.NotifyClose() {
		logs.Tracef("ok")
	} else {
		logs.Errorf("err")
		s.notify()
	}
}

func (s *AsyncUploader) Remove(md5 string) {
	if s.state.Remove(md5) && s.state.AllDone() {
		s.NotifyClose()
	}
}

func (s *AsyncUploader) Clear() {
	msgs := []string{}
	s.state.Range(func(md5 string, ok bool) {
		// if !ok {
		// 	////// 任务退出，移除未决的文件
		// 	if msg, ok := RemovePendingFile(s.uuid, md5); ok {
		// 		msgs = append(msgs, msg)
		// 	}
		// } else {
		// 	////// 任务退出，移除校验失败的文件
		// 	if msg, ok := RemoveCheckErrFile(s.uuid, md5); ok {
		// 		msgs = append(msgs, msg)
		// 	}
		// }
	})
	tg_bot.TgWarnMsg(msgs...)
}

func (s *AsyncUploader) onQuit(slot run.Slot) {
	s.Clear()
	global.Uploaders.Remove(s.uuid).Put()
}

func (s *AsyncUploader) handler(msg any, args ...any) (exit bool) {
	s.uploading(msg.(*global.Req))
	exit = s.state.AllDone()
	if exit {
		logs.Tracef("--------------------- ****** 无待上传文件，结束任务 %v ...", s.uuid)
	}
	s.notify()
	return
}

func (s *AsyncUploader) Upload(req *global.Req) {
	for _, key := range req.Keys {
		s.state.TryAdd(key.Md5)
	}
	s.start()
	s.pipe.Do(req)
	/// http.ResponseWriter 生命周期原因，不支持异步，所以加了 wait
	s.wait()
}

func (s *AsyncUploader) uploading(req *global.Req) {
	logs.Tracef("%#x", *req)
	s.update()
	resp := req.Resp
	result := req.Result
	for _, k := range req.Keys {
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
		if info.Done(false) {
			logs.Fatalf("%v %v[%v] %v %v/%v finished\nurl[%v]", info.Uuid(), info.SrcName(), info.Md5(), info.DstName(), info.Now(false), info.Total(false), info.Url(false))
		}
		////// 校验uuid
		if req.Uuid != info.Uuid() {
			logs.Fatalf("%v %v[%v] %v", info.Uuid(), info.SrcName(), info.Md5(), req.Uuid)
		}
		////// 校验MD5
		if k.Md5 != info.Md5() {
			logs.Fatalf("%v %v[%v] md5:%v", info.Uuid(), info.SrcName(), info.Md5(), k.Md5)
		}
		////// 校验数据大小
		if k.Total != strconv.FormatInt(info.Total(false), 10) {
			logs.Fatalf("%v %v[%v] info.total:%v total:%v", info.Uuid, info.SrcName, info.Md5, info.Total, k.Total)
		}
		////// 校验文件offset
		if k.Offset != strconv.FormatInt(info.Now(false), 10) {
			result = append(result,
				global.Result{
					Uuid:    info.Uuid(),
					File:    info.SrcName(),
					Md5:     info.Md5(),
					Now:     info.Now(false),
					Total:   info.Total(false),
					Expired: s.Get().Add(time.Duration(config.Config.Upload.PendingTimeout) * time.Second).Unix(),
					ErrCode: global.ErrCheckReUpload.ErrCode,
					ErrMsg:  global.ErrCheckReUpload.ErrMsg,
					Message: strings.Join([]string{info.Uuid(), " check reuploading ", info.DstName(), " progress:", strconv.FormatInt(info.Now(false), 10), "/", k.Total}, ""),
				})
			logs.Errorf("%v %v[%v] %v/%v offset:%v", info.Uuid(), info.SrcName(), info.Md5(), info.Now(false), info.Total(false), k.Offset)
			offset_n, _ := strconv.ParseInt(k.Offset, 10, 0)
			logs.Debugf("--------------------- ****** checking re-upload %v %v[%v] %v/%v offset:%v seg_size[%d]", info.Uuid(), header.Filename, k.Md5, info.Now(false), k.Total, offset_n, header.Size)
			continue
		}
		////// 检查上传目录
		_, err = os.Stat(config.Config.Upload.Dir)
		if err != nil && os.IsNotExist(err) {
			os.MkdirAll(config.Config.Upload.Dir, 0777)
		}
		////// 检查上传文件
		f := config.Config.Upload.Dir + info.DstName()
		switch config.Config.Upload.WriteFile > 0 {
		case true:
			_, err = os.Stat(f)
			if err != nil && os.IsNotExist(err) {
			} else {
				/// 第一次写如果文件存在则删除
				if info.Now(false) == int64(0) {
					os.Remove(f)
					logs.Fatalf("error")
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
						Expired: s.Get().Add(time.Duration(config.Config.Upload.PendingTimeout) * time.Second).Unix(),
						ErrCode: global.ErrCheckReUpload.ErrCode,
						ErrMsg:  global.ErrCheckReUpload.ErrMsg,
						Message: strings.Join([]string{info.Uuid(), " check reuploading ", info.DstName(), " progress:", strconv.FormatInt(info.Now(false), 10), "/", k.Total}, ""),
					})
				logs.Errorf(err.Error())
				offset_n, _ := strconv.ParseInt(k.Offset, 10, 0)
				logs.Debugf("--------------------- ****** checking re-upload %v %v[%v] %v/%v offset:%v seg_size[%d]", info.Uuid(), header.Filename, k.Md5, info.Now(false), k.Total, offset_n, header.Size)
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
						Expired: s.Get().Add(time.Duration(config.Config.Upload.PendingTimeout) * time.Second).Unix(),
						ErrCode: global.ErrCheckReUpload.ErrCode,
						ErrMsg:  global.ErrCheckReUpload.ErrMsg,
						Message: strings.Join([]string{info.Uuid(), " check reuploading ", info.DstName(), " progress:", strconv.FormatInt(info.Now(false), 10), "/", k.Total}, ""),
					})
				logs.Errorf(err.Error())
				err = fd.Close()
				if err != nil {
					logs.Errorf(err.Error())
				}
				offset_n, _ := strconv.ParseInt(k.Offset, 10, 0)
				logs.Debugf("--------------------- ****** checking re-upload %v %v[%v] %v/%v offset:%v seg_size[%d]", info.Uuid(), header.Filename, k.Md5, info.Now(false), k.Total, offset_n, header.Size)
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
				switch config.Config.Upload.WriteFile > 0 {
				case true:
					switch config.Config.Upload.CheckMd5 > 0 {
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
				logs.Debugf("%v %v[%v] %v ==>>> %v/%v +%v last_segment[finished] checking md5 ...", s.uuid, header.Filename, k.Md5, info.DstName(), info.Now(false), k.Total, header.Size)
				if ok {
					// global.FileInfos.Remove(info.Md5()).Put()
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
							Message: strings.Join([]string{info.Uuid(), " uploading ", info.DstName(), " progress:", strconv.FormatInt(info.Now(false), 10) + "/" + k.Total + " 上传成功!"}, "")})
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
							Now:     info.Now(false),
							Total:   info.Total(false),
							ErrCode: global.ErrFileMd5.ErrCode,
							ErrMsg:  global.ErrFileMd5.ErrMsg,
							Message: strings.Join([]string{info.Uuid(), " uploading ", info.DstName(), " progress:", strconv.FormatInt(info.Now(false), 10) + "/" + k.Total + " 上传完毕 MD5校验失败!"}, "")})
					logs.Errorf("%v %v[%v] %v chkmd5 [Err] elapsed:%vms", req.Uuid, header.Filename, k.Md5, info.DstName(), time.Since(start).Milliseconds())
					tg_bot.TgErrMsg(fmt.Sprintf("%v\n%v[%v]\n%v chkmd5 [Err] elapsed:%vms", req.Uuid, header.Filename, k.Md5, info.DstName(), time.Since(start).Milliseconds()))
				}
			} else {
				result = append(result,
					global.Result{
						Uuid:    req.Uuid,
						File:    header.Filename,
						Md5:     info.Md5(),
						Now:     info.Now(false),
						Total:   info.Total(false),
						Expired: s.Get().Add(time.Duration(config.Config.Upload.PendingTimeout) * time.Second).Unix(),
						ErrCode: global.ErrSegOk.ErrCode,
						ErrMsg:  global.ErrSegOk.ErrMsg,
						Message: strings.Join([]string{info.Uuid(), " uploading ", info.DstName(), " progress:", strconv.FormatInt(info.Now(false), 10) + "/" + k.Total}, "")})
				if info.Now(false) == header.Size {
					logs.Tracef("%v %v[%v] %v ==>>> %v/%v +%v first_segment", req.Uuid, header.Filename, k.Md5, info.DstName(), info.Now(false), k.Total, header.Size)
				} else {
					logs.Warnf("%v %v[%v] %v ==>>> %v/%v +%v continue_segment", req.Uuid, header.Filename, k.Md5, info.DstName(), info.Now(false), k.Total, header.Size)
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
							Now:     info.Now(false),
							Total:   info.Total(false),
							ErrCode: global.ErrFileMd5.ErrCode,
							ErrMsg:  global.ErrFileMd5.ErrMsg,
							Message: strings.Join([]string{info.Uuid(), " uploading ", info.DstName(), " progress:", strconv.FormatInt(info.Now(false), 10) + "/" + k.Total + " 上传完毕 MD5校验失败!"}, "")})
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
						Now:     info.Now(false),
						Total:   info.Total(false),
						ErrCode: global.ErrFileMd5.ErrCode,
						ErrMsg:  global.ErrFileMd5.ErrMsg,
						Message: strings.Join([]string{info.Uuid(), " uploading ", info.DstName(), " progress:", strconv.FormatInt(info.Now(false), 10) + "/" + k.Total + " 上传完毕 MD5校验失败!"}, "")})
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
		/// http.ResponseWriter 生命周期原因，不支持异步，所以加了 notify
		httpsrv.WriteResponse(req.W, req.R, resp)
		// logs.Errorf("%v %v", req.Uuid, string(j))
	} else {
		/// http.ResponseWriter 生命周期原因，不支持异步，所以加了 notify
		httpsrv.WriteResponse(req.W, req.R, &global.Resp{})
		logs.Fatalf("%v", req.Uuid)
	}
}
