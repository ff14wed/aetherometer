package hub

import (
	"sync"
	"time"
)

type EventBatcher struct {
	batchInterval time.Duration
	events        chan struct{}

	lock sync.Mutex
}

func NewEventBatcher(batchInterval time.Duration) *EventBatcher {
	return &EventBatcher{
		batchInterval: batchInterval,
		events:        make(chan struct{}, 10),
	}
}

// Notify posts an event to the EventBatcher. If the EventBatcher is currently
// sleeping, then Notify will wake the EventBatcher and the EventBatcher will
// post a batched event after the batchInterval duration
func (b *EventBatcher) Notify() {
	success := b.lock.TryLock()
	if !success {
		return
	}
	go b.sleepUntilBatchInterval()
}

func (b *EventBatcher) sleepUntilBatchInterval() {
	defer b.lock.Unlock()
	time.Sleep(b.batchInterval)

	// Don't worry if the channel is blocked
	select {
	case b.events <- struct{}{}:
	default:
	}
}

// BatchedEvents allows consumers to receive batched events from a channel
func (b *EventBatcher) BatchedEvents() chan struct{} {
	return b.events
}
