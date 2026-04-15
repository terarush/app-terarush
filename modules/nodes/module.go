package nodes

import (
	"go-modular/internal/pkg/bus"
	"go-modular/internal/pkg/logger"
	"go-modular/internal/pkg/middleware"
	"go-modular/internal/pkg/validator"
	"go-modular/modules/nodes/domain/entity"
	"go-modular/modules/nodes/domain/repository"
	"go-modular/modules/nodes/domain/service"
	"go-modular/modules/nodes/handler"
	productRepo "go-modular/modules/products/domain/repository"

	"github.com/labstack/echo"
	"gorm.io/gorm"
)

// Module implements the application Module interface for the nodes module
type Module struct {
	db          *gorm.DB
	logger      *logger.Logger
	event       *bus.EventBus
	nodeService service.NodeService
	nodeHandler *handler.NodeHandler
}

// Name returns the name of the module
func (m *Module) Name() string {
	return "nodes"
}

// Initialize initializes the module
func (m *Module) Initialize(db *gorm.DB, log *logger.Logger, event *bus.EventBus) error {
	m.db = db
	m.logger = log
	m.event = event

	m.logger.Info("Initializing nodes module")

	// Get validator from echo context
	v := validator.NewCustomValidator()

	// Initialize repositories
	nodeRepo := repository.NewNodeRepository(db)
	productRepository := productRepo.NewProductRepository(db)
	m.logger.Debug("Node repository initialized")

	// Initialize services
	nodeService, err := service.NewNodeService(nodeRepo, productRepository)
	if err != nil {
		m.logger.Error("Failed to initialize node service", "error", err)
		return err
	}
	m.nodeService = nodeService
	m.logger.Debug("Node service initialized")

	// Initialize handlers
	m.nodeHandler = handler.NewNodeHandler(m.logger, v, m.nodeService)
	m.logger.Debug("Node handler initialized")

	m.logger.Info("Nodes module initialized successfully")
	return nil
}

// RegisterRoutes registers the module's routes
func (m *Module) RegisterRoutes(e *echo.Echo, basePath string) {
	m.logger.Info("Registering node routes at %s/nodes", basePath)

	group := e.Group(basePath+"/nodes", middleware.Auth)

	// Node CRUD operations
	group.POST("", m.nodeHandler.CreateNode)
	group.GET("", m.nodeHandler.GetNodes)
	group.GET("/:id", m.nodeHandler.GetNode)
	group.PUT("/:id", m.nodeHandler.UpdateNode)
	group.DELETE("/:id", m.nodeHandler.DeleteNode)

	// Node actions (start, stop, restart)
	group.POST("/:id/action", m.nodeHandler.NodeAction)

	m.logger.Debug("Node routes registered successfully")
}

// Migrations returns the module's migrations
func (m *Module) Migrations() error {
	m.logger.Info("Registering nodes module migrations")
	return m.db.AutoMigrate(&entity.Node{})
}

// Logger returns the module's logger
func (m *Module) Logger() *logger.Logger {
	return m.logger
}

// NewModule creates a new nodes module
func NewModule() *Module {
	return &Module{}
}
