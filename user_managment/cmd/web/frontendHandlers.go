package main

import (
	"html/template"
	"log"
	"net/http"
)

func (app *aplication) home(w http.ResponseWriter, r *http.Request) {

	ts, err := template.ParseFiles("ui/html/index.html")
	if err != nil {
		log.Println(err.Error())
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	err = ts.Execute(w, nil)
	if err != nil {
		log.Println(err.Error())
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
}

func (app *aplication) profile(w http.ResponseWriter, r *http.Request) {
	err := app.templates.ExecuteTemplate(w, "profile", nil)
	if err != nil {
		log.Println(err.Error())
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
}
func (app *aplication) settingsPage(w http.ResponseWriter, r *http.Request) {
	data := struct {
		CSSLink string
	}{
		CSSLink: "static/css/settings.css",
	}
	err := app.templates.ExecuteTemplate(w, "settings", data)
	if err != nil {
		log.Println(err.Error())
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
}

func (app *aplication) LoginPage(w http.ResponseWriter, r *http.Request) {
	ts, err := template.ParseFiles("ui/html/login.html")
	if err != nil {
		log.Println(err.Error())
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	err = ts.Execute(w, nil)
	if err != nil {
		log.Println(err.Error())
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
}

func (app *aplication) completeProfile(w http.ResponseWriter, r *http.Request) {
	userID := r.URL.Query().Get("id")
	data := struct {
		UserID string
	}{
		UserID: userID,
	}
	ts, err := template.ParseFiles("ui/html/complete_profile.html")
	if err != nil {
		log.Println(err.Error())
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	err = ts.Execute(w, data)
	if err != nil {
		log.Println(err.Error())
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
}

func (app *aplication) imageUploader(w http.ResponseWriter, r *http.Request) {
	ts, err := template.ParseFiles("ui/html/image_uploader.html")
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	err = ts.Execute(w, nil)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

}

func (app *aplication) forgotPassword(w http.ResponseWriter, r *http.Request) {
	ts, err := template.ParseFiles("ui/html/reset_password.html")
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	err = ts.Execute(w, nil)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
}

func (app *aplication) newPasswordView(w http.ResponseWriter, r *http.Request) {
	ts, err := template.ParseFiles("ui/html/newPassword.html")
	if err != nil {
		log.Println(err.Error())
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	err = ts.Execute(w, nil)
	if err != nil {
		log.Println(err.Error())
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
}
func (app *aplication) validated(w http.ResponseWriter, r *http.Request) {
	err := app.templates.ExecuteTemplate(w, "validation", nil)
	if err != nil {
		log.Println(err.Error())
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
}
