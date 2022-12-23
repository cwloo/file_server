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
	ID     []string
	CONF   []string
	LOG    []string
	SERVER []string
	RPC    []string
}

func (s *ARG) FormatId(id int) string {
	return strings.Join([]string{"--", Cmd.Arg.ID[0], "=", strconv.Itoa(id)}, "")
}

func (s *ARG) FormatConf(conf string) string {
	return strings.Join([]string{"--", Cmd.Arg.CONF[0], "=", conf}, "")
}

func (s *ARG) FormatLog(dir string) string {
	return strings.Join([]string{"--", Cmd.Arg.LOG[0], "=", dir}, "")
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
	Server   string
	Rpc      string
}

func (s *CMD) ParseArgs() {
	s.ID, s.Dir, s.Conf_Dir, s.Log_Dir, s.Server, s.Rpc = s.Arg.parse()
}

func replaceG(old string) (new string) {
	new = old
	exist := true
LOOP:
	for {
		switch len(new) >= 2 && new[0:2] == "--" {
		case true:
			new = strings.Replace(new, "--", "", 1)
		}
		switch len(new) > 0 && new[0:1] == "-" {
		case true:
			new = strings.Replace(new, "-", "", 1)
		default:
			exist = false
		}
		switch exist {
		case false:
			break LOOP
		}
	}
	return
}

func (s *ARG) parse() (id int, dir, conf, log, server, rpc string) {
	logs.Warnf("%v", os.Args)
	d := map[string]string{}
	for _, v := range os.Args {
		m := strings.Split(v, "=")
		switch len(m) == 2 {
		case true:
			m[0] = replaceG(m[0])
			m[0] = strings.ToLower(m[0])
			d[m[0]] = m[1]
		}
	}
	dir = filepath.Dir(Dir)
	dir = filepath.Dir(dir)
	for _, c := range s.ID {
		v, ok := d[c]
		switch ok {
		case true:
			id = utils.Atoi(v)
		}
	}
	for _, c := range s.CONF {
		v, ok := d[c]
		switch ok {
		case true:
			conf = v
		}
	}
	switch conf == "" {
	case true:
		p, _, _ := utils.G()
		conf = strings.Join([]string{dir, p, "config", p, "conf.ini"}, "")
	}
	for _, c := range s.LOG {
		v, ok := d[c]
		switch ok {
		case true:
			log = v
		}
	}
	for _, c := range s.SERVER {
		v, ok := d[c]
		switch ok {
		case true:
			server = v
		}
	}
	for _, c := range s.RPC {
		v, ok := d[c]
		switch ok {
		case true:
			rpc = v
		}
	}
	return
}
