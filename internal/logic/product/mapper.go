package product

import (
	"fmt"
	"strings"

	"github.com/google/uuid"

	productmodel "products/internal/model/product"
)

func toNewProduct(request productmodel.CreateProductRequest) (productmodel.NewProduct, error) {
	name := strings.TrimSpace(request.Name)
	slug := strings.ToLower(strings.TrimSpace(request.Slug))
	brand := strings.TrimSpace(request.Brand)
	category := strings.ToLower(strings.TrimSpace(request.Category))
	gender := strings.ToLower(strings.TrimSpace(request.Gender))
	currency := strings.ToUpper(strings.TrimSpace(request.Currency))
	typeName := strings.ToLower(strings.TrimSpace(request.Type))

	switch {
	case name == "":
		return productmodel.NewProduct{}, fmt.Errorf("%w: name is required", ErrInvalidProduct)
	case slug == "":
		return productmodel.NewProduct{}, fmt.Errorf("%w: slug is required", ErrInvalidProduct)
	case brand == "":
		return productmodel.NewProduct{}, fmt.Errorf("%w: brand is required", ErrInvalidProduct)
	case category == "":
		return productmodel.NewProduct{}, fmt.Errorf("%w: category is required", ErrInvalidProduct)
	case gender == "":
		return productmodel.NewProduct{}, fmt.Errorf("%w: gender is required", ErrInvalidProduct)
	case request.Price <= 0:
		return productmodel.NewProduct{}, fmt.Errorf("%w: price must be greater than zero", ErrInvalidProduct)
	case request.DiscountPrice != nil && (*request.DiscountPrice <= 0 || *request.DiscountPrice >= request.Price):
		return productmodel.NewProduct{}, fmt.Errorf("%w: discount price must be lower than price", ErrInvalidProduct)
	case request.DiscountExpiresAt != nil && request.DiscountPrice == nil:
		return productmodel.NewProduct{}, fmt.Errorf("%w: discount expiry requires a discount price", ErrInvalidProduct)
	case len(request.Variants) == 0:
		return productmodel.NewProduct{}, fmt.Errorf("%w: at least one variant is required", ErrInvalidProduct)
	case currency == "":
		return productmodel.NewProduct{}, fmt.Errorf("%w: currency is required", ErrInvalidProduct)
	}

	variants, err := toNewProductVariants(request.Variants)
	if err != nil {
		return productmodel.NewProduct{}, err
	}

	return productmodel.NewProduct{
		PublicID:         uuid.New(),
		Name:             name,
		Slug:             slug,
		Brand:            brand,
		Category:         category,
		Gender:           gender,
		Material:         strings.TrimSpace(request.Material),
		ShortDescription: strings.TrimSpace(request.ShortDescription),
		Description:      strings.TrimSpace(request.Description),
		Tags:             normalizeStrings(request.Tags),
		Price:            request.Price,
		DiscountPrice:    request.DiscountPrice,
		DiscountExpiresAt: request.DiscountExpiresAt,
		Currency:         currency,
		MainImageURL:     strings.TrimSpace(request.MainImageURL),
		Gallery:          normalizeStrings(request.Gallery),
		Type:             typeName,
		IsFeatured:       request.IsFeatured,
		IsActive:         true,
		Variants:         variants,
	}, nil
}

func toNewProductVariants(requests []productmodel.CreateProductVariantRequest) ([]productmodel.NewProductVariant, error) {
	variants := make([]productmodel.NewProductVariant, 0, len(requests))
	seenSKUs := make(map[string]struct{}, len(requests))
	seenCombinations := make(map[string]struct{}, len(requests))

	for index, request := range requests {
		sku := strings.ToUpper(strings.TrimSpace(request.SKU))
		size := strings.ToUpper(strings.TrimSpace(request.Size))
		color := strings.TrimSpace(request.Color)
		colorCode := strings.ToUpper(strings.TrimSpace(request.ColorCode))
		barcode := strings.TrimSpace(request.Barcode)

		switch {
		case sku == "":
			return nil, fmt.Errorf("%w: variant %d sku is required", ErrInvalidProduct, index)
		case size == "":
			return nil, fmt.Errorf("%w: variant %d size is required", ErrInvalidProduct, index)
		case color == "":
			return nil, fmt.Errorf("%w: variant %d color is required", ErrInvalidProduct, index)
		case request.Price != nil && *request.Price <= 0:
			return nil, fmt.Errorf("%w: variant %d price must be greater than zero", ErrInvalidProduct, index)
		}

		combination := strings.ToLower(size) + "\x00" + strings.ToLower(color)
		if _, exists := seenSKUs[sku]; exists {
			return nil, fmt.Errorf("%w: duplicate sku %q", ErrInvalidProduct, sku)
		}
		if _, exists := seenCombinations[combination]; exists {
			return nil, fmt.Errorf("%w: duplicate size/color combination", ErrInvalidProduct)
		}
		seenSKUs[sku] = struct{}{}
		seenCombinations[combination] = struct{}{}

		variants = append(variants, productmodel.NewProductVariant{
			PublicID:  uuid.New(),
			SKU:       sku,
			Size:      size,
			Color:     color,
			ColorCode: colorCode,
			Price:     request.Price,
			Barcode:   barcode,
		})
	}

	return variants, nil
}

func toNewProductVariant(request productmodel.CreateProductVariantRequest) (productmodel.NewProductVariant, error) {
	variants, err := toNewProductVariants([]productmodel.CreateProductVariantRequest{request})
	if err != nil {
		return productmodel.NewProductVariant{}, err
	}
	return variants[0], nil
}

func toUpdatedProductVariant(current productmodel.ProductVariant, request productmodel.UpdateProductVariantRequest) (productmodel.ProductVariant, error) {
	if request.SKU != nil {
		current.SKU = strings.ToUpper(strings.TrimSpace(*request.SKU))
	}
	if request.Size != nil {
		current.Size = strings.ToUpper(strings.TrimSpace(*request.Size))
	}
	if request.Color != nil {
		current.Color = strings.TrimSpace(*request.Color)
	}
	if request.ColorCode != nil {
		current.ColorCode = strings.ToUpper(strings.TrimSpace(*request.ColorCode))
	}
	if request.Price != nil {
		current.Price = request.Price
	}
	if request.Barcode != nil {
		current.Barcode = strings.TrimSpace(*request.Barcode)
	}

	switch {
	case current.SKU == "":
		return productmodel.ProductVariant{}, fmt.Errorf("%w: variant sku is required", ErrInvalidProduct)
	case current.Size == "":
		return productmodel.ProductVariant{}, fmt.Errorf("%w: variant size is required", ErrInvalidProduct)
	case current.Color == "":
		return productmodel.ProductVariant{}, fmt.Errorf("%w: variant color is required", ErrInvalidProduct)
	case current.Price != nil && *current.Price <= 0:
		return productmodel.ProductVariant{}, fmt.Errorf("%w: variant price must be greater than zero", ErrInvalidProduct)
	}

	return current, nil
}

func toUpdatedProduct(current productmodel.Product, request productmodel.UpdateProductRequest) (productmodel.Product, error) {
	if request.Name != nil {
		current.Name = strings.TrimSpace(*request.Name)
	}
	if request.Slug != nil {
		current.Slug = strings.ToLower(strings.TrimSpace(*request.Slug))
	}
	if request.Brand != nil {
		current.Brand = strings.TrimSpace(*request.Brand)
	}
	if request.Category != nil {
		current.Category = strings.ToLower(strings.TrimSpace(*request.Category))
	}
	if request.Gender != nil {
		current.Gender = strings.ToLower(strings.TrimSpace(*request.Gender))
	}
	if request.Material != nil {
		current.Material = strings.TrimSpace(*request.Material)
	}
	if request.ShortDescription != nil {
		current.ShortDescription = strings.TrimSpace(*request.ShortDescription)
	}
	if request.Description != nil {
		current.Description = strings.TrimSpace(*request.Description)
	}
	if request.Tags != nil {
		current.Tags = normalizeStrings(request.Tags)
	}
	if request.Price != nil {
		current.Price = *request.Price
	}
	if request.DiscountPrice != nil {
		current.DiscountPrice = request.DiscountPrice
	}
	if request.DiscountExpiresAt != nil {
		current.DiscountExpiresAt = request.DiscountExpiresAt
	}
	if request.Currency != nil {
		current.Currency = strings.ToUpper(strings.TrimSpace(*request.Currency))
	}
	if request.MainImageURL != nil {
		current.MainImageURL = strings.TrimSpace(*request.MainImageURL)
	}
	if request.Gallery != nil {
		current.Gallery = normalizeStrings(request.Gallery)
	}
	if request.Type != nil {
		current.Type = strings.ToLower(strings.TrimSpace(*request.Type))
	}
	if request.IsFeatured != nil {
		current.IsFeatured = *request.IsFeatured
	}

	switch {
	case current.Name == "", current.Slug == "", current.Brand == "", current.Category == "", current.Gender == "":
		return productmodel.Product{}, fmt.Errorf("%w: required product field is empty", ErrInvalidProduct)
	case current.Price <= 0:
		return productmodel.Product{}, fmt.Errorf("%w: price must be greater than zero", ErrInvalidProduct)
	case current.DiscountPrice != nil && (*current.DiscountPrice <= 0 || *current.DiscountPrice >= current.Price):
		return productmodel.Product{}, fmt.Errorf("%w: discount price must be lower than price", ErrInvalidProduct)
	case current.DiscountExpiresAt != nil && current.DiscountPrice == nil:
		return productmodel.Product{}, fmt.Errorf("%w: discount expiry requires a discount price", ErrInvalidProduct)
	case current.Currency == "":
		return productmodel.Product{}, fmt.Errorf("%w: currency is required", ErrInvalidProduct)
	}
	return current, nil
}

func calculateDiscountPercent(price int64, discountPrice *int64) *int32 {
	if discountPrice == nil || price <= 0 {
		return nil
	}
	percent := int32((price - *discountPrice) * 100 / price)
	return &percent
}

func toProductDetailResponse(storedProduct productmodel.Product) productmodel.ProductDetailResponse {
	return productmodel.ProductDetailResponse{
		PublicID:         storedProduct.PublicID.String(),
		Name:             storedProduct.Name,
		Slug:             storedProduct.Slug,
		Brand:            storedProduct.Brand,
		Category:         storedProduct.Category,
		Gender:           storedProduct.Gender,
		Material:         storedProduct.Material,
		ShortDescription: storedProduct.ShortDescription,
		Description:      storedProduct.Description,
		Tags:             storedProduct.Tags,
		Price:            storedProduct.Price,
		DiscountPrice:    storedProduct.DiscountPrice,
		DiscountPercent:  calculateDiscountPercent(storedProduct.Price, storedProduct.DiscountPrice),
		DiscountExpiresAt: storedProduct.DiscountExpiresAt,
		Currency:         storedProduct.Currency,
		MainImageURL:     storedProduct.MainImageURL,
		Gallery:          storedProduct.Gallery,
		Type:             storedProduct.Type,
		IsFeatured:       storedProduct.IsFeatured,
		Variants:         toProductVariantResponses(storedProduct.Variants),
	}
}

func toProductVariantResponses(variants []productmodel.ProductVariant) []productmodel.ProductVariantResponse {
	result := make([]productmodel.ProductVariantResponse, 0, len(variants))
	for _, variant := range variants {
		result = append(result, toProductVariantResponse(variant))
	}
	return result
}

func toProductVariantResponse(variant productmodel.ProductVariant) productmodel.ProductVariantResponse {
	return productmodel.ProductVariantResponse{
		PublicID:  variant.PublicID.String(),
		SKU:       variant.SKU,
		Size:      variant.Size,
		Color:     variant.Color,
		ColorCode: variant.ColorCode,
		Price:     variant.Price,
		Barcode:   variant.Barcode,
	}
}

func toProductListItemDTO(storedProduct productmodel.Product) productmodel.ProductListItemDTO {
	return productmodel.ProductListItemDTO{
		PublicID:      storedProduct.PublicID.String(),
		Name:          storedProduct.Name,
		Slug:          storedProduct.Slug,
		Brand:         storedProduct.Brand,
		Category:      storedProduct.Category,
		Gender:        storedProduct.Gender,
		Price:         storedProduct.Price,
		DiscountPrice: storedProduct.DiscountPrice,
		Currency:      storedProduct.Currency,
		MainImageURL:  storedProduct.MainImageURL,
	}
}

func normalizeStrings(values []string) []string {
	if len(values) == 0 {
		return nil
	}

	result := make([]string, 0, len(values))
	for _, value := range values {
		if value = strings.TrimSpace(value); value != "" {
			result = append(result, value)
		}
	}
	return result
}
