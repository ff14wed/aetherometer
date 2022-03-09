package main

import (
	"path/filepath"

	"github.com/ff14wed/aetherometer/core/config"
)

func defaultConfig() (config.Config, error) {
	dirPath, err := getCurrentDirectory()
	if err != nil {
		return config.Config{}, err
	}

	return config.Config{
		APIPort: 0,
		Sources: config.Sources{
			DataPath: filepath.Join(dirPath, "resources", "datasheets"),
			Maps: config.MapConfig{
				Cache: filepath.Join(dirPath, "resources", "maps"),
			},
		},
		Adapters: config.Adapters{
			Hook: config.HookConfig{
				Enabled:      true,
				DLLPath:      filepath.Join(dirPath, "resources", "win", "xivhook.dll"),
				FFXIVProcess: "ffxiv_dx11.exe",
			},
		},
	}, nil
}
