package global

// <summary>
// RPCServer
// <summary>
type RPCServer interface {
	Addr() string
	Port() int
	Node() string
	EtcdSchema() string
	EtcdAddr() []string
	Target() string
	Run(id int)
}
