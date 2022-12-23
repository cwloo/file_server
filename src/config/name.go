package config

import (
	"github.com/cwloo/gonet/logs"
	"github.com/cwloo/uploader/src/global"
)

func setServiceName(cb func(*IniConfig) string, c *IniConfig) {
	switch global.Name {
	case "":
		switch cb {
		case nil:
		default:
			global.Name = cb(c)
		}
	}
	switch global.Name == "" {
	case true:
		logs.Fatalf("error")
	}
}

func serviceName() string {
	switch global.Name {
	case "":
		logs.Fatalf("error")
	}
	return global.Name
}
