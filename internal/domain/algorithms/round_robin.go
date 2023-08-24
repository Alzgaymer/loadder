package algorithms

import (
	"container/ring"
	"fmt"
	"log"
	"net/http"
)

type RoundRobin struct {
	client *http.Client
	hosts  *ring.Ring
}

func NewRoundRobin(client *http.Client, address string, ports ...string) *RoundRobin {
	r := ring.New(len(ports))

	for i := 0; i < len(ports); i++ {
		r.Value = string(append([]byte(address+":"), ports[i]...))
		r = r.Next()
	}

	return &RoundRobin{client: client, hosts: r}
}

func (rr *RoundRobin) Handler() func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		r.Host = rr.iterate()

		req, err := http.NewRequestWithContext(r.Context(), r.Method, fmt.Sprintf("http://%s%s", r.Host, r.URL.Path), http.NoBody)
		if err != nil {
			// rr.logger
			log.Printf("%v\n", err)
			w.WriteHeader(http.StatusBadRequest)

			return
		}

		resp, err := rr.client.Do(req)
		if err != nil {
			// rr.logger
			log.Printf("%v\n", err)
			w.WriteHeader(http.StatusInternalServerError)

			return
		}
		defer resp.Body.Close()

		w.WriteHeader(resp.StatusCode)
		resp.Write(w)
	}
}

func (rr *RoundRobin) iterate() string {
	host := rr.hosts.Value
	rr.hosts = rr.hosts.Next()

	return host.(string)
}
