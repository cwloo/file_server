package main

import (
	"mime/multipart"

	"github.com/cwloo/uploader/file_server/config"
)

// <summary>
// OSS 云存储
// <summary>
type OSS interface {
	UploadFile(info FileInfo, header *multipart.FileHeader) (string, string, error)
	DeleteFile(key string) error
	Put()
}

func NewOss(info FileInfo) OSS {
	switch config.Config.OssType {
	// case "local":
	// 	return &Local{}
	// case "qiniu":
	// 	return &Qiniu{}
	// case "tencent-cos":
	// 	return &TencentCOS{}
	case "aliyun-oss":
		return NewAliyun(info)
	// case "huawei-obs":
	// 	return HuaWeiObs
	// case "aws-s3":
	// 	return &AwsS3{}
	// default:
	// 	return &Local{}
	default:
		return NewAliyun(info)
	}
}

func UploadDomain() string {
	switch config.Config.OssType {
	// case "local":
	// 	return ""
	// case "qiniu":
	// 	return config.Config.Qiniu.Bucket + "/"
	// case "tencent-cos":
	// 	return config.Config.TencentCOS.BaseURL + "/"
	case "aliyun-oss":
		return config.Config.Aliyun_BucketUrl + "/"
	// case "huawei-obs":
	// 	return config.Config.HuaWeiObs.Path + "/"
	// case "aws-s3":
	// 	return config.Config.Aliyun.BucketUrl + "/"
	default:
		return ""
	}
}
