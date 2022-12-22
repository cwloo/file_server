package global

import (
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/cwloo/gonet/logs"
	"github.com/cwloo/gonet/utils"
)

// <summary>
// ARG
// <summary>
type ARG struct {
	ID   string
	CONF string
	LOG  string
	C    string
}

func (s *ARG) FormatId(id int) string {
	return strings.Join([]string{"--", Cmd.Arg.ID, "=", strconv.Itoa(id)}, "")
}

func (s *ARG) FormatConf(conf string) string {
	return strings.Join([]string{"--", Cmd.Arg.CONF, "=", conf}, "")
}

func (s *ARG) FormatC(conf string) string {
	return strings.Join([]string{"--", Cmd.Arg.C, "=", conf}, "")
}

func (s *ARG) FormatLogDir(dir string) string {
	return strings.Join([]string{"--", Cmd.Arg.LOG, "=", dir}, "")
}

func (s *ARG) ParseArgs() (id int, dir, conf, log_dir string) {
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
	v, ok := d[s.ID]
	switch ok {
	case true:
		id = utils.Atoi(v)
	}
	v, ok = d[s.CONF]
	switch ok {
	case true:
		conf = v
	}
	v, ok = d[s.C]
	switch ok {
	case true:
		conf = v
	}
	v, ok = d[s.LOG]
	switch ok {
	case true:
		log_dir = v
	}
	switch conf == "" {
	case true:
		p, _, _ := utils.G()
		conf = strings.Join([]string{dir, p, "config", p, "conf.ini"}, "")
	}
	return
}

// <summary>
// CMD
// <summary>
type CMD struct {
	Arg      ARG
	ID       int
	Dir      string
	Conf_Dir string
	Log_Dir  string
}

func (s *CMD) ParseArgs() {
	s.ID, s.Dir, s.Conf_Dir, s.Log_Dir = s.Arg.ParseArgs()
}
