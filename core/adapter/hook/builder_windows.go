package hook

import (
	"github.com/ff14wed/aetherometer/core/config"
	"github.com/ff14wed/aetherometer/core/stream"
	"github.com/ff14wed/aetherometer/core/win32"
	"go.uber.org/zap"
)

// GetInfo returns the adapter's name and a builder for the adapter.
func GetInfo() stream.AdapterInfo {
	return stream.AdapterInfo{
		Name:    "Hook",
		Builder: &builder{},
	}
}

type builder struct {
	cfg config.Config
}

// LoadConfig loads the configuration for the adapter into the builder.
func (b *builder) LoadConfig(cfg config.Config) error {
	b.cfg = cfg
	return nil
}

// Build returns a new instance of the hook adapter with access to the Windows
// API
func (b *builder) Build(
	streamUp chan<- stream.Provider,
	streamDown chan<- int,
	logger *zap.Logger,
) stream.Adapter {
	p := win32.Provider{}

	return NewAdapter(
		AdapterConfig{
			HookConfig:            b.cfg.Adapters.Hook,
			StreamUp:              streamUp,
			StreamDown:            streamDown,
			RemoteProcessProvider: p,
			ProcessEnumerator:     p,
		},
		logger,
	)
}
