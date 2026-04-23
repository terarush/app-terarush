package service

import (
	"context"
	"errors"
	"go-modular/modules/blogs/domain/entity"
	"go-modular/modules/blogs/domain/repository"
)

// Errors
var (
	ErrBlogNotFound = errors.New("blog not found")
)

// BlogService handles blog domain logic
type BlogService struct {
	blogRepo repository.BlogRepository
}

// NewBlogService creates a new blog service
func NewBlogService(blogRepo repository.BlogRepository) *BlogService {
	return &BlogService{
		blogRepo: blogRepo,
	}
}

// GetAllBlogs gets all blogs
func (s *BlogService) GetAllBlogs(ctx context.Context) ([]*entity.Blog, error) {
	return s.blogRepo.FindAll(ctx)
}

// GetBlogByID gets a blog by ID
func (s *BlogService) GetBlogByID(ctx context.Context, id uint) (*entity.Blog, error) {
	blog, err := s.blogRepo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if blog == nil {
		return nil, ErrBlogNotFound
	}
	return blog, nil
}

// GetBlogBySlug gets a blog by slug and increments view count
func (s *BlogService) GetBlogBySlug(ctx context.Context, slug string) (*entity.Blog, error) {
	blog, err := s.blogRepo.FindBySlug(ctx, slug)
	if err != nil {
		if err == repository.ERR_RECORD_NOT_FOUND {
			return nil, ErrBlogNotFound
		}
		return nil, err
	}
	if blog == nil {
		return nil, ErrBlogNotFound
	}

	// Increment view count
	s.blogRepo.IncrementViewCount(ctx, blog.ID)

	return blog, nil
}

// GetPublishedBlogs gets all published blogs
func (s *BlogService) GetPublishedBlogs(ctx context.Context) ([]*entity.Blog, error) {
	return s.blogRepo.FindPublished(ctx)
}

// CreateBlog creates a new blog
func (s *BlogService) CreateBlog(ctx context.Context, blog *entity.Blog) error {
	return s.blogRepo.Create(ctx, blog)
}

// UpdateBlog updates a blog
func (s *BlogService) UpdateBlog(ctx context.Context, blog *entity.Blog) error {
	existingBlog, err := s.blogRepo.FindByID(ctx, blog.ID)
	if err != nil {
		return err
	}
	if existingBlog == nil {
		return ErrBlogNotFound
	}

	return s.blogRepo.Update(ctx, blog)
}

// DeleteBlog deletes a blog
func (s *BlogService) DeleteBlog(ctx context.Context, id uint) error {
	existingBlog, err := s.blogRepo.FindByID(ctx, id)
	if err != nil {
		return err
	}
	if existingBlog == nil {
		return ErrBlogNotFound
	}

	return s.blogRepo.Delete(ctx, id)
}
