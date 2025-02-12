package webapp

import (
	"context"
	"net/http"
	"os"
	"path"
	"time"

	"github.com/euiko/go-fullstack-boilerplate/pkg/log"
	"github.com/euiko/go-fullstack-boilerplate/pkg/signal"
	"github.com/go-chi/chi/v5"
)

type (
	App struct {
		name     string
		settings Settings

		registry           registry
		modules            []Module
		defaultMiddlewares []Middleware
	}

	Middleware func(http.Handler) http.Handler

	Option func(*App)

	registry []func(*Settings) Module
)

func WithDefaultMiddlewares(middlewares ...Middleware) Option {
	return func(a *App) {
		a.defaultMiddlewares = middlewares
	}
}

func NewApp(name string, opts ...Option) *App {
	app := &App{
		name:               name,
		settings:           loadSettings(name),
		modules:            []Module{},
		defaultMiddlewares: []Middleware{},
	}

	// apply options
	for _, opt := range opts {
		opt(app)
	}

	return app
}

// Register a module factory function to the app
func (a *App) Register(f func(*Settings) Module) {
	a.registry = append(a.registry, f)
}

// Run the app
func (a *App) Run(ctx context.Context) error {
	// initialize logger
	initializeLogger(a.settings.Log)

	// create and initialize modules
	log.Info("initializing modules...")
	a.modules = make([]Module, len(a.registry))
	for i, factory := range a.registry {
		a.modules[i] = factory(&a.settings)
		a.modules[i].Init(ctx)
	}

	// create and run server
	log.Info("starting the server...")
	server := a.createServer()
	go server.ListenAndServe()
	log.Info("server started", log.WithField("addr", a.settings.Server.Addr))

	// wait for signal to be done
	signal := signal.NewSignalNotifier()
	signal.OnSignal(func(ctx context.Context, sig os.Signal) bool {
		return true // exit on receiving any signal
	})
	signal.Wait(ctx)

	// close all modules
	log.Info("closing modules...")
	for _, module := range a.modules {
		module.Close()
	}

	// close the server within 120s
	log.Info("closing the server...")
	ctx, cancel := context.WithTimeout(context.Background(), 120*time.Second)
	defer cancel() // ensure no context leak on graceful shutdown
	return server.Shutdown(ctx)
}

// internal createServer function
func (a *App) createServer() http.Server {
	// use chi as the router
	router := chi.NewRouter()

	// use default middlewares
	for _, middleware := range a.defaultMiddlewares {
		router.Use(middleware)
	}

	// register static routes
	if a.settings.StaticServer.Enabled {
		createStaticRoutes(&a.settings.StaticServer, router)
	}

	// register routes
	for _, module := range a.modules {
		// register routes
		if service, ok := module.(Service); ok {
			service.Route(router)
		}
	}

	// creates http server
	// TODO: add https support
	return http.Server{
		Addr:         a.settings.Server.Addr,
		Handler:      router,
		ReadTimeout:  a.settings.Server.ReadTimeout,
		WriteTimeout: a.settings.Server.WriteTimeout,
		IdleTimeout:  a.settings.Server.IdleTimeout,
	}
}

func loadSettings(name string) Settings {
	// default settings
	settings := Settings{
		Log: LogSettings{
			Level: "info",
		},
		Server: ServerSettings{
			Addr:         ":8080",
			ReadTimeout:  10 * time.Second,
			WriteTimeout: 10 * time.Second,
			IdleTimeout:  0,
			// TODO: add https support
		},
		StaticServer: StaticServerSettings{
			Enabled:    true,
			Path:       "/",
			IndexPath:  "index.html",
			AssetsPath: "assets",
		},
		config: nil,
	}

	// use viper as config provider
	viperOpts := []ViperOptions{}
	homeDir := os.Getenv("HOME")
	if homeDir != "" {
		viperOpts = append(viperOpts,
			ViperPaths(homeDir),
			ViperPaths(path.Join(homeDir, ".config", name)),
		)
	}

	// load settings
	config := NewViper(name, viperOpts...)
	config.Scan(&settings)

	// set the setting's config
	settings.config = config
	return settings
}

func initializeLogger(settings LogSettings) {
	// use LogrusLogger as default logger
	level := log.ParseLevel(settings.Level)
	log.SetDefault(log.NewLogrusLogger(level))
}
