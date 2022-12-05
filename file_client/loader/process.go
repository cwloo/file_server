package main

import (
	"os"
	"sync"

	"github.com/cwloo/gonet/logs"
)

var (
	Process = map[int]*os.Process{}
	Lock    = &sync.Mutex{}
	wg      = sync.WaitGroup{}
)

func startProcess(name string, args []string) bool {
	attr := &os.ProcAttr{
		Files: []*os.File{os.Stdin, os.Stdout, os.Stderr},
	}
	p, err := os.StartProcess(name, args, attr)
	if err != nil {
		logs.LogError("%v", err)
		return false
	}
	wg.Add(1)
	Lock.Lock()
	Process[p.Pid] = p
	Lock.Unlock()
	go monitor(p)
	return true
}

func monitor(p *os.Process) {
	sta, err := p.Wait()
	if err != nil {
		logs.LogError(err.Error())
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

func kill(pid int) {
	var err error
	Lock.Lock()
	if p, ok := Process[pid]; ok {
		if pid == p.Pid {
			err = p.Kill()
			if err != nil {
				Lock.Unlock()
				goto ERR
			}
			delete(Process, p.Pid)
			Lock.Unlock()
			goto OK
		}
	}
	Lock.Unlock()
	return
ERR:
	logs.LogError(err.Error())
	return
OK:
	logs.LogError("%v", pid)
}

func killAll() {
	Lock.Lock()
	for _, p := range Process {
		err := p.Kill()
		if err != nil {
			logs.LogError(err.Error())
		} else {
			delete(Process, p.Pid)
		}
	}
	Lock.Unlock()
}

func waitAll() {
	wg.Wait()
}
