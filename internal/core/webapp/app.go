package webapp

import (
	"context"
	"net/http"
	"os"
	"time"

	"github.com/euiko/go-fullstack-boilerplate/internal/core/log"
	"github.com/euiko/go-fullstack-boilerplate/internal/core/signal"
	"github.com/go-chi/chi/v5"
	"github.com/spf13/cobra"
)

type (
	App struct {
		name      string
		shortName string
		settings  *Settings

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

func NewApp(name string, shortName string, opts ...Option) *App {
	app := App{
		name:               name,
		shortName:          shortName,
		settings:           nil,
		modules:            []Module{},
		defaultMiddlewares: []Middleware{},
	}

	// apply options
	for _, opt := range opts {
		opt(&app)
	}

	return &app
}

// Register a module factory function to the app
func (a *App) Register(f func(*Settings) Module) {
	a.registry = append(a.registry, f)
}

// Run the app
func (a *App) Run(ctx context.Context) error {
	// load settings
	settings := loadSettings(a.name, a.shortName)
	a.settings = &settings

	// initialize logger
	initializeLogger(settings.Log)

	// create and initialize modules
	log.Trace("initializing modules...")
	a.modules = make([]Module, len(a.registry))
	for i, factory := range a.registry {
		a.modules[i] = factory(&settings)
		a.modules[i].Init(ctx)
	}

	rootCmd := a.initializeCli()
	if err := rootCmd.ExecuteContext(ctx); err != nil {
		return err
	}

	// close all modules
	log.Trace("closing modules...")
	for _, module := range a.modules {
		module.Close()
	}
	return nil
}

func (a *App) Start(ctx context.Context) error {
	// create and initialize server
	log.Info("starting the server...", log.WithField("addr", a.settings.Server.Addr))
	server := a.createServer()
	if err := initializeDB(a.settings.DB); err != nil {
		return err
	}

	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			panic(err)
		}
	}()

	// wait for signal to be done
	signal := signal.NewSignalNotifier()
	signal.OnSignal(func(ctx context.Context, sig os.Signal) bool {
		return true // exit on receiving any signal
	})
	signal.Wait(ctx)

	// close the server within 120s
	log.Info("closing the server...")
	defer closeDB()
	ctx, cancel := context.WithTimeout(context.Background(), 120*time.Second)
	defer cancel() // ensure no context leak on graceful shutdown
	return server.Shutdown(ctx)
}

func (a *App) initializeCli() *cobra.Command {
	rootCmd := cobra.Command{
		Use: a.name,
	}

	for m := range a.modules {
		if cli, ok := a.modules[m].(CLI); ok {
			cli.Command(&rootCmd)
		}
	}

	return &rootCmd
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
		createStaticRoutes(router)
	}

	// register routes
	router.Route("/api", func(r chi.Router) {
		for _, module := range a.modules {
			// register routes
			if service, ok := module.(APIService); ok {
				service.APIRoute(r)
			}
		}
	})

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

func initializeLogger(settings LogSettings) {
	// use LogrusLogger as default logger
	level := log.ParseLevel(settings.Level)
	log.SetDefault(log.NewLogrusLogger(level))
}
