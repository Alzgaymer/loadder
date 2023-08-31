package config

import (
	"gopkg.in/yaml.v3"
	"io"
	"log"
	"time"
)

type Config struct {
	Services            map[string]*Service `yaml:"services"`
	Algorithm           string              `yaml:"algorithm"`
	LoadBalancerAddress string              `yaml:"load-balancer-address"`
}

type Healthcheck struct {
	Path               string `yaml:"path"`
	Interval           string `yaml:"interval"`
	Timeout            string `yaml:"timeout"`
	UnhealthyThreshold int    `yaml:"unhealthy-threshold"`
	TimeoutThreshold   int    `yaml:"timeout-threshold"`
}

func (h *Healthcheck) IntervalDuration() time.Duration {
	duration, err := time.ParseDuration(h.Interval)
	if err != nil {
		log.Fatalf("invalid interval duration: %s", err)
	}

	return duration
}

func (h *Healthcheck) TimoutDuration() time.Duration {
	duration, err := time.ParseDuration(h.Timeout)
	if err != nil {
		log.Fatalf("invalid interval duration: %s", err)
	}

	return duration
}

type Service struct {
	Name        string      `yaml:"name"`
	Address     string      `yaml:"address"`
	Weight      float64     `yaml:"weight"`
	Healthcheck Healthcheck `yaml:"healthcheck"`
}

func Parse(file io.Reader) (*Config, error) {
	cfg := &Config{}

	if err := yaml.NewDecoder(file).Decode(cfg); err != nil {
		return nil, err
	}

	return cfg, nil
}
