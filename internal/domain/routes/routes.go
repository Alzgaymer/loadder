package routes

import (
	"fmt"
	"github.com/go-chi/chi/v5"
	"loadder/internal/config"
	"loadder/internal/domain/algorithms"
	"net/http"
)

func ExtractRoutes(c *config.Config) (chi.Router, error) {
	router := chi.NewRouter()
	for serviceID, service := range c.Services {
		err := extractRoutesFromService(service, router)
		if err != nil {
			return nil, fmt.Errorf("service %s: %w", serviceID, err)
		}
	}

	return router, nil
}

func extractRoutesFromService(service *config.Service, router chi.Router) error {
	r := service.Requests.HTTP
	builder := algorithms.NewBuilder(http.DefaultClient, service.Address, service.Ports)

	methods := []map[string]string{
		r.GET,
		r.HEAD,
		r.POST,
		r.PUT,
		r.PATCH,
		r.DELETE,
		r.TRACE,
		r.OPTIONS,
		r.CONNECT,
	}

	for _, method := range methods {
		for path, algo := range method {
			h, err := builder.DefineAlgorithm(algo)
			if err != nil {
				return err
			}

			router.Get(path, h.Handler())
		}
	}

	return nil
}
