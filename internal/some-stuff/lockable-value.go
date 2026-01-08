package some_stuff

import "sync"

type LockableValue[Type int8 | int | int64 | int32 | int16] struct {
	Value Type
	mutex sync.RWMutex
}

func NewLockableValue[Type int8 | int | int64 | int32 | int16]() LockableValue[Type] {
	return LockableValue[Type]{
		mutex: sync.RWMutex{},
	}
}

func (l *LockableValue[Type]) Read() Type {
	l.mutex.RLock()
	defer l.mutex.RUnlock()

	readValue := l.Value

	return readValue
}

func (l *LockableValue[Type]) Write(v Type) {
	l.mutex.Lock()
	defer l.mutex.Unlock()

	l.Value = v
}

func (l *LockableValue[Type]) Add(v Type) {
	l.mutex.Lock()
	defer l.mutex.Unlock()

	l.Value += v
}

func (l *LockableValue[Type]) Lock() {
	l.mutex.Lock()
}

func (l *LockableValue[Type]) Unlock() {
	l.mutex.Unlock()
}
