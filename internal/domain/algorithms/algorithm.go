package algorithms

import (
	"errors"
	"net/http"
)

var (
	ErrInvalidAlgorithm = errors.New("invalid algorithm")
)

const (

	// RR defines Round-robin algorithm.
	RR = "RR"

	// WRR defines Weighted Round-robin algorithm.
	WRR = "WRR"

	// LC defines Least Connections algorithm.
	LC = "LC"

	// LRT defines Least Response Time algorithm.
	LRT = "LRT"

	// LB defines Least Bandwidth algorithm.
	LB = "LB"

	// H defines Hashing algorithm.
	H = "H"
)

type Algorithm interface {
	Handler() func(w http.ResponseWriter, r *http.Request)
}

type Builder struct {
	client  *http.Client
	address string
	ports   []string
}

func NewBuilder(client *http.Client, address string, ports []string) *Builder {
	return &Builder{client: client, address: address, ports: ports}
}

func (b *Builder) DefineAlgorithm(algo string) (Algorithm, error) {
	switch algo {
	case RR:
		return NewRoundRobin(b.client, b.address, b.ports...), nil
	case WRR:
	case LC:
	case LRT:
	case LB:
	case H:
	default:
		return nil, ErrInvalidAlgorithm
	}

	return nil, nil
}
