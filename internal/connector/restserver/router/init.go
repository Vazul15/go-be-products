package router

import (
	"fmt"

	"products/internal/connector/restserver/endpoint"
)

// New creates a Router.
func New(getProduct *endpoint.GetProduct, listProducts *endpoint.ListProducts, createProduct *endpoint.CreateProduct, updateProduct *endpoint.UpdateProduct, deleteProduct *endpoint.DeleteProduct, createProductVariant *endpoint.CreateProductVariant, listProductVariants *endpoint.ListProductVariants, getProductVariant *endpoint.GetProductVariant, updateProductVariant *endpoint.UpdateProductVariant, deleteProductVariant *endpoint.DeleteProductVariant) (*Router, error) {
	if getProduct == nil {
		return nil, fmt.Errorf("router: get product endpoint must not be nil")
	}
	if listProducts == nil {
		return nil, fmt.Errorf("router: list products endpoint must not be nil")
	}
	if createProduct == nil {
		return nil, fmt.Errorf("router: create product endpoint must not be nil")
	}
	if updateProduct == nil {
		return nil, fmt.Errorf("router: update product endpoint must not be nil")
	}
	if deleteProduct == nil {
		return nil, fmt.Errorf("router: delete product endpoint must not be nil")
	}
	if createProductVariant == nil {
		return nil, fmt.Errorf("router: create product variant endpoint must not be nil")
	}
	if listProductVariants == nil {
		return nil, fmt.Errorf("router: list product variants endpoint must not be nil")
	}
	if getProductVariant == nil {
		return nil, fmt.Errorf("router: get product variant endpoint must not be nil")
	}
	if updateProductVariant == nil {
		return nil, fmt.Errorf("router: update product variant endpoint must not be nil")
	}
	if deleteProductVariant == nil {
		return nil, fmt.Errorf("router: delete product variant endpoint must not be nil")
	}
	return &Router{
		getProduct:           getProduct,
		listProducts:         listProducts,
		createProduct:        createProduct,
		updateProduct:        updateProduct,
		deleteProduct:        deleteProduct,
		createProductVariant: createProductVariant,
		listProductVariants:  listProductVariants,
		getProductVariant:    getProductVariant,
		updateProductVariant: updateProductVariant,
		deleteProductVariant: deleteProductVariant,
	}, nil
}

// MustNew creates a Router and panics if initialization fails.
func MustNew(
	getProduct *endpoint.GetProduct,
	listProducts *endpoint.ListProducts,
	createProduct *endpoint.CreateProduct,
	updateProduct *endpoint.UpdateProduct,
	deleteProduct *endpoint.DeleteProduct,
	createProductVariant *endpoint.CreateProductVariant,
	listProductVariants *endpoint.ListProductVariants,
	getProductVariant *endpoint.GetProductVariant,
	updateProductVariant *endpoint.UpdateProductVariant,
	deleteProductVariant *endpoint.DeleteProductVariant,
) *Router {

	r, err := New(
		getProduct,
		listProducts,
		createProduct,
		updateProduct,
		deleteProduct,
		createProductVariant,
		listProductVariants,
		getProductVariant,
		updateProductVariant,
		deleteProductVariant,
	)
	if err != nil {
		panic(err)
	}
	return r
}
