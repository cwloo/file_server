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

func queryFileinfoCache(md5 string) (*global.FileInfoResp, bool) {
	info := fileInfos.Get(md5)
	if info == nil {
		return &global.FileInfoResp{Md5: md5, ErrCode: 5, ErrMsg: "not found"}, false
	}
	return &global.FileInfoResp{
		Uuid:    info.Uuid(),
		File:    info.SrcName(),
		Md5:     md5,
		Now:     info.Now(false),
		Total:   info.Total(false),
		ErrCode: 0,
		ErrMsg:  "ok"}, true
}

func handlerFileinfoJsonReq(body []byte) (*global.FileInfoResp, bool) {
	if len(body) == 0 {
		return &global.FileInfoResp{ErrCode: 3, ErrMsg: "no body"}, false
	}
	logs.LogWarn("%v", string(body))
	req := global.FileInfoReq{}
	err := json.Unmarshal(body, &req)
	if err != nil {
		logs.LogError(err.Error())
		return &global.FileInfoResp{ErrCode: 4, ErrMsg: "parse body error"}, false
	}
	if req.Md5 == "" && len(req.Md5) != 32 {
		return &global.FileInfoResp{Md5: req.Md5, ErrCode: 1, ErrMsg: "parse param error"}, false
	}
	return queryFileinfoCache(req.Md5)
}

func handlerFileinfoQuery(query url.Values) (*global.FileInfoResp, bool) {
	var md5 string
	if query.Has("md5") && len(query["md5"]) > 0 {
		md5 = query["md5"][0]
	}
	if md5 == "" && len(md5) != 32 {
		return &global.FileInfoResp{Md5: md5, ErrCode: 1, ErrMsg: "parse param error"}, false
	}
	return queryFileinfoCache(md5)
}

func handlerFileinfo(w http.ResponseWriter, r *http.Request) {
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
			resp, _ := handlerFileinfoJsonReq(body)
			writeResponse(w, r, resp)
		default:
			resp, _ := handlerFileinfoQuery(r.URL.Query())
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
			resp, _ := handlerFileinfoJsonReq(body)
			writeResponse(w, r, resp)
		default:
			resp, _ := handlerFileinfoQuery(r.URL.Query())
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
			resp, _ := handlerFileinfoJsonReq(body)
			writeResponse(w, r, resp)
		default:
			resp, _ := handlerFileinfoQuery(r.URL.Query())
			writeResponse(w, r, resp)
		}
	}
}
