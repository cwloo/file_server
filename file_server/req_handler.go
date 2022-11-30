package main

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"github.com/cwloo/gonet/logs"
	"github.com/cwloo/gonet/utils"
)

func handlerUploadFile(w http.ResponseWriter, r *http.Request) {
	err := r.ParseMultipartForm(MaxMemory)
	if err != nil {
		w.WriteHeader(http.StatusOK)
		obj := &Resp{
			ErrCode: ErrParsePartData.ErrCode,
			ErrMsg:  ErrParsePartData.ErrMsg,
		}
		j, _ := json.Marshal(obj)
		w.Write(j)
		logs.LogError("%v", err.Error())
		return
	}
	uuid := ""
	form := r.MultipartForm
	for k := range form.Value {
		switch k {
		case "uuid":
			uuid = r.FormValue(k)
		}
		// logs.LogTrace("%v=%v", k, v)
	}
	if !checkUUID(uuid) {
		w.WriteHeader(http.StatusOK)
		obj := &Resp{
			ErrCode: ErrParamsUUID.ErrCode,
			ErrMsg:  ErrParamsUUID.ErrMsg,
		}
		j, _ := json.Marshal(obj)
		w.Write(j)
		logs.LogError("uuid=%v", uuid)
		return
	}
	uploader := uploaders.Get(uuid)
	if uploader == nil {
		allTotal := int64(0)
		///////////////////////////// 新的上传任务 /////////////////////////////
		keys := []string{}
		ignore := []*FileInfo{}
		for k := range form.File {
			total := r.FormValue(k + ".total")
			md5 := strings.ToLower(r.FormValue(k + ".md5"))
			/// header检查
			_, header, err := r.FormFile(k)
			if err != nil {
				logs.LogError("%v", err.Error())
				return
			}
			/// header.size检查
			if !checkMultiPartFileHeader(header) {
				w.WriteHeader(http.StatusOK)
				obj := &Resp{
					ErrCode: ErrParamsSegSizeLimit.ErrCode,
					ErrMsg:  ErrParamsSegSizeLimit.ErrMsg,
				}
				j, _ := json.Marshal(obj)
				w.Write(j)
				logs.LogError("uuid:%v %v=%v[%v] %v seg_size[%v]", uuid, k, header.Filename, md5, total, header.Size)
				continue
			}
			/// header.size检查
			if header.Size <= 0 {
				w.WriteHeader(http.StatusOK)
				obj := &Resp{
					ErrCode: ErrParamsSegSizeZero.ErrCode,
					ErrMsg:  ErrParamsSegSizeZero.ErrMsg,
				}
				j, _ := json.Marshal(obj)
				w.Write(j)
				logs.LogError("uuid:%v %v=%v[%v] %v seg_size[%v]", uuid, k, header.Filename, md5, total, header.Size)
				continue
			}
			/// total检查
			if !checkTotal(total) {
				w.WriteHeader(http.StatusOK)
				obj := &Resp{
					ErrCode: ErrParamsTotalLimit.ErrCode,
					ErrMsg:  ErrParamsTotalLimit.ErrMsg,
				}
				j, _ := json.Marshal(obj)
				w.Write(j)
				logs.LogError("uuid:%v %v=%v[%v] %v seg_size[%v]", uuid, k, header.Filename, md5, total, header.Size)
				continue
			}
			/// md5检查
			if !checkMD5(md5) {
				w.WriteHeader(http.StatusOK)
				obj := &Resp{
					ErrCode: ErrParamsMD5.ErrCode,
					ErrMsg:  ErrParamsMD5.ErrMsg,
				}
				j, _ := json.Marshal(obj)
				w.Write(j)
				logs.LogError("uuid:%v %v=%v[%v] %v seg_size[%v]", uuid, k, header.Filename, md5, total, header.Size)
				continue
			}
			info := fileInfos.Get(md5)
			if info == nil {
				size, _ := strconv.ParseInt(total, 10, 0)
				allTotal += size
				/// 没有上传，等待上传
				keys = append(keys, k)
				info = &FileInfo{
					Uuid:    uuid,
					Md5:     md5,
					SrcName: header.Filename,
					DstName: utils.RandomString(10) + "." + header.Filename,
					Total:   size,
				}
				fileInfos.Add(md5, info)
				logs.LogWarn("--------------------- 没有上传，等待上传 uuid:%v %v=%v[%v] %v seg_size[%v]", uuid, k, header.Filename, md5, total, header.Size)
			} else {
				/// 已在其它上传任务中
				info.Assert()
				////// 校验MD5
				if md5 != info.Md5 {
					logs.LogFatal("uuid:%v:%v(%v) conflict md5:%v", info.Uuid, info.SrcName, info.Md5, md5)
				}
				////// 校验数据大小
				if total != strconv.FormatInt(info.Total, 10) {
					logs.LogFatal("uuid:%v:%v(%v) conflict %v:%v", info.Uuid, info.SrcName, info.Md5, total, info.Total)
				}
				////// 校验uuid
				if uuid == info.Uuid {
					logs.LogFatal("uuid:%v:%v(%v) conflict uuid:%v", info.Uuid, info.SrcName, info.Md5, uuid)
				}
				////// 还未接收完
				if info.Finished() {
					logs.LogError("uuid:%v:%v(%v) finished uuid:%v", info.Uuid, info.SrcName, info.Md5, uuid)
				}
				ignore = append(ignore, info)
				logs.LogError("重复上传，忽略 uuid:%v %v=%v[%v] seg_size[%v] uuid:%v %v/%v", uuid, k, header.Filename, md5, header.Size, info.Uuid, info.Now, total)
			}
		}
		if !checkAlltotal(allTotal) {
			w.WriteHeader(http.StatusOK)
			obj := &Resp{
				ErrCode: ErrParamsAllTotalLimit.ErrCode,
				ErrMsg:  ErrParamsAllTotalLimit.ErrMsg,
			}
			j, _ := json.Marshal(obj)
			w.Write(j)
			logs.LogError("uuid=%v", uuid)
			return
		}
		if len(keys) > 0 {
			/// 有待上传文件，启动新任务
			uploader := NewUploader(uuid)
			uploaders.Add(uuid, uploader)
			j, _ := json.Marshal(keys)
			logs.LogTrace("--------------------- ****** 有待上传文件，启动任务 uuid:%v ... %v", uuid, string(j))
			uploader.Do(&Req{uuid: uuid, keys: keys, ignore: ignore, w: w, r: r})
		} else {
			/// 无待上传文件，直接返回
			result := []Result{}
			for _, info := range ignore {
				now := strconv.FormatInt(info.Now, 10)
				total := strconv.FormatInt(info.Total, 10)
				result = append(result,
					Result{
						Uuid:    uuid,
						File:    info.SrcName,
						Md5:     info.Md5,
						ErrCode: ErrRepeat.ErrCode,
						ErrMsg:  ErrRepeat.ErrMsg,
						Result:  info.Uuid + " 正在上传 " + info.SrcName + " 进度 " + now + "/" + total})
			}
			if len(result) > 0 {
				w.WriteHeader(http.StatusOK)
				obj := &Resp{
					Data: result,
				}
				j, _ := json.Marshal(obj)
				w.Write(j)
			} else {
				logs.LogTrace("ignore:%#v", ignore)
			}
		}
	} else {
		allTotal := int64(0)
		///////////////////////////// 当前上传任务 /////////////////////////////
		keys := []string{}
		for k := range form.File {
			total := r.FormValue(k + ".total")
			md5 := strings.ToLower(r.FormValue(k + ".md5"))
			/// header检查
			_, header, err := r.FormFile(k)
			if err != nil {
				logs.LogError("%v", err.Error())
				// uploaders.Remove(uuid).Close()
				continue
			}
			/// header.size检查
			if !checkMultiPartFileHeader(header) {
				w.WriteHeader(http.StatusOK)
				obj := &Resp{
					ErrCode: ErrParamsSegSizeLimit.ErrCode,
					ErrMsg:  ErrParamsSegSizeLimit.ErrMsg,
				}
				j, _ := json.Marshal(obj)
				w.Write(j)
				logs.LogError("uuid:%v %v=%v[%v] %v seg_size[%v]", uuid, k, header.Filename, md5, total, header.Size)
				// uploaders.Remove(uuid).Close()
				continue
			}
			/// header.size检查
			if header.Size <= 0 {
				w.WriteHeader(http.StatusOK)
				obj := &Resp{
					ErrCode: ErrParamsSegSizeZero.ErrCode,
					ErrMsg:  ErrParamsSegSizeZero.ErrMsg,
				}
				j, _ := json.Marshal(obj)
				w.Write(j)
				logs.LogError("uuid:%v %v=%v[%v] %v seg_size[%v]", uuid, k, header.Filename, md5, total, header.Size)
				// uploaders.Remove(uuid).Close()
				continue
			}
			/// total检查
			if !checkTotal(total) {
				w.WriteHeader(http.StatusOK)
				obj := &Resp{
					ErrCode: ErrParamsTotalLimit.ErrCode,
					ErrMsg:  ErrParamsTotalLimit.ErrMsg,
				}
				j, _ := json.Marshal(obj)
				w.Write(j)
				logs.LogError("uuid:%v %v=%v[%v] %v seg_size[%v]", uuid, k, header.Filename, md5, total, header.Size)
				// uploaders.Remove(uuid).Close()
				continue
			}
			/// md5检查
			if !checkMD5(md5) {
				w.WriteHeader(http.StatusOK)
				obj := &Resp{
					ErrCode: ErrParamsMD5.ErrCode,
					ErrMsg:  ErrParamsMD5.ErrMsg,
				}
				j, _ := json.Marshal(obj)
				w.Write(j)
				logs.LogError("uuid:%v %v=%v[%v] %v seg_size[%v]", uuid, k, header.Filename, md5, total, header.Size)
				// uploaders.Remove(uuid).Close()
				continue
			}
			info := fileInfos.Get(md5)
			if info == nil {
				size, _ := strconv.ParseInt(total, 10, 0)
				allTotal += size
				/// 没有上传，等待上传
				keys = append(keys, k)
				info = &FileInfo{
					Uuid:    uuid,
					Md5:     md5,
					SrcName: header.Filename,
					DstName: utils.RandomString(10) + "." + header.Filename,
					Total:   size,
				}
				fileInfos.Add(md5, info)
				logs.LogWarn("--------------------- 没有上传，等待上传 uuid:%v %v=%v[%v] %v seg_size[%v]", uuid, k, header.Filename, md5, total, header.Size)
			} else {
				size, _ := strconv.ParseInt(total, 10, 0)
				allTotal += size
				/// 已在当前上传任务中
				info.Assert()
				////// 校验MD5
				if md5 != info.Md5 {
					logs.LogFatal("uuid:%v:%v(%v) conflict md5:%v", info.Uuid, info.SrcName, info.Md5, md5)
				}
				////// 校验数据大小
				if total != strconv.FormatInt(info.Total, 10) {
					logs.LogFatal("uuid:%v:%v(%v) conflict %v:%v", info.Uuid, info.SrcName, info.Md5, total, info.Total)
				}
				////// 校验uuid
				if uuid != info.Uuid {
					logs.LogFatal("uuid:%v:%v(%v) conflict uuid:%v", info.Uuid, info.SrcName, info.Md5, uuid)
				}
				////// 校验filename
				if header.Filename != info.SrcName {
					logs.LogFatal("uuid:%v:%v(%v) conflict %v", info.Uuid, info.SrcName, info.Md5, header.Filename)
				}
				////// 还未接收完
				if info.Finished() {
					logs.LogFatal("uuid:%v:%v(%v) finished", info.Uuid, info.SrcName, info.Md5)
				}
				keys = append(keys, k)
				logs.LogInfo("继续上传中 uuid:%v %v=%v[%v] %v/%v seg_size[%d]", uuid, k, header.Filename, md5, info.Now, total, header.Size)
			}
		}
		if !checkAlltotal(allTotal) {
			w.WriteHeader(http.StatusOK)
			obj := &Resp{
				ErrCode: ErrParamsAllTotalLimit.ErrCode,
				ErrMsg:  ErrParamsAllTotalLimit.ErrMsg,
			}
			j, _ := json.Marshal(obj)
			w.Write(j)
			logs.LogError("uuid=%v", uuid)
			//需要移除任务
			uploaders.Remove(uuid).Close()
			return
		}
		if len(keys) > 0 {
			/// 有待上传文件，加入当前任务
			uploader.Do(&Req{uuid: uuid, keys: keys, w: w, r: r})
		} else {
			/// 无待上传文件
			logs.LogTrace("--------------------- ****** 无待上传文件，当前任务 uuid:%v ...", uuid)
		}
	}
}
