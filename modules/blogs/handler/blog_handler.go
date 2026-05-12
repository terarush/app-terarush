package handler

import (
	"net/http"
	"strconv"
	"time"

	"go-modular/internal/pkg/bus"
	"go-modular/internal/pkg/logger"
	"go-modular/internal/pkg/middleware"
	"go-modular/modules/blogs/domain/entity"
	"go-modular/modules/blogs/domain/service"
	"go-modular/modules/blogs/dto/request"
	"go-modular/modules/blogs/dto/response"

	"github.com/labstack/echo"
)

// BlogHandler handles HTTP requests for blogs
type BlogHandler struct {
	blogService *service.BlogService
	log         *logger.Logger
	event       *bus.EventBus
}

// NewBlogHandler creates a new blog handler
func NewBlogHandler(log *logger.Logger, event *bus.EventBus, blogService *service.BlogService) *BlogHandler {
	return &BlogHandler{
		blogService: blogService,
		log:         log,
		event:       event,
	}
}

// GetAllBlogs gets all blogs (admin only)
func (h *BlogHandler) GetAllBlogs(c echo.Context) error {
	ctx := c.Request().Context()

	blogs, err := h.blogService.GetAllBlogs(ctx)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, response.FromEntities(blogs))
}

// GetPublishedBlogs gets all published blogs (public)
func (h *BlogHandler) GetPublishedBlogs(c echo.Context) error {
	ctx := c.Request().Context()

	blogs, err := h.blogService.GetPublishedBlogs(ctx)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, response.FromEntities(blogs))
}

// GetBlogByID gets a blog by ID (admin only)
func (h *BlogHandler) GetBlogByID(c echo.Context) error {
	ctx := c.Request().Context()

	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid blog ID"})
	}

	blog, err := h.blogService.GetBlogByID(ctx, uint(id))
	if err != nil {
		if err == service.ErrBlogNotFound {
			return c.JSON(http.StatusNotFound, map[string]string{"error": "Blog not found"})
		}
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, response.FromEntity(blog))
}

// GetBlogBySlug gets a blog by slug (public)
func (h *BlogHandler) GetBlogBySlug(c echo.Context) error {
	ctx := c.Request().Context()

	slug := c.Param("slug")
	if slug == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Slug is required"})
	}

	blog, err := h.blogService.GetBlogBySlug(ctx, slug)
	if err != nil {
		if err == service.ErrBlogNotFound {
			return c.JSON(http.StatusNotFound, map[string]string{"error": "Blog not found"})
		}
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, response.FromEntity(blog))
}

// CreateBlog creates a new blog (admin only)
func (h *BlogHandler) CreateBlog(c echo.Context) error {
	ctx := c.Request().Context()

	req := new(request.CreateBlogRequest)
	if err := c.Bind(req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	if err := c.Validate(req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	blog := &entity.Blog{
		Title:       req.Title,
		Slug:        req.Slug,
		Content:     req.Content,
		Excerpt:     req.Excerpt,
		Author:      req.Author,
		Category:    req.Category,
		Tags:        req.Tags,
		Image:       req.Image,
		IsPublished: req.IsPublished,
	}

	if req.IsPublished {
		now := time.Now()
		blog.PublishedAt = &now
	}

	err := h.blogService.CreateBlog(ctx, blog)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	h.event.Publish(bus.Event{Type: "blog.created", Payload: blog})

	return c.JSON(http.StatusCreated, response.FromEntity(blog))
}

// UpdateBlog updates a blog (admin only)
func (h *BlogHandler) UpdateBlog(c echo.Context) error {
	ctx := c.Request().Context()

	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid blog ID"})
	}

	req := new(request.UpdateBlogRequest)
	if err := c.Bind(req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	if err := c.Validate(req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	blog, err := h.blogService.GetBlogByID(ctx, uint(id))
	if err != nil {
		if err == service.ErrBlogNotFound {
			return c.JSON(http.StatusNotFound, map[string]string{"error": "Blog not found"})
		}
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	blog.Title = req.Title
	blog.Slug = req.Slug
	blog.Content = req.Content
	blog.Excerpt = req.Excerpt
	blog.Author = req.Author
	blog.Category = req.Category
	blog.Tags = req.Tags
	blog.Image = req.Image
	blog.IsPublished = req.IsPublished

	if req.IsPublished && blog.PublishedAt == nil {
		now := time.Now()
		blog.PublishedAt = &now
	}

	err = h.blogService.UpdateBlog(ctx, blog)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	h.event.Publish(bus.Event{Type: "blog.updated", Payload: blog})

	return c.JSON(http.StatusOK, response.FromEntity(blog))
}

// DeleteBlog deletes a blog (admin only)
func (h *BlogHandler) DeleteBlog(c echo.Context) error {
	ctx := c.Request().Context()

	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid blog ID"})
	}

	err = h.blogService.DeleteBlog(ctx, uint(id))
	if err != nil {
		if err == service.ErrBlogNotFound {
			return c.JSON(http.StatusNotFound, map[string]string{"error": "Blog not found"})
		}
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	h.event.Publish(bus.Event{Type: "blog.deleted", Payload: map[string]uint{"id": uint(id)}})

	return c.NoContent(http.StatusNoContent)
}

// RegisterRoutes registers the blog routes
func (h *BlogHandler) RegisterRoutes(e *echo.Echo, basePath string) {
	// Public routes
	public := e.Group(basePath + "/blogs")
	public.GET("", h.GetPublishedBlogs)
	public.GET("/:slug", h.GetBlogBySlug)

	// Admin routes
	admin := e.Group(basePath+"/admin/blogs", middleware.Auth, middleware.AdminOnly)
	admin.GET("", h.GetAllBlogs)
	admin.GET("/:id", h.GetBlogByID)
	admin.POST("", h.CreateBlog)
	admin.PUT("/:id", h.UpdateBlog)
	admin.DELETE("/:id", h.DeleteBlog)
	admin.POST("/upload/image", h.UploadBlogImage)
}
