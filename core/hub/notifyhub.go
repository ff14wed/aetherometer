package hub

import (
	"fmt"
	"sync"
	"sync/atomic"
)

// NotifyHub is responsible for broadcasting generic events to subscribers
// Example usage:
// 	notifyHub := hub.NewNotifyHub[string](20)
// 	go func() {
// 		sub, id := notifyHub.Subscribe()
// 		defer notifyHub.Unsubscribe(id)
// 		for payload := range sub {
// 			fmt.Printf("Received payload %#v\n", payload)
// 		}
// 	}()
// 	time.Sleep(10 * time.Millisecond)
// 	notifyHub.Broadcast("foobar")
// 	time.Sleep(10 * time.Millisecond)
//
// Expected output:
// 	Received payload "foobar"
type NotifyHub[T interface{}] struct {
	subscribers map[uint64]chan T
	baseSubID   uint64
	chanSize    int

	subLock sync.Mutex
}

// NewNotifyHub creates a new hub for broadcasting events to subscribers
func NewNotifyHub[T interface{}](chanSize int) *NotifyHub[T] {
	return &NotifyHub[T]{
		subscribers: make(map[uint64]chan T),
		baseSubID:   0,
		chanSize:    chanSize,
	}
}

// Broadcast sends the message to all subscribers of this hub
func (h *NotifyHub[T]) Broadcast(payload T) {
	subsList := []chan T{}
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
func (h *NotifyHub[T]) Subscribe() (chan T, uint64) {
	sub := make(chan T, h.chanSize)
	id := atomic.AddUint64(&h.baseSubID, 1)
	h.subLock.Lock()
	h.subscribers[id] = sub
	h.subLock.Unlock()
	return sub, id
}

// Unsubscribe removes a hub subscriber
func (h *NotifyHub[T]) Unsubscribe(id uint64) {
	h.subLock.Lock()
	if sub, ok := h.subscribers[id]; ok {
		close(sub)
		delete(h.subscribers, id)
	}
	h.subLock.Unlock()
}
