package cors

import (
	"net/http"
)

func EnableCors(w *http.ResponseWriter, r *http.Request) {

	allowList := map[string]bool{
		"http://77.136.126.254:8081": true,
		"http://192.168.0.27:8081":   true,
		"http://localhost:8081":      true,
	}

	if origin := r.Header.Get("Origin"); allowList[origin] {
		(*w).Header().Set("Access-Control-Allow-Origin", origin)
	}
}
