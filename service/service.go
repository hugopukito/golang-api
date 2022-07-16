package service

import (
	"encoding/json"
	"fmt"
	"net/http"

	"module.com/webServer/cors"
	"module.com/webServer/db"
	"module.com/webServer/user"
)

func GetUsers(w http.ResponseWriter, r *http.Request) {
	cors.EnableCors(&w, r)
	w.Header().Set("Content-Type", "application/json")
	users := db.GetAll()
	dto, _ := json.Marshal(users)
	fmt.Fprintf(w, string(dto))
}

func PostUserMessage(w http.ResponseWriter, r *http.Request) {
	cors.EnableCors(&w, r)

	var user user.User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	db.InsertUserMessage(user.Name, user.Message)
}
