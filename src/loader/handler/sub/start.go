package sub

import (
	"os/exec"
	"strconv"
	"strings"

	"github.com/cwloo/gonet/core/base/sub"
	"github.com/cwloo/gonet/core/base/sys"
	"github.com/cwloo/gonet/core/base/sys/cmd"
	"github.com/cwloo/gonet/logs"
	"github.com/cwloo/uploader/src/config"
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
			Node string
		}
	}
	Cmd      string
	Exec     string
	Dir      string
	Conf     string
	Log      string
	Filelist []string
}

func SubFilelist() (m map[int][]string) {
	m = map[int][]string{}
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
	return
}

func Start() {
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
				Node string
			}
		}
	}{
		config.Config.Client.Name: {
			Num:  config.Config.Sub.Client.Num,
			Cmd:  strings.Join([]string{sys.Cmd, config.Config.Sub.Client.Exec, sys.Ext}, ""),
			Dir:  sys.CorrectPath(strings.Join([]string{cmd.Root(), sys.P, config.Config.Sub.Client.Dir, sys.P}, "")),
			Exec: config.Config.Sub.Client.Exec + sys.Ext,
			Conf: cmd.Conf(),
			Log:  cmd.Log()},
		config.Config.Gate.Name: {
			Server: struct {
				Ip   string
				Port []int
				Rpc  struct {
					Ip   string
					Port []int
					Node string
				}
			}{
				Ip:   config.Config.Gate.Ip,
				Port: config.Config.Gate.Port,
				Rpc: struct {
					Ip   string
					Port []int
					Node string
				}{
					Ip:   config.Config.Rpc.Ip,
					Port: config.Config.Rpc.Gate.Port,
					Node: config.Config.Rpc.Gate.Node,
				},
			},
			Num:  config.Config.Sub.Gate.Num,
			Cmd:  strings.Join([]string{sys.Cmd, config.Config.Sub.Gate.Exec, sys.Ext}, ""),
			Dir:  sys.CorrectPath(strings.Join([]string{cmd.Root(), sys.P, config.Config.Sub.Gate.Dir, sys.P}, "")),
			Exec: config.Config.Sub.Gate.Exec + sys.Ext,
			Conf: cmd.Conf(),
			Log:  cmd.Log()},
		config.Config.HttpGate.Name: {
			Server: struct {
				Ip   string
				Port []int
				Rpc  struct {
					Ip   string
					Port []int
					Node string
				}
			}{
				Ip:   config.Config.HttpGate.Ip,
				Port: config.Config.HttpGate.Port,
				Rpc: struct {
					Ip   string
					Port []int
					Node string
				}{
					Ip:   config.Config.Rpc.Ip,
					Port: config.Config.Rpc.HttpGate.Port,
					Node: config.Config.Rpc.HttpGate.Node,
				},
			},
			Num:  config.Config.Sub.HttpGate.Num,
			Cmd:  strings.Join([]string{sys.Cmd, config.Config.Sub.HttpGate.Exec, sys.Ext}, ""),
			Dir:  sys.CorrectPath(strings.Join([]string{cmd.Root(), sys.P, config.Config.Sub.HttpGate.Dir, sys.P}, "")),
			Exec: config.Config.Sub.HttpGate.Exec + sys.Ext,
			Conf: cmd.Conf(),
			Log:  cmd.Log()},
		config.Config.File.Name: {
			Server: struct {
				Ip   string
				Port []int
				Rpc  struct {
					Ip   string
					Port []int
					Node string
				}
			}{
				Ip:   config.Config.File.Ip,
				Port: config.Config.File.Port,
				Rpc: struct {
					Ip   string
					Port []int
					Node string
				}{
					Ip:   config.Config.Rpc.Ip,
					Port: config.Config.Rpc.File.Port,
					Node: config.Config.Rpc.File.Node,
				},
			},
			Num:  config.Config.Sub.File.Num,
			Cmd:  strings.Join([]string{sys.Cmd, config.Config.Sub.File.Exec, sys.Ext}, ""),
			Dir:  sys.CorrectPath(strings.Join([]string{cmd.Root(), sys.P, config.Config.Sub.File.Dir, sys.P}, "")),
			Exec: config.Config.Sub.File.Exec + sys.Ext,
			Conf: cmd.Conf(),
			Log:  cmd.Log()},
	}
	n := 0
	m := SubFilelist()
	for name, Exec := range subs {
		id := 0
		for i := 0; i < Exec.Num; {
			f, err := exec.LookPath(sys.CorrectPath(strings.Join([]string{Exec.Dir, sys.P, Exec.Exec}, "")))
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
					switch len(filelist) > 0 {
					case true:
						args = append(args, cmd.FormatArg("n", strconv.Itoa(len(filelist))))
						args = append(args, filelist...)
						_, ok = sub.Start(f, args, Succ, Monitor, &PID{
							Id:       id,
							Name:     name,
							Cmd:      Exec.Cmd,
							Exec:     Exec.Exec,
							Dir:      Exec.Dir,
							Conf:     Exec.Conf,
							Log:      Exec.Log,
							Filelist: filelist,
						})
						switch ok {
						case true:
							id++
							i++
							n++
						}
					}
				}
			default:
				_, ok := sub.Start(f, args, Succ, Monitor, &PID{
					Id:   id,
					Name: name,
					Server: struct {
						Ip   string
						Port int
						Rpc  struct {
							Ip   string
							Port int
							Node string
						}
					}{
						Ip:   Exec.Server.Ip,
						Port: Exec.Server.Port[id],
						Rpc: struct {
							Ip   string
							Port int
							Node string
						}{
							Ip:   Exec.Server.Rpc.Ip,
							Port: Exec.Server.Rpc.Port[id],
							Node: Exec.Server.Rpc.Node,
						},
					},
					Cmd:  Exec.Cmd,
					Exec: Exec.Exec,
					Dir:  Exec.Dir,
					Conf: Exec.Conf,
					Log:  Exec.Log,
				})
				switch ok {
				case true:
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
