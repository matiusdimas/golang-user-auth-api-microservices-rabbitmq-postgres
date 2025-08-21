package middleware

import (
	"net/http"
	"os"
	"strings"
)

func CORSRestrictedMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		allowedOrigins := os.Getenv("ALLOWED_ORIGINS")
		if allowedOrigins == "" {
			allowedOrigins = "http://localhost:3000,http://localhost:8080"
		}

		origins := strings.Split(allowedOrigins, ",")
		requestOrigin := r.Header.Get("Origin")
		var allowedOrigin string
		for _, origin := range origins {
			if origin == requestOrigin {
				allowedOrigin = origin
				break
			}
		}
		if allowedOrigin == "" {
			allowedOrigin = origins[0] 
		}
		w.Header().Set("Access-Control-Allow-Origin", allowedOrigin)
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS, PATCH")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, X-Requested-With, Accept, Origin")
		w.Header().Set("Access-Control-Allow-Credentials", "true")
		w.Header().Set("Access-Control-Max-Age", "86400")
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}
		next.ServeHTTP(w, r)
	})
}