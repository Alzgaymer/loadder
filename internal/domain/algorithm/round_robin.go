package algorithm

import (
	"container/ring"
	lb "loadder/internal/domain/load_balancer"
	"net/http"
)

type RoundRobin struct {
	backends lb.Backends
	hosts    *ring.Ring
}

func (rr *RoundRobin) Set(b lb.Backends) {
	rr.backends = b
}

func (rr *RoundRobin) ServeHTTP(w http.ResponseWriter, r *http.Request) {

}
