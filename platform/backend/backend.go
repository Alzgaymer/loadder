package backend

import (
	"errors"
	"loadder/internal/domain/algorithm"
	lb "loadder/internal/domain/load_balancer"
	"net/http"
	"net/http/httputil"
	"net/url"
	"sync"
)

type Backend struct { //nolint:govet
	Alive     bool
	api       *httputil.ReverseProxy
	mux       sync.RWMutex
	Algorithm algorithm.Algorithm
}

func NewBackend(api *httputil.ReverseProxy) *Backend {
	return &Backend{api: api}
}

type Backends map[string]lb.Backend

func (b Backends) Add(ports []string, backends ...lb.Backend) error {
	if len(ports) != len(backends) {
		return errors.New("length of ports and backends doesn't match")
	}

	for i, port := range ports {
		b[port] = backends[i]
	}

	return nil
}

func (b *Backend) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	b.api.ServeHTTP(w, r)
}

func ParseBackends(urls ...*url.URL) Backends {
	var apis Backends

	for _, u := range urls {
		_ = apis.Add([]string{u.Host}, NewBackend(httputil.NewSingleHostReverseProxy(u)))
	}

	return apis
}
