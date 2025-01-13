package kmutex

import (
	"hash/fnv"
	"runtime"
	"sync"
)

type KMutexHasher[T comparable] func(key T) uint64

type Kmutex[T comparable] struct {
	mutexes []sync.Mutex
	hasher  KMutexHasher[T]
}

func NewKmutex[T comparable](hasher KMutexHasher[T], n int) *Kmutex[T] {
	if n <= 0 {
		n = runtime.NumCPU()
	}
	return &Kmutex[T]{
		mutexes: make([]sync.Mutex, n),
		hasher:  hasher,
	}
}

func (km *Kmutex[T]) Lock(key T) {
	km.mutexes[km.hasher(key)%uint64(len(km.mutexes))].Lock()
}

func (km *Kmutex[T]) Unlock(key T) {
	km.mutexes[km.hasher(key)%uint64(len(km.mutexes))].Unlock()
}

func StringHasher(key string) uint64 {
	h := fnv.New64a()
	h.Write([]byte(key))
	return h.Sum64()
}
