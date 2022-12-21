package file_server

import (
	"sync"

	"github.com/cwloo/uploader/src/global"
)

var (
	wg sync.WaitGroup
)

func Run(id int) {
	global.Server = &Router{}
	global.RpcServer = &RPCServer{}
	wg.Add(2)
	go global.Server.Run(id)
	go global.RpcServer.Run(id)
	wg.Wait()
}
