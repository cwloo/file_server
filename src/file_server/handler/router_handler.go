package handler

import (
	"os"

	"github.com/cwloo/gonet/core/base/sys/cmd"
	pb_public "github.com/cwloo/uploader/proto/public"
	"github.com/cwloo/uploader/src/config"
	"github.com/cwloo/uploader/src/global"
)

func QueryRouter(md5 string) (*pb_public.RouterResp, error) {
	info := global.FileInfos.Get(md5)
	switch info {
	case nil:
		return &pb_public.RouterResp{
			Node: &pb_public.NodeInfo{
				Pid:        int32(os.Getpid()),
				Name:       global.Name,
				Id:         int32(cmd.Id()) + 1,
				NumOfPends: int32(PendingNum()),
				NumOfFiles: int32(FinishedNum()),
				NumOfLoads: int32(global.Uploaders.Len()),
				Ip:         config.Config.File.Ip,
				Port:       int32(config.Config.File.Port[cmd.Id()]),
				// Domain: strings.Join([]string{"http://", config.Config.File.Ip, ":", strconv.Itoa(config.Config.File.Port[cmd.Id()])}, ""),
				Domain: config.Config.File.Domain[cmd.Id()],
				Rpc: &pb_public.NodeInfo_Rpc{
					Ip:   config.Config.Rpc.Ip,
					Port: int32(config.Config.Rpc.File.Port[cmd.Id()]),
				},
			},
			Md5:     md5,
			ErrCode: 6,
			ErrMsg:  "not exist"}, nil
	default:
		return &pb_public.RouterResp{
			Node: &pb_public.NodeInfo{
				Pid:        int32(os.Getpid()),
				Name:       global.Name,
				Id:         int32(cmd.Id()) + 1,
				NumOfPends: int32(PendingNum()),
				NumOfFiles: int32(FinishedNum()),
				NumOfLoads: int32(global.Uploaders.Len()),
				Ip:         config.Config.File.Ip,
				Port:       int32(config.Config.File.Port[cmd.Id()]),
				// Domain: strings.Join([]string{"http://", config.Config.File.Ip, ":", strconv.Itoa(config.Config.File.Port[cmd.Id()])}, ""),
				Domain: config.Config.File.Domain[cmd.Id()],
				Rpc: &pb_public.NodeInfo_Rpc{
					Ip:   config.Config.Rpc.Ip,
					Port: int32(config.Config.Rpc.File.Port[cmd.Id()]),
				},
			},
			Md5:     md5,
			ErrCode: 0,
			ErrMsg:  "ok"}, nil
	}
}

func GetNodeInfo() (*pb_public.NodeInfoResp, error) {
	return &pb_public.NodeInfoResp{
		Node: &pb_public.NodeInfo{
			Pid:        int32(os.Getpid()),
			Name:       global.Name,
			Id:         int32(cmd.Id()) + 1,
			NumOfPends: int32(PendingNum()),
			NumOfFiles: int32(FinishedNum()),
			NumOfLoads: int32(global.Uploaders.Len()),
			Ip:         config.Config.File.Ip,
			Port:       int32(config.Config.File.Port[cmd.Id()]),
			// Domain: strings.Join([]string{"http://", config.Config.File.Ip, ":", strconv.Itoa(config.Config.File.Port[cmd.Id()])}, ""),
			Domain: config.Config.File.Domain[cmd.Id()],
			Rpc: &pb_public.NodeInfo_Rpc{
				Ip:   config.Config.Rpc.Ip,
				Port: int32(config.Config.Rpc.File.Port[cmd.Id()]),
			},
		},
		ErrCode: 0,
		ErrMsg:  "ok"}, nil
}
