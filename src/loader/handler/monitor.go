package handler

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/cwloo/gonet/core/base/sub"
	"github.com/cwloo/gonet/logs"
)

func List() {
	sub.Range(func(pid int, args ...any) {
		logs.DebugfP("%v %v", pid, args)
	})
}

func restart(id int, cmd, dir, Exec, conf, log_dir string) {
	logs.Warnf("%v %v %v %v %v %v", id, cmd, dir, Exec, conf, log_dir)
	f, err := exec.LookPath(dir + Exec)
	if err != nil {
		logs.Fatalf(err.Error())
		return
	}
	args := []string{
		cmd,
		fmt.Sprintf("--id=%v", id),
		fmt.Sprintf("--config=%v", conf),
		fmt.Sprintf("--log_dir=%v", log_dir),
	}
	sub.Start(f, args, Monitor, id, cmd, dir, Exec, conf, log_dir)
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
			id := v[0].(int)
			cmd := v[1].(string)
			dir := v[2].(string)
			Exec := v[3].(string)
			conf := v[4].(string)
			log_dir := v[5].(string)
			restart(id, cmd, dir, Exec, conf, log_dir)
		}
	}
}
