package backend

import (
	"net"
	"net/http"
	"net/http/httputil"
	"net/url"
	"sync"
)

type Backends []*Backend

type Backend struct { //nolint:govet
	alive bool
	host  string
	api   *httputil.ReverseProxy
	mux   sync.RWMutex
}

func NewBackend(host string, api *httputil.ReverseProxy) *Backend {
	return &Backend{host: host, api: api}
}

func (b *Backend) Alive() bool {
	b.mux.RLock()
	defer b.mux.RUnlock()

	return b.alive
}
func (b *Backend) IsAlive() bool {
	b.mux.Lock()
	defer b.mux.Unlock()

	dial, err := net.Dial("tcp", b.host)
	if err != nil {
		return false
	}
	defer dial.Close()

	b.alive = true

	return b.alive
}

func (b *Backend) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	b.api.ServeHTTP(w, r)
}

func CreateBackends(u ...*url.URL) Backends {
	backs := make(Backends, len(u))

	for i, url := range u {
		backs[i] = NewBackend(url.Host, httputil.NewSingleHostReverseProxy(url))
	}

	return backs
}
