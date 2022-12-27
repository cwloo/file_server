package uploader

import (
	"sync"
)

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
	Len() int
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

func (s *uploaderState) Len() (c int) {
	s.l.RLock()
	c = len(s.m)
	s.l.RUnlock()
	return
}

func (s *uploaderState) exist(md5 string) (ok bool) {
	s.l.RLock()
	_, ok = s.m[md5]
	s.l.RUnlock()
	return
}

func (s *uploaderState) tryAdd(md5 string) {
	s.l.Lock()
	_, ok := s.m[md5]
	switch !ok {
	case true:
		s.m[md5] = false
	}
	s.l.Unlock()
}

func (s *uploaderState) TryAdd(md5 string) {
	switch !s.exist(md5) {
	case true:
		s.tryAdd(md5)
	}
}

func (s *uploaderState) SetDone(md5 string) {
	s.l.Lock()
	_ = s.remove(md5)
	s.l.Unlock()
}

func (s *uploaderState) AllDone() (ok bool) {
	s.l.RLock()
	ok = len(s.m) == 0
	s.l.RUnlock()
	return
}

func (s *uploaderState) reset() {
	s.m = nil
}

func (s *uploaderState) Put() {
	s.reset()
	uploaderStates.Put(s)
}

func (s *uploaderState) remove(md5 string) (ok bool) {
	_, ok = s.m[md5]
	switch ok {
	case true:
		delete(s.m, md5)
	}
	return
}

func (s *uploaderState) Remove(md5 string) (ok bool) {
	s.l.Lock()
	ok = s.remove(md5)
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
