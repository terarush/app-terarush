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

// List gets all activities
// @Summary List activities
// @Description List all activities with optional type filter
// @Tags Activity
// @Produce json
// @Param type query string false "Filter by activity type"
// @Success 200 {object} utils.SuccessResponse "Activities list"
// @Failure 500 {object} utils.ErrorResponse "Internal server error"
// @Router /api/v1/activity [get]
// @Security Bearer
func (h *ActivityHandler) List(c echo.Context) error {
	filter := map[string]any{}
	if t := c.QueryParam("type"); t != "" { filter["type"] = t }
	items, err := h.svc.List(c.Request().Context(), filter)
	if err != nil { return h.r.InternalServerErrorResponse(c, err.Error()) }
	return h.r.SuccessResponse(c, items, "")
}

// GetFeed gets activity feed
// @Summary Get activity feed
// @Description Get paginated activity feed for current user
// @Tags Activity
// @Produce json
// @Param limit query int false "Number of items (default 20)"
// @Success 200 {object} utils.SuccessResponse "Feed items"
// @Failure 500 {object} utils.ErrorResponse "Internal server error"
// @Router /api/v1/activity/feed [get]
// @Security Bearer
func (h *ActivityHandler) GetFeed(c echo.Context) error {
	limit, _ := strconv.Atoi(c.QueryParam("limit"))
	if limit < 1 { limit = 20 }
	items, err := h.svc.GetFeed(c.Request().Context(), 0, limit)
	if err != nil { return h.r.InternalServerErrorResponse(c, err.Error()) }
	return h.r.SuccessResponse(c, items, "")
}

// MarkFeedRead marks a feed item as read
// @Summary Mark feed item read
// @Description Mark a single feed item as read by ID
// @Tags Activity
// @Produce json
// @Param id path int true "Feed item ID"
// @Success 200 {object} utils.SuccessResponse "Marked read"
// @Failure 500 {object} utils.ErrorResponse "Internal server error"
// @Router /api/v1/activity/feed/{id}/read [patch]
// @Security Bearer
func (h *ActivityHandler) MarkFeedRead(c echo.Context) error {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 32)
	if err := h.svc.MarkFeedRead(c.Request().Context(), uint(id)); err != nil {
		return h.r.InternalServerErrorResponse(c, err.Error())
	}
	return h.r.SuccessResponse(c, nil, "marked read")
}

// UnreadCount gets unread feed count
// @Summary Get unread count
// @Description Get the count of unread feed items
// @Tags Activity
// @Produce json
// @Success 200 {object} utils.SuccessResponse "Unread count"
// @Failure 500 {object} utils.ErrorResponse "Internal server error"
// @Router /api/v1/activity/feed/unread [get]
// @Security Bearer
func (h *ActivityHandler) UnreadCount(c echo.Context) error {
	count, _ := h.svc.GetUnreadCount(c.Request().Context(), 0)
	return h.r.SuccessResponse(c, map[string]int64{"count": count}, "")
}

// Follow follows a user
// @Summary Follow user
// @Description Follow another user by user ID
// @Tags Activity
// @Produce json
// @Param userId path int true "User ID to follow"
// @Success 200 {object} utils.SuccessResponse "Following"
// @Failure 500 {object} utils.ErrorResponse "Internal server error"
// @Router /api/v1/activity/follow/{userId} [post]
// @Security Bearer
func (h *ActivityHandler) Follow(c echo.Context) error {
	id, _ := strconv.ParseUint(c.Param("userId"), 10, 32)
	if err := h.svc.Follow(c.Request().Context(), 0, uint(id)); err != nil {
		return h.r.InternalServerErrorResponse(c, err.Error())
	}
	return h.r.SuccessResponse(c, nil, "following")
}

// Unfollow unfollows a user
// @Summary Unfollow user
// @Description Unfollow a user by user ID
// @Tags Activity
// @Produce json
// @Param userId path int true "User ID to unfollow"
// @Success 200 {object} utils.SuccessResponse "Unfollowed"
// @Failure 500 {object} utils.ErrorResponse "Internal server error"
// @Router /api/v1/activity/follow/{userId} [delete]
// @Security Bearer
func (h *ActivityHandler) Unfollow(c echo.Context) error {
	id, _ := strconv.ParseUint(c.Param("userId"), 10, 32)
	if err := h.svc.Unfollow(c.Request().Context(), 0, uint(id)); err != nil {
		return h.r.InternalServerErrorResponse(c, err.Error())
	}
	return h.r.SuccessResponse(c, nil, "unfollowed")
}

// GetFollowers gets followers list
// @Summary Get followers
// @Description Get list of users following the current user
// @Tags Activity
// @Produce json
// @Success 200 {object} utils.SuccessResponse "Followers list"
// @Failure 500 {object} utils.ErrorResponse "Internal server error"
// @Router /api/v1/activity/followers [get]
// @Security Bearer
func (h *ActivityHandler) GetFollowers(c echo.Context) error {
	items, err := h.svc.GetFollowers(c.Request().Context(), 0)
	if err != nil { return h.r.InternalServerErrorResponse(c, err.Error()) }
	return h.r.SuccessResponse(c, items, "")
}

// GetFollowing gets following list
// @Summary Get following
// @Description Get list of users the current user follows
// @Tags Activity
// @Produce json
// @Success 200 {object} utils.SuccessResponse "Following list"
// @Failure 500 {object} utils.ErrorResponse "Internal server error"
// @Router /api/v1/activity/following [get]
// @Security Bearer
func (h *ActivityHandler) GetFollowing(c echo.Context) error {
	items, err := h.svc.GetFollowing(c.Request().Context(), 0)
	if err != nil { return h.r.InternalServerErrorResponse(c, err.Error()) }
	return h.r.SuccessResponse(c, items, "")
}

// GetMetrics gets activity metrics
// @Summary Get activity metrics
// @Description Get aggregated activity metrics for the current user
// @Tags Activity
// @Produce json
// @Success 200 {object} utils.SuccessResponse "Activity metrics"
// @Failure 500 {object} utils.ErrorResponse "Internal server error"
// @Router /api/v1/activity/metrics [get]
// @Security Bearer
func (h *ActivityHandler) GetMetrics(c echo.Context) error {
	m, err := h.svc.GetMetrics(c.Request().Context(), 0)
	if err != nil { return h.r.InternalServerErrorResponse(c, err.Error()) }
	return h.r.SuccessResponse(c, m, "")
}

// GetLeaderboard gets leaderboard
// @Summary Get leaderboard
// @Description Get the leaderboard of most active users
// @Tags Activity
// @Produce json
// @Param limit query int false "Number of top users (default 10)"
// @Success 200 {object} utils.SuccessResponse "Leaderboard"
// @Failure 500 {object} utils.ErrorResponse "Internal server error"
// @Router /api/v1/activity/leaderboard [get]
// @Security Bearer
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
