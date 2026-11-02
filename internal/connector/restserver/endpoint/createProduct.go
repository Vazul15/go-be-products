package endpoint

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"

	"products/internal/logic/product"
	productmodel "products/internal/model/product"
)

type CreateProduct struct {
	productHandler product.ProductHandler
}

func NewCreateProductEndpoint(handler product.ProductHandler) (*CreateProduct, error) {
	if handler == nil {
		return nil, fmt.Errorf("endpoint: product handler must not be nil")
	}
	return &CreateProduct{productHandler: handler}, nil
}

func MustNewCreateProductEndpoint(handler product.ProductHandler) *CreateProduct {
	endpoint, err := NewCreateProductEndpoint(handler)
	if err != nil {
		panic(err)
	}
	return endpoint
}

func (endpoint *CreateProduct) Endpoint(c *gin.Context) {
	var request productmodel.CreateProductRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		respondBadRequest(c, err)
		return
	}

	created, err := endpoint.productHandler.CreateProduct(c.Request.Context(), request)
	if err != nil {
		respondError(c, err)
		return
	}

	c.JSON(http.StatusCreated, created)
}
