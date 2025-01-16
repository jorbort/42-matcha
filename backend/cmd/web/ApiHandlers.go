package main

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/grahms/godantic"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jorbort/42-matcha/backend/internals/models"
)

type loginData struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
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
	if !passRegexp.MatchString(string(user.Password)) {
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
	code := r.URL.Query().Get("code")
	userInfo, err := app.models.UserValidation(r.Context(), code)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	for k, v := range userInfo {
		profileURI := fmt.Sprintf("/complete-profile?id=%d", k)
		if v == false {
			http.Redirect(w, r, profileURI, http.StatusSeeOther)
		} else {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
		}
		break
	}
}

func (app *aplication) UserLogin(w http.ResponseWriter, r *http.Request) {
	var loginData loginData
	var user *models.User

	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	validator := godantic.Validate{}
	err = validator.BindJSON(body, &loginData)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	user, err = app.models.GetUserByUsername(r.Context(), loginData.Username)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if !user.Validated {
		http.Error(w, "user not validated", http.StatusBadRequest)
		return
	}
	if !app.models.VerifyPassword([]byte(loginData.Password), []byte(user.Password)) {
		http.Error(w, "invalid password", http.StatusBadRequest)
		return
	}
	tokenstring, err := app.generateJWT(user.Username, time.Now().Add(time.Hour*24))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	refreshToken, err := app.generateJWT(user.Username, time.Now().Add(time.Hour*24*7))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	http.SetCookie(w, &http.Cookie{
		Name:  "access-token",
		Value: tokenstring,
	})
	http.SetCookie(w, &http.Cookie{
		Name:  "refresh-token",
		Value: refreshToken,
	})
	w.WriteHeader(http.StatusOK)
}

func (app *aplication) generateJWT(username string, exp time.Time) (string, error) {
	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"username": username,
		"exp":      exp.Unix(),
	})
	tokenstring, err := accessToken.SignedString([]byte(os.Getenv("SECRET_KEY")))
	if err != nil {
		return "", err
	}
	return tokenstring, nil
}

func (app *aplication) completeUserProfile(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	defer r.Body.Close()
	var profile models.ProfileInfo

	validator := godantic.Validate{}
	err = validator.BindJSON(body, &profile)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	err = app.models.InsertProfileInfo(r.Context(), &profile)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	err = app.models.UpdateUserCompleted(r.Context(), profile.ID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
}

func (app *aplication) ImageEndpoint(w http.ResponseWriter, r *http.Request) {

	r.Body = http.MaxBytesReader(w, r.Body, 1048576)
	err := r.ParseMultipartForm(1048576)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	userIDStr := r.FormValue("user_id")
	pictureNumberStr := r.FormValue("picture_number")
	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}
	pictureNumber, err := strconv.Atoi(pictureNumberStr)
	if err != nil || pictureNumber < 1 || pictureNumber > 5 {
		http.Error(w, "Invalid picture number", http.StatusBadRequest)
		return
	}
	file, header, err := r.FormFile("image")
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	defer file.Close()

	valid, err := app.models.ValidateImage(file)
	if !valid || err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	extension := filepath.Ext(header.Filename)
	fileName, err := app.models.GenerateFileName(extension)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	fileURI, err := app.models.SaveFile(file, fileName)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = app.models.InsertImage(r.Context(), userID, pictureNumber, fileURI)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
}
