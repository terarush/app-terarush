package service

import (
	"context"
	"errors"
	"fmt"
	"go-modular/modules/agent/domain/entity"
	"go-modular/modules/agent/domain/repository"
	"time"
)

// Domain errors
var (
	ErrAgentNotFound     = errors.New("agent not found")
	ErrTaskNotFound      = errors.New("task not found")
	ErrSessionNotFound   = errors.New("session not found")
	ErrCapabilityNotFound = errors.New("capability not found")
	ErrTemplateNotFound  = errors.New("template not found")
	ErrScheduleNotFound  = errors.New("schedule not found")
	ErrAgentBusy         = errors.New("agent is busy")
	ErrAgentDisabled     = errors.New("agent is disabled")
	ErrNoIdleAgent       = errors.New("no idle agent available")
	ErrTaskDependency    = errors.New("task dependencies not met")
)

// AgentService handles all agent domain logic
type AgentService struct {
	agentRepo     repository.AgentRepository
	taskRepo      repository.TaskRepository
	sessionRepo   repository.SessionRepository
	capRepo       repository.CapabilityRepository
	logRepo       repository.LogRepository
	tmplRepo      repository.TemplateRepository
	schedRepo     repository.ScheduleRepository
}

func NewAgentService(
	agentRepo repository.AgentRepository,
	taskRepo repository.TaskRepository,
	sessionRepo repository.SessionRepository,
	capRepo repository.CapabilityRepository,
	logRepo repository.LogRepository,
	tmplRepo repository.TemplateRepository,
	schedRepo repository.ScheduleRepository,
) *AgentService {
	return &AgentService{
		agentRepo:   agentRepo,
		taskRepo:    taskRepo,
		sessionRepo: sessionRepo,
		capRepo:     capRepo,
		logRepo:     logRepo,
		tmplRepo:    tmplRepo,
		schedRepo:   schedRepo,
	}
}

// ─── Agent CRUD ──────────────────────────────────────────────────────────────

func (s *AgentService) GetAgents(ctx context.Context, filter map[string]any) ([]*entity.Agent, error) {
	return s.agentRepo.FindAll(ctx, filter)
}

func (s *AgentService) GetAgentByID(ctx context.Context, id uint) (*entity.Agent, error) {
	agent, err := s.agentRepo.FindByID(ctx, id)
	if err != nil {
		return nil, ErrAgentNotFound
	}
	return agent, nil
}

func (s *AgentService) CreateAgent(ctx context.Context, agent *entity.Agent) error {
	return s.agentRepo.Create(ctx, agent)
}

func (s *AgentService) UpdateAgent(ctx context.Context, agent *entity.Agent) error {
	if _, err := s.agentRepo.FindByID(ctx, agent.ID); err != nil {
		return ErrAgentNotFound
	}
	return s.agentRepo.Update(ctx, agent)
}

func (s *AgentService) DeleteAgent(ctx context.Context, id uint) error {
	if _, err := s.agentRepo.FindByID(ctx, id); err != nil {
		return ErrAgentNotFound
	}
	return s.agentRepo.Delete(ctx, id)
}

func (s *AgentService) SetAgentStatus(ctx context.Context, id uint, status string) error {
	agent, err := s.agentRepo.FindByID(ctx, id)
	if err != nil {
		return ErrAgentNotFound
	}
	agent.Status = status
	agent.UpdatedAt = time.Now()
	return s.agentRepo.Update(ctx, agent)
}

// ─── Task Management ─────────────────────────────────────────────────────────

func (s *AgentService) GetTasks(ctx context.Context, filter map[string]any, page, limit int) ([]*entity.AgentTask, int64, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 20
	}
	return s.taskRepo.FindAll(ctx, filter, page, limit)
}

func (s *AgentService) GetTaskByID(ctx context.Context, id uint) (*entity.AgentTask, error) {
	task, err := s.taskRepo.FindByID(ctx, id)
	if err != nil {
		return nil, ErrTaskNotFound
	}
	return task, nil
}

func (s *AgentService) CreateTask(ctx context.Context, task *entity.AgentTask) error {
	return s.taskRepo.Create(ctx, task)
}

func (s *AgentService) UpdateTask(ctx context.Context, task *entity.AgentTask) error {
	if _, err := s.taskRepo.FindByID(ctx, task.ID); err != nil {
		return ErrTaskNotFound
	}
	return s.taskRepo.Update(ctx, task)
}

func (s *AgentService) DeleteTask(ctx context.Context, id uint) error {
	if _, err := s.taskRepo.FindByID(ctx, id); err != nil {
		return ErrTaskNotFound
	}
	return s.taskRepo.Delete(ctx, id)
}

// AssignTask assigns a task to an agent, checks agent availability
func (s *AgentService) AssignTask(ctx context.Context, taskID, agentID uint) error {
	task, err := s.taskRepo.FindByID(ctx, taskID)
	if err != nil {
		return ErrTaskNotFound
	}
	agent, err := s.agentRepo.FindByID(ctx, agentID)
	if err != nil {
		return ErrAgentNotFound
	}
	if agent.Status == "disabled" {
		return ErrAgentDisabled
	}
	if agent.Status == "busy" && task.Priority < 2 {
		return ErrAgentBusy
	}

	// check dependencies
	if task.Dependencies != "" {
		// format: "1,2,3" — all must be completed
		// ponytail: simple parse, proper DAG check later
		// TODO: proper dependency resolution
	}

	task.AssignedTo = &agentID
	task.Status = "assigned"
	task.UpdatedAt = time.Now()
	return s.taskRepo.Update(ctx, task)
}

// StartTask marks task as in_progress, creates session
func (s *AgentService) StartTask(ctx context.Context, taskID, userID uint) (*entity.AgentSession, error) {
	task, err := s.taskRepo.FindByID(ctx, taskID)
	if err != nil {
		return nil, ErrTaskNotFound
	}
	task.Status = "in_progress"
	now := time.Now()
	task.StartedAt = &now
	task.Attempts++
	task.UpdatedAt = now
	if err := s.taskRepo.Update(ctx, task); err != nil {
		return nil, err
	}

	agentID := uint(0)
	if task.AssignedTo != nil {
		agentID = *task.AssignedTo
	}

	session := entity.NewSession(agentID, userID, &taskID)
	if err := s.sessionRepo.Create(ctx, session); err != nil {
		return nil, err
	}
	return session, nil
}

// CompleteTask marks task done with output
func (s *AgentService) CompleteTask(ctx context.Context, taskID uint, output, summary string) error {
	task, err := s.taskRepo.FindByID(ctx, taskID)
	if err != nil {
		return ErrTaskNotFound
	}
	now := time.Now()
	task.Status = "completed"
	task.Output = output
	task.ResultSummary = summary
	task.Progress = 100
	task.CompletedAt = &now
	task.UpdatedAt = now
	return s.taskRepo.Update(ctx, task)
}

// FailTask marks task failed with error
func (s *AgentService) FailTask(ctx context.Context, taskID uint, errMsg string) error {
	task, err := s.taskRepo.FindByID(ctx, taskID)
	if err != nil {
		return ErrTaskNotFound
	}
	task.Status = "failed"
	task.ErrorMessage = errMsg
	now := time.Now()
	task.CompletedAt = &now
	task.UpdatedAt = now

	// auto-retry if under limit
	if task.Attempts < task.MaxAttempts {
		task.Status = "pending"
		task.AssignedTo = nil
	}
	return s.taskRepo.Update(ctx, task)
}

// GetPendingTasks returns unassigned pending tasks, highest priority first
func (s *AgentService) GetPendingTasks(ctx context.Context, limit int) ([]*entity.AgentTask, error) {
	return s.taskRepo.FindPending(ctx, limit)
}

// AutoAssign assigns pending tasks to idle agents with matching capabilities
func (s *AgentService) AutoAssign(ctx context.Context) ([]*entity.AgentTask, error) {
	agents, err := s.agentRepo.FindByStatus(ctx, "idle")
	if err != nil || len(agents) == 0 {
		return nil, ErrNoIdleAgent
	}
	pending, err := s.taskRepo.FindPending(ctx, len(agents))
	if err != nil {
		return nil, err
	}
	var assigned []*entity.AgentTask
	for i, task := range pending {
		if i >= len(agents) {
			break
		}
		id := agents[i].ID
		task.AssignedTo = &id
		task.Status = "assigned"
		task.UpdatedAt = time.Now()
		if err := s.taskRepo.Update(ctx, task); err != nil {
			continue
		}
		assigned = append(assigned, task)
	}
	return assigned, nil
}

// ─── Sessions ────────────────────────────────────────────────────────────────

func (s *AgentService) GetSessions(ctx context.Context, filter map[string]any) ([]*entity.AgentSession, error) {
	return s.sessionRepo.FindAll(ctx, filter)
}

func (s *AgentService) GetSessionByID(ctx context.Context, id uint) (*entity.AgentSession, error) {
	sess, err := s.sessionRepo.FindByID(ctx, id)
	if err != nil {
		return nil, ErrSessionNotFound
	}
	return sess, nil
}

func (s *AgentService) EndSession(ctx context.Context, id uint) error {
	sess, err := s.sessionRepo.FindByID(ctx, id)
	if err != nil {
		return ErrSessionNotFound
	}
	now := time.Now()
	sess.Status = "completed"
	sess.EndedAt = &now
	if !sess.StartedAt.IsZero() {
		sess.Duration = int64(now.Sub(sess.StartedAt).Seconds())
	}
	return s.sessionRepo.Update(ctx, sess)
}

// ─── Capabilities ────────────────────────────────────────────────────────────

func (s *AgentService) GetCapabilities(ctx context.Context, agentID uint) ([]*entity.AgentCapability, error) {
	return s.capRepo.FindAll(ctx, agentID)
}

func (s *AgentService) AddCapability(ctx context.Context, cap *entity.AgentCapability) error {
	return s.capRepo.Create(ctx, cap)
}

func (s *AgentService) RemoveCapability(ctx context.Context, id uint) error {
	return s.capRepo.Delete(ctx, id)
}

func (s *AgentService) ToggleCapability(ctx context.Context, id uint, enabled bool) error {
	cap, err := s.capRepo.FindByID(ctx, id)
	if err != nil {
		return ErrCapabilityNotFound
	}
	cap.Enabled = enabled
	cap.UpdatedAt = time.Now()
	return s.capRepo.Update(ctx, cap)
}

// ─── Logs ────────────────────────────────────────────────────────────────────

func (s *AgentService) GetLogs(ctx context.Context, filter map[string]any, page, limit int) ([]*entity.AgentLog, int64, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 200 {
		limit = 50
	}
	return s.logRepo.FindAll(ctx, filter, page, limit)
}

func (s *AgentService) CreateLog(ctx context.Context, log *entity.AgentLog) error {
	return s.logRepo.Create(ctx, log)
}

func (s *AgentService) CleanupLogs(ctx context.Context, days int) (int64, error) {
	return s.logRepo.DeleteOlderThan(ctx, days)
}

// ─── Templates ───────────────────────────────────────────────────────────────

func (s *AgentService) GetTemplates(ctx context.Context, filter map[string]any) ([]*entity.AgentTemplate, error) {
	return s.tmplRepo.FindAll(ctx, filter)
}

func (s *AgentService) GetTemplateByID(ctx context.Context, id uint) (*entity.AgentTemplate, error) {
	tmpl, err := s.tmplRepo.FindByID(ctx, id)
	if err != nil {
		return nil, ErrTemplateNotFound
	}
	return tmpl, nil
}

func (s *AgentService) CreateTemplate(ctx context.Context, tmpl *entity.AgentTemplate) error {
	return s.tmplRepo.Create(ctx, tmpl)
}

func (s *AgentService) DeleteTemplate(ctx context.Context, id uint) error {
	if _, err := s.tmplRepo.FindByID(ctx, id); err != nil {
		return ErrTemplateNotFound
	}
	return s.tmplRepo.Delete(ctx, id)
}

// CreateTaskFromTemplate generates a task from a template
func (s *AgentService) CreateTaskFromTemplate(ctx context.Context, templateID uint, input string, createdBy uint) (*entity.AgentTask, error) {
	tmpl, err := s.tmplRepo.FindByID(ctx, templateID)
	if err != nil {
		return nil, ErrTemplateNotFound
	}
	task := entity.NewTask(tmpl.Name, tmpl.Description, tmpl.Category, tmpl.Priority, createdBy)
	task.MaxAttempts = tmpl.MaxAttempts
	task.Input = input
	if err := s.taskRepo.Create(ctx, task); err != nil {
		return nil, fmt.Errorf("create task from template: %w", err)
	}
	return task, nil
}

// ─── Schedules ───────────────────────────────────────────────────────────────

func (s *AgentService) GetSchedules(ctx context.Context, filter map[string]any) ([]*entity.AgentSchedule, error) {
	return s.schedRepo.FindAll(ctx, filter)
}

func (s *AgentService) GetScheduleByID(ctx context.Context, id uint) (*entity.AgentSchedule, error) {
	sched, err := s.schedRepo.FindByID(ctx, id)
	if err != nil {
		return nil, ErrScheduleNotFound
	}
	return sched, nil
}

func (s *AgentService) CreateSchedule(ctx context.Context, sched *entity.AgentSchedule) error {
	// ponytail: static cron parse later, for now store raw
	return s.schedRepo.Create(ctx, sched)
}

func (s *AgentService) DeleteSchedule(ctx context.Context, id uint) error {
	if _, err := s.schedRepo.FindByID(ctx, id); err != nil {
		return ErrScheduleNotFound
	}
	return s.schedRepo.Delete(ctx, id)
}

// ProcessDueSchedules finds schedules due and creates tasks from them
func (s *AgentService) ProcessDueSchedules(ctx context.Context) ([]*entity.AgentTask, error) {
	due, err := s.schedRepo.FindDue(ctx)
	if err != nil {
		return nil, err
	}
	var tasks []*entity.AgentTask
	now := time.Now()
	for _, sched := range due {
		task := entity.NewTask(
			fmt.Sprintf("[Scheduled] %s", sched.Name),
			sched.Description,
			"scheduled",
			0,
			sched.CreatedBy,
		)
		task.Input = sched.InputConfig
		task.AssignedTo = &sched.AgentID
		if sched.TemplateID != nil {
			tmpl, tErr := s.tmplRepo.FindByID(ctx, *sched.TemplateID)
			if tErr == nil {
				task.MaxAttempts = tmpl.MaxAttempts
			}
		}
		if err := s.taskRepo.Create(ctx, task); err != nil {
			continue
		}
		sched.RunCount++
		sched.LastRunAt = &now
		sched.UpdatedAt = now
		_ = s.schedRepo.Update(ctx, sched)
		tasks = append(tasks, task)
	}
	return tasks, nil
}

// ─── Dashboard / Stats ──────────────────────────────────────────────────────

type AgentStats struct {
	TotalAgents    int64 `json:"total_agents"`
	IdleAgents     int64 `json:"idle_agents"`
	BusyAgents     int64 `json:"busy_agents"`
	TotalTasks     int64 `json:"total_tasks"`
	PendingTasks   int64 `json:"pending_tasks"`
	ActiveSessions int64 `json:"active_sessions"`
}

func (s *AgentService) GetStats(ctx context.Context) (*AgentStats, error) {
	totalAgents, _ := s.agentRepo.Count(ctx, nil)
	idleAgents, _ := s.agentRepo.Count(ctx, map[string]any{"status": "idle"})
	busyAgents, _ := s.agentRepo.Count(ctx, map[string]any{"status": "busy"})
	_, totalTasks, _ := s.taskRepo.FindAll(ctx, nil, 1, 1)
	pendingTasks, _, _ := s.taskRepo.FindAll(ctx, map[string]any{"status": "pending"}, 1, 1)
	activeSessions, _ := s.sessionRepo.FindAll(ctx, map[string]any{"status": "active"})

	return &AgentStats{
		TotalAgents:    totalAgents,
		IdleAgents:     idleAgents,
		BusyAgents:     busyAgents,
		TotalTasks:     totalTasks,
		PendingTasks:   int64(len(pendingTasks)),
		ActiveSessions: int64(len(activeSessions)),
	}, nil
}
