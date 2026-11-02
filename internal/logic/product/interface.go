package product

import (
	"context"

	productmodel "products/internal/model/product"
)

type ProductHandler interface {
	CreateProduct(ctx context.Context, req productmodel.CreateProductRequest) (productmodel.ProductDetailResponse, error)
	GetProduct(ctx context.Context, publicID string) (productmodel.ProductDetailResponse, error)
	ListProducts(ctx context.Context, page, limit int, filter productmodel.ProductFilter) (productmodel.ProductListResponse, error)
	UpdateProduct(ctx context.Context, publicID string, req productmodel.UpdateProductRequest) (productmodel.ProductDetailResponse, error)
	DeleteProduct(ctx context.Context, publicID string) error
	CreateProductVariant(ctx context.Context, productPublicID string, req productmodel.CreateProductVariantRequest) (productmodel.ProductVariantResponse, error)
	ListProductVariants(ctx context.Context, productPublicID string) ([]productmodel.ProductVariantResponse, error)
	GetProductVariant(ctx context.Context, publicID string) (productmodel.ProductVariantResponse, error)
	UpdateProductVariant(ctx context.Context, publicID string, req productmodel.UpdateProductVariantRequest) (productmodel.ProductVariantResponse, error)
	DeleteProductVariant(ctx context.Context, publicID string) error
}
