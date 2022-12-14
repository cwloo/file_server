package main

import (
	"net/http"

	"github.com/cwloo/uploader/file_server/config"
)

func UploadReq(w http.ResponseWriter, r *http.Request) {
	switch config.Config.MultiFile > 0 {
	case true:
		handlerMultiUpload(w, r)
	default:
		handlerUpload(w, r)
	}
}

func GetReq(w http.ResponseWriter, r *http.Request) {
	// resp := &Resp{
	// 	ErrCode: 0,
	// 	ErrMsg:  "OK",
	// }
	// writeResponse(w, r, resp)
	handlerFileinfo(w, r)
}

func DelCacheFileReq(w http.ResponseWriter, r *http.Request) {
	handlerDelCacheFile(w, r)
}

func GetFileinfoReq(w http.ResponseWriter, r *http.Request) {
	handlerFileinfo(w, r)
}

func UpdateConfigReq(w http.ResponseWriter, r *http.Request) {
	handlerUpdateCfg(w, r)
}
