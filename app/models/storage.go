package models

import (
	"strings"
	"sync"
	"time"
)

type Entry struct {
	Value      string
	Expiration time.Time
}

type InMemoryStorage struct {
	Index       uint64
	Data        map[string]Entry
	mu          sync.RWMutex
	size        int
	sizeWithTTL int
}

func NewInMemoryStorage(dbIndex uint64) *InMemoryStorage {
	return &InMemoryStorage{
		Data:        make(map[string]Entry),
		size:        0,
		sizeWithTTL: 0,
		Index:       dbIndex,
	}
}

func (s *InMemoryStorage) Get(key string) (string, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	val, ok := s.Data[key]
	if ok && !val.Expiration.IsZero() && val.Expiration.Before(time.Now()) {
		delete(s.Data, key)
		return "", false
	}
	return val.Value, ok
}

func (s *InMemoryStorage) SetKey(key string, value string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.Data[key] = Entry{
		Value: value,
	}
	s.size += 1
}

func (s *InMemoryStorage) SetKeyWithTTL(key string, value string, ttl time.Duration) {
	s.mu.Lock()
	defer s.mu.Unlock()
	expire := time.Now().Add(ttl)
	s.Data[key] = Entry{
		Value:      value,
		Expiration: expire,
	}
	s.size += 1
	s.sizeWithTTL += 1
}

func (s *InMemoryStorage) Keys(pattern string) []string {
	var res []string
	patterns := strings.Split(pattern, "*")
	prefix := patterns[0]
	suffix := patterns[len(patterns)-1]

	if prefix == "*" {
		prefix = ""
	}
	if suffix == "*" {
		suffix = ""
	}

	for key, _ := range s.Data {
		res = append(res, key)
	}
	return res
}
