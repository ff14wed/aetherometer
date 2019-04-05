package stream

import (
	"fmt"

	"github.com/ff14wed/sibyl/backend/config"
	"go.uber.org/zap"
)

// BuildAdapterInventory creates all of the enabled adapters provided in the
// inventory list
func BuildAdapterInventory(
	inventory []AdapterInfo,
	cfg config.Config,
	streamUp chan<- Provider,
	streamDown chan<- int,
	logger *zap.Logger,
) (map[string]Adapter, error) {
	adapters := make(map[string]Adapter)
	for _, info := range inventory {
		if !cfg.Adapters.IsEnabled(info.Name) {
			continue
		}
		err := info.Builder.LoadConfig(cfg)
		if err != nil {
			return nil, fmt.Errorf("error creating adapter %s: %s", info.Name, err)
		}
		adapters[info.Name] = info.Builder.Build(streamUp, streamDown, logger)
	}
	return adapters, nil
}
