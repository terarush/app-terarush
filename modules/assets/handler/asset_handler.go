package handler

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"go-modular/internal/pkg/logger"
	"go-modular/modules/assets/domain/entity"
	"go-modular/modules/assets/domain/service"
	"go-modular/modules/assets/dto/response"

	"github.com/labstack/echo/v4"
)

// AssetHandler handles HTTP requests for assets
type AssetHandler struct {
	assetService *service.AssetService
	log          *logger.Logger
}

// NewAssetHandler creates a new asset handler
func NewAssetHandler(log *logger.Logger, assetService *service.AssetService) *AssetHandler {
	return &AssetHandler{
		assetService: assetService,
		log:          log,
	}
}

// UploadAsset uploads a file and returns asset URL
// @Summary Upload an asset file
// @Description Upload a file (image, document, or video) and get a URL. Maximum file size: 10MB. Allowed types: images, PDFs, documents, and videos (admin only)
// @Tags Assets
// @Accept multipart/form-data
// @Produce json
// @Param file formData file true "File to upload"
// @Success 201 {object} response.UploadAssetResponse "File uploaded successfully"
// @Failure 400 {object} map[string]string "Invalid file or file too large"
// @Failure 500 {object} map[string]string "Failed to upload file"
// @Router /api/v1/admin/assets/upload [post]
// @Security Bearer
func (h *AssetHandler) UploadAsset(c echo.Context) error {
	// Get file from multipart form
	file, err := c.FormFile("file")
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "No file provided"})
	}

	// Validate file size (max 10MB)
	maxSize := int64(10 * 1024 * 1024)
	if file.Size > maxSize {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "File size must be less than 10MB"})
	}

	// Open file
	src, err := file.Open()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to read file"})
	}
	defer src.Close()

	// Read file header to validate MIME type
	header := make([]byte, 512)
	if _, err := src.Read(header); err != nil && err != io.EOF {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to read file header"})
	}

	mimeType := http.DetectContentType(header)

	// Allowed MIME types: images, documents, videos, etc.
	allowedTypes := map[string]bool{
		"image/jpeg":  true,
		"image/png":   true,
		"image/gif":   true,
		"image/webp":  true,
		"image/svg+xml": true,
		"application/pdf": true,
		"text/plain":  true,
		"application/msword": true,
		"application/vnd.openxmlformats-officedocument.wordprocessingml.document": true,
		"application/vnd.ms-excel": true,
		"application/vnd.openxmlformats-officedocument.spreadsheetml.sheet": true,
		"video/mp4":   true,
		"video/webm":  true,
	}

	if !allowedTypes[mimeType] {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "File type not allowed. Allowed types: images, PDFs, documents, and videos",
		})
	}

	// Reset file pointer
	src.Seek(0, 0)

	// Create unique filename
	timestamp := strconv.FormatInt(time.Now().UnixNano(), 10)
	ext := filepath.Ext(file.Filename)
	filename := fmt.Sprintf("%s%s", timestamp, ext)

	// Create upload directory if not exists
	uploadDir := "public/uploads/assets"
	if err := os.MkdirAll(uploadDir, 0755); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to create upload directory"})
	}

	// Save file
	filePath := filepath.Join(uploadDir, filename)
	dst, err := os.Create(filePath)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to save file"})
	}
	defer dst.Close()

	if _, err := io.Copy(dst, src); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to save file"})
	}

	// Generate asset URL
	assetURL := fmt.Sprintf("/public/uploads/assets/%s", filename)

	// Create asset record in database
	asset := entity.NewAsset(assetURL)
	asset.FileName = file.Filename
	asset.FileSize = file.Size
	asset.MimeType = mimeType
	ctx := c.Request().Context()
	if err := h.assetService.CreateAsset(ctx, asset); err != nil {
		h.log.Error("Error creating asset record", err)
		// Don't fail the upload, just log the error
	}

	return c.JSON(http.StatusCreated, response.UploadAssetResponse{
		ID:   asset.ID,
		URL:  assetURL,
		Path: filePath,
	})
}

// GetAssets retrieves all assets with pagination
// @Summary Get all assets
// @Description Get all assets with pagination (admin only)
// @Tags Assets
// @Produce json
// @Param page query int false "Page number (default: 1)"
// @Param page_size query int false "Page size (default: 20)"
// @Success 200 {object} response.ListAssetsResponse "Assets retrieved successfully"
// @Failure 500 {object} map[string]string "Failed to retrieve assets"
// @Router /api/v1/admin/assets [get]
// @Security Bearer
func (h *AssetHandler) GetAssets(c echo.Context) error {
	page := 1
	pageSize := 20

	// Get query parameters
	if p := c.QueryParam("page"); p != "" {
		if parsed, err := strconv.Atoi(p); err == nil && parsed > 0 {
			page = parsed
		}
	}
	if ps := c.QueryParam("page_size"); ps != "" {
		if parsed, err := strconv.Atoi(ps); err == nil && parsed > 0 && parsed <= 100 {
			pageSize = parsed
		}
	}

	ctx := c.Request().Context()
	assets, total, err := h.assetService.GetAllAssets(ctx, page, pageSize)
	if err != nil {
		h.log.Error("Error retrieving assets", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to retrieve assets"})
	}

	// Convert to response objects
	assetResponses := make([]response.AssetResponse, len(assets))
	for i, asset := range assets {
		assetResponses[i] = response.FromEntity(asset)
	}

	return c.JSON(http.StatusOK, response.ListAssetsResponse{
		Assets:   assetResponses,
		Total:    total,
		Page:     page,
		PageSize: pageSize,
	})
}

// DeleteAsset deletes an asset
// @Summary Delete an asset
// @Description Delete an asset by ID (admin only)
// @Tags Assets
// @Produce json
// @Param id path int true "Asset ID"
// @Success 204 "Asset deleted successfully"
// @Failure 400 {object} map[string]string "Invalid asset ID"
// @Failure 404 {object} map[string]string "Asset not found"
// @Failure 500 {object} map[string]string "Failed to delete asset"
// @Router /api/v1/admin/assets/{id} [delete]
// @Security Bearer
func (h *AssetHandler) DeleteAsset(c echo.Context) error {
	ctx := c.Request().Context()
	id := c.Param("id")

	assetID := uint(0)
	if _, err := fmt.Sscanf(id, "%d", &assetID); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid asset ID"})
	}

	// Get asset to find file path
	asset, err := h.assetService.GetAssetByID(ctx, assetID)
	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "Asset not found"})
	}

	// Delete file from disk
	if err := os.Remove(asset.URL); err != nil && !os.IsNotExist(err) {
		h.log.Error("Error deleting asset file", err)
		// Don't fail, continue with database deletion
	}

	// Delete asset record from database
	if err := h.assetService.DeleteAsset(ctx, assetID); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to delete asset"})
	}

	return c.NoContent(http.StatusNoContent)
}

// GetAsset serves an asset file
// @Summary Get asset file
// @Description Retrieve an asset file by filename (public endpoint)
// @Tags Assets
// @Produce image/jpeg,image/png,image/gif,image/webp,image/svg+xml,application/pdf,video/mp4,video/webm
// @Param filename path string true "Asset filename"
// @Success 200 "Asset file"
// @Failure 400 {object} map[string]string "Invalid filename"
// @Failure 404 {object} map[string]string "Asset not found"
// @Failure 500 {object} map[string]string "Failed to retrieve asset"
// @Router /api/v1/assets/{filename} [get]
func (h *AssetHandler) GetAsset(c echo.Context) error {
	filename := c.Param("filename")
	if filename == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Filename is required"})
	}

	// Security: prevent directory traversal attacks
	if filepath.Base(filename) != filename {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid filename"})
	}

	filePath := filepath.Join("public/uploads/assets", filename)

	// Check if file exists
	if _, err := os.Stat(filePath); err != nil {
		if os.IsNotExist(err) {
			return c.JSON(http.StatusNotFound, map[string]string{"error": "Asset not found"})
		}
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to read asset"})
	}

	// Set cache headers
	c.Response().Header().Set("Cache-Control", "public, max-age=2592000") // 30 days

	// Serve the file
	return c.File(filePath)
}

// RegisterRoutes registers the asset routes
func (h *AssetHandler) RegisterRoutes(e *echo.Echo, basePath string) {
	// Public routes
	public := e.Group(basePath + "/assets")
	public.GET("/:filename", h.GetAsset)

	// Admin routes - upload, delete, and list
	admin := e.Group(basePath + "/admin/assets")
	admin.GET("", h.GetAssets)
	admin.POST("/upload", h.UploadAsset)
	admin.DELETE("/:id", h.DeleteAsset)
}
