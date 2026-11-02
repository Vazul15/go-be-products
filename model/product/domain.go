package product

import (
	"time"

	"github.com/google/uuid"
)

type Product struct {
	ID               int64
	PublicID         uuid.UUID
	Name             string
	Slug             string
	Brand            string
	Category         string
	Gender           string
	Material         string
	ShortDescription string
	Description      string
	Tags             []string
	Price            int64
	DiscountPrice    *int64
	DiscountExpiresAt *time.Time
	Currency         string
	MainImageURL     string
	Gallery          []string
	Type             string
	IsFeatured       bool
	IsActive         bool
	Variants         []ProductVariant
}

type NewProduct struct {
	PublicID         uuid.UUID
	Name             string
	Slug             string
	Brand            string
	Category         string
	Gender           string
	Material         string
	ShortDescription string
	Description      string
	Tags             []string
	Price            int64
	DiscountPrice    *int64
	DiscountExpiresAt *time.Time
	Currency         string
	MainImageURL     string
	Gallery          []string
	Type             string
	IsFeatured       bool
	IsActive         bool
	Variants         []NewProductVariant
}

type ProductVariant struct {
	ID        int64
	PublicID  uuid.UUID
	SKU       string
	Size      string
	Color     string
	ColorCode string
	Price     *int64
	Barcode   string
}

type NewProductVariant struct {
	PublicID  uuid.UUID
	SKU       string
	Size      string
	Color     string
	ColorCode string
	Price     *int64
	Barcode   string
}

type ProductPage struct {
	Items []Product
	Total int64
}

type ProductFilter struct {
	Search   string
	Category string
	Gender   string
	Type     string
}
