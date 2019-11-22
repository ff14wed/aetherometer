package pcap

import (
	"github.com/thejerf/suture"
	"go.uber.org/zap"

	"github.com/ff14wed/aetherometer/core/config"
	"github.com/ff14wed/aetherometer/core/stream"
)

// Adapter defines the implementation of the pcap Adapter
type Adapter struct {
	*suture.Supervisor

}

// AdapterConfig provides commonly accessed configuration for the services
// that make up the hook Adapter
type AdapterConfig struct {
	PcapConfig config.PcapConfig

	StreamUp   chan<- stream.Provider
	StreamDown chan<- int
}

// NewAdapter creates a new instance of the pcap Adapter
func NewAdapter(cfg AdapterConfig, logger *zap.Logger) *Adapter {
	hookLogger := logger.Named("pcap-adapter")
	supervisorLogger := hookLogger.Named("supervisor")
	a := &Adapter{
		Supervisor: suture.New("pcap-adapter", suture.Spec{
			Log: func(line string) {
				supervisorLogger.Info(line)
			},
		}),
	}

	streamSupervisorLogger := hookLogger.Named("stream-supervisor")
	streamSupervisor := suture.New("stream-supervisor", suture.Spec{
		Log: func(line string) {
			streamSupervisorLogger.Info(line)
		},
	})

	providerPool := NewProviderPool()

	streamFactory := &ffxivStreamFactory{
		cfg: cfg,
		providerPool: providerPool,
		streamSupervisor: streamSupervisor,
		logger: logger,
	}

	streamPool := tcpassembly.NewStreamPool(streamFactory)

	manager := NewManager(
		streamPool,
		[]string{"get", "list", "of", "devices", "somehow"},
	)

	a.Add(streamSupervisor)
	a.Add(manager)

	return a
}
