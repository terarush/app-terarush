package response

import (
	"time"

	"go-modular/modules/blogs/domain/entity"
	userEntity "go-modular/modules/users/domain/entity"
)

// UserInfo represents basic user information for blog context
type UserInfo struct {
	ID     uint   `json:"id"`
	Name   string `json:"name"`
	Email  string `json:"email"`
	Avatar string `json:"avatar"`
	Bio    string `json:"bio"`
}

// BlogResponse is the response for a blog
type BlogResponse struct {
	ID          uint       `json:"id"`
	Title       string     `json:"title"`
	Slug        string     `json:"slug"`
	Content     string     `json:"content"`
	Excerpt     string     `json:"excerpt"`
	Author      string     `json:"author"`
	Category    string     `json:"category"`
	Tags        string     `json:"tags"`
	Image       string     `json:"image"`
	IsPublished bool       `json:"is_published"`
	ViewCount   int64      `json:"view_count"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
	PublishedAt *time.Time `json:"published_at"`
}

// BlogDetailResponse is the detailed response for a blog with user data
type BlogDetailResponse struct {
	ID          uint       `json:"id"`
	Title       string     `json:"title"`
	Slug        string     `json:"slug"`
	Content     string     `json:"content"`
	Excerpt     string     `json:"excerpt"`
	Author      string     `json:"author"`
	User        *UserInfo  `json:"user"`
	Category    string     `json:"category"`
	Tags        string     `json:"tags"`
	Image       string     `json:"image"`
	IsPublished bool       `json:"is_published"`
	ViewCount   int64      `json:"view_count"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
	PublishedAt *time.Time `json:"published_at"`
}

// FromEntity converts a blog entity to a response
func FromEntity(blog *entity.Blog) *BlogResponse {
	if blog == nil {
		return nil
	}

	return &BlogResponse{
		ID:          blog.ID,
		Title:       blog.Title,
		Slug:        blog.Slug,
		Content:     blog.Content,
		Excerpt:     blog.Excerpt,
		Author:      blog.Author,
		Category:    blog.Category,
		Tags:        blog.Tags,
		Image:       blog.Image,
		IsPublished: blog.IsPublished,
		ViewCount:   blog.ViewCount,
		CreatedAt:   blog.CreatedAt,
		UpdatedAt:   blog.UpdatedAt,
		PublishedAt: blog.PublishedAt,
	}
}

// FromEntityWithUser converts a blog entity with user to a detailed response
func FromEntityWithUser(blog *entity.Blog, user *userEntity.User) *BlogDetailResponse {
	if blog == nil {
		return nil
	}

	var userInfo *UserInfo
	if user != nil {
		userInfo = &UserInfo{
			ID:     user.ID,
			Name:   user.Name,
			Email:  user.Email,
			Avatar: user.Avatar,
			Bio:    user.Bio,
		}
	}

	return &BlogDetailResponse{
		ID:          blog.ID,
		Title:       blog.Title,
		Slug:        blog.Slug,
		Content:     blog.Content,
		Excerpt:     blog.Excerpt,
		Author:      blog.Author,
		User:        userInfo,
		Category:    blog.Category,
		Tags:        blog.Tags,
		Image:       blog.Image,
		IsPublished: blog.IsPublished,
		ViewCount:   blog.ViewCount,
		CreatedAt:   blog.CreatedAt,
		UpdatedAt:   blog.UpdatedAt,
		PublishedAt: blog.PublishedAt,
	}
}

// FromEntities converts multiple blog entities to responses
func FromEntities(blogs []*entity.Blog) []*BlogResponse {
	responses := make([]*BlogResponse, len(blogs))
	for i, blog := range blogs {
		responses[i] = FromEntity(blog)
	}
	return responses
}
