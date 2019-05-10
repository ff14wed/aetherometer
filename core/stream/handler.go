package stream

import (
	"fmt"

	"github.com/ff14wed/aetherometer/core/models"
	"github.com/ff14wed/aetherometer/core/store"
	"github.com/ff14wed/aetherometer/core/store/update"
	"github.com/ff14wed/xivnet/v3"
	"go.uber.org/zap"
)

// handler is a process responsible for taking in xivnet Frames from
// adapters and generating updates with them
type handler struct {
	streamID    int
	ingressChan <-chan *xivnet.Frame
	egressChan  <-chan *xivnet.Frame
	updateChan  chan<- store.Update
	generator   update.Generator
	logger      *zap.Logger

	stop     chan struct{}
	stopDone chan struct{}
}

// NewHandler returns a new stream Handler.
func NewHandler(args HandlerFactoryArgs) Handler {
	return &handler{
		streamID:    args.StreamID,
		ingressChan: args.IngressChan,
		egressChan:  args.EgressChan,
		updateChan:  args.UpdateChan,
		generator:   args.Generator,
		logger:      args.Logger.Named(fmt.Sprintf("stream-handler-%d", args.StreamID)),

		stop:     make(chan struct{}),
		stopDone: make(chan struct{}),
	}
}

// Serve runs the main loop for the stream handler. It runs inside a goroutine
// as a service and is responsible for processing data from streams and
// generating updates with that data.
func (h *handler) Serve() {
	defer close(h.stopDone)
	h.logger.Info("Running")
	h.updateChan <- addStreamUpdate{streamID: h.streamID}
	for {
		select {
		case inFrame := <-h.ingressChan:
			for _, parsedBlock := range inFrame.Blocks {
				h.updateChan <- h.generator.Generate(h.streamID, false, parsedBlock)
			}
		case outFrame := <-h.egressChan:
			for _, parsedBlock := range outFrame.Blocks {
				h.updateChan <- h.generator.Generate(h.streamID, true, parsedBlock)
			}
		case <-h.stop:
			h.logger.Info("Stopping...")
			h.updateChan <- removeStreamUpdate{streamID: h.streamID}
			return
		}
	}
}

// Stop will shutdown this service and wait on it to stop before returning
func (h *handler) Stop() {
	close(h.stop)
	<-h.stopDone
}

type addStreamUpdate struct {
	streamID int
}

func (u addStreamUpdate) ModifyStore(streams *store.Streams) ([]models.StreamEvent, []models.EntityEvent, error) {
	s := models.Stream{
		ID:          u.streamID,
		EntitiesMap: make(map[uint64]*models.Entity),
	}

	streams.Map[u.streamID] = &s
	streams.KeyOrder = append(streams.KeyOrder, u.streamID)

	return []models.StreamEvent{models.StreamEvent{
		StreamID: u.streamID,
		Type:     models.AddStream{Stream: s},
	}}, nil, nil
}

type removeStreamUpdate struct {
	streamID int
}

func (u removeStreamUpdate) ModifyStore(streams *store.Streams) ([]models.StreamEvent, []models.EntityEvent, error) {
	delete(streams.Map, u.streamID)
	streamIDX := -1
	for i, v := range streams.KeyOrder {
		if v == u.streamID {
			streamIDX = i
		}
	}
	if streamIDX >= 0 {
		streams.KeyOrder = append(streams.KeyOrder[:streamIDX], streams.KeyOrder[streamIDX+1:]...)
	}

	return []models.StreamEvent{models.StreamEvent{
		StreamID: u.streamID,
		Type:     models.RemoveStream{ID: u.streamID},
	}}, nil, nil
}
