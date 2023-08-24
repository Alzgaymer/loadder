package config

import (
	"fmt"
	"github.com/go-chi/chi/v5"
)

type Config struct {
	Services         map[string]*Service `yaml:"services"`
	LoadBalancerPort string              `yaml:"-"`
}

func (c *Config) Router() (chi.Router, error) {
	router := chi.NewRouter()
	for serviceID, service := range c.Services {
		err := service.Requests.HTTP.Router(router)
		if err != nil {
			return nil, fmt.Errorf("service %s: %w", serviceID, err)
		}
	}

	return router, nil
}
