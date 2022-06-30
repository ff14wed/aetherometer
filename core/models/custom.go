package models

import (
	"encoding/json"
	"fmt"
	"io"
	"reflect"
	"sort"
	"strconv"
	"time"

	"github.com/99designs/gqlgen/graphql"
	"gopkg.in/dealancer/validate.v2"
)

//go:generate go run ../scripts/gqlgen.go

// Stream represents state reconstructed from the live stream of data from a
// running FFXIV instance.
type Stream struct {
	ID          int `json:"id"`
	ServerID    int `json:"serverID"`
	InstanceNum int `json:"instanceNum"`

	CharacterID  uint64 `json:"characterID"`
	HomeWorld    World  `json:"homeWorld"`
	CurrentWorld World  `json:"currentWorld"`

	Place        Place         `json:"place"`
	Enmity       Enmity        `json:"enmity"`
	CraftingInfo *CraftingInfo `json:"craftingInfo"`

	Stats *Stats `json:"stats"`

	EntitiesMap map[uint64]*Entity `json:"entities"`
}

// Entities returns all the entities from the stream, sorted in order by index.
func (s *Stream) Entities() []Entity {
	var entities []Entity
	for _, e := range s.EntitiesMap {
		if e == nil {
			continue
		}
		entities = append(entities, *e)
	}
	sort.SliceStable(entities, func(i, j int) bool {
		return entities[i].Index < entities[j].Index
	})
	return entities
}

// MarshalTimestamp converts the provided time to milliseconds since the Unix
// epoch.
func MarshalTimestamp(t time.Time) graphql.Marshaler {
	return graphql.WriterFunc(func(w io.Writer) {
		_, _ = io.WriteString(w, strconv.FormatInt(getTimeInMs(t), 10))
	})
}

// UnmarshalTimestamp is currently unimplemented.
func UnmarshalTimestamp(v interface{}) (time.Time, error) {
	panic("not implemented")
}

func getTimeInMs(t time.Time) int64 {
	return t.UnixNano() / (int64(time.Millisecond) / int64(time.Nanosecond))
}

// MarshalUint marshals the provided uint64 to a string
func MarshalUint(u uint64) graphql.Marshaler {
	return graphql.WriterFunc(func(w io.Writer) {
		_, _ = io.WriteString(w, strconv.FormatUint(u, 10))
	})
}

// UnmarshalUint converts the string representation for an unsigned integer to a
// uint64.
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

// Validate implements a way to validate this StreamEvent, since the library
// used doesn't dive into interfaces
func (s *StreamEvent) Validate() error {
	return validate.Validate(reflect.ValueOf(s.Type).Interface())
}

// Validate implements a way to validate this EntityEvent, since the library
// used doesn't dive into interfaces
func (e *EntityEvent) Validate() error {
	return validate.Validate(reflect.ValueOf(e.Type).Interface())
}
