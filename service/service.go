package service

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"golang.org/x/crypto/bcrypt"
	"module.com/webServer/cors"
	"module.com/webServer/db"
	"module.com/webServer/entity"
)

func GetMessages(w http.ResponseWriter, r *http.Request) {
	cors.EnableCors(&w, r)
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
}

func SignUp(w http.ResponseWriter, r *http.Request) {
	cors.EnableCors(&w, r)

	var user entity.User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	db.FindUser(user)

	pwd, err := GeneratehashPassword(user.Password)
	if err != nil {
		log.Fatalln("error in password hash")
	}

	fmt.Fprintf(w, pwd)
}

func SignIn(w http.ResponseWriter, r *http.Request) {

}

func GeneratehashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}
