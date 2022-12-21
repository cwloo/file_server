package global

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/cwloo/gonet/logs"
	"github.com/cwloo/gonet/utils"
)

func ParseArgs() (id int, parentDir, conf string) {
	d := map[string]string{}
	for _, v := range os.Args {
		m := strings.Split(v, "=")
		if len(m) == 2 {
			m[0] = strings.ReplaceAll(m[0], "-", "")
			d[m[0]] = m[1]
		}
	}
	if v, ok := d["id"]; ok {
		id = utils.Atoi(v)
	}
	if v, ok := d["config"]; ok {
		conf = v
	}
	if v, ok := d["c"]; ok {
		conf = v
	}
	if conf == "" {
		p, _, _ := utils.G()
		dir := filepath.Dir(Dir)
		conf = strings.Join([]string{filepath.Dir(dir), p, "config", p, "conf.ini"}, "")
		parentDir = filepath.Dir(dir)
	}
	logs.Warnf("%v", os.Args)
	return
}
