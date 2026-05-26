package comments

import (
	"go-modular/internal/pkg/bus"
	"go-modular/internal/pkg/logger"
	"go-modular/modules/comments/domain/entity"
	"go-modular/modules/comments/domain/repository"
	"go-modular/modules/comments/domain/service"
	"go-modular/modules/comments/handler"

	"github.com/labstack/echo"
	"gorm.io/gorm"
)

type Module struct {
	db             *gorm.DB
	logger         *logger.Logger
	commentService service.CommentService
	commentHandler *handler.CommentHandler
	event          *bus.EventBus
}

func (m *Module) Name() string {
	return "comments"
}

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

	m.logger.Info("Comments module initialized successfully")
	return nil
}

func (m *Module) RegisterRoutes(e *echo.Echo, basePath string) {
	m.logger.Info("Registering comments module routes")
	m.commentHandler.RegisterRoutes(e, basePath)
}

func (m *Module) Migrations() error {
	m.logger.Info("Running comments module migrations")
	return m.db.AutoMigrate(&entity.Comment{})
}

func (m *Module) Logger() *logger.Logger {
	return m.logger
}

func NewModule() *Module {
	return &Module{}
}
