package sub

import (
	"fmt"
	"os/exec"
	"strings"

	"github.com/cwloo/gonet/core/base/sub"
	"github.com/cwloo/gonet/logs"
	"github.com/cwloo/gonet/utils"
	"github.com/cwloo/uploader/src/config"
	"github.com/cwloo/uploader/src/global"
	"github.com/cwloo/uploader/src/loader/handler"
)

func Start() {
	p, cmd, Ext := utils.G()
	subs := map[string]struct {
		Num    int
		Cmd    string
		Dir    string
		Exec   string
		Conf   string
		LogDir string
	}{
		"gate": {
			Num:    config.Config.Sub.Gate.Num,
			Cmd:    strings.Join([]string{cmd, config.Config.Sub.Gate.Exec, Ext}, ""),
			Dir:    strings.Join([]string{global.Cmd.Dir, p, config.Config.Sub.Gate.Dir, p}, ""),
			Exec:   config.Config.Sub.Gate.Exec + Ext,
			Conf:   global.Cmd.Conf_Dir,
			LogDir: global.Cmd.Log_Dir},
		"gate.http": {
			Num:    config.Config.Sub.Gate.Http.Num,
			Cmd:    strings.Join([]string{cmd, config.Config.Sub.Gate.Http.Exec, Ext}, ""),
			Dir:    strings.Join([]string{global.Cmd.Dir, p, config.Config.Sub.Gate.Http.Dir, p}, ""),
			Exec:   config.Config.Sub.Gate.Http.Exec + Ext,
			Conf:   global.Cmd.Conf_Dir,
			LogDir: global.Cmd.Log_Dir},
		"file": {
			Num:    config.Config.Sub.File.Num,
			Cmd:    strings.Join([]string{cmd, config.Config.Sub.File.Exec, Ext}, ""),
			Dir:    strings.Join([]string{global.Cmd.Dir, p, config.Config.Sub.File.Dir, p}, ""),
			Exec:   config.Config.Sub.File.Exec + Ext,
			Conf:   global.Cmd.Conf_Dir,
			LogDir: global.Cmd.Log_Dir},
	}
	n := 0
	for _, Exec := range subs {
		id := 0
		for i := 0; i < Exec.Num; {
			f, err := exec.LookPath(Exec.Dir + Exec.Exec)
			if err != nil {
				logs.Fatalf(err.Error())
				return
			}
			// args := strings.Split(strings.Join([]string{
			// 	m.Cmd,
			// 	fmt.Sprintf("--id=%v", id),
			// 	fmt.Sprintf("--config=%v", Exec.Conf),
			// }, " "), " ")
			args := []string{
				Exec.Cmd,
				fmt.Sprintf("--id=%v", id),
				fmt.Sprintf("--config=%v", Exec.Conf),
				fmt.Sprintf("--log_dir=%v", Exec.LogDir),
			}
			if sub.Start(f, args, handler.Monitor, id, Exec.Cmd, Exec.Dir, Exec.Exec, Exec.Conf, Exec.LogDir) {
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
