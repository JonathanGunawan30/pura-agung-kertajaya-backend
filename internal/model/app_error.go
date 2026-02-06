package model

import "net/http"

type ResponseError struct {
	Message string
	Code    int
}

func (e *ResponseError) Error() string {
	return e.Message
}

func NewError(code int, message string) *ResponseError {
	return &ResponseError{
		Code:    code,
		Message: message,
	}
}

var (
	ErrNotFound     = func(msg string) *ResponseError { return NewError(http.StatusNotFound, msg) }
	ErrConflict     = func(msg string) *ResponseError { return NewError(http.StatusConflict, msg) }
	ErrBadRequest   = func(msg string) *ResponseError { return NewError(http.StatusBadRequest, msg) }
	ErrUnauthorized = func(msg string) *ResponseError { return NewError(http.StatusUnauthorized, msg) }
	ErrInternal     = func(msg string) *ResponseError { return NewError(http.StatusInternalServerError, msg) }
	ErrForbidden    = func(msg string) *ResponseError {
		return NewError(http.StatusForbidden, msg)
	}
)
