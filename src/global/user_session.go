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
	l *sync.RWMutex
	m map[string]Uploader
}

func NewSessionToHandler() *SessionToHandler {
	return &SessionToHandler{m: map[string]Uploader{}, l: &sync.RWMutex{}}
}

func (s *SessionToHandler) Len() (c int) {
	s.l.RLock()
	c = len(s.m)
	s.l.RUnlock()
	return
}

func (s *SessionToHandler) Get(uuid string) (handler Uploader, ok bool) {
	s.l.RLock()
	handler, ok = s.m[uuid]
	s.l.RUnlock()
	return
}

func (s *SessionToHandler) Do(uuid string, cb func(Uploader)) {
	var handler Uploader
	s.l.RLock()
	c, ok := s.m[uuid]
	switch ok {
	case true:
		handler = c
		s.l.RUnlock()
		goto OK
	}
	s.l.RUnlock()
	return
OK:
	cb(handler)
}

func (s *SessionToHandler) GetAdd(uuid string, async bool, new NewUploader) (handler Uploader, ok bool) {
	handler, ok = s.Get(uuid)
	switch ok {
	case true:
	default:
		handler, ok = s.getAdd(uuid, async, new)
	}
	return
}

func (s *SessionToHandler) getAdd(uuid string, async bool, new NewUploader) (handler Uploader, ok bool) {
	n := 0
	s.l.Lock()
	handler, ok = s.m[uuid]
	switch ok {
	case true:
	default:
		switch new == nil {
		case true:
			s.l.Unlock()
			goto ERR
		}
		handler = new(async, uuid)
		switch handler == nil {
		case true:
			s.l.Unlock()
			goto ERR
		}
		s.m[uuid] = handler
		n = len(s.m)
		s.l.Unlock()
		ok = true
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
	s.l.RLock()
	logs.Debugf("---------------------------------------------------------------------------------")
	for uuid := range s.m {
		logs.Errorf("%v", uuid)
	}
	logs.Debugf("---------------------------------------------------------------------------------")
	s.l.RUnlock()
}

func (s *SessionToHandler) Remove(uuid string) (handler Uploader) {
	s.List()
	handler, _ = s.remove_(uuid)
	return
}

func (s *SessionToHandler) remove_(uuid string) (c Uploader, ok bool) {
	_, ok = s.Get(uuid)
	switch ok {
	case true:
		c, ok = s.remove(uuid)
	default:
	}
	return
}

func (s *SessionToHandler) remove(uuid string) (c Uploader, ok bool) {
	n := 0
	s.l.Lock()
	c, ok = s.m[uuid]
	switch ok {
	case true:
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
	handler, _ = s.removeWithCond_(uuid, cond, cb)
	return
}

func (s *SessionToHandler) removeWithCond_(uuid string, cond func(Uploader) bool, cb func(Uploader)) (c Uploader, ok bool) {
	_, ok = s.Get(uuid)
	switch ok {
	case true:
		c, ok = s.removeWithCond(uuid, cond, cb)
	}
	return
}

func (s *SessionToHandler) removeWithCond(uuid string, cond func(Uploader) bool, cb func(Uploader)) (c Uploader, ok bool) {
	n := 0
	s.l.Lock()
	c, ok = s.m[uuid]
	switch ok {
	case true:
		switch cond(c) {
		case true:
			cb(c)
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
	s.l.RLock()
	for uuid, c := range s.m {
		cb(uuid, c)
	}
	s.l.RUnlock()
}

func (s *SessionToHandler) RangeRemoveWithCond(cond func(Uploader) bool, cb func(Uploader)) {
	n := 0
	list := []string{}
	s.l.Lock()
	for uuid, c := range s.m {
		switch cond(c) {
		case true:
			cb(c)
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
