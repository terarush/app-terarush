package agent

import (
	"go-modular/internal/pkg/bus"
	"go-modular/internal/pkg/logger"
	"go-modular/modules/agent/domain/entity"
	"go-modular/modules/agent/domain/repository"
	"go-modular/modules/agent/domain/service"
	"go-modular/modules/agent/handler"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

type Module struct {
	db   *gorm.DB
	log  *logger.Logger
	evt  *bus.EventBus
	svc  *service.AgentService
	hdl  *handler.AgentHandler
}

func (m *Module) Name() string { return "agent" }

func (m *Module) Initialize(db *gorm.DB, log *logger.Logger, evt *bus.EventBus) error {
	m.db = db
	m.log = log
	m.evt = evt

	m.log.Info("Initializing agent module")

	agentRepo := repository.NewAgentRepository()
	taskRepo := repository.NewTaskRepository()
	sessionRepo := repository.NewSessionRepository()
	capRepo := repository.NewCapabilityRepository()
	logRepo := repository.NewLogRepository()
	tmplRepo := repository.NewTemplateRepository()
	schedRepo := repository.NewScheduleRepository()

	m.svc = service.NewAgentService(agentRepo, taskRepo, sessionRepo, capRepo, logRepo, tmplRepo, schedRepo)
	m.hdl = handler.NewAgentHandler(m.log, m.evt, m.svc)

	m.evt.SubscribeFunc("agent.created", m.hdl.HandleAgentEvent)
	m.evt.SubscribeFunc("agent.task.created", m.hdl.HandleAgentEvent)
	m.evt.SubscribeFunc("agent.task.completed", m.hdl.HandleAgentEvent)

	m.log.Info("Agent module initialized")
	return nil
}

func (m *Module) RegisterRoutes(e *echo.Echo, basePath string) {
	m.hdl.RegisterRoutes(e, basePath)
}

func (m *Module) Migrations() error {
	m.log.Info("Running agent module migrations")
	return m.db.AutoMigrate(
		&entity.Agent{},
		&entity.AgentTask{},
		&entity.AgentSession{},
		&entity.AgentCapability{},
		&entity.AgentLog{},
		&entity.AgentTemplate{},
		&entity.AgentSchedule{},
	)
}

func (m *Module) Logger() *logger.Logger { return m.log }

func NewModule() *Module { return &Module{} }
