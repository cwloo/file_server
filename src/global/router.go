package global

import "github.com/cwloo/uploader/src/global/httpsrv"

// <summary>
// Router
// <summary>
type Router interface {
	Server() httpsrv.HttpServer
	Run(id int, name string)
}
