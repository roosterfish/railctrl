package api

import "net/http"

type APIError struct {
	Code    int
	Message string
}

func (e *APIError) Error() string {
	return e.Message
}

func BadRequest(err error) error {
	return &APIError{
		Code:    http.StatusBadRequest,
		Message: err.Error(),
	}
}

func NotFound(err error) error {
	return &APIError{
		Code:    http.StatusNotFound,
		Message: err.Error(),
	}
}

func Internal(err error) error {
	return &APIError{
		Code:    http.StatusInternalServerError,
		Message: err.Error(),
	}
}
