package config

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"sync"
	"sync/atomic"
	"time"

	"github.com/BurntSushi/toml"
	"github.com/ff14wed/aetherometer/core/hub"
	"github.com/fsnotify/fsnotify"
	"go.uber.org/zap"
)

// Provider provides access to a constantly updating config file. It runs as a
// long running service that hot reloads the config in response to file changes.
// It returns a current snapshot of the config whenever it is polled for a
// config file.
type Provider struct {
	UpdateEvents *hub.NotifyHub[struct{}]
	ErrorEvents  *hub.NotifyHub[string]

	logger *zap.Logger

	configFile string

	savedConfig Config
	configLock  sync.RWMutex

	eventBatcher       *hub.EventBatcher
	internalWriteEvent int32

	ready    chan struct{}
	stop     chan struct{}
	stopDone chan struct{}
}

// NewProvider creates a new config provider.
func NewProvider(
	configFile string,
	defaultConfig Config,
	logger *zap.Logger,
) *Provider {
	return &Provider{
		UpdateEvents: hub.NewNotifyHub[struct{}](10),
		ErrorEvents:  hub.NewNotifyHub[string](10),

		logger: logger.Named("config-provider"),

		configFile: configFile,

		savedConfig: defaultConfig,
		configLock:  sync.RWMutex{},

		eventBatcher: hub.NewEventBatcher(20 * time.Millisecond),

		ready:    make(chan struct{}),
		stop:     make(chan struct{}),
		stopDone: make(chan struct{}),
	}
}

// EnsureConfigFile ensure the config file exists and validates the initial
// configuration passed to the provider
func (p *Provider) EnsureConfigFile() error {
	if _, err := os.Stat(p.configFile); errors.Is(err, os.ErrNotExist) {
		// Config file doesn't exist, so write the config to disk first
		p.logger.Info("Writing default config")
		writeErr := p.writeConfig()
		if writeErr != nil {
			return fmt.Errorf("unable to write config file: %s", err)
		}
	} else if err != nil {
		return fmt.Errorf("unable to check for config file: %s", err)
	} else {
		readErr := p.readConfig()
		if readErr != nil {
			return fmt.Errorf("unable to read config file: %s", readErr)
		}
	}
	return nil
}

// broadcastError emits a message to the log and all notifyhub subscribers
func (p *Provider) broadcastError(message string, err error) {
	p.logger.Error(message, zap.Error(err))
	p.ErrorEvents.Broadcast(fmt.Sprintf("%s: %s", message, err))
}

// Serve runs the main loop for the provider. It updates the saved
// configuration in response to file changes.
func (p *Provider) Serve() {
	if err := p.EnsureConfigFile(); err != nil {
		p.broadcastError("Error loading config file", err)
		return
	}
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		p.broadcastError("Unable to setup FS watcher", err)
		return
	}
	defer watcher.Close()

	err = watcher.Add(p.configFile)
	if err != nil {
		p.broadcastError("Unable to setup FS watcher", err)
		return
	}
	defer close(p.stopDone)

	p.logger.Info("Running")
	close(p.ready)

	for {
		select {
		case event, ok := <-watcher.Events:
			if !ok {
				return
			}
			if event.Op&fsnotify.Write == fsnotify.Write {
				p.eventBatcher.Notify()
			}
		case err, ok := <-watcher.Errors:
			if !ok {
				return
			}
			p.broadcastError("FS watcher error", err)
		case <-p.eventBatcher.BatchedEvents():
			// Ignore write event if it was an internal write
			swapped := atomic.CompareAndSwapInt32(&p.internalWriteEvent, 1, 0)
			if swapped {
				break
			}
			p.logger.Info("Detected config file change")
			err = p.readConfig()
			if err != nil {
				p.broadcastError("Unable to read config file", err)
				break
			}
			p.logger.Info("Successfully applied config change")
		case <-p.stop:
			p.logger.Info("Stopping...")
			return
		}
	}
}

// WaitUntilReady blocks until the config provider is up and running
func (p *Provider) WaitUntilReady() {
	<-p.ready
}

// Stop will shutdown this service and wait on it to stop before returning
func (p *Provider) Stop() {
	close(p.stop)
	<-p.stopDone
}

// Config returns a the stored configuration from the provider.
func (p *Provider) Config() Config {
	p.configLock.RLock()
	defer p.configLock.RUnlock()
	return p.savedConfig
}

// updateConfig overwrites the savedConfig with cfg
// It is expected to be used inside a critical section
func (p *Provider) updateConfig(cfg Config) {
	p.savedConfig = cfg
	p.UpdateEvents.Broadcast(struct{}{})
}

// readConfig reads the saved config file from disk
func (p *Provider) readConfig() error {
	cfg := Config{}
	_, err := toml.DecodeFile(p.configFile, &cfg)
	if err != nil {
		return err
	}
	err = cfg.Validate()
	if err != nil {
		return err
	}
	p.configLock.Lock()
	defer p.configLock.Unlock()
	p.updateConfig(cfg)

	return nil
}

// writeConfig writes the given config file to the disk
func (p *Provider) writeConfig() error {
	p.configLock.RLock()
	defer p.configLock.RUnlock()
	configBytes := bytes.Buffer{}
	encoder := toml.NewEncoder(&configBytes)
	err := encoder.Encode(p.savedConfig)
	if err != nil {
		return err
	}
	return os.WriteFile(p.configFile, configBytes.Bytes(), 0644)
}

// AddPlugin adds the given plugin to the configuration
// It errors if the plugin name already exists.
func (p *Provider) AddPlugin(name string, pluginURL string) error {
	return p.MutateConfig(func(cfg Config) (Config, error) {
		if _, ok := cfg.Plugins[name]; ok {
			return cfg, fmt.Errorf(`plugin "%s" already exists`, name)
		}

		if cfg.Plugins == nil {
			cfg.Plugins = make(map[string]string)
		} else {
			cfg.Plugins = copyMap(cfg.Plugins)
		}
		cfg.Plugins[name] = pluginURL

		return cfg, nil
	})
}

// RemovePlugin removes the plugin with the given name from the configuration.
// If the plugin doesn't exist, it is a no-op.
func (p *Provider) RemovePlugin(name string) error {
	return p.MutateConfig(func(cfg Config) (Config, error) {
		if cfg.Plugins != nil {
			if _, ok := cfg.Plugins[name]; ok {
				cfg.Plugins = copyMap(cfg.Plugins)
				delete(cfg.Plugins, name)
				return cfg, nil
			}
		}
		return cfg, nil
	})
}

// MutateConfig provides a callback that allows the caller to safely mutate
// the configuration
func (p *Provider) MutateConfig(mutate func(Config) (Config, error)) error {
	err := func() error {
		p.configLock.Lock()
		defer p.configLock.Unlock()

		newCfg, err := mutate(p.savedConfig)
		if err != nil {
			return err
		}

		p.updateConfig(newCfg)
		return nil
	}()

	if err != nil {
		return err
	}

	atomic.StoreInt32(&p.internalWriteEvent, 1)
	return p.writeConfig()

}

func copyMap(src map[string]string) map[string]string {
	target := make(map[string]string)
	for key, value := range src {
		target[key] = value
	}
	return target
}
