package main

import (
	"net/http"
)

func (app *aplication) routes() *http.ServeMux{
	serv := http.NewServeMux()

	fileServer := http.FileServer(http.Dir("./ui/static/"))

	serv.Handle("GET /static/", http.StripPrefix("/static", fileServer))
	
	serv.HandleFunc("GET /{$}", app.home)

	serv.HandleFunc("POST /create_user", app.CreateUser)

	return serv
}