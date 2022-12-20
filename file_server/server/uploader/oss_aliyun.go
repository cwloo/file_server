package uploader

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
	"github.com/cwloo/uploader/file_server/global"
	"github.com/cwloo/uploader/file_server/tg_bot"
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
}

func NewAliyun(info global.FileInfo) global.Oss {
	bucket, err := NewBucket()
	if err != nil {
		errMsg := strings.Join([]string{info.Uuid(), " ", info.SrcName(), "[", info.Md5(), "] ", info.YunName(), "\n", "NewBucket:", err.Error()}, "")
		logs.Errorf(errMsg)
		tg_bot.TgErrMsg(errMsg)
		return aliyums.Get().(*Aliyun)
	}
	yunPath := strings.Join([]string{config.Config.Oss.Aliyun.BasePath, "/uploads/", info.Date(), "/", info.YunName()}, "")
	imur, err := bucket.InitiateMultipartUpload(yunPath)
	if err != nil {
		errMsg := strings.Join([]string{info.Uuid(), " ", info.SrcName(), "[", info.Md5(), "] ", info.YunName(), "\n", "InitiateMultipartUpload:", err.Error()}, "")
		logs.Errorf(errMsg)
		tg_bot.TgErrMsg(errMsg)
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

func (s *Aliyun) UploadFile(info global.FileInfo, header *multipart.FileHeader) (string, string, *global.ErrorMsg) {
	switch s.valid() {
	case true:
		switch uploadFromFile {
		case true:
			switch config.Config.Upload.WriteFile > 0 {
			case true:
				return s.uploadFromFile(info, header)
			default:
				return s.uploadFromHeader(info, header)
			}
		default:
			return s.uploadFromHeader(info, header)
		}
	default:
		return "", "", nil
	}
}

func (s *Aliyun) uploadFromHeader(info global.FileInfo, header *multipart.FileHeader) (string, string, *global.ErrorMsg) {
	yunPath := ""
	part, err := header.Open()
	if err != nil {
		errMsg := strings.Join([]string{info.Uuid(), " ", info.SrcName(), "[", info.Md5(), "] ", info.YunName(), "\n", "Open:", err.Error()}, "")
		logs.Errorf(errMsg)
		tg_bot.TgErrMsg(errMsg)
		return "", "", &global.ErrorMsg{ErrCode: global.ErrRetry.ErrCode, ErrMsg: errMsg}
	}
	start := time.Now()
	part_oss, err := s.bucket.UploadPart(*s.imur, part, header.Size, len(s.parts)+1, oss.Routines(config.Config.Oss.Aliyun.Routines))
	if err != nil {
		_ = part.Close()
		errMsg := strings.Join([]string{info.Uuid(), " ", info.SrcName(), "[", info.Md5(), "] ", info.YunName(), "\n", "UploadPart:", err.Error()}, "")
		logs.Errorf(errMsg)
		tg_bot.TgErrMsg(errMsg)
		return "", "", &global.ErrorMsg{ErrCode: global.ErrRetry.ErrCode, ErrMsg: errMsg}
	}
	_ = part.Close()
	s.parts = append(s.parts, part_oss)
	logs.Warnf("%v %v[%v] %v elapsed:%v", info.Uuid(), info.SrcName(), info.Md5(), info.YunName(), time.Since(start))
	switch info.Last(false, header.Size) {
	case true:
		_, err := s.bucket.CompleteMultipartUpload(*s.imur, s.parts)
		if err != nil {
			errMsg := strings.Join([]string{info.Uuid(), " ", info.SrcName(), "[", info.Md5(), "] ", info.YunName(), "\n", "CompleteMultipartUpload:", err.Error()}, "")
			logs.Errorf(errMsg)
			tg_bot.TgErrMsg(errMsg)
			s.reset()
			return "", "", &global.ErrorMsg{ErrCode: global.ErrFatal.ErrCode, ErrMsg: errMsg}
		}
		yunPath = s.yunPath
		s.reset()
	default:
		return "", "", nil
	}
	return strings.Join([]string{config.Config.Oss.Aliyun.BucketUrl, "/", yunPath}, ""), yunPath, nil
}

func (s *Aliyun) uploadFromFile(info global.FileInfo, header *multipart.FileHeader) (string, string, *global.ErrorMsg) {
	yunPath := ""
	f := config.Config.Upload.Dir + info.DstName()
	fd, err := os.OpenFile(f, os.O_RDONLY, 0)
	if err != nil {
		errMsg := strings.Join([]string{info.Uuid(), " ", info.SrcName(), "[", info.Md5(), "] ", info.YunName(), "\n", "OpenFile:", err.Error()}, "")
		logs.Errorf(errMsg)
		tg_bot.TgErrMsg(errMsg)
		return "", "", &global.ErrorMsg{ErrCode: global.ErrRetry.ErrCode, ErrMsg: errMsg}
	}
	// _, err = fd.Seek(info.Now(false), io.SeekStart)
	_, err = fd.Seek(header.Size, io.SeekEnd)
	if err != nil {
		_ = fd.Close()
		errMsg := strings.Join([]string{info.Uuid(), " ", info.SrcName(), "[", info.Md5(), "] ", info.YunName(), "\n", "Seek:", err.Error()}, "")
		logs.Errorf(errMsg)
		tg_bot.TgErrMsg(errMsg)
		return "", "", &global.ErrorMsg{ErrCode: global.ErrRetry.ErrCode, ErrMsg: errMsg}
	}
	start := time.Now()
	// part_oss, err := s.bucket.UploadPartFromFile(*s.imur, f, info.Now(false), header.Size, len(s.parts)+1, oss.Routines(config.Config.Aliyun_Routines))
	part_oss, err := s.bucket.UploadPart(*s.imur, fd, header.Size, len(s.parts)+1, oss.Routines(config.Config.Oss.Aliyun.Routines))
	if err != nil {
		_ = fd.Close()
		errMsg := strings.Join([]string{info.Uuid(), " ", info.SrcName(), "[", info.Md5(), "] ", info.YunName(), "\n", "UploadPart:", err.Error()}, "")
		logs.Errorf(errMsg)
		tg_bot.TgErrMsg(errMsg)
		return "", "", &global.ErrorMsg{ErrCode: global.ErrRetry.ErrCode, ErrMsg: errMsg}
	}
	_ = fd.Close()
	s.parts = append(s.parts, part_oss)
	logs.Warnf("%v %v[%v] %v elapsed:%v", info.Uuid(), info.SrcName(), info.Md5(), info.YunName(), time.Since(start))
	switch info.Last(false, header.Size) {
	case true:
		_, err := s.bucket.CompleteMultipartUpload(*s.imur, s.parts)
		if err != nil {
			errMsg := strings.Join([]string{info.Uuid(), " ", info.SrcName(), "[", info.Md5(), "] ", info.YunName(), "\n", "CompleteMultipartUpload:", err.Error()}, "")
			logs.Errorf(errMsg)
			tg_bot.TgErrMsg(errMsg)
			s.reset()
			return "", "", &global.ErrorMsg{ErrCode: global.ErrFatal.ErrCode, ErrMsg: errMsg}
		}
		yunPath = s.yunPath
		s.reset()
	default:
		return "", "", nil
	}
	return strings.Join([]string{config.Config.Oss.Aliyun.BucketUrl, "/", yunPath}, ""), yunPath, nil
}

func (s *Aliyun) reset() {
	s.bucket = nil
	s.imur = nil
	s.parts = nil
	s.yunPath = ""
}

func (s *Aliyun) Put() {
	s.reset()
	aliyums.Put(s)
}

func NewBucket() (*oss.Bucket, error) {
	client, err := oss.New(config.Config.Oss.Aliyun.EndPoint,
		config.Config.Oss.Aliyun.AccessKeyId,
		config.Config.Oss.Aliyun.AccessKeySecret, oss.Timeout(120000, 120000))
	if err != nil {
		return nil, err
	}
	bucket, err := client.Bucket(config.Config.Oss.Aliyun.BucketName)
	if err != nil {
		return nil, err
	}
	return bucket, nil
}

func (s *Aliyun) DeleteFile(key string) error {
	err := s.bucket.DeleteObject(key)
	if err != nil {
		logs.Errorf(err.Error())
		return err
	}
	return nil
}
