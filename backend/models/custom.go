package models

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"strconv"
	"time"

	"github.com/99designs/gqlgen/graphql"
)

//go:generate go run ../scripts/gqlgen.go -v

// DB encompasses the entire internal store for the current state of all
// streams. There is no normalization of the data in the store, so each
// stream has its own independent state. Querying for any data requires
// walking down the data hierarchy.
type DB struct {
	StreamsMap        map[int]*Stream
	StreamKeys        []int
	StreamEventSource StreamEventSource
	EntityEventSource EntityEventSource
}

// Streams returns all the streams from the internal store.
func (db *DB) Streams() []Stream {
	streams := make([]Stream, len(db.StreamKeys))
	for i, k := range db.StreamKeys {
		streams[i] = *db.StreamsMap[k]
	}
	return streams
}

// Stream returns a specific stream from the store, queried by streamID.
func (db *DB) Stream(streamID int) (Stream, error) {
	s, ok := db.StreamsMap[streamID]
	if !ok {
		return Stream{}, fmt.Errorf("stream ID %d not found", streamID)
	}
	return *s, nil
}

// Entity returns a specific entity in a specific from the store, queried by
// streamID and entityID. It returns an error if the stream ID is not found or
// if the entityID is not found in the stream.
func (db *DB) Entity(streamID int, entityID uint64) (Entity, error) {
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

// StreamEvents returns an event channel that can be used for subscriptions to
// Stream events
func (db *DB) StreamEvents(ctx context.Context) (<-chan StreamEvent, error) {
	ch, id := db.StreamEventSource.Subscribe()
	go func() {
		<-ctx.Done()
		db.StreamEventSource.Unsubscribe(id)
	}()
	return ch, nil
}

// EntityEvents returns an event channel that can be used for subscriptions to
// Entity events
func (db *DB) EntityEvents(ctx context.Context) (<-chan EntityEvent, error) {
	ch, id := db.EntityEventSource.Subscribe()
	go func() {
		<-ctx.Done()
		db.EntityEventSource.Unsubscribe(id)
	}()
	return ch, nil
}

// Stream represents state reconstructed from the live stream of data from a
// running FFXIV instance.
type Stream struct {
	Pid        int    `json:"pid"`
	MyEntityID uint64 `json:"myEntityID"`

	Place        Place         `json:"place"`
	Enmity       Enmity        `json:"enmity"`
	CraftingInfo *CraftingInfo `json:"craftingInfo"`

	EntitiesMap  map[uint64]*Entity `json:"entities"`
	EntitiesKeys []uint64
}

// Entity returns a specific entity from the stream, queried by entityID.
func (s *Stream) Entity(entityID uint64) (Entity, error) {
	e, ok := s.EntitiesMap[entityID]
	if !ok {
		return Entity{}, fmt.Errorf("entity ID %d not found", entityID)
	}
	return *e, nil
}

// Entities returns all the entities from the stream.
func (s *Stream) Entities() []Entity {
	entities := make([]Entity, len(s.EntitiesKeys))
	for i, k := range s.EntitiesKeys {
		entities[i] = *s.EntitiesMap[k]
	}
	return entities
}

// MarshalTimestamp converts the provided time to milliseconds since the Unix
// epoch.
func MarshalTimestamp(t time.Time) graphql.Marshaler {
	return graphql.WriterFunc(func(w io.Writer) {
		io.WriteString(w, strconv.FormatInt(getTimeInMs(t), 10))
	})
}

func getTimeInMs(t time.Time) int64 {
	return t.UnixNano() / (int64(time.Millisecond) / int64(time.Nanosecond))
}

// MarshalUint marshals the provided uint64 to a string
func MarshalUint(u uint64) graphql.Marshaler {
	return graphql.WriterFunc(func(w io.Writer) {
		io.WriteString(w, strconv.FormatUint(u, 10))
	})
}

func UnmarshalUint(v interface{}) (uint64, error) {
	switch v := v.(type) {
	case string:
		return strconv.ParseUint(v, 10, 64)
	case int:
		return uint64(v), nil
	case int32:
		return uint64(v), nil
	case int64:
		return uint64(v), nil
	case uint:
		return uint64(v), nil
	case uint32:
		return uint64(v), nil
	case uint64:
		return v, nil
	case json.Number:
		return strconv.ParseUint(string(v), 10, 64)
	default:
		return 0, fmt.Errorf("%T is not a supported integer type", v)
	}
}
