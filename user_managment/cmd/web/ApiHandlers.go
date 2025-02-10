package main

import (
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/grahms/godantic"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jorbort/42-matcha/user_managment/internals/models"
)

type loginData struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}
type loginResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}
type uploadResponse struct {
	Message string `json:"message"`
	URL     string `json:"url"`
}
type ErrorResponse struct {
	Message string `json:"message"`
	Code    int    `json:"code"`
}

type PasswordReset struct {
	Email string `json:"email" binding:"required" format:"email"`
}
type NewPassword struct {
	Password string `json:"password" binding:"required"`
	Code     string `json:"code" binding:"required"`
}

func (app *aplication) CreateUser(w http.ResponseWriter, r *http.Request) {
	var user models.User
	response := ErrorResponse{}
	r.Body = http.MaxBytesReader(w, r.Body, 1048576)

	body, err := io.ReadAll(r.Body)
	if err != nil {
		writeJsonError(w, http.StatusBadRequest, err.Error())
		return
	}
	defer r.Body.Close()

	validator := godantic.Validate{}
	err = validator.BindJSON(body, &user)
	if err != nil {
		writeJsonError(w, http.StatusBadRequest, err.Error())
		return
	}
	passRegexp := regexp.MustCompile(`^[a-zA-Z0-9_#$]{4,25}$`)
	if !passRegexp.MatchString(string(user.Password)) {
		writeJsonError(w, http.StatusBadRequest, "password must be 4-25 characters long, and alphanumeric values")
		return
	}
	if len(user.Username) < 3 || len(user.Username) > 20 {
		writeJsonError(w, http.StatusBadRequest, "username must be 3-20 characters long")
		return
	}
	user.Validated = false
	user.Completed = false

	var sender EmailSender
	sender.Destiantion = user.Email
	validationStr := sender.generateValidationURI()
	user.ValidationCode = validationStr
	err = app.models.InsertUser(r.Context(), &user)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			switch pgErr.SQLState() {
			case "23505":
				writeJsonError(w, http.StatusBadRequest, "Username or email already exists")
				return
			default:
				writeJsonError(w, http.StatusInternalServerError, "Internal server error")
				return
			}
		}
		writeJsonError(w, http.StatusInternalServerError, "Internal server error")
		return
	}
	err = sender.sendValidationEmail("Validate your account", "validate", "Click the link below to validate your account!!")
	if err != nil {
		writeJsonError(w, http.StatusInternalServerError, err.Error())
		return
	}
	w.WriteHeader(http.StatusCreated)
	response.Message = "User created successfully"
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func (app *aplication) ValidateUser(w http.ResponseWriter, r *http.Request) {
	code := r.URL.Query().Get("code")
	err := app.models.UserValidation(r.Context(), code)
	if err != nil {
		writeJsonError(w, http.StatusBadRequest, err.Error())
		return
	}
	http.Redirect(w, r, "http://localhost:3000/validated", http.StatusSeeOther)
}

func (app *aplication) UserLogin(w http.ResponseWriter, r *http.Request) {
	var loginData loginData
	var user *models.User

	body, err := io.ReadAll(r.Body)
	if err != nil {
		log.Println(err.Error())
		writeJsonError(w, http.StatusBadRequest, err.Error())
		return
	}
	defer r.Body.Close()

	validator := godantic.Validate{}
	err = validator.BindJSON(body, &loginData)
	if err != nil {
		log.Println(err.Error())
		writeJsonError(w, http.StatusBadRequest, err.Error())
		return
	}
	user, err = app.models.GetUserByUsername(r.Context(), loginData.Username)
	if err != nil {
		log.Println(err.Error())
		writeJsonError(w, http.StatusNotFound, err.Error())
		return
	}
	if !user.Validated {
		log.Println("user not validated")
		writeJsonError(w, http.StatusBadRequest, "user not validated")
		return
	}
	if !app.models.VerifyPassword([]byte(loginData.Password), []byte(user.Password)) {
		writeJsonError(w, http.StatusBadRequest, "invalid password")
		return
	}
	tokenstring, err := app.generateJWT(user.Username, time.Now().Add(time.Hour*24))
	if err != nil {
		log.Println(err.Error())
		writeJsonError(w, http.StatusInternalServerError, err.Error())
		return
	}
	refreshToken, err := app.generateJWT(user.Username, time.Now().Add(time.Hour*24*7))
	if err != nil {
		log.Println(err.Error())
		writeJsonError(w, http.StatusInternalServerError, err.Error())
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

	http.Redirect(w, r, "http://localhost:3000/profile", http.StatusSeeOther)
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
		writeJsonError(w, http.StatusBadRequest, err.Error())
		return
	}
	defer r.Body.Close()
	var profile models.ProfileInfo

	validator := godantic.Validate{}
	err = validator.BindJSON(body, &profile)
	if err != nil {
		writeJsonError(w, http.StatusBadRequest, err.Error())
		return
	}
	if profile.Latitude < -90 || profile.Latitude > 90 {
		writeJsonError(w, http.StatusBadRequest, "invalid latitude")
		return
	}
	if profile.Longitude < -180 || profile.Longitude > 180 {
		writeJsonError(w, http.StatusBadRequest, "invalid longitude")
		return
	}
	err = app.models.InsertProfileInfo(r.Context(), &profile)
	if err != nil {
		writeJsonError(w, http.StatusInternalServerError, err.Error())
		return
	}
	err = app.models.UpdateUserCompleted(r.Context(), profile.ID)
	if err != nil {
		writeJsonError(w, http.StatusInternalServerError, err.Error())
		return
	}
	w.WriteHeader(http.StatusCreated)
}

func (app *aplication) ImageEndpoint(w http.ResponseWriter, r *http.Request) {

	r.Body = http.MaxBytesReader(w, r.Body, 1048576)
	err := r.ParseMultipartForm(1048576)
	if err != nil {
		log.Println("error parseando multipart-form", err)
		writeJsonError(w, http.StatusBadRequest, err.Error())
		return
	}
	userIDStr := r.FormValue("user_id")
	pictureNumberStr := r.FormValue("picture_number")
	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		writeJsonError(w, http.StatusBadRequest, "Invalid user ID")
		return
	}
	pictureNumber, err := strconv.Atoi(pictureNumberStr)
	if err != nil || pictureNumber < 1 || pictureNumber > 5 {
		writeJsonError(w, http.StatusBadRequest, "Invalid picture number")
		return
	}
	file, header, err := r.FormFile("image")
	if err != nil {
		log.Println("error retieving the file", err)
		writeJsonError(w, http.StatusBadRequest, err.Error())
		return
	}
	defer file.Close()

	valid, err := app.models.ValidateImage(file)
	if !valid || err != nil {
		log.Println("invalid image", err)
		writeJsonError(w, http.StatusBadRequest, err.Error())
		return
	}
	extension := filepath.Ext(header.Filename)
	fileName, err := app.models.GenerateFileName(extension)
	if err != nil {
		log.Println("error generando el file name", err)
		writeJsonError(w, http.StatusInternalServerError, err.Error())
		return
	}
	fileURI, err := app.models.SaveFile(file, fileName)
	if err != nil {
		log.Println("error saving file", err)
		writeJsonError(w, http.StatusInternalServerError, err.Error())
		return
	}

	err = app.models.InsertImage(r.Context(), userID, pictureNumber, fileURI)
	if err != nil {
		log.Println("error inserting image to db", err)
		writeJsonError(w, http.StatusInternalServerError, err.Error())
		return
	}
	response := uploadResponse{
		Message: "Image uploaded successfully",
		URL:     fileURI,
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func (app *aplication) ResetPassword(w http.ResponseWriter, r *http.Request) {
	var email PasswordReset

	body, err := io.ReadAll(r.Body)
	if err != nil {
		writeJsonError(w, http.StatusBadRequest, err.Error())
		return
	}
	defer r.Body.Close()

	validator := godantic.Validate{}
	err = validator.BindJSON(body, &email)
	if err != nil {
		writeJsonError(w, http.StatusBadRequest, err.Error())
		return
	}
	var sender EmailSender
	sender.Destiantion = email.Email
	err = app.models.UpdateUser(r.Context(), sender.generateValidationURI(), email.Email)
	if err != nil {
		writeJsonError(w, http.StatusInternalServerError, err.Error())
		return
	}
	err = sender.sendValidationEmail("Reset your password", "reset", "Click the link below to reset your password")
	if err != nil {
		writeJsonError(w, http.StatusInternalServerError, err.Error())
		return
	}
	w.WriteHeader(http.StatusCreated)
}

func (app *aplication) updatePassword(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		writeJsonError(w, http.StatusBadRequest, err.Error())
		return
	}
	defer r.Body.Close()
	validator := godantic.Validate{}
	var newPassword NewPassword
	err = validator.BindJSON(body, &newPassword)
	if err != nil {
		log.Println(err.Error())
		writeJsonError(w, http.StatusBadRequest, err.Error())
		return
	}
	log.Println(newPassword.Password)
	log.Println(newPassword.Code)
	err = app.models.UpdatePassword(r.Context(), newPassword.Code, newPassword.Password)
	if err != nil {
		if err.Error() == "code not found" {
			writeJsonError(w, 404, "invalid code")
			return
		}
		writeJsonError(w, http.StatusInternalServerError, err.Error())
		return
	}
	http.Redirect(w, r, "/login", http.StatusSeeOther)
}

func writeJsonError(w http.ResponseWriter, statusCode int, message string) {
	response := ErrorResponse{
		Code:    statusCode,
		Message: message,
	}
	jsonResponse, err := json.Marshal(response)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	w.Write(jsonResponse)
}
