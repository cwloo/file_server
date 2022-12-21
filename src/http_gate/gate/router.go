package gate

import (
	"net/http"

	"github.com/cwloo/gonet/logs"
	"github.com/cwloo/uploader/src/global/handler"
	"github.com/cwloo/uploader/src/global/httpsrv"
	"github.com/cwloo/uploader/src/http_gate/config"
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

}
