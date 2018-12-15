package models

//go:generate counterfeiter . StreamEventSource
//go:generate counterfeiter . EntityEventSource

// StreamEventSource describes the expected interface of a source for stream
// events
type StreamEventSource interface {
	Subscribe() (channel chan StreamEventsPayload, subscriberID uint64)
	Unsubscribe(id uint64)
}

// EntityEventSource describes the expected interface of a source for entity
// events
type EntityEventSource interface {
	Subscribe() (channel chan EntityEventsPayload, subscriberID uint64)
	Unsubscribe(id uint64)
}
