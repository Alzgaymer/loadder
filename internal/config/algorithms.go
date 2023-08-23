package config

import "errors"

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
