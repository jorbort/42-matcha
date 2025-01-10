package main

import (
	"html/template"
	"log"
	"net/http"
	"regexp"
	"io"
	"errors"

	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jorbort/42-matcha/backend/internals/models"
	"github.com/grahms/godantic"
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

func (app *aplication) CreateUser(w http.ResponseWriter, r *http.Request) {
	var user models.User
	
	r.Body = http.MaxBytesReader(w, r.Body, 1048576)

	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	validator := godantic.Validate{}
	err = validator.BindJSON(body, &user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	passRegexp := regexp.MustCompile(`^[a-zA-Z0-9_#$]{4,25}$`)
	if !passRegexp.MatchString(user.Password){
		http.Error(w, "password must be 4-25 characters long, and alphanumeric values", http.StatusBadRequest)
		return
	}
	if len(user.Username) < 3 || len(user.Username) > 20 {
		http.Error(w, "username must be 3-20 characters long", http.StatusBadRequest)
		return
	} 
	user.Validated = false
	user.Completed = false
	
	var sender EmailSender
	sender.destiantion = user.Email
	validationStr := sender.generateValidationURI()
	user.ValidationCode = validationStr
	err = app.models.InsertUser(r.Context(), &user)
	if err != nil {
		var pgErr *pgconn.PgError
    	if errors.As(err, &pgErr) {
        	switch pgErr.SQLState() {
       		case "23505": 
            	http.Error(w, "Username or email already exists", http.StatusBadRequest)
            	return
        	default:
            	http.Error(w, err.Error(), http.StatusInternalServerError)
            	return
        }
    }
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	err = sender.sendValidationEmail()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	
	w.WriteHeader(http.StatusCreated)
}

func (app *aplication) ValidateUser(w http.ResponseWriter, r *http.Request) {
	code := r.URL.Query.Get("code")
	userInfo , err := app.models.userValidation(r.Context(), code)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	for k , v := range userInfo {
		profileURI := fmt.Sprintf("/complete-profile$id=%d", k)
		if v == false {
			http.Redirect(w, r, profileURI, http.StatusSeeOther)
		}
		else {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
		}
		break;
	}
}

func (app *aplication) Login(w http.ResponseWriter, r *http.Request){
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
	userID := r.URL.Query.Get("id")
	ts , err := template.ParseFiles("ui/html/complete_profile.html")
	if err != nil {
		log.Println(err.Error())
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	err = ts.Execute(w, userID)
	if err != nil {
		log.Println(err.Error())
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
}
