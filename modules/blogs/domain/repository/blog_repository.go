package repository

import (
	"context"
	"go-modular/modules/blogs/domain/entity"
	userEntity "go-modular/modules/users/domain/entity"
)

// BlogRepository defines the blog repository interface
type BlogRepository interface {
	FindAll(ctx context.Context) ([]*entity.Blog, error)
	FindAllWithPagination(ctx context.Context, page, pageSize int) ([]*entity.Blog, int64, error)
	FindAllWithPaginationAndSearch(ctx context.Context, page, pageSize int, search string) ([]*entity.Blog, int64, error)
	FindByID(ctx context.Context, id uint) (*entity.Blog, error)
	FindBySlug(ctx context.Context, slug string) (*entity.Blog, error)
	FindBySlugWithUser(ctx context.Context, slug string) (*entity.Blog, *userEntity.User, error)
	FindPublished(ctx context.Context) ([]*entity.Blog, error)
	FindPublishedWithPagination(ctx context.Context, page, pageSize int) ([]*entity.Blog, int64, error)
	FindPublishedWithPaginationAndSearch(ctx context.Context, page, pageSize int, search string) ([]*entity.Blog, int64, error)
	Create(ctx context.Context, blog *entity.Blog) error
	Update(ctx context.Context, blog *entity.Blog) error
	Delete(ctx context.Context, id uint) error
	IncrementViewCount(ctx context.Context, id uint) error
}
