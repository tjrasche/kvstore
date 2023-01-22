package store

import (
	"sync"
)

type Store struct {
	m *sync.RWMutex
	v map[string]string
}

func NewStore() Store {
	return Store{m: &sync.RWMutex{}, v: make(map[string]string)}
}

func (s Store) Put(key string, value string) {
	s.m.Lock()
	defer s.m.Unlock()
	s.v[key] = value
}

func (s Store) Get(key string) (value string, ok bool) {
	s.m.RLock()
	defer s.m.RUnlock()
	value, ok = s.v[key]
	return value, ok
}

func (s Store) Delete(key string) {
	s.m.Lock()
	defer s.m.Unlock()
	delete(s.v, key)
}

func (s Store) CountElements() (count uint64) {
	s.m.RLock()
	defer s.m.RUnlock()
	for range s.v {
		count++
	}
	return count
}
