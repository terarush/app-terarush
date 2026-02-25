package repository

import (
	"context"
	"go-modular/modules/products/domain/entity"
)

// ProductRepository defines the interface for product data operations
type ProductRepository interface {
	// Admin operations
	Create(ctx context.Context, product *entity.Product) error
	Update(ctx context.Context, product *entity.Product) error
	Delete(ctx context.Context, id uint) error

	// Read operations
	FindByID(ctx context.Context, id uint) (*entity.Product, error)
	FindAll(ctx context.Context, includeInactive bool) ([]*entity.Product, error)
	FindActive(ctx context.Context) ([]*entity.Product, error)

	// Stock management
	UpdateStock(ctx context.Context, id uint, quantity int) error
}
