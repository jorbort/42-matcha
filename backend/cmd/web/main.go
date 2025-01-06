package main

import (
	"log"
	"net/http"

	// _"github.com/jackc/pgx/v5"
)

func main() {
	serv := http.NewServeMux()

	serv.HandleFunc("GET /{$}", home)

	log.Println("Starting server on :3000")
	err := http.ListenAndServe(":3000", serv)
	log.Fatal(err)
}
