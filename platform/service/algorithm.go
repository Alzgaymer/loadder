package service

import (
	"errors"
	"loadder/platform/backend"
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
	http.Handler
}

func DefineAlgorithm(algo string, b ...*backend.Backend) (Algorithm, error) {
	switch algo {
	case RR:
		return NewRoundRobinAlgorithm(b...), nil
	case WRR:
	case LC:
	case LRT:
	case LB:
	case H:
	}
	return nil, ErrInvalidAlgorithm
}
