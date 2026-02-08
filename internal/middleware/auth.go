package middleware

import (
	"encoding/json"
	"log"
	"net/http"
	"time"
)

const APIKey = "secret12345"

func Auth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("X-API-KEY") != APIKey {
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(map[string]string{
				"error": "unauthorized",
			})
			return
		}
		next.ServeHTTP(w, r)
	})
}

func Logger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf(
			"%s %s %s server-log",
			time.Now().Format(time.RFC3339),
			r.Method,
			r.URL.Path,
		)
		next.ServeHTTP(w, r)
	})
}
