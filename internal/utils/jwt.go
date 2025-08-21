package utils

import (
	"github.com/golang-jwt/jwt/v5"
	"time"
	"net/http"
)

func SetAuthCookie(w http.ResponseWriter, token string, isProduction bool) {
	cookie := &http.Cookie{
		Name:     "auth_token",
		Value:    token,
		Path:     "/",
		HttpOnly: true,
		Secure:   isProduction,
		SameSite: http.SameSiteStrictMode,
		Expires:  time.Now().Add(365 * 24 * time.Hour),
	}
	http.SetCookie(w, cookie)
}

func GenerateJWT(userID, email, name, secret string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": userID,
		"email":   email,
		"name": name,
	})

	return token.SignedString([]byte(secret))
}

func ValidateJWT(tokenString, secret string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(secret), nil
	}, jwt.WithoutClaimsValidation())

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, jwt.ErrSignatureInvalid
}

func GetUserIDFromToken(tokenString, secret string) (string, error) {
	claims, err := ValidateJWT(tokenString, secret)
	if err != nil {
		return "", err
	}

	userID, ok := claims["user_id"].(string)
	if !ok {
		return "", jwt.ErrTokenInvalidId
	}

	return userID, nil
}