package productRepository

import (
	"context"

	"github.com/google/uuid"

	"products/internal/model/product"
)

// Repository is the persistence contract this connector publishes for
// internal/logic/product to depend on directly.
type Repository interface {
	ListProducts(ctx context.Context, offset, limit int, filter product.ProductFilter) (product.ProductPage, error)
	GetProduct(ctx context.Context, publicID uuid.UUID) (product.Product, error)
	CreateProduct(ctx context.Context, draft product.NewProduct) (product.Product, error)
	CreateProductVariant(ctx context.Context, productPublicID uuid.UUID, variant product.NewProductVariant) (product.ProductVariant, error)
	ListProductVariants(ctx context.Context, productPublicID uuid.UUID) ([]product.ProductVariant, error)
	GetProductVariant(ctx context.Context, publicID uuid.UUID) (product.ProductVariant, error)
	UpdateProductVariant(ctx context.Context, publicID uuid.UUID, variant product.ProductVariant) (product.ProductVariant, error)
	DeleteProductVariant(ctx context.Context, publicID uuid.UUID) error
	UpdateProduct(ctx context.Context, publicID uuid.UUID, storedProduct product.Product) (product.Product, error)
	DeleteProduct(ctx context.Context, publicID uuid.UUID) error
}
