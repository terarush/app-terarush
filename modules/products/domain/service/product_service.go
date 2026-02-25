package service

import (
	"context"
	"errors"
	"go-modular/modules/products/domain/entity"
	"go-modular/modules/products/domain/repository"
)

var (
	ErrProductNotFound     = errors.New("product not found")
	ErrProductNotAvailable = errors.New("product not available")
	ErrInsufficientStock   = errors.New("insufficient stock")
)

type ProductService struct {
	repo repository.ProductRepository
}

// NewProductService creates a new product service
func NewProductService(repo repository.ProductRepository) *ProductService {
	return &ProductService{repo: repo}
}

// CreateProduct creates a new product (Admin only)
func (s *ProductService) CreateProduct(ctx context.Context, product *entity.Product) error {
	return s.repo.Create(ctx, product)
}

// UpdateProduct updates an existing product (Admin only)
func (s *ProductService) UpdateProduct(ctx context.Context, product *entity.Product) error {
	// Check if product exists
	existing, err := s.repo.FindByID(ctx, product.ID)
	if err != nil {
		if err == repository.ERR_RECORD_NOT_FOUND {
			return ErrProductNotFound
		}
		return err
	}

	// Update fields
	existing.Name = product.Name
	existing.Description = product.Description
	existing.Price = product.Price
	existing.CPUCores = product.CPUCores
	existing.RAMMB = product.RAMMB
	existing.StorageGB = product.StorageGB
	existing.BandwidthGB = product.BandwidthGB
	existing.IsActive = product.IsActive
	existing.Stock = product.Stock
	existing.ImageURL = product.ImageURL

	return s.repo.Update(ctx, existing)
}

// DeleteProduct soft deletes a product (Admin only)
func (s *ProductService) DeleteProduct(ctx context.Context, id uint) error {
	err := s.repo.Delete(ctx, id)
	if err != nil {
		if err == repository.ERR_RECORD_NOT_FOUND {
			return ErrProductNotFound
		}
		return err
	}
	return nil
}

// GetProductByID gets a product by ID
func (s *ProductService) GetProductByID(ctx context.Context, id uint) (*entity.Product, error) {
	product, err := s.repo.FindByID(ctx, id)
	if err != nil {
		if err == repository.ERR_RECORD_NOT_FOUND {
			return nil, ErrProductNotFound
		}
		return nil, err
	}
	return product, nil
}

// GetAllProducts gets all products (for admin)
func (s *ProductService) GetAllProducts(ctx context.Context, includeInactive bool) ([]*entity.Product, error) {
	return s.repo.FindAll(ctx, includeInactive)
}

// GetActiveProducts gets all active products (for users)
func (s *ProductService) GetActiveProducts(ctx context.Context) ([]*entity.Product, error) {
	return s.repo.FindActive(ctx)
}

// ValidateProductForPurchase checks if product can be purchased
func (s *ProductService) ValidateProductForPurchase(ctx context.Context, productID uint, quantity int) (*entity.Product, error) {
	product, err := s.repo.FindByID(ctx, productID)
	if err != nil {
		if err == repository.ERR_RECORD_NOT_FOUND {
			return nil, ErrProductNotFound
		}
		return nil, err
	}

	if !product.IsAvailable() {
		return nil, ErrProductNotAvailable
	}

	if product.Stock != -1 && product.Stock < quantity {
		return nil, ErrInsufficientStock
	}

	return product, nil
}
