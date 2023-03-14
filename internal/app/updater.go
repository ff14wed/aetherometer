package app

import (
	"bytes"
	"crypto/sha256"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"syscall"

	"github.com/ff14wed/aetherometer/core/config"

	"github.com/BurntSushi/toml"
	"github.com/creativeprojects/go-selfupdate"
	"github.com/sqweek/dialog"
	"go.uber.org/zap"
)

const datasheetsUpstream string = "https://raw.githubusercontent.com/ff14wed/aetherometer/master/resources/datasheets/"
const hookSha256Sum string = "8d592a4ae901047477b65d4275f1e05242a9e5cc98e72aabcf2468205606ef20"

var dataFiles = []string{
	"Action.csv",
	"BNpcBase.csv",
	"BNpcName.csv",
	"ClassJob.csv",
	"CraftAction.csv",
	"ENpcResident.csv",
	"Item.csv",
	"Map.csv",
	"ModelChara.csv",
	"ModelSkeleton.csv",
	"Omen.csv",
	"PlaceName.csv",
	"Recipe.csv",
	"RecipeLevelTable.csv",
	"Status.csv",
	"TerritoryType.csv",
	"World.csv",
}

func ensureDirectory(dirPath string, name string) error {
	err := os.MkdirAll(dirPath, 0755)
	if err != nil {
		return fmt.Errorf("unable to ensure %s directory exists: %s", name, err)
	}
	return nil
}

func readDataFromURL(url string) ([]byte, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return body, nil
}

func updateVersion(dataPath string) (updated bool, err error) {
	upstreamVersion, err := readDataFromURL(fmt.Sprintf("%s/VERSION", datasheetsUpstream))
	if err != nil {
		return false, err
	}
	localVersionFile := filepath.Join(dataPath, "VERSION")
	localVersionBytes, err := ioutil.ReadFile(localVersionFile)
	if err == nil && bytes.Equal(upstreamVersion, localVersionBytes) {
		return false, nil
	}
	err = ioutil.WriteFile(localVersionFile, upstreamVersion, 0755)
	if err != nil {
		return false, err
	}
	return true, nil
}

func updateDataFiles(dataPath string, versionUpdated bool, logger *zap.Logger) error {
	if versionUpdated {
		go func() {
			msgBuilder := dialog.Message("Downloading all datasheets... This may take some time.")
			msgBuilder.Title("Aetherometer update in progress").Info()
		}()
	}
	for _, df := range dataFiles {
		shouldUpdate := versionUpdated

		dfPath := filepath.Join(dataPath, df)
		if _, err := os.Stat(dfPath); errors.Is(err, os.ErrNotExist) {
			shouldUpdate = true
		}

		if !shouldUpdate {
			continue
		}

		logger.Info("Downloading datasheet", zap.String("file", dfPath))

		dfBytes, err := readDataFromURL(fmt.Sprintf("%s/%s", datasheetsUpstream, df))
		if err != nil {
			return err
		}
		if err := ioutil.WriteFile(dfPath, dfBytes, 0755); err != nil {
			return err
		}
	}

	return nil
}

func ensureHookAdapter(hookPath string, logger *zap.Logger) error {
	logger.Info("Checking if hook DLL exists locally...")
	if _, err := os.Stat(hookPath); errors.Is(err, os.ErrNotExist) {
		logger.Info("Downloading hook DLL...")

		go func() {
			msgBuilder := dialog.Message("Downloading hook DLL.... This may take some time.")
			msgBuilder.Title("Aetherometer update in progress").Info()
		}()

		hookDLLBytes, err := readDataFromURL("https://github.com/ff14wed/deucalion/releases/download/0.9.1/deucalion.dll")
		if err != nil {
			return err
		}
		buf := bytes.NewBuffer(hookDLLBytes)
		hasher := sha256.New()
		if _, err := io.Copy(hasher, buf); err != nil {
			return err
		}
		downloadedSha256Sum := fmt.Sprintf("%x", hasher.Sum(nil))
		if downloadedSha256Sum != hookSha256Sum {
			return fmt.Errorf(
				"downloaded Hook DLL failed SHA256 sum check: %s vs %s",
				downloadedSha256Sum,
				hookSha256Sum,
			)
		}
		if err := ioutil.WriteFile(hookPath, hookDLLBytes, 0755); err != nil {
			return err
		}
	} else {
		logger.Info("Hook DLL already exists. Nothing to do.")
	}
	return nil
}

func checkResourceUpdates(configFile string, logger *zap.Logger) error {
	cfg := config.Config{}
	if _, err := os.Stat(configFile); errors.Is(err, os.ErrNotExist) {
		// Config file doesn't exist, so write the config to disk first
		var cfgErr error
		cfg, cfgErr = defaultConfig()
		if cfgErr != nil {
			return fmt.Errorf("unable to get default config: %s", cfgErr)
		}
	} else if err != nil {
		return fmt.Errorf("unable to check for config file: %s", err)
	} else {
		_, decodeErr := toml.DecodeFile(configFile, &cfg)
		if decodeErr != nil {
			return fmt.Errorf("unable to read config file: %s", decodeErr)
		}
	}
	if !cfg.AutoUpdate {
		logger.Info("Skipping auto-update of resources because it is disabled in config")
		return nil
	}

	if err := ensureDirectory(cfg.Sources.DataPath, "sources.data_path"); err != nil {
		return err
	}

	if err := ensureDirectory(cfg.Sources.Maps.Cache, "sources.maps.cache"); err != nil {
		return err
	}

	if cfg.Adapters.Hook.Enabled && cfg.Adapters.Hook.DLLPath != "" {
		if err := ensureDirectory(filepath.Dir(cfg.Adapters.Hook.DLLPath), "sources.adapters.hook.dll_path"); err != nil {
			return err
		}
		if err := ensureHookAdapter(cfg.Adapters.Hook.DLLPath, logger); err != nil {
			return fmt.Errorf("could not download hook library: %s", err)
		}
	}
	versionUpdated, err := updateVersion(cfg.Sources.DataPath)
	if err != nil {
		return fmt.Errorf("unable to update version file: %s", err)
	}
	if err := updateDataFiles(cfg.Sources.DataPath, versionUpdated, logger); err != nil {
		return fmt.Errorf("unable to update data files: %s", err)
	}

	logger.Info("Finished checking for updates.")
	return nil
}

func restartApp(exe string) error {
	if runtime.GOOS == "windows" {
		cmd := exec.Command(exe, os.Args[1:]...)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		cmd.Stdin = os.Stdin
		cmd.Env = os.Environ()
		err := cmd.Run()
		if err == nil {
			os.Exit(0)
		}
		return err
	}

	return syscall.Exec(exe, append([]string{exe}, os.Args[1:]...), os.Environ())
}

func checkAppUpdate(version string, logger *zap.Logger) error {
	if version == "development" {
		logger.Info("Skipping app update for development version.")
		return nil
	}
	latest, found, err := selfupdate.DetectLatest("ff14wed/aetherometer")
	if err != nil {
		logger.Error("Unable to detect latest version... Skipping app update.", zap.Error(err))
		return nil
	}
	if !found {
		logger.Error("No version could be found from github repository... Skipping app update.", zap.Error(err))
		return nil
	}
	if latest.LessOrEqual(version) {
		logger.Info("Current version is the latest.", zap.String("version", version))
		return nil
	}

	logger.Info("Found newer Aetheromter version.", zap.String("version", latest.Version()))

	updateDlg := dialog.Message("A new version of Aetherometer (%s) is available. Update?", latest.Version())
	updateDlg.Title("Aetherometer Update")
	if !updateDlg.YesNo() {
		logger.Info("Skipping app update due to user cancellation.")
		return nil
	}

	exe, err := os.Executable()
	if err != nil {
		return fmt.Errorf("could not locate executable path: %s", err)
	}
	if err := selfupdate.UpdateTo(latest.AssetURL, latest.AssetName, exe); err != nil {
		return fmt.Errorf("error occurred while updating binary: %s", err)
	}
	logger.Info("Successfully updated to latest version", zap.String("version", latest.Version()))

	if err := restartApp(exe); err != nil {
		return fmt.Errorf("error restarting app: %s", err)
	}
	return nil
}

func CheckUpdates(version string, configFile string, logger *zap.Logger) error {
	logger.Info("Checking for updates...")

	if err := checkAppUpdate(version, logger); err != nil {
		return err
	}
	if err := checkResourceUpdates(configFile, logger); err != nil {
		return err
	}

	return nil
}
