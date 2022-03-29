package app

import (
	"os"
	"path/filepath"
)

func GetCurrentDirectory() (string, error) {
	configDir, err := os.UserConfigDir()
	if err != nil {
		return "", err
	}
	appPath := filepath.Join(configDir, "aetherometer.exe")

	err = os.Mkdir(appPath, 0755)
	if err != nil && !os.IsExist(err) {
		return "", err
	}

	return filepath.Clean(appPath), nil
}
