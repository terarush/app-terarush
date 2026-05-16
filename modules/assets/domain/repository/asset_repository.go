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
	Delete(ctx context.Context, id uint) error
}
