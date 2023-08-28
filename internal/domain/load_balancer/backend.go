package lb

import (
	"net/http/httputil"
	"net/url"
	"sync"
)

type Backend struct {
	Alive bool
	api   *httputil.ReverseProxy
	mux   sync.RWMutex
}

func NewBackend(api *httputil.ReverseProxy) *Backend {
	return &Backend{api: api}
}

type Backends struct {
	apis map[*url.URL]*Backend
}

func NewBackends(apis map[*url.URL]*Backend) *Backends {
	return &Backends{apis: apis}
}

func ParseBackends(urls ...*url.URL) *Backends {
	apis := make(map[*url.URL]*Backend)

	for _, u := range urls {
		apis[u] = NewBackend(httputil.NewSingleHostReverseProxy(u))
	}

	return NewBackends(apis)
}
