package main

import (
	"time"

	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"github.com/cwloo/gonet/logs"
	"github.com/cwloo/uploader/file_server/config"
)

type Aliyun struct{}

func (*Aliyun) UploadFile(info FileInfo) (string, string, error) {
	bucket, err := NewBucket()
	if err != nil {
		logs.LogError(err.Error())
		return "", "", err
	}
	localFile := config.Config.UploadlDir + info.DstName()
	// yunFilePath := filepath.Join("uploads", time.Now().Format("2006-01-02")) + "/" + info.SrcName()
	yunFilePath := config.Config.Aliyun_BasePath + "/uploads/" + time.Now().Format("2006-01-02") + "/" + info.SrcName()
	start := time.Now()
	logs.LogWarn("start oss %v", start)
	err = bucket.UploadFile(yunFilePath, localFile, 1000*1024, oss.Routines(5)) //bucket.PutObject(yunFilePath, f)
	if err != nil {
		logs.LogError(err.Error())
		TgErrMsg(err.Error())
		return "", "", err
	}
	logs.LogWarn("finished oss elapsed:%vs", time.Since(start))
	return config.Config.Aliyun_BucketUrl + "/" + yunFilePath, yunFilePath, nil
}

func (*Aliyun) DeleteFile(key string) error {
	bucket, err := NewBucket()
	if err != nil {
		logs.LogError(err.Error())
		return err
	}
	err = bucket.DeleteObject(key)
	if err != nil {
		logs.LogError(err.Error())
		return err
	}

	return nil
}

func NewBucket() (*oss.Bucket, error) {
	client, err := oss.New(config.Config.Aliyun_Endpoint,
		config.Config.Aliyun_AccessKeyId,
		config.Config.Aliyun_AccessKeySecret, oss.Timeout(120000, 120000))
	if err != nil {
		logs.LogError(err.Error())
		return nil, err
	}
	bucket, err := client.Bucket(config.Config.Aliyun_BucketName)
	if err != nil {
		logs.LogError(err.Error())
		return nil, err
	}

	return bucket, nil
}
