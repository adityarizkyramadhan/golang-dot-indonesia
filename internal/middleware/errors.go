package middleware

import (
	"fmt"

	custom_error "github.com/adityarizkyramadhan/golang-dot-indonesia/pkg/errors"
	"github.com/adityarizkyramadhan/golang-dot-indonesia/pkg/response"
	"github.com/gin-gonic/gin"
)

func ErrorHandler() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx.Next()
		if len(ctx.Errors) > 0 {
			err := ctx.Errors.Last()
			if err != nil {
				errParse := custom_error.ParseError(err)
				response := response.Error(fmt.Sprintf("%s: %s", errParse.Key, errParse.Message))
				ctx.JSON(errParse.StatusCode, response)
				ctx.Abort()
			}
		}
	}
}
