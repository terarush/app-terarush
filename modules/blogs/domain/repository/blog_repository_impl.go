package repository

import (
	"context"
	"errors"
	"go-modular/internal/pkg/database"
	"go-modular/modules/blogs/domain/entity"
	userEntity "go-modular/modules/users/domain/entity"
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

// FindAllWithPagination finds all blogs with pagination
func (r BlogRepositoryImpl) FindAllWithPagination(ctx context.Context, page, pageSize int) ([]*entity.Blog, int64, error) {
	var blogs []*entity.Blog
	var total int64

	// Get total count
	if err := database.DB.WithContext(ctx).Model(&entity.Blog{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Calculate offset
	offset := (page - 1) * pageSize

	// Get paginated results
	result := database.DB.WithContext(ctx).
		Order("created_at DESC").
		Offset(offset).
		Limit(pageSize).
		Find(&blogs)

	if result.Error != nil {
		return nil, 0, result.Error
	}

	return blogs, total, nil
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

// FindBySlugWithUser finds a blog by slug and joins with user data
func (r BlogRepositoryImpl) FindBySlugWithUser(ctx context.Context, slug string) (*entity.Blog, *userEntity.User, error) {
	var blog entity.Blog
	var user userEntity.User

	result := database.DB.WithContext(ctx).
		Where("slug = ?", slug).
		Joins("LEFT JOIN users ON blogs.user_id = users.id").
		Select("blogs.*").
		First(&blog)

	if result.Error != nil {
		if result.RowsAffected == 0 {
			return nil, nil, ERR_RECORD_NOT_FOUND
		}
		return nil, nil, result.Error
	}

	// Fetch user if user_id is set
	var userPtr *userEntity.User
	if blog.UserID > 0 {
		userResult := database.DB.WithContext(ctx).First(&user, blog.UserID)
		if userResult.Error == nil {
			userPtr = &user
		}
	}

	return &blog, userPtr, nil
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

// FindPublishedWithPagination finds published blogs with pagination
func (r BlogRepositoryImpl) FindPublishedWithPagination(ctx context.Context, page, pageSize int) ([]*entity.Blog, int64, error) {
	var blogs []*entity.Blog
	var total int64

	// Get total count of published blogs
	if err := database.DB.WithContext(ctx).Model(&entity.Blog{}).Where("is_published = ?", true).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Calculate offset
	offset := (page - 1) * pageSize

	// Get paginated results
	result := database.DB.WithContext(ctx).
		Where("is_published = ?", true).
		Order("published_at DESC").
		Offset(offset).
		Limit(pageSize).
		Find(&blogs)

	if result.Error != nil {
		return nil, 0, result.Error
	}

	return blogs, total, nil
}

// FindAllWithPaginationAndSearch finds all blogs with pagination and search
func (r BlogRepositoryImpl) FindAllWithPaginationAndSearch(ctx context.Context, page, pageSize int, search string) ([]*entity.Blog, int64, error) {
	var blogs []*entity.Blog
	var total int64

	query := database.DB.WithContext(ctx)
	
	// Apply search filter if provided
	if search != "" {
		query = query.Where("title LIKE ? OR excerpt LIKE ? OR content LIKE ?", "%"+search+"%", "%"+search+"%", "%"+search+"%")
	}

	// Get total count
	if err := query.Model(&entity.Blog{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Calculate offset
	offset := (page - 1) * pageSize

	// Get paginated results
	result := query.
		Order("created_at DESC").
		Offset(offset).
		Limit(pageSize).
		Find(&blogs)

	if result.Error != nil {
		return nil, 0, result.Error
	}

	return blogs, total, nil
}

// FindPublishedWithPaginationAndSearch finds published blogs with pagination and search
func (r BlogRepositoryImpl) FindPublishedWithPaginationAndSearch(ctx context.Context, page, pageSize int, search string) ([]*entity.Blog, int64, error) {
	var blogs []*entity.Blog
	var total int64

	query := database.DB.WithContext(ctx).Where("is_published = ?", true)
	
	// Apply search filter if provided
	if search != "" {
		query = query.Where("title LIKE ? OR excerpt LIKE ? OR content LIKE ?", "%"+search+"%", "%"+search+"%", "%"+search+"%")
	}

	// Get total count of published blogs matching search
	if err := query.Model(&entity.Blog{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Calculate offset
	offset := (page - 1) * pageSize

	// Get paginated results
	result := query.
		Order("published_at DESC").
		Offset(offset).
		Limit(pageSize).
		Find(&blogs)

	if result.Error != nil {
		return nil, 0, result.Error
	}

	return blogs, total, nil
}

// Update implements BlogRepository.
func (r BlogRepositoryImpl) Update(ctx context.Context, blog *entity.Blog) error {
	return database.DB.WithContext(ctx).Save(blog).Error
}

// IncrementViewCount increments the view count for a blog
func (r BlogRepositoryImpl) IncrementViewCount(ctx context.Context, id uint) error {
	return database.DB.WithContext(ctx).Model(&entity.Blog{}).Where("id = ?", id).Update("view_count", database.DB.Raw("view_count + 1")).Error
}

func NewBlogRepositoryImpl() BlogRepository {
	return BlogRepositoryImpl{}
}
