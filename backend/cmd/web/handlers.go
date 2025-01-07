package main

import(
	"net/http"
	"html/template"
	"log"
)

func home(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "text/html")

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

// func CreateUser(w http.ResponseWriter, r *http.Request) {
	
// }