package lb

import "net/http"

type RoundRobin struct {
	services     []*Service
	currentIndex int
}

func (rr *RoundRobin) Add(services ...*Service) {
	rr.services = append(rr.services, services...)
}

func NewRoundRobin() *RoundRobin {
	return &RoundRobin{
		services:     nil,
		currentIndex: 0,
	}
}

func (rr *RoundRobin) NextAlive() *Service {
	numServices := len(rr.services)

	for i := 0; i < numServices; i++ {
		currentIndex := rr.currentIndex
		rr.currentIndex = (rr.currentIndex + 1) % numServices

		service := rr.services[currentIndex]
		if service.Alive() {
			return service
		}
	}

	return nil
}

func (rr *RoundRobin) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	rr.NextAlive().ServeHTTP(w, r)
}
