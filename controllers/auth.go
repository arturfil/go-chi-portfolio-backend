package controllers

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"portfolio-api/helpers"
	"portfolio-api/models"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

// this struct is fine to be here since we only use it here
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
	var payload models.JsonResponse
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
		helpers.ErrorJSON(w, errors.New("invalid, no user found"))
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

// GetUserByToken - this method will serve to get the user once a token is provided. This
// is really usefull when ever you want to preserve login state in the front-end
func GetUserByToken(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Reached [GetUserByToken] method")
	var myKey = []byte(os.Getenv("SECRET_KEY"))
	type TokenClaim struct {
		Authorized bool   `json:"authorized"`
		Email      string `json:"email"`
		Exp        int    `json:"exp"`
		Name       string `json:"name"`
		jwt.StandardClaims
	}

	claims := &TokenClaim{}
	token, err := jwt.ParseWithClaims(r.Header["Authorization"][0], claims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("there was an error")
		}
		return myKey, nil
	})
	if token.Valid {
		if err != nil {
			helpers.MessageLogs.Infolog.Print(err)
		}

		if err != nil {
			helpers.MessageLogs.Infolog.Print(err)
			return
		}

		user, err := mod.User.GetByEmail(claims.Email)
		if err != nil {
			helpers.MessageLogs.Infolog.Print(err)
		}
		user.Password = ""
		helpers.WriteJSON(w, http.StatusOK, user)
		return
	}
}
