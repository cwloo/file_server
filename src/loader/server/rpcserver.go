package loader

import (
	"context"
	"net"
	"strconv"
	"strings"

	config "github.com/cwloo/uploader/src/config"
	"github.com/cwloo/uploader/src/global"

	"github.com/cwloo/gonet/core/net/conn"
	"github.com/cwloo/gonet/logs"
	"github.com/cwloo/gonet/utils"
	pb_monitor "github.com/cwloo/uploader/proto/monitor"
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

func (s *RPCServer) Run(id int, name string) {
	switch global.Cmd.Rpc {
	case "":
		if id >= len(config.Config.Rpc.Monitor.Port) {
			logs.Fatalf("error id=%v Rpc.Monitor.Port.size=%v", id, len(config.Config.Rpc.Monitor.Port))
		}
		s.addr = config.Config.Rpc.Ip
		s.port = config.Config.Rpc.Monitor.Port[id]
	default:
		addr := conn.ParseAddress(global.Cmd.Rpc)
		switch addr {
		case nil:
			logs.Fatalf("error")
		default:
			s.addr = addr.Ip
			s.port = utils.Atoi(addr.Port)
		}
	}
	s.node = config.Config.Rpc.Monitor.Node
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
	pb_monitor.RegisterMonitorServer(server, s)

	logs.Warnf("%v:%v etcd%v schema=%v node=%v:%v:%v", name, id, s.etcdAddr, s.etcdSchema, s.node, s.addr, s.port)
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

func (r *RPCServer) GetFileServer(_ context.Context, in *pb_monitor.FileServerReq) (*pb_monitor.FileServerResp, error) {
	logs.Debugf("")
	return &pb_monitor.FileServerResp{}, nil
}
