package endpoint

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"

	"products/internal/logic/product"
	"products/internal/model/apierror"
	productmodel "products/internal/model/product"
)

type productHandlerStub struct {
	product.ProductHandler
	err error
}

func (stub productHandlerStub) GetProduct(context.Context, string) (productmodel.ProductDetailResponse, error) {
	return productmodel.ProductDetailResponse{}, stub.err
}

func TestRespondError_MapsErrorsToStatus(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name       string
		err        error
		wantStatus int
		leakedText string
	}{
		{"not found", product.ErrNotFound, http.StatusNotFound, ""},
		{"invalid product", fmt.Errorf("%w: price must be greater than zero", product.ErrInvalidProduct), http.StatusBadRequest, ""},
		{"conflict", apierror.Conflict("a variant with this sku already exists"), http.StatusConflict, ""},
		{"unknown", errors.New("pq: connection refused to 10.0.0.5"), http.StatusInternalServerError, "10.0.0.5"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			endpoint := MustNewGetProductEndpoint(productHandlerStub{err: tt.err})
			engine := gin.New()
			engine.GET("/:public_id", endpoint.Endpoint)

			rec := httptest.NewRecorder()
			engine.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/abc", nil))

			if rec.Code != tt.wantStatus {
				t.Fatalf("status = %d, want %d (body %s)", rec.Code, tt.wantStatus, rec.Body.String())
			}
			if tt.leakedText != "" && strings.Contains(rec.Body.String(), tt.leakedText) {
				t.Fatalf("response leaked internal detail %q: %s", tt.leakedText, rec.Body.String())
			}
		})
	}
}

func TestCreateProduct_MalformedBody_ReturnsBadRequest(t *testing.T) {
	gin.SetMode(gin.TestMode)
	endpoint := MustNewCreateProductEndpoint(productHandlerStub{})
	engine := gin.New()
	engine.POST("/", endpoint.Endpoint)

	rec := httptest.NewRecorder()
	engine.ServeHTTP(rec, httptest.NewRequest(http.MethodPost, "/", strings.NewReader("{not json")))

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusBadRequest)
	}
}
