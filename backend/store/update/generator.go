package update

import (
	"errors"
	"math"
	"reflect"
	"time"

	"github.com/ff14wed/sibyl/backend/datasheet"
	"github.com/ff14wed/xivnet"
)

var ErrorStreamNotFound = errors.New("stream not found")
var ErrorEntityNotFound = errors.New("entity not found")

type updateFactory func(pid int, b *xivnet.Block, data *datasheet.Collection) Update

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

// Generator generates database updates for a given process and direction of
// data traffic
type Generator struct {
	pid      int
	isEgress bool
	data     *datasheet.Collection
}

// NewGenerator returns a new Generator that converts blocks to updates.
// - pid is the process ID for which updates are generated
// - isEgress is true if the input blocks were sent from the client to the server
// - data stores the lookup tables for game information needed by
func NewGenerator(pid int, isEgress bool, data *datasheet.Collection) Generator {
	return Generator{pid: pid, isEgress: isEgress, data: data}
}

// Generate creates an update based off the block received
func (g *Generator) Generate(b *xivnet.Block) Update {
	registry := ingressRegistry
	if g.isEgress {
		registry = egressRegistry
	}
	if factory, found := registry[reflect.TypeOf(b.Data)]; found {
		return factory(g.pid, b, g.data)
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
