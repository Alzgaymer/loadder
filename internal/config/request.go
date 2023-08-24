package config

type Requests struct {
	HTTP HTTPRequest `yaml:"http"`
}

type HTTPRequest struct {
	GET     map[string]string `yaml:"GET"`
	HEAD    map[string]string `yaml:"HEAD"`
	POST    map[string]string `yaml:"POST"`
	PUT     map[string]string `yaml:"PUT"`
	PATCH   map[string]string `yaml:"PATCH"`
	DELETE  map[string]string `yaml:"DELETE"`
	CONNECT map[string]string `yaml:"CONNECT"`
	OPTIONS map[string]string `yaml:"OPTIONS"`
	TRACE   map[string]string `yaml:"TRACE"`
}
