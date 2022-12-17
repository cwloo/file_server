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

func handlerUuidListJsonReq(body []byte) (*global.UuidListResp, bool) {
	if len(body) == 0 {
		return &global.UuidListResp{ErrCode: 3, ErrMsg: "no body"}, false
	}
	logs.LogWarn("%v", string(body))
	req := global.UuidListReq{}
	err := json.Unmarshal(body, &req)
	if err != nil {
		logs.LogError(err.Error())
		return &global.UuidListResp{ErrCode: 4, ErrMsg: "parse body error"}, false
	}
	// if req.Md5 == "" && len(req.Md5) != 32 {
	// 	return &global.UuidListResp{ErrCode: 1, ErrMsg: "parse param error"}, false
	// }
	return QueryCacheUuidList()
}

func handlerUuidListQuery(query url.Values) (*global.UuidListResp, bool) {
	// var md5 string
	// if query.Has("md5") && len(query["md5"]) > 0 {
	// 	md5 = query["md5"][0]
	// }
	// if md5 == "" && len(md5) != 32 {
	// 	return &global.UuidListResp{ErrCode: 1, ErrMsg: "parse param error"}, false
	// }
	return QueryCacheUuidList()
}

func handlerUuidList(w http.ResponseWriter, r *http.Request) {
	logs.LogInfo("%v %v %#v", r.Method, r.URL.String(), r.Header)
	switch strings.ToUpper(r.Method) {
	case "POST":
		switch r.Header.Get("Content-Type") {
		case "application/json":
			body, err := ioutil.ReadAll(r.Body)
			if err != nil {
				logs.LogError(err.Error())
				resp := &global.UuidListResp{ErrCode: 2, ErrMsg: "read body error"}
				writeResponse(w, r, resp)
				return
			}
			resp, _ := handlerUuidListJsonReq(body)
			writeResponse(w, r, resp)
		default:
			resp, _ := handlerUuidListQuery(r.URL.Query())
			writeResponse(w, r, resp)
		}
	case "GET":
		switch r.Header.Get("Content-Type") {
		case "application/json":
			body, err := ioutil.ReadAll(r.Body)
			if err != nil {
				logs.LogError(err.Error())
				resp := &global.UuidListResp{ErrCode: 2, ErrMsg: "read body error"}
				writeResponse(w, r, resp)
				return
			}
			resp, _ := handlerUuidListJsonReq(body)
			writeResponse(w, r, resp)
		default:
			resp, _ := handlerUuidListQuery(r.URL.Query())
			writeResponse(w, r, resp)
		}
	case "OPTIONS":
		switch r.Header.Get("Content-Type") {
		case "application/json":
			body, err := ioutil.ReadAll(r.Body)
			if err != nil {
				logs.LogError(err.Error())
				resp := &global.ListResp{ErrCode: 2, ErrMsg: "read body error"}
				writeResponse(w, r, resp)
				return
			}
			resp, _ := handlerUuidListJsonReq(body)
			writeResponse(w, r, resp)
		default:
			resp, _ := handlerUuidListQuery(r.URL.Query())
			writeResponse(w, r, resp)
		}
	}
}

func handlerListJsonReq(body []byte) (*global.ListResp, bool) {
	if len(body) == 0 {
		return &global.ListResp{ErrCode: 3, ErrMsg: "no body"}, false
	}
	logs.LogWarn("%v", string(body))
	req := global.ListReq{}
	err := json.Unmarshal(body, &req)
	if err != nil {
		logs.LogError(err.Error())
		return &global.ListResp{ErrCode: 4, ErrMsg: "parse body error"}, false
	}
	// if req.Md5 == "" && len(req.Md5) != 32 {
	// 	return &global.ListResp{ErrCode: 1, ErrMsg: "parse param error"}, false
	// }
	return QueryCacheList()
}

func handlerListQuery(query url.Values) (*global.ListResp, bool) {
	// var md5 string
	// if query.Has("md5") && len(query["md5"]) > 0 {
	// 	md5 = query["md5"][0]
	// }
	// if md5 == "" && len(md5) != 32 {
	// 	return &global.ListResp{ErrCode: 1, ErrMsg: "parse param error"}, false
	// }
	return QueryCacheList()
}

func handlerList(w http.ResponseWriter, r *http.Request) {
	logs.LogInfo("%v %v %#v", r.Method, r.URL.String(), r.Header)
	switch strings.ToUpper(r.Method) {
	case "POST":
		switch r.Header.Get("Content-Type") {
		case "application/json":
			body, err := ioutil.ReadAll(r.Body)
			if err != nil {
				logs.LogError(err.Error())
				resp := &global.ListResp{ErrCode: 2, ErrMsg: "read body error"}
				writeResponse(w, r, resp)
				return
			}
			resp, _ := handlerListJsonReq(body)
			writeResponse(w, r, resp)
		default:
			resp, _ := handlerListQuery(r.URL.Query())
			writeResponse(w, r, resp)
		}
	case "GET":
		switch r.Header.Get("Content-Type") {
		case "application/json":
			body, err := ioutil.ReadAll(r.Body)
			if err != nil {
				logs.LogError(err.Error())
				resp := &global.ListResp{ErrCode: 2, ErrMsg: "read body error"}
				writeResponse(w, r, resp)
				return
			}
			resp, _ := handlerListJsonReq(body)
			writeResponse(w, r, resp)
		default:
			resp, _ := handlerListQuery(r.URL.Query())
			writeResponse(w, r, resp)
		}
	case "OPTIONS":
		switch r.Header.Get("Content-Type") {
		case "application/json":
			body, err := ioutil.ReadAll(r.Body)
			if err != nil {
				logs.LogError(err.Error())
				resp := &global.ListResp{ErrCode: 2, ErrMsg: "read body error"}
				writeResponse(w, r, resp)
				return
			}
			resp, _ := handlerListJsonReq(body)
			writeResponse(w, r, resp)
		default:
			resp, _ := handlerListQuery(r.URL.Query())
			writeResponse(w, r, resp)
		}
	}
}
