//go:build embed

package webapp

import (
	"embed"
	"io/fs"
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"
)

const indexPath = "index.html"

var StaticFS embed.FS

type subFS struct {
	path string
	fs   fs.FS
}

// Open implements fs.FS for subFS that prepends the path to the name
func (s subFS) Open(name string) (fs.File, error) {
	return s.fs.Open(s.path + "/" + name)
}

func createStaticRoutes(r chi.Router) {
	staticFs := newSubFS(StaticFS, "ui/dist")

	// serve other files from the embedded StaticFS
	entries, _ := StaticFS.ReadDir("ui/dist")
	for _, entry := range entries {
		// skip index.html
		if entry.Name() == indexPath {
			continue
		}

		// use absolute route and staticfs for files
		fs := staticFs
		route := "/" + entry.Name()

		// use wildcard route and subfs for directories
		if entry.IsDir() {
			route = "/" + entry.Name() + "/*"
			fs = newSubFS(staticFs, entry.Name())
		}

		// register route
		r.Get(route, func(w http.ResponseWriter, r *http.Request) {
			httpFs := http.FileServer(http.FS(fs))

			// trim the directory prefix for directories
			if entry.IsDir() {
				prefix := strings.TrimSuffix(route, "/*")
				httpFs = http.StripPrefix(prefix, httpFs)
			}

			// serve the httpFs
			httpFs.ServeHTTP(w, r)
		})
	}

	// serve index.html from embedded static
	r.Get("/*", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFileFS(w, r, staticFs, indexPath)
	})
}

func newSubFS(fs fs.FS, path string) fs.FS {
	return &subFS{
		path: path,
		fs:   fs,
	}
}
