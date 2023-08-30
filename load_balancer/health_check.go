package lb

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"sync"
)

func HealthCheck(ctx context.Context, services []*Service) <-chan error {

	err := make(chan error, len(services))
	wg := &sync.WaitGroup{}
	wg.Add(len(services))

	for _, service := range services {
		go func(ctx context.Context, service *Service, errs chan<- error) {
			defer wg.Done()
			_, timeoutDuration := service.Timeout()
			intervalCh, _ := service.Interval()
			client := &http.Client{}
			// infinite health check only way to escape is

			for service.UnhealthyThreshold() <= service.MaxUnhealthyThreshold() || service.TimeoutThreshold() <= service.MaxTimeoutThreshold() {
				select {
				case <-ctx.Done():
					errs <- ctx.Err()
					return
				case <-intervalCh:
					ctx, _ = context.WithTimeout(ctx, timeoutDuration)

					request, err := http.NewRequestWithContext(ctx, http.MethodGet, service.HealthEndpoint(), http.NoBody)
					if err != nil {
						return
					}

					resp, err := client.Do(request)
					if errors.Is(err, context.DeadlineExceeded) {
						service.IncreaseTimeoutThreshold()
						continue
					} else if err != nil || resp.StatusCode != http.StatusOK {
						service.IncreaseUnhealthyThreshold()
						continue
					}

					_ = resp.Body.Close()

					service.Healthy()
				}
			}

			errs <- fmt.Errorf("service: %s unavailable(unhealthy threshold:%d, timeout threshold:%d)",
				service.HealthEndpoint(), service.UnhealthyThreshold(), service.TimeoutThreshold())

		}(ctx, service, err)
	}

	go func() {
		wg.Wait()
		close(err)
	}()

	return err
}
