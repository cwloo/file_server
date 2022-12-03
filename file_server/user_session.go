package main

import (
	"sync"

	"github.com/cwloo/gonet/logs"
)

// <summary>
// SessionToHandler [uuid]=handler
// <summary>
type SessionToHandler struct {
	l *sync.Mutex
	m map[string]Uploader
}

func NewSessionToHandler() *SessionToHandler {
	return &SessionToHandler{m: map[string]Uploader{}, l: &sync.Mutex{}}
}

func (s *SessionToHandler) Len() int {
	s.l.Lock()
	c := len(s.m)
	s.l.Unlock()
	return c
}

func (s *SessionToHandler) Get(uuid string) (handler Uploader) {
	s.l.Lock()
	if c, ok := s.m[uuid]; ok {
		handler = c
	}
	s.l.Unlock()
	return
}

func (s *SessionToHandler) Do(uuid string, cb func(Uploader)) {
	var handler Uploader
	s.l.Lock()
	if c, ok := s.m[uuid]; ok {
		handler = c
		s.l.Unlock()
		goto OK
	}
	s.l.Unlock()
	return
OK:
	cb(handler)
}

func (s *SessionToHandler) GetAdd(uuid string, async bool) (handler Uploader, ok bool) {
	n := 0
	s.l.Lock()
	handler, ok = s.m[uuid]
	if !ok {
		switch async {
		case true:
			handler = NewAsyncUploader(uuid)
		default:
			handler = NewSyncUploader(uuid)
		}
		s.m[uuid] = handler
		n = len(s.m)
		s.l.Unlock()
		goto OK
	}
	s.l.Unlock()
	return
OK:
	logs.LogError("uuid:%v size=%v", uuid, n)
	return
}

func (s *SessionToHandler) Remove(uuid string) (handler Uploader) {
	n := 0
	s.l.Lock()
	if c, ok := s.m[uuid]; ok {
		handler = c
		delete(s.m, uuid)
		n = len(s.m)
		s.l.Unlock()
		goto OK
	}
	s.l.Unlock()
	return
OK:
	logs.LogError("uuid:%v size=%v", uuid, n)
	return
}

func (s *SessionToHandler) RemoveWithCond(uuid string, cond func(Uploader) bool, cb func(Uploader)) (handler Uploader) {
	n := 0
	s.l.Lock()
	if c, ok := s.m[uuid]; ok {
		if cond(c) {
			handler = c
			cb(handler)
			delete(s.m, uuid)
			n = len(s.m)
			s.l.Unlock()
			goto OK
		}
	}
	s.l.Unlock()
	return
OK:
	logs.LogError("uuid:%v size=%v", uuid, n)
	return
}

func (s *SessionToHandler) Range(cb func(string, Uploader)) {
	s.l.Lock()
	for uuid, handler := range s.m {
		cb(uuid, handler)
	}
	s.l.Unlock()
}

func (s *SessionToHandler) RangeRemoveWithCond(cond func(Uploader) bool, cb func(Uploader)) {
	s.l.Lock()
	for uuid, handler := range s.m {
		if cond(handler) {
			cb(handler)
			delete(s.m, uuid)
		}
	}
	s.l.Unlock()
}
