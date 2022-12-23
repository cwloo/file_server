package handler

import (
	"os"
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
	case true:
		return &pb_file.FileServerResp{
			Resp: &pb_file.CommonResp{
				Pid:  int32(os.Getpid()),
				Name: global.Name,
				Id:   int32(cmd.Id()) + 1,
				Server: &pb_file.CommonResp_Server{
					Ip:   config.Config.File.Ip,
					Port: int32(config.Config.File.Port[cmd.Id()]),
					Rpc: &pb_file.CommonResp_Server_Rpc{
						Ip:   config.Config.Rpc.Ip,
						Port: int32(config.Config.Rpc.File.Port[cmd.Id()]),
					},
				},
			},
			Md5:     md5,
			ErrCode: 6,
			ErrMsg:  "not founded"}, nil
	default:
		return &pb_file.FileServerResp{
			Resp: &pb_file.CommonResp{
				Pid:  int32(os.Getpid()),
				Name: global.Name,
				Id:   int32(cmd.Id()) + 1,
				Server: &pb_file.CommonResp_Server{
					Ip:   config.Config.File.Ip,
					Port: int32(config.Config.File.Port[cmd.Id()]),
					Rpc: &pb_file.CommonResp_Server_Rpc{
						Ip:   config.Config.Rpc.Ip,
						Port: int32(config.Config.Rpc.File.Port[cmd.Id()]),
					},
				},
			},
			Md5:     md5,
			Dns:     strings.Join([]string{config.Config.File.Ip, strconv.Itoa(config.Config.File.Port[cmd.Id()])}, ""),
			ErrCode: 0,
			ErrMsg:  "ok"}, nil
	}
}
