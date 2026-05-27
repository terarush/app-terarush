package service

import (
	"context"
	"errors"
	"go-modular/modules/favorites/domain/entity"
	"go-modular/modules/favorites/domain/repository"
)

var (
	ERR_COMMENT_NOT_FOUND = errors.New("favorite not found")
)

type FavoriteService struct {
	favoriteRepo repository.FavoriteRepository
}

func NewFavoriteService(favoriteRepo repository.FavoriteRepository) *FavoriteService {
	return &FavoriteService{favoriteRepo: favoriteRepo}
}

func (s *FavoriteService) CreateFavorite(ctx context.Context, favorite *entity.Favorite, userId uint) error {
	favorite.UserID = userId
	return s.favoriteRepo.Create(ctx, favorite)
}

func (s *FavoriteService) DeleteFavorite(ctx context.Context, id uint) error {
	return s.favoriteRepo.Delete(ctx, id)
}

func (s *FavoriteService) FindFavoritesByUserID(ctx context.Context, userId uint) ([]entity.Favorite, error) {
	return s.favoriteRepo.FindByUserID(ctx, userId)
}
