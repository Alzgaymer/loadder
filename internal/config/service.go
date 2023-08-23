package config

type Service struct { //nolint:govet
	Name     string    `yaml:"service-name"`
	Address  string    `yaml:"address"`
	Ports    []string  `yaml:"ports"`
	Requests *Requests `yaml:"requests"`
}
