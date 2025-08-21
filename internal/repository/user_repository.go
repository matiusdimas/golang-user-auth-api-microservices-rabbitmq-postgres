package repository

import (
	"User-api/internal/models"
	"database/sql"
	"time"
)

type UserRepository interface {
	CreateUser(user *models.User) error
	GetUserByEmail(email string) (*models.User, error)
	GetUserByID(id string) (*models.User, error)
}

type userRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) CreateUser(user *models.User) error {
	query := `INSERT INTO users (email, name, password_hash, created_at, updated_at) 
	          VALUES ($1, $2, $3, $4, $5) RETURNING id`
	
	now := time.Now()
	err := r.db.QueryRow(query, user.Email, user.Name, user.PasswordHash, now, now).Scan(&user.ID)
	if err != nil {
		return err
	}
	user.CreatedAt = now
	user.UpdatedAt = now
	
	return nil
}

func (r *userRepository) GetUserByEmail(email string) (*models.User, error) {
	query := `SELECT id, email, name, password_hash, created_at, updated_at 
	          FROM users WHERE email = $1`
	
	user := &models.User{}
	err := r.db.QueryRow(query, email).Scan(
		&user.ID, &user.Email, &user.Name, &user.PasswordHash, &user.CreatedAt, &user.UpdatedAt,
	)
	
	if err == sql.ErrNoRows {
		return nil, nil
	}
	
	return user, err
}

func (r *userRepository) GetUserByID(id string) (*models.User, error) {
	query := `SELECT id, email, name, password_hash, created_at, updated_at 
	          FROM users WHERE id = $1`
	
	user := &models.User{}
	err := r.db.QueryRow(query, id).Scan(
		&user.ID, &user.Email, &user.Name, &user.PasswordHash, &user.CreatedAt, &user.UpdatedAt,
	)
	
	if err == sql.ErrNoRows {
		return nil, nil
	}
	
	return user, err
}