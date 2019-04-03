package config

// HookConfig stores the configuration for the hook adapter
type HookConfig struct {
	Enabled      bool   `toml:"enabled"`
	DLLPath      string `toml:"dll_path" validate:"file"`
	FFXIVProcess string `toml:"ffxiv_process" validate:"nonempty"`
}
