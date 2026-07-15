package handler

import (
	"go-modular/internal/pkg/bus"
	"go-modular/internal/pkg/logger"
	"go-modular/internal/pkg/middleware"
	"go-modular/internal/pkg/utils"
	"go-modular/modules/activity/domain/service"
	"strconv"

	"github.com/labstack/echo/v4"
)

type ActivityHandler struct {
	svc *service.ActivityService
	log *logger.Logger
	evt *bus.EventBus
	r   *utils.Response
}

func NewActivityHandler(log *logger.Logger, evt *bus.EventBus, svc *service.ActivityService) *ActivityHandler {
	return &ActivityHandler{svc: svc, log: log, evt: evt, r: &utils.Response{}}
}

func (h *ActivityHandler) List(c echo.Context) error {
	filter := map[string]any{}
	if t := c.QueryParam("type"); t != "" { filter["type"] = t }
	items, err := h.svc.List(c.Request().Context(), filter)
	if err != nil { return h.r.InternalServerErrorResponse(c, err.Error()) }
	return h.r.SuccessResponse(c, items, "")
}

func (h *ActivityHandler) GetFeed(c echo.Context) error {
	limit, _ := strconv.Atoi(c.QueryParam("limit"))
	if limit < 1 { limit = 20 }
	items, err := h.svc.GetFeed(c.Request().Context(), 0, limit)
	if err != nil { return h.r.InternalServerErrorResponse(c, err.Error()) }
	return h.r.SuccessResponse(c, items, "")
}

func (h *ActivityHandler) MarkFeedRead(c echo.Context) error {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 32)
	if err := h.svc.MarkFeedRead(c.Request().Context(), uint(id)); err != nil {
		return h.r.InternalServerErrorResponse(c, err.Error())
	}
	return h.r.SuccessResponse(c, nil, "marked read")
}

func (h *ActivityHandler) UnreadCount(c echo.Context) error {
	count, _ := h.svc.GetUnreadCount(c.Request().Context(), 0)
	return h.r.SuccessResponse(c, map[string]int64{"count": count}, "")
}

func (h *ActivityHandler) Follow(c echo.Context) error {
	id, _ := strconv.ParseUint(c.Param("userId"), 10, 32)
	if err := h.svc.Follow(c.Request().Context(), 0, uint(id)); err != nil {
		return h.r.InternalServerErrorResponse(c, err.Error())
	}
	return h.r.SuccessResponse(c, nil, "following")
}

func (h *ActivityHandler) Unfollow(c echo.Context) error {
	id, _ := strconv.ParseUint(c.Param("userId"), 10, 32)
	if err := h.svc.Unfollow(c.Request().Context(), 0, uint(id)); err != nil {
		return h.r.InternalServerErrorResponse(c, err.Error())
	}
	return h.r.SuccessResponse(c, nil, "unfollowed")
}

func (h *ActivityHandler) GetFollowers(c echo.Context) error {
	items, err := h.svc.GetFollowers(c.Request().Context(), 0)
	if err != nil { return h.r.InternalServerErrorResponse(c, err.Error()) }
	return h.r.SuccessResponse(c, items, "")
}

func (h *ActivityHandler) GetFollowing(c echo.Context) error {
	items, err := h.svc.GetFollowing(c.Request().Context(), 0)
	if err != nil { return h.r.InternalServerErrorResponse(c, err.Error()) }
	return h.r.SuccessResponse(c, items, "")
}

func (h *ActivityHandler) GetMetrics(c echo.Context) error {
	m, err := h.svc.GetMetrics(c.Request().Context(), 0)
	if err != nil { return h.r.InternalServerErrorResponse(c, err.Error()) }
	return h.r.SuccessResponse(c, m, "")
}

func (h *ActivityHandler) GetLeaderboard(c echo.Context) error {
	limit, _ := strconv.Atoi(c.QueryParam("limit"))
	if limit < 1 { limit = 10 }
	items, err := h.svc.GetLeaderboard(c.Request().Context(), limit)
	if err != nil { return h.r.InternalServerErrorResponse(c, err.Error()) }
	return h.r.SuccessResponse(c, items, "")
}

func (h *ActivityHandler) RegisterRoutes(e *echo.Echo, basePath string) {
	g := e.Group(basePath+"/activity", middleware.Auth)
	g.GET("", h.List)
	g.GET("/feed", h.GetFeed)
	g.PATCH("/feed/:id/read", h.MarkFeedRead)
	g.GET("/feed/unread", h.UnreadCount)

	g.POST("/follow/:userId", h.Follow)
	g.DELETE("/follow/:userId", h.Unfollow)
	g.GET("/followers", h.GetFollowers)
	g.GET("/following", h.GetFollowing)

	g.GET("/metrics", h.GetMetrics)
	g.GET("/leaderboard", h.GetLeaderboard)
}
