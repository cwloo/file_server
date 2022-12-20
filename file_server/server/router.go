package file_server

import (
	"net/http"

	"github.com/cwloo/uploader/file_server/config"
	"github.com/cwloo/uploader/file_server/handler"
	"github.com/cwloo/uploader/file_server/httpsrv"
	"github.com/cwloo/uploader/file_server/server/uploader"
)

// <summary>
// Router
// <summary>
type Router struct {
	server httpsrv.HttpServer
}

func NewRouter() *Router {
	s := &Router{
		server: httpsrv.NewHttpServer(),
	}
	return s
}

func (s *Router) Run() {
	s.server.Router(config.Config.Upload.Path.Upload, s.UploadReq)
	s.server.Router(config.Config.Upload.Path.Get, s.GetReq)
	s.server.Router(config.Config.Upload.Path.Get, s.DelCacheFileReq)
	s.server.Router(config.Config.Upload.Path.Fileinfo, s.GetFileinfoReq)
	s.server.Router(config.Config.Upload.Path.UpdateCfg, s.UpdateConfigReq)
	s.server.Router(config.Config.Upload.Path.GetCfg, s.GetConfigReq)
	s.server.Router(config.Config.Upload.Path.FileDetail, s.FileDetailReq)
	s.server.Router(config.Config.Upload.Path.UuidList, s.UuidListReq)
	s.server.Router(config.Config.Upload.Path.List, s.ListReq)
	s.server.Run()
}

func (s *Router) UploadReq(w http.ResponseWriter, r *http.Request) {
	switch config.Config.Upload.MultiFile > 0 {
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

func (s *Router) UpdateConfigReq(w http.ResponseWriter, r *http.Request) {
	handler.UpdateCfgReq(w, r)
}

func (s *Router) GetConfigReq(w http.ResponseWriter, r *http.Request) {
	handler.GetCfgReq(w, r)
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
