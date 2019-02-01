package models

//go:generate counterfeiter . StoreProvider

// StoreProvider describes the expected interface of a datastore that can
// provide the backing API requests
type StoreProvider interface {
	Streams() ([]Stream, error)
	Stream(streamID int) (Stream, error)
	Entity(streamID int, entityID uint64) (Entity, error)
	StreamEventSource() StreamEventSource
	EntityEventSource() EntityEventSource
}
