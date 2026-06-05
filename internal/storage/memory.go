package storage

import (
	"context"
	"sync"
)

type InMemoryStorage[K comparable, V any] struct {
	mu    sync.RWMutex
	store map[K]V
}

func NewInMemoryStorage[K comparable, V any]() *InMemoryStorage[K, V] {
	return &InMemoryStorage[K, V]{
		store: make(map[K]V),
	}
}

func (s *InMemoryStorage[K, V]) Put(ctx context.Context, key K, value V) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.store[key] = value
	return nil
}

func (s *InMemoryStorage[K, V]) Get(ctx context.Context, key K) (V, bool, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	val, ok := s.store[key]
	return val, ok, nil
}

func (s *InMemoryStorage[K, V]) Delete(ctx context.Context, key K) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.store, key)
	return nil
}

func (s *InMemoryStorage[K, V]) Scan(ctx context.Context, filter func(K) bool) ([]V, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	var result []V
	for k, v := range s.store {
		if filter(k) {
			result = append(result, v)
		}
	}
	return result, nil
}

func (s *InMemoryStorage[K, V]) Keys(ctx context.Context) ([]K, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	keys := make([]K, 0, len(s.store))
	for k := range s.store {
		keys = append(keys, k)
	}
	return keys, nil
}

func (s *InMemoryStorage[K, V]) Flush(ctx context.Context) error { return nil }
func (s *InMemoryStorage[K, V]) Load(ctx context.Context) error  { return nil }

func (s *InMemoryStorage[K, V]) Clear(ctx context.Context) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.store = make(map[K]V)
	return nil
}
