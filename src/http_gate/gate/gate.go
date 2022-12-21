package gate

import (
	"sync"
)

var (
	wg        sync.WaitGroup
	router    = &Router{}
	rpcserver = &RPCServer{}
)

func Run(id int) {
	wg.Add(2)
	go router.Run(id)
	go rpcserver.Run(id)
	wg.Wait()
}
