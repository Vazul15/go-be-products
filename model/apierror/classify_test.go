package apierror

import (
	"errors"
	"fmt"
	"net/http"
	"testing"
)

func TestClassify_AccessDeniedIdentity_PreservesCode(t *testing.T) {
	err := AccessDenied("Access denied: missing permission 'x'")

	got := Classify(err)

	if got.Code != CodeAccessDenied {
		t.Errorf("Code = %q, want %q", got.Code, CodeAccessDenied)
	}
	if got.Status != http.StatusForbidden {
		t.Errorf("Status = %d, want %d", got.Status, http.StatusForbidden)
	}
	if got.Message != "Access denied: missing permission 'x'" {
		t.Errorf("Message = %q, want the original message", got.Message)
	}
}

func TestClassify_WrappedAccessDenied_StillDetected(t *testing.T) {
	err := fmt.Errorf("wrapping: %w", AccessDenied("nested"))

	got := Classify(err)

	if got.Code != CodeAccessDenied {
		t.Errorf("Code = %q, want %q", got.Code, CodeAccessDenied)
	}
	if got.Status != http.StatusForbidden {
		t.Errorf("Status = %d, want %d", got.Status, http.StatusForbidden)
	}
}

func TestClassify_APIErrors_KeepOwnStatus(t *testing.T) {
	tests := []struct {
		err        error
		wantCode   string
		wantStatus int
	}{
		{Generic("bad input"), CodeGeneric, http.StatusBadRequest},
		{NotFound("product not found"), CodeNotFound, http.StatusNotFound},
		{Conflict("sku already exists"), CodeConflict, http.StatusConflict},
		{fmt.Errorf("wrapped: %w", NotFound("nested")), CodeNotFound, http.StatusNotFound},
	}
	for _, tt := range tests {
		got := Classify(tt.err)
		if got.Code != tt.wantCode || got.Status != tt.wantStatus {
			t.Errorf("Classify(%v) = (%q, %d), want (%q, %d)", tt.err, got.Code, got.Status, tt.wantCode, tt.wantStatus)
		}
	}
}

func TestClassify_UnknownError_IsInternalWithoutLeakingMessage(t *testing.T) {
	tests := []error{
		errors.New("pq: password authentication failed for user \"products\""),
		fmt.Errorf("some wrapped: %w", errors.New("cause")),
		emptyMessageError{},
		nil,
	}
	for _, err := range tests {
		got := Classify(err)
		if got.Code != CodeInternal {
			t.Errorf("Classify(%v).Code = %q, want %q", err, got.Code, CodeInternal)
		}
		if got.Status != http.StatusInternalServerError {
			t.Errorf("Classify(%v).Status = %d, want %d", err, got.Status, http.StatusInternalServerError)
		}
		if got.Message != "Internal server error" {
			t.Errorf("Classify(%v).Message = %q, want the generic message", err, got.Message)
		}
	}
}

// emptyMessageError is an error whose Error() text is empty.
type emptyMessageError struct{}

func (emptyMessageError) Error() string { return "" }

func TestClassify_AccessDeniedWithEmptyMessage_StillFallsBackToUnknownError(t *testing.T) {
	err := &APIError{Code: CodeAccessDenied, Status: http.StatusForbidden}

	got := Classify(err)

	if got.Code != CodeAccessDenied {
		t.Errorf("Code = %q, want %q", got.Code, CodeAccessDenied)
	}
	if got.Message != "Unknown error" {
		t.Errorf("Message = %q, want %q", got.Message, "Unknown error")
	}
}
