package config

// HookConfig stores the configuration for the hook adapter
type HookConfig struct {
	// Enabled toggles whether or not the Hook adapter is enabled.
	Enabled bool `toml:"enabled"`

	// DLLPath sets the path of the Hook DLL on the system.
	DLLPath string `toml:"dll_path" validate:"file"`

	// FFXIVProcess is the name of the exe file for the game.
	FFXIVProcess string `toml:"ffxiv_process" validate:"nonempty"`

	// DialRetryInterval controls how long to wait before retrying
	// failures to make a connection with the hook DLL.
	// Defaults to 500 milliseconds.
	DialRetryInterval Duration `toml:"dial_retry_interval"`

	// PingInterval controls the interval between liveness checks to
	// hook. Defaults to 1 second.
	PingInterval Duration `toml:"ping_interval"`
}

type PcapConfig struct {
	Enabled bool `toml:"enabled"`

	Put        string
	Your       string
	Parameters string
	Here       string
}
