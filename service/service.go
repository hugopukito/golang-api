package service

import (
	"encoding/json"
	"log"
	"net/http"

	"module.com/webServer/cors"
	"module.com/webServer/db"
	"module.com/webServer/entity"
)

var user entity.User

func SignUp(w http.ResponseWriter, r *http.Request) {
	cors.EnableCors(&w, r)

	var user entity.User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if retrieveUserWithEmail(user.Email).Email != "" {
		w.WriteHeader(http.StatusConflict)
		return
	}

	pwd, err := generateHashPassword(user.Password)
	if err != nil {
		log.Fatalln("error in password hash")
	}

	user.Password = pwd
	db.InsertUser(user)
	w.WriteHeader(http.StatusCreated)
}

func SignIn(w http.ResponseWriter, r *http.Request) {
	cors.EnableCors(&w, r)

	var authdetails entity.Authentication
	err := json.NewDecoder(r.Body).Decode(&authdetails)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	authUser := db.FindUser(authdetails.Email)
	if authUser.Email == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	check := checkPasswordHash(authdetails.Password, authUser.Password)

	if !check {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	validToken, err := generateJWT(authUser.Email)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var token entity.Token
	token.Email = authUser.Email
	token.Name = authUser.Name
	token.TokenString = validToken
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(token)
}

func retrieveUserWithEmail(email string) entity.User {
	return db.FindUser(email)
}

func parseAuthorization(w http.ResponseWriter, r *http.Request) {
	token := r.Header["Authorization"]
	if len(token) > 0 {
		claims := parseJwt(w, token[0])
		user.Email = retrieveEmail(claims)
		user.Name = retrieveUserWithEmail(user.Email).Name
	}
}
