package repository

import (
	"context"
	"errors"
	"go-modular/internal/pkg/database"
	"go-modular/modules/blogs/domain/entity"
)

var (
	ERR_RECORD_NOT_FOUND = errors.New("record not found")
)

type BlogRepositoryImpl struct{}

// Create implements BlogRepository.
func (r BlogRepositoryImpl) Create(ctx context.Context, blog *entity.Blog) error {
	return database.DB.WithContext(ctx).Create(blog).Error
}

// Delete implements BlogRepository.
func (r BlogRepositoryImpl) Delete(ctx context.Context, id uint) error {
	return database.DB.WithContext(ctx).Delete(&entity.Blog{}, id).Error
}

// FindAll finds all blogs
func (r BlogRepositoryImpl) FindAll(ctx context.Context) ([]*entity.Blog, error) {
	var blogs []*entity.Blog
	result := database.DB.WithContext(ctx).Order("created_at DESC").Find(&blogs)
	if result.Error != nil {
		return nil, result.Error
	}
	return blogs, nil
}

// FindByID implements BlogRepository.
func (r BlogRepositoryImpl) FindByID(ctx context.Context, id uint) (*entity.Blog, error) {
	var blog entity.Blog
	result := database.DB.WithContext(ctx).First(&blog, id)
	if result.Error != nil {
		return nil, result.Error
	}
	return &blog, nil
}

// FindBySlug implements BlogRepository.
func (r BlogRepositoryImpl) FindBySlug(ctx context.Context, slug string) (*entity.Blog, error) {
	var blog entity.Blog
	result := database.DB.WithContext(ctx).Where("slug = ?", slug).First(&blog)
	if result.Error != nil {
		if result.RowsAffected == 0 {
			return nil, ERR_RECORD_NOT_FOUND
		}
		return nil, result.Error
	}
	return &blog, nil
}

// FindPublished finds all published blogs
func (r BlogRepositoryImpl) FindPublished(ctx context.Context) ([]*entity.Blog, error) {
	var blogs []*entity.Blog
	result := database.DB.WithContext(ctx).Where("is_published = ?", true).Order("published_at DESC").Find(&blogs)
	if result.Error != nil {
		return nil, result.Error
	}
	return blogs, nil
}

// Update implements BlogRepository.
func (r BlogRepositoryImpl) Update(ctx context.Context, blog *entity.Blog) error {
	return database.DB.WithContext(ctx).Save(blog).Error
}

// IncrementViewCount increments the view count for a blog
func (r BlogRepositoryImpl) IncrementViewCount(ctx context.Context, id uint) error {
	return database.DB.WithContext(ctx).Model(&entity.Blog{}).Where("id = ?", id).Update("view_count", "view_count + 1").Error
}

func NewBlogRepositoryImpl() BlogRepository {
	return BlogRepositoryImpl{}
}
