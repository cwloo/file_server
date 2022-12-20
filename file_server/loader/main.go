package main

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"

	"github.com/cwloo/gonet/core/base/sub"
	"github.com/cwloo/gonet/logs"
	"github.com/cwloo/gonet/server/gate/config"
	"github.com/cwloo/gonet/server/gate/global"
	"github.com/cwloo/gonet/utils"
)

var (
	cmd string
	f   string
)

func monitor(id int, sta *os.ProcessState) {
	if !sta.Success() {
		args := []string{
			cmd,
			fmt.Sprintf("i=%v", id),
		}
		sub.Start(id, f, args, monitor)
	}
}

func main() {
	config.InitConfig()
	logs.SetTimezone(logs.TimeZone(config.Config.Log.Timezone))
	logs.SetMode(logs.Mode(config.Config.Log.Mode))
	logs.SetStyle(logs.Style(config.Config.Log.Style))
	logs.Init(config.Config.Log.Dir, logs.Level(config.Config.Log.Level), global.Exe, 100000000)

	go func() {
		utils.ReadConsole(onInput)
	}()
	// dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	// if err != nil {
	// 	logs.Fatalf("%v", err)
	// }
	dir := global.Dir

	if runtime.GOOS == "linux" {
		dir += "../"
		cmd = "./" + config.Config.Sub.Exec
	} else if runtime.GOOS == "windows" {
		dir += "..\\"
		cmd = config.Config.Sub.Exec + ".exe"
	}
	var err error
	f, err = exec.LookPath(dir + config.Config.Sub.Exec)
	if err != nil {
		logs.Errorf(err.Error())
		return
	}
	n := config.Config.Sub.Num
	for id := 0; id < n; {
		// args := strings.Split(strings.Join([]string{
		// 	cmd,
		// 	fmt.Sprintf("i=%v", id),
		// }, " "), " ")
		args := []string{
			cmd,
			fmt.Sprintf("i=%v", id),
		}
		if sub.Start(id, f, args, monitor) {
			id++
		}
	}
	logs.Debugf("Children = Succ[%03d]", n)
	sub.WaitAll()
	logs.Debugf("exit...")
	logs.Close()
}
