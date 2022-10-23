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

	router.HandleFunc("/image/{id}", service.GetImage).Methods("GET")

	router.HandleFunc("/signup", service.SignUp).Methods("POST", "OPTIONS")
	router.HandleFunc("/signin", service.SignIn).Methods("POST", "OPTIONS")

	router.HandleFunc("/websocket", service.HandleChatConnections)

	log.Fatal(http.ListenAndServe(":8080", router))
}
