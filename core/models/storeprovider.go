package models

// StoreProvider describes the expected interface of a datastore that can
// provide the backing API requests.
// There is no normalization of the data expected in the store, so each
// stream has its own independent state. Querying for any data requires
// walking down the data hierarchy.
type StoreProvider interface {
	Streams() ([]Stream, error)
	Stream(streamID int) (*Stream, error)
	Entity(streamID int, entityID uint64) (*Entity, error)
	StreamEventSource() StreamEventSource
	EntityEventSource() EntityEventSource
}
