package lb

import (
	"net/http"
	"net/http/httputil"
	"net/url"
	"sync"
	"time"
)

type Backend struct { //nolint:govet
	alive       bool
	u           *url.URL
	api         *httputil.ReverseProxy
	mux         sync.RWMutex
	healthCheck *HealthCheck
}

func (b *Backend) Interval() <-chan time.Time {
	return b.healthCheck.Interval.C
}

func NewBackend(u *url.URL, api *httputil.ReverseProxy, check *HealthCheck) *Backend {
	return &Backend{u: u, api: api, healthCheck: check}
}

func (b *Backend) Alive() bool {
	b.mux.RLock()
	defer b.mux.RUnlock()

	return b.alive
}

func (b *Backend) HealthCheck() error {
	err := b.healthCheck.HealthCheck()
	if err != nil {
		b.mux.Lock()
		b.alive = false
		b.mux.Unlock()

		return err
	}

	b.mux.Lock()
	b.alive = true
	b.mux.Unlock()

	return nil
}

func (b *Backend) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	b.api.ServeHTTP(w, r)
}
