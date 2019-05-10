package hook

import (
	"fmt"

	"github.com/thejerf/suture"
	"go.uber.org/zap"
)

// Manager is responsible for starting up new hook streams whenever it detects
// that a new instance of a watched process is created. It is also responsible
// for shutting down those streams when a watched process is closed.
// Additionally, the Manager notifies the StreamUp and StreamDown channels when
// the watched process is created and closed, respectively.
type Manager struct {
	cfg AdapterConfig

	addProcEventChan <-chan uint32
	remProcEventChan <-chan uint32

	streamBuilder    func(streamID uint32) Stream
	streamSupervisor *suture.Supervisor
	streamTokens     map[uint32]suture.ServiceToken

	logger *zap.Logger

	stop     chan struct{}
	stopDone chan struct{}
}

// NewManager creates a new hook Stream Manager
func NewManager(
	cfg AdapterConfig,
	addProcEventChan <-chan uint32,
	remProcEventChan <-chan uint32,
	streamBuilder func(streamID uint32) Stream,
	streamSupervisor *suture.Supervisor,
	logger *zap.Logger,
) *Manager {
	return &Manager{
		cfg: cfg,

		addProcEventChan: addProcEventChan,
		remProcEventChan: remProcEventChan,

		streamBuilder:    streamBuilder,
		streamSupervisor: streamSupervisor,
		streamTokens:     make(map[uint32]suture.ServiceToken),

		logger: logger.Named("hook-manager"),

		stop:     make(chan struct{}),
		stopDone: make(chan struct{}),
	}
}

// Serve runs the service responsible for handling process add and
// remove events.
func (m *Manager) Serve() {
	defer close(m.stopDone)

	m.logger.Info("Running")
	for {
		select {
		case streamID := <-m.addProcEventChan:
			m.handleProcessAdd(streamID)
		case streamID := <-m.remProcEventChan:
			m.handleProcessRemove(streamID)
		case <-m.stop:
			m.logger.Info("Stopping...")
			return
		}
	}
}

// Stop will shutdown this service and wait on it to stop before returning.
func (m *Manager) Stop() {
	close(m.stop)
	<-m.stopDone
}

func (m *Manager) handleProcessAdd(streamID uint32) {
	s := m.streamBuilder(streamID)
	if s != nil {
		m.streamTokens[streamID] = m.streamSupervisor.Add(s)
		m.cfg.StreamUp <- s
	}
}

func (m *Manager) handleProcessRemove(streamID uint32) {
	err := m.streamSupervisor.Remove(m.streamTokens[streamID])
	if err != nil {
		m.logger.Error("Error removing process group",
			zap.Uint32("streamID", streamID),
			zap.String("token", fmt.Sprintf("%#v", m.streamTokens[streamID])),
			zap.Error(err),
		)
		return
	}
	m.cfg.StreamDown <- int(streamID)
	delete(m.streamTokens, streamID)
}
