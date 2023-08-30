package lb

import (
	"net/http"
)

type RoundRobinAlgorithm struct {
	liveCount int
	iterator  int
	backends  []Service
}

func NewRoundRobinAlgorithm(b ...Service) *RoundRobinAlgorithm {
	return &RoundRobinAlgorithm{backends: b}
}

func (rr *RoundRobinAlgorithm) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if rr.liveCount == 0 {
		return
	}

	for {
		if rr.iterator == len(rr.backends) {
			rr.iterator = 0
		}

		if rr.backends[rr.iterator].Alive() {
			rr.iterator++
			rr.backends[rr.iterator].ServeHTTP(w, r)
		}

		rr.iterator++
	}
}
