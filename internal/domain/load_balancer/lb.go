package lb

import (
	"context"
	"errors"
	"log"
	"net/http"
	"sync"
)

type LoadBalancer struct {
	server   *http.Server
	services []Service
}

type Service interface {
	Name() string
	Shutdown(ctx context.Context) error
	ListenAndServe() error
}

type Backend interface {
	http.Handler
}

type Backends interface {
	Add(port []string, b ...Backend) error
}

func (l *LoadBalancer) Add(v ...Service) {
	l.services = append(l.services, v...)
}

func NewLoadBalancer(server *http.Server, services ...Service) *LoadBalancer {
	return &LoadBalancer{server: server, services: services}
}

func (l *LoadBalancer) Stop(ctx context.Context) error {
	return l.server.Shutdown(ctx)
}

type response struct {
	ServiceName string
	err         error
}

func (l *LoadBalancer) Run(ctx context.Context) error {
	res := make(chan *response, len(l.services))
	wg := sync.WaitGroup{}
	wg.Add(len(l.services))

	for _, service := range l.services {
		go func(ctx context.Context, service Service, res chan *response) {
			defer wg.Done()

			errs := make(chan error)

			go func(err chan<- error) {
				err <- service.ListenAndServe()
				close(err)
			}(errs)

			for {
				select {
				case <-ctx.Done():
					err := service.Shutdown(ctx)
					if err != nil {
						res <- &response{
							ServiceName: service.Name(),
							err:         errors.Join(err, ctx.Err()),
						}
					}

					res <- &response{
						ServiceName: service.Name(),
						err:         ctx.Err(),
					}
				case err := <-errs:
					if err != nil && !errors.Is(err, http.ErrServerClosed) {
						res <- &response{
							ServiceName: service.Name(),
							err:         err,
						}
					}

					res <- nil
				}
			}
		}(ctx, service, res)
	}

	go func() {
		wg.Wait()
		close(res)
	}()

	for r := range res {
		log.Println(r)
	}

	return nil
}
