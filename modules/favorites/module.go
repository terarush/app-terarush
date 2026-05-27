package favorites

import (
	"go-modular/internal/pkg/bus"
	"go-modular/internal/pkg/logger"
	"go-modular/modules/favorites/domain/service"
	"go-modular/modules/favorites/handler"

	"gorm.io/gorm"
)

type Module struct {
	db              *gorm.DB
	logger          *logger.Logger
	favoriteService service.FavoriteService
	favoriteHandler *handler.FavoriteHandler
	event           *bus.EventBus
}
