package main

import (
	"github.com/go-chi/chi/v5"
	"log"
	"net/http"
)

func main() {
	router := chi.NewRouter()

	router.Get("/hello-world", func(writer http.ResponseWriter, request *http.Request) {
		log.Printf("Hello World!")
	})

	server := &http.Server{Addr: ":9001", Handler: router}

	log.Fatalln(server.ListenAndServe())
}
