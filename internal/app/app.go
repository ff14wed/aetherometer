package app

import (
	"context"
	"fmt"
	"net/http"
	"net/http/pprof"
	"runtime/debug"
	"time"

	"github.com/ff14wed/aetherometer/core/adapter"
	"github.com/ff14wed/aetherometer/core/config"
	"github.com/ff14wed/aetherometer/core/datasheet"
	"github.com/ff14wed/aetherometer/core/models"
	"github.com/ff14wed/aetherometer/core/server"
	"github.com/ff14wed/aetherometer/core/server/handlers"
	"github.com/ff14wed/aetherometer/core/store"
	"github.com/ff14wed/aetherometer/core/store/update"
	"github.com/ff14wed/aetherometer/core/stream"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/handler/extension"
	"github.com/99designs/gqlgen/graphql/handler/lru"
	"github.com/99designs/gqlgen/graphql/handler/transport"
	"github.com/gorilla/websocket"
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
	cfgPath string
	version string
	logger  *zap.Logger

	ctx         context.Context
	cfgProvider *config.Provider

	appSupervisor    *suture.Supervisor
	streamSupervisor *suture.Supervisor

	collection     *datasheet.Collection
	srv            *server.Server
	storeProvider  *store.Provider
	authHandler    *handlers.Auth
	streamManager  *stream.Manager
	streamAdapters map[string]stream.Adapter

	ready chan struct{}
}

// NewApp creates a new App application struct
func NewApp(cfgPath string, version string, logger *zap.Logger) *App {
	return &App{
		cfgPath: cfgPath,
		version: version,
		logger:  logger,
		ready:   make(chan struct{}),
	}
}

// Initialize is called before the application is started
func (b *App) Initialize() error {
	defaultCfg, err := defaultConfig()
	if err != nil {
		return fmt.Errorf("could not setup default config: %v", err)
	}

	if err = CheckUpdates(b.cfgPath, b.logger); err != nil {
		return fmt.Errorf("error updating: %s", err)
	}

	b.cfgProvider = config.NewProvider(b.cfgPath, defaultCfg, b.logger)

	err = b.cfgProvider.EnsureConfigFile()
	if err != nil {
		return fmt.Errorf("config error: %s", err)
	}

	b.appSupervisor = suture.New("main", suture.Spec{
		Log: func(line string) {
			b.logger.Named("supervisor").Info(line)
		},
	})

	b.collection = new(datasheet.Collection)
	err = b.reloadDatasheets(b.collection)
	if err != nil {
		return fmt.Errorf("could not populate data from datasheets: %s", err)
	}

	generator := update.NewGenerator(b.collection)

	b.storeProvider = store.NewProvider(b.logger)

	b.authHandler, err = handlers.NewAuth(b.cfgProvider, b.logger)
	if err != nil {
		return fmt.Errorf("could not initialize Auth handler: %s", err)
	}

	err = b.authHandler.RefreshConfig()
	if err != nil {
		return fmt.Errorf("config error: %s", err)
	}

	b.streamSupervisor = suture.New("stream-supervisor", suture.Spec{
		Log: func(line string) {
			b.logger.Named("stream-supervisor").Info(line)
		},
	})

	b.streamManager = stream.NewManager(
		generator,
		b.storeProvider.UpdatesChan(),
		b.streamSupervisor,
		stream.NewHandler,
		b.logger,
	)

	cfg := b.cfgProvider.Config()
	b.streamAdapters, err = stream.BuildAdapterInventory(
		adapter.Inventory(),
		cfg,
		b.streamManager.StreamUp(),
		b.streamManager.StreamDown(),
		b.logger,
	)
	if err != nil {
		return fmt.Errorf("could not initialize adapters: %s", err)
	}

	b.srv = b.initializeServer()

	return nil
}

// Startup is called at application startup
func (b *App) Startup(ctx context.Context) {
	b.ctx = ctx

	b.appSupervisor.ServeBackground()
	b.appSupervisor.Add(b.cfgProvider)
	b.cfgProvider.WaitUntilReady()

	b.appSupervisor.Add(b.storeProvider)

	appEventWatcher := NewEventWatcher(
		b.storeProvider.StreamEventSource(),
		b.cfgProvider,
		b.authHandler,
		ctx,
		b.logger,
	)
	b.appSupervisor.Add(appEventWatcher)

	b.appSupervisor.Add(b.streamSupervisor)

	b.appSupervisor.Add(b.streamManager)

	for _, adapter := range b.streamAdapters {
		b.appSupervisor.Add(adapter)
	}
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

// reloadDatasheets reloads datasheets from the filepath
func (b *App) reloadDatasheets(collection *datasheet.Collection) error {
	defer func() {
		// Since loading datasheets takes up a lot of memory for some reason
		debug.FreeOSMemory()
	}()
	cfg := b.cfgProvider.Config()
	err := collection.Populate(cfg.Sources.DataPath)
	if err != nil {
		return err
	}
	return nil
}

func (b *App) initializeServer() *server.Server {
	cfg := b.cfgProvider.Config()
	srv := server.New(cfg, b.logger)

	streamRequestHandler := func(streamID int, request []byte) (string, error) {
		b, err := b.streamManager.SendRequest(streamID, request)
		return string(b), err
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

	addDebugHandlers(srv)

	srv.AddHandler("/playground", handlers.Playground("GraphQL playground", "/query"))
	srv.AddHandler("/query", queryHandler)
	srv.AddHandler("/map/", mapHandler)

	return srv
}
