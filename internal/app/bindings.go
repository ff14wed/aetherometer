package app

import (
	"errors"
	"fmt"
	"strings"

	"github.com/ff14wed/aetherometer/core/config"
	"github.com/ff14wed/aetherometer/core/models"
	"github.com/ff14wed/aetherometer/core/server/handlers"
	"github.com/skratchdot/open-golang/open"
)

type Bindings struct {
	app *App
}

func GetBindings(app *App) *Bindings {
	return &Bindings{app}
}

func (b *Bindings) WaitForStartup() {
	<-b.app.ready
}

func (b *Bindings) GetVersion() string {
	return b.app.version
}

func (b *Bindings) GetAPIVersion() string {
	return models.AetherometerAPIVersion
}

func (b *Bindings) GetAPIURL() string {
	addr := b.app.srv.Address()
	if addr == nil {
		return ""
	}
	return fmt.Sprintf("http://localhost:%d/query", addr.Port)
}

func (b *Bindings) GetAppDirectory() string {
	dirPath, _ := GetAppDirectory()
	return dirPath
}

func (b *Bindings) OpenAppDirectory() {
	dirPath := b.GetAppDirectory()
	if dirPath != "" {
		open.Start(dirPath)
	}
}

type StreamInfo struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

func (b *Bindings) GetStreams() []StreamInfo {
	streams, err := b.app.storeProvider.Streams()
	if err != nil {
		return nil
	}
	var infos []StreamInfo
	for _, s := range streams {
		if s.EntitiesMap == nil {
			continue
		}
		if char, ok := s.EntitiesMap[s.CharacterID]; ok && char != nil {
			infos = append(infos, StreamInfo{
				ID:   s.ID,
				Name: char.Name,
			})
		} else {
			infos = append(infos, StreamInfo{
				ID:   s.ID,
				Name: fmt.Sprintf("Stream %d", s.ID),
			})
		}
	}
	return infos
}

func (b *Bindings) GetPlugins() map[string]handlers.PluginInfo {
	return b.app.authHandler.GetRegisteredPlugins()
}

func (b *Bindings) GetConfig() config.Config {
	return b.app.cfgProvider.Config()
}

func (b *Bindings) AddPlugin(name string, url string) error {
	name = strings.TrimSpace(name)
	if name == "" {
		return errors.New("plugin name must not be empty")
	}

	url = strings.TrimSpace(url)
	if url == "" {
		return errors.New("plugin URL must not be empty")
	}
	return b.app.cfgProvider.AddPlugin(name, url)
}

func (b *Bindings) RemovePlugin(name string) error {
	return b.app.cfgProvider.RemovePlugin(name)
}
