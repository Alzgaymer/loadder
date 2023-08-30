package lb

import (
	"net/http"
)

type RoundRobinAlgorithm struct {
	backends []Service
	iterator int
}

func NewRoundRobinAlgorithm(b ...Service) *RoundRobinAlgorithm {
	return &RoundRobinAlgorithm{backends: b}
}

func (rr *RoundRobinAlgorithm) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if rr.iterator == len(rr.backends) {
		rr.iterator = 0
	}

	if rr.backends[rr.iterator].Alive() {
		rr.backends[rr.iterator].ServeHTTP(w, r)
	}

	rr.iterator++
}
