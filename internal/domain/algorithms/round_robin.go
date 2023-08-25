package algorithms

import (
	"container/ring"
	"net/http"
	"net/http/httputil"
)

type RoundRobin struct {
	client *httputil.ReverseProxy
	hosts  *ring.Ring
}

func NewRoundRobin(client *http.Client, address string, ports ...string) *RoundRobin {

}

func (rr *RoundRobin) Handler() func(http.ResponseWriter, *http.Request) {

}

func (rr *RoundRobin) iterate() string {
	host := rr.hosts.Value
	rr.hosts = rr.hosts.Next()

	return host.(string)
}
