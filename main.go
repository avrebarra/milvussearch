package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/avrebarra/milvus-dating/server"
)

func main() {
	s, err := server.New(server.Config{})
	catcherr(err)

	log.Println("starting server in http://localhost:5678/")
	err = http.ListenAndServe(":5678", s.GetHandler())
	catcherr(err)
}

func catcherr(err error) {
	if err != nil {
		err = fmt.Errorf("unexpected error: %w", err)
		log.Fatal(err)
	}
}
