package handler

import (
	"go-modular/internal/pkg/bus"
	"go-modular/internal/pkg/logger"
	"go-modular/internal/pkg/middleware"
	"go-modular/internal/pkg/utils"
	"go-modular/modules/favorites/domain/service"

	"github.com/labstack/echo/v4"
)

type FavoriteHandler struct {
	favoriteService service.FavoriteService
	log             *logger.Logger
	event           *bus.EventBus
	r               *utils.Response
}

func NewFavoriteHandler(log *logger.Logger, event *bus.EventBus, favoriteService service.FavoriteService) *FavoriteHandler {
	return &FavoriteHandler{
		favoriteService: favoriteService,
		log:             log,
		event:           event,
		r:               nil,
	}
}

func (h *FavoriteHandler) CreateFavorite(c echo.Context) error {
	return nil
}

func (h *FavoriteHandler) DeleteFavorite(c echo.Context) error {
	return nil
}

func (h *FavoriteHandler) GetByUserId(c echo.Context) error {
	return nil
}

func (h *FavoriteHandler) RegisterRoutes(e *echo.Echo, basePath string) {
	favorite := e.Group(basePath+"/favorite", middleware.Auth)
	favorite.GET("/:userId", h.GetByUserId)
	favorite.POST("", h.CreateFavorite)
	favorite.DELETE("/:favoriteId", h.DeleteFavorite)
}
