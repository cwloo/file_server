package main

import (
	"sync"

	"github.com/cwloo/gonet/logs"
)

// <summary>
// SessionToHandler [uuid]=handler
// <summary>
type SessionToHandler struct {
	l *sync.RWMutex
	m map[string]*Uploader
}

func NewSessionToHandler() *SessionToHandler {
	return &SessionToHandler{m: map[string]*Uploader{}, l: &sync.RWMutex{}}
}

func (s *SessionToHandler) Len() int {
	s.l.RLock()
	c := len(s.m)
	s.l.RUnlock()
	return c
}

func (s *SessionToHandler) Get(sessionId string) (handler *Uploader) {
	s.l.RLock()
	if c, ok := s.m[sessionId]; ok {
		handler = c
	}
	s.l.RUnlock()
	return
}

func (s *SessionToHandler) Do(sessionId string, cb func(*Uploader)) {
	var handler *Uploader
	s.l.RLock()
	if c, ok := s.m[sessionId]; ok {
		handler = c
		s.l.RUnlock()
		goto end
	}
	s.l.RUnlock()
	return
end:
	cb(handler)
}

func (s *SessionToHandler) Add(sessionId string, handler *Uploader) (old *Uploader) {
	s.l.Lock()
	if c, ok := s.m[sessionId]; ok {
		old = c
	}
	logs.LogError("uuid:%v", sessionId)
	s.m[sessionId] = handler
	s.l.Unlock()
	return
}

func (s *SessionToHandler) Remove(sessionId string) (handler *Uploader) {
	s.l.Lock()
	if c, ok := s.m[sessionId]; ok {
		handler = c
		delete(s.m, sessionId)
		logs.LogError("uuid:%v", sessionId)
	}
	s.l.Unlock()
	return
}

func (s *SessionToHandler) Range(cb func(string, *Uploader)) {
	s.l.RLock()
	for sessionId, handler := range s.m {
		cb(sessionId, handler)
	}
	s.l.RUnlock()
}
