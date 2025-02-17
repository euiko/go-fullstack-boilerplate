package hello

import (
	"context"
	"net/http"

	"github.com/euiko/go-fullstack-boilerplate/internal/core/webapp"
	"github.com/go-chi/chi/v5"
)

type Service struct {
}

func NewService(settings *webapp.Settings) webapp.Module {
	return &Service{}
}

func (svc *Service) Init(ctx context.Context) error {
	return nil
}

func (svc *Service) Close() error {
	return nil
}

func (svc *Service) Route(router chi.Router) {
	router.Get("/hello", svc.hello)
}

func (svc *Service) hello(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Hello, World!"))
}
