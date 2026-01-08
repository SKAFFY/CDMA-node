package cdma_node

import (
	"log"
)

type Handler interface {
	Handle(data []int64) error
}

type Node struct {
	sendChan    chan<- []int64
	receiveChan <-chan []int64

	codeIdentifier []int64

	onReceiveHandler Handler
}

func NewNode(
	receiveChan <-chan []int64,
	sendChan chan<- []int64,
	onReceiveHandler Handler,
	codeIdentifier []int64,
) *Node {
	return &Node{
		receiveChan:      receiveChan,
		sendChan:         sendChan,
		onReceiveHandler: onReceiveHandler,
		codeIdentifier:   codeIdentifier,
	}
}

func (n *Node) Run() {

	go func() {
		for k := range n.receiveChan {
			log.Print("NODE RECEIVED")

			decodedData, err := n.decode(k)
			if err != nil {
				log.Printf("error decoding raw data in Node %d, err: %s", n.codeIdentifier, err)
			}

			log.Print("Node Receiving")
			err = n.onReceiveHandler.Handle(decodedData)
			if err != nil {
				log.Printf("error handling data in Node %d, err: %s", n.codeIdentifier, err)
			}
		}
	}()

}

func (n *Node) Send(data []int64) error {
	n.sendChan <- n.encode(data)

	return nil
}

//func (n *Node) decode(rawData []int64) ([]int64, error) {
//	result := make([]int64, len(rawData)/len(n.codeIdentifier))
//
//	for k := range result {
//		var sum int64
//		sum = 0
//
//		for i := range n.codeIdentifier {
//			rawIndex := k*len(n.codeIdentifier) + i
//
//			sum += rawData[rawIndex] * n.codeIdentifier[i]
//		}
//
//		if sum > 0 {
//			sum = 1
//		}
//		if sum < 0 {
//			sum = 0
//		}
//
//		result[k] = sum
//
//		sum = 0
//	}
//
//	return result, nil
//}

func (n *Node) decode(rawData []int64) ([]int64, error) {
	result := make([]int64, len(rawData)/len(n.codeIdentifier))

	for k := range result {
		var sum int64
		for i := range n.codeIdentifier {
			rawIndex := k*len(n.codeIdentifier) + i
			bipolar := 2*rawData[rawIndex] - 1
			sum += bipolar * n.codeIdentifier[i]
		}

		if sum > 0 {
			result[k] = 1
		} else {
			result[k] = 0
		}
	}

	return result, nil
}

//func (n *Node) encode(data []int64) []int64 {
//	result := make([]int64, 0, len(n.codeIdentifier)*len(data))
//
//	for k := range n.codeIdentifier {
//		piece := make([]int64, len(data))
//
//		for i := range data {
//			piece[i] = (data[i] - 1) * n.codeIdentifier[k]
//		}
//
//		result = append(result, piece...)
//	}
//
//	return result
//}

func (n *Node) encode(data []int64) []int64 {
	result := make([]int64, 0, len(n.codeIdentifier)*len(data))

	for i := range data {
		bipolarData := 2*data[i] - 1

		for _, code := range n.codeIdentifier {
			result = append(result, bipolarData*code)
		}
	}

	return result
}
