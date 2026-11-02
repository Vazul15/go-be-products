package product

import (
	"fmt"

	productrepo "products/internal/connector/db/productRepository"
)

// New creates a ProductOrchestrator with the required dependencies.
func New(repo productrepo.Repository) (*ProductOrchestrator, error) {
	if repo == nil {
		return nil, fmt.Errorf("product: repository must not be nil")
	}
	return &ProductOrchestrator{repo: repo}, nil
}

// MustNew creates a ProductOrchestrator and panics if initialization fails.
func MustNew(repo productrepo.Repository) *ProductOrchestrator {
	o, err := New(repo)
	if err != nil {
		panic(err)
	}
	return o
}
