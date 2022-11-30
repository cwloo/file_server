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

func (s *SessionToHandler) Add(sessionId string, handler Uploader) (old Uploader) {
	s.l.Lock()
	if c, ok := s.m[sessionId]; ok {
		old = c
		logs.LogFatal("error")
	}
	s.m[sessionId] = handler
	s.l.Unlock()
	logs.LogError("uuid:%v", sessionId)
	return
}

func (s *SessionToHandler) Remove(sessionId string) (handler Uploader) {
	s.l.Lock()
	if c, ok := s.m[sessionId]; ok {
		handler = c
		delete(s.m, sessionId)
		s.l.Unlock()
		goto end
	}
	s.l.Unlock()
	return
end:
	logs.LogError("uuid:%v", sessionId)
	return
}

func (s *SessionToHandler) Range(cb func(string, Uploader)) {
	s.l.Lock()
	for sessionId, handler := range s.m {
		cb(sessionId, handler)
	}
	s.l.Unlock()
}

func (s *SessionToHandler) CHeckRemove(cond func(Uploader) bool) {
	s.l.Lock()
	for sessionId, handler := range s.m {
		if cond(handler) {
			handler.Clear()
			delete(s.m, sessionId)
		}
	}
	s.l.Unlock()
}
