package activity

import (
	"go-modular/internal/pkg/bus"
	"go-modular/internal/pkg/logger"
	"go-modular/modules/activity/domain/entity"
	"go-modular/modules/activity/domain/repository"
	"go-modular/modules/activity/domain/service"
	"go-modular/modules/activity/handler"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

type Module struct {
	db  *gorm.DB
	log *logger.Logger
	evt *bus.EventBus
	svc *service.ActivityService
	hdl *handler.ActivityHandler
}

func (m *Module) Name() string { return "activity" }

func (m *Module) Initialize(db *gorm.DB, log *logger.Logger, evt *bus.EventBus) error {
	m.db = db; m.log = log; m.evt = evt
	ar := repository.NewActivityRepository()
	fr := repository.NewActivityFeedRepository()
	flr := repository.NewFollowRepository()
	mr := repository.NewEngagementMetricRepository()
	m.svc = service.NewActivityService(ar, fr, flr, mr)
	m.hdl = handler.NewActivityHandler(m.log, m.evt, m.svc)
	return nil
}

func (m *Module) RegisterRoutes(e *echo.Echo, basePath string) { m.hdl.RegisterRoutes(e, basePath) }

func (m *Module) Migrations() error {
	return m.db.AutoMigrate(&entity.Activity{}, &entity.ActivityFeed{}, &entity.Follow{}, &entity.EngagementMetric{})
}

func (m *Module) Logger() *logger.Logger { return m.log }

func NewModule() *Module { return &Module{} }
