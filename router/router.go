package router

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"module.com/webServer/service"
)

func InitRouter() {
	router := mux.NewRouter()

	router.HandleFunc("/users", service.GetUsers).Methods("GET")
	router.HandleFunc("/users", service.PostUserMessage).Methods("POST")

	log.Fatal(http.ListenAndServe(":8080", router))
}
