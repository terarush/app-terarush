package handler

import (
	"encoding/json"
	"go-modular/internal/pkg/logger"
	"go-modular/internal/pkg/utils"
	"go-modular/modules/products/domain/entity"
	"go-modular/modules/products/domain/service"
	"go-modular/modules/products/dto/request"
	"go-modular/modules/products/dto/response"
	"net/http"
	"strconv"

	"github.com/labstack/echo"
)

type ProductHandler struct {
	productService *service.ProductService
	log            *logger.Logger
	r              *utils.Response
}

func NewProductHandler(productService *service.ProductService, log *logger.Logger) *ProductHandler {
	return &ProductHandler{
		productService: productService,
		log:            log,
		r:              &utils.Response{},
	}
}

// GetProducts godoc
// @Summary Get all products
// @Description Get all products (admin sees all, user sees active only)
// @Tags products
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param include_inactive query bool false "Include inactive products (admin only)"
// @Success 200 {object} utils.SuccessResponse{data=[]response.ProductResponse}
// @Router /products [get]
func (h *ProductHandler) GetProducts(c echo.Context) error {
	// Check if user is admin
	userRole, _ := c.Get("role").(string)
	isAdmin := userRole == "admin"

	includeInactive := false
	if isAdmin {
		includeInactive = c.QueryParam("include_inactive") == "true"
	}

	var products []*entity.Product
	var err error

	if isAdmin {
		products, err = h.productService.GetAllProducts(c.Request().Context(), includeInactive)
	} else {
		products, err = h.productService.GetActiveProducts(c.Request().Context())
	}

	if err != nil {
		h.log.Error("Failed to get products:", err)
		return h.r.ErrorResponse(c, http.StatusInternalServerError, "Failed to get products")
	}

	return h.r.SuccessResponse(c, response.FromEntities(products), "Products retrieved successfully")
}

// GetProductByID godoc
// @Summary Get product by ID
// @Description Get product details by ID
// @Tags products
// @Accept json
// @Produce json
// @Param id path int true "Product ID"
// @Success 200 {object} utils.SuccessResponse{data=response.ProductResponse}
// @Router /products/{id} [get]
func (h *ProductHandler) GetProductByID(c echo.Context) error {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		return h.r.ErrorResponse(c, http.StatusBadRequest, "Invalid product ID")
	}

	product, err := h.productService.GetProductByID(c.Request().Context(), uint(id))
	if err != nil {
		if err == service.ErrProductNotFound {
			return h.r.ErrorResponse(c, http.StatusNotFound, "Product not found")
		}
		h.log.Error("Failed to get product:", err)
		return h.r.ErrorResponse(c, http.StatusInternalServerError, "Failed to get product")
	}

	return h.r.SuccessResponse(c, response.FromEntity(product), "Product retrieved successfully")
}

// CreateProduct godoc
// @Summary Create product (Admin only)
// @Description Create a new product
// @Tags products
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param product body request.CreateProductRequest true "Product data"
// @Success 201 {object} utils.SuccessResponse{data=response.ProductResponse}
// @Router /products [post]
func (h *ProductHandler) CreateProduct(c echo.Context) error {
	var req request.CreateProductRequest
	if err := c.Bind(&req); err != nil {
		return h.r.ErrorResponse(c, http.StatusBadRequest, "Invalid request body")
	}

	if err := c.Validate(&req); err != nil {
		return h.r.ErrorResponse(c, http.StatusBadRequest, err.Error())
	}

	// Marshal ImageTags to JSON string
	imageTagsJSON, err := json.Marshal(req.ImageTags)
	if err != nil {
		return h.r.ErrorResponse(c, http.StatusBadRequest, "Invalid image tags")
	}

	product := &entity.Product{
		Name:        req.Name,
		Description: req.Description,
		Price:       req.Price,
		CPUCores:    req.CPUCores,
		RAMMB:       req.RAMMB,
		StorageGB:   req.StorageGB,
		BandwidthGB: req.BandwidthGB,
		IsActive:    req.IsActive,
		Stock:       req.Stock,
		ImageURL:    req.ImageURL,
		DockerImage: req.DockerImage,
		ImageTags:   string(imageTagsJSON),
	}

	err = h.productService.CreateProduct(c.Request().Context(), product)
	if err != nil {
		h.log.Error("Failed to create product:", err)
		return h.r.ErrorResponse(c, http.StatusInternalServerError, "Failed to create product")
	}

	return h.r.SuccessResponse(c, response.FromEntity(product), "Product created successfully")
}

// UpdateProduct godoc
// @Summary Update product (Admin only)
// @Description Update an existing product
// @Tags products
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Product ID"
// @Param product body request.UpdateProductRequest true "Product data"
// @Success 200 {object} utils.SuccessResponse{data=response.ProductResponse}
// @Router /products/{id} [put]
func (h *ProductHandler) UpdateProduct(c echo.Context) error {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		return h.r.ErrorResponse(c, http.StatusBadRequest, "Invalid product ID")
	}

	var req request.UpdateProductRequest
	if err := c.Bind(&req); err != nil {
		return h.r.ErrorResponse(c, http.StatusBadRequest, "Invalid request body")
	}

	if err := c.Validate(&req); err != nil {
		return h.r.ErrorResponse(c, http.StatusBadRequest, err.Error())
	}

	// Marshal ImageTags to JSON string
	imageTagsJSON, err := json.Marshal(req.ImageTags)
	if err != nil {
		return h.r.ErrorResponse(c, http.StatusBadRequest, "Invalid image tags")
	}

	product := &entity.Product{
		ID:          uint(id),
		Name:        req.Name,
		Description: req.Description,
		Price:       req.Price,
		CPUCores:    req.CPUCores,
		RAMMB:       req.RAMMB,
		StorageGB:   req.StorageGB,
		BandwidthGB: req.BandwidthGB,
		IsActive:    req.IsActive,
		Stock:       req.Stock,
		ImageURL:    req.ImageURL,
		DockerImage: req.DockerImage,
		ImageTags:   string(imageTagsJSON),
	}

	err = h.productService.UpdateProduct(c.Request().Context(), product)
	if err != nil {
		if err == service.ErrProductNotFound {
			return h.r.ErrorResponse(c, http.StatusNotFound, "Product not found")
		}
		h.log.Error("Failed to update product:", err)
		return h.r.ErrorResponse(c, http.StatusInternalServerError, "Failed to update product")
	}

	return h.r.SuccessResponse(c, response.FromEntity(product), "Product updated successfully")
}

// DeleteProduct godoc
// @Summary Delete product (Admin only)
// @Description Delete a product (soft delete)
// @Tags products
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Product ID"
// @Success 200 {object} utils.SuccessResponse
// @Router /products/{id} [delete]
func (h *ProductHandler) DeleteProduct(c echo.Context) error {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		return h.r.ErrorResponse(c, http.StatusBadRequest, "Invalid product ID")
	}

	err = h.productService.DeleteProduct(c.Request().Context(), uint(id))
	if err != nil {
		if err == service.ErrProductNotFound {
			return h.r.ErrorResponse(c, http.StatusNotFound, "Product not found")
		}
		h.log.Error("Failed to delete product:", err)
		return h.r.ErrorResponse(c, http.StatusInternalServerError, "Failed to delete product")
	}

	return h.r.SuccessResponse(c, nil, "Product deleted successfully")
}

// RegisterRoutes registers product routes
func (h *ProductHandler) RegisterRoutes(e *echo.Echo, basePath string, authMiddleware echo.MiddlewareFunc, adminMiddleware echo.MiddlewareFunc) {
	group := e.Group(basePath + "/products")

	// Public routes (no auth required)
	group.GET("", h.GetProducts)
	group.GET("/:id", h.GetProductByID)

	// Admin only routes - apply middleware per route to avoid leaking
	group.POST("", h.CreateProduct, authMiddleware, adminMiddleware)
	group.PUT("/:id", h.UpdateProduct, authMiddleware, adminMiddleware)
	group.DELETE("/:id", h.DeleteProduct, authMiddleware, adminMiddleware)
}
