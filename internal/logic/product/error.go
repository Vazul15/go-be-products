package product

import "errors"

var (
	ErrNotFound       = errors.New("product not found")
	ErrInvalidProduct = errors.New("invalid product")
)
