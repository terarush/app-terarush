package favorites

import (
	"go-modular/internal/pkg/bus"
	"go-modular/internal/pkg/logger"
	"go-modular/modules/favorites/domain/entity"
	"go-modular/modules/favorites/domain/repository"
	"go-modular/modules/favorites/domain/service"
	"go-modular/modules/favorites/handler"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

// Module implements the application Module interface for the favorites module
type Module struct {
	db              *gorm.DB
	logger          *logger.Logger
	favoriteService service.FavoriteService
	favoriteHandler *handler.FavoriteHandler
	event           *bus.EventBus
}

// Name returns the name of the module
func (m *Module) Name() string {
	return "favorites"
}

// Initialize initializes the module
func (m *Module) Initialize(db *gorm.DB, log *logger.Logger, event *bus.EventBus) error {
	m.db = db
	m.logger = log
	m.event = event

	m.logger.Info("Initializing favorites module")

	// Initialize repositories
	favoriteRepo := repository.NewFavoriteRepositoryImpl()
	m.logger.Debug("Favorite repository initialized")

	// Initialize services
	m.favoriteService = *service.NewFavoriteService(favoriteRepo)
	m.logger.Debug("Favorite service initialized")

	// Initialize handlers
	m.favoriteHandler = handler.NewFavoriteHandler(m.logger, m.event, m.favoriteService)
	m.logger.Debug("Favorite handler initialized")

	// register event listeners
	m.logger.Info("Registering favorites module event listeners")
	m.event.SubscribeFunc("favorite.created", m.handleFavoriteCreated)
	m.event.SubscribeFunc("favorite.deleted", m.handleFavoriteDeleted)

	m.logger.Info("Favorites module initialized successfully")
	return nil
}

// RegisterRoutes registers the module's routes
func (m *Module) RegisterRoutes(e *echo.Echo, basePath string) {
	m.logger.Info("Registering favorites routes at %s", basePath)
	m.favoriteHandler.RegisterRoutes(e, basePath)
	m.logger.Debug("Favorites routes registered successfully")
}

// Migrations returns the module's migrations
func (m *Module) Migrations() error {
	m.logger.Info("Registering favorites module migrations")
	return m.db.AutoMigrate(&entity.Favorite{})
}

// Logger returns the module's logger
func (m *Module) Logger() *logger.Logger {
	return m.logger
}

// Event handlers
func (m *Module) handleFavoriteCreated(event bus.Event) {
	m.logger.Info("Favorite created event received: %v", event.Payload)
}

func (m *Module) handleFavoriteDeleted(event bus.Event) {
	m.logger.Info("Favorite deleted event received: %v", event.Payload)
}

// NewModule creates a new favorites module
func NewModule() *Module {
	return &Module{}
}
