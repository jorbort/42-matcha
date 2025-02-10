package main

import (
	"net/http"

	"github.com/justinas/alice"
)

func (app *aplication) routes() http.Handler {
	serv := http.NewServeMux()
	dynamicMiddleware := alice.New(authHandler)

	fileServer := http.FileServer(http.Dir("./ui/static/"))
	serv.Handle("GET /static/", http.StripPrefix("/static", fileServer))
	//routes that serve html

	//api and ws endpoints
	serv.Handle("GET /ws", dynamicMiddleware.ThenFunc(app.handleWebSocket))
	standardMiddleware := alice.New(logRequest, commonHeaders)
	return standardMiddleware.Then(serv)
}
