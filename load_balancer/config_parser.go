package lb

import (
	"fmt"
	"loadder/config"
	"net/http/httputil"
	"net/url"
	"time"
)

func Parse(c *config.Config) ([]*Service, error) {
	services := make([]*Service, 0, len(c.Services))

	for _, service := range c.Services {
		address := fmt.Sprintf("http://%s", service.Address)

		u, err := url.Parse(address)
		if err != nil {
			return nil, err
		}

		var wieght float64
		if c.Algorithm != WightedRoundRobinAlgorithm {
			wieght = NoWeight
		} else {
			wieght = service.Weight
		}

		services = append(services, NewService(
			httputil.NewSingleHostReverseProxy(u),
			fmt.Sprintf("%s%s", address, service.Healthcheck.Path),
			&HealthStat{
				UnhealthyThreshold: service.Healthcheck.UnhealthyThreshold,
				TimeoutThreshold:   service.Healthcheck.TimeoutThreshold,
				interval:           time.NewTicker(service.Healthcheck.IntervalDuration()),
				timeout:            time.NewTicker(service.Healthcheck.TimoutDuration()),
				intervalDuration:   service.Healthcheck.IntervalDuration(),
				timeoutDuration:    service.Healthcheck.TimoutDuration(),
			},
			&AlgorithmStat{
				weight: wieght,
			}))
	}

	return services, nil
}
