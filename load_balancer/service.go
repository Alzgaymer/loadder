package lb

import (
	"net/http"
	"net/http/httputil"
	"sync"
	"time"
)

type HealthStat struct {
	UnhealthyThreshold int
	TimeoutThreshold   int
	interval           *time.Ticker
	timeout            *time.Ticker
	intervalDuration   time.Duration
	timeoutDuration    time.Duration

	currentTimeout   int
	currentUnhealthy int
}

type Service struct {
	api            *httputil.ReverseProxy
	healthEndpoint string
	alive          bool
	mux            sync.RWMutex
	healthStat     *HealthStat
}

func (s *Service) Healthy() {
	s.healthStat.currentUnhealthy = 0
	s.healthStat.currentTimeout = 0
}

func (s *Service) IncreaseTimeoutThreshold() {
	s.healthStat.currentTimeout++
}

func (s *Service) IncreaseUnhealthyThreshold() {
	s.healthStat.currentTimeout++
}

func (s *Service) TimeoutThreshold() int {
	return s.healthStat.currentTimeout
}
func (s *Service) UnhealthyThreshold() int {
	return s.healthStat.currentUnhealthy
}

func (s *Service) MaxTimeoutThreshold() int {
	return s.healthStat.TimeoutThreshold
}

func (s *Service) MaxUnhealthyThreshold() int {
	return s.healthStat.UnhealthyThreshold
}

func (s *Service) Timeout() (timeoutCh <-chan time.Time, timeoutDuration time.Duration) {
	return s.healthStat.timeout.C, s.healthStat.timeoutDuration
}

func (s *Service) Interval() (intervalCh <-chan time.Time, intervalDuration time.Duration) {
	return s.healthStat.interval.C, s.healthStat.intervalDuration
}

func (s *Service) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.api.ServeHTTP(w, r)
}

func (s *Service) HealthEndpoint() string {
	return s.healthEndpoint
}

func (s *Service) Alive() bool {
	s.mux.Lock()
	defer s.mux.Unlock()

	return s.alive
}

func (s *Service) SetAlive() {
	s.mux.RLock()
	defer s.mux.RUnlock()

	s.alive = true
}
func (s *Service) SetNotAlive() {
	s.mux.RLock()
	defer s.mux.RUnlock()

	s.alive = false
}
