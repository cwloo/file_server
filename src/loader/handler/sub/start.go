package sub

import (
	"os/exec"
	"strconv"
	"strings"

	"github.com/cwloo/gonet/core/base/sub"
	"github.com/cwloo/gonet/core/base/sys/cmd"
	"github.com/cwloo/gonet/logs"
	"github.com/cwloo/gonet/utils"
	"github.com/cwloo/uploader/src/config"
)

var (
	P, Cmd, Ext = utils.G()
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
	Cmd      string
	Exec     string
	Dir      string
	Conf     string
	Log      string
	Filelist []string
}

func Start() {
	m := map[int][]string{}
	{
		id := 0
		for _, f := range config.Config.Client.Upload.Filelist {
			m[id] = append(m[id], f)
			id++
			if id >= config.Config.Sub.Client.Num {
				id = 0
			}
		}
	}

	subs := map[string]struct {
		Num    int
		Cmd    string
		Dir    string
		Exec   string
		Conf   string
		Log    string
		Server struct {
			Ip   string
			Port []int
			Rpc  struct {
				Ip   string
				Port []int
			}
		}
	}{
		config.Config.Client.Name: {
			Num:  config.Config.Sub.Client.Num,
			Cmd:  strings.Join([]string{Cmd, config.Config.Sub.Client.Exec, Ext}, ""),
			Dir:  cmd.CorrectPath(strings.Join([]string{cmd.Dir(), P, config.Config.Sub.Client.Dir, P}, "")),
			Exec: config.Config.Sub.Client.Exec + Ext,
			Conf: cmd.Conf(),
			Log:  cmd.Log()},
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
			Num:  config.Config.Sub.Gate.Num,
			Cmd:  strings.Join([]string{Cmd, config.Config.Sub.Gate.Exec, Ext}, ""),
			Dir:  cmd.CorrectPath(strings.Join([]string{cmd.Dir(), P, config.Config.Sub.Gate.Dir, P}, "")),
			Exec: config.Config.Sub.Gate.Exec + Ext,
			Conf: cmd.Conf(),
			Log:  cmd.Log()},
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
			Num:  config.Config.Sub.Gate.Http.Num,
			Cmd:  strings.Join([]string{Cmd, config.Config.Sub.Gate.Http.Exec, Ext}, ""),
			Dir:  cmd.CorrectPath(strings.Join([]string{cmd.Dir(), P, config.Config.Sub.Gate.Http.Dir, P}, "")),
			Exec: config.Config.Sub.Gate.Http.Exec + Ext,
			Conf: cmd.Conf(),
			Log:  cmd.Log()},
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
			Num:  config.Config.Sub.File.Num,
			Cmd:  strings.Join([]string{Cmd, config.Config.Sub.File.Exec, Ext}, ""),
			Dir:  cmd.CorrectPath(strings.Join([]string{cmd.Dir(), P, config.Config.Sub.File.Dir, P}, "")),
			Exec: config.Config.Sub.File.Exec + Ext,
			Conf: cmd.Conf(),
			Log:  cmd.Log()},
	}
	n := 0
	for name, Exec := range subs {
		id := 0
		for i := 0; i < Exec.Num; {
			f, err := exec.LookPath(cmd.CorrectPath(strings.Join([]string{Exec.Dir, P, Exec.Exec}, "")))
			if err != nil {
				logs.Fatalf(err.Error())
				return
			}
			// args := strings.Split(strings.Join([]string{
			// 	Exec.Cmd,
			// 	global.FormatId(id),
			// 	global.FormatConf(Exec.Conf),
			// 	global.FormatLog(Exec.Log),
			// }, " "), " ")
			args := []string{
				Exec.Cmd,
				cmd.FormatId(id),
				cmd.FormatConf(Exec.Conf),
				cmd.FormatLog(Exec.Log),
			}
			filelist := []string{}
			switch name {
			case config.Config.Client.Name:
				v, ok := m[id]
				switch ok {
				case true:
					for i, f := range v {
						filelist = append(filelist, cmd.FormatArg(strings.Join([]string{"file", strconv.Itoa(i)}, ""), f))
					}
				}
				args = append(args, cmd.FormatArg("n", strconv.Itoa(len(filelist))))
				args = append(args, filelist...)
				if _, ok := sub.Start(f, args, func(pid int, v ...any) {
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
				}, Monitor, &PID{
					Id:       id,
					Name:     name,
					Cmd:      Exec.Cmd,
					Exec:     Exec.Exec,
					Dir:      Exec.Dir,
					Conf:     Exec.Conf,
					Log:      Exec.Log,
					Filelist: filelist,
				}); ok {
					id++
					i++
					n++
				}
			default:
				if _, ok := sub.Start(f, args, func(pid int, v ...any) {
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
					Log:  Exec.Log,
				}); ok {
					id++
					i++
					n++
				}
			}

		}
	}
	logs.Debugf("Children = Succ[%03d]", n)
}

func WaitAll() {
	sub.WaitAll()
}
