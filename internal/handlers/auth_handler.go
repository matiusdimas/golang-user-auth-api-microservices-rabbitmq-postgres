package handlers

import (
	"User-api/internal/repository"
	"User-api/internal/config"
	"encoding/json"
	"net/http"

	amqp "github.com/rabbitmq/amqp091-go"
)

type Response struct {
	Error   bool        `json:"error"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

type AuthHandler struct {
	userRepo   repository.UserRepository
	jwtSecret  string
	rabbitConn *amqp.Connection
	cfg        *config.Config 
}

func NewAuthHandler(userRepo repository.UserRepository, jwtSecret string, rabbitConn *amqp.Connection, cfg *config.Config) *AuthHandler {
	return &AuthHandler{
		userRepo:   userRepo,
		jwtSecret:  jwtSecret,
		rabbitConn: rabbitConn,
		cfg:        cfg, 
	}
}

func writeJSONResponse(w http.ResponseWriter, statusCode int, response Response) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(response)
}


