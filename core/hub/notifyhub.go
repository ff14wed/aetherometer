package hub

import (
	"fmt"
	"sync"
	"sync/atomic"
)

// NotifyHub is responsible for broadcasting generic events to subscribers
// Example usage:
// 	notifyHub := hub.Notifyhub(20)
// 	go func() {
// 		sub, id := notifyHub.Subscribe()
// 		defer notifyHub.Unsubscribe(id)
// 		for {
// 			select {
// 			case <-sub:
// 				fmt.Println("Received event")
// 			}
// 		}
// 	}
// 	notifyHub.Broadcast()
//
// Expected output:
// 	Received event
type NotifyHub struct {
	subscribers map[uint64]chan struct{}
	baseSubID   uint64
	chanSize    int

	subLock sync.Mutex
}

// NewNotifyHub creates a new hub for broadcasting events to subscribers
func NewNotifyHub(chanSize int) *NotifyHub {
	return &NotifyHub{
		subscribers: make(map[uint64]chan struct{}),
		baseSubID:   0,
		chanSize:    chanSize,
	}
}

// Broadcast sends the message to all subscribers of this hub
func (h *NotifyHub) Broadcast() {
	subsList := []chan struct{}{}
	h.subLock.Lock()
	for _, sub := range h.subscribers {
		subsList = append(subsList, sub)
	}
	h.subLock.Unlock()

	for _, sub := range subsList {
		select {
		case sub <- struct{}{}:
		default:
			fmt.Println("Channel is blocked. Dropping message.")
		}
	}
}

// Subscribe adds a new hub subscriber
func (h *NotifyHub) Subscribe() (chan struct{}, uint64) {
	sub := make(chan struct{}, h.chanSize)
	id := atomic.AddUint64(&h.baseSubID, 1)
	h.subLock.Lock()
	h.subscribers[id] = sub
	h.subLock.Unlock()
	return sub, id
}

// Unsubscribe removes a hub subscriber
func (h *NotifyHub) Unsubscribe(id uint64) {
	h.subLock.Lock()
	if sub, ok := h.subscribers[id]; ok {
		close(sub)
		delete(h.subscribers, id)
	}
	h.subLock.Unlock()
}
