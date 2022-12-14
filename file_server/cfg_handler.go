package main

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"

	"github.com/cwloo/gonet/logs"
	"github.com/cwloo/uploader/file_server/global"
)

func handlerUpdateCfgJsonReq(body []byte) (*global.UpdateCfgResp, bool) {
	if len(body) == 0 {
		return &global.UpdateCfgResp{ErrCode: 3, ErrMsg: "no body"}, false
	}
	logs.LogWarn("%v", string(body))
	req := global.UpdateCfgReq{}
	err := json.Unmarshal(body, &req)
	if err != nil {
		logs.LogError(err.Error())
		return &global.UpdateCfgResp{ErrCode: 4, ErrMsg: "parse body error"}, false
	}
	logs.LogDebug("%#v", req)
	return UpdateCfg(&req)
}

func handlerUpdateCfgQuery(query url.Values) (*global.UpdateCfgResp, bool) {
	req := &global.UpdateCfgReq{}
	if query.Has("interval") && len(query["interval"]) > 0 {
		req.Interval = query["interval"][0]
	}
	if query.Has("maxMemory") && len(query["maxMemory"]) > 0 {
		req.MaxMemory = query["maxMemory"][0]
	}
	if query.Has("maxSegmentSize") && len(query["maxSegmentSize"]) > 0 {
		req.MaxSegmentSize = query["maxSegmentSize"][0]
	}
	if query.Has("maxSingleSize") && len(query["maxSingleSize"]) > 0 {
		req.MaxSingleSize = query["maxSingleSize"][0]
	}
	if query.Has("maxTotalSize") && len(query["maxTotalSize"]) > 0 {
		req.MaxTotalSize = query["maxTotalSize"][0]
	}
	if query.Has("pendingTimeout") && len(query["pendingTimeout"]) > 0 {
		req.PendingTimeout = query["pendingTimeout"][0]
	}
	if query.Has("fileExpiredTimeout") && len(query["fileExpiredTimeout"]) > 0 {
		req.FileExpiredTimeout = query["fileExpiredTimeout"][0]
	}
	if query.Has("checkMd5") && len(query["checkMd5"]) > 0 {
		req.CheckMd5 = query["checkMd5"][0]
	}
	if query.Has("writeFile") && len(query["writeFile"]) > 0 {
		req.WriteFile = query["writeFile"][0]
	}
	logs.LogDebug("%#v", req)
	return UpdateCfg(req)
}

func handlerUpdateCfg(w http.ResponseWriter, r *http.Request) {
	logs.LogInfo("%v %v %#v", r.Method, r.URL.String(), r.Header)
	switch strings.ToUpper(r.Method) {
	case "POST":
		switch r.Header.Get("Content-Type") {
		case "application/json":
			body, err := ioutil.ReadAll(r.Body)
			if err != nil {
				logs.LogError(err.Error())
				resp := &global.FileInfoResp{ErrCode: 2, ErrMsg: "read body error"}
				writeResponse(w, r, resp)
				return
			}
			resp, _ := handlerUpdateCfgJsonReq(body)
			writeResponse(w, r, resp)
		default:
			resp, _ := handlerUpdateCfgQuery(r.URL.Query())
			writeResponse(w, r, resp)
		}
	case "GET":
		switch r.Header.Get("Content-Type") {
		case "application/json":
			body, err := ioutil.ReadAll(r.Body)
			if err != nil {
				logs.LogError(err.Error())
				resp := &global.FileInfoResp{ErrCode: 2, ErrMsg: "read body error"}
				writeResponse(w, r, resp)
				return
			}
			resp, _ := handlerUpdateCfgJsonReq(body)
			writeResponse(w, r, resp)
		default:
			resp, _ := handlerUpdateCfgQuery(r.URL.Query())
			writeResponse(w, r, resp)
		}
	case "OPTIONS":
		switch r.Header.Get("Content-Type") {
		case "application/json":
			body, err := ioutil.ReadAll(r.Body)
			if err != nil {
				logs.LogError(err.Error())
				resp := &global.FileInfoResp{ErrCode: 2, ErrMsg: "read body error"}
				writeResponse(w, r, resp)
				return
			}
			resp, _ := handlerUpdateCfgJsonReq(body)
			writeResponse(w, r, resp)
		default:
			resp, _ := handlerUpdateCfgQuery(r.URL.Query())
			writeResponse(w, r, resp)
		}
	}
}
