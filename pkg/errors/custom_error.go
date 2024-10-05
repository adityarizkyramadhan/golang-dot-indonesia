package custom_error

import (
	"fmt"
	"strings"
)

type CustomError struct {
	Key     string
	Message string
}

type ParsingError struct {
	Key        string `json:"key"`
	StatusCode int    `json:"status_code"`
	Message    string `json:"message"`
}

const (
	ErrInternalServer = "internal_server_error"
	ErrNotFound       = "not_found"
	ErrConflict       = "conflict"
	ErrBadRequest     = "bad_request"
	ErrUnauthorized   = "unauthorized"
	ErrForbidden      = "forbidden"
	ErrValidation     = "validation"
	ErrUnknown        = "unknown"
)

func (e *CustomError) Error() string {
	return fmt.Sprintf("%s: %s", e.Key, e.Message)
}

func NewError(key, message string) error {
	return &CustomError{
		Key:     key,
		Message: message,
	}
}

func ParseError(err error) ParsingError {

	splitStr := strings.Split(err.Error(), ":")
	if len(splitStr) < 2 {
		return ParsingError{
			Key:        ErrUnknown,
			Message:    err.Error(),
			StatusCode: 500,
		}
	}

	key := strings.TrimSpace(splitStr[0])
	message := strings.TrimSpace(splitStr[1])

	statusCode := mapErrorKeyToStatusCode(key)

	return ParsingError{
		Key:        key,
		Message:    message,
		StatusCode: statusCode,
	}
}

func mapErrorKeyToStatusCode(key string) int {
	switch key {
	case ErrInternalServer:
		return 500
	case ErrNotFound:
		return 404
	case ErrConflict:
		return 409
	case ErrBadRequest:
		return 400
	case ErrUnauthorized:
		return 401
	case ErrForbidden:
		return 403
	case ErrValidation:
		return 422
	default:
		return 500
	}
}
