package main

import (
	"context"
	"embed"
	"errors"
	"flag"
	"fmt"
	"net/url"
	"os"
	"os/signal"
	"path/filepath"
	"runtime"
	"syscall"

	"github.com/ff14wed/aetherometer/core/config"
	"github.com/ff14wed/aetherometer/internal/app"

	"github.com/apenwarr/fixconsole"
	"github.com/sqweek/dialog"
	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/logger"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/windows"
	"go.uber.org/zap"
)

//go:embed frontend/dist
var assets embed.FS

var Version = "development"

// Workaround for Windows support for zap from
// https://github.com/uber-go/zap/issues/621
func newWinFileSink(u *url.URL) (zap.Sink, error) {
	// Remove leading slash left by url.Parse()
	return os.OpenFile(u.Path[1:], os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0644)
}
func NewLogger(outputLogPath string, useStdout bool) (*zap.Logger, error) {
	if useStdout {
		outputLogPath = "stdout"
	} else if runtime.GOOS == "windows" {
		err := zap.RegisterSink("winfile", newWinFileSink)
		if err != nil {
			return nil, fmt.Errorf("couldn't register winfile log sink: %s", err)
		}
		outputLogPath = "winfile:///" + filepath.ToSlash(outputLogPath)
	}
	zapCfg := zap.NewDevelopmentConfig()
	zapCfg.DisableStacktrace = true
	zapCfg.DisableCaller = true
	zapCfg.OutputPaths = []string{outputLogPath}
	zapCfg.ErrorOutputPaths = []string{outputLogPath}
	return zapCfg.Build()
}

func startup() error {
	flag.Usage = func() {
		fixconsole.FixConsoleIfNeeded()
		fmt.Fprintf(flag.CommandLine.Output(), "Usage of %s:\n", os.Args[0])
		flag.PrintDefaults()
	}

	dirPath, err := app.GetCurrentDirectory()
	if err != nil {
		return err
	}

	cfgPath := flag.String("c", filepath.Join(dirPath, "config.toml"), "optional path to TOML config file")

	headless := flag.Bool("headless", false, "run Aetherometer in headless mode.")

	helpFlag := flag.Bool("h", false, "displays usage information")

	versionFlag := flag.Bool("v", false, "displays version information")

	flag.Parse()

	if *helpFlag {
		flag.Usage()
		return nil
	}

	if *versionFlag {
		fixconsole.FixConsoleIfNeeded()
		fmt.Println("Aetherometer version:", Version)
		return nil
	}

	outputLogPath := filepath.Join(dirPath, "aetherometer.log")
	if *headless {
		fixconsole.FixConsoleIfNeeded()
	}

	zapLogger, err := NewLogger(outputLogPath, *headless)
	if err != nil {
		return fmt.Errorf("can't initialize zap logger: %v", err)
	}
	defer func() {
		_ = zapLogger.Sync()
	}()
	zap.ReplaceGlobals(zapLogger)

	if len(*cfgPath) == 0 {
		return errors.New("config path cannot be empty")
	}

	defaultCfg, err := app.DefaultConfig()
	if err != nil {
		return fmt.Errorf("could not setup default config: %v", err)
	}

	cfgProvider := config.NewProvider(*cfgPath, defaultCfg, zapLogger)

	a := app.NewApp(cfgProvider, Version, zapLogger)

	// Do not start the GUI in headless mode.
	if *headless {
		a.Startup(context.Background())

		signals := make(chan os.Signal, 32)
		signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)

		sig := <-signals
		zapLogger.Info("Received signal, shutting down...", zap.Stringer("signal", sig))

		a.Shutdown(context.Background())
		return nil
	}

	// Start run wails app if not headless mode
	err = wails.Run(&options.App{
		Title:             "aetherometer",
		Width:             1280,
		Height:            800,
		MinWidth:          1024,
		MinHeight:         768,
		DisableResize:     false,
		Fullscreen:        false,
		Frameless:         true,
		StartHidden:       false,
		HideWindowOnClose: false,
		RGBA:              &options.RGBA{R: 33, G: 37, B: 43, A: 255},
		Assets:            assets,
		LogLevel:          logger.DEBUG,
		OnStartup:         a.Startup,
		OnDomReady:        a.DomReady,
		OnShutdown:        a.Shutdown,
		Bind: []interface{}{
			app.GetBindings(a),
		},
		// Windows platform specific options
		Windows: &windows.Options{
			WebviewIsTransparent: false,
			WindowIsTranslucent:  false,
			DisableWindowIcon:    false,
		},
	})

	return err
}

func main() {
	err := startup()
	if err != nil {
		msgBuilder := dialog.Message("Fatal error starting Aetherometer: %s", err)
		msgBuilder.Error()
	}
}
