package main

import (
	"context"

	"github.com/euiko/go-fullstack-boilerplate/pkg/webapp"

	"github.com/euiko/go-fullstack-boilerplate/internal/service/hello"
)

func main() {
	app := webapp.NewApp("go-fullstack-boilerplate")
	app.Register(hello.NewService)
	app.Run(context.Background())
}
