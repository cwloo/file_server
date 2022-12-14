package main

import "sync"

var (
	uploaderStates = sync.Pool{
		New: func() any {
			return &uploaderState{}
		},
	}
)

// <summary>
// State
// <summary>
type State interface {
	TryAdd(md5 string)
	SetDone(md5 string)
	AllDone() bool
	Range(cb func(string, bool))
	Remove(md5 string) bool
	Put()
}

// <summary>
// uploaderState
// <summary>
type uploaderState struct {
	m map[string]bool
	l *sync.RWMutex
}

func NewUploaderState() State {
	s := uploaderStates.Get().(*uploaderState)
	s.m = map[string]bool{}
	s.l = &sync.RWMutex{}
	return s
}

func (s *uploaderState) TryAdd(md5 string) {
	s.l.Lock()
	if _, ok := s.m[md5]; !ok {
		s.m[md5] = false
	}
	s.l.Unlock()
}

func (s *uploaderState) SetDone(md5 string) {
	s.l.Lock()
	if _, ok := s.m[md5]; ok {
		s.m[md5] = true
	}
	s.l.Unlock()
}

func (s *uploaderState) AllDone() bool {
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

func (s *uploaderState) reset() {
	s.m = nil
}

func (s *uploaderState) Put() {
	s.reset()
	uploaderStates.Put(s)
}

func (s *uploaderState) Remove(md5 string) (ok bool) {
	s.l.Lock()
	_, ok = s.m[md5]
	if ok {
		delete(s.m, md5)
	}
	s.l.Unlock()
	return
}

func (s *uploaderState) Range(cb func(string, bool)) {
	s.l.RLock()
	for md5, ok := range s.m {
		cb(md5, ok)
	}
	s.l.RUnlock()
}
