package global

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/cwloo/gonet/logs"
	"github.com/cwloo/gonet/utils"
)

func ParseArgs() (id int, dir, conf, log_dir string) {
	logs.Warnf("%v", os.Args)
	d := map[string]string{}
	for _, v := range os.Args {
		m := strings.Split(v, "=")
		if len(m) == 2 {
			m[0] = strings.ToLower(strings.ReplaceAll(m[0], "-", ""))
			d[m[0]] = m[1]
		}
	}
	dir = filepath.Dir(Dir)
	dir = filepath.Dir(dir)
	v, ok := d["id"]
	switch ok {
	case true:
		id = utils.Atoi(v)
	}
	v, ok = d["config"]
	switch ok {
	case true:
		conf = v
	}
	v, ok = d["c"]
	switch ok {
	case true:
		conf = v
	}
	v, ok = d["log_dir"]
	switch ok {
	case true:
		log_dir = v
	}
	switch conf == "" {
	case true:
		p, _, _ := utils.G()
		conf = strings.Join([]string{dir, p, "config", p, "conf.ini"}, "")
	}
	// switch log_dir == "" {
	// case true:
	// 	p, _, _ := utils.G()
	// 	log_dir = strings.Join([]string{dir, p, "log"}, "")
	// }
	return
}
