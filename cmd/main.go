package main

import (
	"User-api/internal/config"
	"User-api/internal/database"
	"User-api/internal/handlers"
	"User-api/internal/middleware"
	"User-api/internal/messaging"
	"User-api/internal/repository"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
)

func main() {
	cfg := config.LoadConfig()
	db, err := database.InitDB(cfg)
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	defer db.Close()
	rabbitConn, err := messaging.InitRabbitMQ(cfg)
	if err != nil {
		log.Fatal("Failed to connect to RabbitMQ:", err)
	}
	defer rabbitConn.Close()
	userRepo := repository.NewUserRepository(db)
	authHandler := handlers.NewAuthHandler(userRepo, cfg.JWTSecret, rabbitConn, cfg)
	verifyHandler := handlers.NewVerifyHandler(cfg.JWTSecret) 
	router := mux.NewRouter()
	router.Use(middleware.CORSRestrictedMiddleware)
	router.HandleFunc("/api/register", authHandler.Register).Methods("POST", "OPTIONS")
	router.HandleFunc("/api/login", authHandler.Login).Methods("POST", "OPTIONS")
	router.HandleFunc("/api/verify", verifyHandler.VerifyToken).Methods("GET", "OPTIONS")
	protectedRouter := router.PathPrefix("/api").Subrouter()
	protectedRouter.Use(middleware.AuthMiddleware(cfg.JWTSecret))
	protectedRouter.HandleFunc("/logout", authHandler.Logout).Methods("DELETE", "OPTIONS")
	router.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status": "ok", "service": "user-api", "environment": "` + cfg.Environment() + `"}`))
	}).Methods("GET")
	log.Printf("Server starting on port %s in %s environment", cfg.Port, cfg.Environment())
	log.Fatal(http.ListenAndServe(":"+cfg.Port, router))
}

func (c *Config) Environment() string {
	if c.IsProduction {
		return "production"
	}
	return "development"
}