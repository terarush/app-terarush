package products

import (
	"go-modular/internal/pkg/bus"
	"go-modular/internal/pkg/logger"
	"go-modular/internal/pkg/middleware"
	"go-modular/modules/products/domain/entity"
	"go-modular/modules/products/domain/repository"
	"go-modular/modules/products/domain/service"
	"go-modular/modules/products/handler"

	"github.com/labstack/echo"
	"gorm.io/gorm"
)

type Module struct {
	db             *gorm.DB
	logger         *logger.Logger
	productService *service.ProductService
	productHandler *handler.ProductHandler
	event          *bus.EventBus
}

func (m *Module) Name() string {
	return "products"
}

func (m *Module) Initialize(db *gorm.DB, log *logger.Logger, event *bus.EventBus) error {
	m.db = db
	m.logger = log
	m.event = event

	// Initialize repositories
	productRepo := repository.NewProductRepository(db)

	// Initialize services
	m.productService = service.NewProductService(productRepo)

	// Initialize handlers
	m.productHandler = handler.NewProductHandler(m.productService, m.logger)

	m.logger.Info("Products module initialized successfully")
	return nil
}

func (m *Module) RegisterRoutes(e *echo.Echo, basePath string) {
	if m.productHandler == nil {
		m.logger.Error("ProductHandler is nil, cannot register routes")
		return
	}

	// Get middleware
	authMiddleware := middleware.Auth
	adminMiddleware := middleware.RequireRole("admin")

	m.productHandler.RegisterRoutes(e, basePath, authMiddleware, adminMiddleware)
}

func (m *Module) Migrations() error {
	m.logger.Info("Registerin products module migrations")
	return m.db.AutoMigrate(&entity.Product{})
}

func (m *Module) Logger() *logger.Logger {
	return m.logger
}

func NewModule() *Module {
	return &Module{}
}
