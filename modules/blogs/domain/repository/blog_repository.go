package repository

import (
	"context"
	"go-modular/modules/blogs/domain/entity"
)

// BlogRepository defines the blog repository interface
type BlogRepository interface {
	FindAll(ctx context.Context) ([]*entity.Blog, error)
	FindByID(ctx context.Context, id uint) (*entity.Blog, error)
	FindBySlug(ctx context.Context, slug string) (*entity.Blog, error)
	FindPublished(ctx context.Context) ([]*entity.Blog, error)
	Create(ctx context.Context, blog *entity.Blog) error
	Update(ctx context.Context, blog *entity.Blog) error
	Delete(ctx context.Context, id uint) error
	IncrementViewCount(ctx context.Context, id uint) error
}
