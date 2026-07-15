package moderation

import (
	"go-modular/internal/pkg/bus"
	"go-modular/internal/pkg/logger"
	"go-modular/modules/moderation/domain/entity"
	"go-modular/modules/moderation/domain/repository"
	"go-modular/modules/moderation/domain/service"
	"go-modular/modules/moderation/handler"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

type Module struct {
	db  *gorm.DB
	log *logger.Logger
	evt *bus.EventBus
	svc *service.ModerationService
	hdl *handler.ModerationHandler
}

func (m *Module) Name() string { return "moderation" }

func (m *Module) Initialize(db *gorm.DB, log *logger.Logger, evt *bus.EventBus) error {
	m.db = db; m.log = log; m.evt = evt
	rr := repository.NewReportRepository()
	ar := repository.NewModerationActionRepository()
	br := repository.NewBannedUserRepository()
	fr := repository.NewContentFilterRepository()
	m.svc = service.NewModerationService(rr, ar, br, fr)
	m.hdl = handler.NewModerationHandler(m.log, m.evt, m.svc)
	return nil
}

func (m *Module) RegisterRoutes(e *echo.Echo, basePath string) { m.hdl.RegisterRoutes(e, basePath) }

func (m *Module) Migrations() error {
	return m.db.AutoMigrate(&entity.Report{}, &entity.ModerationAction{}, &entity.BannedUser{}, &entity.ContentFilter{})
}

func (m *Module) Logger() *logger.Logger { return m.log }

func NewModule() *Module { return &Module{} }
