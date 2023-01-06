package http_gate

import (
	"sync"

	"github.com/cwloo/uploader/src/global"
)

var (
	wg sync.WaitGroup
)

func Run(id int, name string) {
	global.Router = &Router{}
	global.RpcServer = &RPCServer{}
	global.Router.Init(id, name)
	global.RpcServer.Init(id, name)
	wg.Add(2)
	go global.Router.Run(id, name)
	go global.RpcServer.Run(id, name)
	wg.Wait()
}
