package request

// CreateBlogRequest is the request to create a blog
type CreateBlogRequest struct {
	Title       string `json:"title" validate:"required,max=255"`
	Slug        string `json:"slug" validate:"required,max=255"`
	Content     string `json:"content" validate:"required"`
	Excerpt     string `json:"excerpt" validate:"max=500"`
	Category    string `json:"category" validate:"max=100"`
	Tags        string `json:"tags" validate:"max=255"`
	Image       string `json:"image" validate:"max=255"`
	IsPublished bool   `json:"is_published"`
}

// UpdateBlogRequest is the request to update a blog
type UpdateBlogRequest struct {
	Title       string `json:"title" validate:"required,max=255"`
	Slug        string `json:"slug" validate:"required,max=255"`
	Content     string `json:"content" validate:"required"`
	Excerpt     string `json:"excerpt" validate:"max=500"`
	Category    string `json:"category" validate:"max=100"`
	Tags        string `json:"tags" validate:"max=255"`
	Image       string `json:"image" validate:"max=255"`
	IsPublished bool   `json:"is_published"`
}
