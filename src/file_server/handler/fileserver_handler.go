package handler

import (
	"strconv"
	"strings"

	"github.com/cwloo/gonet/core/base/sys/cmd"
	pb_file "github.com/cwloo/uploader/proto/file"
	"github.com/cwloo/uploader/src/config"
	"github.com/cwloo/uploader/src/global"
)

func QueryFileServer(md5 string) (*pb_file.FileServerResp, error) {
	info := global.FileInfos.Get(md5)
	switch info == nil {
	case false:
		return &pb_file.FileServerResp{
			Md5:     md5,
			Dns:     strings.Join([]string{config.Config.File.Ip, strconv.Itoa(config.Config.File.Port[cmd.Id()])}, ""),
			ErrCode: 0,
			ErrMsg:  "ok"}, nil
	}
	return &pb_file.FileServerResp{
		Md5:     md5,
		ErrCode: 6,
		ErrMsg:  "not founded"}, nil
}
