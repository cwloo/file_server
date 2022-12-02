package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"

	"github.com/cwloo/gonet/logs"
	"github.com/cwloo/gonet/utils"
)

var (
	path, _  = os.Executable()
	dir, exe = filepath.Split(path)
)

func main() {
	InitConfig()
	logs.LogTimezone(logs.MY_CST)
	logs.LogInit(dir+"logs", logs.LVL_DEBUG, exe, 100000000)
	logs.LogMode(logs.M_STDOUT_FILE)
	go func() {
		utils.ReadConsole(onInput)
	}()
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		logs.LogFatal("%v", err)
	}
	var execname, execStr string
	if runtime.GOOS == "linux" {
		dir += "../"
		execname = "file_client"
		execStr = "./" + execname
	} else if runtime.GOOS == "windows" {
		dir += "..\\"
		execname = "file_client.exe"
		execStr = execname
	}
	f, err := exec.LookPath(dir + execname)
	if err != nil {
		logs.LogError("%v", err)
		return
	}
	//子进程数量
	n := Config.Sub
	sub := map[int][]string{}
	// 给子进程均匀分配要上传的文件
	id := 0
	for _, f := range Config.FileList {
		sub[id] = append(sub[id], f)
		id++
		if id >= n {
			id = 0
		}
	}
	// 创建若干子进程
	for id := 0; id < n; {
		// 当前子进程要上传的文件
		c := 0
		file := []string{}
		if v, ok := sub[id]; ok {
			c = len(v)
			for i, f := range v {
				file = append(file, fmt.Sprintf("-file%v=%v", i, f))
			}
		}
		// 子进程参数
		// args := strings.Split(strings.Join([]string{
		// 	execStr,
		// 	fmt.Sprintf("-id=%v", id),
		// 	fmt.Sprintf("-c=%v", c),
		// }, " "), " ")
		args := []string{
			execStr,
			fmt.Sprintf("-id=%v", id),
			fmt.Sprintf("-c=%v", c),
		}
		args = append(args, file...)
		// 启动子进程
		if startProcess(f, args) {
			id++
		}
	}
	waitAll()
	logs.LogDebug("exit...")
	logs.LogClose()
}
