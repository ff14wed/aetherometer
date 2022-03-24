package app

import (
	"os"
	"path/filepath"
)

func GetCurrentDirectory() (string, error) {
	execPath, err := os.Executable()
	if err != nil {
		return "", err
	}

	cleanPath, err := filepath.EvalSymlinks(execPath)
	if err != nil {
		return "", err
	}
	return filepath.Dir(cleanPath), nil
}
