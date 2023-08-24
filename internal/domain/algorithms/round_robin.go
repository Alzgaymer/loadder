package algorithms

import (
	"container/ring"
	"net/http"
)

type RoundRobin struct {
	client *http.Client
	hosts  *ring.Ring
}

func NewRoundRobin(client *http.Client, address string, ports ...string) *RoundRobin {
	r := ring.New(len(ports))

	for i := 0; i < len(ports); i++ {
		r.Value = string(append([]byte(address), ports[i]...))
		r = r.Next()
	}

	return &RoundRobin{client: client, hosts: r}
}

func (rr *RoundRobin) Handler() func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		host := rr.hosts.Value
		rr.hosts = rr.hosts.Next()

		r.URL.Host = host.(string)

		resp, _ := rr.client.Get(r.URL.String())
		defer resp.Body.Close()

		w.WriteHeader(resp.StatusCode)
		resp.Write(w)
	}
}
