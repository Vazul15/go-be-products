package apierror

import (
	"errors"
	"net/http"
	"testing"
)

func TestAccessDeniedIdentity(t *testing.T) {
	err := AccessDenied("Access denied: missing permission 'x'")

	if !errors.Is(err, ErrAccessDenied) {
		t.Fatal("expected errors.Is(err, ErrAccessDenied) to be true")
	}
	if errors.Is(err, ErrGeneric) {
		t.Fatal("expected errors.Is(err, ErrGeneric) to be false")
	}

	var apiErr *APIError
	if !errors.As(err, &apiErr) {
		t.Fatal("expected errors.As to succeed")
	}
	if apiErr.Code != CodeAccessDenied {
		t.Errorf("Code = %q, want %q", apiErr.Code, CodeAccessDenied)
	}
	if apiErr.Status != http.StatusForbidden {
		t.Errorf("Status = %d, want %d", apiErr.Status, http.StatusForbidden)
	}
}

func TestGenericIdentity(t *testing.T) {
	err := Generic("boom")

	if !errors.Is(err, ErrGeneric) {
		t.Fatal("expected errors.Is(err, ErrGeneric) to be true")
	}
	if errors.Is(err, ErrAccessDenied) {
		t.Fatal("expected errors.Is(err, ErrAccessDenied) to be false")
	}

	var apiErr *APIError
	if !errors.As(err, &apiErr) {
		t.Fatal("expected errors.As to succeed")
	}
	if apiErr.Code != CodeGeneric {
		t.Errorf("Code = %q, want %q", apiErr.Code, CodeGeneric)
	}
	if apiErr.Status != http.StatusBadRequest {
		t.Errorf("Status = %d, want %d", apiErr.Status, http.StatusBadRequest)
	}
}

func TestErrorMessage(t *testing.T) {
	err := AccessDenied("missing permission")
	want := "FAT_0114: missing permission"
	if err.Error() != want {
		t.Errorf("Error() = %q, want %q", err.Error(), want)
	}
}

func TestCodesMatchJava(t *testing.T) {
	if CodeAccessDenied != "FAT_0114" {
		t.Errorf("CodeAccessDenied = %q, want %q", CodeAccessDenied, "FAT_0114")
	}
	if CodeGeneric != "FAT_0202" {
		t.Errorf("CodeGeneric = %q, want %q", CodeGeneric, "FAT_0202")
	}
}

func TestDifferentInstancesSameCodeAreEqualIdentity(t *testing.T) {
	a := AccessDenied("first message")
	b := AccessDenied("second, different message")

	if !errors.Is(a, b) {
		t.Fatal("expected two AccessDenied errors with different messages to share identity via Is")
	}
}
