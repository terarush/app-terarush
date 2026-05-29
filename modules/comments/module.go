package comments

import (
	"go-modular/internal/pkg/bus"
	"go-modular/internal/pkg/logger"
	"go-modular/modules/comments/domain/entity"
	"go-modular/modules/comments/domain/repository"
	"go-modular/modules/comments/domain/service"
	"go-modular/modules/comments/handler"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

// Module implements the application Module interface for the comments module
type Module struct {
	db             *gorm.DB
	logger         *logger.Logger
	commentService service.CommentService
	commentHandler *handler.CommentHandler
	event          *bus.EventBus
}

// Name returns the name of the module
func (m *Module) Name() string {
	return "comments"
}

// Initialize initializes the module
func (m *Module) Initialize(db *gorm.DB, log *logger.Logger, event *bus.EventBus) error {
	m.db = db
	m.logger = log
	m.event = event

	m.logger.Info("Initializing comments module")

	// Initialize repositories
	commentRepo := repository.NewCommentRepository()
	m.logger.Debug("Comment repository initialized")

	// Initialize services
	m.commentService = *service.NewCommentService(commentRepo)
	m.logger.Debug("Comment service initialized")

	// Initialize handlers
	m.commentHandler = handler.NewCommentHandler(m.logger, m.event, m.commentService)
	m.logger.Debug("Comment handler initialized")

	// register event listeners
	m.logger.Info("Registering comments module event listeners")
	m.event.SubscribeFunc("comment.created", m.handleCommentCreated)
	m.event.SubscribeFunc("comment.updated", m.handleCommentUpdated)
	m.event.SubscribeFunc("comment.deleted", m.handleCommentDeleted)

	m.logger.Info("Comments module initialized successfully")
	return nil
}

// RegisterRoutes registers the module's routes
func (m *Module) RegisterRoutes(e *echo.Echo, basePath string) {
	m.logger.Info("Registering comments routes at %s", basePath)
	m.commentHandler.RegisterRoutes(e, basePath)
	m.logger.Debug("Comments routes registered successfully")
}

// Migrations returns the module's migrations
func (m *Module) Migrations() error {
	m.logger.Info("Registering comments module migrations")
	if err := m.db.AutoMigrate(&entity.Comment{}); err != nil {
		return err
	}
	// Add foreign key constraint with cascade delete
	return m.db.Migrator().CreateConstraint(&entity.Comment{}, "post_id")
}

// Logger returns the module's logger
func (m *Module) Logger() *logger.Logger {
	return m.logger
}

// Event handlers
func (m *Module) handleCommentCreated(event bus.Event) {
	m.logger.Info("Comment created event received: %v", event.Payload)
}

func (m *Module) handleCommentUpdated(event bus.Event) {
	m.logger.Info("Comment updated event received: %v", event.Payload)
}

func (m *Module) handleCommentDeleted(event bus.Event) {
	m.logger.Info("Comment deleted event received: %v", event.Payload)
}

// NewModule creates a new comments module
func NewModule() *Module {
	return &Module{}
}
