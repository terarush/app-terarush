package subscriptions

import (
	"go-modular/internal/pkg/bus"
	"go-modular/internal/pkg/logger"
	"go-modular/modules/subscriptions/domain/entity"
	"go-modular/modules/subscriptions/domain/repository"
	"go-modular/modules/subscriptions/domain/service"
	"go-modular/modules/subscriptions/handler"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

type Module struct {
	db  *gorm.DB
	log *logger.Logger
	evt *bus.EventBus
	svc *service.SubscriptionService
	hdl *handler.SubscriptionHandler
}

func (m *Module) Name() string { return "subscriptions" }

func (m *Module) Initialize(db *gorm.DB, log *logger.Logger, evt *bus.EventBus) error {
	m.db = db; m.log = log; m.evt = evt
	sr := repository.NewSubscriptionRepository()
	nr := repository.NewNewsletterRepository()
	sbr := repository.NewNewsletterSubscriberRepository()
	m.svc = service.NewSubscriptionService(sr, nr, sbr)
	m.hdl = handler.NewSubscriptionHandler(m.log, m.evt, m.svc)
	return nil
}

func (m *Module) RegisterRoutes(e *echo.Echo, basePath string) { m.hdl.RegisterRoutes(e, basePath) }

func (m *Module) Migrations() error {
	return m.db.AutoMigrate(&entity.Subscription{}, &entity.Newsletter{}, &entity.NewsletterSubscriber{})
}

func (m *Module) Logger() *logger.Logger { return m.log }

func NewModule() *Module { return &Module{} }
