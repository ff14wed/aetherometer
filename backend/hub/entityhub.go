package hub

import (
	"fmt"
	"sync"
	"sync/atomic"

	"github.com/ff14wed/sibyl/backend/models"
)

// EntityHub is responsible for broadcasting Entity events to subscribers
// Example usage:
// 	entityHub := hub.NewEntityHub(20)
// 	go func() {
// 		sub, id := entityHub.Subscribe()
// 		defer entityHub.Unsubscribe(id)
// 		for {
// 			select {
// 			case payload := <-sub:
// 				fmt.Printf("%#v\n", payload)
// 			}
// 		}
// 	}
// 	entityHub.Broadcast(models.EntityEvent{StreamID: 1})
//
// Expected output:
// 	models.EntityEvent{StreamID:1, EntityID:0, Type:models.EntityEventType(nil)}
type EntityHub struct {
	subscribers map[uint64]chan models.EntityEvent
	baseSubID   uint64
	chanSize    int

	subLock sync.Mutex
}

// NewEntityHub creates a new hub for broadcasting Entity events to subscribers
func NewEntityHub(chanSize int) *EntityHub {
	return &EntityHub{
		subscribers: make(map[uint64]chan models.EntityEvent),
		baseSubID:   0,
		chanSize:    chanSize,
	}
}

// Broadcast sends the message to all subscribers of this hub
func (h *EntityHub) Broadcast(payload models.EntityEvent) {
	subsList := []chan models.EntityEvent{}
	h.subLock.Lock()
	for _, sub := range h.subscribers {
		subsList = append(subsList, sub)
	}
	h.subLock.Unlock()

	for _, sub := range subsList {
		select {
		case sub <- payload:
		default:
			fmt.Println("Channel is blocked. Dropping message:", payload)
		}
	}
}

// Subscribe adds a new hub subscriber
func (h *EntityHub) Subscribe() (chan models.EntityEvent, uint64) {
	sub := make(chan models.EntityEvent, h.chanSize)
	id := atomic.AddUint64(&h.baseSubID, 1)
	h.subLock.Lock()
	h.subscribers[id] = sub
	h.subLock.Unlock()
	return sub, id
}

// Unsubscribe removes a hub subscriber
func (h *EntityHub) Unsubscribe(id uint64) {
	h.subLock.Lock()
	if sub, ok := h.subscribers[id]; ok {
		close(sub)
		delete(h.subscribers, id)
	}
	h.subLock.Unlock()
}
