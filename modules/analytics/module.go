package analytics

import (
	"go-modular/internal/pkg/bus"
	"go-modular/internal/pkg/logger"
	"go-modular/modules/analytics/domain/entity"
	"go-modular/modules/analytics/domain/repository"
	"go-modular/modules/analytics/domain/service"
	"go-modular/modules/analytics/handler"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

type Module struct {
	db  *gorm.DB
	log *logger.Logger
	evt *bus.EventBus
	svc *service.AnalyticsService
	hdl *handler.AnalyticsHandler
}

func (m *Module) Name() string { return "analytics" }

func (m *Module) Initialize(db *gorm.DB, log *logger.Logger, evt *bus.EventBus) error {
	m.db = db; m.log = log; m.evt = evt
	er := repository.NewAnalyticsEventRepository()
	br := repository.NewBlogStatsRepository()
	ur := repository.NewUserAnalyticsRepository()
	m.svc = service.NewAnalyticsService(er, br, ur)
	m.hdl = handler.NewAnalyticsHandler(m.log, m.evt, m.svc)
	return nil
}

func (m *Module) RegisterRoutes(e *echo.Echo, basePath string) { m.hdl.RegisterRoutes(e, basePath) }

func (m *Module) Migrations() error {
	return m.db.AutoMigrate(&entity.AnalyticsEvent{}, &entity.AnalyticsAggregate{}, &entity.BlogStats{}, &entity.UserAnalytics{})
}

func (m *Module) Logger() *logger.Logger { return m.log }

func NewModule() *Module { return &Module{} }
