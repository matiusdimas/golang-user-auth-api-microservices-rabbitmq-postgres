package utils

import (
	"fmt"
	"net/mail"
	"regexp"
	"strings"
	"unicode"
	"User-api/internal/models"

	"github.com/microcosm-cc/bluemonday"
)

type ValidationResult struct {
	IsValid bool
	Errors  map[string]string
}
func SanitizeInput(input string) string {
	p := bluemonday.UGCPolicy()
	p.AllowAttrs("class").OnElements("code")
	sanitized := p.Sanitize(input)
	sanitized = strings.TrimSpace(sanitized)
	sanitized = strings.ToValidUTF8(sanitized, "")
	
	return sanitized
}
func ValidateEmail(email string) error {
	email = SanitizeInput(email)
	
	if email == "" {
		return fmt.Errorf("email is required")
	}
	
	if len(email) > 254 {
		return fmt.Errorf("email too long")
	}
	
	_, err := mail.ParseAddress(email)
	if err != nil {
		return fmt.Errorf("invalid email format")
	}
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	if !emailRegex.MatchString(email) {
		return fmt.Errorf("invalid email format")
	}
	return nil
}
func ValidateName(name string) error {
	name = SanitizeInput(name)
	if name == "" {
		return fmt.Errorf("name is required")
	}
	
	if len(name) < 2 {
		return fmt.Errorf("name must be at least 2 characters")
	}
	if len(name) > 100 {
		return fmt.Errorf("name too long")
	}
	nameRegex := regexp.MustCompile(`^[a-zA-Z0-9\s\-'.]+$`)
	if !nameRegex.MatchString(name) {
		return fmt.Errorf("name contains invalid characters")
	}
	
	return nil
}
func ValidatePassword(password string) error {
	if password == "" {
		return fmt.Errorf("password is required")
	}
	
	if len(password) < 8 {
		return fmt.Errorf("password must be at least 8 characters")
	}
	
	if len(password) > 72 { 
		return fmt.Errorf("password too long")
	}
	
	var (
		hasUpper   = false
		hasLower   = false
		hasNumber  = false
		hasSpecial = false
	)
	
	for _, char := range password {
		switch {
		case unicode.IsUpper(char):
			hasUpper = true
		case unicode.IsLower(char):
			hasLower = true
		case unicode.IsDigit(char):
			hasNumber = true
		case unicode.IsPunct(char) || unicode.IsSymbol(char):
			hasSpecial = true
		}
	}
	
	if !hasUpper {
		return fmt.Errorf("password must contain at least one uppercase letter")
	}
	
	if !hasLower {
		return fmt.Errorf("password must contain at least one lowercase letter")
	}
	
	if !hasNumber {
		return fmt.Errorf("password must contain at least one number")
	}
	
	if !hasSpecial {
		return fmt.Errorf("password must contain at least one special character")
	}
	
	return nil
}
func ValidateRegisterRequest(req *models.RegisterRequest) ValidationResult {
	errors := make(map[string]string)
	req.Email = SanitizeInput(req.Email)
	req.Name = SanitizeInput(req.Name)
	if err := ValidateEmail(req.Email); err != nil {
		errors["email"] = err.Error()
	}
	if err := ValidateName(req.Name); err != nil {
		errors["name"] = err.Error()
	}
	if err := ValidatePassword(req.Password); err != nil {
		errors["password"] = err.Error()
	}
	return ValidationResult{
		IsValid: len(errors) == 0,
		Errors:  errors,
	}
}
func ValidateLoginRequest(req *models.LoginRequest) ValidationResult {
	errors := make(map[string]string)
	req.Email = SanitizeInput(req.Email)
	if err := ValidateEmail(req.Email); err != nil {
		errors["email"] = err.Error()
	}
	if req.Password == "" {
		errors["password"] = "password is required"
	} else if len(req.Password) < 1 {
		errors["password"] = "password is required"
	}
	return ValidationResult{
		IsValid: len(errors) == 0,
		Errors:  errors,
	}
}