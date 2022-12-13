package main

import "sync"

var (
	uploaderDatas = sync.Pool{
		New: func() any {
			return &uploaderData{}
		},
	}
)

// <summary>
// Data
// <summary>
type Data interface {
	TryAdd(md5 string)
	SetDone(md5 string)
	AllDone() bool
	Range(cb func(string, bool))
	Remove(md5 string) bool
	Put()
}

// <summary>
// uploaderData
// <summary>
type uploaderData struct {
	m map[string]bool
	l *sync.RWMutex
}

func NewUploaderData() Data {
	s := uploaderDatas.Get().(*uploaderData)
	s.m = map[string]bool{}
	s.l = &sync.RWMutex{}
	return s
}

func (s *uploaderData) TryAdd(md5 string) {
	s.l.Lock()
	if _, ok := s.m[md5]; !ok {
		s.m[md5] = false
	}
	s.l.Unlock()
}

func (s *uploaderData) SetDone(md5 string) {
	s.l.Lock()
	if _, ok := s.m[md5]; ok {
		s.m[md5] = true
	}
	s.l.Unlock()
}

func (s *uploaderData) AllDone() bool {
	s.l.RLock()
	for _, v := range s.m {
		if !v {
			s.l.RUnlock()
			return false
		}
	}
	s.l.RUnlock()
	return true
}

func (s *uploaderData) reset() {
	s.m = nil
}

func (s *uploaderData) Put() {
	s.reset()
	uploaderDatas.Put(s)
}

func (s *uploaderData) Remove(md5 string) (ok bool) {
	s.l.Lock()
	_, ok = s.m[md5]
	if ok {
		delete(s.m, md5)
	}
	s.l.Unlock()
	return
}

func (s *uploaderData) Range(cb func(string, bool)) {
	s.l.RLock()
	for md5, ok := range s.m {
		cb(md5, ok)
	}
	s.l.RUnlock()
}
