package handler

import (
	"go-modular/internal/pkg/bus"
	"go-modular/internal/pkg/logger"
	"go-modular/internal/pkg/middleware"
	"go-modular/internal/pkg/utils"
	"go-modular/modules/moderation/domain/entity"
	"go-modular/modules/moderation/domain/service"
	"strconv"

	"github.com/labstack/echo/v4"
)

type ModerationHandler struct {
	svc *service.ModerationService
	log *logger.Logger
	evt *bus.EventBus
	r   *utils.Response
}

func NewModerationHandler(log *logger.Logger, evt *bus.EventBus, svc *service.ModerationService) *ModerationHandler {
	return &ModerationHandler{svc: svc, log: log, evt: evt, r: &utils.Response{}}
}

func (h *ModerationHandler) ListReports(c echo.Context) error {
	filter := map[string]any{}
	if s := c.QueryParam("status"); s != "" { filter["status"] = s }
	if s := c.QueryParam("severity"); s != "" { filter["severity"] = s }
	items, err := h.svc.ListReports(c.Request().Context(), filter)
	if err != nil { return h.r.InternalServerErrorResponse(c, err.Error()) }
	return h.r.SuccessResponse(c, items, "")
}

func (h *ModerationHandler) GetReport(c echo.Context) error {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 32)
	r, err := h.svc.GetReport(c.Request().Context(), uint(id))
	if err != nil { return h.r.NotFoundResponse(c, "report not found") }
	return h.r.SuccessResponse(c, r, "")
}

func (h *ModerationHandler) ReviewReport(c echo.Context) error {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 32)
	req := struct {
		Status string `json:"status" validate:"required"`
		Notes  string `json:"notes"`
		Action string `json:"action"`
	}{}
	if err := c.Bind(&req); err != nil { return h.r.BadRequestResponse(c, err.Error()) }
	if err := h.svc.ReviewReport(c.Request().Context(), uint(id), 0, req.Status, req.Notes, req.Action); err != nil {
		return h.r.InternalServerErrorResponse(c, err.Error())
	}
	return h.r.SuccessResponse(c, nil, "report reviewed")
}

func (h *ModerationHandler) ListFilters(c echo.Context) error {
	items, err := h.svc.ListFilters(c.Request().Context(), true)
	if err != nil { return h.r.InternalServerErrorResponse(c, err.Error()) }
	return h.r.SuccessResponse(c, items, "")
}

func (h *ModerationHandler) CreateFilter(c echo.Context) error {
	var f entity.ContentFilter
	if err := c.Bind(&f); err != nil { return h.r.BadRequestResponse(c, err.Error()) }
	if err := h.svc.CreateFilter(c.Request().Context(), &f); err != nil { return h.r.InternalServerErrorResponse(c, err.Error()) }
	return h.r.CreatedResponse(c, f, "")
}

func (h *ModerationHandler) DeleteFilter(c echo.Context) error {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 32)
	if err := h.svc.DeleteFilter(c.Request().Context(), uint(id)); err != nil {
		return h.r.InternalServerErrorResponse(c, err.Error())
	}
	return h.r.NoContentResponse(c)
}

func (h *ModerationHandler) GetStats(c echo.Context) error {
	pending, _ := h.svc.GetPendingCount(c.Request().Context())
	return h.r.SuccessResponse(c, map[string]int64{"pending_reports": pending}, "")
}

func (h *ModerationHandler) RegisterRoutes(e *echo.Echo, basePath string) {
	adm := e.Group(basePath+"/admin/moderation", middleware.Auth, middleware.AdminOnly)
	adm.GET("/reports", h.ListReports)
	adm.GET("/reports/:id", h.GetReport)
	adm.POST("/reports/:id/review", h.ReviewReport)
	adm.GET("/filters", h.ListFilters)
	adm.POST("/filters", h.CreateFilter)
	adm.DELETE("/filters/:id", h.DeleteFilter)
	adm.GET("/stats", h.GetStats)

	pub := e.Group(basePath + "/reports")
	pub.POST("", func(c echo.Context) error { return h.r.SuccessResponse(c, nil, "report submitted") })
}
