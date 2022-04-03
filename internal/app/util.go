package app

import (
	"context"
	"os"
	"path/filepath"
)

type headlessCtxType string

const headlessCtx headlessCtxType = "headless"

func GetAppDirectory() (string, error) {
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

func HeadlessContext() context.Context {
	ctx := context.Background()
	return context.WithValue(ctx, headlessCtx, true)
}

func IsHeadless(ctx context.Context) bool {
	return ctx.Value(headlessCtx) == true
}
