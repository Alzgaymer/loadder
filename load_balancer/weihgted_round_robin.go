package lb

import "net/http"

type WeightedRoundRobin struct {
	services     []*Service
	currentIndex int
}

func NewWeightedRoundRobin() *WeightedRoundRobin {
	return &WeightedRoundRobin{
		services:     nil,
		currentIndex: 0,
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
		totalWeight = 0.0
		numServices = len(wrr.services)
	)

	for _, service := range wrr.services {
		if service.Alive() {
			totalWeight += service.Weight()
		}
	}

	for i := 0; i < numServices; i++ {
		currentIndex := wrr.currentIndex
		wrr.currentIndex = (wrr.currentIndex + 1) % numServices

		service := wrr.services[currentIndex]
		if service.Alive() {
			totalWeight -= service.Weight()
			if totalWeight <= 0 {
				return service
			}
		}
	}

	return nil
}
