package http_gate

import (
	"context"
	"net"
	"strconv"
	"strings"

	config "github.com/cwloo/uploader/src/config"

	"github.com/cwloo/gonet/logs"
	pb_httpgate "github.com/cwloo/uploader/proto/gate.http"
	getcdv3 "github.com/cwloo/uploader/src/global/pkg/grpc-etcdv3/getcdv3"

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

func (s *RPCServer) Addr() string {
	return s.addr
}

func (s *RPCServer) Port() int {
	return s.port
}

func (s *RPCServer) Node() string {
	return s.node
}

func (s *RPCServer) EtcdSchema() string {
	return s.etcdSchema
}

func (s *RPCServer) EtcdAddr() []string {
	return s.etcdAddr
}

func (s *RPCServer) Target() string {
	return s.target
}

func (s *RPCServer) Run(id int) {
	if id >= len(config.Config.Rpc.Gate.Http.Port) {
		logs.Fatalf("error id=%v Rpc.Gate.Http.Port.size=%v", id, len(config.Config.Rpc.Gate.Http.Port))
	}
	s.addr = config.Config.Rpc.Ip
	s.port = config.Config.Rpc.Gate.Http.Port[id]
	s.node = config.Config.Rpc.Gate.Http.Node
	s.etcdSchema = config.Config.Etcd.Schema
	s.etcdAddr = config.Config.Etcd.Addr
	listener, err := net.Listen("tcp", strings.Join([]string{s.addr, strconv.Itoa(s.port)}, ":"))
	if err != nil {
		logs.Fatalf(err.Error())
	}
	defer listener.Close()
	var opts []grpc.ServerOption
	server := grpc.NewServer(opts...)
	defer server.GracefulStop()
	pb_httpgate.RegisterHttpGateServer(server, s)

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

func (r *RPCServer) GetFileServer(_ context.Context, in *pb_httpgate.FileServerReq) (*pb_httpgate.FileServerResp, error) {
	logs.Debugf("")
	return &pb_httpgate.FileServerResp{}, nil
}
