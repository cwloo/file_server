package global

import (
	"sync"

	"github.com/cwloo/gonet/logs"
)

var Uploaders = NewSessionToHandler()

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

func (s *SessionToHandler) Len() (c int) {
	s.l.Lock()
	c = len(s.m)
	s.l.Unlock()
	return
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

func (s *SessionToHandler) GetAdd(uuid string, async bool, new func(bool, string) Uploader) (handler Uploader, ok bool) {
	n := 0
	s.l.Lock()
	handler, ok = s.m[uuid]
	if !ok {
		if new == nil {
			s.l.Unlock()
			goto ERR
		}
		handler = new(async, uuid)
		s.m[uuid] = handler
		n = len(s.m)
		s.l.Unlock()
		goto OK
	}
	s.l.Unlock()
	return
ERR:
	logs.Fatalf("error")
	return
OK:
	logs.Errorf("%v size=%v", uuid, n)
	return
}

func (s *SessionToHandler) List() {
	s.l.Lock()
	logs.Debugf("---------------------------------------------------------------------------------")
	for uuid := range s.m {
		logs.Errorf("%v", uuid)
	}
	logs.Debugf("---------------------------------------------------------------------------------")
	s.l.Unlock()
}

func (s *SessionToHandler) Remove(uuid string) (handler Uploader) {
	s.List()
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
	logs.Errorf("%v size=%v", uuid, n)
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
	logs.Errorf("%v size=%v", uuid, n)
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
	n := 0
	list := []string{}
	s.l.Lock()
	for uuid, handler := range s.m {
		if cond(handler) {
			cb(handler)
			delete(s.m, uuid)
			n = len(s.m)
			list = append(list, uuid)
		}
	}
	s.l.Unlock()
	if len(list) > 0 {
		logs.Errorf("removed:%v size=%v", len(list), n)
	}
}
