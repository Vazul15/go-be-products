package endpoint

import (
	"fmt"
	"net/http"

	"products/internal/logic/product"

	"github.com/gin-gonic/gin"
)

type GetProduct struct {
	productHandler product.ProductHandler
}

func NewGetProductEndpoint(handler product.ProductHandler) (*GetProduct, error) {
	if handler == nil {
		return nil, fmt.Errorf("endpoint: product handler must not be nil")
	}
	return &GetProduct{productHandler: handler}, nil
}

func MustNewGetProductEndpoint(handler product.ProductHandler) *GetProduct {
	e, err := NewGetProductEndpoint(handler)
	if err != nil {
		panic(err)
	}
	return e
}

func (handler *GetProduct) Endpoint(c *gin.Context) {
	publicID := c.Param("public_id")

	product, err := handler.productHandler.GetProduct(c.Request.Context(), publicID)
	if err != nil {
		respondError(c, err)
		return
	}

	c.JSON(http.StatusOK, product)
}
