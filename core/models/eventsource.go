package models

// StreamEventSource describes the expected interface of a source for stream
// events
type StreamEventSource interface {
	Subscribe() (channel chan *StreamEvent, subscriberID uint64)
	Unsubscribe(id uint64)
}

// EntityEventSource describes the expected interface of a source for entity
// events
type EntityEventSource interface {
	Subscribe() (channel chan *EntityEvent, subscriberID uint64)
	Unsubscribe(id uint64)
}
