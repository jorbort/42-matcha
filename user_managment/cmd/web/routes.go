package main

import (
	"net/http"

	"github.com/justinas/alice"
)

func (app *aplication) routes() http.Handler {
	serv := http.NewServeMux()

	dynamicMiddleware := alice.New(authHandler)

	fileServer := http.FileServer(http.Dir("./ui/static/"))

	serv.Handle("GET /ui/static/", http.StripPrefix("/ui/static", fileServer))
	serv.Handle("GET /static/", http.StripPrefix("/static", fileServer))
	//frontend routes that serve html
	serv.HandleFunc("GET /login", app.LoginPage)
	serv.HandleFunc("GET /complete-profile", app.completeProfile)
	serv.HandleFunc("GET /{$}", app.home)
	serv.Handle("GET /profile", dynamicMiddleware.ThenFunc(app.profile))
	serv.HandleFunc("GET /imageUpload", app.imageUploader)
	serv.HandleFunc("GET /forgotPassword", app.forgotPassword)
	serv.HandleFunc("GET /resetPassword", app.newPasswordView)
	serv.HandleFunc("GET /validated", app.validated)
	serv.Handle("GET /settings", dynamicMiddleware.ThenFunc(app.settingsPage))

	// api routes
	serv.Handle("POST /uploadImg", dynamicMiddleware.ThenFunc(app.ImageEndpoint))
	serv.HandleFunc("GET /validate", app.ValidateUser)
	serv.HandleFunc("POST /login", app.UserLogin)
	serv.HandleFunc("POST /create_user", app.CreateUser)
	serv.Handle("POST /complete_profile", dynamicMiddleware.ThenFunc(app.completeUserProfile))
	serv.HandleFunc("POST /SendResetPassword", app.ResetPassword)
	serv.HandleFunc("POST /updatePassword", app.updatePassword)
	standardMiddleware := alice.New(logRequest, commonHeaders)

	return standardMiddleware.Then(serv)
}
