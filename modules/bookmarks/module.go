package bookmarks

import (
	"go-modular/internal/pkg/bus"
	"go-modular/internal/pkg/logger"
	"go-modular/modules/bookmarks/domain/entity"
	"go-modular/modules/bookmarks/domain/repository"
	"go-modular/modules/bookmarks/domain/service"
	"go-modular/modules/bookmarks/handler"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

type Module struct {
	db  *gorm.DB
	log *logger.Logger
	evt *bus.EventBus
	svc *service.BookmarkService
	hdl *handler.BookmarkHandler
}

func (m *Module) Name() string { return "bookmarks" }

func (m *Module) Initialize(db *gorm.DB, log *logger.Logger, evt *bus.EventBus) error {
	m.db = db; m.log = log; m.evt = evt
	br := repository.NewBookmarkRepository()
	cr := repository.NewCollectionRepository()
	gr := repository.NewReadingGoalRepository()
	m.svc = service.NewBookmarkService(br, cr, gr)
	m.hdl = handler.NewBookmarkHandler(m.log, m.evt, m.svc)
	return nil
}

func (m *Module) RegisterRoutes(e *echo.Echo, basePath string) { m.hdl.RegisterRoutes(e, basePath) }

func (m *Module) Migrations() error {
	return m.db.AutoMigrate(&entity.Bookmark{}, &entity.Collection{}, &entity.ReadingGoal{})
}

func (m *Module) Logger() *logger.Logger { return m.log }

func NewModule() *Module { return &Module{} }
