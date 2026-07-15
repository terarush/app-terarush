package handler

import (
	"go-modular/internal/pkg/bus"
	"go-modular/internal/pkg/logger"
	"go-modular/internal/pkg/middleware"
	"go-modular/internal/pkg/utils"
	"go-modular/modules/notifications/domain/service"
	"go-modular/modules/notifications/dto/response"
	"strconv"

	"github.com/labstack/echo/v4"
)

type NotificationHandler struct {
	svc *service.NotificationService
	log *logger.Logger
	evt *bus.EventBus
	r   *utils.Response
}

func NewNotificationHandler(log *logger.Logger, evt *bus.EventBus, svc *service.NotificationService) *NotificationHandler {
	return &NotificationHandler{svc: svc, log: log, evt: evt, r: &utils.Response{}}
}

func (h *NotificationHandler) List(c echo.Context) error {
	filter := map[string]any{}
	if s := c.QueryParam("status"); s != "" {
		filter["status"] = s
	}
	items, err := h.svc.List(c.Request().Context(), filter)
	if err != nil {
		return h.r.InternalServerErrorResponse(c, err.Error())
	}
	return h.r.SuccessResponse(c, response.NotifFromEntities(items), "")
}

func (h *NotificationHandler) Get(c echo.Context) error {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 32)
	n, err := h.svc.Get(c.Request().Context(), uint(id))
	if err != nil {
		return h.r.NotFoundResponse(c, "notification not found")
	}
	return h.r.SuccessResponse(c, response.NotifFromEntity(n), "")
}

func (h *NotificationHandler) MarkRead(c echo.Context) error {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 32)
	if err := h.svc.MarkRead(c.Request().Context(), uint(id)); err != nil {
		return h.r.InternalServerErrorResponse(c, err.Error())
	}
	return h.r.SuccessResponse(c, nil, "marked as read")
}

func (h *NotificationHandler) MarkAllRead(c echo.Context) error {
	// ponytail: needs user ID from context
	return h.r.SuccessResponse(c, nil, "all marked read")
}

func (h *NotificationHandler) Delete(c echo.Context) error {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 32)
	if err := h.svc.Delete(c.Request().Context(), uint(id)); err != nil {
		return h.r.InternalServerErrorResponse(c, err.Error())
	}
	return h.r.NoContentResponse(c)
}

func (h *NotificationHandler) GetUnreadCount(c echo.Context) error {
	// ponytail: needs user ID from context
	return h.r.SuccessResponse(c, response.UnreadCountResponse{Count: 0}, "")
}

func (h *NotificationHandler) RegisterRoutes(e *echo.Echo, basePath string) {
	g := e.Group(basePath+"/notifications", middleware.Auth)
	g.GET("", h.List)
	g.GET("/:id", h.Get)
	g.PATCH("/:id/read", h.MarkRead)
	g.POST("/read-all", h.MarkAllRead)
	g.DELETE("/:id", h.Delete)
	g.GET("/unread-count", h.GetUnreadCount)

	adm := e.Group(basePath+"/admin/notifications", middleware.Auth, middleware.AdminOnly)
	adm.GET("/templates", func(c echo.Context) error { return h.r.SuccessResponse(c, nil, "") })
	adm.POST("/templates", func(c echo.Context) error { return h.r.SuccessResponse(c, nil, "") })
}
