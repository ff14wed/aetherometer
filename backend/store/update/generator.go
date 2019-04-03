package update

import (
	"errors"
	"math"
	"reflect"
	"time"

	"github.com/ff14wed/sibyl/backend/datasheet"
	"github.com/ff14wed/sibyl/backend/models"
	"github.com/ff14wed/sibyl/backend/store"
	"github.com/ff14wed/xivnet/v3"
)

// ErrorStreamNotFound is returned if the given stream ID is not found in the
// list of streams
var ErrorStreamNotFound = errors.New("stream not found")

// ErrorEntityNotFound is returned if the given entity ID is not found in the
// list of entities for the given stream
var ErrorEntityNotFound = errors.New("entity not found")

type updateFactory func(pid int, b *xivnet.Block, data *datasheet.Collection) store.Update

var ingressRegistry = make(map[reflect.Type]updateFactory)
var egressRegistry = make(map[reflect.Type]updateFactory)

func registerIngressHandler(example interface{}, f updateFactory) {
	t := reflect.TypeOf(example)
	ingressRegistry[t] = f
}

func registerEgressHandler(example interface{}, f updateFactory) {
	t := reflect.TypeOf(example)
	egressRegistry[t] = f
}

// Generator generates store updates for a given process and direction of
// data traffic
type Generator struct {
	data *datasheet.Collection
}

// NewGenerator returns a new Generator that converts blocks to updates.
// 	- data stores the lookup tables for game information needed by
func NewGenerator(data *datasheet.Collection) Generator {
	return Generator{data: data}
}

// Generate creates an update based off the block received
// 	- pid is the process ID for which updates are generated
// 	- isEgress is true if the input blocks were sent from the client to the server
func (g *Generator) Generate(pid int, isEgress bool, b *xivnet.Block) store.Update {
	registry := ingressRegistry
	if isEgress {
		registry = egressRegistry
	}
	if factory, found := registry[reflect.TypeOf(b.Data)]; found {
		return factory(pid, b, g.data)
	}
	return nil
}

func getCanonicalOrientation(d uint32, max uint32) float64 {
	factor := float32(d) / float32(max)
	return float64(factor) * 2 * math.Pi
}

func getTimeForDuration(secs float32) time.Time {
	nsecDuration := float64(secs) * float64(int64(time.Second)/int64(time.Nanosecond))
	return time.Unix(0, int64(nsecDuration))
}

type entityUpdateFunc func(s *models.Stream, e *models.Entity) ([]models.StreamEvent, []models.EntityEvent, error)

func validateEntityUpdate(streams *store.Streams, pid int, entityID uint64, u entityUpdateFunc) ([]models.StreamEvent, []models.EntityEvent, error) {
	stream, found := streams.Map[pid]
	if !found {
		return nil, nil, ErrorStreamNotFound
	}
	if stream.CharacterID == 0 {
		return nil, nil, nil
	}
	entity, found := stream.EntitiesMap[entityID]
	if !found {
		return nil, nil, ErrorEntityNotFound
	}
	if entity == nil {
		return nil, nil, nil
	}
	return u(stream, entity)
}
