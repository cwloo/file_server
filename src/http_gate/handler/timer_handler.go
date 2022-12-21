package handler

import (
	"time"

	"github.com/cwloo/gonet/core/base/task"
	"github.com/cwloo/gonet/core/cb"
	"github.com/cwloo/uploader/src/config"
	"github.com/cwloo/uploader/src/global"
)

func ReadConfig() {
	config.ReadConfig(global.Name, global.Cmd.Conf_Dir)
	task.After(time.Duration(config.Config.Interval)*time.Second, cb.NewFunctor00(func() {
		ReadConfig()
	}))
}
