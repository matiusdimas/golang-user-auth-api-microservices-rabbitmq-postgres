package middleware

import (
	"User-api/internal/utils"
	"encoding/json"
	"net/http"
	"strings"
)

type Response struct {
	Error   bool   `json:"error"`
	Message string `json:"message"`
}

func AuthMiddleware(jwtSecret string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			cookie, err := r.Cookie("auth_token")
			var tokenString string

			if err == nil {
				tokenString = cookie.Value
			} else {
				authHeader := r.Header.Get("Authorization")
				if authHeader == "" {
					w.Header().Set("Content-Type", "application/json")
					w.WriteHeader(http.StatusUnauthorized)
					json.NewEncoder(w).Encode(Response{
						Error:   true,
						Message: "Authorization header missing",
					})
					return
				}

				parts := strings.Split(authHeader, " ")
				if len(parts) != 2 || parts[0] != "Bearer" {
					w.Header().Set("Content-Type", "application/json")
					w.WriteHeader(http.StatusUnauthorized)
					json.NewEncoder(w).Encode(Response{
						Error:   true,
						Message: "Invalid authorization header",
					})
					return
				}

				tokenString = parts[1]
			}

			_, err = utils.ValidateJWT(tokenString, jwtSecret)
			if err != nil {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusUnauthorized)
				json.NewEncoder(w).Encode(Response{
					Error:   true,
					Message: "Invalid token",
				})
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}