package endpoint

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"

	"products/internal/logic/product"
)

type DeleteProduct struct {
	productHandler product.ProductHandler
}

func NewDeleteProductEndpoint(handler product.ProductHandler) (*DeleteProduct, error) {
	if handler == nil {
		return nil, fmt.Errorf("endpoint: product handler must not be nil")
	}
	return &DeleteProduct{productHandler: handler}, nil
}

func MustNewDeleteProductEndpoint(handler product.ProductHandler) *DeleteProduct {
	endpoint, err := NewDeleteProductEndpoint(handler)
	if err != nil {
		panic(err)
	}
	return endpoint
}

func (endpoint *DeleteProduct) Endpoint(c *gin.Context) {
	if err := endpoint.productHandler.DeleteProduct(c.Request.Context(), c.Param("public_id")); err != nil {
		respondError(c, err)
		return
	}

	c.Status(http.StatusNoContent)
}
