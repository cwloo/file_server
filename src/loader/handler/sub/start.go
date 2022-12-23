package sub

import (
	"os/exec"
	"strings"

	"github.com/cwloo/gonet/core/base/sub"
	"github.com/cwloo/gonet/logs"
	"github.com/cwloo/gonet/utils"
	"github.com/cwloo/uploader/src/config"
	"github.com/cwloo/uploader/src/global"
)

// <summary>
// PID
// <summary>
type PID struct {
	Id     int
	Name   string
	Server struct {
		Ip   string
		Port int
		Rpc  struct {
			Ip   string
			Port int
		}
	}
	Cmd  string
	Exec string
	Dir  string
	Conf string
	Log  string
}

func Start() {
	p, cmd, Ext := utils.G()
	subs := map[string]struct {
		Num    int
		Cmd    string
		Dir    string
		Exec   string
		Conf   string
		LogDir string
		Server struct {
			Ip   string
			Port []int
			Rpc  struct {
				Ip   string
				Port []int
			}
		}
	}{
		config.Config.Gate.Name: {
			Server: struct {
				Ip   string
				Port []int
				Rpc  struct {
					Ip   string
					Port []int
				}
			}{
				Ip:   config.Config.Gate.Ip,
				Port: config.Config.Gate.Port,
				Rpc: struct {
					Ip   string
					Port []int
				}{
					Ip:   config.Config.Rpc.Ip,
					Port: config.Config.Rpc.Gate.Port,
				},
			},
			Num:    config.Config.Sub.Gate.Num,
			Cmd:    strings.Join([]string{cmd, config.Config.Sub.Gate.Exec, Ext}, ""),
			Dir:    strings.Join([]string{global.Cmd.Dir, p, config.Config.Sub.Gate.Dir, p}, ""),
			Exec:   config.Config.Sub.Gate.Exec + Ext,
			Conf:   global.Cmd.Conf_Dir,
			LogDir: global.Cmd.Log_Dir},
		config.Config.Gate.Http.Name: {
			Server: struct {
				Ip   string
				Port []int
				Rpc  struct {
					Ip   string
					Port []int
				}
			}{
				Ip:   config.Config.Gate.Http.Ip,
				Port: config.Config.Gate.Http.Port,
				Rpc: struct {
					Ip   string
					Port []int
				}{
					Ip:   config.Config.Rpc.Ip,
					Port: config.Config.Rpc.Gate.Http.Port,
				},
			},
			Num:    config.Config.Sub.Gate.Http.Num,
			Cmd:    strings.Join([]string{cmd, config.Config.Sub.Gate.Http.Exec, Ext}, ""),
			Dir:    strings.Join([]string{global.Cmd.Dir, p, config.Config.Sub.Gate.Http.Dir, p}, ""),
			Exec:   config.Config.Sub.Gate.Http.Exec + Ext,
			Conf:   global.Cmd.Conf_Dir,
			LogDir: global.Cmd.Log_Dir},
		config.Config.File.Name: {
			Server: struct {
				Ip   string
				Port []int
				Rpc  struct {
					Ip   string
					Port []int
				}
			}{
				Ip:   config.Config.File.Ip,
				Port: config.Config.File.Port,
				Rpc: struct {
					Ip   string
					Port []int
				}{
					Ip:   config.Config.Rpc.Ip,
					Port: config.Config.Rpc.File.Port,
				},
			},
			Num:    config.Config.Sub.File.Num,
			Cmd:    strings.Join([]string{cmd, config.Config.Sub.File.Exec, Ext}, ""),
			Dir:    strings.Join([]string{global.Cmd.Dir, p, config.Config.Sub.File.Dir, p}, ""),
			Exec:   config.Config.Sub.File.Exec + Ext,
			Conf:   global.Cmd.Conf_Dir,
			LogDir: global.Cmd.Log_Dir},
	}
	n := 0
	for name, Exec := range subs {
		id := 0
		for i := 0; i < Exec.Num; {
			f, err := exec.LookPath(Exec.Dir + Exec.Exec)
			if err != nil {
				logs.Fatalf(err.Error())
				return
			}
			// args := strings.Split(strings.Join([]string{
			// 	Exec.Cmd,
			// 	global.FormatId(id),
			// 	global.FormatConf(Exec.Conf),
			// 	global.FormatLog(Exec.LogDir),
			// }, " "), " ")
			args := []string{
				Exec.Cmd,
				global.Cmd.Arg.FormatId(id),
				global.Cmd.Arg.FormatConf(Exec.Conf),
				global.Cmd.Arg.FormatLog(Exec.LogDir),
			}
			if _, ok := sub.Start(f, args, func(pid int, v ...any) {
				p := v[0].(*PID)
				logs.DebugfP("%v [%v:%v %v:%v rpc:%v:%v %v %v %v %v]",
					pid,
					p.Name,
					p.Id,
					p.Server.Ip,
					p.Server.Port,
					p.Server.Rpc.Ip,
					p.Server.Rpc.Port,
					p.Dir,
					p.Cmd,
					global.Cmd.Arg.FormatConf(p.Conf),
					global.Cmd.Arg.FormatLog(p.Log))
			}, Monitor, &PID{
				Id:   id,
				Name: name,
				Server: struct {
					Ip   string
					Port int
					Rpc  struct {
						Ip   string
						Port int
					}
				}{
					Ip:   Exec.Server.Ip,
					Port: Exec.Server.Port[id],
					Rpc: struct {
						Ip   string
						Port int
					}{
						Ip:   Exec.Server.Rpc.Ip,
						Port: Exec.Server.Rpc.Port[id],
					},
				},
				Cmd:  Exec.Cmd,
				Exec: Exec.Exec,
				Dir:  Exec.Dir,
				Conf: Exec.Conf,
				Log:  Exec.LogDir,
			}); ok {
				id++
				i++
				n++
			}
		}
	}
	logs.Debugf("Children = Succ[%03d]", n)
}

func WaitAll() {
	sub.WaitAll()
}
