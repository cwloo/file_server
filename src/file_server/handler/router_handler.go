package handler

import (
	"os"

	"github.com/cwloo/gonet/core/base/sys/cmd"
	pb_file "github.com/cwloo/uploader/proto/file"
	"github.com/cwloo/uploader/src/config"
	"github.com/cwloo/uploader/src/global"
)

func QueryRouter(md5 string) (*pb_file.RouterResp, error) {
	info := global.FileInfos.Get(md5)
	switch info {
	case nil:
		return &pb_file.RouterResp{
			Node: &pb_file.NodeInfo{
				Pid:        int32(os.Getpid()),
				Name:       global.Name,
				Id:         int32(cmd.Id()) + 1,
				NumOfLoads: int32(global.Uploaders.Len()),
				Ip:         config.Config.File.Ip,
				Port:       int32(config.Config.File.Port[cmd.Id()]),
				// Domain: strings.Join([]string{"http://", config.Config.File.Ip, ":", strconv.Itoa(config.Config.File.Port[cmd.Id()])}, ""),
				Domain: config.Config.File.Domain[cmd.Id()],
				Rpc: &pb_file.NodeInfo_Rpc{
					Ip:   config.Config.Rpc.Ip,
					Port: int32(config.Config.Rpc.File.Port[cmd.Id()]),
				},
			},
			Md5:     md5,
			ErrCode: 6,
			ErrMsg:  "not exist"}, nil
	default:
		return &pb_file.RouterResp{
			Node: &pb_file.NodeInfo{
				Pid:        int32(os.Getpid()),
				Name:       global.Name,
				Id:         int32(cmd.Id()) + 1,
				NumOfLoads: int32(global.Uploaders.Len()),
				Ip:         config.Config.File.Ip,
				Port:       int32(config.Config.File.Port[cmd.Id()]),
				// Domain: strings.Join([]string{"http://", config.Config.File.Ip, ":", strconv.Itoa(config.Config.File.Port[cmd.Id()])}, ""),
				Domain: config.Config.File.Domain[cmd.Id()],
				Rpc: &pb_file.NodeInfo_Rpc{
					Ip:   config.Config.Rpc.Ip,
					Port: int32(config.Config.Rpc.File.Port[cmd.Id()]),
				},
			},
			Md5:     md5,
			ErrCode: 0,
			ErrMsg:  "ok"}, nil
	}
}
