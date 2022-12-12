package main

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strconv"

	"github.com/cwloo/gonet/logs"
)

func delCache(delType int, md5 string) {
	switch delType {
	case 1:
		// 1-移除未决的文件
		fileInfos.RemoveWithCond(md5, func(info FileInfo) bool {
			return !info.Done(false)
		}, func(info FileInfo) {
			os.Remove(dir_upload + info.DstName())
			uploaders.Get(info.Uuid()).Remove(md5)
			info.Put()
		})
	case 2:
		// 2-移除已上传的文件
		fileInfos.RemoveWithCond(md5, func(info FileInfo) bool {
			if ok, _ := info.Ok(false); ok {
				return true
			}
			return false
		}, func(info FileInfo) {
			os.Remove(dir_upload + info.DstName())
			info.Put()
		})
	}
}

func handlerJsonReq(body []byte) (*DelResp, bool) {
	if len(body) == 0 {
		return &DelResp{ErrCode: 3, ErrMsg: "no body"}, false
	}
	logs.LogWarn("%v", string(body))
	req := DelReq{}
	err := json.Unmarshal(body, &req)
	if err != nil {
		logs.LogError(err.Error())
		return &DelResp{ErrCode: 4, ErrMsg: "parse body error"}, false
	}
	if req.Type != 1 && req.Type != 2 && req.Md5 == "" && len(req.Md5) != 32 {
		return &DelResp{Type: req.Type, Md5: req.Md5, ErrCode: 1, ErrMsg: "parse param error"}, false
	}
	delCache(req.Type, req.Md5)
	return &DelResp{Type: req.Type, Md5: req.Md5, ErrCode: 0, ErrMsg: "ok"}, true
}

func handlerQuery(query url.Values) (*DelResp, bool) {
	var delType int
	var md5 string
	if query.Has("type") && len(query["type"]) > 0 {
		delType, _ = strconv.Atoi(query["type"][0])
	}
	if query.Has("md5") && len(query["md5"]) > 0 {
		md5 = query["md5"][0]
	}
	if delType != 1 && delType != 2 && md5 == "" && len(md5) != 32 {
		return &DelResp{Type: delType, Md5: md5, ErrCode: 1, ErrMsg: "parse param error"}, false
	}
	delCache(delType, md5)
	return &DelResp{Type: delType, Md5: md5, ErrCode: 0, ErrMsg: "ok"}, true
}

func handlerDelCache(w http.ResponseWriter, r *http.Request) {
	logs.LogInfo("%v %v %#v", r.Method, r.URL.String(), r.Header)
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		logs.LogError(err.Error())
		resp := &DelResp{ErrCode: 2, ErrMsg: "read body error"}
		reponse(w, r, resp)
		return
	}
	switch r.Method {
	case "POST":
		switch r.Header.Get("Content-Type") {
		case "application/json":
			resp, _ := handlerJsonReq(body)
			reponse(w, r, resp)
		default:
			resp, _ := handlerQuery(r.URL.Query())
			reponse(w, r, resp)
		}
	case "GET":
		switch r.Header.Get("Content-Type") {
		case "application/json":
			resp, _ := handlerJsonReq(body)
			reponse(w, r, resp)
		default:
			resp, _ := handlerQuery(r.URL.Query())
			reponse(w, r, resp)
		}
	}
}

func reponse(w http.ResponseWriter, r *http.Request, resp *DelResp) {
	j, _ := json.Marshal(resp)
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
