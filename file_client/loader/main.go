package main

import (
	"flag"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"sync"

	"github.com/cwloo/gonet/logs"
	"github.com/cwloo/gonet/utils"
)

///////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
//.\loader -children= -httpaddr= -wsaddr= -mailboxs= -clients= -baseTest= -deltaClients= -deltaTime= -interval= -timeout=

var children = flag.Int("children", 5, "")
var Process = map[int]*os.Process{}
var Lock = &sync.Mutex{}
var wg sync.WaitGroup

func onInput(str string) int {
	switch str {
	case "c":
		{
			utils.ClearScreen[runtime.GOOS]()
			killAllChild()
			return 0
		}
	case "q":
		{
			utils.ClearScreen[runtime.GOOS]()
			killAllChild()
			return -1
		}
	}
	return 0
}

func killAllChild() {
	Lock.Lock()
	for _, p := range Process {
		err := p.Kill()
		if err != nil {
			logs.LogError("%v", err.Error())
		}
	}
	Process = map[int]*os.Process{}
	Lock.Unlock()
}

// func killChild(p *os.Process) {
// 	err := p.Kill()
// 	if err != nil {
// 		logs.LogError("%v", err.Error())
// 		return
// 	}
// 	Lock.Lock()
// 	delete(Process, p.Pid)
// 	Lock.Unlock()
// }

func waitChild(p *os.Process) {
	sta, err := p.Wait()
	if err != nil {
		logs.LogError("%v", err.Error())
	}
	if p.Pid != sta.Pid() {
		logs.LogFatal("%v %v", p.Pid, sta.Pid())
	}
	if sta.Success() {
		logs.LogDebug("pid:%v exit(%v) succ = %v", sta.Pid(), sta.ExitCode(), sta.String())
	} else {
		logs.LogError("pid:%v exit(%v) failed = %v", sta.Pid(), sta.ExitCode(), sta.String())
	}
	Lock.Lock()
	delete(Process, p.Pid)
	Lock.Unlock()
	wg.Done()
}

func startProcess(name, execStr string) {
	for i := 0; i < *children; i++ {
		cmdLine := execStr
		args := strings.Split(cmdLine, " ")
		attr := &os.ProcAttr{
			Files: []*os.File{os.Stdin, os.Stdout, os.Stderr},
		}
		p, err := os.StartProcess(name, args, attr)
		if err != nil {
			logs.LogError("%v", err)
			continue
		}
		wg.Add(1)
		go waitChild(p)
		Lock.Lock()
		Process[p.Pid] = p
		Lock.Unlock()
	}
}

func main() {
	path, _ := os.Executable()
	dir, exe := filepath.Split(path)
	logs.LogTimezone(logs.MY_CST)
	logs.LogInit(dir+"/logs", logs.LVL_DEBUG, exe, 100000000)
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
		dir += "/../"
		execname = "file_client"
		execStr = "./" + execname
	} else if runtime.GOOS == "windows" {
		dir += "\\..\\"
		execname = "file_client.exe"
		execStr = execname
	}
	f, err := exec.LookPath(dir + execname)
	if err != nil {
		logs.LogError("%v", err)
		return
	}
	startProcess(f, execStr)
	wg.Wait()
	logs.LogDebug("exit...")
	logs.LogClose()
}
