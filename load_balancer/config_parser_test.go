package lb

import (
	"loadder/config"
	"net/http/httputil"
	"net/url"
	"testing"
	"time"
)

func TestParse_RightConfigWithoutAlgorithm(t *testing.T) {
	c := &config.Config{
		Services: map[string]*config.Service{
			"first-service": {
				Name:    "hello-world-1",
				Address: "http://localhost:9000",
				Weight:  2.0,
				Healthcheck: &config.Healthcheck{
					Path:               "/healthcheck",
					Interval:           "1s",
					Timeout:            "2s",
					UnhealthyThreshold: 2,
					TimeoutThreshold:   2,
				},
			},
		},
	}

	u, _ := url.Parse("http://localhost:9000")

	service := []*Service{
		{
			api: httputil.NewSingleHostReverseProxy(u),
			healthStat: &HealthStat{
				interval:           time.NewTicker(1 * time.Second),
				timeout:            time.NewTicker(2 * time.Second),
				healthEndpoint:     "http://localhost:9000/healthcheck",
				intervalDuration:   1 * time.Second,
				timeoutDuration:    2 * time.Second,
				TimeoutThreshold:   2,
				UnhealthyThreshold: 2,
				currentUnhealthy:   0,
				currentTimeout:     0,
			},
			algorithmStat: &AlgorithmStat{weight: NoWeight},
		},
	}

	got, _ := Parse(c)

	for i, s := range service {
		if s.MaxTimeoutThreshold() != got[i].MaxTimeoutThreshold() {
			t.Errorf("Parse() = %v, want %v", s.MaxTimeoutThreshold(), got[i].MaxTimeoutThreshold())
		}

		if s.MaxUnhealthyThreshold() != got[i].MaxUnhealthyThreshold() {
			t.Errorf("Parse() = %v, want %v", s.MaxUnhealthyThreshold(), got[i].MaxUnhealthyThreshold())
		}

		if s.Weight() != got[i].Weight() {
			t.Errorf("Parse() = %v, want %v", s.Weight(), got[i].Weight())
		}

		if s.HealthEndpoint() != got[i].HealthEndpoint() {
			t.Errorf("Parse() = %v, want %v", s.HealthEndpoint(), got[i].HealthEndpoint())
		}

		_, intervalDuration := s.Interval()
		_, gotIntervalDuration := got[i].Interval()

		if intervalDuration != gotIntervalDuration {
			t.Errorf("Parse() = %v, want %v", intervalDuration, gotIntervalDuration)
		}
	}
}

func TestParse_RightConfigWithAlgorithm(t *testing.T) {
	c := &config.Config{
		Algorithm: WightedRoundRobinAlgorithm,
		Services: map[string]*config.Service{
			"first-service": {
				Name:    "hello-world-1",
				Address: "http://localhost:9000",
				Weight:  2.0,
				Healthcheck: &config.Healthcheck{
					Path:               "/healthcheck",
					Interval:           "1s",
					Timeout:            "2s",
					UnhealthyThreshold: 2,
					TimeoutThreshold:   2,
				},
			},
		},
	}

	u, _ := url.Parse("http://localhost:9000")

	service := []*Service{
		{
			api: httputil.NewSingleHostReverseProxy(u),
			healthStat: &HealthStat{
				interval:           time.NewTicker(1 * time.Second),
				timeout:            time.NewTicker(2 * time.Second),
				healthEndpoint:     "http://localhost:9000/healthcheck",
				intervalDuration:   1 * time.Second,
				timeoutDuration:    2 * time.Second,
				TimeoutThreshold:   2,
				UnhealthyThreshold: 2,
				currentUnhealthy:   0,
				currentTimeout:     0,
			},
			algorithmStat: &AlgorithmStat{weight: 2.0},
		},
	}

	got, _ := Parse(c)

	for i, s := range service {
		if s.MaxTimeoutThreshold() != got[i].MaxTimeoutThreshold() {
			t.Errorf("Parse() = %v, want %v", s.MaxTimeoutThreshold(), got[i].MaxTimeoutThreshold())
		}

		if s.MaxUnhealthyThreshold() != got[i].MaxUnhealthyThreshold() {
			t.Errorf("Parse() = %v, want %v", s.MaxUnhealthyThreshold(), got[i].MaxUnhealthyThreshold())
		}

		if s.Weight() != got[i].Weight() {
			t.Errorf("Parse() = %v, want %v", s.Weight(), got[i].Weight())
		}

		if s.HealthEndpoint() != got[i].HealthEndpoint() {
			t.Errorf("Parse() = %v, want %v", s.HealthEndpoint(), got[i].HealthEndpoint())
		}

		_, intervalDuration := s.Interval()
		_, gotIntervalDuration := got[i].Interval()

		if intervalDuration != gotIntervalDuration {
			t.Errorf("Parse() = %v, want %v", intervalDuration, gotIntervalDuration)
		}
	}
}

func TestParse_RightConfigWithNonWeightedRRAlgorithm(t *testing.T) {
	c := &config.Config{
		Algorithm: RoundRobinAlgorithm,
		Services: map[string]*config.Service{
			"first-service": {
				Name:    "hello-world-1",
				Address: "http://localhost:9000",
				Weight:  2.0,
				Healthcheck: &config.Healthcheck{
					Path:               "/healthcheck",
					Interval:           "1s",
					Timeout:            "2s",
					UnhealthyThreshold: 2,
					TimeoutThreshold:   2,
				},
			},
		},
	}

	u, _ := url.Parse("http://localhost:9000")

	service := []*Service{
		{
			api: httputil.NewSingleHostReverseProxy(u),
			healthStat: &HealthStat{
				interval:           time.NewTicker(1 * time.Second),
				timeout:            time.NewTicker(2 * time.Second),
				healthEndpoint:     "http://localhost:9000/healthcheck",
				intervalDuration:   1 * time.Second,
				timeoutDuration:    2 * time.Second,
				TimeoutThreshold:   2,
				UnhealthyThreshold: 2,
				currentUnhealthy:   0,
				currentTimeout:     0,
			},
			algorithmStat: &AlgorithmStat{weight: NoWeight},
		},
	}

	got, _ := Parse(c)

	for i, s := range service {
		if s.MaxTimeoutThreshold() != got[i].MaxTimeoutThreshold() {
			t.Errorf("Parse() = %v, want %v", s.MaxTimeoutThreshold(), got[i].MaxTimeoutThreshold())
		}

		if s.MaxUnhealthyThreshold() != got[i].MaxUnhealthyThreshold() {
			t.Errorf("Parse() = %v, want %v", s.MaxUnhealthyThreshold(), got[i].MaxUnhealthyThreshold())
		}

		if s.Weight() != got[i].Weight() {
			t.Errorf("Parse() = %v, want %v", s.Weight(), got[i].Weight())
		}

		if s.HealthEndpoint() != got[i].HealthEndpoint() {
			t.Errorf("Parse() = %v, want %v", s.HealthEndpoint(), got[i].HealthEndpoint())
		}

		_, intervalDuration := s.Interval()
		_, gotIntervalDuration := got[i].Interval()

		if intervalDuration != gotIntervalDuration {
			t.Errorf("Parse() = %v, want %v", intervalDuration, gotIntervalDuration)
		}
	}
}

func TestParse_NilConfig(t *testing.T) {
	var c *config.Config

	var service []*Service

	got, _ := Parse(c)

	if len(got) != len(service) {
		t.Errorf("Parse() = %v, want %v", got, service)
	}
}

func TestParse_EmptyConfig(t *testing.T) {
	c := &config.Config{}

	var services []*Service

	got, _ := Parse(c)

	if len(got) != len(services) {
		t.Errorf("Parse() = %v, want %v", got, services)
	}
}
