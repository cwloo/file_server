package main

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/cwloo/gonet/core/base/sub"
	"github.com/cwloo/gonet/logs"
)

func Restart(id int, cmd, dir, Exec, conf string) {
	logs.Warnf("%v %v %v %v %v", id, cmd, dir, Exec, conf)
	f, err := exec.LookPath(dir + Exec)
	if err != nil {
		logs.Fatalf(err.Error())
		return
	}
	args := []string{
		cmd,
		fmt.Sprintf("--id=%v", id),
		fmt.Sprintf("--config=%v", conf),
	}
	sub.Start(f, args, Monitor, id, cmd, dir, Exec, conf)
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
			Restart(id, cmd, dir, Exec, conf)
		}
	}
}
