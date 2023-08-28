package config

type Service struct { //nolint:govet
	Name      string   `yaml:"service-name"`
	Address   string   `yaml:"service-address"`
	Ports     string   `yaml:"service-ports"`
	Exclude   []string `yaml:"exclude"`
	ProxyPort string   `yaml:"proxy-port"`
	Algorithm string   `yaml:"algorithm"`
}
