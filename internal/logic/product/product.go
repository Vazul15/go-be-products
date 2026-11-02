package product

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"math"

	"github.com/google/uuid"

	productrepo "products/internal/connector/db/productRepository"
	productmodel "products/internal/model/product"
)

type ProductOrchestrator struct {
	repo productrepo.Repository
}

func (orch *ProductOrchestrator) ListProducts(ctx context.Context, page, limit int, filter productmodel.ProductFilter) (productmodel.ProductListResponse, error) {
	if page < 1 || limit < 1 || limit > 100 {
		return productmodel.ProductListResponse{}, fmt.Errorf("%w: invalid pagination", ErrInvalidProduct)
	}

	productPage, err := orch.repo.ListProducts(ctx, (page-1)*limit, limit, filter)
	if err != nil {
		return productmodel.ProductListResponse{}, err
	}

	items := make([]productmodel.ProductListItemDTO, 0, len(productPage.Items))
	for _, storedProduct := range productPage.Items {
		items = append(items, toProductListItemDTO(storedProduct))
	}

	totalPages := int(math.Ceil(float64(productPage.Total) / float64(limit)))
	return productmodel.ProductListResponse{
		Items:      items,
		Page:       page,
		Limit:      limit,
		Total:      productPage.Total,
		TotalPages: totalPages,
		HasNext:    page < totalPages,
	}, nil
}

func (orch *ProductOrchestrator) GetProduct(ctx context.Context, publicID string) (productmodel.ProductDetailResponse, error) {
	id, err := uuid.Parse(publicID)
	if err != nil {
		return productmodel.ProductDetailResponse{}, ErrNotFound
	}

	storedProduct, err := orch.repo.GetProduct(ctx, id)
	if errors.Is(err, sql.ErrNoRows) {
		return productmodel.ProductDetailResponse{}, ErrNotFound
	}
	if err != nil {
		return productmodel.ProductDetailResponse{}, err
	}

	return toProductDetailResponse(storedProduct), nil
}

func (orch *ProductOrchestrator) CreateProduct(ctx context.Context, req productmodel.CreateProductRequest) (productmodel.ProductDetailResponse, error) {
	draft, err := toNewProduct(req)
	if err != nil {
		return productmodel.ProductDetailResponse{}, err
	}

	storedProduct, err := orch.repo.CreateProduct(ctx, draft)
	if err != nil {
		return productmodel.ProductDetailResponse{}, err
	}

	return toProductDetailResponse(storedProduct), nil
}

func (orch *ProductOrchestrator) CreateProductVariant(ctx context.Context, productPublicID string, req productmodel.CreateProductVariantRequest) (productmodel.ProductVariantResponse, error) {
	productID, err := uuid.Parse(productPublicID)
	if err != nil {
		return productmodel.ProductVariantResponse{}, ErrNotFound
	}

	variant, err := toNewProductVariant(req)
	if err != nil {
		return productmodel.ProductVariantResponse{}, err
	}

	storedVariant, err := orch.repo.CreateProductVariant(ctx, productID, variant)
	if errors.Is(err, sql.ErrNoRows) {
		return productmodel.ProductVariantResponse{}, ErrNotFound
	}
	if err != nil {
		return productmodel.ProductVariantResponse{}, err
	}

	return toProductVariantResponse(storedVariant), nil
}

func (orch *ProductOrchestrator) ListProductVariants(ctx context.Context, productPublicID string) ([]productmodel.ProductVariantResponse, error) {
	productID, err := uuid.Parse(productPublicID)
	if err != nil {
		return nil, ErrNotFound
	}

	variants, err := orch.repo.ListProductVariants(ctx, productID)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, err
	}

	return toProductVariantResponses(variants), nil
}

func (orch *ProductOrchestrator) GetProductVariant(ctx context.Context, publicID string) (productmodel.ProductVariantResponse, error) {
	variantID, err := uuid.Parse(publicID)
	if err != nil {
		return productmodel.ProductVariantResponse{}, ErrNotFound
	}

	variant, err := orch.repo.GetProductVariant(ctx, variantID)
	if errors.Is(err, sql.ErrNoRows) {
		return productmodel.ProductVariantResponse{}, ErrNotFound
	}
	if err != nil {
		return productmodel.ProductVariantResponse{}, err
	}

	return toProductVariantResponse(variant), nil
}

func (orch *ProductOrchestrator) UpdateProductVariant(ctx context.Context, publicID string, req productmodel.UpdateProductVariantRequest) (productmodel.ProductVariantResponse, error) {
	variantID, err := uuid.Parse(publicID)
	if err != nil {
		return productmodel.ProductVariantResponse{}, ErrNotFound
	}

	current, err := orch.repo.GetProductVariant(ctx, variantID)
	if errors.Is(err, sql.ErrNoRows) {
		return productmodel.ProductVariantResponse{}, ErrNotFound
	}
	if err != nil {
		return productmodel.ProductVariantResponse{}, err
	}

	updated, err := toUpdatedProductVariant(current, req)
	if err != nil {
		return productmodel.ProductVariantResponse{}, err
	}

	storedVariant, err := orch.repo.UpdateProductVariant(ctx, variantID, updated)
	if errors.Is(err, sql.ErrNoRows) {
		return productmodel.ProductVariantResponse{}, ErrNotFound
	}
	if err != nil {
		return productmodel.ProductVariantResponse{}, err
	}

	return toProductVariantResponse(storedVariant), nil
}

func (orch *ProductOrchestrator) DeleteProductVariant(ctx context.Context, publicID string) error {
	variantID, err := uuid.Parse(publicID)
	if err != nil {
		return ErrNotFound
	}

	if err := orch.repo.DeleteProductVariant(ctx, variantID); errors.Is(err, sql.ErrNoRows) {
		return ErrNotFound
	} else {
		return err
	}
}

func (orch *ProductOrchestrator) UpdateProduct(ctx context.Context, publicID string, req productmodel.UpdateProductRequest) (productmodel.ProductDetailResponse, error) {
	id, err := uuid.Parse(publicID)
	if err != nil {
		return productmodel.ProductDetailResponse{}, ErrNotFound
	}

	current, err := orch.repo.GetProduct(ctx, id)
	if errors.Is(err, sql.ErrNoRows) {
		return productmodel.ProductDetailResponse{}, ErrNotFound
	}
	if err != nil {
		return productmodel.ProductDetailResponse{}, err
	}

	updated, err := toUpdatedProduct(current, req)
	if err != nil {
		return productmodel.ProductDetailResponse{}, err
	}

	storedProduct, err := orch.repo.UpdateProduct(ctx, id, updated)
	if errors.Is(err, sql.ErrNoRows) {
		return productmodel.ProductDetailResponse{}, ErrNotFound
	}
	if err != nil {
		return productmodel.ProductDetailResponse{}, err
	}

	return toProductDetailResponse(storedProduct), nil
}

func (orch *ProductOrchestrator) DeleteProduct(ctx context.Context, publicID string) error {
	id, err := uuid.Parse(publicID)
	if err != nil {
		return ErrNotFound
	}

	if err := orch.repo.DeleteProduct(ctx, id); errors.Is(err, sql.ErrNoRows) {
		return ErrNotFound
	} else {
		return err
	}
}
