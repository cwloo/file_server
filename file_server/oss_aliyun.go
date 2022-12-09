package main

import (
	"errors"
	"time"

	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"github.com/cwloo/gonet/logs"
	"github.com/cwloo/uploader/file_server/config"
)

type AliyunOSS struct{}

func (*AliyunOSS) UploadFile(info FileInfo) (string, string, error) {
	bucket, err := NewBucket()
	if err != nil {
		logs.LogError(err.Error())
		return "", "", errors.New("function AliyunOSS.NewBucket() Failed, err:" + err.Error())
	}
	localFile := config.Config.UploadlDir + info.DstName()
	// yunFilePath := filepath.Join("uploads", time.Now().Format("2006-01-02")) + "/" + info.SrcName()
	yunFilePath := config.Config.AliyunOSS_BasePath + "/uploads/" + time.Now().Format("2006-01-02") + "/" + info.SrcName()
	start := time.Now()
	logs.LogWarn("start oss %v", start)
	err = bucket.UploadFile(yunFilePath, localFile, 1000*1024, oss.Routines(5)) //bucket.PutObject(yunFilePath, f)
	if err != nil {
		logs.LogError(err.Error())
		TgErrMsg(err.Error())
		return "", "", errors.New("function formUploader.Put() Failed, err:" + err.Error())
	}
	logs.LogWarn("finished oss elapsed:%vs", time.Since(start))
	return config.Config.AliyunOSS_BucketUrl + "/" + yunFilePath, yunFilePath, nil
}

func (*AliyunOSS) DeleteFile(key string) error {
	bucket, err := NewBucket()
	if err != nil {
		logs.LogError(err.Error())
		return errors.New("function AliyunOSS.NewBucket() Failed, err:" + err.Error())
	}
	err = bucket.DeleteObject(key)
	if err != nil {
		logs.LogError(err.Error())
		return errors.New("function bucketManager.Delete() Filed, err:" + err.Error())
	}

	return nil
}

func NewBucket() (*oss.Bucket, error) {
	client, err := oss.New(config.Config.AliyunOSS_Endpoint,
		config.Config.AliyunOSS_AccessKeyId,
		config.Config.AliyunOSS_AccessKeySecret, oss.Timeout(120000, 120000))
	if err != nil {
		logs.LogError(err.Error())
		return nil, err
	}
	bucket, err := client.Bucket(config.Config.AliyunOSS_BucketName)
	if err != nil {
		logs.LogError(err.Error())
		return nil, err
	}

	return bucket, nil
}
