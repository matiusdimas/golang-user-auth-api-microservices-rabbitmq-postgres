package handlers

import (
	"User-api/internal/messaging"
	"User-api/internal/models"
	"User-api/internal/utils"
	"net/http"
	"strings"
	"time"
)

func (h *AuthHandler) Logout(w http.ResponseWriter, r *http.Request) {
	var tokenString string
	if cookie, err := r.Cookie("auth_token"); err == nil {
		tokenString = cookie.Value
	} else {
		authHeader := r.Header.Get("Authorization")
		if authHeader != "" {
			parts := strings.Split(authHeader, " ")
			if len(parts) == 2 && parts[0] == "Bearer" {
				tokenString = parts[1]
			}
		}
	}
	http.SetCookie(w, &http.Cookie{
		Name:     "auth_token",
		Value:    "",
		Path:     "/",
		HttpOnly: true,
		Secure:   false,
		SameSite: http.SameSiteStrictMode,
		Expires:  time.Now().Add(-time.Hour),
		MaxAge:   -1,
	})
	if tokenString != "" {
		claims, err := utils.ValidateJWT(tokenString, h.jwtSecret)
		if err == nil {
			if userID, ok := claims["user_id"].(string); ok {
				user, err := h.userRepo.GetUserByID(userID)
				if err == nil && user != nil {
					// Publish user logged out event
					go messaging.PublishUserEvent(h.rabbitConn, "user_logged_out", models.AuthResponse{
						ID:        user.ID,
						Email:     user.Email,
						Name:      user.Name,
						CreatedAt: user.CreatedAt,
					})
				}
			}
		}
	}
	writeJSONResponse(w, http.StatusOK, Response{
		Error:   false,
		Message: "Logged out successfully",
	})
}