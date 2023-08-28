package lb

type Service struct { //nolint:govet
	backends *Backends
}

func NewService(backends *Backends) *Service {
	return &Service{backends: backends}
}
