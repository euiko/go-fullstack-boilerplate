//go:build embed

package main

import (
	"embed"

	"github.com/euiko/go-fullstack-boilerplate/internal/core/webapp"
)

//go:embed ui/dist
var static embed.FS

func init() {
	// inject static files into webapp
	webapp.StaticFS = static
}
