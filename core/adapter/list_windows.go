package adapter

import (
	"github.com/ff14wed/aetherometer/core/adapter/hook"
	"github.com/ff14wed/aetherometer/core/stream"
)

// Inventory enumerates the adapters that are compatible with Windows.
func Inventory() []stream.AdapterInfo {
	return []stream.AdapterInfo{
		hook.GetInfo(),
	}
}
