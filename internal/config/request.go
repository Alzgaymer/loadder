package config

import (
	"github.com/go-chi/chi/v5"
	"net/http"
)

type Requests struct {
	HTTP HTTPRequest `yaml:"http"`
}

type HTTPRequest struct {
	GET     map[string]string `yaml:"GET"`
	HEAD    map[string]string `yaml:"HEAD"`
	POST    map[string]string `yaml:"POST"`
	PUT     map[string]string `yaml:"PUT"`
	PATCH   map[string]string `yaml:"PATCH"`
	DELETE  map[string]string `yaml:"DELETE"`
	CONNECT map[string]string `yaml:"CONNECT"`
	OPTIONS map[string]string `yaml:"OPTIONS"`
	TRACE   map[string]string `yaml:"TRACE"`
}

func (r *HTTPRequest) Router() (chi.Router, error) {
	router := chi.NewRouter()

	for path, algo := range r.GET {
		h, err := DefineAlgorithm(algo)
		if err != nil {
			return nil, err
		}

		router.Get(path, h)
	}

	for path, algo := range r.HEAD {
		h, err := DefineAlgorithm(algo)
		if err != nil {
			return nil, err
		}

		router.Head(path, h)
	}

	for path, algo := range r.POST {
		h, err := DefineAlgorithm(algo)
		if err != nil {
			return nil, err
		}

		router.Post(path, h)
	}

	for path, algo := range r.PUT {
		h, err := DefineAlgorithm(algo)
		if err != nil {
			return nil, err
		}

		router.Put(path, h)
	}

	for path, algo := range r.PATCH {
		h, err := DefineAlgorithm(algo)
		if err != nil {
			return nil, err
		}

		router.Patch(path, h)
	}

	for path, algo := range r.DELETE {
		h, err := DefineAlgorithm(algo)
		if err != nil {
			return nil, err
		}

		router.Delete(path, h)
	}

	for path, algo := range r.CONNECT {
		h, err := DefineAlgorithm(algo)
		if err != nil {
			return nil, err
		}

		router.Connect(path, h)
	}

	for path, algo := range r.OPTIONS {
		h, err := DefineAlgorithm(algo)
		if err != nil {
			return nil, err
		}

		router.Options(path, h)
	}

	for path, algo := range r.TRACE {
		h, err := DefineAlgorithm(algo)
		if err != nil {
			return nil, err
		}

		router.Trace(path, h)
	}

	return router, nil
}

func DefineAlgorithm(algo string) (http.HandlerFunc, error) {
	switch algo {
	case RR:
	case WRR:
	case LC:
	case LRT:
	case LB:
	case H:
	default:
		return nil, ErrInvalidAlgorithm
	}

	return nil, nil
}
