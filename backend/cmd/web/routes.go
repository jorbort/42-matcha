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
	//frontend routes that serve html
	serv.HandleFunc("GET /login", app.LoginPage)
	serv.HandleFunc("GET /complete-profile", app.completeProfile)
	serv.HandleFunc("GET /{$}", app.home)
	serv.Handle("GET /testPage", dynamicMiddleware.ThenFunc(app.home))
	serv.HandleFunc("GET /imageUpload", app.imageUploader)

	// api routes
	serv.Handle("POST /uploadImg", dynamicMiddleware.ThenFunc(app.ImageEndpoint))
	serv.HandleFunc("GET /validate", app.ValidateUser)
	serv.HandleFunc("POST /login", app.UserLogin)
	serv.HandleFunc("POST /create_user", app.CreateUser)
	serv.Handle("POST /complete_profile", dynamicMiddleware.ThenFunc(app.completeUserProfile))
	serv.HandleFunc("GET /SendResetPassword", app.ResetPassword)
	serv.HandleFunc("POST /updatePassword", app.updatePassword)
	standardMiddleware := alice.New(logRequest, commonHeaders)

	return standardMiddleware.Then(serv)
}
