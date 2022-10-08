package cors

import (
	"net/http"
)

func EnableCors(w *http.ResponseWriter, r *http.Request) {
	header := (*w).Header()

	allowList := map[string]bool{
		"http://www.hugopukito.com": true,
		"http://hugopukito.com":     true,
		"http://localhost:8082":     true,
	}

	if origin := r.Header.Get("Origin"); allowList[origin] {
		header.Add("Access-Control-Allow-Origin", origin)
	}

	header.Add("Access-Control-Allow-Headers", "Authorization")

	if r.Method == "OPTIONS" {
		(*w).Header().Add("Access-Control-Max-Age", "3600")
		(*w).WriteHeader(http.StatusOK)
		return
	}
}
