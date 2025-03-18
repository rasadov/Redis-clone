package main

import (
	"fmt"
	"sync"
	"time"
)

type entry struct {
	value      string
	expiration time.Time
}

type InMemoryStorage struct {
	data map[string]entry
	mu   sync.RWMutex
}

func NewInMemoryStorage() *InMemoryStorage {
	return &InMemoryStorage{
		data: make(map[string]entry),
	}
}

func (s *InMemoryStorage) Get(key string) (string, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	val, ok := s.data[key]
	fmt.Println(val.expiration)
	fmt.Println(time.Now())
	if ok && !val.expiration.IsZero() && val.expiration.Before(time.Now()) {
		delete(s.data, key)
		return "", false
	}
	return val.value, ok
}

func (s *InMemoryStorage) SetKey(key string, value string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.data[key] = entry{
		value: value,
	}
}

func (s *InMemoryStorage) SetKeyWithTTL(key string, value string, ttl time.Duration) {
	s.mu.Lock()
	defer s.mu.Unlock()
	expire := time.Now().Add(ttl)
	fmt.Println("TTL: ", ttl)
	fmt.Println("EXPIRE ", expire)
	s.data[key] = entry{
		value:      value,
		expiration: expire,
	}
}
