package product

import "time"

type CreateProductRequest struct {
	Name              string                        `json:"name" binding:"required,min=3,max=255"`
	Slug              string                        `json:"slug" binding:"required,min=3,max=255"`
	Brand             string                        `json:"brand" binding:"required,min=1,max=255"`
	Category          string                        `json:"category" binding:"required,oneof=tops bottoms dresses outerwear footwear accessories underwear other"`
	Gender            string                        `json:"gender" binding:"required,oneof=women men unisex kids"`
	Material          string                        `json:"material,omitempty" binding:"omitempty,max=255"`
	ShortDescription  string                        `json:"short_description,omitempty"`
	Description       string                        `json:"description" binding:"required,min=10"`
	Tags              []string                      `json:"tags,omitempty"`
	Price             int64                         `json:"price" binding:"required,gt=0"`
	DiscountPrice     *int64                        `json:"discount_price,omitempty" binding:"omitempty,gt=0,ltfield=Price"`
	DiscountExpiresAt *time.Time                    `json:"discount_expires_at,omitempty"`
	Currency          string                        `json:"currency" binding:"required"`
	MainImageURL      string                        `json:"main_image_url" binding:"required,url"`
	Gallery           []string                      `json:"gallery,omitempty" binding:"omitempty,dive,url"`
	Variants          []CreateProductVariantRequest `json:"variants" binding:"required,min=1,dive"`
	Type              string                        `json:"type" binding:"required,oneof=t-shirt shirt blouse trousers jeans dress skirt jacket coat hoodie sweater underwear shoes accessory other"`
	IsFeatured        bool                          `json:"is_featured"`
}

type CreateProductVariantRequest struct {
	SKU       string `json:"sku" binding:"required,min=1,max=100"`
	Size      string `json:"size" binding:"required,max=30"`
	Color     string `json:"color" binding:"required,max=100"`
	ColorCode string `json:"color_code,omitempty" binding:"omitempty,max=7"`
	Price     *int64 `json:"price,omitempty" binding:"omitempty,gt=0"`
	Barcode   string `json:"barcode,omitempty" binding:"omitempty,max=100"`
}

type UpdateProductVariantRequest struct {
	SKU       *string `json:"sku,omitempty" binding:"omitempty,min=1,max=100"`
	Size      *string `json:"size,omitempty" binding:"omitempty,max=30"`
	Color     *string `json:"color,omitempty" binding:"omitempty,max=100"`
	ColorCode *string `json:"color_code,omitempty" binding:"omitempty,max=7"`
	Price     *int64  `json:"price,omitempty" binding:"omitempty,gt=0"`
	Barcode   *string `json:"barcode,omitempty" binding:"omitempty,max=100"`
}

type UpdateProductRequest struct {
	Name              *string    `json:"name,omitempty" binding:"omitempty,min=3,max=255"`
	Slug              *string    `json:"slug,omitempty" binding:"omitempty,min=3,max=255"`
	Brand             *string    `json:"brand,omitempty" binding:"omitempty,min=1,max=255"`
	Category          *string    `json:"category,omitempty" binding:"omitempty,oneof=tops bottoms dresses outerwear footwear accessories underwear other"`
	Gender            *string    `json:"gender,omitempty" binding:"omitempty,oneof=women men unisex kids"`
	Material          *string    `json:"material,omitempty" binding:"omitempty,max=255"`
	ShortDescription  *string    `json:"short_description,omitempty"`
	Description       *string    `json:"description,omitempty" binding:"omitempty,min=10"`
	Tags              []string   `json:"tags,omitempty"`
	Price             *int64     `json:"price,omitempty" binding:"omitempty,gt=0"`
	DiscountPrice     *int64     `json:"discount_price,omitempty" binding:"omitempty,gt=0,ltfield=Price"`
	DiscountExpiresAt *time.Time `json:"discount_expires_at,omitempty"`
	Currency          *string    `json:"currency,omitempty"`
	MainImageURL      *string    `json:"main_image_url,omitempty" binding:"omitempty,url"`
	Gallery           []string   `json:"gallery,omitempty" binding:"omitempty,dive,url"`
	Type              *string    `json:"type,omitempty" binding:"omitempty,oneof=t-shirt shirt blouse trousers jeans dress skirt jacket coat hoodie sweater underwear shoes accessory other"`
	IsFeatured        *bool      `json:"is_featured,omitempty"`
}
