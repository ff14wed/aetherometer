package adapter

import (
	"github.com/ff14wed/sibyl/backend/adapter/hook"
	"github.com/ff14wed/sibyl/backend/stream"
)

// Inventory enumerates the adapters that are compatible with Windows.
func Inventory() []stream.AdapterInfo {
	return []stream.AdapterInfo{
		hook.GetInfo(),
	}
}
