package repository

import (
	"context"
	"errors"
	"go-modular/internal/pkg/database"
	"go-modular/modules/assets/domain/entity"
)

var (
	ERR_RECORD_NOT_FOUND = errors.New("record not found")
)

type AssetRepositoryImpl struct{}

// Create creates a new asset record
func (r AssetRepositoryImpl) Create(ctx context.Context, asset *entity.Asset) error {
	return database.DB.WithContext(ctx).Create(asset).Error
}

// FindByID finds an asset by ID
func (r AssetRepositoryImpl) FindByID(ctx context.Context, id uint) (*entity.Asset, error) {
	var asset entity.Asset
	result := database.DB.WithContext(ctx).First(&asset, id)
	if result.Error != nil {
		return nil, result.Error
	}
	return &asset, nil
}

// FindByURL finds an asset by URL
func (r AssetRepositoryImpl) FindByURL(ctx context.Context, url string) (*entity.Asset, error) {
	var asset entity.Asset
	result := database.DB.WithContext(ctx).Where("url = ?", url).First(&asset)
	if result.Error != nil {
		if result.RowsAffected == 0 {
			return nil, ERR_RECORD_NOT_FOUND
		}
		return nil, result.Error
	}
	return &asset, nil
}

// Delete deletes an asset
func (r AssetRepositoryImpl) Delete(ctx context.Context, id uint) error {
	return database.DB.WithContext(ctx).Delete(&entity.Asset{}, id).Error
}

func NewAssetRepositoryImpl() AssetRepository {
	return AssetRepositoryImpl{}
}
