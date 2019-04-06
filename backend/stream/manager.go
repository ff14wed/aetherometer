package stream

import (
	"fmt"
	"sync"

	"github.com/ff14wed/sibyl/backend/store"
	"github.com/ff14wed/sibyl/backend/store/update"
	"github.com/ff14wed/xivnet/v3"
	"github.com/thejerf/suture"
	"go.uber.org/zap"
)

type HandlerFactoryArgs struct {
	PID         int
	IngressChan <-chan *xivnet.Frame
	EgressChan  <-chan *xivnet.Frame
	UpdateChan  chan<- store.Update
	Generator   update.Generator
	Logger      *zap.Logger
}

// Handler defines the interface that a stream handler must implement
type Handler interface {
	suture.Service
}

// HandlerFactory defines a factory for creating a handler capable of processing
// updates from the stream provider
type HandlerFactory func(h HandlerFactoryArgs) Handler

// Manager is a process responsible for watching stream created or stream
// closed events from all adapters
type Manager struct {
	generator        update.Generator
	updateChan       chan<- store.Update
	streamSupervisor *suture.Supervisor
	handlerFactory   HandlerFactory
	logger           *zap.Logger

	stop     chan struct{}
	stopDone chan struct{}

	streamTokens  map[int]suture.ServiceToken
	providers     map[int]Provider
	providersLock sync.Mutex

	streamUp   chan Provider
	streamDown chan int
}

// NewManager returns a new stream Manager.
func NewManager(
	generator update.Generator,
	updateChan chan<- store.Update,
	streamSupervisor *suture.Supervisor,
	handlerFactory HandlerFactory,
	logger *zap.Logger,
) *Manager {
	return &Manager{
		generator:        generator,
		updateChan:       updateChan,
		streamSupervisor: streamSupervisor,
		handlerFactory:   handlerFactory,
		logger:           logger.Named("stream-manager"),

		stop:     make(chan struct{}),
		stopDone: make(chan struct{}),

		streamTokens: make(map[int]suture.ServiceToken),
		providers:    make(map[int]Provider),

		streamUp:   make(chan Provider, 64),
		streamDown: make(chan int, 64),
	}
}

// Serve runs the stream manager. It is responsible for spinning up stream
// handlers whenever new streams are created and shutting down stream handlers
// when streams are closed.
func (m *Manager) Serve() {
	m.logger.Info("Running")
	for {
		select {
		case sp := <-m.streamUp:
			pid := sp.PID()
			ingressChan := sp.SubscribeIngress()
			egressChan := sp.SubscribeEgress()
			sh := m.handlerFactory(HandlerFactoryArgs{
				PID:         pid,
				IngressChan: ingressChan,
				EgressChan:  egressChan,
				UpdateChan:  m.updateChan,
				Generator:   m.generator,
				Logger:      m.logger,
			})
			token := m.streamSupervisor.Add(sh)
			m.streamTokens[pid] = token

			m.providersLock.Lock()
			m.providers[pid] = sp
			m.providersLock.Unlock()
		case pid := <-m.streamDown:
			err := m.streamSupervisor.Remove(m.streamTokens[pid])
			if err != nil {
				m.logger.Error("Error removing stream", zap.Int("streamID", pid), zap.Error(err))
			}
			delete(m.streamTokens, pid)

			m.providersLock.Lock()
			delete(m.providers, pid)
			m.providersLock.Unlock()
		case <-m.stop:
			m.logger.Info("Stopping...")
			close(m.stopDone)
			return
		}
	}
}

// Stop will shutdown this service eventually, but will not wait on it to stop
func (m *Manager) Stop() {
	close(m.stop)
	<-m.stopDone
}

// SendRequest forwards a request for a given stream ID to the correct stream
// Provider.
func (m *Manager) SendRequest(pid int, req []byte) ([]byte, error) {
	m.providersLock.Lock()
	provider, found := m.providers[pid]
	m.providersLock.Unlock()

	if found {
		return provider.SendRequest(req)
	}
	return nil, fmt.Errorf("stream provider %d not found", pid)
}

// StreamUp returns a channel that allows an upstream service to notify the
// manager that a new stream has been created.
func (m *Manager) StreamUp() chan<- Provider {
	return m.streamUp
}

// StreamDown returns a channel that allows an upstream service to notify the
// manager that a stream has closed.
func (m *Manager) StreamDown() chan<- int {
	return m.streamDown
}
