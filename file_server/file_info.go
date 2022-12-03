package main

import (
	"sync"
	"time"

	"github.com/cwloo/gonet/logs"
)

// <summary>
// FileInfo
// <summary>
type FileInfo struct {
	Uuid     string
	Md5      string
	SrcName  string
	DstName  string
	Now      int64
	Total    int64
	DoneTime int64
	HitTime  int64
	Md5Ok    bool
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

func (s *FileInfo) Update(size int64) {
	if size <= 0 {
		logs.LogFatal("error")
	}
	s.Now += size
	if s.Now > s.Total {
		logs.LogFatal("error")
	}
}

func (s *FileInfo) Finished() bool {
	return s.Now == s.Total
}

func (s *FileInfo) UpdateHitTime(time time.Time) {
	s.HitTime = time.Unix()
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
	n := 0
	s.l.Lock()
	info, ok = s.m[md5]
	if !ok {
		info = &FileInfo{}
		s.m[md5] = info
		n = len(s.m)
		s.l.Unlock()
		goto end
	}
	s.l.Unlock()
	return
end:
	logs.LogError("md5:%v size=%v", md5, n)
	return
}

func (s *FileInfos) Remove(md5 string) (info *FileInfo) {
	n := 0
	s.l.Lock()
	if c, ok := s.m[md5]; ok {
		info = c
		delete(s.m, md5)
		n = len(s.m)
		s.l.Unlock()
		goto end
	}
	s.l.Unlock()
	return
end:
	logs.LogError("md5:%v size=%v", md5, n)
	return
}

func (s *FileInfos) Range(cb func(string, *FileInfo)) {
	s.l.Lock()
	for md5, info := range s.m {
		cb(md5, info)
	}
	s.l.Unlock()
}

func (s *FileInfos) RangeRemove(cond func(*FileInfo) bool, cb func(*FileInfo)) {
	s.l.Lock()
	for md5, info := range s.m {
		if cond(info) {
			cb(info)
			delete(s.m, md5)
		}
	}
	s.l.Unlock()
}
