package utils

import (
	"User-api/internal/models"

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
func CheckAuth(
    r *http.Request,
    w http.ResponseWriter,
    jwtSecret string,
    getUserByIDFunc func(string) (*models.User, error),
    isProduction bool,
) (bool, error) {
    cookie, err := r.Cookie("auth_token")
    if err != nil {
        return false, nil
    }

    claims, err := ValidateJWT(cookie.Value, jwtSecret)
    if err != nil {
        DeleteAuthCookie(w, isProduction)
        return false, nil
    }

    userID, ok := claims["user_id"].(string)
    if !ok {
        DeleteAuthCookie(w, isProduction)
        return false, nil
    }

    user, err := getUserByIDFunc(userID)
    if err != nil {
        return false, err
    }
    if user == nil {
        DeleteAuthCookie(w, isProduction)
        return false, nil
    }

    return true, nil
}


func DeleteAuthCookie(w http.ResponseWriter, isProduction bool) {
	cookie := &http.Cookie{
		Name:     "auth_token",
		Value:    "",
		Path:     "/",
		HttpOnly: true,
		Secure:   isProduction,
		SameSite: http.SameSiteStrictMode,
		MaxAge:   -1,
	}
	http.SetCookie(w, cookie)
}