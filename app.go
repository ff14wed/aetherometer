package main

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
	ctx    context.Context
	cfg    config.Config
	logger *zap.Logger

	appSupervisor *suture.Supervisor
}

// NewApp creates a new App application struct
func NewApp(cfg config.Config, logger *zap.Logger) *App {
	return &App{
		cfg:    cfg,
		logger: logger,
	}
}

// startup is called at application startup
func (b *App) startup(ctx context.Context) {
	// Perform your setup here
	b.ctx = ctx

	collection := new(datasheet.Collection)
	err := collection.Populate(b.cfg.DataPath)
	if err != nil {
		b.logger.Fatal("Error populating data", zap.Error(err))
	}

	// Since loading datasheets takes up a lot of memory for some reason
	debug.FreeOSMemory()

	srv := server.New(b.cfg, b.logger)

	generator := update.NewGenerator(collection)

	b.appSupervisor = suture.New("main", suture.Spec{
		Log: func(line string) {
			b.logger.Named("supervisor").Info(line)
		},
	})

	b.appSupervisor.ServeBackground()

	storeProvider := store.NewProvider(b.logger)
	b.appSupervisor.Add(storeProvider)

	streamSupervisor := suture.New("stream-supervisor", suture.Spec{
		Log: func(line string) {
			b.logger.Named("stream-supervisor").Info(line)
		},
	})
	b.appSupervisor.Add(streamSupervisor)

	sm := stream.NewManager(
		generator,
		storeProvider.UpdatesChan(),
		streamSupervisor,
		stream.NewHandler,
		b.logger,
	)

	streamRequestHandler := func(streamID int, request []byte) (string, error) {
		b, err := sm.SendRequest(streamID, request)
		return string(b), err
	}

	b.appSupervisor.Add(sm)

	adapters, err := stream.BuildAdapterInventory(adapter.Inventory(), b.cfg, sm.StreamUp(), sm.StreamDown(), b.logger)
	if err != nil {
		b.logger.Fatal("Error creating adapter", zap.Error(err))
	}
	for _, adapter := range adapters {
		b.appSupervisor.Add(adapter)
	}

	authHandler, err := handlers.NewAuth(b.cfg, transport.GetInitPayload, b.logger)
	if err != nil {
		b.logger.Fatal("Error initializing Auth handler", zap.Error(err))
	}

	queryResolver := models.NewResolver(storeProvider, authHandler, streamRequestHandler)

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

	queryHandler := authHandler.Handler(gqlServer)
	mapHandler := authHandler.Handler(handlers.NewMapHandler("/map/", b.cfg, b.logger))

	addDebugHandlers(srv)

	srv.AddHandler("/playground", handlers.Playground("GraphQL playground", "/query"))
	srv.AddHandler("/query", queryHandler)
	srv.AddHandler("/map/", mapHandler)

	b.appSupervisor.Add(srv)
}

// domReady is called after the front-end dom has been loaded
func (b *App) domReady(ctx context.Context) {
	// Add your action here
}

// shutdown is called at application termination
func (b *App) shutdown(ctx context.Context) {
	// Perform your teardown here

	b.appSupervisor.Stop()
}

// Greet returns a greeting for the given name
func (b *App) Greet(name string) string {
	return fmt.Sprintf("Hello %s!", name)
}
