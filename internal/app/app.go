package app

import (
	"context"
	"github.com/go-chi/chi/v5"
	"loadder/internal/config"
)

type App struct {
	cfg *config.Config
	mux chi.Router
}

func Run(ctx context.Context) error {
	return nil
}

func Stop(ctx context.Context) error {
	return nil
}
