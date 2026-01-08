package gorutine_fans

import "sync"

const BufferSize = 16

type FanInOut struct {
	inChans  []<-chan int32
	outChans []chan<- int32
}

func NewFanInOut(inChans []<-chan int32, outChans []chan<- int32) *FanInOut {
	return &FanInOut{
		inChans:  inChans,
		outChans: outChans,
	}
}

func (f *FanInOut) Start() {
	var wg sync.WaitGroup

	for _, inCh := range f.inChans {
		wg.Add(1)

		go func() {
			defer wg.Done()

			for k := range inCh {
				for _, outChan := range f.outChans {
					outChan <- k
				}
			}
		}()

	}

	go func() {
		wg.Wait()
		for _, outChan := range f.outChans {
			close(outChan)
		}
	}()

}

func FanOut(inputChan <-chan int32, n int32) []chan<- int32 {
	outChans := make([]chan<- int32, n)
	for i := range outChans {
		outChans[i] = make(chan int32, BufferSize)
	}

	go func() {
		defer func() {
			for _, outChan := range outChans {
				close(outChan)
			}
		}()

		for k := range inputChan {
			for _, outChan := range outChans {
				outChan <- k
			}
		}
	}()

	return outChans
}

func FanIn(inputChans []<-chan int32) chan<- int32 {
	outChan := make(chan int32, BufferSize)
	var wg sync.WaitGroup

	for _, inCh := range inputChans {
		wg.Add(1)

		go func() {
			defer wg.Done()
			for k := range inCh {
				outChan <- k
			}
		}()

	}

	go func() {
		wg.Wait()
		close(outChan)
	}()

	return outChan
}
