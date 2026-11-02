package apierror

import (
	"errors"
	"net/http"
)

const (
	unknownErrorMessage  = "Unknown error"
	internalErrorMessage = "Internal server error"
)

// Classify maps err to the APIError rendered to the client. Errors that are
// not an *APIError become a generic 500 so internal details never leak.
func Classify(err error) *APIError {
	var apiErr *APIError
	if !errors.As(err, &apiErr) {
		return &APIError{Code: CodeInternal, Status: http.StatusInternalServerError, Message: internalErrorMessage}
	}

	message := apiErr.Message
	if message == "" {
		message = unknownErrorMessage
	}
	return &APIError{Code: apiErr.Code, Status: apiErr.Status, Message: message}
}
