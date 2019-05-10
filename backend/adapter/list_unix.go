// +build darwin dragonfly freebsd linux netbsd openbsd

package adapter

import "github.com/ff14wed/sibyl/backend/stream"

// Inventory enumerates the adapters that are compatible with Unix based systems.
func Inventory() []stream.AdapterInfo {
	return []stream.AdapterInfo{}
}
