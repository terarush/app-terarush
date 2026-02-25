package transactions

import (
	"go-modular/internal/pkg/bus"
	"go-modular/internal/pkg/logger"
	"go-modular/internal/pkg/middleware"
	productRepo "go-modular/modules/products/domain/repository"
	"go-modular/modules/transactions/domain/entity"
	"go-modular/modules/transactions/domain/repository"
	"go-modular/modules/transactions/domain/service"
	"go-modular/modules/transactions/handler"

	"github.com/labstack/echo"
	"gorm.io/gorm"
)

type Module struct {
	db                 *gorm.DB
	logger             *logger.Logger
	transactionService *service.TransactionService
	transactionHandler *handler.TransactionHandler
	event              *bus.EventBus
}

func (m *Module) Name() string {
	return "transactions"
}

func (m *Module) Initialize(db *gorm.DB, log *logger.Logger, event *bus.EventBus) error {
	m.db = db
	m.logger = log
	m.event = event

	// Initialize repositories
	transactionRepo := repository.NewTransactionRepository(db)
	productRepository := productRepo.NewProductRepository(db)

	// Initialize services
	m.transactionService = service.NewTransactionService(transactionRepo, productRepository)

	// Initialize handlers
	m.transactionHandler = handler.NewTransactionHandler(m.transactionService, m.logger)

	m.logger.Info("Transactions module initialized successfully")
	return nil
}

func (m *Module) RegisterRoutes(e *echo.Echo, basePath string) {
	if m.transactionHandler == nil {
		m.logger.Error("TransactionHandler is nil, cannot register routes")
		return
	}

	// Get middleware
	authMiddleware := middleware.Auth
	adminMiddleware := middleware.RequireRole("admin")

	m.transactionHandler.RegisterRoutes(e, basePath, authMiddleware, adminMiddleware)
}

func (m *Module) Migrations() error {
	m.logger.Info("Registering transactions module migrations")
	return m.db.AutoMigrate(&entity.Transaction{})
}

func (m *Module) Logger() *logger.Logger {
	return m.logger
}

func NewModule() *Module {
	return &Module{}
}
