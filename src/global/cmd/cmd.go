package cmd

import (
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/cwloo/gonet/logs"
	"github.com/cwloo/gonet/utils"
	"github.com/cwloo/uploader/src/global"
)

var (
	cmd CMD
)

func init() {
	cmd.Arg.ID = []string{"id", "i"}
	cmd.Arg.CONF = []string{"config", "conf", "c"}
	cmd.Arg.LOG = []string{"log_dir", "log-dir", "logdir", "log", "l"}
	cmd.Arg.SERVER = []string{"server"}
	cmd.Arg.RPC = []string{"rpc"}
}

func Dir() string {
	return cmd.Dir
}

func Id() int {
	return cmd.ID
}

func FormatId(id int) string {
	return cmd.Arg.formatId(id)
}

func Conf() string {
	return cmd.Conf
}

func FormatConf(conf string) string {
	return cmd.Arg.formatConf(conf)
}

func Log() string {
	return cmd.Log
}

func FormatLog(log string) string {
	return cmd.Arg.formatLog(log)
}

func Server() string {
	return cmd.Server
}

func FormatServer(server string) string {
	return cmd.Arg.formatServer(server)
}

func Rpc() string {
	return cmd.Rpc
}

func FormatRpc(rpc string) string {
	return cmd.Arg.formatRpc(rpc)
}

func ParseArgs() {
	cmd.parseArgs()
}

// <summary>
// CMD
// <summary>
type CMD struct {
	Arg    ARG
	ID     int
	Dir    string
	Conf   string
	Log    string
	Server string
	Rpc    string
}

func (s *CMD) parseArgs() {
	s.ID, s.Dir, s.Conf, s.Log, s.Server, s.Rpc = s.Arg.parse()
}

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

func (s *ARG) assertId() {
	switch len(cmd.Arg.ID) == 0 {
	case true:
		logs.Fatalf("error")
	}
}

func (s *ARG) assertConf() {
	switch len(cmd.Arg.CONF) == 0 {
	case true:
		logs.Fatalf("error")
	}
}

func (s *ARG) assertLog() {
	switch len(cmd.Arg.LOG) == 0 {
	case true:
		logs.Fatalf("error")
	}
}

func (s *ARG) formatId(id int) string {
	s.assertId()
	return strings.Join([]string{"--", cmd.Arg.ID[0], "=", strconv.Itoa(id)}, "")
}

func (s *ARG) formatConf(conf string) string {
	s.assertConf()
	return strings.Join([]string{"--", cmd.Arg.CONF[0], "=", conf}, "")
}

func (s *ARG) formatLog(log string) string {
	s.assertLog()
	return strings.Join([]string{"--", cmd.Arg.LOG[0], "=", log}, "")
}

func (s *ARG) formatServer(server string) string {
	s.assertLog()
	return strings.Join([]string{"--", cmd.Arg.SERVER[0], "=", server}, "")
}

func (s *ARG) formatRpc(rpc string) string {
	s.assertLog()
	return strings.Join([]string{"--", cmd.Arg.RPC[0], "=", rpc}, "")
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
	dir = filepath.Dir(global.Dir)
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
