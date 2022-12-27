package handler

import (
	"os"

	"github.com/cwloo/gonet/core/base/sys/cmd"
	pb_public "github.com/cwloo/uploader/proto/public"
	"github.com/cwloo/uploader/src/config"
	"github.com/cwloo/uploader/src/global"
)

func GetNodeInfo() (*pb_public.NodeInfoResp, error) {
	return &pb_public.NodeInfoResp{
		Node: &pb_public.NodeInfo{
			Pid:        int32(os.Getpid()),
			Name:       global.Name,
			Id:         int32(cmd.Id()) + 1,
			NumOfPends: int32(PendingNum()),
			NumOfFiles: int32(FinishedNum()),
			NumOfLoads: int32(global.Uploaders.Len()),
			Ip:         config.Config.Monitor.Ip,
			Port:       int32(config.Config.Monitor.Port[cmd.Id()]),
			Rpc: &pb_public.NodeInfo_Rpc{
				Ip:   config.Config.Rpc.Ip,
				Port: int32(config.Config.Rpc.Monitor.Port[cmd.Id()]),
			},
		},
		ErrCode: 0,
		ErrMsg:  "ok"}, nil
}
