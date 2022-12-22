package httpsrv

import (
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/cwloo/gonet/logs"
)

type Handler func(http.ResponseWriter, *http.Request)

// <summary>
// HttpServer
// <summary>
type HttpServer interface {
	Router(pattern string, handler Handler)
	Run(id int, name string)
}

// <summary>
// httpserver
// <summary>
type httpserver struct {
	server *http.Server
}

func NewHttpServer(ip string, port int, timeout int) HttpServer {
	s := &httpserver{
		server: &http.Server{
			Addr:              strings.Join([]string{ip, strconv.Itoa(port)}, ":"),
			Handler:           http.NewServeMux(),
			ReadTimeout:       time.Duration(timeout) * time.Second,
			ReadHeaderTimeout: time.Duration(timeout) * time.Second,
			WriteTimeout:      time.Duration(timeout) * time.Second,
			IdleTimeout:       time.Duration(timeout) * time.Second,
		},
	}
	return s
}

func (s *httpserver) Router(pattern string, handler Handler) {
	if !s.valid() {
		logs.Errorf("error")
		return
	}
	s.mux().HandleFunc(pattern, handler)
}

func (s *httpserver) Run(id int, name string) {
	logs.Infof("%v:%v %v", name, id, s.server.Addr)
	s.server.SetKeepAlivesEnabled(true)
	err := s.server.ListenAndServe()
	if err != nil {
		logs.Fatalf(err.Error())
	}
}

func (s *httpserver) valid() bool {
	return s.server != nil && s.server.Handler != nil
}

func (s httpserver) mux() *http.ServeMux {
	return s.server.Handler.(*http.ServeMux)
}
