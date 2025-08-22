package handlers

import (
	"User-api/internal/messaging"
	"User-api/internal/models"
	"User-api/internal/utils"
	"encoding/json"
	"net/http"
)

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	getUserFunc := func(userID string) (*models.User, error) {
    	return h.userRepo.GetUserByID(userID)
	}

	isLoggedIn, err := utils.CheckAuth(r, w, h.jwtSecret, getUserFunc, h.cfg.IsProduction)
	if err != nil {
		writeJSONResponse(w, http.StatusInternalServerError, Response{
			Error:   true,
			Message: "Database error",
		})
		return
	}
	if isLoggedIn {
		writeJSONResponse(w, http.StatusForbidden, Response{
			Error:   true,
			Message: "User already logged in. Please logout first.",
		})
		return
	}

	var req models.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSONResponse(w, http.StatusBadRequest, Response{
			Error:   true,
			Message: "Invalid request body",
		})
		return
	}

	validationResult := utils.ValidateLoginRequest(&req)
	if !validationResult.IsValid {
		writeJSONResponse(w, http.StatusBadRequest, Response{
			Error:   true,
			Message: "Validation failed",
			Data:    validationResult.Errors,
		})
		return
	}
	user, err := h.userRepo.GetUserByEmail(req.Email)
	if err != nil {
		writeJSONResponse(w, http.StatusInternalServerError, Response{
			Error:   true,
			Message: "Database error",
		})
		return
	}
	if user == nil {
		writeJSONResponse(w, http.StatusUnauthorized, Response{
			Error:   true,
			Message: "Invalid credentials",
		})
		return
	}
	if !utils.CheckPasswordHash(req.Password, user.PasswordHash) {
		writeJSONResponse(w, http.StatusUnauthorized, Response{
			Error:   true,
			Message: "Invalid credentials",
		})
		return
	}
	token, err := utils.GenerateJWT(user.ID, user.Email, user.Name, h.jwtSecret)
	if err != nil {
		writeJSONResponse(w, http.StatusInternalServerError, Response{
			Error:   true,
			Message: "Failed to generate token",
		})
		return
	}
	utils.SetAuthCookie(w, token, h.cfg.IsProduction)
	go messaging.PublishUserEvent(h.rabbitConn, "user_logged_in", models.AuthResponse{
		ID:        user.ID,
		Email:     user.Email,
		Name:      user.Name,
		CreatedAt: user.CreatedAt,
	})
	responseData := models.AuthResponse{
		ID:        user.ID,
		Email:     user.Email,
		Name:      user.Name,
		CreatedAt: user.CreatedAt,
	}
	writeJSONResponse(w, http.StatusOK, Response{
		Error:   false,
		Message: "Login successful",
		Data:    responseData,
	})
}