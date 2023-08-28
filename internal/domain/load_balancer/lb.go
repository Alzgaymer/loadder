package lb

import (
	"context"
	"errors"
	"log"
	"net/http"
	"sync"
)

type LoadBalancer struct {
	services []Service
}

func NewLoadBalancer(services ...Service) *LoadBalancer {
	return &LoadBalancer{services: services}
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

func (l *LoadBalancer) Add(v Service) {
	l.services = append(l.services, v)
}

type response struct {
	ServiceName string
	err         error
}

func (l *LoadBalancer) Run(ctx context.Context) error {
	// res stands for every service proxy server error channel
	res := make(chan *response, len(l.services))
	// wg stands for closing every proxy server and close res channel
	wg := sync.WaitGroup{}

	wg.Add(len(l.services))

	// for every service we start goroutine with ListenAndServer server and waiting until it close or ctx.Done channel
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

	// waits in parallel for every proxy goroutine. Only closing response channel means
	go func() {
		wg.Wait()
		close(res)
	}()

	// listening for responses from proxy servers
	for r := range res {
		log.Println(r)
	}

	return nil
}
