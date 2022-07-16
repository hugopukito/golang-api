package router

import (
	"net/http"

	"module.com/webServer/service"
)

func InitRouter() {
	http.HandleFunc("/", service.GetUsers)
	http.ListenAndServe(":8080", nil)
}
