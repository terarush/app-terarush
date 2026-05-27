package service

import (
	"errors"
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
