package lb

import (
	"net/http/httputil"
	"net/url"
	"sync"
)

type Backend struct { //nolint:govet
	URL   *url.URL
	Alive bool
	mux   sync.RWMutex
	proxy *httputil.ReverseProxy
}
