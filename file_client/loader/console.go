package main

import (
	"runtime"
	"strconv"
	"strings"

	"github.com/cwloo/gonet/utils"
)

func onInput(str string) int {
	if str == "" {
		return 0
	}
	switch str[0] {
	case 'c':
		utils.ClearScreen[runtime.GOOS]()
	case 'q':
		utils.ClearScreen[runtime.GOOS]()
		killAll()
		return -1
	case 'k':
		str = strings.ReplaceAll(str, " ", "")
		if len(str) > 2 {
			pid, _ := strconv.Atoi(str[1:])
			kill(pid)
		}
	}
	return 0
}
