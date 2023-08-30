package lb

import (
	"context"
	"errors"
	"go.uber.org/zap"
	"loadder/config"
	"net/http"
	"net/http/httputil"
	"net/url"
	"time"
)

type options func(balancer *LoadBalancer) error

type LoadBalancer struct {
	logger   *zap.Logger
	server   *http.Server
	services []Service
}

func (b *LoadBalancer) Shutdown(ctx context.Context) error {
	return b.server.Shutdown(ctx)
}

func (b *LoadBalancer) Run(ctx context.Context) error {
	b.logger.Info("starting load balancer", zap.String("port", b.server.Addr))

	done := make(chan struct{})
	go func() {
		if err := b.server.ListenAndServe(); err != nil {
			b.logger.Panic("failed to start load balancer", zap.Error(err), zap.String("port", b.server.Addr))
		}
		done <- struct{}{}
	}()

	for _, service := range b.services {
		go func(ctx context.Context, service Service) {
			for {
				select {
				case <-service.Interval():
					if err := service.HealthCheck(); err != nil {
						b.logger.Error("failed to perform health check", zap.Error(err))
					}
				case <-ctx.Done():
					return
				}
			}
		}(ctx, service)
	}

	select {
	case <-ctx.Done():
		b.logger.Info("stopping load balancer: context exceeded", zap.Error(ctx.Err()))
		return errors.Join(ctx.Err(), b.server.Shutdown(ctx))
	case <-done:
		b.logger.Info("stopping load balancer: done")
		return nil
	}
}

func WithPort(port string) options {
	return func(balancer *LoadBalancer) error {
		balancer.server.Addr = port
		return nil
	}
}

func WithAlgorithm(algorithm string) options {
	return func(balancer *LoadBalancer) error {
		definedAlgorithm, err := DefineAlgorithm(algorithm, balancer.services...)
		if err != nil {
			return err
		}

		balancer.server.Handler = definedAlgorithm

		return nil
	}
}

func NewLoadBalancer(server *http.Server, logger *zap.Logger, services []Service, opts ...options) (*LoadBalancer, error) {
	balancer := &LoadBalancer{server: server, logger: logger, services: services}

	for _, opt := range opts {
		err := opt(balancer)
		if err != nil {
			return nil, err
		}
	}

	return balancer, nil
}

type Service interface {
	http.Handler
	Alive() bool
	HealthCheck() error
	Interval() <-chan time.Time
}

type HealthCheck struct {
	Path               string
	Interval           *time.Ticker
	Timeout            *time.Ticker
	intervalDuration   time.Duration
	timeoutDuration    time.Duration
	UnhealthyThreshold int
	TimeoutThreshold   int
	currentTimeout     int
	currentUnhealthy   int
}

func (c *HealthCheck) HealthCheck() error {
	for (c.currentTimeout < c.TimeoutThreshold) && (c.currentUnhealthy < c.UnhealthyThreshold) {
		<-c.Interval.C
		resp, err := http.Get(c.Path)
		if errors.Is(err, context.DeadlineExceeded) {
			c.currentTimeout++
		} else if err != nil || resp.StatusCode != http.StatusOK {
			c.currentUnhealthy++
			return err
		}

		c.currentTimeout = 0
		c.currentUnhealthy = 0
	}

	return nil
}

func Parse(c *config.Config) ([]Service, error) {
	backends := make([]Service, 0, len(c.Services))

	for _, service := range c.Services {
		u, err := url.Parse(service.Address)
		if err != nil {
			continue
		}

		timeoutDuration, err := service.Healthcheck.TimeoutDuration()
		if err != nil {
			continue
		}

		intervalDuration, err := service.Healthcheck.IntervalDuration()
		if err != nil {
			continue
		}

		backends = append(backends, NewBackend(u, httputil.NewSingleHostReverseProxy(u), &HealthCheck{
			Path:               service.Address + service.Healthcheck.Path,
			Interval:           time.NewTicker(intervalDuration),
			Timeout:            time.NewTicker(timeoutDuration),
			intervalDuration:   intervalDuration,
			timeoutDuration:    timeoutDuration,
			UnhealthyThreshold: service.Healthcheck.UnhealthyThreshold,
			TimeoutThreshold:   service.Healthcheck.TimeoutThreshold,
		}))
	}

	return backends, nil
}
