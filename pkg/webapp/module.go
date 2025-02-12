package webapp

import (
	"context"

	"github.com/go-chi/chi/v5"
)

type (
	Module interface {
		Init(ctx context.Context) error
		Close() error
	}

	Service interface {
		Route(router chi.Router)
	}
)
