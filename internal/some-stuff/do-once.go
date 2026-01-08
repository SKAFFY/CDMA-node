package some_stuff

import (
	"context"
	"log"
	"sync"
)

type OnceFunc func(ctx context.Context, params ...any) error

type FuncSync struct {
	did   []bool
	n     int32
	mutex sync.Mutex
}

func NewFuncSync() *FuncSync {
	return &FuncSync{
		did:   make([]bool, 0),
		n:     0,
		mutex: sync.Mutex{},
	}
}

func (f *FuncSync) GetDoOnce(v func(ctx context.Context, params ...any) error) OnceFunc {
	f.mutex.Lock()
	defer f.mutex.Unlock()

	funcIndex := len(f.did)
	f.did = append(f.did, false)
	localMutex := sync.Mutex{}

	log.Print("GetDoOnce finished")

	return func(ctx context.Context, params ...any) error {
		localMutex.Lock()
		defer localMutex.Unlock()
		log.Print("GetDoOnce trying")

		if f.did[funcIndex] {
			log.Print("f.did[funcIndex] == true, func ind: ", funcIndex)
			return nil
		}

		f.did[funcIndex] = true

		log.Print("GetDoOnce doing")
		return v(ctx, params...)
	}
}

func (f *FuncSync) Reset() {
	f.mutex.Lock()
	defer f.mutex.Unlock()

	log.Print("resetting funcs")

	for i, _ := range f.did {
		f.did[i] = false
	}
}
