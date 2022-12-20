package handler

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/cwloo/gonet/logs"
	"github.com/cwloo/uploader/file_server/global"
	"github.com/cwloo/uploader/file_server/httpsrv"
)

func handlerCacheFileJsonReq(body []byte) (*global.DelResp, bool) {
	if len(body) == 0 {
		return &global.DelResp{ErrCode: 3, ErrMsg: "no body"}, false
	}
	logs.Warnf("%v", string(body))
	req := global.DelReq{}
	err := json.Unmarshal(body, &req)
	if err != nil {
		logs.Errorf(err.Error())
		return &global.DelResp{ErrCode: 4, ErrMsg: "parse body error"}, false
	}
	if req.Type != 1 && req.Type != 2 && req.Md5 == "" && len(req.Md5) != 32 {
		return &global.DelResp{Type: req.Type, Md5: req.Md5, ErrCode: 1, ErrMsg: "parse param error"}, false
	}
	DelCacheFile(req.Type, req.Md5)
	return &global.DelResp{Type: req.Type, Md5: req.Md5, ErrCode: 0, ErrMsg: "ok"}, true
}

func handlerCacheFileQuery(query url.Values) (*global.DelResp, bool) {
	var delType int
	var md5 string
	if query.Has("type") && len(query["type"]) > 0 {
		delType, _ = strconv.Atoi(query["type"][0])
	}
	if query.Has("md5") && len(query["md5"]) > 0 {
		md5 = query["md5"][0]
	}
	if delType != 1 && delType != 2 && md5 == "" && len(md5) != 32 {
		return &global.DelResp{Type: delType, Md5: md5, ErrCode: 1, ErrMsg: "parse param error"}, false
	}
	DelCacheFile(delType, md5)
	return &global.DelResp{Type: delType, Md5: md5, ErrCode: 0, ErrMsg: "ok"}, true
}

func DelCacheFileReq(w http.ResponseWriter, r *http.Request) {
	logs.Infof("%v %v %#v", r.Method, r.URL.String(), r.Header)
	switch strings.ToUpper(r.Method) {
	case "POST":
		switch r.Header.Get("Content-Type") {
		case "application/json":
			body, err := ioutil.ReadAll(r.Body)
			if err != nil {
				logs.Errorf(err.Error())
				resp := &global.DelResp{ErrCode: 2, ErrMsg: "read body error"}
				httpsrv.WriteResponse(w, r, resp)
				return
			}
			resp, _ := handlerCacheFileJsonReq(body)
			httpsrv.WriteResponse(w, r, resp)
		default:
			resp, _ := handlerCacheFileQuery(r.URL.Query())
			httpsrv.WriteResponse(w, r, resp)
		}
	case "GET":
		switch r.Header.Get("Content-Type") {
		case "application/json":
			body, err := ioutil.ReadAll(r.Body)
			if err != nil {
				logs.Errorf(err.Error())
				resp := &global.DelResp{ErrCode: 2, ErrMsg: "read body error"}
				httpsrv.WriteResponse(w, r, resp)
				return
			}
			resp, _ := handlerCacheFileJsonReq(body)
			httpsrv.WriteResponse(w, r, resp)
		default:
			resp, _ := handlerCacheFileQuery(r.URL.Query())
			httpsrv.WriteResponse(w, r, resp)
		}
	case "OPTIONS":
		switch r.Header.Get("Content-Type") {
		case "application/json":
			body, err := ioutil.ReadAll(r.Body)
			if err != nil {
				logs.Errorf(err.Error())
				resp := &global.DelResp{ErrCode: 2, ErrMsg: "read body error"}
				httpsrv.WriteResponse(w, r, resp)
				return
			}
			resp, _ := handlerCacheFileJsonReq(body)
			httpsrv.WriteResponse(w, r, resp)
		default:
			resp, _ := handlerCacheFileQuery(r.URL.Query())
			httpsrv.WriteResponse(w, r, resp)
		}
	}
}
