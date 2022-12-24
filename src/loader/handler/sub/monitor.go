package sub

import (
	"os"
	"os/exec"
	"strings"

	"github.com/cwloo/gonet/core/base/sub"
	"github.com/cwloo/gonet/core/base/sys/cmd"
	"github.com/cwloo/gonet/logs"
)

func List() {
	sub.Range(func(pid int, v ...any) {
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
	f, err := exec.LookPath(cmd.CorrectPath(strings.Join([]string{p.Dir, P, p.Exec}, "")))
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
