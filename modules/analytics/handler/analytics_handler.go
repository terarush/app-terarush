package handler

import (
	"go-modular/internal/pkg/bus"
	"go-modular/internal/pkg/logger"
	"go-modular/internal/pkg/middleware"
	"go-modular/internal/pkg/utils"
	"go-modular/modules/analytics/domain/service"
	"strconv"

	"github.com/labstack/echo/v4"
)

type AnalyticsHandler struct {
	svc *service.AnalyticsService
	log *logger.Logger
	evt *bus.EventBus
	r   *utils.Response
}

func NewAnalyticsHandler(log *logger.Logger, evt *bus.EventBus, svc *service.AnalyticsService) *AnalyticsHandler {
	return &AnalyticsHandler{svc: svc, log: log, evt: evt, r: &utils.Response{}}
}

func (h *AnalyticsHandler) ListEvents(c echo.Context) error {
	filter := map[string]any{}
	if t := c.QueryParam("event_type"); t != "" {
		filter["event_type"] = t
	}
	items, err := h.svc.ListEvents(c.Request().Context(), filter)
	if err != nil {
		return h.r.InternalServerErrorResponse(c, err.Error())
	}
	return h.r.SuccessResponse(c, items, "")
}

func (h *AnalyticsHandler) GetDashboard(c echo.Context) error {
	stats, err := h.svc.GetDashboardStats(c.Request().Context())
	if err != nil {
		return h.r.InternalServerErrorResponse(c, err.Error())
	}
	return h.r.SuccessResponse(c, stats, "")
}

func (h *AnalyticsHandler) GetTopBlogs(c echo.Context) error {
	limit, _ := strconv.Atoi(c.QueryParam("limit"))
	if limit < 1 { limit = 10 }
	blogs, err := h.svc.GetTopBlogs(c.Request().Context(), limit)
	if err != nil {
		return h.r.InternalServerErrorResponse(c, err.Error())
	}
	return h.r.SuccessResponse(c, blogs, "")
}

func (h *AnalyticsHandler) RegisterRoutes(e *echo.Echo, basePath string) {
	adm := e.Group(basePath+"/admin/analytics", middleware.Auth, middleware.AdminOnly)
	adm.GET("/events", h.ListEvents)
	adm.GET("/dashboard", h.GetDashboard)
	adm.GET("/top-blogs", h.GetTopBlogs)
	adm.POST("/events/track", func(c echo.Context) error {
		return h.r.SuccessResponse(c, nil, "tracked")
	})
}
