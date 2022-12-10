package main

import (
	"mime/multipart"
	"strings"
	"time"

	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"github.com/cwloo/gonet/logs"
	"github.com/cwloo/uploader/file_server/config"
)

type Aliyun struct {
	bucket      *oss.Bucket
	yunFilePath string
	imur        oss.InitiateMultipartUploadResult
	parts       []oss.UploadPart
	num         int
}

func NewAliyun(info FileInfo) OSS {
	bucket, err := NewBucket()
	if err != nil {
		logs.LogFatal(err.Error())
	}
	yunFilePath := strings.Join([]string{config.Config.Aliyun_BasePath, "/uploads/", info.Date(), "/", info.YunName()}, "")
	imur, err := bucket.InitiateMultipartUpload(yunFilePath)
	if err != nil {
		logs.LogFatal(err.Error())
	}
	s := &Aliyun{
		bucket:      bucket,
		yunFilePath: yunFilePath,
		imur:        imur,
		parts:       []oss.UploadPart{}}
	return s
}

func (s *Aliyun) UploadFile(info FileInfo, header *multipart.FileHeader, done bool) (string, string, error) {
	part, err := header.Open()
	if err != nil {
		logs.LogError(err.Error())
		return "", "", err
	}
	// f := dir_upload_tmp + info.DstName()
	// fd, err := os.OpenFile(f, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0666)
	// if err != nil {
	// 	_ = part.Close()
	// 	logs.LogError(err.Error())
	// 	return "", "", err
	// }
	// _, err = io.Copy(fd, part)
	// if err != nil {
	// 	_ = part.Close()
	// 	_ = fd.Close()
	// 	os.Remove(f)
	// 	logs.LogError(err.Error())
	// 	return "", "", err
	// }
	s.num++
	start := time.Now()
	logs.LogWarn("start oss %v", start)
	// err = s.bucket.UploadFile(s.yunFilePath, f, header.Size, oss.Routines(5))
	part_oss, err := s.bucket.UploadPart(s.imur, part, header.Size, s.num, oss.Routines(5))
	if err != nil {
		_ = part.Close()
		// _ = fd.Close()
		// os.Remove(f)
		logs.LogError(err.Error())
		TgErrMsg(err.Error())
		return "", "", err
	}
	_ = part.Close()
	// _ = fd.Close()
	// os.Remove(f)
	s.parts = append(s.parts, part_oss)
	if done {
		_, err := s.bucket.CompleteMultipartUpload(s.imur, s.parts)
		if err != nil {
			logs.LogError(err.Error())
			TgErrMsg(err.Error())
			return "", "", err
		}
	}
	logs.LogWarn("finished oss elapsed:%vs", time.Since(start))
	return config.Config.Aliyun_BucketUrl + "/" + s.yunFilePath, s.yunFilePath, nil
}

func NewBucket() (*oss.Bucket, error) {
	client, err := oss.New(config.Config.Aliyun_Endpoint,
		config.Config.Aliyun_AccessKeyId,
		config.Config.Aliyun_AccessKeySecret, oss.Timeout(120000, 120000))
	if err != nil {
		return nil, err
	}
	bucket, err := client.Bucket(config.Config.Aliyun_BucketName)
	if err != nil {
		return nil, err
	}
	return bucket, nil
}

func (s *Aliyun) DeleteFile(key string) error {
	err := s.bucket.DeleteObject(key)
	if err != nil {
		logs.LogError(err.Error())
		return err
	}
	return nil
}
