package handler

import (
	"runtime"
	"strconv"
	"strings"

	"github.com/cwloo/gonet/core/base/sub"
	"github.com/cwloo/gonet/utils"
)

func OnInput(str string) int {
	if str == "" {
		return 0
	}
	switch str[0] {
	case 'c':
		utils.ClearScreen[runtime.GOOS]()
	case 'q':
		utils.ClearScreen[runtime.GOOS]()
		sub.KillAll()
		return -1
	case 'k':
		str = strings.ReplaceAll(str, " ", "")
		if len(str) > 2 {
			pid, _ := strconv.Atoi(str[1:])
			sub.Kill(pid)
		}
	}
	return 0
}
