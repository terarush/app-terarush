package handler

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/labstack/echo"
)

// UploadImageResponse is the response for image upload
type UploadImageResponse struct {
	URL  string `json:"url"`
	Path string `json:"path"`
}

// UploadBlogImage uploads an image for blog post
func (h *BlogHandler) UploadBlogImage(c echo.Context) error {
	// Get file from multipart form
	file, err := c.FormFile("image")
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "No file provided"})
	}

	// Validate file size (max 5MB)
	if file.Size > 5*1024*1024 {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "File size must be less than 5MB"})
	}

	// Validate file type
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
	if mimeType != "image/jpeg" && mimeType != "image/png" && mimeType != "image/gif" && mimeType != "image/webp" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid file type. Only JPEG, PNG, GIF, and WebP are allowed"})
	}

	// Reset file pointer
	src.Seek(0, 0)

	// Create unique filename
	timestamp := strconv.FormatInt(time.Now().UnixNano(), 10)
	ext := filepath.Ext(file.Filename)
	filename := fmt.Sprintf("blog_%s%s", timestamp, ext)

	// Create upload directory if not exists
	uploadDir := "public/uploads/blogs"
	if err := os.MkdirAll(uploadDir, 0755); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to create upload directory"})
	}

	// Save file
	filepath := filepath.Join(uploadDir, filename)
	dst, err := os.Create(filepath)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to save file"})
	}
	defer dst.Close()

	if _, err := io.Copy(dst, src); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to save file"})
	}

	// Return URL with image ID (just the filename)
	url := fmt.Sprintf("/images/%s", filename)

	return c.JSON(http.StatusOK, UploadImageResponse{
		URL:  url,
		Path: filepath,
	})
}

// GetBlogImage serves a blog image
func (h *BlogHandler) GetBlogImage(c echo.Context) error {
	filename := c.Param("filename")
	if filename == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Filename is required"})
	}

	// Security: prevent directory traversal attacks
	if filepath.Base(filename) != filename {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid filename"})
	}

	filepath := filepath.Join("public/uploads/blogs", filename)

	// Check if file exists
	if _, err := os.Stat(filepath); err != nil {
		if os.IsNotExist(err) {
			return c.JSON(http.StatusNotFound, map[string]string{"error": "Image not found"})
		}
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to read image"})
	}

	// Set cache headers
	c.Response().Header().Set("Cache-Control", "public, max-age=604800") // 7 days

	// Serve the file
	return c.File(filepath)
}
