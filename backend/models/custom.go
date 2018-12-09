package models

import (
	"fmt"
)

// DB encompasses the entire internal store for the current state of all
// streams. There is no normalization of the data in the store, so each
// stream has its own independent state. Querying for any data requires
// walking down the data hierarchy.
type DB struct {
	StreamsMap map[int]Stream
	StreamKeys []int
}

// Streams returns all the streams from the internal store.
func (db *DB) Streams() []Stream {
	streams := make([]Stream, len(db.StreamKeys))
	for i, k := range db.StreamKeys {
		streams[i] = db.StreamsMap[k]
	}
	return streams
}

// Stream returns a specific stream from the store, queried by streamID.
func (db *DB) Stream(streamID int) (Stream, error) {
	s, ok := db.StreamsMap[streamID]
	if !ok {
		return Stream{}, fmt.Errorf("stream ID %d not found", streamID)
	}
	return s, nil
}

// Entity returns a specific entity in a specific from the store, queried by
// streamID and entityID. It returns an error if the stream ID is not found or
// if the entityID is not found in the stream.
func (db *DB) Entity(streamID int, entityID int) (Entity, error) {
	s, err := db.Stream(streamID)
	if err != nil {
		return Entity{}, err
	}
	e, err := s.Entity(entityID)
	if err != nil {
		return Entity{}, fmt.Errorf("stream id %d: %s", streamID, err)
	}
	return e, nil
}

// Stream represents state reconstructed from the live stream of data from a
// running FFXIV instance.
type Stream struct {
	Pid               int            `json:"pid"`
	Maps              []MapInfo      `json:"maps"`
	TerritoryID       int            `json:"territoryID"`
	MapID             int            `json:"mapID"`
	TimestampMs       int            `json:"timestampMs"`
	EntitiesMap       map[int]Entity `json:"entities"`
	EntitiesKeys      []int
	MyEntityID        int           `json:"myEntityID"`
	TargetHateRanking []HateRanking `json:"targetHateRanking"`
	NearbyEnemyHate   []HateEntry   `json:"nearbyEnemyHate"`
	CraftingInfo      *CraftingInfo `json:"craftingInfo"`
}

// Entity returns a specific entity from the stream, queried by entityID.
func (s *Stream) Entity(entityID int) (Entity, error) {
	e, ok := s.EntitiesMap[entityID]
	if !ok {
		return Entity{}, fmt.Errorf("entity ID %d not found", entityID)
	}
	return e, nil
}

// Entities returns all the entities from the stream.
func (s *Stream) Entities() []Entity {
	entities := make([]Entity, len(s.EntitiesKeys))
	for i, k := range s.EntitiesKeys {
		entities[i] = s.EntitiesMap[k]
	}
	return entities
}
