package handler

import (
	"log"
	"net/http"
	"time"
)

const validAPIKey = "my-secret-key-123"

// LoggingMiddleware
func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		ww := &responseWriter{w, http.StatusOK}
		next.ServeHTTP(ww, r)
		log.Printf("[%s] %s %s %d %s",
			time.Now().Format(time.RFC3339),
			r.Method,
			r.URL.Path,
			ww.statusCode,
			time.Since(start),
		)
	})
}

// APIKeyMiddleware
func APIKeyMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		key := r.Header.Get("X-API-KEY")
		if key != validAPIKey {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		next.ServeHTTP(w, r)
	})
}

type responseWriter struct {
	http.ResponseWriter
	statusCode int
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}
