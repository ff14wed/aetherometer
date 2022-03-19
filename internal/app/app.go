package app

import (
	"context"
	"fmt"
	"net/http"
	"net/http/pprof"
	"runtime/debug"
	"time"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/handler/extension"
	"github.com/99designs/gqlgen/graphql/handler/lru"
	"github.com/99designs/gqlgen/graphql/handler/transport"
	"github.com/ff14wed/aetherometer/core/adapter"
	"github.com/ff14wed/aetherometer/core/config"
	"github.com/ff14wed/aetherometer/core/datasheet"
	"github.com/ff14wed/aetherometer/core/models"
	"github.com/ff14wed/aetherometer/core/server"
	"github.com/ff14wed/aetherometer/core/server/handlers"
	"github.com/ff14wed/aetherometer/core/store"
	"github.com/ff14wed/aetherometer/core/store/update"
	"github.com/ff14wed/aetherometer/core/stream"
	"github.com/gorilla/websocket"
	"github.com/skratchdot/open-golang/open"
	"github.com/thejerf/suture"
	"go.uber.org/zap"
)

func addDebugHandlers(srv *server.Server) {
	srv.AddHandler("/debug/pprof/", http.HandlerFunc(pprof.Index))
	srv.AddHandler("/debug/pprof/cmdline", http.HandlerFunc(pprof.Cmdline))
	srv.AddHandler("/debug/pprof/profile", http.HandlerFunc(pprof.Profile))
	srv.AddHandler("/debug/pprof/symbol", http.HandlerFunc(pprof.Symbol))
	srv.AddHandler("/debug/pprof/trace", http.HandlerFunc(pprof.Trace))
}

// App application struct
type App struct {
	ctx         context.Context
	cfgProvider *config.Provider
	logger      *zap.Logger

	version string

	appSupervisor *suture.Supervisor

	collection    *datasheet.Collection
	srv           *server.Server
	storeProvider *store.Provider
	authHandler   *handlers.Auth

	ready chan struct{}
}

// NewApp creates a new App application struct
func NewApp(cfgProvider *config.Provider, version string, logger *zap.Logger) *App {
	return &App{
		cfgProvider: cfgProvider,
		version:     version,
		logger:      logger,
		ready:       make(chan struct{}),
	}
}

// Startup is called at application startup
func (b *App) Startup(ctx context.Context) {
	b.logger.Info("====================================")
	b.logger.Info("Starting Aetherometer...")
	b.ctx = ctx

	b.appSupervisor = suture.New("main", suture.Spec{
		Log: func(line string) {
			b.logger.Named("supervisor").Info(line)
		},
	})

	b.appSupervisor.ServeBackground()
	b.appSupervisor.Add(b.cfgProvider)
	b.cfgProvider.WaitUntilReady()

	cfg := b.cfgProvider.Config()

	b.srv = server.New(cfg, b.logger)

	b.collection = new(datasheet.Collection)
	b.reloadDatasheets(b.collection)

	generator := update.NewGenerator(b.collection)

	b.storeProvider = store.NewProvider(b.logger)
	b.appSupervisor.Add(b.storeProvider)

	var err error

	b.authHandler, err = handlers.NewAuth(b.cfgProvider, b.logger)
	if err != nil {
		b.logger.Fatal("Error initializing Auth handler", zap.Error(err))
	}

	appEventWatcher := NewEventWatcher(
		b.storeProvider.StreamEventSource(),
		b.cfgProvider.NotifyHub,
		b.authHandler,
		ctx,
		b.logger,
	)
	b.appSupervisor.Add(appEventWatcher)

	streamSupervisor := suture.New("stream-supervisor", suture.Spec{
		Log: func(line string) {
			b.logger.Named("stream-supervisor").Info(line)
		},
	})
	b.appSupervisor.Add(streamSupervisor)

	sm := stream.NewManager(
		generator,
		b.storeProvider.UpdatesChan(),
		streamSupervisor,
		stream.NewHandler,
		b.logger,
	)

	streamRequestHandler := func(streamID int, request []byte) (string, error) {
		b, err := sm.SendRequest(streamID, request)
		return string(b), err
	}

	b.appSupervisor.Add(sm)

	adapters, err := stream.BuildAdapterInventory(adapter.Inventory(), cfg, sm.StreamUp(), sm.StreamDown(), b.logger)
	if err != nil {
		b.logger.Fatal("Error creating adapter", zap.Error(err))
	}
	for _, adapter := range adapters {
		b.appSupervisor.Add(adapter)
	}

	queryResolver := models.NewResolver(b.storeProvider, b.authHandler, streamRequestHandler)

	upgrader := websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
		EnableCompression: true,
	}

	gqlServer := handler.New(models.NewExecutableSchema(models.Config{
		Resolvers: queryResolver,
	}))

	gqlServer.AddTransport(transport.Websocket{
		Upgrader:              upgrader,
		InitFunc:              b.authHandler.WebsocketInitFunc,
		KeepAlivePingInterval: 10 * time.Second,
	})
	gqlServer.AddTransport(transport.Options{})
	gqlServer.AddTransport(transport.GET{})
	gqlServer.AddTransport(transport.POST{})
	gqlServer.AddTransport(transport.MultipartForm{})

	gqlServer.Use(extension.Introspection{})

	gqlServer.Use(extension.AutomaticPersistedQuery{
		Cache: lru.New(100),
	})

	queryHandler := b.authHandler.Handler(gqlServer)
	mapHandler := b.authHandler.Handler(handlers.NewMapHandler("/map/", cfg, b.logger))

	addDebugHandlers(b.srv)

	b.srv.AddHandler("/playground", handlers.Playground("GraphQL playground", "/query"))
	b.srv.AddHandler("/query", queryHandler)
	b.srv.AddHandler("/map/", mapHandler)

	b.appSupervisor.Add(b.srv)

	close(b.ready)
}

// DomReady is called after the front-end dom has been loaded
func (b *App) DomReady(ctx context.Context) {
	// Add your action here
}

// Shutdown is called at application termination
func (b *App) Shutdown(ctx context.Context) {
	// Perform your teardown here

	b.appSupervisor.Stop()
}

func (b *App) WaitForStartup() {
	<-b.ready
}

// reloadDatasheets reloads datasheets from the filepath
func (b *App) reloadDatasheets(collection *datasheet.Collection) {
	cfg := b.cfgProvider.Config()
	err := collection.Populate(cfg.Sources.DataPath)
	if err != nil {
		b.logger.Fatal("Error populating data", zap.Error(err))
	}

	// Since loading datasheets takes up a lot of memory for some reason
	debug.FreeOSMemory()
}

func (b *App) GetVersion() string {
	return b.version
}

func (b *App) GetAPIVersion() string {
	return models.AetherometerAPIVersion
}

func (b *App) GetAPIURL() string {
	addr := b.srv.Address()
	if addr == nil {
		return ""
	}
	return fmt.Sprintf("http://localhost:%d/query", addr.Port)
}

func (b *App) GetAppDirectory() string {
	dirPath, _ := GetCurrentDirectory()
	return dirPath
}

func (b *App) OpenAppDirectory() {
	dirPath := b.GetAppDirectory()
	if dirPath != "" {
		open.Start(dirPath)
	}
}

type StreamInfo struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

func (b *App) GetStreams() []StreamInfo {
	streams, err := b.storeProvider.Streams()
	if err != nil {
		return nil
	}
	var infos []StreamInfo
	for _, s := range streams {
		if char, ok := s.EntitiesMap[s.CharacterID]; ok {
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

func (b *App) GetPlugins() map[string]handlers.PluginInfo {
	return b.authHandler.GetRegisteredPlugins()
}

func (b *App) GetConfig() config.Config {
	return b.cfgProvider.Config()
}

func (b *App) AddPlugin(name string, url string) error {
	return b.cfgProvider.AddPlugin(name, url)
}

func (b *App) RemovePlugin(name string) error {
	return b.cfgProvider.RemovePlugin(name)
}
