package repository

import (
	"context"
	"errors"
	"go-modular/modules/products/domain/entity"

	"gorm.io/gorm"
)

var (
	ERR_RECORD_NOT_FOUND = errors.New("record not found")
)

type productRepositoryImpl struct {
	db *gorm.DB
}

// NewProductRepository creates a new product repository
func NewProductRepository(db *gorm.DB) ProductRepository {
	return &productRepositoryImpl{db: db}
}

// Create inserts a new product
func (r *productRepositoryImpl) Create(ctx context.Context, product *entity.Product) error {
	return r.db.WithContext(ctx).Create(product).Error
}

// Update updates an existing product
func (r *productRepositoryImpl) Update(ctx context.Context, product *entity.Product) error {
	return r.db.WithContext(ctx).Save(product).Error
}

// Delete soft deletes a product
func (r *productRepositoryImpl) Delete(ctx context.Context, id uint) error {
	result := r.db.WithContext(ctx).Delete(&entity.Product{}, id)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return ERR_RECORD_NOT_FOUND
	}
	return nil
}

// FindByID finds a product by ID
func (r *productRepositoryImpl) FindByID(ctx context.Context, id uint) (*entity.Product, error) {
	var product entity.Product
	err := r.db.WithContext(ctx).First(&product, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ERR_RECORD_NOT_FOUND
		}
		return nil, err
	}
	return &product, nil
}

// FindAll finds all products (including inactive if specified)
func (r *productRepositoryImpl) FindAll(ctx context.Context, includeInactive bool) ([]*entity.Product, error) {
	var products []*entity.Product
	query := r.db.WithContext(ctx)

	if !includeInactive {
		query = query.Where("is_active = ?", true)
	}

	err := query.Order("created_at DESC").Find(&products).Error
	if err != nil {
		return nil, err
	}
	return products, nil
}

// FindActive finds all active products
func (r *productRepositoryImpl) FindActive(ctx context.Context) ([]*entity.Product, error) {
	var products []*entity.Product
	err := r.db.WithContext(ctx).
		Where("is_active = ?", true).
		Where("deleted_at IS NULL").
		Order("price ASC").
		Find(&products).Error

	if err != nil {
		return nil, err
	}
	return products, nil
}

// UpdateStock updates product stock
func (r *productRepositoryImpl) UpdateStock(ctx context.Context, id uint, quantity int) error {
	return r.db.WithContext(ctx).
		Model(&entity.Product{}).
		Where("id = ?", id).
		UpdateColumn("stock", gorm.Expr("stock - ?", quantity)).
		Error
}
