package sub

import (
	"context"
	"os"
	"os/exec"
	"strconv"
	"strings"

	"github.com/cwloo/gonet/core/base/sub"
	"github.com/cwloo/gonet/core/base/sys"
	"github.com/cwloo/gonet/core/base/sys/cmd"
	"github.com/cwloo/gonet/logs"
	"github.com/cwloo/gonet/utils"
	"github.com/cwloo/grpc-etcdv3/getcdv3"
	"github.com/cwloo/grpc-etcdv3/getcdv3/gRPCs"
	pb_public "github.com/cwloo/uploader/proto/public"
	"github.com/cwloo/uploader/src/config"
)

func List() {
	sub.Range(func(pid int, v ...any) {
		p := v[0].(*PID)
		uploaders := 0
		pends := 0
		files := 0
		utils.CheckPanic()
		switch p.Name {
		case config.Config.Gate.Name:
			// logs.ErrorfL("etcd[%v] %v:///%v:%v:%v/", strings.Join(config.Config.Etcd.Addr, ","), config.Config.Etcd.Schema, p.Server.Rpc.Node, p.Server.Rpc.Ip, p.Server.Rpc.Port)
			v, _ := getcdv3.GetConn(config.Config.Etcd.Schema, p.Server.Rpc.Node, p.Server.Rpc.Ip, p.Server.Rpc.Port)
			switch v {
			case nil:
			default:
				client := pb_public.NewPeerClient(v.Conn())
				req := &pb_public.NodeInfoReq{}
				resp, err := client.GetNodeInfo(context.Background(), req)
				if err != nil {
					logs.Errorf("%v [%v:%v rpc=%v:%v:%v %v", pid, p.Name, p.Id+1, p.Server.Rpc.Node, p.Server.Rpc.Ip, p.Server.Rpc.Port, err.Error())
					gRPCs.Conns().RemoveBy(err)
					v.Free()
					return
				}
				pends = int(resp.Node.NumOfPends)
				files = int(resp.Node.NumOfFiles)
				uploaders = int(resp.Node.NumOfLoads)
			}
		case config.Config.Gate.Http.Name:
			// logs.ErrorfL("etcd[%v] %v:///%v:%v:%v/", strings.Join(config.Config.Etcd.Addr, ","), config.Config.Etcd.Schema, p.Server.Rpc.Node, p.Server.Rpc.Ip, p.Server.Rpc.Port)
			v, _ := getcdv3.GetConn(config.Config.Etcd.Schema, p.Server.Rpc.Node, p.Server.Rpc.Ip, p.Server.Rpc.Port)
			switch v {
			case nil:
			default:
				client := pb_public.NewPeerClient(v.Conn())
				req := &pb_public.NodeInfoReq{}
				resp, err := client.GetNodeInfo(context.Background(), req)
				if err != nil {
					logs.Errorf("%v [%v:%v rpc=%v:%v:%v %v", pid, p.Name, p.Id+1, p.Server.Rpc.Node, p.Server.Rpc.Ip, p.Server.Rpc.Port, err.Error())
					gRPCs.Conns().RemoveBy(err)
					v.Free()
					return
				}
				pends = int(resp.Node.NumOfPends)
				files = int(resp.Node.NumOfFiles)
				uploaders = int(resp.Node.NumOfLoads)
			}
		case config.Config.File.Name:
			// logs.ErrorfL("etcd[%v] %v:///%v:%v:%v/", strings.Join(config.Config.Etcd.Addr, ","), config.Config.Etcd.Schema, p.Server.Rpc.Node, p.Server.Rpc.Ip, p.Server.Rpc.Port)
			v, _ := getcdv3.GetConn(config.Config.Etcd.Schema, p.Server.Rpc.Node, p.Server.Rpc.Ip, p.Server.Rpc.Port)
			switch v {
			case nil:
			default:
				client := pb_public.NewPeerClient(v.Conn())
				req := &pb_public.NodeInfoReq{}
				resp, err := client.GetNodeInfo(context.Background(), req)
				if err != nil {
					logs.Errorf("%v [%v:%v rpc=%v:%v:%v %v", pid, p.Name, p.Id+1, p.Server.Rpc.Node, p.Server.Rpc.Ip, p.Server.Rpc.Port, err.Error())
					gRPCs.Conns().RemoveBy(err)
					v.Free()
					return
				}
				pends = int(resp.Node.NumOfPends)
				files = int(resp.Node.NumOfFiles)
				uploaders = int(resp.Node.NumOfLoads)
			}
		}
		logs.DebugfP("%v [%v:%v uploaders:%v pending:%v files:%v %v:%v rpc:%v:%v %v %v %v %v]",
			pid,
			p.Name,
			p.Id+1,
			uploaders,
			pends,
			files,
			p.Server.Ip,
			p.Server.Port,
			p.Server.Rpc.Ip,
			p.Server.Rpc.Port,
			p.Dir,
			p.Cmd,
			cmd.FormatConf(p.Conf),
			cmd.FormatLog(p.Log))
	})
}

func restart(pid int, v ...any) {
	p := v[0].(*PID)
	logs.Warnf("%v [%v:%v %v:%v rpc:%v:%v %v %v %v %v]",
		pid,
		p.Name,
		p.Id+1,
		p.Server.Ip,
		p.Server.Port,
		p.Server.Rpc.Ip,
		p.Server.Rpc.Port,
		p.Dir,
		p.Cmd,
		cmd.FormatConf(p.Conf),
		cmd.FormatLog(p.Log))
	f, err := exec.LookPath(sys.CorrectPath(strings.Join([]string{p.Dir, sys.P, p.Exec}, "")))
	if err != nil {
		logs.Fatalf(err.Error())
		return
	}
	args := []string{
		p.Cmd,
		cmd.FormatId(p.Id),
		cmd.FormatConf(p.Conf),
		cmd.FormatLog(p.Log),
	}
	switch p.Name {
	case config.Config.Client.Name:
		switch len(p.Filelist) > 0 {
		case true:
			args = append(args, cmd.FormatArg("n", strconv.Itoa(len(p.Filelist))))
			args = append(args, p.Filelist...)
		}
	}
	sub.Start(f, args, func(pid int, v ...any) {
		p := v[0].(*PID)
		logs.DebugfP("%v [%v:%v %v:%v rpc:%v:%v %v %v %v %v]",
			pid,
			p.Name,
			p.Id+1,
			p.Server.Ip,
			p.Server.Port,
			p.Server.Rpc.Ip,
			p.Server.Rpc.Port,
			p.Dir,
			p.Cmd,
			cmd.FormatConf(p.Conf),
			cmd.FormatLog(p.Log))
	}, Monitor, p)
}

func Monitor(sta *os.ProcessState, v ...any) {
	logs.Infof("")
	switch sta.Success() {
	case false:
		switch sta.ExitCode() {
		case 2:
		case -1:
			fallthrough
		default:
			restart(sta.Pid(), v...)
		}
	}
}

func Succ(pid int, v ...any) {
	p := v[0].(*PID)
	logs.DebugfP("%v [%v:%v %v:%v rpc:%v:%v %v %v %v %v]",
		pid,
		p.Name,
		p.Id+1,
		p.Server.Ip,
		p.Server.Port,
		p.Server.Rpc.Ip,
		p.Server.Rpc.Port,
		p.Dir,
		p.Cmd,
		cmd.FormatConf(p.Conf),
		cmd.FormatLog(p.Log))
}
