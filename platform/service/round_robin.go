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
		r.Next()
	}

	return &RoundRobinAlgorithm{store: r}
}

func (rr *RoundRobinAlgorithm) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	defer rr.store.Next()
	b := rr.store.Value.(*backend.Backend)
	b.ServeHTTP(w, r)
}
