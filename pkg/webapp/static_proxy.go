//go:build !embed

package webapp

import (
	"net/http"
	"net/http/httputil"
	"net/url"

	"github.com/euiko/go-fullstack-boilerplate/pkg/log"
	"github.com/go-chi/chi/v5"
)

func createStaticRoutes(settings *StaticServerSettings, r chi.Router) {
	target := "http://localhost:5173"
	url, err := url.Parse(target)
	if err != nil {
		log.Fatal("invalid target", log.WithField("target", target))
	}

	r.HandleFunc("/*", func(w http.ResponseWriter, r *http.Request) {
		httputil.NewSingleHostReverseProxy(url).ServeHTTP(w, r)
	})
}
