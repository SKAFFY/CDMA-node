package main

import (
	broadcast_environment "CDMA-telecom-lab1/internal/broadcast-environment"
	cdma_node "CDMA-telecom-lab1/internal/cdma-node"
	some_stuff "CDMA-telecom-lab1/internal/some-stuff"
	"CDMA-telecom-lab1/internal/walsh"
	"fmt"
	"io"
	"log"
	"time"
	"unicode/utf8"

	"simonwaldherr.de/go/golibs/bitmask"
)

const (
	numOfCorrespondents = 4
	bufferSize          = 16
	pkgLength           = 8
	cyclicBufferSize    = 5
)

var (
	messageAlpha = []rune{'A', 'L', 'P', 'H', 'A'}
	messageBeta  = []rune{'B', 'E', 'T', 'A'}
	messageGamma = []rune{'G', 'A', 'M', 'M', 'A'}
	messagePhi   = []rune{'P', 'H', 'I'}
)

type LogWriteHandler struct {
	name string
}

func NewLogWriteHandler(name string) *LogWriteHandler {
	return &LogWriteHandler{
		name: name,
	}
}

func (l *LogWriteHandler) Handle(data []int64) error {
	bm := bitmask.New(0)

	for i := range data {
		if data[i] > 0 {
			bm.Set(i, true)
		}
		if data[i] == 0 {
			bm.Set(i, false)
		}
	}

	r, _ := utf8.DecodeRune(bm.Byte())

	fmt.Printf("name: %s, received data: %q \n", l.name, r)

	return nil
}

type CyclicBufferHandler struct {
	name string

	buffer *some_stuff.CyclicBuffer[rune]
}

func NewCyclicBufferHandler(name string) *CyclicBufferHandler {
	return &CyclicBufferHandler{
		name:   name,
		buffer: some_stuff.NewCyclicBuffer[rune](cyclicBufferSize),
	}
}

func (c *CyclicBufferHandler) Handle(data []int64) error {
	bm := bitmask.New(0)

	for i := range data {
		if data[i] > 0 {
			bm.Set(i, true)
		}
		if data[i] == 0 {
			bm.Set(i, false)
		}
	}

	r, _ := utf8.DecodeRune(bm.Byte())

	if r == '\x00' {
		return nil
	}

	c.buffer.Push(r)

	fmt.Printf("name: %s, received data: %s \n", c.name, string(c.buffer.Body()))

	return nil
}

func main() {
	log.SetOutput(io.Discard)

	inChans := make([]chan []int64, numOfCorrespondents)
	outChans := make([]chan []int64, numOfCorrespondents)
	cdmaNodes := make([]*cdma_node.Node, numOfCorrespondents)

	walshCodes := walsh.NewWalshTable(numOfCorrespondents)

	for i := range numOfCorrespondents {
		inChans[i] = make(chan []int64, bufferSize)
		outChans[i] = make(chan []int64, bufferSize)

		handler := NewCyclicBufferHandler(fmt.Sprint("Node №", i+1))

		cdmaNodes[i] = cdma_node.NewNode(outChans[i], inChans[i], handler, walshCodes[i])
	}

	broadcastInChans := make([]<-chan []int64, numOfCorrespondents)
	broadcastOutChans := make([]chan<- []int64, numOfCorrespondents)

	for i := range numOfCorrespondents {
		broadcastInChans[i] = inChans[i]
		broadcastOutChans[i] = outChans[i]
	}

	broadcastEnvironment := broadcast_environment.NewBroadcastEnvironment(numOfCorrespondents, pkgLength, broadcastInChans, broadcastOutChans)

	broadcastEnvironment.StartBroadcast()

	for i := range numOfCorrespondents {
		cdmaNodes[i].Run()

		var message []rune

		switch i {
		case 0:
			message = messageAlpha
		case 1:
			message = messageBeta
		case 2:
			message = messageGamma
		case 3:
			message = messagePhi
		}

		go func() {
			localNode := cdmaNodes[i]

		outer:
			for {
				for _, messageRune := range message {
					convertedData := convertRuneToData(messageRune)

					log.Print("Node sending")
					err := localNode.Send(convertedData)
					if err != nil {
						log.Printf("Node error: %v", err)
						break outer
					}
				}
			}

			close(inChans[i])
		}()
	}

	time.Sleep(100 * time.Second)

	for i := range numOfCorrespondents {
		close(outChans[i])
		close(inChans[i])
	}

}

func convertRuneToData(r rune) []int64 {
	data := make([]byte, pkgLength)

	utf8.EncodeRune(data, r)

	bm := bitmask.New(int(data[0]))

	convertedData := make([]int64, pkgLength)

	for i := range convertedData {
		if bm.Get(i) {
			convertedData[i] = 1
		} else {
			convertedData[i] = 0
		}
	}

	return convertedData
}
