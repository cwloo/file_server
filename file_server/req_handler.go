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
	// w.WriteHeader(http.StatusOK)
	err := r.ParseMultipartForm(MaxMemory)
	if err != nil {
		obj := &Resp{
			ErrCode: ErrParsePartData.ErrCode,
			ErrMsg:  ErrParsePartData.ErrMsg,
		}
		j, _ := json.Marshal(obj)
		w.Header().Set("Content-Length", strconv.Itoa(len(j)))
		_, err := w.Write(j)
		if err != nil {
			logs.LogError(err.Error())
		}
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
		obj := &Resp{
			ErrCode: ErrParamsUUID.ErrCode,
			ErrMsg:  ErrParamsUUID.ErrMsg,
		}
		j, _ := json.Marshal(obj)
		w.Header().Set("Content-Length", strconv.Itoa(len(j)))
		_, err := w.Write(j)
		if err != nil {
			logs.LogError(err.Error())
		}
		logs.LogError("uuid=%v", uuid)
		return
	}
	var resp *Resp
	result := []Result{}
	uploader := uploaders.Get(uuid)
	if uploader == nil {
		allTotal := int64(0)
		///////////////////////////// 新的上传任务 /////////////////////////////
		keys := []string{}
		for k := range form.File {
			total := r.FormValue(k + ".total")
			md5 := strings.ToLower(r.FormValue(k + ".md5"))
			/// header检查
			_, header, err := r.FormFile(k)
			if err != nil {
				logs.LogError("%v", err.Error())
				result = append(result,
					Result{
						Uuid:    uuid,
						Key:     k,
						File:    "",
						Md5:     md5,
						ErrCode: ErrParseFormFile.ErrCode,
						ErrMsg:  ErrParseFormFile.ErrMsg,
						Result:  ""})
				continue
			}
			/// header.size检查
			if !checkMultiPartFileHeader(header) {
				result = append(result,
					Result{
						Uuid:    uuid,
						Key:     k,
						File:    header.Filename,
						Md5:     md5,
						ErrCode: ErrParamsSegSizeLimit.ErrCode,
						ErrMsg:  ErrParamsSegSizeLimit.ErrMsg,
						Result:  ""})
				logs.LogError("uuid:%v %v=%v[%v] %v seg_size[%v]", uuid, k, header.Filename, md5, total, header.Size)
				continue
			}
			/// header.size检查
			if header.Size <= 0 {
				result = append(result,
					Result{
						Uuid:    uuid,
						Key:     k,
						File:    header.Filename,
						Md5:     md5,
						ErrCode: ErrParamsSegSizeZero.ErrCode,
						ErrMsg:  ErrParamsSegSizeZero.ErrMsg,
						Result:  ""})
				logs.LogError("uuid:%v %v=%v[%v] %v seg_size[%v]", uuid, k, header.Filename, md5, total, header.Size)
				continue
			}
			/// total检查
			if !checkTotal(total) {
				result = append(result,
					Result{
						Uuid:    uuid,
						Key:     k,
						File:    header.Filename,
						Md5:     md5,
						ErrCode: ErrParamsTotalLimit.ErrCode,
						ErrMsg:  ErrParamsTotalLimit.ErrMsg,
						Result:  ""})
				logs.LogError("uuid:%v %v=%v[%v] %v seg_size[%v]", uuid, k, header.Filename, md5, total, header.Size)
				continue
			}
			/// md5检查
			if !checkMD5(md5) {
				result = append(result,
					Result{
						Uuid:    uuid,
						Key:     k,
						File:    header.Filename,
						Md5:     md5,
						ErrCode: ErrParamsMD5.ErrCode,
						ErrMsg:  ErrParamsMD5.ErrMsg,
						Result:  ""})
				logs.LogError("uuid:%v %v=%v[%v] %v seg_size[%v]", uuid, k, header.Filename, md5, total, header.Size)
				continue
			}
			info := fileInfos.Get(md5)
			if info == nil {
				size, _ := strconv.ParseInt(total, 10, 0)
				allTotal += size
				/// 没有上传，等待上传
				info = &FileInfo{
					Uuid:    uuid,
					Md5:     md5,
					SrcName: header.Filename,
					DstName: strings.Join([]string{uuid, utils.RandomString(10), ".", header.Filename}, ""),
					Total:   size,
				}
				fileInfos.Add(md5, info)
				keys = append(keys, k)
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
				result = append(result,
					Result{
						Uuid:    uuid,
						File:    info.SrcName,
						Md5:     info.Md5,
						ErrCode: ErrRepeat.ErrCode,
						ErrMsg:  ErrRepeat.ErrMsg,
						Result:  strings.Join([]string{"uuid:", info.Uuid, " uploading ", info.DstName, " progress:", strconv.FormatInt(info.Now, 10), "/", total}, ""),
					})
				logs.LogError("--------------------- ****** 忽略重复上传 uuid:%v %v=%v[%v] seg_size[%v] uuid:%v uploading %v progress:%v/%v", uuid, k, header.Filename, md5, header.Size, info.Uuid, info.DstName, info.Now, total)
			}
		} /// {{{ end for range form.File
		if !checkAlltotal(allTotal) {
			resp = &Resp{
				ErrCode: ErrParamsAllTotalLimit.ErrCode,
				ErrMsg:  ErrParamsAllTotalLimit.ErrMsg,
			}
		}
		if len(keys) > 0 {
			/// 有待上传文件，启动新任务
			j, _ := json.Marshal(keys)
			logs.LogTrace("--------------------- ****** 有待上传文件，启动任务 uuid:%v ... %v", uuid, string(j))
			switch UseAsyncUploader {
			case true:
				uploader := NewAsyncUploader(uuid)
				uploaders.Add(uuid, uploader)
				uploader.Upload(&Req{uuid: uuid, keys: keys, w: w, r: r, resp: resp, result: result})
			default:
				uploader := NewSyncUploader(uuid)
				uploaders.Add(uuid, uploader)
				uploader.Upload(&Req{uuid: uuid, keys: keys, w: w, r: r, resp: resp, result: result})
			}
		} else {
			/// 无待上传文件，直接返回
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
				w.Header().Set("Content-Length", strconv.Itoa(len(j)))
				_, err := w.Write(j)
				if err != nil {
					logs.LogError(err.Error())
				}
				// logs.LogError("uuid:%v %v", uuid, string(j))
			} else {
				resp = &Resp{}
				j, _ := json.Marshal(resp)
				w.Header().Set("Content-Length", strconv.Itoa(len(j)))
				_, err := w.Write(j)
				if err != nil {
					logs.LogError(err.Error())
				}
				logs.LogFatal("--------------------- ****** 无待上传文件，未分配任务 uuid:%v ...", uuid)
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
				result = append(result,
					Result{
						Uuid:    uuid,
						Key:     k,
						File:    "",
						Md5:     md5,
						ErrCode: ErrParseFormFile.ErrCode,
						ErrMsg:  ErrParseFormFile.ErrMsg,
						Result:  ""})
				continue
			}
			/// header.size检查
			if !checkMultiPartFileHeader(header) {
				result = append(result,
					Result{
						Uuid:    uuid,
						Key:     k,
						File:    header.Filename,
						Md5:     md5,
						ErrCode: ErrParamsSegSizeLimit.ErrCode,
						ErrMsg:  ErrParamsSegSizeLimit.ErrMsg,
						Result:  ""})
				logs.LogError("uuid:%v %v=%v[%v] %v seg_size[%v]", uuid, k, header.Filename, md5, total, header.Size)
				// uploaders.Remove(uuid).Close()
				continue
			}
			/// header.size检查
			if header.Size <= 0 {
				result = append(result,
					Result{
						Uuid:    uuid,
						Key:     k,
						File:    header.Filename,
						Md5:     md5,
						ErrCode: ErrParamsSegSizeZero.ErrCode,
						ErrMsg:  ErrParamsSegSizeZero.ErrMsg,
						Result:  ""})
				logs.LogError("uuid:%v %v=%v[%v] %v seg_size[%v]", uuid, k, header.Filename, md5, total, header.Size)
				// uploaders.Remove(uuid).Close()
				continue
			}
			/// total检查
			if !checkTotal(total) {
				result = append(result,
					Result{
						Uuid:    uuid,
						Key:     k,
						File:    header.Filename,
						Md5:     md5,
						ErrCode: ErrParamsTotalLimit.ErrCode,
						ErrMsg:  ErrParamsTotalLimit.ErrMsg,
						Result:  ""})
				logs.LogError("uuid:%v %v=%v[%v] %v seg_size[%v]", uuid, k, header.Filename, md5, total, header.Size)
				// uploaders.Remove(uuid).Close()
				continue
			}
			/// md5检查
			if !checkMD5(md5) {
				result = append(result,
					Result{
						Uuid:    uuid,
						Key:     k,
						File:    header.Filename,
						Md5:     md5,
						ErrCode: ErrParamsMD5.ErrCode,
						ErrMsg:  ErrParamsMD5.ErrMsg,
						Result:  ""})
				logs.LogError("uuid:%v %v=%v[%v] %v seg_size[%v]", uuid, k, header.Filename, md5, total, header.Size)
				// uploaders.Remove(uuid).Close()
				continue
			}
			info := fileInfos.Get(md5)
			if info == nil {
				size, _ := strconv.ParseInt(total, 10, 0)
				allTotal += size
				/// 没有上传，等待上传
				info = &FileInfo{
					Uuid:    uuid,
					Md5:     md5,
					SrcName: header.Filename,
					DstName: strings.Join([]string{uuid, utils.RandomString(10), ".", header.Filename}, ""),
					Total:   size,
				}
				fileInfos.Add(md5, info)
				keys = append(keys, k)
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
		} /// {{{ end for range form.File
		if !checkAlltotal(allTotal) {
			resp = &Resp{
				ErrCode: ErrParamsAllTotalLimit.ErrCode,
				ErrMsg:  ErrParamsAllTotalLimit.ErrMsg,
			}
			// uploaders.Remove(uuid).Close()
		}
		if len(keys) > 0 {
			/// 有待上传文件，加入当前任务
			j, _ := json.Marshal(keys)
			logs.LogTrace("--------------------- ****** 有待上传文件，加入任务 uuid:%v ... %v", uuid, string(j))
			uploader.Upload(&Req{uuid: uuid, keys: keys, w: w, r: r, resp: resp, result: result})
		} else {
			/// 无待上传文件，直接返回
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
				w.Header().Set("Content-Length", strconv.Itoa(len(j)))
				_, err := w.Write(j)
				if err != nil {
					logs.LogError(err.Error())
				}
				// logs.LogError("uuid:%v %v", uuid, string(j))
			} else {
				resp = &Resp{}
				j, _ := json.Marshal(resp)
				w.Header().Set("Content-Length", strconv.Itoa(len(j)))
				_, err := w.Write(j)
				if err != nil {
					logs.LogError(err.Error())
				}
				logs.LogTrace("--------------------- ****** 无待上传文件，当前任务 uuid:%v ...", uuid)
			}
		}
	}
}
