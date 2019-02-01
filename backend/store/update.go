package store

import "github.com/ff14wed/sibyl/backend/models"

// Update defines the interface for making modifications to the database
// and emitting the resulting stream and entity events.
// ModifyDB is expected to run in a thread safe environment so it doesn't
// compete with other updates running at the same time. It is assumed that
// the resulting stream events are applied first and in order. Then the
// resulting entity events are applied in order.
type Update interface {
	ModifyDB(db *models.DB) ([]models.StreamEvent, []models.EntityEvent, error)
}
