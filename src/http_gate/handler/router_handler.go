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

func QueryRouter(md5 string) (*global.RouterResp, bool) {
	// v := getcdv3.GetDefaultConn(config.Config.Etcd.Schema, strings.Join(config.Config.Etcd.Addr, ","), config.Config.Rpc.File.Node)
	grpcCons := getcdv3.GetDefaultConn4Unique(config.Config.Etcd.Schema, strings.Join(config.Config.Etcd.Addr, ","), config.Config.Rpc.File.Node)
	logs.Infof("%v grpcCons.size=%v", md5, len(grpcCons))
	NumOfLoads := map[string]*pb_file.RouterResp{}
	for _, v := range grpcCons {
		client := pb_file.NewFileClient(v)
		req := &pb_file.RouterReq{
			Md5: md5,
		}
		resp, err := client.GetRouter(context.Background(), req)
		if err != nil {
			logs.Errorf(err.Error())
			continue
		}
		switch resp.ErrCode {
		default:
			logs.Errorf("%v %v [%v:%v %v:%v rpc:%v:%v NumOfLoads:%v]", v.Target(),
				resp.Resp.Pid,
				resp.Resp.Name,
				resp.Resp.Id,
				resp.Resp.Server.Ip, resp.Resp.Server.Port,
				resp.Resp.Server.Rpc.Ip, resp.Resp.Server.Rpc.Port,
				resp.NumOfLoads)
			NumOfLoads[resp.Dns] = resp
			continue
		case 0:
			logs.Infof("%v %v [%v:%v %v:%v rpc:%v:%v NumOfLoads:%v]", v.Target(),
				resp.Resp.Pid,
				resp.Resp.Name,
				resp.Resp.Id,
				resp.Resp.Server.Ip, resp.Resp.Server.Port,
				resp.Resp.Server.Rpc.Ip, resp.Resp.Server.Rpc.Port,
				resp.NumOfLoads)
			return &global.RouterResp{
				Node: &global.NodeInfo{
					Pid:  int(resp.Resp.Pid),
					Name: resp.Resp.Name,
					Id:   int(resp.Resp.Id),
					Server: struct {
						Ip   string `json:"ip" form:"ip"`
						Port int    `json:"port" form:"port"`
						Rpc  struct {
							Ip   string `json:"ip" form:"ip"`
							Port int    `json:"port" form:"port"`
						} `json:"rpc" form:"rpc"`
					}{
						Ip:   resp.Resp.Server.Ip,
						Port: int(resp.Resp.Server.Port),
						Rpc: struct {
							Ip   string `json:"ip" form:"ip"`
							Port int    `json:"port" form:"port"`
						}{
							Ip:   resp.Resp.Server.Rpc.Ip,
							Port: int(resp.Resp.Server.Rpc.Port),
						},
					},
				},
				Md5:        md5,
				Dns:        resp.Dns,
				NumOfLoads: int(resp.NumOfLoads),
				ErrCode:    0,
				ErrMsg:     "ok"}, true
		}
	}
	var minRouter *pb_file.RouterResp
	minLoads := -1
	for _, v := range NumOfLoads {
		switch minLoads {
		case -1:
			minRouter = v
			minLoads = int(v.GetNumOfLoads())
		default:
			switch int(v.GetNumOfLoads()) < minLoads {
			case true:
				minRouter = v
				minLoads = int(v.GetNumOfLoads())
			}
		}
	}
	switch minRouter {
	case nil:
		return &global.RouterResp{
			Md5:     md5,
			ErrCode: 6,
			ErrMsg:  "no file_server"}, false
	default:
		return &global.RouterResp{
			Node: &global.NodeInfo{
				Pid:  int(minRouter.Resp.Pid),
				Name: minRouter.Resp.Name,
				Id:   int(minRouter.Resp.Id),
				Server: struct {
					Ip   string `json:"ip" form:"ip"`
					Port int    `json:"port" form:"port"`
					Rpc  struct {
						Ip   string `json:"ip" form:"ip"`
						Port int    `json:"port" form:"port"`
					} `json:"rpc" form:"rpc"`
				}{
					Ip:   minRouter.Resp.Server.Ip,
					Port: int(minRouter.Resp.Server.Port),
					Rpc: struct {
						Ip   string `json:"ip" form:"ip"`
						Port int    `json:"port" form:"port"`
					}{
						Ip:   minRouter.Resp.Server.Rpc.Ip,
						Port: int(minRouter.Resp.Server.Rpc.Port),
					},
				},
			},
			Md5:        md5,
			Dns:        minRouter.Dns,
			NumOfLoads: int(minRouter.NumOfLoads),
			ErrCode:    0,
			ErrMsg:     "ok"}, true
	}
}

func handlerRouterJsonReq(body []byte) (*global.RouterResp, bool) {
	if len(body) == 0 {
		return &global.RouterResp{ErrCode: 3, ErrMsg: "no body"}, false
	}
	logs.Warnf("%v", string(body))
	req := global.RouterReq{}
	err := json.Unmarshal(body, &req)
	if err != nil {
		logs.Errorf(err.Error())
		return &global.RouterResp{ErrCode: 4, ErrMsg: "parse body error"}, false
	}
	if req.Md5 == "" && len(req.Md5) != 32 {
		return &global.RouterResp{Md5: req.Md5, ErrCode: 1, ErrMsg: "parse param error"}, false
	}
	return QueryRouter(req.Md5)
}

func handlerRouterQuery(query url.Values) (*global.RouterResp, bool) {
	var md5 string
	if query.Has("md5") && len(query["md5"]) > 0 {
		md5 = query["md5"][0]
	}
	if md5 == "" && len(md5) != 32 {
		return &global.RouterResp{Md5: md5, ErrCode: 1, ErrMsg: "parse param error"}, false
	}
	return QueryRouter(md5)
}

func RouterReq(w http.ResponseWriter, r *http.Request) {
	logs.Infof("%v %v %#v", r.Method, r.URL.String(), r.Header)
	switch strings.ToUpper(r.Method) {
	case "POST":
		switch r.Header.Get("Content-Type") {
		case "application/json":
			body, err := ioutil.ReadAll(r.Body)
			if err != nil {
				logs.Errorf(err.Error())
				resp := &global.RouterResp{ErrCode: 2, ErrMsg: "read body error"}
				httpsrv.WriteResponse(w, r, resp)
				return
			}
			resp, _ := handlerRouterJsonReq(body)
			httpsrv.WriteResponse(w, r, resp)
		default:
			resp, _ := handlerRouterQuery(r.URL.Query())
			httpsrv.WriteResponse(w, r, resp)
		}
	case "GET":
		switch r.Header.Get("Content-Type") {
		case "application/json":
			body, err := ioutil.ReadAll(r.Body)
			if err != nil {
				logs.Errorf(err.Error())
				resp := &global.RouterResp{ErrCode: 2, ErrMsg: "read body error"}
				httpsrv.WriteResponse(w, r, resp)
				return
			}
			resp, _ := handlerRouterJsonReq(body)
			httpsrv.WriteResponse(w, r, resp)
		default:
			resp, _ := handlerRouterQuery(r.URL.Query())
			httpsrv.WriteResponse(w, r, resp)
		}
	case "OPTIONS":
		switch r.Header.Get("Content-Type") {
		case "application/json":
			body, err := ioutil.ReadAll(r.Body)
			if err != nil {
				logs.Errorf(err.Error())
				resp := &global.RouterResp{ErrCode: 2, ErrMsg: "read body error"}
				httpsrv.WriteResponse(w, r, resp)
				return
			}
			resp, _ := handlerRouterJsonReq(body)
			httpsrv.WriteResponse(w, r, resp)
		default:
			resp, _ := handlerRouterQuery(r.URL.Query())
			httpsrv.WriteResponse(w, r, resp)
		}
	}
}
