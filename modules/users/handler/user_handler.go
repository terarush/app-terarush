// internal/modules/user/interfaces/handler/user_handler.go

package handler

import (
	"fmt"
	"go-modular/internal/pkg/bus"
	"go-modular/internal/pkg/logger"
	"go-modular/internal/pkg/middleware"
	"go-modular/internal/pkg/utils"
	"go-modular/modules/users/domain/entity"
	"go-modular/modules/users/domain/service"
	"go-modular/modules/users/dto/request"
	"go-modular/modules/users/dto/response"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/labstack/echo/v4"
)

// UserHandler handles HTTP requests for users
type UserHandler struct {
	userService *service.UserService
	log         *logger.Logger
	event       *bus.EventBus
	r           *utils.Response
}

// NewUserHandler creates a new user handler
func NewUserHandler(log *logger.Logger, event *bus.EventBus, userService *service.UserService) *UserHandler {
	return &UserHandler{
		userService: userService,
		log:         log,
		event:       event,
		r:           &utils.Response{},
	}
}

// Event Bus Event user created
func (h *UserHandler) Handle(event bus.Event) {
	fmt.Printf("User created: %v", event.Payload)
}

// getUserIDFromContext extracts user ID from JWT claims in context
func getUserIDFromContext(c echo.Context) (uint, error) {
	claims, ok := c.Get("user").(map[string]interface{})
	if !ok {
		return 0, fmt.Errorf("user not found in context")
	}

	userIDFloat, ok := claims["user_id"].(float64)
	if !ok {
		return 0, fmt.Errorf("invalid user ID in token")
	}

	return uint(userIDFloat), nil
}

// GetAllUsers gets all users
func (h *UserHandler) GetAllUsers(c echo.Context) error {
	ctx := c.Request().Context()

	users, err := h.userService.GetAllUsers(ctx)
	if err != nil {
		return h.r.InternalServerErrorResponse(c, err.Error())
	}

	return h.r.SuccessResponse(c, response.FromEntities(users), "")
}

// GetUser gets a user by ID
func (h *UserHandler) GetUser(c echo.Context) error {
	ctx := c.Request().Context()

	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		return h.r.BadRequestResponse(c, "Invalid user ID")
	}

	user, err := h.userService.GetUserByID(ctx, uint(id))
	if err != nil {
		if err == service.ErrUserNotFound {
			return h.r.NotFoundResponse(c, "User not found")
		}
		return h.r.InternalServerErrorResponse(c, err.Error())
	}

	return h.r.SuccessResponse(c, response.FromEntity(user), "")
}

// CreateUser creates a new user
func (h *UserHandler) CreateUser(c echo.Context) error {
	ctx := c.Request().Context()

	req := new(request.CreateUserRequest)
	if err := c.Bind(req); err != nil {
		return h.r.BadRequestResponse(c, err.Error())
	}

	if err := c.Validate(req); err != nil {
		return h.r.ValidationErrorResponse(c, err.Error())
	}

	user := entity.NewUser(req.Name, req.Email, req.Password)
	err := h.userService.CreateUser(ctx, user)
	if err != nil {
		if err == service.ErrEmailAlreadyUsed {
			return h.r.ConflictResponse(c, "Email already in use")
		}
		return h.r.InternalServerErrorResponse(c, err.Error())
	}

	// event bus publish
	h.event.Publish(bus.Event{Type: "user.created", Payload: user})

	return h.r.CreatedResponse(c, response.FromEntity(user), "")
}

// UpdateUser updates a user
func (h *UserHandler) UpdateUser(c echo.Context) error {
	ctx := c.Request().Context()

	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		return h.r.BadRequestResponse(c, "Invalid user ID")
	}

	req := new(request.UpdateUserRequest)
	if err := c.Bind(req); err != nil {
		return h.r.BadRequestResponse(c, err.Error())
	}

	if err := c.Validate(req); err != nil {
		return h.r.ValidationErrorResponse(c, err.Error())
	}

	user, err := h.userService.GetUserByID(ctx, uint(id))
	if err != nil {
		if err == service.ErrUserNotFound {
			return h.r.NotFoundResponse(c, "User not found")
		}
		return h.r.InternalServerErrorResponse(c, err.Error())
	}

	user.Name = req.Name
	user.Email = req.Email
	if req.Role != "" {
		user.Role = req.Role
	}

	err = h.userService.UpdateUser(ctx, user)
	if err != nil {
		return h.r.InternalServerErrorResponse(c, err.Error())
	}

	return h.r.SuccessResponse(c, response.FromEntity(user), "")
}

// DeleteUser deletes a user
func (h *UserHandler) DeleteUser(c echo.Context) error {
	ctx := c.Request().Context()

	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		return h.r.BadRequestResponse(c, "Invalid user ID")
	}

	err = h.userService.DeleteUser(ctx, uint(id))
	if err != nil {
		if err == service.ErrUserNotFound {
			return h.r.NotFoundResponse(c, "User not found")
		}
		return h.r.InternalServerErrorResponse(c, err.Error())
	}

	return h.r.NoContentResponse(c)
}

// GetCurrentUser gets the current authenticated user
func (h *UserHandler) GetCurrentUser(c echo.Context) error {
	ctx := c.Request().Context()

	// Get user ID from context
	userID, err := getUserIDFromContext(c)
	if err != nil {
		return h.r.UnauthorizedResponse(c, err.Error())
	}

	user, err := h.userService.GetUserByID(ctx, userID)
	if err != nil {
		if err == service.ErrUserNotFound {
			return h.r.NotFoundResponse(c, "User not found")
		}
		return h.r.InternalServerErrorResponse(c, err.Error())
	}

	return h.r.SuccessResponse(c, response.FromEntity(user), "")
}

// UpdateProfile updates the current user's profile
func (h *UserHandler) UpdateProfile(c echo.Context) error {
	ctx := c.Request().Context()

	// Get user ID from context
	userID, err := getUserIDFromContext(c)
	if err != nil {
		return h.r.UnauthorizedResponse(c, err.Error())
	}

	req := new(request.UpdateProfileRequest)
	if err := c.Bind(req); err != nil {
		return h.r.BadRequestResponse(c, err.Error())
	}

	if err := c.Validate(req); err != nil {
		return h.r.ValidationErrorResponse(c, err.Error())
	}

	user, err := h.userService.GetUserByID(ctx, userID)
	if err != nil {
		if err == service.ErrUserNotFound {
			return h.r.NotFoundResponse(c, "User not found")
		}
		return h.r.InternalServerErrorResponse(c, err.Error())
	}

	user.Name = req.Name
	if req.Avatar != "" {
		user.Avatar = req.Avatar
	}
	if req.Bio != "" {
		user.Bio = req.Bio
	}
	if req.Banner != "" {
		user.Banner = req.Banner
	}

	err = h.userService.UpdateUser(ctx, user)
	if err != nil {
		return h.r.InternalServerErrorResponse(c, err.Error())
	}

	return h.r.SuccessResponse(c, response.FromEntity(user), "")
}

// UploadBanner handles banner file upload
func (h *UserHandler) UploadBanner(c echo.Context) error {
	// Get user ID from context
	userID, err := getUserIDFromContext(c)
	if err != nil {
		return h.r.UnauthorizedResponse(c, err.Error())
	}

	// Get file from form
	file, err := c.FormFile("banner")
	if err != nil {
		return h.r.BadRequestResponse(c, "No file uploaded")
	}

	// Validate file type
	if !strings.HasPrefix(file.Header.Get("Content-Type"), "image/") {
		return h.r.BadRequestResponse(c, "File must be an image")
	}

	// Validate file size (max 10MB)
	if file.Size > 10*1024*1024 {
		return h.r.BadRequestResponse(c, "File size must be less than 10MB")
	}

	// Open uploaded file
	src, err := file.Open()
	if err != nil {
		return h.r.InternalServerErrorResponse(c, "Failed to open file")
	}
	defer src.Close()

	// Create uploads directory if not exists
	uploadDir := "./public/uploads/banners"
	if err := os.MkdirAll(uploadDir, 0755); err != nil {
		return h.r.InternalServerErrorResponse(c, "Failed to create upload directory")
	}

	// Generate unique filename
	ext := filepath.Ext(file.Filename)
	filename := fmt.Sprintf("%d_%d%s", userID, time.Now().Unix(), ext)
	filePath := filepath.Join(uploadDir, filename)

	// Create destination file
	dst, err := os.Create(filePath)
	if err != nil {
		return h.r.InternalServerErrorResponse(c, "Failed to create file")
	}
	defer dst.Close()

	// Copy file
	if _, err := io.Copy(dst, src); err != nil {
		return h.r.InternalServerErrorResponse(c, "Failed to save file")
	}

	// Return the URL path
	bannerURL := fmt.Sprintf("/public/uploads/banners/%s", filename)

	return h.r.SuccessResponse(c, map[string]string{
		"banner_url": bannerURL,
		"message":    "Banner uploaded successfully",
	}, "")
}

// UploadAvatar handles avatar file upload
func (h *UserHandler) UploadAvatar(c echo.Context) error {
	// Get user ID from context
	userID, err := getUserIDFromContext(c)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": err.Error()})
	}

	// Get file from form
	file, err := c.FormFile("avatar")
	if err != nil {
		return h.r.BadRequestResponse(c, "No file uploaded")
	}

	// Validate file type
	if !strings.HasPrefix(file.Header.Get("Content-Type"), "image/") {
		return h.r.BadRequestResponse(c, "File must be an image")
	}

	// Validate file size (max 5MB)
	if file.Size > 5*1024*1024 {
		return h.r.BadRequestResponse(c, "File size must be less than 5MB")
	}

	// Open uploaded file
	src, err := file.Open()
	if err != nil {
		return h.r.InternalServerErrorResponse(c, "Failed to open file")
	}
	defer src.Close()

	// Create uploads directory if not exists
	uploadDir := "./public/uploads/avatars"
	if err := os.MkdirAll(uploadDir, 0755); err != nil {
		return h.r.InternalServerErrorResponse(c, "Failed to create upload directory")
	}

	// Generate unique filename
	ext := filepath.Ext(file.Filename)
	filename := fmt.Sprintf("%d_%d%s", userID, time.Now().Unix(), ext)
	filePath := filepath.Join(uploadDir, filename)

	// Create destination file
	dst, err := os.Create(filePath)
	if err != nil {
		return h.r.InternalServerErrorResponse(c, "Failed to create file")
	}
	defer dst.Close()

	// Copy file
	if _, err := io.Copy(dst, src); err != nil {
		return h.r.InternalServerErrorResponse(c, "Failed to save file")
	}

	// Return the URL path
	avatarURL := fmt.Sprintf("/public/uploads/avatars/%s", filename)

	return h.r.SuccessResponse(c, map[string]string{
		"avatar_url": avatarURL,
		"message":    "Avatar uploaded successfully",
	}, "")
}

// RegisterRoutes registers the user routes
func (h *UserHandler) RegisterRoutes(e *echo.Echo, basePath string) {
	group := e.Group(basePath + "/users")

	// Admin only routes - full CRUD access
	admin := group.Group("", middleware.Auth, middleware.AdminOnly)
	admin.GET("", h.GetAllUsers)
	admin.POST("", h.CreateUser)
	admin.PUT("/:id", h.UpdateUser)
	admin.DELETE("/:id", h.DeleteUser)

	// User or Admin can view their own profile or admin can view any
	userOrAdmin := group.Group("", middleware.Auth)
	userOrAdmin.GET("/:id", h.GetUser)

	// Current user profile routes
	profile := e.Group(basePath+"/profile", middleware.Auth)
	profile.GET("", h.GetCurrentUser)
	profile.PUT("", h.UpdateProfile)
	profile.POST("/avatar", h.UploadAvatar)
	profile.POST("/banner", h.UploadBanner)
}
