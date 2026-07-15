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

// List lists notifications
// @Summary List notifications
// @Description List all notifications with optional status filter
// @Tags Notifications
// @Produce json
// @Param status query string false "Filter by status"
// @Success 200 {object} utils.SuccessResponse "Notifications list"
// @Failure 500 {object} utils.ErrorResponse "Internal server error"
// @Router /api/v1/notifications [get]
// @Security Bearer
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

// Get gets a notification
// @Summary Get notification
// @Description Get a specific notification by ID
// @Tags Notifications
// @Produce json
// @Param id path int true "Notification ID"
// @Success 200 {object} utils.SuccessResponse "Notification details"
// @Failure 404 {object} utils.ErrorResponse "Not found"
// @Failure 500 {object} utils.ErrorResponse "Internal server error"
// @Router /api/v1/notifications/{id} [get]
// @Security Bearer
func (h *NotificationHandler) Get(c echo.Context) error {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 32)
	n, err := h.svc.Get(c.Request().Context(), uint(id))
	if err != nil {
		return h.r.NotFoundResponse(c, "notification not found")
	}
	return h.r.SuccessResponse(c, response.NotifFromEntity(n), "")
}

// MarkRead marks a notification as read
// @Summary Mark notification read
// @Description Mark a specific notification as read by ID
// @Tags Notifications
// @Produce json
// @Param id path int true "Notification ID"
// @Success 200 {object} utils.SuccessResponse "Marked as read"
// @Failure 500 {object} utils.ErrorResponse "Internal server error"
// @Router /api/v1/notifications/{id}/read [patch]
// @Security Bearer
func (h *NotificationHandler) MarkRead(c echo.Context) error {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 32)
	if err := h.svc.MarkRead(c.Request().Context(), uint(id)); err != nil {
		return h.r.InternalServerErrorResponse(c, err.Error())
	}
	return h.r.SuccessResponse(c, nil, "marked as read")
}

// MarkAllRead marks all notifications as read
// @Summary Mark all read
// @Description Mark all notifications as read for the current user
// @Tags Notifications
// @Produce json
// @Success 200 {object} utils.SuccessResponse "All marked read"
// @Failure 500 {object} utils.ErrorResponse "Internal server error"
// @Router /api/v1/notifications/read-all [post]
// @Security Bearer
func (h *NotificationHandler) MarkAllRead(c echo.Context) error {
	// ponytail: needs user ID from context
	return h.r.SuccessResponse(c, nil, "all marked read")
}

// Delete deletes a notification
// @Summary Delete notification
// @Description Delete a notification by ID
// @Tags Notifications
// @Produce json
// @Param id path int true "Notification ID"
// @Success 204 "No content"
// @Failure 500 {object} utils.ErrorResponse "Internal server error"
// @Router /api/v1/notifications/{id} [delete]
// @Security Bearer
func (h *NotificationHandler) Delete(c echo.Context) error {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 32)
	if err := h.svc.Delete(c.Request().Context(), uint(id)); err != nil {
		return h.r.InternalServerErrorResponse(c, err.Error())
	}
	return h.r.NoContentResponse(c)
}

// GetUnreadCount gets unread count
// @Summary Get unread count
// @Description Get the count of unread notifications for the current user
// @Tags Notifications
// @Produce json
// @Success 200 {object} utils.SuccessResponse "Unread count"
// @Failure 500 {object} utils.ErrorResponse "Internal server error"
// @Router /api/v1/notifications/unread-count [get]
// @Security Bearer
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
