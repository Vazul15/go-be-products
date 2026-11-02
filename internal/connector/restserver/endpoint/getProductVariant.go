package endpoint

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"

	"products/internal/logic/product"
)

type GetProductVariant struct {
	productHandler product.ProductHandler
}

func NewGetProductVariantEndpoint(handler product.ProductHandler) (*GetProductVariant, error) {
	if handler == nil {
		return nil, fmt.Errorf("endpoint: product handler must not be nil")
	}
	return &GetProductVariant{productHandler: handler}, nil
}

func MustNewGetProductVariantEndpoint(handler product.ProductHandler) *GetProductVariant {
	endpoint, err := NewGetProductVariantEndpoint(handler)
	if err != nil {
		panic(err)
	}
	return endpoint
}

func (endpoint *GetProductVariant) Endpoint(c *gin.Context) {
	variant, err := endpoint.productHandler.GetProductVariant(c.Request.Context(), c.Param("public_id"))
	if err != nil {
		respondError(c, err)
		return
	}

	c.JSON(http.StatusOK, variant)
}
