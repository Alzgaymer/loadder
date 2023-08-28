package service

import (
	"container/ring"
	"loadder/platform/backend"
	"net/http"
)

type RoundRobinAlgorithm struct {
	store *ring.Ring
}

func NewRoundRobinAlgorithm(b ...*backend.Backend) *RoundRobinAlgorithm {
	r := ring.New(len(b))
	for _, back := range b {
		r.Value = back
		r = r.Next()
	}

	return &RoundRobinAlgorithm{store: r}
}

func (rr *RoundRobinAlgorithm) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	defer func() {
		rr.store = rr.store.Next()
	}()

	b, _ := rr.store.Value.(*backend.Backend)
	b.ServeHTTP(w, r)
}
