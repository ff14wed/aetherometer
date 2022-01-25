package main

import (
	"context"
	"embed"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/BurntSushi/toml"
	"github.com/apenwarr/fixconsole"
	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/logger"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/windows"
	"go.uber.org/zap"
)

//go:embed frontend/dist
var assets embed.FS

//go:embed build/appicon.png
var icon []byte

func main() {
	flag.Usage = func() {
		fixconsole.FixConsoleIfNeeded()
		fmt.Fprintf(flag.CommandLine.Output(), "Usage of %s:\n", os.Args[0])
		flag.PrintDefaults()
	}

	cfgPath := flag.String("c", "", "optional path to TOML config file")

	headless := flag.Bool("headless", false, "run Aetherometer in headless mode.")

	helpFlag := flag.Bool("h", false, "displays usage information")

	flag.Parse()

	if *helpFlag {
		flag.Usage()
		return
	}

	loggingOutput := "aetherometer.log"
	if *headless {
		loggingOutput = "stdout"
		fixconsole.FixConsoleIfNeeded()
	}

	zapCfg := zap.NewDevelopmentConfig()
	zapCfg.DisableStacktrace = true
	zapCfg.DisableCaller = true
	zapCfg.OutputPaths = []string{loggingOutput}
	zapCfg.ErrorOutputPaths = []string{loggingOutput}
	zapLogger, err := zapCfg.Build()
	if err != nil {
		log.Fatalf("can't initialize zap logger: %v\n", err)
	}
	defer func() {
		_ = zapLogger.Sync()
	}()
	zap.ReplaceGlobals(zapLogger)

	cfg, err := defaultConfig()
	if err != nil {
		zapLogger.Fatal("Error setting up config", zap.Error(err))
	}
	zapLogger.Info("Using config", zap.Any("config", cfg))

	if len(*cfgPath) != 0 {
		_, err = toml.DecodeFile(*cfgPath, &cfg)
		if err != nil {
			zapLogger.Fatal("Error reading config file", zap.Error(err))
		}
		err = cfg.Validate()
		if err != nil {
			zapLogger.Fatal("Error validating config file", zap.Error(err))
		}
	}

	app := NewApp(cfg, zapLogger)

	// Do not start the GUI in headless mode.
	if *headless {
		app.startup(context.Background())

		signals := make(chan os.Signal, 32)
		signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)

		sig := <-signals
		zapLogger.Info("Received signal, shutting down...", zap.Stringer("signal", sig))

		app.shutdown(context.Background())
		return
	}

	// Start run wails app if not headless mode
	err = wails.Run(&options.App{
		Title:             "aetherometer",
		Width:             720,
		Height:            570,
		MinWidth:          720,
		MinHeight:         570,
		MaxWidth:          1280,
		MaxHeight:         740,
		DisableResize:     false,
		Fullscreen:        false,
		Frameless:         false,
		StartHidden:       false,
		HideWindowOnClose: false,
		RGBA:              &options.RGBA{R: 33, G: 37, B: 43, A: 255},
		Assets:            assets,
		LogLevel:          logger.DEBUG,
		OnStartup:         app.startup,
		OnDomReady:        app.domReady,
		OnShutdown:        app.shutdown,
		Bind: []interface{}{
			app,
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
