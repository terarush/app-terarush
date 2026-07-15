package handler

import (
	"go-modular/internal/pkg/bus"
	"go-modular/internal/pkg/logger"
	"go-modular/internal/pkg/middleware"
	"go-modular/internal/pkg/utils"
	"go-modular/modules/bookmarks/domain/entity"
	"go-modular/modules/bookmarks/domain/service"
	"go-modular/modules/bookmarks/dto/response"
	"strconv"

	"github.com/labstack/echo/v4"
)

type BookmarkHandler struct {
	svc *service.BookmarkService
	log *logger.Logger
	evt *bus.EventBus
	r   *utils.Response
}

func NewBookmarkHandler(log *logger.Logger, evt *bus.EventBus, svc *service.BookmarkService) *BookmarkHandler {
	return &BookmarkHandler{svc: svc, log: log, evt: evt, r: &utils.Response{}}
}

// List lists all bookmarks
// @Summary List bookmarks
// @Description List all bookmarks with optional status and collection_id filters
// @Tags Bookmarks
// @Produce json
// @Param status query string false "Filter by status"
// @Param collection_id query string false "Filter by collection ID"
// @Success 200 {object} utils.SuccessResponse "Bookmarks list"
// @Failure 500 {object} utils.ErrorResponse "Internal server error"
// @Router /api/v1/bookmarks [get]
// @Security Bearer
func (h *BookmarkHandler) List(c echo.Context) error {
	filter := map[string]any{}
	if s := c.QueryParam("status"); s != "" { filter["status"] = s }
	if s := c.QueryParam("collection_id"); s != "" { filter["collection_id"] = s }
	items, err := h.svc.List(c.Request().Context(), filter)
	if err != nil { return h.r.InternalServerErrorResponse(c, err.Error()) }
	return h.r.SuccessResponse(c, response.BmFromEntities(items), "")
}

// Add adds a new bookmark
// @Summary Add bookmark
// @Description Add a new bookmark for a resource
// @Tags Bookmarks
// @Accept json
// @Produce json
// @Param request body object true "Bookmark data" Schema({Properties:{resource_type:{type:string},resource_id:{type:integer},note:{type:string},tags:{type:string}}})
// @Success 201 {object} utils.CreatedResponse "Bookmark created"
// @Failure 400 {object} utils.ErrorResponse "Invalid request"
// @Failure 500 {object} utils.ErrorResponse "Internal server error"
// @Router /api/v1/bookmarks [post]
// @Security Bearer
func (h *BookmarkHandler) Add(c echo.Context) error {
	req := struct {
		ResourceType string `json:"resource_type" validate:"required"`
		ResourceID   uint   `json:"resource_id" validate:"required"`
		Note         string `json:"note"`
		Tags         string `json:"tags"`
	}{}
	if err := c.Bind(&req); err != nil { return h.r.BadRequestResponse(c, err.Error()) }
	bm, err := h.svc.Add(c.Request().Context(), 0, req.ResourceType, req.ResourceID, nil, req.Note, req.Tags)
	if err != nil { return h.r.InternalServerErrorResponse(c, err.Error()) }
	return h.r.CreatedResponse(c, response.BmFromEntity(bm), "")
}

// UpdateProgress updates bookmark progress
// @Summary Update bookmark progress
// @Description Update the reading progress and status of a bookmark
// @Tags Bookmarks
// @Accept json
// @Produce json
// @Param id path int true "Bookmark ID"
// @Param request body object true "Progress data" Schema({Properties:{progress:{type:number},status:{type:string}}})
// @Success 200 {object} utils.SuccessResponse "Updated"
// @Failure 400 {object} utils.ErrorResponse "Invalid request"
// @Failure 500 {object} utils.ErrorResponse "Internal server error"
// @Router /api/v1/bookmarks/{id}/progress [patch]
// @Security Bearer
func (h *BookmarkHandler) UpdateProgress(c echo.Context) error {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 32)
	req := struct {
		Progress float64 `json:"progress"`
		Status   string  `json:"status"`
	}{}
	if err := c.Bind(&req); err != nil { return h.r.BadRequestResponse(c, err.Error()) }
	if err := h.svc.UpdateProgress(c.Request().Context(), uint(id), req.Progress, req.Status); err != nil {
		return h.r.InternalServerErrorResponse(c, err.Error())
	}
	return h.r.SuccessResponse(c, nil, "updated")
}

// Remove removes a bookmark
// @Summary Remove bookmark
// @Description Remove a bookmark by ID
// @Tags Bookmarks
// @Produce json
// @Param id path int true "Bookmark ID"
// @Success 204 "No content"
// @Failure 500 {object} utils.ErrorResponse "Internal server error"
// @Router /api/v1/bookmarks/{id} [delete]
// @Security Bearer
func (h *BookmarkHandler) Remove(c echo.Context) error {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 32)
	if err := h.svc.Remove(c.Request().Context(), uint(id)); err != nil {
		return h.r.InternalServerErrorResponse(c, err.Error())
	}
	return h.r.NoContentResponse(c)
}

// Collections
// ListCollections lists bookmark collections
// @Summary List collections
// @Description List all bookmark collections for the current user
// @Tags Bookmarks
// @Produce json
// @Success 200 {object} utils.SuccessResponse "Collections list"
// @Failure 500 {object} utils.ErrorResponse "Internal server error"
// @Router /api/v1/bookmarks/collections [get]
// @Security Bearer
func (h *BookmarkHandler) ListCollections(c echo.Context) error {
	items, err := h.svc.ListCollections(c.Request().Context(), 0)
	if err != nil { return h.r.InternalServerErrorResponse(c, err.Error()) }
	return h.r.SuccessResponse(c, items, "")
}

// CreateCollection creates a new collection
// @Summary Create collection
// @Description Create a new bookmark collection
// @Tags Bookmarks
// @Accept json
// @Produce json
// @Param request body entity.Collection true "Collection data"
// @Success 201 {object} utils.CreatedResponse "Collection created"
// @Failure 400 {object} utils.ErrorResponse "Invalid request"
// @Failure 500 {object} utils.ErrorResponse "Internal server error"
// @Router /api/v1/bookmarks/collections [post]
// @Security Bearer
func (h *BookmarkHandler) CreateCollection(c echo.Context) error {
	var col entity.Collection
	if err := c.Bind(&col); err != nil { return h.r.BadRequestResponse(c, err.Error()) }
	if err := h.svc.CreateCollection(c.Request().Context(), &col); err != nil {
		return h.r.InternalServerErrorResponse(c, err.Error())
	}
	return h.r.CreatedResponse(c, col, "")
}

// DeleteCollection deletes a collection
// @Summary Delete collection
// @Description Delete a bookmark collection by ID
// @Tags Bookmarks
// @Produce json
// @Param id path int true "Collection ID"
// @Success 204 "No content"
// @Failure 500 {object} utils.ErrorResponse "Internal server error"
// @Router /api/v1/bookmarks/collections/{id} [delete]
// @Security Bearer
func (h *BookmarkHandler) DeleteCollection(c echo.Context) error {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 32)
	if err := h.svc.DeleteCollection(c.Request().Context(), uint(id)); err != nil {
		return h.r.InternalServerErrorResponse(c, err.Error())
	}
	return h.r.NoContentResponse(c)
}

// Goal
// GetGoal gets reading goal
// @Summary Get reading goal
// @Description Get the reading goal for a specific year
// @Tags Bookmarks
// @Produce json
// @Param year query int false "Year (default 2026)"
// @Success 200 {object} utils.SuccessResponse "Reading goal"
// @Failure 500 {object} utils.ErrorResponse "Internal server error"
// @Router /api/v1/bookmarks/goal [get]
// @Security Bearer
func (h *BookmarkHandler) GetGoal(c echo.Context) error {
	year, _ := strconv.Atoi(c.QueryParam("year"))
	if year == 0 { year = 2026 }
	g, err := h.svc.GetGoal(c.Request().Context(), 0, year)
	if err != nil { return h.r.InternalServerErrorResponse(c, err.Error()) }
	return h.r.SuccessResponse(c, g, "")
}

// SetGoal sets reading goal
// @Summary Set reading goal
// @Description Set the reading goal count for a specific year
// @Tags Bookmarks
// @Accept json
// @Produce json
// @Param request body object true "Goal data" Schema({Properties:{year:{type:integer},count:{type:integer}}})
// @Success 200 {object} utils.SuccessResponse "Goal set"
// @Failure 400 {object} utils.ErrorResponse "Invalid request"
// @Failure 500 {object} utils.ErrorResponse "Internal server error"
// @Router /api/v1/bookmarks/goal [post]
// @Security Bearer
func (h *BookmarkHandler) SetGoal(c echo.Context) error {
	req := struct {
		Year  int `json:"year" validate:"required"`
		Count int `json:"count" validate:"required"`
	}{}
	if err := c.Bind(&req); err != nil { return h.r.BadRequestResponse(c, err.Error()) }
	g, err := h.svc.SetGoal(c.Request().Context(), 0, req.Year, req.Count)
	if err != nil { return h.r.InternalServerErrorResponse(c, err.Error()) }
	return h.r.SuccessResponse(c, g, "")
}

func (h *BookmarkHandler) RegisterRoutes(e *echo.Echo, basePath string) {
	g := e.Group(basePath+"/bookmarks", middleware.Auth)
	g.GET("", h.List)
	g.POST("", h.Add)
	g.PATCH("/:id/progress", h.UpdateProgress)
	g.DELETE("/:id", h.Remove)
	g.GET("/collections", h.ListCollections)
	g.POST("/collections", h.CreateCollection)
	g.DELETE("/collections/:id", h.DeleteCollection)
	g.GET("/goal", h.GetGoal)
	g.POST("/goal", h.SetGoal)
}
