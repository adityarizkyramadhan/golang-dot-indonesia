package middleware

import (
	"os"
	"strings"

	custom_error "github.com/adityarizkyramadhan/golang-dot-indonesia/pkg/errors"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func JWTMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		token := ctx.GetHeader("Authorization")
		if token == "" {
			err := custom_error.NewError(custom_error.ErrUnauthorized, "token not found")
			ctx.Error(err).SetType(gin.ErrorTypePublic)
			ctx.Abort()
			return
		}

		token = strings.Replace(token, "Bearer ", "", 1)
		claims, err := verifyJWT(token)
		if err != nil {
			ctx.Error(err).SetType(gin.ErrorTypePublic)
			ctx.Abort()
			return
		}
		mapClaims, ok := claims.(jwt.MapClaims)
		if !ok {
			err := custom_error.NewError(custom_error.ErrInternalServer, "unexpected claims")
			ctx.Error(err).SetType(gin.ErrorTypePrivate)
			ctx.Abort()
			return
		}
		// mapClaims["id"] buat jadi int
		idInt := int(mapClaims["id"].(float64))
		ctx.Set("id", idInt)
		ctx.Next()
	}
}

func verifyJWT(tokenString string) (interface{}, error) {
	secretKey := os.Getenv("SECRET_KEY")
	if secretKey == "" {
		return nil, custom_error.NewError(custom_error.ErrInternalServer, "secret key not found")
	}
	jwtSecret := []byte(secretKey)
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, custom_error.NewError(custom_error.ErrInternalServer, "unexpected signing method")
		}
		return jwtSecret, nil
	})

	if err != nil {
		return nil, custom_error.NewError(custom_error.ErrInternalServer, err.Error())
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, custom_error.NewError(custom_error.ErrInternalServer, "invalid token")
}
