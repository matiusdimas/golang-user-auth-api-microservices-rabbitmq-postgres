package handlers

import (
	"User-api/internal/utils"
	"encoding/json"
	"net/http"
	"strings"
)

type VerifyHandler struct {
	jwtSecret string
}

func NewVerifyHandler(jwtSecret string) *VerifyHandler {
	return &VerifyHandler{
		jwtSecret: jwtSecret,
	}
}

type VerifyResponse struct {
	Error   bool   `json:"error"`
	Message string `json:"message"`
	Valid   bool   `json:"valid"`
	UserID  string `json:"user_id,omitempty"`
	Email   string `json:"email,omitempty"`
	Name   	string `json:"name,omitempty"`
}

func (h *VerifyHandler) VerifyToken(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	cookie, err := r.Cookie("auth_token")
	var tokenString string

	if err == nil {
		tokenString = cookie.Value
	} else {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(VerifyResponse{
				Error:   true,
				Message: "No authentication token found",
				Valid:   false,
			})
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(VerifyResponse{
				Error:   true,
				Message: "Invalid authorization header format",
				Valid:   false,
			})
			return
		}

		tokenString = parts[1]
	}

	claims, err := utils.ValidateJWT(tokenString, h.jwtSecret)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(VerifyResponse{
			Error:   true,
			Message: "Invalid or expired token: " + err.Error(),
			Valid:   false,
		})
		return
	}

	userID, _ := claims["user_id"].(string)
	email, _ := claims["email"].(string)
	name, _ := claims["name"].(string)

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(VerifyResponse{
		Error:   false,
		Message: "Token is valid.",
		Valid:   true,
		UserID:  userID,
		Email:   email,
		Name:   name,
	})
}