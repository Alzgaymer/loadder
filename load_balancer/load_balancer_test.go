package lb

import (
	"context"
	"github.com/go-chi/chi"
	"go.uber.org/zap/zaptest"
	"net/http"
	"testing"
	"time"
)

func TestLoadBalancer_RunSuccess(t *testing.T) {
	router := chi.NewRouter()
	router.Get("/healthcheck", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	balancer, _ := NewBalancer(func(balancer *LoadBalancer) error {
		balancer.server.Handler = router
		return nil
	}, WithLogger(zaptest.NewLogger(t)), WithAddress("localhost:9000"))

	go func() {
		_ = balancer.Run(context.Background())
	}()
	time.Sleep(500 * time.Millisecond)
	client := &http.Client{}
	resp, err := client.Get("http://localhost:9000/healthcheck")

	if err != nil && resp.StatusCode != http.StatusOK {
		t.Errorf("LoadBalancer.Run() = %v, want %v", err, nil)
	}

	_ = balancer.Stop(context.Background())
}

func TestLoadBalancer_RunNoAddress(t *testing.T) {
	balancer, _ := NewBalancer(WithLogger(zaptest.NewLogger(t)))

	if err := balancer.Run(context.Background()); err == nil {
		t.Errorf("LoadBalancer.Run() = %v, want %v", err, nil)
	}

}
