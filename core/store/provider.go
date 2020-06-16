package store

import (
	"fmt"
	"time"

	"github.com/ff14wed/aetherometer/core/hub"
	"github.com/ff14wed/aetherometer/core/models"
	"go.uber.org/zap"
)

// Provider provides access to the store. It runs as a long running service
// that updates the store in response to update events. All updates and
// accesses are handled in an evented loop to serialize changes to the store
// for thread safety.
// Provider also emits events for updates made to the store.
type Provider struct {
	queryTimeout time.Duration
	logger       *zap.Logger

	streams   Streams
	streamHub *hub.StreamHub
	entityHub *hub.EntityHub

	updatesChan         chan Update
	internalRequestChan chan internalRequest

	stop     chan struct{}
	stopDone chan struct{}
}

// NewProvider creates a new store provider. It will also initialize the
// store.
// Options can optionally be provided like follows:
// 	provider = store.NewProvider(
// 		logger,
// 		store.WithQueryTimeout(10*time.Millisecond),
// 		store.WithUpdateBufferSize(10),
// 		store.WithEventBufferSize(10),
// 		store.WithRequestBufferSize(10),
// 	)
func NewProvider(
	logger *zap.Logger,
	opts ...Option,
) *Provider {
	cfg := providerConfig{
		queryTimeout:      5 * time.Second,
		updateBufferSize:  10000,
		eventBufferSize:   10000,
		requestBufferSize: 10,
	}
	for _, opt := range opts {
		opt(&cfg)
	}

	return &Provider{
		queryTimeout: cfg.queryTimeout,
		logger:       logger.Named("store-provider"),

		streams:   Streams{Map: make(map[int]*models.Stream)},
		streamHub: hub.NewStreamHub(cfg.eventBufferSize),
		entityHub: hub.NewEntityHub(cfg.eventBufferSize),

		updatesChan:         make(chan Update, cfg.updateBufferSize),
		internalRequestChan: make(chan internalRequest, cfg.requestBufferSize),

		stop:     make(chan struct{}),
		stopDone: make(chan struct{}),
	}
}

// Serve runs the main loop for the provider. It runs inside a goroutine
// as a service and is responsible for exclusively handling all reads and
// updates, including serving read requests.
func (p *Provider) Serve() {
	defer close(p.stopDone)
	p.logger.Info("Running")
	for {
		select {
		case u := <-p.updatesChan:
			p.handleUpdate(u)
		case r := <-p.internalRequestChan:
			switch v := r.(type) {
			case streamsRequest:
				p.handleStreamsRequest(v)
			case streamRequest:
				p.handleStreamRequest(v)
			case entityRequest:
				p.handleEntityRequest(v)
			}
		case <-p.stop:
			p.logger.Info("Stopping...")
			return
		}
	}
}

// Stop will shutdown this service and wait on it to stop before returning
func (p *Provider) Stop() {
	close(p.stop)
	<-p.stopDone
}

// UpdatesChan returns a channel on which other services can send store updates
func (p *Provider) UpdatesChan() chan<- Update {
	return p.updatesChan
}

func (p *Provider) handleUpdate(u Update) {
	if u == nil {
		return
	}
	streamEvents, entityEvents, err := u.ModifyStore(&p.streams)
	if err != nil {
		p.logger.Error("Error applying update",
			zap.String("update", fmt.Sprintf("%#v", u)),
			zap.Error(err),
		)
	}
	for _, streamEvent := range streamEvents {
		p.streamHub.Broadcast(streamEvent)
	}
	for _, entityEvent := range entityEvents {
		p.entityHub.Broadcast(entityEvent)
	}
}

// Streams returns all the streams from the internal store. This query will
// return an error if the request exceeds the timeout duration.
func (p *Provider) Streams() ([]models.Stream, error) {
	streamsChan := make(chan []models.Stream, 1)
	p.internalRequestChan <- streamsRequest{respChan: streamsChan}
	select {
	case resp := <-streamsChan:
		return resp, nil
	case <-time.After(p.queryTimeout):
		p.logger.Error("Streams()",
			zap.Error(ErrRequestTimedOut),
			zap.Duration("timeout-duration", p.queryTimeout),
		)
		return nil, ErrRequestTimedOut
	}
}

// Stream returns a specific stream from the store, queried by streamID. This
// query will return an error if the request exceeds the timeout duration.
func (p *Provider) Stream(streamID int) (*models.Stream, error) {
	streamChan := make(chan *models.Stream, 1)
	p.internalRequestChan <- streamRequest{
		respChan: streamChan,
		streamID: streamID,
	}
	select {
	case resp := <-streamChan:
		if resp == nil {
			return nil, fmt.Errorf("stream ID %d not found", streamID)
		}
		return resp, nil
	case <-time.After(p.queryTimeout):
		p.logger.Error("Stream()",
			zap.Error(ErrRequestTimedOut),
			zap.Int("streamID", streamID),
			zap.Duration("timeout-duration", p.queryTimeout),
		)
		return nil, ErrRequestTimedOut
	}
}

// Entity returns a specific entity in a specific from the store, queried by
// streamID and entityID. It returns an error if the stream ID is not found or
// if the entityID is not found in the stream. This query will return an error
// if the request exceeds the timeout duration.
func (p *Provider) Entity(streamID int, entityID uint64) (*models.Entity, error) {
	entityChan := make(chan *models.Entity, 1)
	p.internalRequestChan <- entityRequest{
		respChan: entityChan,
		streamID: streamID,
		entityID: entityID,
	}
	select {
	case resp := <-entityChan:
		if resp == nil {
			return nil, fmt.Errorf("entity ID %d not found in stream %d", entityID, streamID)
		}
		return resp, nil
	case <-time.After(p.queryTimeout):
		p.logger.Error("Entity()",
			zap.Error(ErrRequestTimedOut),
			zap.Int("streamID", streamID),
			zap.Uint64("entityID", entityID),
			zap.Duration("timeout-duration", p.queryTimeout),
		)
		return nil, ErrRequestTimedOut
	}
}

// StreamEventSource returns an event source that allows consumers
// to subscribe to stream events
func (p *Provider) StreamEventSource() models.StreamEventSource {
	return p.streamHub
}

// EntityEventSource returns an event source that allows consumers
// to subscribe to entity events
func (p *Provider) EntityEventSource() models.EntityEventSource {
	return p.entityHub
}
