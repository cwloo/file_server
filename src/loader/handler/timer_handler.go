package handler

import (
	"time"

	"github.com/cwloo/gonet/core/base/task"
	"github.com/cwloo/gonet/core/cb"
	"github.com/cwloo/uploader/src/config"
	"github.com/cwloo/uploader/src/global/cmd"
)

func ReadConfig() {
	config.ReadConfig(cmd.Conf())
	task.After(time.Duration(config.Config.Interval)*time.Second, cb.NewFunctor00(func() {
		ReadConfig()
	}))
}
