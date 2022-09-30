package service

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"module.com/webServer/cors"
	"module.com/webServer/db"
	"module.com/webServer/entity"
)

func GetMessages(w http.ResponseWriter, r *http.Request) {
	cors.EnableCors(&w, r)

	ParseJwt(w, r.Header["Authorization"][0])

	w.Header().Set("Content-Type", "application/json")
	messages := db.GetAllMessages()
	dto, _ := json.Marshal(messages)
	fmt.Fprintf(w, string(dto))
}

func PostMessage(w http.ResponseWriter, r *http.Request) {
	cors.EnableCors(&w, r)

	var message entity.Message
	err := json.NewDecoder(r.Body).Decode(&message)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	db.InsertMessage(message)
	w.WriteHeader(http.StatusCreated)
}

func SignUp(w http.ResponseWriter, r *http.Request) {
	cors.EnableCors(&w, r)

	var user entity.User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if db.FindUser(user.Email).Email != "" {
		w.WriteHeader(http.StatusConflict)
		return
	}

	pwd, err := GeneratehashPassword(user.Password)
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

	check := CheckPasswordHash(authdetails.Password, authUser.Password)

	if !check {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	validToken, err := GenerateJWT(authUser.Email, authUser.Role)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var token entity.Token
	token.Email = authUser.Email
	token.Role = authUser.Role
	token.TokenString = validToken
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(token)
}
