package gate

import (
	"net/http"

	"github.com/cwloo/uploader/file_server/config"
	"github.com/cwloo/uploader/file_server/httpsrv"
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
	s.server.Router(config.Config.Gate.Path.Fileserver, s.FileServerReq)
	s.server.Run()
}

func (s *Router) FileServerReq(w http.ResponseWriter, r *http.Request) {

}
