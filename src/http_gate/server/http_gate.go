package http_gate

import (
	"sync"

	"github.com/cwloo/uploader/src/global"
)

var (
	wg sync.WaitGroup
)

func Run(id int, name string) {
	global.Server = &Router{}
	global.RpcServer = &RPCServer{}
	wg.Add(2)
	go global.Server.Run(id, name)
	go global.RpcServer.Run(id, name)
	wg.Wait()
}
