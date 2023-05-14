package router

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"module.com/webServer/game"
	"module.com/webServer/service"
)

func InitRouter() {
	router := mux.NewRouter()

	router.HandleFunc("/signup", service.SignUp).Methods("POST", "OPTIONS")
	router.HandleFunc("/signin", service.SignIn).Methods("POST", "OPTIONS")

	router.HandleFunc("/chat", service.HandleChatConnections)

	router.HandleFunc("/game", game.HandleGameConnections)

	log.Fatal(http.ListenAndServe(":8080", router))
}
