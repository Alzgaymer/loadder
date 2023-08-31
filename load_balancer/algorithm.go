package lb

import (
	"errors"
	"net/http"
)

type Algorithm interface {
	http.Handler
	NextAlive() *Service
	Add(services ...*Service)
}

const (
	// RoundRobinAlgorithm represents round-robin algorithm
	RoundRobinAlgorithm = "round-robin"

	// WightedRoundRobinAlgorithm represents wighted round-robin algorithm
	WightedRoundRobinAlgorithm = "weighted-round-robin"
)

var ErrAlgorithmNotSupported = errors.New("algorithm not supported")

func DefineAlgorithm(algorithm string) (Algorithm, error) {
	switch algorithm {
	case RoundRobinAlgorithm:
		return NewRoundRobin(), nil
	case WightedRoundRobinAlgorithm:
		return NewWeightedRoundRobin(), nil
	default:
		return nil, ErrAlgorithmNotSupported
	}
}
