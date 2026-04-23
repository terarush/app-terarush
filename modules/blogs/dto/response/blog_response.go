package response

import (
	"time"

	"go-modular/modules/blogs/domain/entity"
)

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

// FromEntities converts multiple blog entities to responses
func FromEntities(blogs []*entity.Blog) []*BlogResponse {
	responses := make([]*BlogResponse, len(blogs))
	for i, blog := range blogs {
		responses[i] = FromEntity(blog)
	}
	return responses
}
