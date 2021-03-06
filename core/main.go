package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"net/http/pprof"
	"os"
	"os/signal"
	"runtime/debug"
	"syscall"
	"time"

	"github.com/99designs/gqlgen/graphql/handler/lru"
	"github.com/BurntSushi/toml"
	"github.com/ff14wed/aetherometer/core/adapter"
	"github.com/ff14wed/aetherometer/core/server/handlers"
	"github.com/ff14wed/aetherometer/core/stream"
	"github.com/gorilla/websocket"
	"go.uber.org/zap"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/handler/extension"
	"github.com/99designs/gqlgen/graphql/handler/transport"
	"github.com/ff14wed/aetherometer/core/config"
	"github.com/ff14wed/aetherometer/core/datasheet"
	"github.com/ff14wed/aetherometer/core/models"
	"github.com/ff14wed/aetherometer/core/server"
	"github.com/ff14wed/aetherometer/core/store"
	"github.com/ff14wed/aetherometer/core/store/update"
	"github.com/thejerf/suture"
)

func addDebugHandlers(srv *server.Server) {
	srv.AddHandler("/debug/pprof/", http.HandlerFunc(pprof.Index))
	srv.AddHandler("/debug/pprof/cmdline", http.HandlerFunc(pprof.Cmdline))
	srv.AddHandler("/debug/pprof/profile", http.HandlerFunc(pprof.Profile))
	srv.AddHandler("/debug/pprof/symbol", http.HandlerFunc(pprof.Symbol))
	srv.AddHandler("/debug/pprof/trace", http.HandlerFunc(pprof.Trace))
}

func main() {
	zapCfg := zap.NewDevelopmentConfig()
	zapCfg.DisableStacktrace = true
	zapCfg.DisableCaller = true
	zapCfg.OutputPaths = []string{"stdout"}
	zapCfg.ErrorOutputPaths = []string{"stdout"}
	logger, err := zapCfg.Build()
	if err != nil {
		log.Fatalf("can't initialize zap logger: %v\n", err)
	}
	defer func() {
		_ = logger.Sync()
	}()
	zap.ReplaceGlobals(logger)

	cfgPath := flag.String("c", "", "path to TOML config file")

	flag.Usage = func() {
		fmt.Fprintf(flag.CommandLine.Output(), "Usage: %s -c [config path]\n", os.Args[0])
		flag.PrintDefaults()
	}

	flag.Parse()
	if len(*cfgPath) == 0 {
		flag.Usage()
		logger.Fatal("Must provide path to config argument.")
	}

	var cfg config.Config
	_, err = toml.DecodeFile(*cfgPath, &cfg)
	if err != nil {
		logger.Fatal("Error reading config file", zap.Error(err))
	}
	err = cfg.Validate()
	if err != nil {
		logger.Fatal("Error validating config file", zap.Error(err))
	}

	collection := new(datasheet.Collection)
	err = collection.Populate(cfg.DataPath)
	if err != nil {
		logger.Fatal("Error populating data", zap.Error(err))
	}

	// Since loading datasheets takes up a lot of memory for some reason
	debug.FreeOSMemory()

	srv := server.New(cfg, logger)

	generator := update.NewGenerator(collection)

	topSupervisor := suture.New("main", suture.Spec{
		Log: func(line string) {
			logger.Named("supervisor").Info(line)
		},
	})

	topSupervisor.ServeBackground()

	storeProvider := store.NewProvider(logger)
	topSupervisor.Add(storeProvider)

	streamSupervisor := suture.New("stream-supervisor", suture.Spec{
		Log: func(line string) {
			logger.Named("stream-supervisor").Info(line)
		},
	})
	topSupervisor.Add(streamSupervisor)

	sm := stream.NewManager(
		generator,
		storeProvider.UpdatesChan(),
		streamSupervisor,
		stream.NewHandler,
		logger,
	)

	streamRequestHandler := func(streamID int, request []byte) (string, error) {
		b, err := sm.SendRequest(streamID, request)
		return string(b), err
	}

	topSupervisor.Add(sm)

	adapters, err := stream.BuildAdapterInventory(adapter.Inventory(), cfg, sm.StreamUp(), sm.StreamDown(), logger)
	if err != nil {
		logger.Fatal("Error creating adapter", zap.Error(err))
	}
	for _, adapter := range adapters {
		topSupervisor.Add(adapter)
	}

	authHandler, err := handlers.NewAuth(cfg, transport.GetInitPayload, logger)
	if err != nil {
		logger.Fatal("Error initializing Auth handler", zap.Error(err))
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
	mapHandler := authHandler.Handler(handlers.NewMapHandler("/map/", cfg, logger))

	addDebugHandlers(srv)

	srv.AddHandler("/playground", handlers.Playground("GraphQL playground", "/query"))
	srv.AddHandler("/query", queryHandler)
	srv.AddHandler("/map/", mapHandler)

	topSupervisor.Add(srv)

	signals := make(chan os.Signal, 32)
	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)

	sig := <-signals
	logger.Info("Received signal, shutting down...", zap.Stringer("signal", sig))

	topSupervisor.Stop()
}
