package utils

import (
	"sync"
)

type safeMap[K comparable, V any] struct {
	lock sync.RWMutex
	data map[K]V
}

func NewSafeMap[K comparable, V any]() *safeMap[K, V] {
	return &safeMap[K, V]{
		lock: sync.RWMutex{},
		data: make(map[K]V),
	}
}

func (sm *safeMap[K, V]) Get(key K) (V, bool) {
	sm.lock.RLock()
	defer sm.lock.RUnlock()
	val, ok := sm.data[key]
	return val, ok
}

func (sm *safeMap[K, V]) Set(key K, val V) {
	sm.lock.Lock()
	defer sm.lock.Unlock()
	sm.data[key] = val
}
