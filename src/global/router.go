package global

import "github.com/cwloo/uploader/src/global/httpsrv"

// <summary>
// HTTPServer
// <summary>
type HTTPServer interface {
	Server() httpsrv.HttpServer
	Init(id int, name string)
	Run(id int, name string)
}
