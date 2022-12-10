package main

import (
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/cwloo/gonet/logs"
)

// <summary>
// FileInfo
// <summary>
type FileInfo interface {
	Uuid() string
	Md5() string
	Now() int64
	Total() int64
	SrcName() string
	DstName() string
	YunName() string
	Date() string
	Assert()
	Update(size int64, cb_seg func(FileInfo, OSS, bool) (string, error), cb func(FileInfo) (time.Time, bool)) (done, ok bool, url string, start time.Time)
	Done() bool
	Ok() (bool, string)
	Url() string
	Time() time.Time
	HitTime() time.Time
	UpdateHitTime(time time.Time)
}

// <summary>
// Fileinfo
// <summary>
type Fileinfo struct {
	uuid    string
	md5     string
	srcName string
	dstName string
	yunName string
	now     int64
	total   int64
	url     string
	date    string
	time    time.Time
	hitTime time.Time
	l       *sync.RWMutex
	oss     OSS
}

func NewFileInfo(uuid, md5, Filename string, total int64) FileInfo {
	now := time.Now()
	YMD := now.Format("2006-01-02")
	YMDHMS := now.Format("20060102150405")
	ext := filepath.Ext(Filename)
	suffix := strings.TrimSuffix(Filename, ext)
	yunName := strings.Join([]string{suffix, "-", YMDHMS, ext}, "")
	dstName := strings.Join([]string{md5, "_", YMDHMS, ext}, "")
	s := &Fileinfo{
		uuid:    uuid,
		md5:     md5,
		date:    YMD,
		srcName: Filename,
		dstName: dstName,
		yunName: yunName,
		total:   total,
		l:       &sync.RWMutex{},
	}
	s.Assert()
	return s
}

func (s *Fileinfo) Uuid() string {
	return s.uuid
}

func (s *Fileinfo) Md5() string {
	return s.md5
}

func (s *Fileinfo) Now() int64 {
	return s.now
}

func (s *Fileinfo) Total() int64 {
	return s.total
}

func (s *Fileinfo) SrcName() string {
	return s.srcName
}

func (s *Fileinfo) DstName() string {
	return s.dstName
}

func (s *Fileinfo) YunName() string {
	return s.yunName
}

func (s *Fileinfo) Date() string {
	return s.date
}

func (s *Fileinfo) Assert() {
	if s.uuid == "" {
		logs.LogFatal("")
	}
	if s.md5 == "" {
		logs.LogFatal("")
	}
	if s.srcName == "" {
		logs.LogFatal("")
	}
	if s.dstName == "" {
		logs.LogFatal("")
	}
	if s.yunName == "" {
		logs.LogFatal("")
	}
	// if s.now == int64(0) {
	// 	logs.LogFatal("")
	// }
	if s.total == int64(0) {
		logs.LogFatal("")
	}
	if s.date == "" {
		logs.LogFatal("")
	}
}

func (s *Fileinfo) Update(size int64, cb_seg func(FileInfo, OSS, bool) (string, error), cb func(FileInfo) (time.Time, bool)) (done, ok bool, url string, start time.Time) {
	if size <= 0 {
		logs.LogFatal("error")
	}
	s.l.Lock()
	if s.now == 0 {
		s.oss = NewOss()
	}
	s.now += size
	if s.now > s.total {
		logs.LogFatal("error")
	}
	done = s.now == s.total
	url, _ = cb_seg(s, s.oss, done)
	if done {
		s.oss = nil
		start, ok = cb(s)
		if ok {
			now := time.Now()
			s.time = now
			s.hitTime = now
			s.url = url
		}
	}
	s.l.Unlock()
	return
}

func (s *Fileinfo) Done() bool {
	s.l.RLock()
	done := s.now == s.total
	s.l.RUnlock()
	return done
}

func (s *Fileinfo) Ok() (bool, string) {
	s.l.RLock()
	ok := s.time.Unix() > 0
	url := s.url
	if ok {
		if s.now != s.total {
			logs.LogFatal("error")
		}
	}
	s.l.RUnlock()
	return ok, url
}

func (s *Fileinfo) Url() string {
	s.l.RLock()
	url := s.url
	s.l.RUnlock()
	return url
}

func (s *Fileinfo) Time() time.Time {
	s.l.RLock()
	t := s.time
	s.l.RUnlock()
	return t
}

func (s *Fileinfo) HitTime() time.Time {
	s.l.RLock()
	t := s.hitTime
	s.l.RUnlock()
	return t
}

func (s *Fileinfo) UpdateHitTime(time time.Time) {
	s.l.Lock()
	s.hitTime = time
	s.l.Unlock()
}

// <summary>
// FileInfos [md5]=FileInfo
// <summary>
type FileInfos struct {
	l *sync.Mutex
	m map[string]FileInfo
}

func NewFileInfos() *FileInfos {
	return &FileInfos{m: map[string]FileInfo{}, l: &sync.Mutex{}}
}

func (s *FileInfos) Len() int {
	s.l.Lock()
	c := len(s.m)
	s.l.Unlock()
	return c
}

func (s *FileInfos) Get(md5 string) (info FileInfo) {
	s.l.Lock()
	if c, ok := s.m[md5]; ok {
		info = c
	}
	s.l.Unlock()
	return
}

func (s *FileInfos) Do(md5 string, cb func(FileInfo)) {
	var info FileInfo
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

func (s *FileInfos) GetAdd(md5 string, uuid, Filename, total string) (info FileInfo, ok bool) {
	n := 0
	s.l.Lock()
	info, ok = s.m[md5]
	if !ok {
		size, _ := strconv.ParseInt(total, 10, 0)
		info = NewFileInfo(uuid, md5, Filename, size)
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

func (s *FileInfos) Remove(md5 string) (info FileInfo) {
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

func (s *FileInfos) RemoveWithCond(md5 string, cond func(FileInfo) bool, cb func(FileInfo)) (info FileInfo) {
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

func (s *FileInfos) Range(cb func(string, FileInfo)) {
	s.l.Lock()
	for md5, info := range s.m {
		cb(md5, info)
	}
	s.l.Unlock()
}

func (s *FileInfos) RangeRemoveWithCond(cond func(FileInfo) bool, cb func(FileInfo)) {
	n := 0
	list := []string{}
	s.l.Lock()
	for md5, info := range s.m {
		if cond(info) {
			cb(info)
			delete(s.m, md5)
			n = len(s.m)
			list = append(list, md5)
		}
	}
	s.l.Unlock()
	if len(list) > 0 {
		logs.LogError("removed:%v size=%v", len(list), n)
	}
}
