//go:build !windows
// +build !windows

package main

import "github.com/ff14wed/aetherometer/core/config"

func defaultConfig() (config.Config, error) {
	execPath, err := os.Executable()
	if err != nil {
		return config.Config{}, err
	}

	cleanPath, err := filepath.EvalSymlinks(execPath)
	if err != nil {
		return config.Config{}, err
	}
	dirPath := filepath.Dir(cleanPath)

	return config.Config{
		APIPort:  0,
		DataPath: filepath.Join(dirPath, "resources", "datasheets"),
		Maps: config.MapConfig{
			Cache: filepath.Join(dirPath, "resources", "maps"),
		},
		Adapters: config.Adapters{
			Hook: config.HookConfig{
				Enabled: false,
			},
		},
	}, nil
}
