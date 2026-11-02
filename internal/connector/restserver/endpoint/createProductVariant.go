package endpoint

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"

	"products/internal/logic/product"
	productmodel "products/internal/model/product"
)

type CreateProductVariant struct {
	productHandler product.ProductHandler
}

func NewCreateProductVariantEndpoint(handler product.ProductHandler) (*CreateProductVariant, error) {
	if handler == nil {
		return nil, fmt.Errorf("endpoint: product handler must not be nil")
	}
	return &CreateProductVariant{productHandler: handler}, nil
}

func MustNewCreateProductVariantEndpoint(handler product.ProductHandler) *CreateProductVariant {
	endpoint, err := NewCreateProductVariantEndpoint(handler)
	if err != nil {
		panic(err)
	}
	return endpoint
}

func (endpoint *CreateProductVariant) Endpoint(c *gin.Context) {
	var request productmodel.CreateProductVariantRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		respondBadRequest(c, err)
		return
	}

	variant, err := endpoint.productHandler.CreateProductVariant(c.Request.Context(), c.Param("public_id"), request)
	if err != nil {
		respondError(c, err)
		return
	}

	c.JSON(http.StatusCreated, variant)
}
