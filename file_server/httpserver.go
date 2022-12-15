package main

import (
	"net/http"
	"time"

	"github.com/cwloo/gonet/logs"
	"github.com/cwloo/uploader/file_server/config"
)

type Handler func(http.ResponseWriter, *http.Request)

// <summary>
// HttpServer
// <summary>
type HttpServer interface {
	Router(pattern string, handler Handler)
	Run()
}

// <summary>
// httpserver
// <summary>
type httpserver struct {
	server *http.Server
}

func NewHttpServer() HttpServer {
	s := &httpserver{
		server: &http.Server{
			Addr:              config.Config.HttpAddr,
			Handler:           http.NewServeMux(),
			ReadTimeout:       time.Duration(config.Config.PendingTimeout) * time.Second,
			ReadHeaderTimeout: time.Duration(config.Config.PendingTimeout) * time.Second,
			WriteTimeout:      time.Duration(config.Config.PendingTimeout) * time.Second,
			IdleTimeout:       time.Duration(config.Config.PendingTimeout) * time.Second,
		},
	}
	return s
}

func (s *httpserver) Router(pattern string, handler Handler) {
	if !s.valid() {
		logs.LogError("error")
		return
	}
	s.mux().HandleFunc(pattern, handler)
}

func (s *httpserver) Run() {
	logs.LogInfo(s.server.Addr)
	s.server.SetKeepAlivesEnabled(true)
	err := s.server.ListenAndServe()
	if err != nil {
		logs.LogFatal(err.Error())
	}
}

func (s *httpserver) valid() bool {
	return s.server != nil && s.server.Handler != nil
}

func (s httpserver) mux() *http.ServeMux {
	return s.server.Handler.(*http.ServeMux)
}
