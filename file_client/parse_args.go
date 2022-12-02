package main

import (
	"os"
	"strconv"
	"strings"

	"github.com/cwloo/gonet/logs"
)

func parseargs() (id int, filelist []string) {
	d := map[string]string{}
	for _, v := range os.Args {
		m := strings.Split(v, "=")
		if len(m) == 2 {
			d[m[0]] = m[1]
		}
	}
	if v, ok := d["i"]; ok {
		id, _ = strconv.Atoi(v)
	}
	num := 0
	if v, ok := d["c"]; ok {
		num, _ = strconv.Atoi(v)
	}
	for i := 0; i < num; i++ {
		if v, ok := d[strings.Join([]string{"file", strconv.Itoa(i)}, "")]; ok {
			filelist = append(filelist, v)
		}
	}
	logs.LogWarn("%v", os.Args)
	return
}
