package utils_jwt

import (
	"os"
	"time"

	custom_error "github.com/adityarizkyramadhan/golang-dot-indonesia/pkg/errors"
	"github.com/golang-jwt/jwt/v5"
)

func GenerateToken(id int) (string, error) {
	secretKey := os.Getenv("SECRET_KEY")
	if secretKey == "" {
		return "", custom_error.NewError(custom_error.ErrInternalServer, "secret key not found")
	}

	claims := jwt.MapClaims{
		"id":  id,
		"exp": time.Now().Add(time.Hour * 72).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString([]byte(secretKey))
	if err != nil {
		return "", custom_error.NewError(custom_error.ErrInternalServer, err.Error())
	}

	return tokenString, nil
}
