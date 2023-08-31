package lb

import "net/http"

type WeightedRoundRobin struct {
	services []*Service
	current  int
}

func NewWeightedRoundRobin() *WeightedRoundRobin {
	return &WeightedRoundRobin{
		services: nil,
		current:  0,
	}
}

func (wrr *WeightedRoundRobin) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	wrr.NextAlive().ServeHTTP(w, r)
}

func (wrr *WeightedRoundRobin) Add(services ...*Service) {
	wrr.services = append(wrr.services, services...)
}

func (wrr *WeightedRoundRobin) NextAlive() *Service {
	var (
		totalWeight   = 0.0
		aliveServices int
	)

	for _, service := range wrr.services {
		if service.Alive() {
			totalWeight += service.Weight()
			aliveServices++
		}
	}

	for i := 0; i < aliveServices; i++ {
		wrr.current = (wrr.current + 1) % len(wrr.services)
		if wrr.services[wrr.current].Alive() {
			return wrr.services[wrr.current]
		}
	}

	return nil
}
