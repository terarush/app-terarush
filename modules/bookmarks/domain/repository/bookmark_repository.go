package repository

import (
	"context"
	"go-modular/modules/bookmarks/domain/entity"
)

type BookmarkRepository interface {
	FindAll(ctx context.Context, filter map[string]any) ([]*entity.Bookmark, error)
	FindByID(ctx context.Context, id uint) (*entity.Bookmark, error)
	FindByUserID(ctx context.Context, userID uint, status string) ([]*entity.Bookmark, error)
	FindByUserAndResource(ctx context.Context, userID uint, resourceType string, resourceID uint) (*entity.Bookmark, error)
	FindByCollectionID(ctx context.Context, collectionID uint) ([]*entity.Bookmark, error)
	Create(ctx context.Context, b *entity.Bookmark) error
	Update(ctx context.Context, b *entity.Bookmark) error
	Delete(ctx context.Context, id uint) error
	CountByUserID(ctx context.Context, userID uint) (int64, error)
}

type CollectionRepository interface {
	FindAll(ctx context.Context, userID uint) ([]*entity.Collection, error)
	FindByID(ctx context.Context, id uint) (*entity.Collection, error)
	Create(ctx context.Context, c *entity.Collection) error
	Update(ctx context.Context, c *entity.Collection) error
	Delete(ctx context.Context, id uint) error
}

type ReadingGoalRepository interface {
	FindByUserAndYear(ctx context.Context, userID uint, year int) (*entity.ReadingGoal, error)
	Upsert(ctx context.Context, g *entity.ReadingGoal) error
	Increment(ctx context.Context, userID uint, year int) error
}
