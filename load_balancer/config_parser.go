package lb

import (
	"fmt"
	"loadder/config"
	"net/http/httputil"
	"net/url"
	"time"
)

func Parse(c *config.Config) ([]*Service, error) {
	if c == nil {
		return nil, fmt.Errorf("config is nil")
	}

	services := make([]*Service, 0, len(c.Services))

	for _, service := range c.Services {
		u, err := url.Parse(service.Address)
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
			&HealthStat{
				UnhealthyThreshold: service.Healthcheck.UnhealthyThreshold,
				TimeoutThreshold:   service.Healthcheck.TimeoutThreshold,
				interval:           time.NewTicker(service.Healthcheck.IntervalDuration()),
				timeout:            time.NewTicker(service.Healthcheck.TimoutDuration()),
				intervalDuration:   service.Healthcheck.IntervalDuration(),
				timeoutDuration:    service.Healthcheck.TimoutDuration(),
				healthEndpoint:     fmt.Sprintf("%s%s", service.Address, service.Healthcheck.Path),
			},
			&AlgorithmStat{
				weight: wieght,
			}))
	}

	return services, nil
}
