package config

type Config struct {
	Services map[string]*Service `yaml:"services"`
}
