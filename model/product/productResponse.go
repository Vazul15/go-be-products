package product

import "time"

type ProductListResponse struct {
	Items      []ProductListItemDTO `json:"items"`
	Page       int                  `json:"page"`
	Limit      int                  `json:"limit"`
	Total      int64                `json:"total"`
	TotalPages int                  `json:"total_pages"`
	HasNext    bool                 `json:"has_next"`
}

type ProductDetailResponse struct {
	PublicID         string                   `json:"public_id"`
	Name             string                   `json:"name"`
	Slug             string                   `json:"slug"`
	Brand            string                   `json:"brand"`
	Category         string                   `json:"category"`
	Gender           string                   `json:"gender"`
	Material         string                   `json:"material,omitempty"`
	ShortDescription string                   `json:"short_description,omitempty"`
	Description      string                   `json:"description"`
	Tags             []string                 `json:"tags,omitempty"`
	Price            int64                    `json:"price"`
	DiscountPrice    *int64                   `json:"discount_price,omitempty"`
	DiscountPercent  *int32                   `json:"discount_percent,omitempty"`
	DiscountExpiresAt *time.Time               `json:"discount_expires_at,omitempty"`
	Currency         string                   `json:"currency"`
	MainImageURL     string                   `json:"main_image_url"`
	Gallery          []string                 `json:"gallery,omitempty"`
	Type             string                   `json:"type"`
	IsFeatured       bool                     `json:"is_featured"`
	Variants         []ProductVariantResponse `json:"variants"`
}

type ProductVariantResponse struct {
	PublicID  string `json:"public_id"`
	SKU       string `json:"sku"`
	Size      string `json:"size"`
	Color     string `json:"color"`
	ColorCode string `json:"color_code,omitempty"`
	Price     *int64 `json:"price,omitempty"`
	Barcode   string `json:"barcode,omitempty"`
}

type ProductListItemDTO struct {
	PublicID      string `json:"public_id"`
	Name          string `json:"name"`
	Slug          string `json:"slug"`
	Brand         string `json:"brand"`
	Category      string `json:"category"`
	Gender        string `json:"gender"`
	Price         int64  `json:"price"`
	DiscountPrice *int64 `json:"discount_price,omitempty"`
	Currency      string `json:"currency"`
	MainImageURL  string `json:"main_image_url"`
}
