package gate

import (
	"net/http"

	"github.com/cwloo/gonet/logs"
	"github.com/cwloo/uploader/src/config"
	"github.com/cwloo/uploader/src/global/httpsrv"
	"github.com/cwloo/uploader/src/http_gate/handler"
)

// <summary>
// Router
// <summary>
type Router struct {
	server httpsrv.HttpServer
}

func (s *Router) Run(id int) {
	if id >= len(config.Config.Gate.Http.Port) {
		logs.Fatalf("error id=%v Gate.Http.Port.size=%v", id, len(config.Config.Gate.Http.Port))
	}
	s.server = httpsrv.NewHttpServer(
		config.Config.Gate.Http.Ip,
		config.Config.Gate.Http.Port[id],
		config.Config.Gate.Http.IdleTimeout)
	s.server.Router(config.Config.Path.UpdateCfg, s.UpdateConfigReq)
	s.server.Router(config.Config.Path.GetCfg, s.GetConfigReq)
	s.server.Router(config.Config.Gate.Http.Path.Fileserver, s.FileServerReq)
	s.server.Run()
}

func (s *Router) UpdateConfigReq(w http.ResponseWriter, r *http.Request) {
	handler.UpdateCfgReq(w, r)
}

func (s *Router) GetConfigReq(w http.ResponseWriter, r *http.Request) {
	handler.GetCfgReq(w, r)
}

func (s *Router) FileServerReq(w http.ResponseWriter, r *http.Request) {

	// grpcCons := getcdv3.GetDefaultGatewayConn4Unique(config.Config.Etcd.Schema, strings.Join(config.Config.Etcd.Addr, ","), operationID)
	// for _, v := range grpcCons {
	// 	if v.Target() == rpcSvr.target {
	// 		logs.LogDebug("Filter self=%v out", rpcSvr.target)
	// 		continue
	// 	}
	// 	client := pbRelay.NewRelayClient(v)
	// 	req := &pbRelay.MultiTerminalLoginCheckReq{
	// 		OperationID: operationID,
	// 		PlatformID:  int32(ctx.GetPlatformId()),
	// 		UserID:      ctx.GetUserId(),
	// 		SessionID:   ctx.GetSession(),
	// 		Token:       ctx.GetToken()}
	// 	resp, err := client.MultiTerminalLoginCheck(context.Background(), req)
	// 	if err != nil {
	// 		logs.LogError("%v", err.Error())
	// 		continue
	// 	}
	// 	if resp.ErrCode != 0 {
	// 		logs.LogError("%v %v", resp.ErrCode, resp.ErrMsg)
	// 		continue
	// 	}
	// 	logs.LogDebug("%v", resp.String())
	// }
	handler.FileServerReq(w, r)
}
