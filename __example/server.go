package main

import (
	"flag"
	"log"
	"net/http"
)

var (
	port = flag.String("port", ":9000", "Specifies server port")
)

func main() {

	flag.Parse()

	mux := http.NewServeMux()

	mux.HandleFunc("/hello-world", func(w http.ResponseWriter, r *http.Request) {
		log.Printf("Hello world!\n")
	})
	mux.HandleFunc("/healthcheck", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	server := &http.Server{Addr: *port, Handler: mux, ReadHeaderTimeout: 0}

	log.Fatalln(server.ListenAndServe())
}
