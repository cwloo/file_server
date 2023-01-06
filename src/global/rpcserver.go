package global

// <summary>
// RPCServer
// <summary>
type RPCServer interface {
	Addr() string
	Port() int
	Node() string
	Schema() string
	Target() string
	Init(id int, name string)
	Run(id int, name string)
}
