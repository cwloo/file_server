package global

import (
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/cwloo/gonet/logs"
)

var FileInfos = NewMd5ToFileInfo()

var (
	fileinfos = sync.Pool{
		New: func() any {
			return &Fileinfo{}
		},
	}
)

type SegmentCallback func(FileInfo, Oss) (string, *ErrorMsg)
type CheckCallback func(FileInfo) (time.Time, bool)

// <summary>
// FileInfo
// <summary>
type FileInfo interface {
	Uuid() string
	Md5() string
	SrcName() string
	DstName() string
	YunName() string
	Date() string
	DateTime() time.Time
	Assert()
	Update(int64, string, SegmentCallback, CheckCallback) (done, ok bool, url string, err *ErrorMsg, start time.Time)
	Now(lock bool) int64
	Total(lock bool) int64
	Last(lock bool, size int64) bool
	Done(lock bool) bool
	Ok(lock bool) (bool, string)
	Url(lock bool) string
	Time(lock bool) time.Time
	HitTime(lock bool) time.Time
	UpdateHitTime(time time.Time)
	Put()
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
	create  time.Time
	time    time.Time
	hitTime time.Time
	l       *sync.RWMutex
	oss     Oss
	new     NewOss
	cancel  bool
}

func NewFileInfo(uuid, md5, Filename string, total int64, useOriginFilename bool, new NewOss) FileInfo {
	now := time.Now()
	YMD := now.Format("2006-01-02")
	YMDHMS := now.Format("20060102150405")
	ext := filepath.Ext(Filename)
	dstName := strings.Join([]string{md5, "_", YMDHMS, ext}, "")
	s := fileinfos.Get().(*Fileinfo)
	s.uuid = uuid
	s.md5 = md5
	s.date = YMD
	s.create = now
	s.srcName = Filename
	s.dstName = dstName
	switch useOriginFilename {
	case true:
		suffix := strings.TrimSuffix(Filename, ext)
		yunName := strings.Join([]string{suffix, "-", YMDHMS, ext}, "")
		s.yunName = yunName
	default:
		s.yunName = dstName
	}
	s.now = 0
	s.total = total
	s.l = &sync.RWMutex{}
	s.new = new
	s.assert()
	return s
}

func (s *Fileinfo) assert() {
	if s.uuid == "" {
		logs.Fatalf("error")
	}
	if s.md5 == "" {
		logs.Fatalf("error")
	}
	if s.srcName == "" {
		logs.Fatalf("error")
	}
	if s.now > int64(0) {
		logs.Fatalf("error")
	}
	if s.total == int64(0) {
		logs.Fatalf("error")
	}
	// if s.url != "" {
	// 	logs.Fatalf("error")
	// }
	// if s.time.Unix() > 0 {
	// 	logs.Fatalf("error")
	// }
	// if s.hitTime.Unix() > 0 {
	// 	logs.Fatalf("error")
	// }
	s.assertNew()
}

func (s *Fileinfo) assertNew() {
	if s.new == nil {
		logs.Fatalf("error")
	}
}

func (s *Fileinfo) assertOss() {
	if s.oss == nil {
		logs.Fatalf("error")
	}
}

func (s *Fileinfo) reset() {
	s.l.Lock()
	s.resetOss(false)
	s.cancel = true
	s.l.Unlock()
}

func (s *Fileinfo) Put() {
	s.reset()
	fileinfos.Put(s)
}

func (s *Fileinfo) resetOss(lock bool) {
	switch lock {
	case true:
		s.l.Lock()
		switch s.oss {
		case nil:
		default:
			s.oss.Put()
			s.oss = nil
		}
		s.l.Unlock()
	default:
		switch s.oss {
		case nil:
		default:
			s.oss.Put()
			s.oss = nil
		}
	}
}

func (s *Fileinfo) Uuid() string {
	return s.uuid
}

func (s *Fileinfo) Md5() string {
	return s.md5
}

func (s *Fileinfo) Now(lock bool) (now int64) {
	switch lock {
	case true:
		s.l.RLock()
		now = s.now
		s.l.RUnlock()
	default:
		now = s.now
	}
	return
}

func (s *Fileinfo) Total(lock bool) (total int64) {
	switch lock {
	case true:
		s.l.RLock()
		total = s.total
		s.l.RUnlock()
	default:
		total = s.total
	}
	return
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

func (s *Fileinfo) DateTime() time.Time {
	return s.create
}

func (s *Fileinfo) Assert() {
	if s.uuid == "" {
		logs.Fatalf("error")
	}
	if s.md5 == "" {
		logs.Fatalf("error")
	}
	if s.srcName == "" {
		logs.Fatalf("error")
	}
	if s.dstName == "" {
		logs.Fatalf("error")
	}
	if s.yunName == "" {
		logs.Fatalf("error")
	}
	// if s.now == int64(0) {
	// 	logs.Fatalf("error")
	// }
	if s.total == int64(0) {
		logs.Fatalf("error")
	}
	if s.date == "" {
		logs.Fatalf("error")
	}
}

func (s *Fileinfo) Update(size int64, ossType string, onSeg SegmentCallback, onCheck CheckCallback) (done, ok bool, url string, err *ErrorMsg, start time.Time) {
	if size <= 0 {
		logs.Fatalf("error")
	}
	s.l.Lock()
	switch s.cancel {
	case true:
		errMsg := strings.Join([]string{s.uuid, " ", s.srcName, "[", s.md5, "] ", s.yunName, "\n", "Cancel"}, "")
		err = &ErrorMsg{ErrCode: ErrCancel.ErrCode, ErrMsg: errMsg}
	default:
		if s.now == 0 {
			s.assertNew()
			s.oss = s.new(s, ossType)
			s.assertOss()
		}
		if s.now+size > s.total {
			s.l.Unlock()
			goto ERR
		}
		url, err = onSeg(s, s.oss)
		switch err {
		case nil:
			s.now += size
			done = s.now == s.total
			if done {
				s.resetOss(false)
				start, ok = onCheck(s)
				if ok {
					now := time.Now()
					s.time = now
					s.hitTime = now
					s.url = url
				}
			}
		default:
			switch err.ErrCode {
			case ErrFatal.ErrCode:
				s.resetOss(false)
			}
		}
	}
	s.l.Unlock()
	return
ERR:
	logs.Fatalf("error")
	return
}

func (s *Fileinfo) Last(lock bool, size int64) (ok bool) {
	switch lock {
	case true:
		s.l.RLock()
		if s.now+size > s.total {
			s.l.RUnlock()
			goto ERR
		}
		ok = s.now+size == s.total
		s.l.RUnlock()
	default:
		if s.now+size > s.total {
			goto ERR
		}
		ok = s.now+size == s.total
	}
	return
ERR:
	logs.Fatalf("error")
	return
}

func (s *Fileinfo) Done(lock bool) (done bool) {
	switch lock {
	case true:
		s.l.RLock()
		done = s.now == s.total
		if done {
			if s.now == 0 {
				s.l.RUnlock()
				goto ERR
			}
		}
		s.l.RUnlock()
	default:
		done = s.now == s.total
		if done {
			if s.now == 0 {
				goto ERR
			}
		}
	}
	return
ERR:
	logs.Fatalf("error")
	return
}

func (s *Fileinfo) Ok(lock bool) (ok bool, url string) {
	switch lock {
	case true:
		s.l.RLock()
		ok = s.time.Unix() > 0
		url = s.url
		if ok {
			if s.now != s.total {
				s.l.RUnlock()
				goto ERR
			}
		}
		s.l.RUnlock()
	default:
		ok = s.time.Unix() > 0
		url = s.url
		if ok {
			if s.now != s.total {
				goto ERR
			}
		}
	}
	return
ERR:
	logs.Fatalf("error")
	return
}

func (s *Fileinfo) Url(lock bool) (url string) {
	switch lock {
	case true:
		s.l.RLock()
		url = s.url
		s.l.RUnlock()
	default:
		url = s.url
	}
	return
}

func (s *Fileinfo) Time(lock bool) (t time.Time) {
	switch lock {
	case true:
		s.l.RLock()
		t = s.time
		s.l.RUnlock()
	default:
		t = s.time
	}
	return
}

func (s *Fileinfo) HitTime(lock bool) (t time.Time) {
	switch lock {
	case true:
		s.l.RLock()
		t = s.hitTime
		s.l.RUnlock()
	default:
		t = s.hitTime
	}
	return
}

func (s *Fileinfo) UpdateHitTime(time time.Time) {
	s.l.Lock()
	s.hitTime = time
	s.l.Unlock()
}

// <summary>
// Md5ToFileInfo [md5]=FileInfo
// <summary>
type Md5ToFileInfo struct {
	l *sync.RWMutex
	m map[string]FileInfo
}

func NewMd5ToFileInfo() *Md5ToFileInfo {
	return &Md5ToFileInfo{m: map[string]FileInfo{}, l: &sync.RWMutex{}}
}

func (s *Md5ToFileInfo) Len() (c int) {
	s.l.RLock()
	c = len(s.m)
	s.l.RUnlock()
	return
}

func (s *Md5ToFileInfo) Get(md5 string) (info FileInfo) {
	s.l.RLock()
	if c, ok := s.m[md5]; ok {
		info = c
	}
	s.l.RUnlock()
	return
}

func (s *Md5ToFileInfo) Do(md5 string, cb func(FileInfo)) {
	var info FileInfo
	s.l.RLock()
	if c, ok := s.m[md5]; ok {
		info = c
		s.l.RUnlock()
		goto OK
	}
	s.l.RUnlock()
	return
OK:
	cb(info)
}

func (s *Md5ToFileInfo) GetAdd(md5 string, uuid, Filename, total string, useOriginFilename bool, new NewOss) (info FileInfo, ok bool) {
	n := 0
	s.l.Lock()
	info, ok = s.m[md5]
	if !ok {
		size, _ := strconv.ParseInt(total, 10, 0)
		info = NewFileInfo(uuid, md5, Filename, size, useOriginFilename, new)
		s.m[md5] = info
		n = len(s.m)
		s.l.Unlock()
		goto OK
	}
	s.l.Unlock()
	return
OK:
	logs.Errorf("md5:%v size=%v", md5, n)
	return
}

func (s *Md5ToFileInfo) Remove(md5 string) (info FileInfo) {
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
	logs.Errorf("md5:%v size=%v", md5, n)
	return
}

func (s *Md5ToFileInfo) RemoveWithCond(md5 string, cond func(FileInfo) bool, cb func(FileInfo)) (info FileInfo) {
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
	logs.Errorf("md5:%v size=%v", md5, n)
	return
}

func (s *Md5ToFileInfo) Range(cb func(string, FileInfo)) {
	s.l.RLock()
	for md5, info := range s.m {
		cb(md5, info)
	}
	s.l.RUnlock()
}

func (s *Md5ToFileInfo) RangeRemoveWithCond(cond func(FileInfo) bool, cb func(FileInfo)) {
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
		logs.Errorf("removed:%v size=%v", len(list), n)
	}
}
