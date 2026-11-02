package endpoint

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"

	"products/internal/logic/product"
	productmodel "products/internal/model/product"
)

type UpdateProductVariant struct {
	productHandler product.ProductHandler
}

func NewUpdateProductVariantEndpoint(handler product.ProductHandler) (*UpdateProductVariant, error) {
	if handler == nil {
		return nil, fmt.Errorf("endpoint: product handler must not be nil")
	}
	return &UpdateProductVariant{productHandler: handler}, nil
}

func MustNewUpdateProductVariantEndpoint(handler product.ProductHandler) *UpdateProductVariant {
	endpoint, err := NewUpdateProductVariantEndpoint(handler)
	if err != nil {
		panic(err)
	}
	return endpoint
}

func (endpoint *UpdateProductVariant) Endpoint(c *gin.Context) {
	var request productmodel.UpdateProductVariantRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		respondBadRequest(c, err)
		return
	}

	variant, err := endpoint.productHandler.UpdateProductVariant(c.Request.Context(), c.Param("public_id"), request)
	if err != nil {
		respondError(c, err)
		return
	}

	c.JSON(http.StatusOK, variant)
}
