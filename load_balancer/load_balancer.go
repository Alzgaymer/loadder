package lb

import (
	"context"
	"go.uber.org/zap"
	"net/http"
)

type LoadBalancer struct {
	server   *http.Server
	services []*Service
	logger   *zap.Logger
}

func (lb *LoadBalancer) Run(ctx context.Context) error {
	healthCheck := HealthCheck(ctx, lb.services)

	go func() {
		for err := range healthCheck {
			lb.logger.Error("health check error", zap.Error(err))
		}
	}()

	lb.logger.Info("load balancer started", zap.String("address", lb.server.Addr))

	return lb.server.ListenAndServe()
}

func WithLogger(logger *zap.Logger) Option {
	return func(balancer *LoadBalancer) error {
		balancer.logger = logger
		return nil
	}
}

func WithAlgorithm(algorithm Algorithm) Option {
	return func(balancer *LoadBalancer) error {
		balancer.server.Handler = algorithm
		return nil
	}
}

func WithAddress(address string) Option {
	return func(balancer *LoadBalancer) error {
		balancer.server.Addr = address
		return nil
	}
}

func WithServices(s ...*Service) Option {
	return func(balancer *LoadBalancer) error {
		balancer.services = s
		return nil
	}
}

type Option func(balancer *LoadBalancer) error

func NewBalancer(opts ...Option) (*LoadBalancer, error) {
	balancer := &LoadBalancer{server: &http.Server{}}
	for _, opt := range opts {
		err := opt(balancer)
		if err != nil {
			return nil, err
		}
	}

	return balancer, nil
}
