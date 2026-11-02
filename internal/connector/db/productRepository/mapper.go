package productRepository

import (
	"database/sql"
	"time"

	"github.com/google/uuid"

	"products/internal/connector/db/sqlc"
	"products/internal/model/product"
)

func toProduct(p sqlc.Product) product.Product {
	return product.Product{
		ID:                p.ID,
		PublicID:          p.PublicID,
		Name:              p.Name,
		Slug:              p.Slug,
		Brand:             p.Brand,
		Category:          p.Category,
		Gender:            p.Gender,
		Material:          NullStringValue(p.Material),
		ShortDescription:  NullStringValue(p.ShortDescription),
		Description:       NullStringValue(p.Description),
		Tags:              p.Tags,
		Price:             p.Price,
		DiscountPrice:     NullInt64Ptr(p.DiscountPrice),
		DiscountExpiresAt: NullTimePtr(p.DiscountExpiresAt),
		Currency:          p.Currency,
		MainImageURL:      p.MainImageUrl,
		Gallery:           p.Gallery,
		Type:              p.Type,
		IsFeatured:        p.IsFeatured,
		IsActive:          p.IsActive,
	}
}

func toProductVariant(p sqlc.ProductVariant) product.ProductVariant {
	return product.ProductVariant{
		ID:        p.ID,
		PublicID:  p.PublicID,
		SKU:       p.Sku,
		Size:      p.Size,
		Color:     p.Color,
		ColorCode: NullStringValue(p.ColorCode),
		Price:     NullInt64Ptr(p.Price),
		Barcode:   NullStringValue(p.Barcode),
	}
}

func NullInt64Ptr(value sql.NullInt64) *int64 {
	if !value.Valid {
		return nil
	}
	result := value.Int64
	return &result
}

func NullTimePtr(value sql.NullTime) *time.Time {
	if !value.Valid {
		return nil
	}
	result := value.Time
	return &result
}

func NullStringValue(value sql.NullString) string {
	if !value.Valid {
		return ""
	}
	return value.String
}

func toCreateProductParams(draft product.NewProduct) sqlc.CreateProductParams {
	return sqlc.CreateProductParams{
		PublicID:          draft.PublicID,
		Name:              draft.Name,
		Slug:              draft.Slug,
		Brand:             draft.Brand,
		Category:          draft.Category,
		Gender:            draft.Gender,
		Material:          sql.NullString{String: draft.Material, Valid: draft.Material != ""},
		ShortDescription:  sql.NullString{String: draft.ShortDescription, Valid: draft.ShortDescription != ""},
		Description:       sql.NullString{String: draft.Description, Valid: draft.Description != ""},
		Tags:              draft.Tags,
		Price:             draft.Price,
		DiscountPrice:     sql.NullInt64{Int64: nullableInt64(draft.DiscountPrice), Valid: draft.DiscountPrice != nil},
		DiscountExpiresAt: sql.NullTime{Time: nullableTime(draft.DiscountExpiresAt), Valid: draft.DiscountExpiresAt != nil},
		Currency:          draft.Currency,
		MainImageUrl:      draft.MainImageURL,
		Gallery:           draft.Gallery,
		Type:              draft.Type,
		IsFeatured:        draft.IsFeatured,
		IsActive:          draft.IsActive,
	}
}

func toUpdateProductParams(publicID uuid.UUID, storedProduct product.Product) sqlc.UpdateProductParams {
	return sqlc.UpdateProductParams{
		PublicID:          publicID,
		Name:              storedProduct.Name,
		Slug:              storedProduct.Slug,
		Brand:             storedProduct.Brand,
		Category:          storedProduct.Category,
		Gender:            storedProduct.Gender,
		Material:          sql.NullString{String: storedProduct.Material, Valid: storedProduct.Material != ""},
		ShortDescription:  sql.NullString{String: storedProduct.ShortDescription, Valid: storedProduct.ShortDescription != ""},
		Description:       sql.NullString{String: storedProduct.Description, Valid: storedProduct.Description != ""},
		Tags:              storedProduct.Tags,
		Price:             storedProduct.Price,
		DiscountPrice:     sql.NullInt64{Int64: nullableInt64(storedProduct.DiscountPrice), Valid: storedProduct.DiscountPrice != nil},
		DiscountExpiresAt: sql.NullTime{Time: nullableTime(storedProduct.DiscountExpiresAt), Valid: storedProduct.DiscountExpiresAt != nil},
		Currency:          storedProduct.Currency,
		MainImageUrl:      storedProduct.MainImageURL,
		Gallery:           storedProduct.Gallery,
		Type:              storedProduct.Type,
		IsFeatured:        storedProduct.IsFeatured,
		IsActive:          storedProduct.IsActive,
	}
}

func toCreateProductVariantParams(productID int64, variant product.NewProductVariant) sqlc.CreateProductVariantParams {
	return sqlc.CreateProductVariantParams{
		PublicID:  variant.PublicID,
		ProductID: productID,
		Sku:       variant.SKU,
		Size:      variant.Size,
		Color:     variant.Color,
		ColorCode: sql.NullString{String: variant.ColorCode, Valid: variant.ColorCode != ""},
		Price:     sql.NullInt64{Int64: nullableInt64(variant.Price), Valid: variant.Price != nil},
		Barcode:   sql.NullString{String: variant.Barcode, Valid: variant.Barcode != ""},
	}
}

func toUpdateProductVariantParams(publicID uuid.UUID, variant product.ProductVariant) sqlc.UpdateProductVariantParams {
	return sqlc.UpdateProductVariantParams{
		PublicID:  publicID,
		Sku:       variant.SKU,
		Size:      variant.Size,
		Color:     variant.Color,
		ColorCode: sql.NullString{String: variant.ColorCode, Valid: variant.ColorCode != ""},
		Price:     sql.NullInt64{Int64: nullableInt64(variant.Price), Valid: variant.Price != nil},
		Barcode:   sql.NullString{String: variant.Barcode, Valid: variant.Barcode != ""},
	}
}

func nullableInt64(v *int64) int64 {
	if v == nil {
		return 0
	}
	return *v
}

func nullableTime(v *time.Time) time.Time {
	if v == nil {
		return time.Time{}
	}
	return *v
}
