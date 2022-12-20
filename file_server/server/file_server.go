package file_server

import (
	"sync"
)

var (
	wg        sync.WaitGroup
	router    = NewRouter()
	rpcserver = &RPCServer{}
)

func Run() {
	wg.Add(2)
	go router.Run()
	go rpcserver.Run(0)
	wg.Wait()
}
