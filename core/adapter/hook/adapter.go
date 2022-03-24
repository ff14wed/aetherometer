package hook

import (
	"net"
	"time"

	"github.com/thejerf/suture"
	"go.uber.org/zap"

	"github.com/ff14wed/aetherometer/core/config"
	"github.com/ff14wed/aetherometer/core/process"
	"github.com/ff14wed/aetherometer/core/stream"
)

// Adapter defines the implementation of the hook Adapter
type Adapter struct {
	*suture.Supervisor
}

// AdapterConfig provides commonly accessed configuration for the hook
// to the services that make up the hook Adapter
type AdapterConfig struct {
	HookConfig config.HookConfig

	StreamUp   chan<- stream.Provider
	StreamDown chan<- int

	RemoteProcessProvider RemoteProcessProvider
	ProcessEnumerator     process.Enumerator
}

//go:generate go run github.com/maxbrunsfeld/counterfeiter/v6 . RemoteProcessProvider

// RemoteProcessProvider defines the interface that exposes methods for interacting
// with other processes on the system
type RemoteProcessProvider interface {
	InjectDLL(processID uint32, payloadPath string) error
	DialPipe(path string, timeout *time.Duration) (net.Conn, error)
	IsPipeClosed(err error) bool
}

// DLLAlreadyInjectedError is an error that indicates the DLL has already been
// injected.
type DLLAlreadyInjectedError interface {
	IsDLLAlreadyInjectedError()
}

// NewAdapter creates a new instance of the hook Adapter
func NewAdapter(cfg AdapterConfig, logger *zap.Logger) *Adapter {
	hookLogger := logger.Named("hook-adapter")
	supervisorLogger := hookLogger.Named("supervisor")
	a := &Adapter{
		Supervisor: suture.New("hook-adapter", suture.Spec{
			Log: func(line string) {
				supervisorLogger.Info(line)
			},
		}),
	}

	scanTicker := time.NewTicker(1 * time.Second)
	scanner := process.NewScanner(
		cfg.HookConfig.FFXIVProcess,
		scanTicker.C,
		cfg.ProcessEnumerator,
		10,
		hookLogger,
	)

	streamSupervisorLogger := hookLogger.Named("stream-supervisor")
	streamSupervisor := suture.New("stream-supervisor", suture.Spec{
		Log: func(line string) {
			streamSupervisorLogger.Info(line)
		},
	})

	streamBuilder := func(streamID uint32) Stream {
		return NewStream(streamID, cfg, hookLogger)
	}

	manager := NewManager(
		cfg,
		scanner.ProcessAddEventListener(),
		scanner.ProcessRemoveEventListener(),
		streamBuilder,
		streamSupervisor,
		hookLogger,
	)

	a.Add(scanner)
	a.Add(streamSupervisor)
	a.Add(manager)

	return a
}
