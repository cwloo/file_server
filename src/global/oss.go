package global

import (
	"mime/multipart"
)

type NewOss func(info FileInfo, ossType string) Oss

// <summary>
// Oss 云存储
// <summary>
type Oss interface {
	UploadFile(info FileInfo, header *multipart.FileHeader) (string, string, *ErrorMsg)
	DeleteFile(key string) error
	Put()
}

// func NewOss(info FileInfo, ossType string) OSS {
// 	switch ossType {
// 	// case "local":
// 	// 	return &Local{}
// 	// case "qiniu":
// 	// 	return &Qiniu{}
// 	// case "tencent-cos":
// 	// 	return &TencentCOS{}
// 	case "aliyun-oss":
// 		return NewAliyun(info)
// 	// case "huawei-obs":
// 	// 	return HuaWeiObs
// 	// case "aws-s3":
// 	// 	return &AwsS3{}
// 	// default:
// 	// 	return &Local{}
// 	default:
// 		return NewAliyun(info)
// 	}
// }

// func UploadDomain(ossType string) string {
// 	switch ossType {
// 	// case "local":
// 	// 	return ""
// 	// case "qiniu":
// 	// 	return config.Config.Qiniu.Bucket + "/"
// 	// case "tencent-cos":
// 	// 	return config.Config.TencentCOS.BaseURL + "/"
// 	case "aliyun-oss":
// 		return config.Config.Oss.Aliyun.BucketUrl + "/"
// 	// case "huawei-obs":
// 	// 	return config.Config.HuaWeiObs.Path + "/"
// 	// case "aws-s3":
// 	// 	return config.Config.Aliyun.BucketUrl + "/"
// 	default:
// 		return ""
// 	}
// }
