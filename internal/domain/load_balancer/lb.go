package lb

import (
	"context"
	"errors"
	"net/http"
)

type LoadBalancer struct {
	server   *http.Server
	backends []*Backend
}

func (l *LoadBalancer) Add(v ...*Backend) {
	l.backends = append(l.backends, v...)
}

func NewLoadBalancer(server *http.Server) *LoadBalancer {
	return &LoadBalancer{server: server}
}

func (l *LoadBalancer) Stop(ctx context.Context) error {
	return l.server.Shutdown(ctx)
}

func (l *LoadBalancer) Run(ctx context.Context) error {
	errs := make(chan error)

	go func() {
		errs <- l.server.ListenAndServe()
		close(errs)
	}()

	for {
		select {
		case <-ctx.Done():
			err := l.server.Shutdown(ctx)
			if err != nil {
				return errors.Join(err, ctx.Err())
			}

			return ctx.Err()
		case err := <-errs:
			if err != nil && !errors.Is(err, http.ErrServerClosed) {
				return err
			}

			return nil
		}
	}
}
