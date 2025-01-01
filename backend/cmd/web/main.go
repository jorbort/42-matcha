package main

import (
	"log"
	"net/http"
)

func home(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello World"))
}

func main() {
	serv := http.NewServeMux()

	serv.HandleFunc("GET /{$}", home)

	err := http.ListenAndServe(":3000", serv)
	log.Fatal(err)
}
