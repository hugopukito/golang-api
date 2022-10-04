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

	w.Header().Set("Content-Type", "application/json")
	messages := db.GetAllMessages()
	dto, err := json.Marshal(messages)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.Write([]byte(string(dto)))
}

func Test(w http.ResponseWriter, r *http.Request) {
	cors.EnableCors(&w, r)

	temp := r.Header["Authorization"]
	var token string
	for i := 0; i < len(temp); i++ {
		token = temp[0]
	}
	fmt.Println(token)
	// parseJwt(w, token)
	/*
		email := retrieveEmail(claims)
		user := retrieveUserWithEmail(email)

		var message entity.Message
		json.NewDecoder(r.Body).Decode(&message)
		fmt.Println(message.Message)
		fmt.Println(user.Email) */
}

func PostMessage(w http.ResponseWriter, r *http.Request) {
	cors.EnableCors(&w, r)

	claims := parseJwt(w, r.Header["Authorization"][0])

	email := retrieveEmail(claims)
	user := retrieveUserWithEmail(email)

	var message entity.Message
	err := json.NewDecoder(r.Body).Decode(&message)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	message.Name = user.Name

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
	token.TokenString = validToken
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(token)
}

func retrieveUserWithEmail(email string) entity.User {
	return db.FindUser(email)
}
