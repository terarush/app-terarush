package service

import (
	"context"
	"errors"
	"go-modular/modules/assets/domain/entity"
	"go-modular/modules/assets/domain/repository"
)

var (
	ErrAssetNotFound = errors.New("asset not found")
)

// AssetService handles asset domain logic
type AssetService struct {
	assetRepo repository.AssetRepository
}

// NewAssetService creates a new asset service
func NewAssetService(assetRepo repository.AssetRepository) *AssetService {
	return &AssetService{
		assetRepo: assetRepo,
	}
}

// CreateAsset creates a new asset record
func (s *AssetService) CreateAsset(ctx context.Context, asset *entity.Asset) error {
	return s.assetRepo.Create(ctx, asset)
}

// GetAssetByID gets an asset by ID
func (s *AssetService) GetAssetByID(ctx context.Context, id uint) (*entity.Asset, error) {
	asset, err := s.assetRepo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if asset == nil {
		return nil, ErrAssetNotFound
	}
	return asset, nil
}

// GetAssetByURL gets an asset by URL
func (s *AssetService) GetAssetByURL(ctx context.Context, url string) (*entity.Asset, error) {
	asset, err := s.assetRepo.FindByURL(ctx, url)
	if err != nil {
		if err == repository.ERR_RECORD_NOT_FOUND {
			return nil, ErrAssetNotFound
		}
		return nil, err
	}
	if asset == nil {
		return nil, ErrAssetNotFound
	}
	return asset, nil
}

// GetAllAssets gets all assets with pagination
func (s *AssetService) GetAllAssets(ctx context.Context, page, pageSize int) ([]*entity.Asset, int64, error) {
	return s.assetRepo.FindAll(ctx, page, pageSize)
}

// DeleteAsset deletes an asset
func (s *AssetService) DeleteAsset(ctx context.Context, id uint) error {
	asset, err := s.assetRepo.FindByID(ctx, id)
	if err != nil {
		return err
	}
	if asset == nil {
		return ErrAssetNotFound
	}
	return s.assetRepo.Delete(ctx, id)
}
