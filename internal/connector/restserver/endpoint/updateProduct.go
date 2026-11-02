package endpoint

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"

	"products/internal/logic/product"
	productmodel "products/internal/model/product"
)

type UpdateProduct struct {
	productHandler product.ProductHandler
}

func NewUpdateProductEndpoint(handler product.ProductHandler) (*UpdateProduct, error) {
	if handler == nil {
		return nil, fmt.Errorf("endpoint: product handler must not be nil")
	}
	return &UpdateProduct{productHandler: handler}, nil
}

func MustNewUpdateProductEndpoint(handler product.ProductHandler) *UpdateProduct {
	endpoint, err := NewUpdateProductEndpoint(handler)
	if err != nil {
		panic(err)
	}
	return endpoint
}

func (endpoint *UpdateProduct) Endpoint(c *gin.Context) {
	var request productmodel.UpdateProductRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		respondBadRequest(c, err)
		return
	}

	updated, err := endpoint.productHandler.UpdateProduct(c.Request.Context(), c.Param("public_id"), request)
	if err != nil {
		respondError(c, err)
		return
	}

	c.JSON(http.StatusOK, updated)
}
