package endpoint

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"products/internal/logic/product"
	productmodel "products/internal/model/product"

	"github.com/gin-gonic/gin"
)

type ListProducts struct {
	productHandler product.ProductHandler
}

func NewListProductsEndpoint(handler product.ProductHandler) (*ListProducts, error) {
	if handler == nil {
		return nil, fmt.Errorf("endpoint: product handler must not be nil")
	}
	return &ListProducts{productHandler: handler}, nil
}

func MustNewListProductsEndpoint(handler product.ProductHandler) *ListProducts {
	endpoint, err := NewListProductsEndpoint(handler)
	if err != nil {
		panic(err)
	}
	return endpoint
}

const (
	defaultLimit = 20
	maxPageSize  = 100
)

func (handler *ListProducts) Endpoint(c *gin.Context) {
	page, err := strconv.Atoi(c.DefaultQuery("page", "1"))
	if err != nil || page < 1 {
		respondBadRequest(c, errors.New("invalid page parameter"))
		return
	}

	limit, err := strconv.Atoi(c.DefaultQuery("limit", strconv.Itoa(defaultLimit)))
	if err != nil || limit < 1 || limit > maxPageSize {
		respondBadRequest(c, errors.New("invalid limit parameter"))
		return
	}

	filter := productmodel.ProductFilter{
		Search:   c.Query("q"),
		Category: c.Query("category"),
		Gender:   c.Query("gender"),
		Type:     c.Query("type"),
	}
	response, err := handler.productHandler.ListProducts(c.Request.Context(), page, limit, filter)
	if err != nil {
		respondError(c, err)
		return
	}

	c.JSON(http.StatusOK, response)
}
