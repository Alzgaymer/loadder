package config

import (
	"github.com/go-chi/chi/v5"
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

func (r *HTTPRequest) Router(router chi.Router) error {
	for path, algo := range r.GET {
		h, err := DefineAlgorithm(algo)
		if err != nil {
			return err
		}

		router.Get(path, h)
	}

	for path, algo := range r.HEAD {
		h, err := DefineAlgorithm(algo)
		if err != nil {
			return err
		}

		router.Head(path, h)
	}

	for path, algo := range r.POST {
		h, err := DefineAlgorithm(algo)
		if err != nil {
			return err
		}

		router.Post(path, h)
	}

	for path, algo := range r.PUT {
		h, err := DefineAlgorithm(algo)
		if err != nil {
			return err
		}

		router.Put(path, h)
	}

	for path, algo := range r.PATCH {
		h, err := DefineAlgorithm(algo)
		if err != nil {
			return err
		}

		router.Patch(path, h)
	}

	for path, algo := range r.DELETE {
		h, err := DefineAlgorithm(algo)
		if err != nil {
			return err
		}

		router.Delete(path, h)
	}

	for path, algo := range r.CONNECT {
		h, err := DefineAlgorithm(algo)
		if err != nil {
			return err
		}

		router.Connect(path, h)
	}

	for path, algo := range r.OPTIONS {
		h, err := DefineAlgorithm(algo)
		if err != nil {
			return err
		}

		router.Options(path, h)
	}

	for path, algo := range r.TRACE {
		h, err := DefineAlgorithm(algo)
		if err != nil {
			return err
		}

		router.Trace(path, h)
	}

	return nil
}
