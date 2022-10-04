package cors

import (
	"net/http"
)

func EnableCors(w *http.ResponseWriter, r *http.Request) {

	allowList := map[string]bool{
		"http://www.hugopukito.com": true,
		"http://hugopukito.com":     true,
		"http://localhost:8082":     true,
	}

	if origin := r.Header.Get("Origin"); allowList[origin] {
		(*w).Header().Set("Access-Control-Allow-Origin", origin)
	}
	(*w).Header().Set("Access-Control-Allow-Headers", "Authorization")
}
