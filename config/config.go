package config

import (
	"gopkg.in/yaml.v3"
	"io"
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

type Service struct {
	Name        string      `yaml:"name"`
	Address     string      `yaml:"address"`
	Healthcheck Healthcheck `yaml:"healthcheck"`
}

func Parse(file io.Reader) (*Config, error) {
	cfg := &Config{}

	if err := yaml.NewDecoder(file).Decode(cfg); err != nil {
		return nil, err
	}

	return cfg, nil
}
