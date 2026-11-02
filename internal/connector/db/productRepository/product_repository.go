package productRepository

import (
	"context"
	"database/sql"

	"github.com/google/uuid"

	"products/internal/connector/db/sqlc"
	"products/internal/model/product"
)

type ProductRepository struct {
	db      *sql.DB
	queries *sqlc.Queries
}

func NewProductRepository(db *sql.DB, queries *sqlc.Queries) *ProductRepository {
	return &ProductRepository{db: db, queries: queries}
}

// Compile-time check that the DB connector implements the published Repository contract.
var _ Repository = (*ProductRepository)(nil)

func (repo *ProductRepository) ListProducts(ctx context.Context, offset, limit int, filter product.ProductFilter) (product.ProductPage, error) {
	filterParams := sqlc.CountProductsParams{
		Search:   filter.Search,
		Category: filter.Category,
		Gender:   filter.Gender,
		Type:     filter.Type,
	}
	total, err := repo.queries.CountProducts(ctx, filterParams)
	if err != nil {
		return product.ProductPage{}, err
	}

	products, err := repo.queries.ListProducts(ctx, sqlc.ListProductsParams{
		Search:     filter.Search,
		Category:   filter.Category,
		Gender:     filter.Gender,
		Type:       filter.Type,
		PageOffset: int32(offset),
		PageLimit:  int32(limit),
	})
	if err != nil {
		return product.ProductPage{}, err
	}

	items := make([]product.Product, 0, len(products))
	for _, storedProduct := range products {
		items = append(items, toProduct(storedProduct))
	}

	return product.ProductPage{
		Items: items,
		Total: total,
	}, nil
}

func (repo *ProductRepository) GetProduct(ctx context.Context, publicID uuid.UUID) (product.Product, error) {
	storedProduct, err := repo.queries.GetProductByPublicID(ctx, publicID)
	if err != nil {
		return product.Product{}, err
	}

	variants, err := repo.queries.ListProductVariantsByProductPublicID(ctx, publicID)
	if err != nil {
		return product.Product{}, err
	}

	result := toProduct(storedProduct)
	result.Variants = make([]product.ProductVariant, 0, len(variants))
	for _, variant := range variants {
		result.Variants = append(result.Variants, toProductVariant(variant))
	}

	return result, nil
}

func (repo *ProductRepository) CreateProductVariant(ctx context.Context, productPublicID uuid.UUID, variant product.NewProductVariant) (product.ProductVariant, error) {
	tx, err := repo.db.BeginTx(ctx, nil)
	if err != nil {
		return product.ProductVariant{}, err
	}

	queries := repo.queries.WithTx(tx)
	productID, err := queries.GetProductIDByPublicID(ctx, productPublicID)
	if err != nil {
		_ = tx.Rollback()
		return product.ProductVariant{}, err
	}

	storedVariant, err := queries.CreateProductVariant(ctx, toCreateProductVariantParams(productID, variant))
	if err != nil {
		_ = tx.Rollback()
		return product.ProductVariant{}, translateError(err)
	}

	if err := tx.Commit(); err != nil {
		return product.ProductVariant{}, err
	}

	return toProductVariant(storedVariant), nil
}

func (repo *ProductRepository) ListProductVariants(ctx context.Context, productPublicID uuid.UUID) ([]product.ProductVariant, error) {
	if _, err := repo.queries.GetProductIDByPublicID(ctx, productPublicID); err != nil {
		return nil, err
	}

	storedVariants, err := repo.queries.ListProductVariantsByProductPublicID(ctx, productPublicID)
	if err != nil {
		return nil, err
	}

	variants := make([]product.ProductVariant, 0, len(storedVariants))
	for _, storedVariant := range storedVariants {
		variants = append(variants, toProductVariant(storedVariant))
	}
	return variants, nil
}

func (repo *ProductRepository) GetProductVariant(ctx context.Context, publicID uuid.UUID) (product.ProductVariant, error) {
	storedVariant, err := repo.queries.GetProductVariantByPublicID(ctx, publicID)
	if err != nil {
		return product.ProductVariant{}, err
	}
	return toProductVariant(storedVariant), nil
}

func (repo *ProductRepository) UpdateProductVariant(ctx context.Context, publicID uuid.UUID, variant product.ProductVariant) (product.ProductVariant, error) {
	storedVariant, err := repo.queries.UpdateProductVariant(ctx, toUpdateProductVariantParams(publicID, variant))
	if err != nil {
		return product.ProductVariant{}, translateError(err)
	}
	return toProductVariant(storedVariant), nil
}

func (repo *ProductRepository) DeleteProductVariant(ctx context.Context, publicID uuid.UUID) error {
	_, err := repo.queries.SoftDeleteProductVariant(ctx, publicID)
	return err
}

func (repo *ProductRepository) CreateProduct(ctx context.Context, draft product.NewProduct) (product.Product, error) {
	tx, err := repo.db.BeginTx(ctx, nil)
	if err != nil {
		return product.Product{}, err
	}

	queries := repo.queries.WithTx(tx)
	storedProduct, err := queries.CreateProduct(ctx, toCreateProductParams(draft))
	if err != nil {
		_ = tx.Rollback()
		return product.Product{}, translateError(err)
	}

	storedVariants := make([]product.ProductVariant, 0, len(draft.Variants))
	for _, variant := range draft.Variants {
		storedVariant, err := queries.CreateProductVariant(ctx, toCreateProductVariantParams(storedProduct.ID, variant))
		if err != nil {
			_ = tx.Rollback()
			return product.Product{}, translateError(err)
		}
		storedVariants = append(storedVariants, toProductVariant(storedVariant))
	}

	if err := tx.Commit(); err != nil {
		return product.Product{}, err
	}

	result := toProduct(storedProduct)
	result.Variants = storedVariants
	return result, nil
}

func (repo *ProductRepository) UpdateProduct(ctx context.Context, publicID uuid.UUID, storedProduct product.Product) (product.Product, error) {
	updatedProduct, err := repo.queries.UpdateProduct(ctx, toUpdateProductParams(publicID, storedProduct))
	if err != nil {
		return product.Product{}, translateError(err)
	}

	variants, err := repo.queries.ListProductVariantsByProductPublicID(ctx, publicID)
	if err != nil {
		return product.Product{}, err
	}

	result := toProduct(updatedProduct)
	result.Variants = make([]product.ProductVariant, 0, len(variants))
	for _, variant := range variants {
		result.Variants = append(result.Variants, toProductVariant(variant))
	}
	return result, nil
}

func (repo *ProductRepository) DeleteProduct(ctx context.Context, publicID uuid.UUID) error {
	_, err := repo.queries.SoftDeleteProduct(ctx, publicID)
	return err
}
