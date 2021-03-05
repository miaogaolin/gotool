package mapx

import "sync"

type SyncMap struct {
	sync.Map
}

func (s *SyncMap) AddInt(key string, num int) {
	if val, ok := s.Load(key); ok {
		s.Store(key, val.(int)+num)
		return
	}
	s.Store(key, num)
}

func (s *SyncMap) GetInt(key string) int {
	val, _ := s.Load(key)
	if val != nil {
		return val.(int)
	}
	return 0
}
