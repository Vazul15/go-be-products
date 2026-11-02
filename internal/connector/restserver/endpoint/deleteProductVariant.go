package endpoint

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"

	"products/internal/logic/product"
)

type DeleteProductVariant struct {
	productHandler product.ProductHandler
}

func NewDeleteProductVariantEndpoint(handler product.ProductHandler) (*DeleteProductVariant, error) {
	if handler == nil {
		return nil, fmt.Errorf("endpoint: product handler must not be nil")
	}
	return &DeleteProductVariant{productHandler: handler}, nil
}

func MustNewDeleteProductVariantEndpoint(handler product.ProductHandler) *DeleteProductVariant {
	endpoint, err := NewDeleteProductVariantEndpoint(handler)
	if err != nil {
		panic(err)
	}
	return endpoint
}

func (endpoint *DeleteProductVariant) Endpoint(c *gin.Context) {
	if err := endpoint.productHandler.DeleteProductVariant(c.Request.Context(), c.Param("public_id")); err != nil {
		respondError(c, err)
		return
	}

	c.Status(http.StatusNoContent)
}
