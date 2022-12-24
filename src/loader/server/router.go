package loader

import (
	"net/http"

	"github.com/cwloo/gonet/core/base/sys/cmd"
	"github.com/cwloo/gonet/core/net/conn"
	"github.com/cwloo/gonet/logs"
	"github.com/cwloo/gonet/utils"
	"github.com/cwloo/uploader/src/config"
	"github.com/cwloo/uploader/src/global/httpsrv"
	"github.com/cwloo/uploader/src/loader/handler"
)

// <summary>
// Router
// <summary>
type Router struct {
	server httpsrv.HttpServer
}

func (s *Router) Server() httpsrv.HttpServer {
	return s.server
}

func (s *Router) Run(id int, name string) {
	switch cmd.Arg("server") {
	case "":
		if id >= len(config.Config.Monitor.Port) {
			logs.Fatalf("error id=%v Monitor.Port.size=%v", id, len(config.Config.Monitor.Port))
		}
		s.server = httpsrv.NewHttpServer(
			config.Config.Monitor.Ip,
			config.Config.Monitor.Port[id],
			config.Config.Monitor.IdleTimeout)
	default:
		addr := conn.ParseAddress(cmd.Arg("server"))
		switch addr {
		case nil:
			logs.Fatalf("error")
		default:
			s.server = httpsrv.NewHttpServer(
				addr.Ip,
				utils.Atoi(addr.Port),
				config.Config.Monitor.IdleTimeout)
		}
	}
	s.server.Router(config.Config.Path.UpdateCfg, s.UpdateConfigReq)
	s.server.Router(config.Config.Path.GetCfg, s.GetConfigReq)
	s.server.Router(config.Config.Monitor.Path.Start, s.startReq)
	s.server.Router(config.Config.Monitor.Path.Kill, s.killReq)
	s.server.Router(config.Config.Monitor.Path.KillAll, s.killAllReq)
	s.server.Router(config.Config.Monitor.Path.SubList, s.subListReq)
	s.server.Run(id, name)
}

func (s *Router) UpdateConfigReq(w http.ResponseWriter, r *http.Request) {
	handler.UpdateCfgReq(w, r)
}

func (s *Router) GetConfigReq(w http.ResponseWriter, r *http.Request) {
	handler.GetCfgReq(w, r)
}

func (s *Router) startReq(w http.ResponseWriter, r *http.Request) {

}

func (s *Router) killReq(w http.ResponseWriter, r *http.Request) {

}

func (s *Router) killAllReq(w http.ResponseWriter, r *http.Request) {

}

func (s *Router) subListReq(w http.ResponseWriter, r *http.Request) {

}
