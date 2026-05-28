package repository

import (
	"context"
	"go-modular/modules/assets/domain/entity"
)

// AssetRepository defines the asset repository interface
type AssetRepository interface {
	Create(ctx context.Context, asset *entity.Asset) error
	FindByID(ctx context.Context, id uint) (*entity.Asset, error)
	FindByURL(ctx context.Context, url string) (*entity.Asset, error)
	FindAll(ctx context.Context, page, pageSize int) ([]*entity.Asset, int64, error)
	Delete(ctx context.Context, id uint) error
}
