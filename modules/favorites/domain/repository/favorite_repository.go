package repository

import (
	"context"
	"go-modular/modules/favorites/domain/entity"
)

type FavoriteRepository interface {
	Create(ctx context.Context, favorite *entity.Favorite) error
	Delete(ctx context.Context, id uint) error
	FindByUserID(ctx context.Context, userID uint) ([]entity.Favorite, error)
}
