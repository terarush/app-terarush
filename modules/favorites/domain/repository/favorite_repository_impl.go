package repository

import (
	"context"
	"errors"
	"go-modular/internal/pkg/database"
	"go-modular/modules/favorites/domain/entity"
)

var (
	ERR_RECORD_NOT_FOUND = errors.New("record not found")
)

type FavoriteRepositoryImpl struct{}

func (r *FavoriteRepositoryImpl) Create(ctx context.Context, favorite *entity.Favorite) error {
	return database.DB.WithContext(ctx).Create(favorite).Error
}

func (r *FavoriteRepositoryImpl) Delete(ctx context.Context, id uint) error {
	return database.DB.WithContext(ctx).Delete(&entity.Favorite{}, "id = ?", id).Error
}

func (r *FavoriteRepositoryImpl) FindByUserID(ctx context.Context, userID uint) ([]entity.Favorite, error) {
	var favorites []entity.Favorite
	return favorites, database.DB.WithContext(ctx).Where("user_id = ?", userID).Find(&favorites).Error
}

func NewFavoriteRepositoryImpl() FavoriteRepository {
	return &FavoriteRepositoryImpl{}
}
