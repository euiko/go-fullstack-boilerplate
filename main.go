package main

import (
	"context"

	"github.com/euiko/go-fullstack-boilerplate/internal/cli"
	"github.com/euiko/go-fullstack-boilerplate/internal/core/webapp"

	"github.com/euiko/go-fullstack-boilerplate/internal/service/hello"
)

func main() {
	app := webapp.NewApp("go-fullstack-boilerplate", "WEBAPP")
	// CLI modules
	app.Register(cli.Server(app))
	app.Register(cli.Migration)

	// Service modules
	app.Register(hello.NewService)
	app.Run(context.Background())
}
