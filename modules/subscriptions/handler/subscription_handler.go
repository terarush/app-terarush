package handler

import (
	"go-modular/internal/pkg/bus"
	"go-modular/internal/pkg/logger"
	"go-modular/internal/pkg/middleware"
	"go-modular/internal/pkg/utils"
	"go-modular/modules/subscriptions/domain/entity"
	"go-modular/modules/subscriptions/domain/service"
	"go-modular/modules/subscriptions/dto/response"
	"strconv"

	"github.com/labstack/echo/v4"
)

type SubscriptionHandler struct {
	svc *service.SubscriptionService
	log *logger.Logger
	evt *bus.EventBus
	r   *utils.Response
}

func NewSubscriptionHandler(log *logger.Logger, evt *bus.EventBus, svc *service.SubscriptionService) *SubscriptionHandler {
	return &SubscriptionHandler{svc: svc, log: log, evt: evt, r: &utils.Response{}}
}

// List lists subscriptions
// @Summary List subscriptions
// @Description List all subscriptions with optional resource_type filter
// @Tags Subscriptions
// @Produce json
// @Param resource_type query string false "Filter by resource type"
// @Success 200 {object} utils.SuccessResponse "Subscriptions list"
// @Failure 500 {object} utils.ErrorResponse "Internal server error"
// @Router /api/v1/subscriptions [get]
// @Security Bearer
func (h *SubscriptionHandler) List(c echo.Context) error {
	filter := map[string]any{}
	if t := c.QueryParam("resource_type"); t != "" { filter["resource_type"] = t }
	items, err := h.svc.List(c.Request().Context(), filter)
	if err != nil { return h.r.InternalServerErrorResponse(c, err.Error()) }
	return h.r.SuccessResponse(c, response.SubFromEntities(items), "")
}

// Get gets a subscription
// @Summary Get subscription
// @Description Get a specific subscription by ID
// @Tags Subscriptions
// @Produce json
// @Param id path int true "Subscription ID"
// @Success 200 {object} utils.SuccessResponse "Subscription details"
// @Failure 404 {object} utils.ErrorResponse "Not found"
// @Failure 500 {object} utils.ErrorResponse "Internal server error"
// @Router /api/v1/subscriptions/{id} [get]
// @Security Bearer
func (h *SubscriptionHandler) Get(c echo.Context) error {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 32)
	s, err := h.svc.Get(c.Request().Context(), uint(id))
	if err != nil { return h.r.NotFoundResponse(c, "subscription not found") }
	return h.r.SuccessResponse(c, response.SubFromEntity(s), "")
}

// Subscribe subscribes to a resource
// @Summary Subscribe to resource
// @Description Subscribe to a resource by type and ID
// @Tags Subscriptions
// @Accept json
// @Produce json
// @Param request body object true "Subscription data" Schema({Properties:{resource_type:{type:string},resource_id:{type:integer}}})
// @Success 201 {object} utils.CreatedResponse "Subscription created"
// @Failure 400 {object} utils.ErrorResponse "Invalid request"
// @Failure 500 {object} utils.ErrorResponse "Internal server error"
// @Router /api/v1/subscriptions [post]
// @Security Bearer
func (h *SubscriptionHandler) Subscribe(c echo.Context) error {
	req := struct {
		ResourceType string `json:"resource_type" validate:"required"`
		ResourceID   uint   `json:"resource_id" validate:"required"`
	}{}
	if err := c.Bind(&req); err != nil { return h.r.BadRequestResponse(c, err.Error()) }
	sub, err := h.svc.Subscribe(c.Request().Context(), 0, req.ResourceType, req.ResourceID)
	if err != nil { return h.r.InternalServerErrorResponse(c, err.Error()) }
	return h.r.CreatedResponse(c, response.SubFromEntity(sub), "")
}

// Unsubscribe removes a subscription
// @Summary Unsubscribe
// @Description Unsubscribe by subscription ID
// @Tags Subscriptions
// @Produce json
// @Param id path int true "Subscription ID"
// @Success 200 {object} utils.SuccessResponse "Unsubscribed"
// @Failure 500 {object} utils.ErrorResponse "Internal server error"
// @Router /api/v1/subscriptions/{id} [delete]
// @Security Bearer
func (h *SubscriptionHandler) Unsubscribe(c echo.Context) error {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 32)
	if err := h.svc.Unsubscribe(c.Request().Context(), uint(id)); err != nil {
		return h.r.InternalServerErrorResponse(c, err.Error())
	}
	return h.r.SuccessResponse(c, nil, "unsubscribed")
}

// ListNewsletters lists newsletters
// @Summary List newsletters
// @Description List all newsletters (admin only)
// @Tags Subscriptions
// @Produce json
// @Success 200 {object} utils.SuccessResponse "Newsletters list"
// @Failure 500 {object} utils.ErrorResponse "Internal server error"
// @Router /api/v1/admin/subscriptions/newsletters [get]
// @Security Bearer
func (h *SubscriptionHandler) ListNewsletters(c echo.Context) error {
	items, err := h.svc.ListNewsletters(c.Request().Context(), nil)
	if err != nil { return h.r.InternalServerErrorResponse(c, err.Error()) }
	return h.r.SuccessResponse(c, items, "")
}

// CreateNewsletter creates a newsletter
// @Summary Create newsletter
// @Description Create a new newsletter (admin only)
// @Tags Subscriptions
// @Accept json
// @Produce json
// @Param request body entity.Newsletter true "Newsletter data"
// @Success 201 {object} utils.CreatedResponse "Newsletter created"
// @Failure 400 {object} utils.ErrorResponse "Invalid request"
// @Failure 500 {object} utils.ErrorResponse "Internal server error"
// @Router /api/v1/admin/subscriptions/newsletters [post]
// @Security Bearer
func (h *SubscriptionHandler) CreateNewsletter(c echo.Context) error {
	var n entity.Newsletter
	if err := c.Bind(&n); err != nil { return h.r.BadRequestResponse(c, err.Error()) }
	if err := h.svc.CreateNewsletter(c.Request().Context(), &n); err != nil {
		return h.r.InternalServerErrorResponse(c, err.Error())
	}
	return h.r.CreatedResponse(c, n, "")
}

// SendNewsletter sends a newsletter
// @Summary Send newsletter
// @Description Send a newsletter to all subscribers by ID (admin only)
// @Tags Subscriptions
// @Produce json
// @Param id path int true "Newsletter ID"
// @Success 200 {object} utils.SuccessResponse "Newsletter sent"
// @Failure 500 {object} utils.ErrorResponse "Internal server error"
// @Router /api/v1/admin/subscriptions/newsletters/{id}/send [post]
// @Security Bearer
func (h *SubscriptionHandler) SendNewsletter(c echo.Context) error {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 32)
	if err := h.svc.SendNewsletter(c.Request().Context(), uint(id)); err != nil {
		return h.r.InternalServerErrorResponse(c, err.Error())
	}
	return h.r.SuccessResponse(c, nil, "newsletter sent")
}

// DeleteNewsletter deletes a newsletter
// @Summary Delete newsletter
// @Description Delete a newsletter by ID (admin only)
// @Tags Subscriptions
// @Produce json
// @Param id path int true "Newsletter ID"
// @Success 204 "No content"
// @Failure 500 {object} utils.ErrorResponse "Internal server error"
// @Router /api/v1/admin/subscriptions/newsletters/{id} [delete]
// @Security Bearer
func (h *SubscriptionHandler) DeleteNewsletter(c echo.Context) error {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 32)
	if err := h.svc.DeleteNewsletter(c.Request().Context(), uint(id)); err != nil {
		return h.r.InternalServerErrorResponse(c, err.Error())
	}
	return h.r.NoContentResponse(c)
}

func (h *SubscriptionHandler) RegisterRoutes(e *echo.Echo, basePath string) {
	usr := e.Group(basePath+"/subscriptions", middleware.Auth)
	usr.GET("", h.List)
	usr.POST("", h.Subscribe)
	usr.DELETE("/:id", h.Unsubscribe)

	adm := e.Group(basePath+"/admin/subscriptions", middleware.Auth, middleware.AdminOnly)
	adm.GET("/newsletters", h.ListNewsletters)
	adm.POST("/newsletters", h.CreateNewsletter)
	adm.POST("/newsletters/:id/send", h.SendNewsletter)
	adm.DELETE("/newsletters/:id", h.DeleteNewsletter)
}
