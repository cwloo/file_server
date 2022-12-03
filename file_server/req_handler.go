package main

import (
	"encoding/json"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/cwloo/gonet/logs"
)

func handlerUploadFile(w http.ResponseWriter, r *http.Request) {
	// w.WriteHeader(http.StatusOK)
	err := r.ParseMultipartForm(MaxMemory)
	if err != nil {
		logs.LogError(err.Error())
		resp := &Resp{
			ErrCode: ErrParsePartData.ErrCode,
			ErrMsg:  ErrParsePartData.ErrMsg,
		}
		j, _ := json.Marshal(resp)
		w.Header().Set("Content-Length", strconv.Itoa(len(j)))
		_, err := w.Write(j)
		if err != nil {
			logs.LogError(err.Error())
		}
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
		resp := &Resp{
			ErrCode: ErrParamsUUID.ErrCode,
			ErrMsg:  ErrParamsUUID.ErrMsg,
		}
		j, _ := json.Marshal(resp)
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
	allTotal := int64(0)
	keys := []string{}
	for k := range form.File {
		offset := r.FormValue(k + ".offset")
		total := r.FormValue(k + ".total")
		md5 := strings.ToLower(k)
		/// header检查
		_, header, err := r.FormFile(k)
		if err != nil {
			logs.LogError("%v", err.Error())
			result = append(result,
				Result{
					Uuid:    uuid,
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
					File:    header.Filename,
					Md5:     md5,
					ErrCode: ErrParamsTotalLimit.ErrCode,
					ErrMsg:  ErrParamsTotalLimit.ErrMsg,
					Result:  ""})
			logs.LogError("uuid:%v %v=%v[%v] %v seg_size[%v]", uuid, k, header.Filename, md5, total, header.Size)
			continue
		}
		/// offset检查
		if !checkOffset(offset, total) {
			result = append(result,
				Result{
					Uuid:    uuid,
					File:    header.Filename,
					Md5:     md5,
					ErrCode: ErrParamsOffset.ErrCode,
					ErrMsg:  ErrParamsOffset.ErrMsg,
					Result:  ""})
			logs.LogError("uuid:%v %v=%v[%v] %v seg_size[%v]", uuid, k, header.Filename, md5, total, header.Size)
			continue
		}
		/// md5检查
		if !checkMD5(md5) {
			result = append(result,
				Result{
					Uuid:    uuid,
					File:    header.Filename,
					Md5:     md5,
					ErrCode: ErrParamsMD5.ErrCode,
					ErrMsg:  ErrParamsMD5.ErrMsg,
					Result:  ""})
			logs.LogError("uuid:%v %v=%v[%v] %v seg_size[%v]", uuid, k, header.Filename, md5, total, header.Size)
			continue
		}
		fi := fileInfos.Get(md5)
		if fi == nil {
			/// 没有上传，判断能否上传
			size, _ := strconv.ParseInt(total, 10, 0)
			allTotal += size
			if !checkAlltotal(allTotal) {
				result = append(result,
					Result{
						Uuid:    uuid,
						File:    header.Filename,
						Md5:     md5,
						ErrCode: ErrParamsAllTotalLimit.ErrCode,
						ErrMsg:  ErrParamsAllTotalLimit.ErrMsg,
						Result:  ""})
				logs.LogError("uuid:%v %v=%v[%v] %v seg_size[%v]", uuid, k, header.Filename, md5, total, header.Size)
				continue
			}
		} else {
			if fi.Uuid == uuid {
				/// 已在当前上传任务中
				size, _ := strconv.ParseInt(total, 10, 0)
				allTotal += size
			}
		}
		info, ok := fileInfos.GetAdd(md5, uuid, header.Filename, total)
		if !ok {
			/// 没有上传，等待上传
			keys = append(keys, k)
			logs.LogWarn("--- *** 没有上传，等待上传 uuid:%v %v=%v[%v] %v seg_size[%v]", uuid, k, header.Filename, md5, total, header.Size)
		} else {
			info.Assert()
			if info.Uuid == uuid {
				/// 已在当前上传任务中

				////// 校验MD5
				if md5 != info.Md5 {
					logs.LogFatal("uuid:%v:%v(%v) md5:%v", info.Uuid, info.SrcName, info.Md5, md5)
				}
				////// 校验数据大小
				if total != strconv.FormatInt(info.Total, 10) {
					logs.LogFatal("uuid:%v:%v(%v) info.total:%v total:%v", info.Uuid, info.SrcName, info.Md5, info.Total, total)
				}
				if info.Ok() {
					if info.Md5Ok {
						info.UpdateHitTime(time.Now())
						// fileInfos.Remove(info.Md5)
						result = append(result,
							Result{
								Uuid:    uuid,
								File:    header.Filename,
								Md5:     info.Md5,
								Now:     info.Now,
								Total:   info.Total,
								ErrCode: ErrOk.ErrCode,
								ErrMsg:  ErrOk.ErrMsg,
								Result:  strings.Join([]string{"uuid:", info.Uuid, " uploading ", info.DstName, " progress:", strconv.FormatInt(info.Now, 10) + "/" + total + " 上传成功!"}, "")})
					} else {
						fileInfos.Remove(info.Md5)
						os.Remove(dir_upload + info.DstName)
						result = append(result,
							Result{
								Uuid:    uuid,
								File:    header.Filename,
								Md5:     info.Md5,
								Now:     info.Now,
								Total:   info.Total,
								ErrCode: ErrFileMd5.ErrCode,
								ErrMsg:  ErrFileMd5.ErrMsg,
								Result:  strings.Join([]string{"uuid:", info.Uuid, " uploading ", info.DstName, " progress:", strconv.FormatInt(info.Now, 10) + "/" + total + " 上传完毕 MD5校验失败!"}, "")})
					}
				} else {
					keys = append(keys, k)
				}
			} else {
				/// 已在其它上传任务中

				if info.Ok() {
					if info.Md5Ok {
						info.UpdateHitTime(time.Now())
						// fileInfos.Remove(info.Md5)
						result = append(result,
							Result{
								Uuid:    uuid,
								File:    header.Filename,
								Md5:     info.Md5,
								Now:     info.Now,
								Total:   info.Total,
								ErrCode: ErrOk.ErrCode,
								ErrMsg:  ErrOk.ErrMsg,
								Result:  strings.Join([]string{"uuid:", info.Uuid, " uploading ", info.DstName, " progress:", strconv.FormatInt(info.Now, 10) + "/" + total + " 别人上传成功!"}, "")})
					} else {
						fileInfos.Remove(info.Md5)
						os.Remove(dir_upload + info.DstName)
						result = append(result,
							Result{
								Uuid:    uuid,
								File:    header.Filename,
								Md5:     info.Md5,
								Now:     info.Now,
								Total:   info.Total,
								ErrCode: ErrFileMd5.ErrCode,
								ErrMsg:  ErrFileMd5.ErrMsg,
								Result:  strings.Join([]string{"uuid:", info.Uuid, " uploading ", info.DstName, " progress:", strconv.FormatInt(info.Now, 10) + "/" + total + " 别人上传完毕 MD5校验失败!"}, "")})
					}
				} else {
					result = append(result,
						Result{
							Uuid:    uuid,
							File:    info.SrcName,
							Md5:     info.Md5,
							ErrCode: ErrRepeat.ErrCode,
							ErrMsg:  ErrRepeat.ErrMsg,
							Result:  strings.Join([]string{"uuid:", info.Uuid, " uploading ", info.DstName, " progress:", strconv.FormatInt(info.Now, 10), "/", total}, " 别人上传中"),
						})
					// logs.LogError("--- *** ignore repeat-upload uuid:%v %v=%v[%v] seg_size[%v] uuid:%v uploading %v progress:%v/%v", uuid, k, header.Filename, md5, header.Size, info.Uuid, info.DstName, info.Now, total)
				}
			}
		}
	} /// {{{ end for range form.File
	var exist bool
	if len(keys) > 0 {
		uploader, ok := uploaders.GetAdd(uuid, UseAsyncUploader)
		if !ok {
			///////////////////////////// 新的上传任务 /////////////////////////////
			/// 有待上传文件，启动新任务
			j, _ := json.Marshal(keys)
			logs.LogTrace("--------------------- ****** 有待上传文件，启动任务 uuid:%v ... %v", uuid, string(j))
			uploader.Upload(&Req{uuid: uuid, keys: keys, w: w, r: r, resp: resp, result: result})
		} else {
			exist = true
			///////////////////////////// 当前上传任务 /////////////////////////////
			/// 有待上传文件，加入当前任务
			// j, _ := json.Marshal(keys)
			// logs.LogTrace("--------------------- ****** 有待上传文件，加入任务 uuid:%v ... %v", uuid, string(j))
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
			if exist {
				logs.LogTrace("--------------------- ****** 无待上传文件，当前任务 uuid:%v ...", uuid)
			} else {
				logs.LogTrace("--------------------- ****** 无待上传文件，未分配任务 uuid:%v ...", uuid)
			}
		}
	}
}
