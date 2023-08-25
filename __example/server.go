package main

import (
	"flag"
	"github.com/go-chi/chi/v5"
	"log"
	"net/http"
)

var (
	port = flag.String("port", ":9000", "Specifies server port")
)

func main() {

	flag.Parse()

	router := chi.NewRouter()

	router.Get("/hello-world", func(writer http.ResponseWriter, request *http.Request) {
		log.Printf("Hello World!")
	})

	server := &http.Server{Addr: *port, Handler: router}

	log.Fatalln(server.ListenAndServe())
}
