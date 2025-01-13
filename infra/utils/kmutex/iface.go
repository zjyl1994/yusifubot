package kmutex

type IKmutex[T comparable] interface {
	Lock(id T)
	Unlock(id T)
}
