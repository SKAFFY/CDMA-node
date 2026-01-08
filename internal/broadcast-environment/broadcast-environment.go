package broadcast_environment

import (
	"context"
	"log"
	"sync"
	"time"

	some_stuff "CDMA-telecom-lab1/internal/some-stuff"
)

type BroadcastEnvironment struct {
	ether               []int64
	correspondentNumber int64

	inChans  []<-chan []int64
	outChans []chan<- []int64

	funcSync  *some_stuff.FuncSync
	etherLock sync.Mutex
}

func NewBroadcastEnvironment(n int64, pkgLen int64, inChans []<-chan []int64, outChans []chan<- []int64) *BroadcastEnvironment {
	return &BroadcastEnvironment{
		correspondentNumber: n,
		ether:               make([]int64, some_stuff.NextPowOf2(n)*pkgLen),
		inChans:             inChans,
		outChans:            outChans,
		funcSync:            some_stuff.NewFuncSync(),
		etherLock:           sync.Mutex{},
	}
}

func (b *BroadcastEnvironment) StartBroadcast() {
	var wg sync.WaitGroup

	ticker := time.NewTicker(100 * time.Millisecond)

	go func() {
		for range ticker.C {
			b.etherLock.Lock()

			log.Print("Broadcast Reset")

			b.funcSync.Reset()
			for _, outCh := range b.outChans {
				outCh <- b.ether
			}

			for i := range b.ether {
				b.ether[i] = 0
			}

			b.etherLock.Unlock()
		}
	}()

	for _, inCh := range b.inChans {
		wg.Add(1)

		oncePerTimeFunc := b.funcSync.GetDoOnce(func(_ context.Context, params ...any) error {
			log.Print("here")
			b.etherLock.Lock()
			defer b.etherLock.Unlock()

			value := params[0].([]int64)

			log.Print("here")

			l := min(len(value), len(b.ether))

			for i := range l {
				b.ether[i] += value[i]
			}
			log.Print("there")

			return nil
		})

		log.Print("starting broadcast routines")

		go func() {
			defer wg.Done()

		infinityLoop:
			for {
				select {
				case in, ok := <-inCh:

					log.Print("inch recieved")
					if !ok {
						break infinityLoop
					}

					log.Print("doing once")

					_ = oncePerTimeFunc(nil, in)

					time.Sleep(time.Millisecond)
				}
			}
		}()
	}

	go func() {
		wg.Wait()
		for _, outChan := range b.outChans {
			close(outChan)
		}

		ticker.Stop()

		log.Print("channels closed")
	}()
}
