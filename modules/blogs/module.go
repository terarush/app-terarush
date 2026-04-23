package blogs

import (
	"go-modular/internal/pkg/bus"
	"go-modular/internal/pkg/logger"
	"go-modular/modules/blogs/domain/entity"
	"go-modular/modules/blogs/domain/repository"
	"go-modular/modules/blogs/domain/service"
	"go-modular/modules/blogs/handler"

	"github.com/labstack/echo"
	"gorm.io/gorm"
)

// Module implements the application Module interface for the blog module
type Module struct {
	db          *gorm.DB
	logger      *logger.Logger
	blogService *service.BlogService
	blogHandler *handler.BlogHandler
	event       *bus.EventBus
}

// Name returns the name of the module
func (m *Module) Name() string {
	return "blog"
}

// Initialize initializes the module
func (m *Module) Initialize(db *gorm.DB, log *logger.Logger, event *bus.EventBus) error {
	m.db = db
	m.logger = log
	m.event = event

	m.logger.Info("Initializing blog module")

	// Initialize repositories
	blogRepo := repository.NewBlogRepositoryImpl()
	m.logger.Debug("Blog repository initialized")

	// Initialize services
	m.blogService = service.NewBlogService(blogRepo)
	m.logger.Debug("Blog service initialized")

	// Initialize handlers
	m.blogHandler = handler.NewBlogHandler(m.logger, m.event, m.blogService)
	m.logger.Debug("Blog handler initialized")

	// register event listeners
	m.logger.Info("Registering blog module event listeners")
	m.event.SubscribeFunc("blog.created", m.handleBlogCreated)
	m.event.SubscribeFunc("blog.updated", m.handleBlogUpdated)
	m.event.SubscribeFunc("blog.deleted", m.handleBlogDeleted)

	m.logger.Info("Blog module initialized successfully")
	return nil
}

// RegisterRoutes registers the module's routes
func (m *Module) RegisterRoutes(e *echo.Echo, basePath string) {
	m.logger.Info("Registering blog routes at %s", basePath)
	m.blogHandler.RegisterRoutes(e, basePath)
	m.logger.Debug("Blog routes registered successfully")
}

// Migrations returns the module's migrations
func (m *Module) Migrations() error {
	m.logger.Info("Registering blog module migrations")
	return m.db.AutoMigrate(&entity.Blog{})
}

// Logger returns the module's logger
func (m *Module) Logger() *logger.Logger {
	return m.logger
}

// Event handlers
func (m *Module) handleBlogCreated(event bus.Event) {
	m.logger.Info("Blog created event received: %v", event.Payload)
}

func (m *Module) handleBlogUpdated(event bus.Event) {
	m.logger.Info("Blog updated event received: %v", event.Payload)
}

func (m *Module) handleBlogDeleted(event bus.Event) {
	m.logger.Info("Blog deleted event received: %v", event.Payload)
}

// NewModule creates a new blog module
func NewModule() *Module {
	return &Module{}
}
