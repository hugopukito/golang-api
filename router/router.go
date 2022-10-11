package router

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"module.com/webServer/service"
)

func InitRouter() {
	router := mux.NewRouter()

	router.HandleFunc("/messages", service.GetMessages).Methods("GET")
	router.HandleFunc("/messages", service.PostMessage).Methods("POST", "OPTIONS")

	router.HandleFunc("/signup", service.SignUp).Methods("POST", "OPTIONS")
	router.HandleFunc("/signin", service.SignIn).Methods("POST", "OPTIONS")

	router.HandleFunc("/socket", service.WsEndpoint)

	log.Fatal(http.ListenAndServe(":8080", router))
}
