package assets

import (
	"go-modular/internal/pkg/bus"
	"go-modular/internal/pkg/logger"
	"go-modular/modules/assets/domain/repository"
	"go-modular/modules/assets/domain/service"
	"go-modular/modules/assets/handler"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

// Module is the assets module
type Module struct {
	handler *handler.AssetHandler
	log     *logger.Logger
}

// Initialize initializes the assets module
func (m *Module) Initialize(db *gorm.DB, log *logger.Logger, event *bus.EventBus) error {
	m.log = log
	m.log.Info("Initializing assets module")

	// Initialize repository
	assetRepo := repository.NewAssetRepositoryImpl()

	// Initialize service
	assetService := service.NewAssetService(assetRepo)

	// Initialize handler
	m.handler = handler.NewAssetHandler(log, assetService)

	m.log.Info("Assets module initialized successfully")
	return nil
}

// RegisterRoutes registers the routes for the assets module
func (m *Module) RegisterRoutes(e *echo.Echo, basePath string) {
	m.log.Info("Registering assets routes at %s", basePath)
	m.handler.RegisterRoutes(e, basePath)
}

// Migrations returns the migrations for the assets module
func (m *Module) Migrations() error {
	// No migrations needed for assets module - using file storage
	return nil
}

// Name returns the name of the module
func (m *Module) Name() string {
	return "assets"
}

// Logger returns the module's logger
func (m *Module) Logger() *logger.Logger {
	return m.log
}

// NewModule creates a new assets module
func NewModule() *Module {
	return &Module{}
}
