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

func testPage(w http.ResponseWriter, r *http.Request){
	ts, err := template.ParseFiles("ui/html/test.html")
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

func (app *aplication) LoginPage(w http.ResponseWriter, r *http.Request){
	ts , err := template.ParseFiles("ui/html/login.html")
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

func (app *aplication) completeProfile(w http.ResponseWriter, r *http.Request){
	userID := r.URL.Query().Get("id")
	data := struct {
		UserID string
	}{
		UserID: userID,
	}
	ts , err := template.ParseFiles("ui/html/complete_profile.html")
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
