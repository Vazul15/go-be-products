package apierror

import (
	"fmt"
	"net/http"
)

// TODO: rename the codes and figure out new ones
const (
	CodeAccessDenied = "FAT_0114"
	CodeGeneric      = "FAT_0202"
	CodeNotFound     = "NOT_FOUND"
	CodeConflict     = "CONFLICT"
	CodeInternal     = "INTERNAL_ERROR"
)

var (
	ErrAccessDenied = &APIError{Code: CodeAccessDenied, Status: http.StatusForbidden}
	ErrGeneric      = &APIError{Code: CodeGeneric, Status: http.StatusBadRequest}
	ErrNotFound     = &APIError{Code: CodeNotFound, Status: http.StatusNotFound}
	ErrConflict     = &APIError{Code: CodeConflict, Status: http.StatusConflict}
)

type APIError struct {
	Code    string
	Status  int
	Message string
}

// Error returns the string representation of the APIError.
func (e *APIError) Error() string {
	return fmt.Sprintf("%s: %s", e.Code, e.Message)
}

// Is checks if the target error is an APIError with the same code as the receiver.
func (e *APIError) Is(target error) bool {
	t, ok := target.(*APIError)
	return ok && e.Code == t.Code
}

// AccessDenied creates a new APIError representing an access denied error.
func AccessDenied(message string) error {
	return &APIError{Code: CodeAccessDenied, Status: http.StatusForbidden, Message: message}
}

// Generic creates a new APIError representing a generic error.
func Generic(message string) error {
	return &APIError{Code: CodeGeneric, Status: http.StatusBadRequest, Message: message}
}

// NotFound creates a new APIError representing a missing resource.
func NotFound(message string) error {
	return &APIError{Code: CodeNotFound, Status: http.StatusNotFound, Message: message}
}

// Conflict creates a new APIError representing a uniqueness or state conflict.
func Conflict(message string) error {
	return &APIError{Code: CodeConflict, Status: http.StatusConflict, Message: message}
}
