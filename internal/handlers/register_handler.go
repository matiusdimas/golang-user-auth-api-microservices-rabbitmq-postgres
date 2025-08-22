package handlers

import (
	"User-api/internal/messaging"
	"User-api/internal/models"
	"User-api/internal/utils"
	"encoding/json"
	"net/http"
)

func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
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
			Message: "User already logged in. Please logout first",
		})
		return
	}

	var req models.RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSONResponse(w, http.StatusBadRequest, Response{
			Error:   true,
			Message: "Invalid request body",
		})
		return
	}

	validationResult := utils.ValidateRegisterRequest(&req)
	if !validationResult.IsValid {
		writeJSONResponse(w, http.StatusBadRequest, Response{
			Error:   true,
			Message: "Validation failed",
			Data:    validationResult.Errors,
		})
		return
	}

	existingUser, err := h.userRepo.GetUserByEmail(req.Email)
	if err != nil {
		writeJSONResponse(w, http.StatusInternalServerError, Response{
			Error:   true,
			Message: "Database error",
		})
		return
	}
	if existingUser != nil {
		writeJSONResponse(w, http.StatusConflict, Response{
			Error:   true,
			Message: "User already exists",
		})
		return
	}

	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		writeJSONResponse(w, http.StatusInternalServerError, Response{
			Error:   true,
			Message: "Failed to hash password",
		})
		return
	}

	user := &models.User{
		Email:        req.Email,
		Name:         req.Name,
		PasswordHash: hashedPassword,
	}

	if err := h.userRepo.CreateUser(user); err != nil {
		writeJSONResponse(w, http.StatusInternalServerError, Response{
			Error:   true,
			Message: "Failed to create user",
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

	go messaging.PublishUserEvent(h.rabbitConn, "user_created", models.AuthResponse{
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

	writeJSONResponse(w, http.StatusCreated, Response{
		Error:   false,
		Message: "User registered successfully",
		Data:    responseData,
	})
}