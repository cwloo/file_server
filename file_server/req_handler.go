package main

import (
	"encoding/json"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/cwloo/gonet/logs"
	"github.com/cwloo/uploader/file_server/config"
	"github.com/cwloo/uploader/file_server/global"
)

func setResponseHeader(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type,AccessToken,X-CSRF-Token,Authorization,Token,X-Token,X-User-Id,C-Token,cz-sdk-key,cz-sdk-sign")
	w.Header().Set("Access-Control-Allow-Methods", "POST,GET,HEAD,TRACE,OPTIONS,DELETE,PUT")
	w.Header().Set("Access-Control-Expose-Headers", "Content-Length,Access-Control-Allow-Origin,Access-Control-Allow-Headers,Access-Control-Request-Headers,Access-Control-Request-Method,Content-Type,New-Token,New-Expires-At,New-C-Token,New-C-Expires-At")
	w.Header().Set("Access-Control-Allow-Credentials", "true")

	w.Header().Set("Host", r.Header.Get("Host"))
	w.Header().Set("X-Real-IP", r.Header.Get("X-Real-IP"))
	w.Header().Set("X-Forwarded-For", r.Header.Get("X-Forwarded-For"))
	w.Header().Set("X-Forwarded-Proto", r.Header.Get("X-Forwarded-Proto"))
	w.Header().Set("Remote-Host", r.Header.Get("Remote-Host"))
	w.Header().Set("User-Agent", r.Header.Get("User-Agent"))
	w.Header().Set("Referer", r.Header.Get("Referer"))
	w.Header().Set("Access-Control-Request-Headers", r.Header.Get("Access-Control-Request-Headers"))
	w.Header().Set("Access-Control-Request-Method", r.Header.Get("Access-Control-Request-Method"))
	w.Header().Set("Origin", r.Header.Get("Origin"))
	w.Header().Set("Sec-Fetch-Dest", r.Header.Get("Sec-Fetch-Dest"))
	w.Header().Set("Sec-Fetch-Mode", r.Header.Get("Sec-Fetch-Mode"))
	w.Header().Set("Sec-Fetch-Site", r.Header.Get("Sec-Fetch-Site"))
	// w.Header().Set("Accept-Encoding", r.Header.Get("Accept-Encoding"))
	// w.Header().Set("Accept-Language", r.Header.Get("Accept-Language"))
}

func writeResponse(w http.ResponseWriter, r *http.Request, v any) {
	j, _ := json.Marshal(v)
	setResponseHeader(w, r)
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Content-Length", strconv.Itoa(len(j)))
	_, err := w.Write(j)
	if err != nil {
		logs.LogError(err.Error())
		return
	}
	logs.LogDebug("%v", string(j))
}

func handlerUpload(w http.ResponseWriter, r *http.Request) {
	// w.WriteHeader(http.StatusOK)
	err := r.ParseMultipartForm(config.Config.MaxMemory)
	if err != nil {
		logs.LogError(err.Error())
		resp := &global.Resp{
			ErrCode: global.ErrParsePartData.ErrCode,
			ErrMsg:  global.ErrParsePartData.ErrMsg,
		}
		writeResponse(w, r, resp)
		return
	}
	uuid := ""
	md5 := ""
	offset := ""
	total := ""
	form := r.MultipartForm
	for k := range form.Value {
		switch strings.ToLower(k) {
		case "uuid":
			uuid = r.FormValue(k)
		case "md5":
			md5 = strings.ToLower(r.FormValue(k))
		case "offset":
			offset = r.FormValue(k)
		case "total":
			total = r.FormValue(k)
		}
		// logs.LogTrace("%v=%v", k, v)
	}
	if !checkUUID(uuid) {
		resp := &global.Resp{
			ErrCode: global.ErrParamsUUID.ErrCode,
			ErrMsg:  global.ErrParamsUUID.ErrMsg,
		}
		writeResponse(w, r, resp)
		logs.LogError("uuid=%v", uuid)
		return
	}
	var resp *global.Resp
	result := []global.Result{}
	allTotal := int64(0)
	keys := []string{}
	if len(form.File) > 1 {
		resp := &global.Resp{
			ErrCode: global.ErrMultiFileNotSupport.ErrCode,
			ErrMsg:  global.ErrMultiFileNotSupport.ErrMsg,
		}
		writeResponse(w, r, resp)
		return
	}
	for k := range form.File {
		/// header检查
		_, header, err := r.FormFile(k)
		if err != nil {
			logs.LogError(err.Error())
			result = append(result,
				global.Result{
					Uuid:    uuid,
					File:    "",
					Md5:     md5,
					ErrCode: global.ErrParseFormFile.ErrCode,
					ErrMsg:  global.ErrParseFormFile.ErrMsg,
					Message: ""})
			continue
		}
		/// header.size检查
		if !checkMultiPartSize(header) {
			result = append(result,
				global.Result{
					Uuid:    uuid,
					File:    header.Filename,
					Md5:     md5,
					ErrCode: global.ErrParamsSegSizeZero.ErrCode,
					ErrMsg:  global.ErrParamsSegSizeZero.ErrMsg,
					Message: ""})
			logs.LogError("%v %v[%v] %v seg_size[%v]", uuid, header.Filename, md5, total, header.Size)
			continue
		}
		/// header.size检查
		if !checkMultiPartSizeLimit(header) {
			result = append(result,
				global.Result{
					Uuid:    uuid,
					File:    header.Filename,
					Md5:     md5,
					ErrCode: global.ErrParamsSegSizeLimit.ErrCode,
					ErrMsg:  global.ErrParamsSegSizeLimit.ErrMsg,
					Message: ""})
			logs.LogError("%v %v[%v] %v seg_size[%v]", uuid, header.Filename, md5, total, header.Size)
			continue
		}
		/// total检查
		if !checkSingle(total) {
			result = append(result,
				global.Result{
					Uuid:    uuid,
					File:    header.Filename,
					Md5:     md5,
					ErrCode: global.ErrParamsTotalLimit.ErrCode,
					ErrMsg:  global.ErrParamsTotalLimit.ErrMsg,
					Message: ""})
			logs.LogError("%v %v[%v] %v seg_size[%v]", uuid, header.Filename, md5, total, header.Size)
			continue
		}
		/// offset检查
		if !checkOffset(offset, total) {
			result = append(result,
				global.Result{
					Uuid:    uuid,
					File:    header.Filename,
					Md5:     md5,
					ErrCode: global.ErrParamsOffset.ErrCode,
					ErrMsg:  global.ErrParamsOffset.ErrMsg,
					Message: ""})
			logs.LogError("%v %v[%v] %v seg_size[%v]", uuid, header.Filename, md5, total, header.Size)
			continue
		}
		/// md5检查
		if !checkMD5(md5) {
			result = append(result,
				global.Result{
					Uuid:    uuid,
					File:    header.Filename,
					Md5:     md5,
					ErrCode: global.ErrParamsMD5.ErrCode,
					ErrMsg:  global.ErrParamsMD5.ErrMsg,
					Message: ""})
			logs.LogError("%v %v[%v] %v seg_size[%v]", uuid, header.Filename, md5, total, header.Size)
			continue
		}
		fi := fileInfos.Get(md5)
		if fi == nil {
			/// 没有上传，判断能否上传
			size, _ := strconv.ParseInt(total, 10, 0)
			allTotal += size
			if !checkTotal(allTotal) {
				result = append(result,
					global.Result{
						Uuid:    uuid,
						File:    header.Filename,
						Md5:     md5,
						ErrCode: global.ErrParamsAllTotalLimit.ErrCode,
						ErrMsg:  global.ErrParamsAllTotalLimit.ErrMsg,
						Message: ""})
				logs.LogError("%v %v[%v] %v seg_size[%v]", uuid, header.Filename, md5, total, header.Size)
				continue
			}
		} else {
			if fi.Uuid() == uuid {
				/// 已在当前上传任务中
				size, _ := strconv.ParseInt(total, 10, 0)
				allTotal += size
			}
		}
		info, ok := fileInfos.GetAdd(md5, uuid, header.Filename, total)
		if !ok {
			/// 没有上传，等待上传
			keys = append(keys, k)
			logs.LogWarn("--- *** 没有上传，等待上传 %v %v[%v] %v/%v seg_size[%v]", uuid, header.Filename, md5, info.Now(false), total, header.Size)
		} else {
			info.Assert()
			if info.Uuid() == uuid {
				/// 已在当前上传任务中

				////// 校验MD5
				if md5 != info.Md5() {
					logs.LogFatal("%v %v(%v) md5:%v", info.Uuid, info.SrcName, info.Md5, md5)
				}
				////// 校验数据大小
				if total != strconv.FormatInt(info.Total(false), 10) {
					logs.LogFatal("%v %v(%v) info.total:%v total:%v", info.Uuid(), info.SrcName(), info.Md5(), info.Total(false), total)
				}
				if info.Done(true) {
					if ok, url := info.Ok(true); ok {
						info.UpdateHitTime(time.Now())
						// fileInfos.Remove(info.Md5()).Put()
						result = append(result,
							global.Result{
								Uuid:    uuid,
								File:    header.Filename,
								Md5:     info.Md5(),
								Now:     info.Now(true),
								Total:   info.Total(false),
								ErrCode: global.ErrOk.ErrCode,
								ErrMsg:  global.ErrOk.ErrMsg,
								Url:     url,
								Message: strings.Join([]string{info.Uuid(), " uploading ", info.DstName(), " progress:", strconv.FormatInt(info.Now(true), 10) + "/" + total + " 上传成功!"}, "")})
						logs.LogWarn("%v %v[%v] %v chkmd5 [ok] %v", uuid, header.Filename, info.Md5(), info.DstName(), url)
					} else {
						fileInfos.Remove(info.Md5()).Put()
						os.Remove(config.Config.UploadlDir + info.DstName())
						result = append(result,
							global.Result{
								Uuid:    uuid,
								File:    header.Filename,
								Md5:     info.Md5(),
								Now:     info.Now(true),
								Total:   info.Total(false),
								ErrCode: global.ErrFileMd5.ErrCode,
								ErrMsg:  global.ErrFileMd5.ErrMsg,
								Message: strings.Join([]string{info.Uuid(), " uploading ", info.DstName(), " progress:", strconv.FormatInt(info.Now(true), 10) + "/" + total + " 上传完毕 MD5校验失败!"}, "")})
						logs.LogError("%v %v[%v] %v chkmd5 [Err]", uuid, header.Filename, md5, info.DstName())
					}
				} else {
					keys = append(keys, k)
				}
			} else {
				/// 已在其它上传任务中

				if info.Done(true) {
					if ok, url := info.Ok(true); ok {
						info.UpdateHitTime(time.Now())
						// fileInfos.Remove(info.Md5).Put()
						result = append(result,
							global.Result{
								Uuid:    uuid,
								File:    header.Filename,
								Md5:     info.Md5(),
								Now:     info.Now(true),
								Total:   info.Total(false),
								ErrCode: global.ErrOk.ErrCode,
								ErrMsg:  global.ErrOk.ErrMsg,
								Url:     url,
								Message: strings.Join([]string{info.Uuid(), " uploading ", info.DstName(), " progress:", strconv.FormatInt(info.Now(true), 10) + "/" + total + " 别人上传成功!"}, "")})
						logs.LogWarn("%v %v[%v] %v chkmd5 [ok] %v", uuid, header.Filename, info.Md5(), info.DstName(), url)
					} else {
						fileInfos.Remove(info.Md5()).Put()
						os.Remove(config.Config.UploadlDir + info.DstName())
						result = append(result,
							global.Result{
								Uuid:    uuid,
								File:    header.Filename,
								Md5:     info.Md5(),
								Now:     info.Now(true),
								Total:   info.Total(false),
								ErrCode: global.ErrFileMd5.ErrCode,
								ErrMsg:  global.ErrFileMd5.ErrMsg,
								Message: strings.Join([]string{info.Uuid(), " uploading ", info.DstName(), " progress:", strconv.FormatInt(info.Now(true), 10) + "/" + total + " 别人上传完毕 MD5校验失败!"}, "")})
						logs.LogError("%v %v[%v] %v chkmd5 [Err]", uuid, header.Filename, md5, info.DstName())
					}
				} else {
					result = append(result,
						global.Result{
							Uuid:    uuid,
							File:    info.SrcName(),
							Md5:     info.Md5(),
							ErrCode: global.ErrRepeat.ErrCode,
							ErrMsg:  global.ErrRepeat.ErrMsg,
							Message: strings.Join([]string{info.Uuid(), " uploading ", info.DstName(), " progress:", strconv.FormatInt(info.Now(true), 10), "/", total}, " 别人上传中"),
						})
					// logs.LogError("--- *** ignore repeat-upload %v %v[%v] seg_size[%v] %v uploading %v progress:%v/%v", uuid, header.Filename, md5, header.Size, info.Uuid(), info.DstName(), info.Now(), total)
				}
			}
		}
	} /// {{{ end for range form.File
	var exist bool
	if len(keys) > 0 {
		uploader, ok := uploaders.GetAdd(uuid, config.Config.UseAsync > 0)
		if !ok {
			///////////////////////////// 新的上传任务 /////////////////////////////
			/// 有待上传文件，启动新任务
			j, _ := json.Marshal(keys)
			logs.LogTrace("--------------------- ****** 有待上传文件，启动任务 %v ... %v", uuid, string(j))
			uploader.Upload(&global.Req{Uuid: uuid, Md5: md5, Offset: offset, Total: total, Keys: keys, W: w, R: r, Resp: resp, Result: result})
		} else {
			exist = true
			///////////////////////////// 当前上传任务 /////////////////////////////
			/// 有待上传文件，加入当前任务
			// j, _ := json.Marshal(keys)
			// logs.LogTrace("--------------------- ****** 有待上传文件，加入任务 %v ... %v", uuid, string(j))
			uploader.Upload(&global.Req{Uuid: uuid, Md5: md5, Offset: offset, Total: total, Keys: keys, W: w, R: r, Resp: resp, Result: result})
		}
	} else {
		/// 无待上传文件，直接返回
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
			writeResponse(w, r, resp)
			// logs.LogError("%v %v", uuid, string(j))
		} else {
			writeResponse(w, r, &global.Resp{})
			if exist {
				logs.LogTrace("--------------------- ****** 无待上传文件，当前任务 %v ...", uuid)
			} else {
				logs.LogTrace("--------------------- ****** 无待上传文件，未分配任务 %v ...", uuid)
			}
		}
	}
}

func handlerMultiUpload(w http.ResponseWriter, r *http.Request) {
	// w.WriteHeader(http.StatusOK)
	err := r.ParseMultipartForm(config.Config.MaxMemory)
	if err != nil {
		logs.LogError(err.Error())
		resp := &global.Resp{
			ErrCode: global.ErrParsePartData.ErrCode,
			ErrMsg:  global.ErrParsePartData.ErrMsg,
		}
		writeResponse(w, r, resp)
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
		resp := &global.Resp{
			ErrCode: global.ErrParamsUUID.ErrCode,
			ErrMsg:  global.ErrParamsUUID.ErrMsg,
		}
		writeResponse(w, r, resp)
		logs.LogError("uuid=%v", uuid)
		return
	}
	var resp *global.Resp
	result := []global.Result{}
	allTotal := int64(0)
	keys := []string{}
	for k := range form.File {
		offset := r.FormValue(k + ".offset")
		total := r.FormValue(k + ".total")
		md5 := strings.ToLower(k)
		/// header检查
		_, header, err := r.FormFile(k)
		if err != nil {
			logs.LogError(err.Error())
			result = append(result,
				global.Result{
					Uuid:    uuid,
					File:    "",
					Md5:     md5,
					ErrCode: global.ErrParseFormFile.ErrCode,
					ErrMsg:  global.ErrParseFormFile.ErrMsg,
					Message: ""})
			continue
		}
		/// header.size检查
		if !checkMultiPartSize(header) {
			result = append(result,
				global.Result{
					Uuid:    uuid,
					File:    header.Filename,
					Md5:     md5,
					ErrCode: global.ErrParamsSegSizeZero.ErrCode,
					ErrMsg:  global.ErrParamsSegSizeZero.ErrMsg,
					Message: ""})
			logs.LogError("%v %v[%v] %v seg_size[%v]", uuid, header.Filename, md5, total, header.Size)
			continue
		}
		/// header.size检查
		if !checkMultiPartSizeLimit(header) {
			result = append(result,
				global.Result{
					Uuid:    uuid,
					File:    header.Filename,
					Md5:     md5,
					ErrCode: global.ErrParamsSegSizeLimit.ErrCode,
					ErrMsg:  global.ErrParamsSegSizeLimit.ErrMsg,
					Message: ""})
			logs.LogError("%v %v[%v] %v seg_size[%v]", uuid, header.Filename, md5, total, header.Size)
			continue
		}
		/// total检查
		if !checkSingle(total) {
			result = append(result,
				global.Result{
					Uuid:    uuid,
					File:    header.Filename,
					Md5:     md5,
					ErrCode: global.ErrParamsTotalLimit.ErrCode,
					ErrMsg:  global.ErrParamsTotalLimit.ErrMsg,
					Message: ""})
			logs.LogError("%v %v[%v] %v seg_size[%v]", uuid, header.Filename, md5, total, header.Size)
			continue
		}
		/// offset检查
		if !checkOffset(offset, total) {
			result = append(result,
				global.Result{
					Uuid:    uuid,
					File:    header.Filename,
					Md5:     md5,
					ErrCode: global.ErrParamsOffset.ErrCode,
					ErrMsg:  global.ErrParamsOffset.ErrMsg,
					Message: ""})
			logs.LogError("%v %v[%v] %v seg_size[%v]", uuid, header.Filename, md5, total, header.Size)
			continue
		}
		/// md5检查
		if !checkMD5(md5) {
			result = append(result,
				global.Result{
					Uuid:    uuid,
					File:    header.Filename,
					Md5:     md5,
					ErrCode: global.ErrParamsMD5.ErrCode,
					ErrMsg:  global.ErrParamsMD5.ErrMsg,
					Message: ""})
			logs.LogError("%v %v[%v] %v seg_size[%v]", uuid, header.Filename, md5, total, header.Size)
			continue
		}
		fi := fileInfos.Get(md5)
		if fi == nil {
			/// 没有上传，判断能否上传
			size, _ := strconv.ParseInt(total, 10, 0)
			allTotal += size
			if !checkTotal(allTotal) {
				result = append(result,
					global.Result{
						Uuid:    uuid,
						File:    header.Filename,
						Md5:     md5,
						ErrCode: global.ErrParamsAllTotalLimit.ErrCode,
						ErrMsg:  global.ErrParamsAllTotalLimit.ErrMsg,
						Message: ""})
				logs.LogError("%v %v[%v] %v seg_size[%v]", uuid, header.Filename, md5, total, header.Size)
				continue
			}
		} else {
			if fi.Uuid() == uuid {
				/// 已在当前上传任务中
				size, _ := strconv.ParseInt(total, 10, 0)
				allTotal += size
			}
		}
		info, ok := fileInfos.GetAdd(md5, uuid, header.Filename, total)
		if !ok {
			/// 没有上传，等待上传
			keys = append(keys, k)
			logs.LogWarn("--- *** 没有上传，等待上传 %v %v[%v] %v/%v seg_size[%v]", uuid, header.Filename, md5, info.Now(false), total, header.Size)
		} else {
			info.Assert()
			if info.Uuid() == uuid {
				/// 已在当前上传任务中

				////// 校验MD5
				if md5 != info.Md5() {
					logs.LogFatal("%v %v(%v) md5:%v", info.Uuid, info.SrcName, info.Md5, md5)
				}
				////// 校验数据大小
				if total != strconv.FormatInt(info.Total(false), 10) {
					logs.LogFatal("%v %v(%v) info.total:%v total:%v", info.Uuid(), info.SrcName(), info.Md5(), info.Total(false), total)
				}
				if info.Done(true) {
					if ok, url := info.Ok(true); ok {
						info.UpdateHitTime(time.Now())
						// fileInfos.Remove(info.Md5()).Put()
						result = append(result,
							global.Result{
								Uuid:    uuid,
								File:    header.Filename,
								Md5:     info.Md5(),
								Now:     info.Now(true),
								Total:   info.Total(false),
								ErrCode: global.ErrOk.ErrCode,
								ErrMsg:  global.ErrOk.ErrMsg,
								Url:     url,
								Message: strings.Join([]string{info.Uuid(), " uploading ", info.DstName(), " progress:", strconv.FormatInt(info.Now(true), 10) + "/" + total + " 上传成功!"}, "")})
						logs.LogWarn("%v %v[%v] %v chkmd5 [ok] %v", uuid, header.Filename, info.Md5(), info.DstName(), url)
					} else {
						fileInfos.Remove(info.Md5()).Put()
						os.Remove(config.Config.UploadlDir + info.DstName())
						result = append(result,
							global.Result{
								Uuid:    uuid,
								File:    header.Filename,
								Md5:     info.Md5(),
								Now:     info.Now(true),
								Total:   info.Total(false),
								ErrCode: global.ErrFileMd5.ErrCode,
								ErrMsg:  global.ErrFileMd5.ErrMsg,
								Message: strings.Join([]string{info.Uuid(), " uploading ", info.DstName(), " progress:", strconv.FormatInt(info.Now(true), 10) + "/" + total + " 上传完毕 MD5校验失败!"}, "")})
						logs.LogError("%v %v[%v] %v chkmd5 [Err]", uuid, header.Filename, md5, info.DstName())
					}
				} else {
					keys = append(keys, k)
				}
			} else {
				/// 已在其它上传任务中

				if info.Done(true) {
					if ok, url := info.Ok(true); ok {
						info.UpdateHitTime(time.Now())
						// fileInfos.Remove(info.Md5).Put()
						result = append(result,
							global.Result{
								Uuid:    uuid,
								File:    header.Filename,
								Md5:     info.Md5(),
								Now:     info.Now(true),
								Total:   info.Total(false),
								ErrCode: global.ErrOk.ErrCode,
								ErrMsg:  global.ErrOk.ErrMsg,
								Url:     url,
								Message: strings.Join([]string{info.Uuid(), " uploading ", info.DstName(), " progress:", strconv.FormatInt(info.Now(true), 10) + "/" + total + " 别人上传成功!"}, "")})
						logs.LogWarn("%v %v[%v] %v chkmd5 [ok] %v", uuid, header.Filename, info.Md5(), info.DstName(), url)
					} else {
						fileInfos.Remove(info.Md5()).Put()
						os.Remove(config.Config.UploadlDir + info.DstName())
						result = append(result,
							global.Result{
								Uuid:    uuid,
								File:    header.Filename,
								Md5:     info.Md5(),
								Now:     info.Now(true),
								Total:   info.Total(false),
								ErrCode: global.ErrFileMd5.ErrCode,
								ErrMsg:  global.ErrFileMd5.ErrMsg,
								Message: strings.Join([]string{info.Uuid(), " uploading ", info.DstName(), " progress:", strconv.FormatInt(info.Now(true), 10) + "/" + total + " 别人上传完毕 MD5校验失败!"}, "")})
						logs.LogError("%v %v[%v] %v chkmd5 [Err]", uuid, header.Filename, md5, info.DstName())
					}
				} else {
					result = append(result,
						global.Result{
							Uuid:    uuid,
							File:    info.SrcName(),
							Md5:     info.Md5(),
							ErrCode: global.ErrRepeat.ErrCode,
							ErrMsg:  global.ErrRepeat.ErrMsg,
							Message: strings.Join([]string{info.Uuid(), " uploading ", info.DstName(), " progress:", strconv.FormatInt(info.Now(true), 10), "/", total}, " 别人上传中"),
						})
					// logs.LogError("--- *** ignore repeat-upload %v %v[%v] seg_size[%v] %v uploading %v progress:%v/%v", uuid, header.Filename, md5, header.Size, info.Uuid(), info.DstName(), info.Now(), total)
				}
			}
		}
	} /// {{{ end for range form.File
	var exist bool
	if len(keys) > 0 {
		uploader, ok := uploaders.GetAdd(uuid, config.Config.UseAsync > 0)
		if !ok {
			///////////////////////////// 新的上传任务 /////////////////////////////
			/// 有待上传文件，启动新任务
			j, _ := json.Marshal(keys)
			logs.LogTrace("--------------------- ****** 有待上传文件，启动任务 %v ... %v", uuid, string(j))
			uploader.Upload(&global.Req{Uuid: uuid, Keys: keys, W: w, R: r, Resp: resp, Result: result})
		} else {
			exist = true
			///////////////////////////// 当前上传任务 /////////////////////////////
			/// 有待上传文件，加入当前任务
			// j, _ := json.Marshal(keys)
			// logs.LogTrace("--------------------- ****** 有待上传文件，加入任务 %v ... %v", uuid, string(j))
			uploader.Upload(&global.Req{Uuid: uuid, Keys: keys, W: w, R: r, Resp: resp, Result: result})
		}
	} else {
		/// 无待上传文件，直接返回
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
			writeResponse(w, r, resp)
			// logs.LogError("%v %v", uuid, string(j))
		} else {
			writeResponse(w, r, &global.Resp{})
			if exist {
				logs.LogTrace("--------------------- ****** 无待上传文件，当前任务 %v ...", uuid)
			} else {
				logs.LogTrace("--------------------- ****** 无待上传文件，未分配任务 %v ...", uuid)
			}
		}
	}
}
