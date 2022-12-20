package main

import (
	"net/http"

	"github.com/cwloo/uploader/file_server/config"
)

// <summary>
// Router
// <summary>
type Router struct {
	server HttpServer
}

func NewRouter() *Router {
	s := &Router{
		server: NewHttpServer(),
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
		handlerMultiUpload(w, r)
	default:
		handlerUpload(w, r)
	}
}

func (s *Router) GetReq(w http.ResponseWriter, r *http.Request) {
	// resp := &Resp{
	// 	ErrCode: 0,
	// 	ErrMsg:  "OK",
	// }
	// writeResponse(w, r, resp)
	handlerFileinfo(w, r)
}

func (s *Router) DelCacheFileReq(w http.ResponseWriter, r *http.Request) {
	handlerDelCacheFile(w, r)
}

func (s *Router) GetFileinfoReq(w http.ResponseWriter, r *http.Request) {
	handlerFileinfo(w, r)
}

func (s *Router) UpdateConfigReq(w http.ResponseWriter, r *http.Request) {
	handlerUpdateCfg(w, r)
}

func (s *Router) GetConfigReq(w http.ResponseWriter, r *http.Request) {
	handlerGetCfg(w, r)
}

func (s *Router) FileDetailReq(w http.ResponseWriter, r *http.Request) {
	handlerFileDetail(w, r)
}

func (s *Router) UuidListReq(w http.ResponseWriter, r *http.Request) {
	handlerUuidList(w, r)
}

func (s *Router) ListReq(w http.ResponseWriter, r *http.Request) {
	handlerList(w, r)
}
