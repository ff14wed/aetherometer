package hub

import (
	"fmt"
	"sync"
	"sync/atomic"

	"github.com/ff14wed/aetherometer/core/models"
)

// StreamHub is responsible for broadcasting Stream events to subscribers
// Example usage:
// 	streamHub := hub.NewStreamHub(20)
// 	go func() {
// 		sub, id := streamHub.Subscribe()
// 		defer streamHub.Unsubscribe(id)
// 		for {
// 			select {
// 			case payload := <-sub:
// 				fmt.Printf("%#v\n", payload)
// 			}
// 		}
// 	}
// 	streamHub.Broadcast(models.StreamEvent{StreamID: 1})
//
// Expected output:
// 	models.StreamEvent{StreamID:1, Type:models.StreamEventType(nil)}
type StreamHub struct {
	subscribers map[uint64]chan *models.StreamEvent
	baseSubID   uint64
	chanSize    int

	subLock sync.Mutex
}

// NewStreamHub creates a new hub for broadcasting Stream events to subscribers
func NewStreamHub(chanSize int) *StreamHub {
	return &StreamHub{
		subscribers: make(map[uint64]chan *models.StreamEvent),
		baseSubID:   0,
		chanSize:    chanSize,
	}
}

// Broadcast sends the message to all subscribers of this hub
func (h *StreamHub) Broadcast(payload models.StreamEvent) {
	subsList := []chan *models.StreamEvent{}
	h.subLock.Lock()
	for _, sub := range h.subscribers {
		subsList = append(subsList, sub)
	}
	h.subLock.Unlock()

	for _, sub := range subsList {
		payloadClone := payload
		select {
		case sub <- &payloadClone:
		default:
			fmt.Println("Channel is blocked. Dropping message:", payload)
		}
	}
}

// Subscribe adds a new hub subscriber
func (h *StreamHub) Subscribe() (chan *models.StreamEvent, uint64) {
	sub := make(chan *models.StreamEvent, h.chanSize)
	id := atomic.AddUint64(&h.baseSubID, 1)
	h.subLock.Lock()
	h.subscribers[id] = sub
	h.subLock.Unlock()
	return sub, id
}

// Unsubscribe removes a hub subscriber
func (h *StreamHub) Unsubscribe(id uint64) {
	h.subLock.Lock()
	if sub, ok := h.subscribers[id]; ok {
		close(sub)
		delete(h.subscribers, id)
	}
	h.subLock.Unlock()
}
