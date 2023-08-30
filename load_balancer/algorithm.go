package lb

import "net/http"

type Algorithm interface {
	http.Handler
	NextAlive() *Service
}
