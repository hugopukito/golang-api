package cors

import (
	"net/http"
)

func EnableCors(w *http.ResponseWriter, r *http.Request) {

	allowList := map[string]bool{
		"http://www.hugopukito.com": true,
		"http://hugopukito.com":     true,
	}

	if origin := r.Header.Get("Origin"); allowList[origin] {
		(*w).Header().Set("Access-Control-Allow-Origin", origin)
	}
}
