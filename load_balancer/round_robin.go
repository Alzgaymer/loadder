package lb

import "net/http"

type RoundRobin struct {
	services []*Service
	iterator int
}

func (rr *RoundRobin) NextAlive() *Service {

}

func (rr *RoundRobin) ServeHTTP(w http.ResponseWriter, r *http.Request) {

}
