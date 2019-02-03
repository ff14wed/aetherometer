package store

import "github.com/ff14wed/sibyl/backend/models"

// Streams defines the structure of the data store in the store provider and
// is consumed by updates for modification
type Streams struct {
	Map      map[int]*models.Stream
	KeyOrder []int
}

// Update defines the interface for making modifications to the streams store
// and emitting the resulting stream and entity events.
// ModifyStore is expected to run in a thread safe environment so it doesn't
// compete with other updates running at the same time. It is assumed that
// the resulting stream events are applied first and in order. Then the
// resulting entity events are applied in order.
type Update interface {
	ModifyStore(streams *Streams) ([]models.StreamEvent, []models.EntityEvent, error)
}
