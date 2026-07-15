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

func (h *BookmarkHandler) List(c echo.Context) error {
	filter := map[string]any{}
	if s := c.QueryParam("status"); s != "" { filter["status"] = s }
	if s := c.QueryParam("collection_id"); s != "" { filter["collection_id"] = s }
	items, err := h.svc.List(c.Request().Context(), filter)
	if err != nil { return h.r.InternalServerErrorResponse(c, err.Error()) }
	return h.r.SuccessResponse(c, response.BmFromEntities(items), "")
}

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

func (h *BookmarkHandler) Remove(c echo.Context) error {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 32)
	if err := h.svc.Remove(c.Request().Context(), uint(id)); err != nil {
		return h.r.InternalServerErrorResponse(c, err.Error())
	}
	return h.r.NoContentResponse(c)
}

// Collections
func (h *BookmarkHandler) ListCollections(c echo.Context) error {
	items, err := h.svc.ListCollections(c.Request().Context(), 0)
	if err != nil { return h.r.InternalServerErrorResponse(c, err.Error()) }
	return h.r.SuccessResponse(c, items, "")
}

func (h *BookmarkHandler) CreateCollection(c echo.Context) error {
	var col entity.Collection
	if err := c.Bind(&col); err != nil { return h.r.BadRequestResponse(c, err.Error()) }
	if err := h.svc.CreateCollection(c.Request().Context(), &col); err != nil {
		return h.r.InternalServerErrorResponse(c, err.Error())
	}
	return h.r.CreatedResponse(c, col, "")
}

func (h *BookmarkHandler) DeleteCollection(c echo.Context) error {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 32)
	if err := h.svc.DeleteCollection(c.Request().Context(), uint(id)); err != nil {
		return h.r.InternalServerErrorResponse(c, err.Error())
	}
	return h.r.NoContentResponse(c)
}

// Goal
func (h *BookmarkHandler) GetGoal(c echo.Context) error {
	year, _ := strconv.Atoi(c.QueryParam("year"))
	if year == 0 { year = 2026 }
	g, err := h.svc.GetGoal(c.Request().Context(), 0, year)
	if err != nil { return h.r.InternalServerErrorResponse(c, err.Error()) }
	return h.r.SuccessResponse(c, g, "")
}

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
