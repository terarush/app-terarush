package entity

import "time"

// Blog represents a blog post
type Blog struct {
	ID          uint       `gorm:"primaryKey" json:"id"`
	Title       string     `gorm:"type:varchar(255);not null" json:"title"`
	Slug        string     `gorm:"type:varchar(255);uniqueIndex;not null" json:"slug"`
	Content     string     `gorm:"type:text;not null" json:"content"`
	Excerpt     string     `gorm:"type:text" json:"excerpt"`
	Author      string     `gorm:"type:varchar(100);not null" json:"author"`
	UserID      uint       `gorm:"index" json:"user_id"` // Foreign key to User
	Category    string     `gorm:"type:varchar(100)" json:"category"`
	Tags        string     `gorm:"type:varchar(255)" json:"tags"` // Comma-separated
	Image       string     `gorm:"type:varchar(255)" json:"image"`
	IsPublished bool       `gorm:"type:boolean;default:false" json:"is_published"`
	ViewCount   int64      `gorm:"type:integer;default:0" json:"view_count"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
	PublishedAt *time.Time `json:"published_at"`
}

// TableName specifies the table name
func (Blog) TableName() string {
	return "blogs"
}

func NewBlog(title, slug, content, excerpt, author string, userID uint, category, tags, image string) *Blog {
	now := time.Now()
	return &Blog{
		Title:       title,
		Slug:        slug,
		Content:     content,
		Excerpt:     excerpt,
		Author:      author,
		UserID:      userID,
		Category:    category,
		Tags:        tags,
		Image:       image,
		IsPublished: false,
		ViewCount:   0,
		CreatedAt:   now,
		UpdatedAt:   now,
	}
}
