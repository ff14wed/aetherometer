package config

// HookConfig stores the configuration for the hook adapter
type HookConfig struct {
	// Enabled toggles whether or not the Hook adapter is enabled.
	Enabled bool `toml:"enabled"`

	// DLLPath sets the path of the Hook DLL on the system.
	DLLPath string `toml:"dll_path" validate:"file"`

	// FFXIVProcess is the name of the exe file for the game.
	FFXIVProcess string `toml:"ffxiv_process" validate:"nonempty"`
}
