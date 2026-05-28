package handler

import (
	"fmt"
	"strconv"
	"time"

	"go-modular/internal/pkg/bus"
	"go-modular/internal/pkg/logger"
	"go-modular/internal/pkg/middleware"
	"go-modular/internal/pkg/utils"
	"go-modular/modules/blogs/domain/entity"
	"go-modular/modules/blogs/domain/service"
	"go-modular/modules/blogs/dto/request"
	"go-modular/modules/blogs/dto/response"

	"github.com/labstack/echo/v4"
)

// BlogHandler handles HTTP requests for blogs
type BlogHandler struct {
	blogService *service.BlogService
	log         *logger.Logger
	event       *bus.EventBus
	r           *utils.Response
}

// NewBlogHandler creates a new blog handler
func NewBlogHandler(log *logger.Logger, event *bus.EventBus, blogService *service.BlogService) *BlogHandler {
	return &BlogHandler{
		blogService: blogService,
		log:         log,
		event:       event,
		r:           &utils.Response{},
	}
}

// getUserIDAndNameFromContext extracts user ID and name from JWT claims in context
func getUserIDAndNameFromContext(c echo.Context) (uint, string, error) {
	claims, ok := c.Get("user").(map[string]interface{})
	if !ok {
		return 0, "", fmt.Errorf("user not found in context")
	}

	userIDFloat, ok := claims["user_id"].(float64)
	if !ok {
		return 0, "", fmt.Errorf("invalid user ID in token")
	}

	userName, ok := claims["name"].(string)
	if !ok {
		return 0, "", fmt.Errorf("invalid user name in token")
	}

	return uint(userIDFloat), userName, nil
}

// GetAllBlogs gets all blogs (admin only) with pagination support
func (h *BlogHandler) GetAllBlogs(c echo.Context) error {
	ctx := c.Request().Context()

	// Get pagination parameters
	page := 1
	pageSize := 10

	if p := c.QueryParam("page"); p != "" {
		if parsedPage, err := strconv.Atoi(p); err == nil && parsedPage > 0 {
			page = parsedPage
		}
	}

	if ps := c.QueryParam("page_size"); ps != "" {
		if parsedSize, err := strconv.Atoi(ps); err == nil && parsedSize > 0 && parsedSize <= 100 {
			pageSize = parsedSize
		}
	}

	// optional search parameter for admin listing
	search := c.QueryParam("search")

	blogs, total, err := h.blogService.GetAllBlogsWithPaginationAndSearch(ctx, page, pageSize, search)
	if err != nil {
		return h.r.InternalServerErrorResponse(c, err.Error())
	}

	return h.r.SuccessResponse(c, map[string]interface{}{
		"blogs":     response.FromEntities(blogs),
		"total":     total,
		"page":      page,
		"page_size": pageSize,
	}, "")
}

// GetPublishedBlogs gets all published blogs (public) with pagination support
func (h *BlogHandler) GetPublishedBlogs(c echo.Context) error {
	ctx := c.Request().Context()

	// Get pagination parameters
	page := 1
	pageSize := 10

	if p := c.QueryParam("page"); p != "" {
		if parsedPage, err := strconv.Atoi(p); err == nil && parsedPage > 0 {
			page = parsedPage
		}
	}

	if ps := c.QueryParam("page_size"); ps != "" {
		if parsedSize, err := strconv.Atoi(ps); err == nil && parsedSize > 0 {
			pageSize = parsedSize
		}
	}

	// optional public search param
	search := c.QueryParam("search")

	blogs, total, err := h.blogService.GetPublishedBlogsWithPaginationAndSearch(ctx, page, pageSize, search)
	if err != nil {
		return h.r.InternalServerErrorResponse(c, err.Error())
	}

	return h.r.SuccessResponse(c, map[string]interface{}{
		"blogs":     response.FromEntities(blogs),
		"total":     total,
		"page":      page,
		"page_size": pageSize,
	}, "")
}

// GetBlogByID gets a blog by ID (admin only)
func (h *BlogHandler) GetBlogByID(c echo.Context) error {
	ctx := c.Request().Context()

	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		return h.r.BadRequestResponse(c, "Invalid blog ID")
	}

	blog, err := h.blogService.GetBlogByID(ctx, uint(id))
	if err != nil {
		if err == service.ErrBlogNotFound {
			return h.r.NotFoundResponse(c, "Blog not found")
		}
		return h.r.InternalServerErrorResponse(c, err.Error())
	}

	return h.r.SuccessResponse(c, response.FromEntity(blog), "")
}

// GetBlogBySlug gets a blog by slug (public)
func (h *BlogHandler) GetBlogBySlug(c echo.Context) error {
	ctx := c.Request().Context()

	slug := c.Param("slug")
	if slug == "" {
		return h.r.BadRequestResponse(c, "Slug is required")
	}

	blog, user, err := h.blogService.GetBlogBySlugWithUser(ctx, slug)
	if err != nil {
		if err == service.ErrBlogNotFound {
			return h.r.NotFoundResponse(c, "Blog not found")
		}
		return h.r.InternalServerErrorResponse(c, err.Error())
	}

	return h.r.SuccessResponse(c, response.FromEntityWithUser(blog, user), "")
}

// CreateBlog creates a new blog (admin only)
func (h *BlogHandler) CreateBlog(c echo.Context) error {
	ctx := c.Request().Context()

	// Get user ID and name from auth context
	userID, userName, err := getUserIDAndNameFromContext(c)
	if err != nil {
		return h.r.UnauthorizedResponse(c, err.Error())
	}

	req := new(request.CreateBlogRequest)
	if err := c.Bind(req); err != nil {
		return h.r.BadRequestResponse(c, err.Error())
	}

	if err := c.Validate(req); err != nil {
		return h.r.ValidationErrorResponse(c, err.Error())
	}

	blog := entity.NewBlog(
		req.Title,
		req.Slug,
		req.Content,
		req.Excerpt,
		userName, // Set author from authenticated user
		userID,   // Set user ID from authenticated user
		req.Category,
		req.Tags,
		req.Image,
	)

	blog.IsPublished = req.IsPublished
	if req.IsPublished {
		now := time.Now()
		blog.PublishedAt = &now
	}

	err = h.blogService.CreateBlog(ctx, blog)
	if err != nil {
		return h.r.InternalServerErrorResponse(c, err.Error())
	}

	h.event.Publish(bus.Event{Type: "blog.created", Payload: blog})

	return h.r.CreatedResponse(c, response.FromEntity(blog), "")
}

// UpdateBlog updates a blog (admin only)
func (h *BlogHandler) UpdateBlog(c echo.Context) error {
	ctx := c.Request().Context()

	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		return h.r.BadRequestResponse(c, "Invalid blog ID")
	}

	req := new(request.UpdateBlogRequest)
	if err := c.Bind(req); err != nil {
		return h.r.BadRequestResponse(c, err.Error())
	}

	if err := c.Validate(req); err != nil {
		return h.r.ValidationErrorResponse(c, err.Error())
	}

	blog, err := h.blogService.GetBlogByID(ctx, uint(id))
	if err != nil {
		if err == service.ErrBlogNotFound {
			return h.r.NotFoundResponse(c, "Blog not found")
		}
		return h.r.InternalServerErrorResponse(c, err.Error())
	}

	blog.Title = req.Title
	blog.Slug = req.Slug
	blog.Content = req.Content
	blog.Excerpt = req.Excerpt
	// Do NOT update Author - it's tied to the original creator via UserID
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
		return h.r.InternalServerErrorResponse(c, err.Error())
	}

	h.event.Publish(bus.Event{Type: "blog.updated", Payload: blog})

	return h.r.SuccessResponse(c, response.FromEntity(blog), "")
}

// DeleteBlog deletes a blog (admin only)
func (h *BlogHandler) DeleteBlog(c echo.Context) error {
	ctx := c.Request().Context()

	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		return h.r.BadRequestResponse(c, "Invalid blog ID")
	}

	err = h.blogService.DeleteBlog(ctx, uint(id))
	if err != nil {
		if err == service.ErrBlogNotFound {
			return h.r.NotFoundResponse(c, "Blog not found")
		}
		return h.r.InternalServerErrorResponse(c, err.Error())
	}

	h.event.Publish(bus.Event{Type: "blog.deleted", Payload: map[string]uint{"id": uint(id)}})

	return h.r.NoContentResponse(c)
}

// RegisterRoutes registers the blog routes
func (h *BlogHandler) RegisterRoutes(e *echo.Echo, basePath string) {
	// Public routes
	public := e.Group(basePath + "/blogs")
	public.GET("", h.GetPublishedBlogs)
	public.GET("/:slug", h.GetBlogBySlug)

	// Public image route
	images := e.Group("/images")
	images.GET("/:filename", h.GetBlogImage)

	// Admin routes
	admin := e.Group(basePath+"/admin/blogs", middleware.Auth, middleware.AdminOnly)
	admin.GET("", h.GetAllBlogs)
	admin.GET("/:id", h.GetBlogByID)
	admin.POST("", h.CreateBlog)
	admin.PUT("/:id", h.UpdateBlog)
	admin.DELETE("/:id", h.DeleteBlog)
	admin.POST("/upload/image", h.UploadBlogImage)
}
