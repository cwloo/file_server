package main

import (
	"sync"

	"github.com/cwloo/gonet/logs"
)

// <summary>
// FileInfo
// <summary>
type FileInfo struct {
	Uuid    string
	Md5     string
	SrcName string
	DstName string
	Now     int64
	Total   int64
}

func (s *FileInfo) Assert() {
	if s.Uuid == "" {
		logs.LogFatal("")
	}
	if s.Md5 == "" {
		logs.LogFatal("")
	}
	if s.SrcName == "" {
		logs.LogFatal("")
	}
	if s.DstName == "" {
		logs.LogFatal("")
	}
	// if s.Now == int64(0) {
	// 	logs.LogFatal("")
	// }
	if s.Total == int64(0) {
		logs.LogFatal("")
	}
}

func (s *FileInfo) Finished() bool {
	return s.Now == s.Total
}

// <summary>
// FileInfos [md5]=FileInfo
// <summary>
type FileInfos struct {
	l *sync.Mutex
	m map[string]*FileInfo
}

func NewFileInfos() *FileInfos {
	return &FileInfos{m: map[string]*FileInfo{}, l: &sync.Mutex{}}
}

func (s *FileInfos) Len() int {
	s.l.Lock()
	c := len(s.m)
	s.l.Unlock()
	return c
}

func (s *FileInfos) Get(md5 string) (info *FileInfo) {
	s.l.Lock()
	if c, ok := s.m[md5]; ok {
		info = c
	}
	s.l.Unlock()
	return
}

func (s *FileInfos) Do(md5 string, cb func(*FileInfo)) {
	var info *FileInfo
	s.l.Lock()
	if c, ok := s.m[md5]; ok {
		info = c
		s.l.Unlock()
		goto end
	}
	s.l.Unlock()
	return
end:
	cb(info)
}

func (s *FileInfos) GetAdd(md5 string) (info *FileInfo, ok bool) {
	s.l.Lock()
	info, ok = s.m[md5]
	if !ok {
		info = &FileInfo{}
		s.m[md5] = info
		s.l.Unlock()
		goto end
	}
	s.l.Unlock()
	return
end:
	logs.LogError("md5:%v", md5)
	return
}

func (s *FileInfos) Remove(md5 string) (info *FileInfo) {
	s.l.Lock()
	if c, ok := s.m[md5]; ok {
		info = c
		delete(s.m, md5)
		s.l.Unlock()
		goto end
	}
	s.l.Unlock()
	return
end:
	logs.LogError("md5:%v", md5)
	return
}

func (s *FileInfos) Range(cb func(string, *FileInfo)) {
	s.l.Lock()
	for md5, info := range s.m {
		cb(md5, info)
	}
	s.l.Unlock()
}
