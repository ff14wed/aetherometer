package store

import (
	"errors"

	"github.com/ff14wed/sibyl/backend/models"
)

// ErrRequestTimedOut is returned if the store is taking too long to return from
// query
var ErrRequestTimedOut = errors.New("request timed out")

type internalRequest interface {
	isInternalRequest()
}

type streamsRequest struct {
	respChan chan []models.Stream
}

func (streamsRequest) isInternalRequest() {}

func (p *Provider) handleStreamsRequest(req streamsRequest) {
	streams := make([]models.Stream, len(p.streams.KeyOrder))
	for i, k := range p.streams.KeyOrder {
		streams[i] = p.streams.Map[k].Clone()
	}
	req.respChan <- streams
}

type streamRequest struct {
	respChan chan *models.Stream
	streamID int
}

func (streamRequest) isInternalRequest() {}

func (p *Provider) handleStreamRequest(req streamRequest) {
	if s, found := p.streams.Map[req.streamID]; found {
		sClone := s.Clone()
		req.respChan <- &sClone
		return
	}
	req.respChan <- nil
}

type entityRequest struct {
	respChan chan *models.Entity
	streamID int
	entityID uint64
}

func (entityRequest) isInternalRequest() {}

func (p *Provider) handleEntityRequest(req entityRequest) {
	if s, found := p.streams.Map[req.streamID]; found {
		if e, found := s.EntitiesMap[req.entityID]; found {
			if e == nil {
				req.respChan <- nil
				return
			}
			eClone := e.Clone()
			req.respChan <- &eClone
			return
		}
	}
	req.respChan <- nil
}
