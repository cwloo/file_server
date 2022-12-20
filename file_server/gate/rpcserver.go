package gate

import (
	"context"
	"net"
	"strconv"
	"strings"

	config "github.com/cwloo/uploader/file_server/config"

	"github.com/cwloo/gonet/logs"
	getcdv3 "github.com/cwloo/gonet/server/pkg/grpc-etcdv3/getcdv3"
	pb_gate "github.com/cwloo/uploader/proto/gate"

	"google.golang.org/grpc"
)

// <summary>
// RPCServer
// <summary>
type RPCServer struct {
	addr       string
	port       int
	node       string
	etcdSchema string
	etcdAddr   []string
	target     string
}

func (s *RPCServer) Run(id int) {
	s.addr = config.Config.Rpc.Ip
	s.port = config.Config.Rpc.Gate.Port[id]
	s.node = config.Config.Rpc.Gate.Node
	s.etcdSchema = config.Config.Etcd.Schema
	s.etcdAddr = config.Config.Etcd.Addr
	listener, err := net.Listen("tcp", s.addr)
	if err != nil {
		logs.Fatalf(err.Error())
	}
	defer listener.Close()
	var opts []grpc.ServerOption
	server := grpc.NewServer(opts...)
	defer server.GracefulStop()
	pb_gate.RegisterGateServer(server, s)

	logs.Warnf("etcd%v schema=%v node=%v:%v:%v", s.etcdAddr, s.etcdSchema, s.node, s.addr, s.port)
	err = getcdv3.RegisterEtcd4Unique(s.etcdSchema, strings.Join(s.etcdAddr, ","), s.addr, s.port, s.node, 10)
	if err != nil {
		errMsg := strings.Join([]string{s.etcdSchema, " ", strings.Join(s.etcdAddr, ","), " ", s.addr, ":", strconv.Itoa(s.port), " ", s.node, " ", err.Error()}, "")
		logs.Fatalf(errMsg)
	}
	s.target = getcdv3.GetTarget(s.etcdSchema, s.addr, s.port, s.node)
	logs.Warnf("target=%v", s.target)
	err = server.Serve(listener)
	if err != nil {
		logs.Fatalf(err.Error())
		return
	}
}

func (r *RPCServer) GetFileServer(_ context.Context, in *pb_gate.FileServerReq) (*pb_gate.FileServerResp, error) {
	logs.Debugf("")
	return nil, nil
}
