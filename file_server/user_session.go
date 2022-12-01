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

func (s *SessionToHandler) Get(sessionId string) (handler Uploader) {
	s.l.Lock()
	if c, ok := s.m[sessionId]; ok {
		handler = c
	}
	s.l.Unlock()
	return
}

func (s *SessionToHandler) Do(sessionId string, cb func(Uploader)) {
	var handler Uploader
	s.l.Lock()
	if c, ok := s.m[sessionId]; ok {
		handler = c
		s.l.Unlock()
		goto end
	}
	s.l.Unlock()
	return
end:
	cb(handler)
}

func (s *SessionToHandler) GetAdd(sessionId string, async bool) (handler Uploader, ok bool) {
	n := 0
	s.l.Lock()
	handler, ok = s.m[sessionId]
	if !ok {
		switch async {
		case true:
			handler = NewAsyncUploader(sessionId)
		default:
			handler = NewSyncUploader(sessionId)
		}
		s.m[sessionId] = handler
		n = len(s.m)
		s.l.Unlock()
		goto end
	}
	s.l.Unlock()
	return
end:
	logs.LogError("uuid:%v size=%v", sessionId, n)
	return
}

func (s *SessionToHandler) Remove(sessionId string) (handler Uploader) {
	n := 0
	s.l.Lock()
	if c, ok := s.m[sessionId]; ok {
		handler = c
		delete(s.m, sessionId)
		n = len(s.m)
		s.l.Unlock()
		goto end
	}
	s.l.Unlock()
	return
end:
	logs.LogError("uuid:%v size=%v", sessionId, n)
	return
}

func (s *SessionToHandler) Range(cb func(string, Uploader)) {
	s.l.Lock()
	for sessionId, handler := range s.m {
		cb(sessionId, handler)
	}
	s.l.Unlock()
}

func (s *SessionToHandler) RangeRemove(cond func(Uploader) bool, cb func(Uploader)) {
	s.l.Lock()
	for sessionId, handler := range s.m {
		if cond(handler) {
			cb(handler)
			delete(s.m, sessionId)
		}
	}
	s.l.Unlock()
}
