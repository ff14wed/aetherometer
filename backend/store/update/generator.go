package update

import (
	"errors"
	"math"
	"reflect"
	"time"

	"github.com/ff14wed/sibyl/backend/datasheet"
	"github.com/ff14wed/sibyl/backend/store"
	"github.com/ff14wed/xivnet"
)

var ErrorStreamNotFound = errors.New("stream not found")
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
// - data stores the lookup tables for game information needed by
func NewGenerator(data *datasheet.Collection) Generator {
	return Generator{data: data}
}

// Generate creates an update based off the block received
// - pid is the process ID for which updates are generated
// - isEgress is true if the input blocks were sent from the client to the server
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
