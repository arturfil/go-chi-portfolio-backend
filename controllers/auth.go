package controllers

import (
	"encoding/json"
	"errors"
	"net/http"
	"os"
	"portfolio-api/helpers"
	"portfolio-api/models"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

type jsonResponse struct {
	Error   bool        `json:"error"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

type tokenResponse struct {
	Token string       `json:"token"`
	User  *models.User `json:"user"`
}

// type envelope map[string]interface{}

// this function first will check that there are no users already signed in
// if there are none, then it will create the ser with hashed password
func Singup(w http.ResponseWriter, r *http.Request) {
	var u models.User
	err := json.NewDecoder(r.Body).Decode(&u)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	helpers.WriteJSON(w, http.StatusOK, u)
	id, err := mod.User.Signup(u)
	if err != nil {
		h.Errorlog.Panicln(err)
		helpers.ErrorJSON(w, err, http.StatusForbidden)
		helpers.MessageLogs.Errorlog.Println("Got back id of", id)
		newUser, _ := mod.User.GetUserById(id)
		helpers.WriteJSON(w, http.StatusOK, newUser)
	}
}

// this method will check that the credentials against the key are
// equal
func Login(w http.ResponseWriter, r *http.Request) {
	var myKey = []byte(os.Getenv("SECRET_KEY"))
	type credentials struct {
		Username string `json:"email"`
		Password string `json:"password"`
	}
	// setup creds & jsonResponse
	var creds credentials
	var payload jsonResponse
	// read Json
	err := helpers.ReadJSON(w, r, &creds)
	if err != nil {
		helpers.MessageLogs.Errorlog.Println(err)
		payload.Error = true
		payload.Message = "Invalid json supplied"
		_ = helpers.WriteJSON(w, http.StatusBadRequest, payload)
	}
	// get user if creds are valid
	user, err := mod.User.GetByEmail(creds.Username)
	if err != nil {
		helpers.ErrorJSON(w, errors.New("invalid no user found"))
		return
	}
	// check if valid
	validPassword, err := user.PasswordMatches(creds.Password)
	if err != nil || !validPassword {
		helpers.ErrorJSON(w, errors.New("invalide username/password"))
		return
	}
	// create new token
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["authorized"] = true
	claims["name"] = user.FirstName
	claims["email"] = user.Email
	claims["exp"] = time.Now().Add(time.Minute * 60 * 4).Unix()

	tokenString, err := token.SignedString(myKey)
	if err != nil {
		helpers.MessageLogs.Errorlog.Println(err)
		helpers.ErrorJSON(w, err)
		return
	}
	user.Password = "hidden"
	// create response
	response := tokenResponse{
		Token: tokenString,
		User:  user,
	}
	// send response if no errors
	err = helpers.WriteJSON(w, http.StatusOK, response)
	if err != nil {
		helpers.MessageLogs.Errorlog.Println(err)
	}
}
