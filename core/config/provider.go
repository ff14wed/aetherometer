package config

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"sync"

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
	NotifyHub *hub.NotifyHub

	logger *zap.Logger

	configFile string

	savedConfig Config
	configLock  sync.RWMutex

	internalWriteEvent chan struct{}

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
		NotifyHub: hub.NewNotifyHub(10),

		logger: logger.Named("config-provider"),

		configFile: configFile,

		savedConfig: defaultConfig,
		configLock:  sync.RWMutex{},

		internalWriteEvent: make(chan struct{}, 10),

		ready:    make(chan struct{}),
		stop:     make(chan struct{}),
		stopDone: make(chan struct{}),
	}
}

// Serve runs the main loop for the provider. It updates the saved
// configuration in response to file changes.
func (p *Provider) Serve() {
	if ok := p.ensureConfigFile(); !ok {
		return
	}
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		p.logger.Error("Unable to setup FS watcher", zap.Error(err))
		return
	}
	defer watcher.Close()

	err = watcher.Add(p.configFile)
	if err != nil {
		p.logger.Error("Unable to setup FS watcher", zap.Error(err))
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
				p.logger.Info("Detected config file change")
				err = p.readConfig()
				if err != nil {
					p.logger.Error("Unable to read config file", zap.Error(err))
					break
				}
				p.logger.Info("Successfully applied config change")
			}
		case err, ok := <-watcher.Errors:
			if !ok {
				return
			}
			p.logger.Error("FS watcher error", zap.Error(err))
		case <-p.internalWriteEvent:
			// If we are writing to disk, consume the next watcher event
			ok := consumeNextWriteEvent(watcher.Events)
			if !ok {
				return
			}
		case <-p.stop:
			p.logger.Info("Stopping...")
			return
		}
	}
}

func consumeNextWriteEvent(fsEvents chan fsnotify.Event) (ok bool) {
	for {
		event, ok := <-fsEvents
		if !ok {
			return false
		}
		if event.Op&fsnotify.Write == fsnotify.Write {
			return true
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

func (p *Provider) ensureConfigFile() (ok bool) {
	if _, err := os.Stat(p.configFile); errors.Is(err, os.ErrNotExist) {
		// Config file doesn't exist, so write the config to disk first
		p.logger.Info("Writing default config")
		writeErr := p.writeConfig()
		if writeErr != nil {
			p.logger.Error("Unable to write config file", zap.Error(err))
			return false
		}
	} else if err != nil {
		p.logger.Error("Unable to check for config file", zap.Error(err))
		return false
	} else {
		readErr := p.readConfig()
		if readErr != nil {
			p.logger.Error("Unable to read config file", zap.Error(readErr))
			return false
		}
	}
	return true
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
	p.NotifyHub.Broadcast()
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

func (p *Provider) sendInternalWriteEvent() {
	// Don't worry if the channel is blocked
	select {
	case p.internalWriteEvent <- struct{}{}:
	default:
	}
}

// AddPlugin adds the given plugin to the configuration
// It errors if the plugin name already exists.
func (p *Provider) AddPlugin(name string, pluginURL string) error {
	err := func() error {
		p.configLock.Lock()
		defer p.configLock.Unlock()

		if _, ok := p.savedConfig.Plugins[name]; ok {
			return fmt.Errorf(`plugin "%s" already exists`, name)
		}

		cfg := p.savedConfig
		if cfg.Plugins == nil {
			cfg.Plugins = make(map[string]string)
		} else {
			cfg.Plugins = copyMap(cfg.Plugins)
		}
		cfg.Plugins[name] = pluginURL
		p.updateConfig(cfg)

		return nil
	}()

	if err != nil {
		return err
	}
	p.sendInternalWriteEvent()
	return p.writeConfig()
}

// RemovePlugin removes the plugin with the given name from the configuration.
// If the plugin doesn't exist, it is a no-op.
func (p *Provider) RemovePlugin(name string) error {
	func() {
		p.configLock.Lock()
		defer p.configLock.Unlock()

		if p.savedConfig.Plugins != nil {
			if _, ok := p.savedConfig.Plugins[name]; ok {
				cfg := p.savedConfig
				cfg.Plugins = copyMap(cfg.Plugins)
				delete(cfg.Plugins, name)
				p.updateConfig(cfg)
			}
		}
	}()

	p.sendInternalWriteEvent()
	return p.writeConfig()
}

// SetDisableAuth sets the value of the DisableAuth field in the configuration.
func (p *Provider) SetDisableAuth(disableAuth bool) error {
	func() {
		p.configLock.Lock()
		defer p.configLock.Unlock()

		p.savedConfig.DisableAuth = disableAuth
		p.NotifyHub.Broadcast()
	}()

	p.sendInternalWriteEvent()
	return p.writeConfig()

}

func copyMap(src map[string]string) map[string]string {
	target := make(map[string]string)
	for key, value := range src {
		target[key] = value
	}
	return target
}
