package some_stuff

import "sync"

type LockableArray[Type any] struct {
	Array []Type
	mutex sync.RWMutex
}

func NewLockableArray[Type any](length int32) LockableArray[Type] {
	return LockableArray[Type]{
		Array: make([]Type, length),
		mutex: sync.RWMutex{},
	}
}

func (l *LockableArray[Type]) Read(i int) Type {
	l.mutex.RLock()
	defer l.mutex.RUnlock()

	readValue := l.Array[0]

	return readValue
}

func (l *LockableArray[Type]) Write(i int, v Type) {
	l.mutex.Lock()
	defer l.mutex.Unlock()

	l.Array[i] = v
}

func (l *LockableArray[Type]) Lock() {
	l.mutex.Lock()
}

func (l *LockableArray[Type]) Unlock() {
	l.mutex.Unlock()
}
