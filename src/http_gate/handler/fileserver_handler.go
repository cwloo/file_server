package handler

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"

	"github.com/cwloo/gonet/logs"
	pb_file "github.com/cwloo/uploader/proto/file"
	"github.com/cwloo/uploader/src/config"
	"github.com/cwloo/uploader/src/global"
	"github.com/cwloo/uploader/src/global/httpsrv"
	"github.com/cwloo/uploader/src/global/pkg/grpc-etcdv3/getcdv3"
)

func QueryFileServer(md5 string) (*global.FileServerResp, bool) {
	grpcCons := getcdv3.GetDefaultConn4Unique(config.Config.Etcd.Schema, strings.Join(config.Config.Etcd.Addr, ","),
		config.Config.Rpc.File.Node, config.Config.Rpc.File.Port)
	logs.Infof("%v grpcCons.size=%v", md5, len(grpcCons))
	for _, v := range grpcCons {
		// if v.Target() == global.RpcServer.Target() {
		// 	continue
		// }
		logs.Warnf("%v", v.Target())
		client := pb_file.NewFileClient(v)
		req := &pb_file.FileServerReq{
			Md5: "",
		}
		resp, err := client.GetFileServer(context.Background(), req)
		if err != nil {
			logs.Errorf(err.Error())
			continue
		}
		if resp.ErrCode != 0 {
			logs.Errorf("%v %v", resp.ErrCode, resp.ErrMsg)
			continue
		}
		if resp.Dns != "" {
			return &global.FileServerResp{
				Md5:     md5,
				Dns:     resp.Dns,
				ErrCode: 0,
				ErrMsg:  "ok"}, true
		}
		logs.Debugf("%v", resp.String())
	}
	return &global.FileServerResp{
		Md5:     md5,
		ErrCode: 6,
		ErrMsg:  "not founded"}, false
}

func handlerFileServerJsonReq(body []byte) (*global.FileServerResp, bool) {
	if len(body) == 0 {
		return &global.FileServerResp{ErrCode: 3, ErrMsg: "no body"}, false
	}
	logs.Warnf("%v", string(body))
	req := global.FileServerReq{}
	err := json.Unmarshal(body, &req)
	if err != nil {
		logs.Errorf(err.Error())
		return &global.FileServerResp{ErrCode: 4, ErrMsg: "parse body error"}, false
	}
	if req.Md5 == "" && len(req.Md5) != 32 {
		return &global.FileServerResp{Md5: req.Md5, ErrCode: 1, ErrMsg: "parse param error"}, false
	}
	return QueryFileServer(req.Md5)
}

func handlerFileServerQuery(query url.Values) (*global.FileServerResp, bool) {
	var md5 string
	if query.Has("md5") && len(query["md5"]) > 0 {
		md5 = query["md5"][0]
	}
	if md5 == "" && len(md5) != 32 {
		return &global.FileServerResp{Md5: md5, ErrCode: 1, ErrMsg: "parse param error"}, false
	}
	return QueryFileServer(md5)
}

func FileServerReq(w http.ResponseWriter, r *http.Request) {
	logs.Infof("%v %v %#v", r.Method, r.URL.String(), r.Header)
	switch strings.ToUpper(r.Method) {
	case "POST":
		switch r.Header.Get("Content-Type") {
		case "application/json":
			body, err := ioutil.ReadAll(r.Body)
			if err != nil {
				logs.Errorf(err.Error())
				resp := &global.FileServerResp{ErrCode: 2, ErrMsg: "read body error"}
				httpsrv.WriteResponse(w, r, resp)
				return
			}
			resp, _ := handlerFileServerJsonReq(body)
			httpsrv.WriteResponse(w, r, resp)
		default:
			resp, _ := handlerFileServerQuery(r.URL.Query())
			httpsrv.WriteResponse(w, r, resp)
		}
	case "GET":
		switch r.Header.Get("Content-Type") {
		case "application/json":
			body, err := ioutil.ReadAll(r.Body)
			if err != nil {
				logs.Errorf(err.Error())
				resp := &global.FileServerResp{ErrCode: 2, ErrMsg: "read body error"}
				httpsrv.WriteResponse(w, r, resp)
				return
			}
			resp, _ := handlerFileServerJsonReq(body)
			httpsrv.WriteResponse(w, r, resp)
		default:
			resp, _ := handlerFileServerQuery(r.URL.Query())
			httpsrv.WriteResponse(w, r, resp)
		}
	case "OPTIONS":
		switch r.Header.Get("Content-Type") {
		case "application/json":
			body, err := ioutil.ReadAll(r.Body)
			if err != nil {
				logs.Errorf(err.Error())
				resp := &global.FileServerResp{ErrCode: 2, ErrMsg: "read body error"}
				httpsrv.WriteResponse(w, r, resp)
				return
			}
			resp, _ := handlerFileServerJsonReq(body)
			httpsrv.WriteResponse(w, r, resp)
		default:
			resp, _ := handlerFileServerQuery(r.URL.Query())
			httpsrv.WriteResponse(w, r, resp)
		}
	}
}
