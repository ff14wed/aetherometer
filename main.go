package main

import (
	"context"
	"embed"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"

	"github.com/apenwarr/fixconsole"
	"github.com/ff14wed/aetherometer/core/config"
	"github.com/ff14wed/aetherometer/internal/app"
	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/logger"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/windows"
	"go.uber.org/zap"
)

//go:embed frontend/dist
var assets embed.FS

var Version = "development"

func main() {
	flag.Usage = func() {
		fixconsole.FixConsoleIfNeeded()
		fmt.Fprintf(flag.CommandLine.Output(), "Usage of %s:\n", os.Args[0])
		flag.PrintDefaults()
	}

	dirPath, err := app.GetCurrentDirectory()
	if err != nil {
		panic(err)
	}

	cfgPath := flag.String("c", filepath.Join(dirPath, "config.toml"), "optional path to TOML config file")

	headless := flag.Bool("headless", false, "run Aetherometer in headless mode.")

	helpFlag := flag.Bool("h", false, "displays usage information")

	versionFlag := flag.Bool("v", false, "displays version information")

	flag.Parse()

	if *helpFlag {
		flag.Usage()
		return
	}

	if *versionFlag {
		fixconsole.FixConsoleIfNeeded()
		fmt.Println("Aetherometer version:", Version)
		return
	}

	var outputLogPath string = "aetherometer.log"
	if *headless {
		outputLogPath = "stdout"
		fixconsole.FixConsoleIfNeeded()
	}

	zapCfg := zap.NewDevelopmentConfig()
	zapCfg.DisableStacktrace = true
	zapCfg.DisableCaller = true
	zapCfg.OutputPaths = []string{outputLogPath}
	zapCfg.ErrorOutputPaths = []string{outputLogPath}
	zapLogger, err := zapCfg.Build()
	if err != nil {
		log.Fatalf("can't initialize zap logger: %v\n", err)
	}
	defer func() {
		_ = zapLogger.Sync()
	}()
	zap.ReplaceGlobals(zapLogger)

	if len(*cfgPath) == 0 {
		zapLogger.Fatal("Config path cannot be empty")
	}

	defaultCfg, err := app.DefaultConfig()
	if err != nil {
		zapLogger.Fatal("Error setting up default config", zap.Error(err))
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
		return
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

	if err != nil {
		log.Fatal(err)
	}
}
