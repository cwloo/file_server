package file_server

import (
	"net/http"

	"github.com/cwloo/gonet/core/net/conn"
	"github.com/cwloo/gonet/logs"
	"github.com/cwloo/gonet/utils"
	"github.com/cwloo/uploader/src/config"
	"github.com/cwloo/uploader/src/file_server/handler"
	"github.com/cwloo/uploader/src/file_server/handler/uploader"
	"github.com/cwloo/uploader/src/global/cmd"
	"github.com/cwloo/uploader/src/global/httpsrv"
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
	switch cmd.Server() {
	case "":
		if id >= len(config.Config.File.Port) {
			logs.Fatalf("error id=%v File.Port.size=%v", id, len(config.Config.File.Port))
		}
		s.server = httpsrv.NewHttpServer(
			config.Config.File.Ip,
			config.Config.File.Port[id],
			config.Config.File.IdleTimeout)
	default:
		addr := conn.ParseAddress(cmd.Server())
		switch addr {
		case nil:
			logs.Fatalf("error")
		default:
			s.server = httpsrv.NewHttpServer(
				addr.Ip,
				utils.Atoi(addr.Port),
				config.Config.File.IdleTimeout)
		}
	}
	s.server.Router(config.Config.Path.UpdateCfg, s.UpdateConfigReq)
	s.server.Router(config.Config.Path.GetCfg, s.GetConfigReq)
	s.server.Router(config.Config.File.Path.Upload, s.UploadReq)
	s.server.Router(config.Config.File.Path.Get, s.GetReq)
	s.server.Router(config.Config.File.Path.Del, s.DelCacheFileReq)
	s.server.Router(config.Config.File.Path.Fileinfo, s.GetFileinfoReq)
	s.server.Router(config.Config.File.Path.FileDetail, s.FileDetailReq)
	s.server.Router(config.Config.File.Path.UuidList, s.UuidListReq)
	s.server.Router(config.Config.File.Path.List, s.ListReq)
	s.server.Run(id, name)
}

func (s *Router) UpdateConfigReq(w http.ResponseWriter, r *http.Request) {
	handler.UpdateCfgReq(w, r)
}

func (s *Router) GetConfigReq(w http.ResponseWriter, r *http.Request) {
	handler.GetCfgReq(w, r)
}

func (s *Router) UploadReq(w http.ResponseWriter, r *http.Request) {
	switch config.Config.File.Upload.MultiFile > 0 {
	case true:
		uploader.MultiUploadReq(w, r)
	default:
		uploader.UploadReq(w, r)
	}
}

func (s *Router) GetReq(w http.ResponseWriter, r *http.Request) {
	// resp := &Resp{
	// 	ErrCode: 0,
	// 	ErrMsg:  "OK",
	// }
	// writeResponse(w, r, resp)
	handler.FileinfoReq(w, r)
}

func (s *Router) DelCacheFileReq(w http.ResponseWriter, r *http.Request) {
	handler.DelCacheFileReq(w, r)
}

func (s *Router) GetFileinfoReq(w http.ResponseWriter, r *http.Request) {
	handler.FileinfoReq(w, r)
}

func (s *Router) FileDetailReq(w http.ResponseWriter, r *http.Request) {
	handler.FileDetailReq(w, r)
}

func (s *Router) UuidListReq(w http.ResponseWriter, r *http.Request) {
	handler.UuidListReq(w, r)
}

func (s *Router) ListReq(w http.ResponseWriter, r *http.Request) {
	handler.ListReq(w, r)
}
