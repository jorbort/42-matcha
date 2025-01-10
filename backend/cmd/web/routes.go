package main

import (
	"net/http"
	"github.com/justinas/alice"
)

func (app *aplication) routes() http.Handler{
	serv := http.NewServeMux()

	fileServer := http.FileServer(http.Dir("./ui/static/"))

	serv.Handle("GET /static/", http.StripPrefix("/static", fileServer))

	serv.HandleFunc("GET /validate", app.ValidateUser)
	
	serv.HandleFunc("GET /login", app.Login)
	serv.HandleFunc("GET /complete-profile", app.completeProfile)

	serv.HandleFunc("GET /{$}", app.home)

	serv.HandleFunc("POST /create_user", app.CreateUser)

	standardMiddleware := alice.New(logRequest, commonHeaders)

	return standardMiddleware.Then(serv)
}