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

// ListReports lists moderation reports
// @Summary List reports
// @Description List all moderation reports with optional status and severity filters (admin only)
// @Tags Moderation
// @Produce json
// @Param status query string false "Filter by status"
// @Param severity query string false "Filter by severity"
// @Success 200 {object} utils.SuccessResponse "Reports list"
// @Failure 500 {object} utils.ErrorResponse "Internal server error"
// @Router /api/v1/admin/moderation/reports [get]
// @Security Bearer
func (h *ModerationHandler) ListReports(c echo.Context) error {
	filter := map[string]any{}
	if s := c.QueryParam("status"); s != "" { filter["status"] = s }
	if s := c.QueryParam("severity"); s != "" { filter["severity"] = s }
	items, err := h.svc.ListReports(c.Request().Context(), filter)
	if err != nil { return h.r.InternalServerErrorResponse(c, err.Error()) }
	return h.r.SuccessResponse(c, items, "")
}

// GetReport gets a report
// @Summary Get report
// @Description Get a specific moderation report by ID (admin only)
// @Tags Moderation
// @Produce json
// @Param id path int true "Report ID"
// @Success 200 {object} utils.SuccessResponse "Report details"
// @Failure 404 {object} utils.ErrorResponse "Report not found"
// @Failure 500 {object} utils.ErrorResponse "Internal server error"
// @Router /api/v1/admin/moderation/reports/{id} [get]
// @Security Bearer
func (h *ModerationHandler) GetReport(c echo.Context) error {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 32)
	r, err := h.svc.GetReport(c.Request().Context(), uint(id))
	if err != nil { return h.r.NotFoundResponse(c, "report not found") }
	return h.r.SuccessResponse(c, r, "")
}

// ReviewReport reviews a report
// @Summary Review report
// @Description Review and take action on a moderation report (admin only)
// @Tags Moderation
// @Accept json
// @Produce json
// @Param id path int true "Report ID"
// @Param request body object true "Review data" Schema({Properties:{status:{type:string},notes:{type:string},action:{type:string}}})
// @Success 200 {object} utils.SuccessResponse "Report reviewed"
// @Failure 400 {object} utils.ErrorResponse "Invalid request"
// @Failure 500 {object} utils.ErrorResponse "Internal server error"
// @Router /api/v1/admin/moderation/reports/{id}/review [post]
// @Security Bearer
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

// ListFilters lists content filters
// @Summary List filters
// @Description List all content moderation filters (admin only)
// @Tags Moderation
// @Produce json
// @Success 200 {object} utils.SuccessResponse "Filters list"
// @Failure 500 {object} utils.ErrorResponse "Internal server error"
// @Router /api/v1/admin/moderation/filters [get]
// @Security Bearer
func (h *ModerationHandler) ListFilters(c echo.Context) error {
	items, err := h.svc.ListFilters(c.Request().Context(), true)
	if err != nil { return h.r.InternalServerErrorResponse(c, err.Error()) }
	return h.r.SuccessResponse(c, items, "")
}

// CreateFilter creates a content filter
// @Summary Create filter
// @Description Create a new content moderation filter (admin only)
// @Tags Moderation
// @Accept json
// @Produce json
// @Param request body entity.ContentFilter true "Filter data"
// @Success 201 {object} utils.CreatedResponse "Filter created"
// @Failure 400 {object} utils.ErrorResponse "Invalid request"
// @Failure 500 {object} utils.ErrorResponse "Internal server error"
// @Router /api/v1/admin/moderation/filters [post]
// @Security Bearer
func (h *ModerationHandler) CreateFilter(c echo.Context) error {
	var f entity.ContentFilter
	if err := c.Bind(&f); err != nil { return h.r.BadRequestResponse(c, err.Error()) }
	if err := h.svc.CreateFilter(c.Request().Context(), &f); err != nil { return h.r.InternalServerErrorResponse(c, err.Error()) }
	return h.r.CreatedResponse(c, f, "")
}

// DeleteFilter deletes a content filter
// @Summary Delete filter
// @Description Delete a content moderation filter by ID (admin only)
// @Tags Moderation
// @Produce json
// @Param id path int true "Filter ID"
// @Success 204 "No content"
// @Failure 500 {object} utils.ErrorResponse "Internal server error"
// @Router /api/v1/admin/moderation/filters/{id} [delete]
// @Security Bearer
func (h *ModerationHandler) DeleteFilter(c echo.Context) error {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 32)
	if err := h.svc.DeleteFilter(c.Request().Context(), uint(id)); err != nil {
		return h.r.InternalServerErrorResponse(c, err.Error())
	}
	return h.r.NoContentResponse(c)
}

// GetStats gets moderation stats
// @Summary Get moderation stats
// @Description Get moderation statistics like pending report count (admin only)
// @Tags Moderation
// @Produce json
// @Success 200 {object} utils.SuccessResponse "Stats"
// @Failure 500 {object} utils.ErrorResponse "Internal server error"
// @Router /api/v1/admin/moderation/stats [get]
// @Security Bearer
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
