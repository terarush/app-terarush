package notifications

import (
	"go-modular/internal/pkg/bus"
	"go-modular/internal/pkg/logger"
	"go-modular/modules/notifications/domain/entity"
	"go-modular/modules/notifications/domain/repository"
	"go-modular/modules/notifications/domain/service"
	"go-modular/modules/notifications/handler"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

type Module struct {
	db  *gorm.DB
	log *logger.Logger
	evt *bus.EventBus
	svc *service.NotificationService
	hdl *handler.NotificationHandler
}

func (m *Module) Name() string { return "notifications" }

func (m *Module) Initialize(db *gorm.DB, log *logger.Logger, evt *bus.EventBus) error {
	m.db = db; m.log = log; m.evt = evt
	nr := repository.NewNotificationRepository()
	tr := repository.NewNotificationTemplateRepository()
	pr := repository.NewNotificationPreferenceRepository()
	m.svc = service.NewNotificationService(nr, tr, pr)
	m.hdl = handler.NewNotificationHandler(m.log, m.evt, m.svc)
	return nil
}

func (m *Module) RegisterRoutes(e *echo.Echo, basePath string) { m.hdl.RegisterRoutes(e, basePath) }

func (m *Module) Migrations() error {
	return m.db.AutoMigrate(&entity.Notification{}, &entity.NotificationTemplate{}, &entity.NotificationPreference{})
}

func (m *Module) Logger() *logger.Logger { return m.log }

func NewModule() *Module { return &Module{} }
