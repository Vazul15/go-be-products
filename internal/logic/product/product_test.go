package product

import (
	"context"
	"errors"
	"testing"

	"github.com/google/uuid"

	productmodel "products/internal/model/product"
)

type repositoryStub struct {
	page         productmodel.ProductPage
	product      productmodel.Product
	createDraft  productmodel.NewProduct
	listOffset   int
	listLimit    int
	createCalled bool
	createErr    error
}

func (stub *repositoryStub) ListProducts(_ context.Context, offset, limit int, _ productmodel.ProductFilter) (productmodel.ProductPage, error) {
	stub.listOffset = offset
	stub.listLimit = limit
	return stub.page, nil
}

func (stub *repositoryStub) GetProduct(_ context.Context, _ uuid.UUID) (productmodel.Product, error) {
	return stub.product, nil
}

func (stub *repositoryStub) CreateProduct(_ context.Context, draft productmodel.NewProduct) (productmodel.Product, error) {
	stub.createCalled = true
	stub.createDraft = draft
	if stub.createErr != nil {
		return productmodel.Product{}, stub.createErr
	}
	return stub.product, nil
}

func (stub *repositoryStub) UpdateProduct(_ context.Context, _ uuid.UUID, storedProduct productmodel.Product) (productmodel.Product, error) {
	return storedProduct, nil
}

func (stub *repositoryStub) DeleteProduct(_ context.Context, _ uuid.UUID) error {
	return nil
}

func (stub *repositoryStub) CreateProductVariant(_ context.Context, _ uuid.UUID, variant productmodel.NewProductVariant) (productmodel.ProductVariant, error) {
	return productmodel.ProductVariant{
		PublicID:  variant.PublicID,
		SKU:       variant.SKU,
		Size:      variant.Size,
		Color:     variant.Color,
		ColorCode: variant.ColorCode,
		Price:     variant.Price,
		Barcode:   variant.Barcode,
	}, nil
}

func (stub *repositoryStub) ListProductVariants(_ context.Context, _ uuid.UUID) ([]productmodel.ProductVariant, error) {
	return nil, nil
}

func (stub *repositoryStub) GetProductVariant(_ context.Context, publicID uuid.UUID) (productmodel.ProductVariant, error) {
	return productmodel.ProductVariant{PublicID: publicID}, nil
}

func (stub *repositoryStub) UpdateProductVariant(_ context.Context, _ uuid.UUID, variant productmodel.ProductVariant) (productmodel.ProductVariant, error) {
	return variant, nil
}

func (stub *repositoryStub) DeleteProductVariant(_ context.Context, _ uuid.UUID) error {
	return nil
}

func TestCreateProductMapsAndNormalizesInLogic(t *testing.T) {
	publicID := uuid.New()
	stub := &repositoryStub{product: productmodel.Product{
		PublicID:      publicID,
		Price:         10000,
		DiscountPrice: pointerToInt64(8000),
	}}
	orchestrator, err := New(stub)
	if err != nil {
		t.Fatalf("New() error = %v", err)
	}

	response, err := orchestrator.CreateProduct(context.Background(), productmodel.CreateProductRequest{
		Name:          "  Blue Shirt  ",
		Slug:          "  Blue-Shirt  ",
		Brand:         "  Acme  ",
		Category:      " TOPS ",
		Gender:        " UNISEX ",
		Description:   " A comfortable everyday shirt. ",
		Price:         10000,
		DiscountPrice: pointerToInt64(8000),
		Currency:      " huf ",
		MainImageURL:  "https://example.com/blue-shirt.jpg",
		Variants: []productmodel.CreateProductVariantRequest{{
			SKU:   " acme-blue-m ",
			Size:  " m ",
			Color: " blue ",
		}},
		Type: " SHIRT ",
		Tags: []string{" summer ", "", " cotton "},
	})
	if err != nil {
		t.Fatalf("CreateProduct() error = %v", err)
	}

	if !stub.createCalled {
		t.Fatal("CreateProduct() did not call repository")
	}
	if stub.createDraft.Slug != "blue-shirt" || stub.createDraft.Category != "tops" || stub.createDraft.Gender != "unisex" {
		t.Fatalf("draft normalization = %#v", stub.createDraft)
	}
	if stub.createDraft.Currency != "HUF" {
		t.Fatalf("draft business values = %#v", stub.createDraft)
	}
	if response.DiscountPercent == nil || *response.DiscountPercent != 20 {
		t.Fatalf("discount percent = %v, want 20", response.DiscountPercent)
	}
	if len(stub.createDraft.Variants) != 1 || stub.createDraft.Variants[0].SKU != "ACME-BLUE-M" {
		t.Fatalf("draft variants = %#v", stub.createDraft.Variants)
	}
	if len(stub.createDraft.Tags) != 2 || stub.createDraft.Tags[0] != "summer" || stub.createDraft.Tags[1] != "cotton" {
		t.Fatalf("draft tags = %#v", stub.createDraft.Tags)
	}
	if response.PublicID != publicID.String() {
		t.Fatalf("response public id = %q, want %q", response.PublicID, publicID)
	}
}

func TestCreateProductRejectsInvalidDiscountBeforeRepository(t *testing.T) {
	stub := &repositoryStub{}
	orchestrator, err := New(stub)
	if err != nil {
		t.Fatalf("New() error = %v", err)
	}

	_, err = orchestrator.CreateProduct(context.Background(), productmodel.CreateProductRequest{
		Name:          "Blue Shirt",
		Slug:          "blue-shirt",
		Brand:         "Acme",
		Category:      "tops",
		Gender:        "unisex",
		Description:   "A comfortable everyday shirt.",
		Price:         10000,
		DiscountPrice: pointerToInt64(10000),
		Currency:      "HUF",
		MainImageURL:  "https://example.com/blue-shirt.jpg",
		Variants: []productmodel.CreateProductVariantRequest{{
			SKU:   "ACME-BLUE-M",
			Size:  "M",
			Color: "blue",
		}},
		Type: "shirt",
	})
	if !errors.Is(err, ErrInvalidProduct) {
		t.Fatalf("error = %v, want ErrInvalidProduct", err)
	}
	if stub.createCalled {
		t.Fatal("repository was called for invalid product")
	}
}

func TestListProductsMapsPageInLogic(t *testing.T) {
	publicID := uuid.New()
	stub := &repositoryStub{
		page: productmodel.ProductPage{
			Total: 21,
			Items: []productmodel.Product{{
				PublicID: publicID,
				Name:     "Blue Shirt",
				Brand:    "Acme",
			}},
		},
	}
	orchestrator, err := New(stub)
	if err != nil {
		t.Fatalf("New() error = %v", err)
	}

	response, err := orchestrator.ListProducts(context.Background(), 2, 10, productmodel.ProductFilter{})
	if err != nil {
		t.Fatalf("ListProducts() error = %v", err)
	}
	if stub.listOffset != 10 || stub.listLimit != 10 {
		t.Fatalf("repository pagination = (%d, %d), want (10, 10)", stub.listOffset, stub.listLimit)
	}
	if response.TotalPages != 3 || !response.HasNext || len(response.Items) != 1 {
		t.Fatalf("response pagination = %#v", response)
	}
	if response.Items[0].PublicID != publicID.String() {
		t.Fatalf("response item id = %q, want %q", response.Items[0].PublicID, publicID)
	}
}

func pointerToInt64(value int64) *int64 {
	return &value
}
