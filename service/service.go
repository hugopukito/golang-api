package service

import (
	"encoding/json"
	"net/http"

	"module.com/webServer/cors"
	"module.com/webServer/db"
)

func GetUsers(w http.ResponseWriter, r *http.Request) {
	cors.EnableCors(&w, r)
	w.Header().Set("Content-Type", "application/json")
	users := db.GetAll()
	json.NewEncoder(w).Encode(users)
}
