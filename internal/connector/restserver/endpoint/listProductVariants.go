package endpoint

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"

	"products/internal/logic/product"
)

type ListProductVariants struct {
	productHandler product.ProductHandler
}

func NewListProductVariantsEndpoint(handler product.ProductHandler) (*ListProductVariants, error) {
	if handler == nil {
		return nil, fmt.Errorf("endpoint: product handler must not be nil")
	}
	return &ListProductVariants{productHandler: handler}, nil
}

func MustNewListProductVariantsEndpoint(handler product.ProductHandler) *ListProductVariants {
	endpoint, err := NewListProductVariantsEndpoint(handler)
	if err != nil {
		panic(err)
	}
	return endpoint
}

func (endpoint *ListProductVariants) Endpoint(c *gin.Context) {
	variants, err := endpoint.productHandler.ListProductVariants(c.Request.Context(), c.Param("public_id"))
	if err != nil {
		respondError(c, err)
		return
	}

	c.JSON(http.StatusOK, variants)
}
