package main

import (
	"io"
	"mime/multipart"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"github.com/cwloo/gonet/logs"
	"github.com/cwloo/uploader/file_server/config"
)

var (
	uploadFromFile = false
	aliyums        = sync.Pool{
		New: func() any {
			return &Aliyun{}
		},
	}
)

// <summary>
// Aliyun
// <summary>
type Aliyun struct {
	bucket  *oss.Bucket
	imur    *oss.InitiateMultipartUploadResult
	parts   []oss.UploadPart
	yunPath string
	num     int
}

func NewAliyun(info FileInfo) OSS {
	bucket, err := NewBucket()
	if err != nil {
		logs.LogError(err.Error())
		return aliyums.Get().(*Aliyun)
	}
	yunPath := strings.Join([]string{config.Config.Aliyun_BasePath, "/uploads/", info.Date(), "/", info.YunName()}, "")
	imur, err := bucket.InitiateMultipartUpload(yunPath)
	if err != nil {
		logs.LogError(err.Error())
		return aliyums.Get().(*Aliyun)
	}
	s := aliyums.Get().(*Aliyun)
	s.bucket = bucket
	s.imur = &imur
	s.parts = []oss.UploadPart{}
	s.yunPath = yunPath
	return s
}

func (s *Aliyun) valid() bool {
	return s.imur != nil
}

func (s *Aliyun) UploadFile(info FileInfo, header *multipart.FileHeader, done bool) (string, string, error) {
	switch s.valid() {
	case true:
		switch uploadFromFile {
		case true:
			switch WriteFile {
			case true:
				return s.uploadFromFile(info, header, done)
			default:
				return s.uploadFromHeader(info, header, done)
			}
		default:
			return s.uploadFromHeader(info, header, done)
		}
	default:
		return "", "", nil
	}
}

func (s *Aliyun) uploadFromHeader(info FileInfo, header *multipart.FileHeader, done bool) (string, string, error) {
	yunPath := ""
	part, err := header.Open()
	if err != nil {
		logs.LogError(err.Error())
		return "", "", err
	}
	s.num++
	start := time.Now()
	logs.LogWarn("start oss %v", start)
	part_oss, err := s.bucket.UploadPart(*s.imur, part, header.Size, s.num, oss.Routines(5))
	if err != nil {
		_ = part.Close()
		logs.LogError(err.Error())
		TgErrMsg(err.Error())
		return "", "", err
	}
	_ = part.Close()
	s.parts = append(s.parts, part_oss)
	logs.LogWarn("elapsed:%v", time.Since(start))
	switch done {
	case true:
		_, err := s.bucket.CompleteMultipartUpload(*s.imur, s.parts)
		if err != nil {
			logs.LogError(err.Error())
			TgErrMsg(err.Error())
			s.reset()
			return "", "", err
		}
		yunPath = s.yunPath
		s.reset()
	default:
		return "", "", nil
	}
	return strings.Join([]string{config.Config.Aliyun_BucketUrl, "/", yunPath}, ""), yunPath, nil
}

func (s *Aliyun) uploadFromFile(info FileInfo, header *multipart.FileHeader, done bool) (string, string, error) {
	yunPath := ""
	f := dir_upload + info.DstName()
	fd, err := os.OpenFile(f, os.O_RDONLY, 0)
	if err != nil {
		logs.LogError(err.Error())
		return "", "", err
	}
	// _, err = fd.Seek(info.Now()-header.Size, io.SeekStart)
	_, err = fd.Seek(header.Size, io.SeekEnd)
	if err != nil {
		_ = fd.Close()
		logs.LogError(err.Error())
		TgErrMsg(err.Error())
		return "", "", err
	}
	s.num++
	start := time.Now()
	// part_oss, err := s.bucket.UploadPartFromFile(*s.imur, f, info.Now()-header.Size, header.Size, s.num, oss.Routines(5))
	part_oss, err := s.bucket.UploadPart(*s.imur, fd, header.Size, s.num, oss.Routines(5))
	if err != nil {
		_ = fd.Close()
		logs.LogError(err.Error())
		TgErrMsg(err.Error())
		return "", "", err
	}
	_ = fd.Close()
	s.parts = append(s.parts, part_oss)
	logs.LogWarn("elapsed:%v", time.Since(start))
	switch done {
	case true:
		_, err := s.bucket.CompleteMultipartUpload(*s.imur, s.parts)
		if err != nil {
			logs.LogError(err.Error())
			TgErrMsg(err.Error())
			s.reset()
			return "", "", err
		}
		yunPath = s.yunPath
		s.reset()
	default:
		return "", "", nil
	}
	return strings.Join([]string{config.Config.Aliyun_BucketUrl, "/", yunPath}, ""), yunPath, nil
}

func (s *Aliyun) reset() {
	s.bucket = nil
	s.imur = nil
	s.parts = nil
	s.yunPath = ""
	s.num = 0
}

func (s *Aliyun) Put() {
	s.reset()
	aliyums.Put(s)
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
