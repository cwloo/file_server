package main

import "github.com/cwloo/uploader/file_server/config"

func Register(server HttpServer) {
	server.Router(config.Config.UploadPath, UploadReq)
	server.Router(config.Config.GetPath, GetReq)
	server.Router(config.Config.DelPath, DelCacheFileReq)
	server.Router(config.Config.FileinfoPath, GetFileinfoReq)
	server.Router(config.Config.UpdateCfgPath, UpdateConfigReq)
	server.Router(config.Config.GetCfgPath, GetConfigReq)
}
