package lb

import (
	"context"
	"fmt"
	"github.com/go-chi/chi"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
	"time"
)

func TestHealthCheck_Success(t *testing.T) {
	ctx := context.Background()

	const serverAddress = "http://localhost:9000"

	services := []*Service{
		{
			healthStat: &HealthStat{
				healthEndpoint:     serverAddress + "/healthcheck",
				UnhealthyThreshold: 2,
				TimeoutThreshold:   2,
				interval:           time.NewTicker(1 * time.Second),
				timeout:            time.NewTicker(1 * time.Second),
				timeoutDuration:    1 * time.Second,
				intervalDuration:   1 * time.Second,
			},
		},
	}

	router := chi.NewRouter()
	router.Get("/healthcheck", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	server := httptest.NewServer(router)

	server.URL = serverAddress

	check := HealthCheck(ctx, services)

	for err := range check {
		if err == nil {
			t.Errorf("HealthCheck() = %v, want %v", err, nil)
		}
	}

	server.Close()
}

func TestHealthCheck_ClosedServer(t *testing.T) {
	ctx := context.Background()

	const serverAddress = "http://localhost:9000"

	services := []*Service{
		{
			healthStat: &HealthStat{
				healthEndpoint:     serverAddress + "/healthcheck",
				UnhealthyThreshold: 2,
				TimeoutThreshold:   2,
				interval:           time.NewTicker(1 * time.Second),
				timeout:            time.NewTicker(1 * time.Second),
				timeoutDuration:    1 * time.Second,
				intervalDuration:   1 * time.Second,
			},
		},
	}

	router := chi.NewRouter()
	router.Get("/healthcheck", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	server := httptest.NewServer(router)

	server.URL = serverAddress
	server.Close()

	check := HealthCheck(ctx, services)

	for err := range check {
		err2 := fmt.Errorf("service: http://localhost:9000/healthcheck unavailable(unhealthy threshold:2, timeout threshold:0)")
		if !reflect.DeepEqual(err, err2) {
			t.Errorf("HealthCheck() = %v, want %v", err, err2)
		}
	}
}

func TestHealthCheck_ContextExceeded(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 0)
	cancel()

	const serverAddress = "http://localhost:9000"

	services := []*Service{
		{
			healthStat: &HealthStat{
				healthEndpoint:     serverAddress + "/healthcheck",
				UnhealthyThreshold: 2,
				TimeoutThreshold:   2,
				interval:           time.NewTicker(1 * time.Second),
				timeout:            time.NewTicker(1 * time.Second),
				timeoutDuration:    1 * time.Second,
				intervalDuration:   1 * time.Second,
			},
		},
	}

	router := chi.NewRouter()
	router.Get("/healthcheck", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	server := httptest.NewServer(router)
	defer server.Close()
	server.URL = serverAddress

	check := HealthCheck(ctx, services)

	for err := range check {
		if !reflect.DeepEqual(err, context.DeadlineExceeded) {
			t.Errorf("HealthCheck() = %v, want %v", err, context.DeadlineExceeded)
		}
	}
}

func TestHealthCheck_MultipleServicesUnavailable(t *testing.T) {
	ctx := context.Background()

	const serverAddress1 = "http://localhost:9000"
	const serverAddress2 = "http://localhost:9001"

	services := []*Service{
		{
			healthStat: &HealthStat{
				healthEndpoint:     serverAddress1 + "/healthcheck",
				UnhealthyThreshold: 2,
				TimeoutThreshold:   2,
				interval:           time.NewTicker(1 * time.Second),
				timeout:            time.NewTicker(1 * time.Second),
				timeoutDuration:    1 * time.Second,
				intervalDuration:   1 * time.Second,
			},
		},
		{
			healthStat: &HealthStat{
				healthEndpoint:     serverAddress2 + "/healthcheck",
				UnhealthyThreshold: 2,
				TimeoutThreshold:   2,
				interval:           time.NewTicker(1 * time.Second),
				timeout:            time.NewTicker(1 * time.Second),
				timeoutDuration:    1 * time.Second,
				intervalDuration:   1 * time.Second,
			},
		},
	}

	router := chi.NewRouter()
	router.Get("/healthcheck", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	server1, server2 := httptest.NewServer(router), httptest.NewServer(router)
	server1.URL, server2.URL = serverAddress1, serverAddress2
	server1.Close()
	server2.Close()

	check := HealthCheck(ctx, services)
	var count int
	for range check {
		count++
	}

	if count != 2 {
		t.Errorf("HealthCheck() = %v, want %v", count, 2)
	}
}
