package productRepository

import (
	"errors"

	"github.com/lib/pq"

	"products/internal/model/apierror"
)

const (
	pqUniqueViolation = "23505"
	pqCheckViolation  = "23514"
)

var uniqueConstraintMessages = map[string]string{
	"products_slug_key":                 "a product with this slug already exists",
	"product_variants_sku_key":          "a variant with this sku already exists",
	"product_variants_barcode_key":      "a variant with this barcode already exists",
	"product_variant_size_color_unique": "this product already has a variant with the same size and color",
}

// translateError turns Postgres constraint violations into client-facing API errors.
func translateError(err error) error {
	var pqErr *pq.Error
	if !errors.As(err, &pqErr) {
		return err
	}

	switch pqErr.Code {
	case pqUniqueViolation:
		if message, ok := uniqueConstraintMessages[pqErr.Constraint]; ok {
			return apierror.Conflict(message)
		}
		return apierror.Conflict("resource already exists")
	case pqCheckViolation:
		return apierror.Generic("value violates a data constraint")
	default:
		return err
	}
}
