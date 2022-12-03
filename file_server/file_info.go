package main

import (
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/cwloo/gonet/logs"
	"github.com/cwloo/gonet/utils"
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
	doneTime time.Time
	hitTime  time.Time
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

func (s *FileInfo) Ok() bool {
	return s.Now == s.Total
}

func (s *FileInfo) DoneTime() time.Time {
	return s.doneTime
}

func (s *FileInfo) UpdateDoneTime(time time.Time) {
	s.doneTime = time
}

func (s *FileInfo) HitTime() time.Time {
	return s.hitTime
}

func (s *FileInfo) UpdateHitTime(time time.Time) {
	s.hitTime = time
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
		goto OK
	}
	s.l.Unlock()
	return
OK:
	cb(info)
}

func (s *FileInfos) GetAdd(md5 string, uuid, Filename, total string) (info *FileInfo, ok bool) {
	n := 0
	s.l.Lock()
	info, ok = s.m[md5]
	if !ok {
		size, _ := strconv.ParseInt(total, 10, 0)
		info = &FileInfo{
			Uuid:    uuid,
			Md5:     md5,
			SrcName: Filename,
			DstName: strings.Join([]string{uuid, ".", utils.RandomCharString(10), ".", Filename}, ""),
			Total:   size,
		}
		s.m[md5] = info
		n = len(s.m)
		s.l.Unlock()
		goto OK
	}
	s.l.Unlock()
	return
OK:
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
		goto OK
	}
	s.l.Unlock()
	return
OK:
	logs.LogError("md5:%v size=%v", md5, n)
	return
}

func (s *FileInfos) RemoveWithCond(md5 string, cond func(*FileInfo) bool, cb func(*FileInfo)) (info *FileInfo) {
	n := 0
	s.l.Lock()
	if c, ok := s.m[md5]; ok {
		if cond(c) {
			info = c
			cb(info)
			delete(s.m, md5)
			n = len(s.m)
			s.l.Unlock()
			goto OK
		}
	}
	s.l.Unlock()
	return
OK:
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

func (s *FileInfos) RangeRemoveWithCond(cond func(*FileInfo) bool, cb func(*FileInfo)) {
	s.l.Lock()
	for md5, info := range s.m {
		if cond(info) {
			cb(info)
			delete(s.m, md5)
		}
	}
	s.l.Unlock()
}
