package productRepository

import (
	"database/sql"
	"errors"
	"net/http"
	"testing"

	"github.com/lib/pq"

	"products/internal/model/apierror"
)

func TestTranslateError(t *testing.T) {
	tests := []struct {
		name       string
		err        error
		wantStatus int
		wantMsg    string
	}{
		{"duplicate sku", &pq.Error{Code: pqUniqueViolation, Constraint: "product_variants_sku_key"}, http.StatusConflict, "a variant with this sku already exists"},
		{"duplicate size/color", &pq.Error{Code: pqUniqueViolation, Constraint: "product_variant_size_color_unique"}, http.StatusConflict, "this product already has a variant with the same size and color"},
		{"unknown unique", &pq.Error{Code: pqUniqueViolation, Constraint: "other_key"}, http.StatusConflict, "resource already exists"},
		{"check violation", &pq.Error{Code: pqCheckViolation}, http.StatusBadRequest, "value violates a data constraint"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var apiErr *apierror.APIError
			if !errors.As(translateError(tt.err), &apiErr) {
				t.Fatalf("translateError(%v) is not an APIError", tt.err)
			}
			if apiErr.Status != tt.wantStatus || apiErr.Message != tt.wantMsg {
				t.Fatalf("got (%d, %q), want (%d, %q)", apiErr.Status, apiErr.Message, tt.wantStatus, tt.wantMsg)
			}
		})
	}
}

func TestTranslateError_PassesThroughOtherErrors(t *testing.T) {
	for _, err := range []error{sql.ErrNoRows, &pq.Error{Code: "08006"}, errors.New("boom")} {
		if got := translateError(err); got != err {
			t.Errorf("translateError(%v) = %v, want unchanged", err, got)
		}
	}
}
