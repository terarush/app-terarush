package response

import (
	"go-modular/modules/agent/domain/entity"
	"time"
)

// ─── Agent ───────────────────────────────────────────────────────────────────
type AgentResponse struct {
	ID            uint      `json:"id"`
	Name          string    `json:"name"`
	Description   string    `json:"description,omitempty"`
	Model         string    `json:"model"`
	SystemPrompt  string    `json:"system_prompt,omitempty"`
	Status        string    `json:"status"`
	MaxConcurrent int       `json:"max_concurrent"`
	Timeout       int64     `json:"timeout"`
	RetryLimit    int       `json:"retry_limit"`
	CreatedBy     uint      `json:"created_by"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

type AgentListResponse struct {
	Agents []*AgentResponse `json:"agents"`
	Total  int64            `json:"total"`
}

func AgentFromEntity(a *entity.Agent) *AgentResponse {
	return &AgentResponse{
		ID:            a.ID,
		Name:          a.Name,
		Description:   a.Description,
		Model:         a.Model,
		SystemPrompt:  a.SystemPrompt,
		Status:        a.Status,
		MaxConcurrent: a.MaxConcurrent,
		Timeout:       a.Timeout,
		RetryLimit:    a.RetryLimit,
		CreatedBy:     a.CreatedBy,
		CreatedAt:     a.CreatedAt,
		UpdatedAt:     a.UpdatedAt,
	}
}

func AgentFromEntities(agents []*entity.Agent) []*AgentResponse {
	res := make([]*AgentResponse, len(agents))
	for i, a := range agents {
		res[i] = AgentFromEntity(a)
	}
	return res
}

// ─── Task ────────────────────────────────────────────────────────────────────
type TaskResponse struct {
	ID            uint       `json:"id"`
	Title         string     `json:"title"`
	Description   string     `json:"description,omitempty"`
	Status        string     `json:"status"`
	Priority      int        `json:"priority"`
	Category      string     `json:"category,omitempty"`
	Tags          string     `json:"tags,omitempty"`
	AssignedTo    *uint      `json:"assigned_to,omitempty"`
	AssignedBy    *uint      `json:"assigned_by,omitempty"`
	ParentTaskID  *uint      `json:"parent_task_id,omitempty"`
	ResultSummary string     `json:"result_summary,omitempty"`
	ErrorMessage  string     `json:"error_message,omitempty"`
	Attempts      int        `json:"attempts"`
	MaxAttempts   int        `json:"max_attempts"`
	Progress      float64    `json:"progress"`
	Dependencies  string     `json:"dependencies,omitempty"`
	StartedAt     *time.Time `json:"started_at,omitempty"`
	CompletedAt   *time.Time `json:"completed_at,omitempty"`
	Deadline      *time.Time `json:"deadline,omitempty"`
	CreatedBy     uint       `json:"created_by"`
	CreatedAt     time.Time  `json:"created_at"`
	UpdatedAt     time.Time  `json:"updated_at"`
}

type TaskListResponse struct {
	Tasks []*TaskResponse `json:"tasks"`
	Total int64           `json:"total"`
	Page  int             `json:"page"`
	Limit int             `json:"limit"`
}

func TaskFromEntity(t *entity.AgentTask) *TaskResponse {
	return &TaskResponse{
		ID:            t.ID,
		Title:         t.Title,
		Description:   t.Description,
		Status:        t.Status,
		Priority:      t.Priority,
		Category:      t.Category,
		Tags:          t.Tags,
		AssignedTo:    t.AssignedTo,
		AssignedBy:    t.AssignedBy,
		ParentTaskID:  t.ParentTaskID,
		ResultSummary: t.ResultSummary,
		ErrorMessage:  t.ErrorMessage,
		Attempts:      t.Attempts,
		MaxAttempts:   t.MaxAttempts,
		Progress:      t.Progress,
		Dependencies:  t.Dependencies,
		StartedAt:     t.StartedAt,
		CompletedAt:   t.CompletedAt,
		Deadline:      t.Deadline,
		CreatedBy:     t.CreatedBy,
		CreatedAt:     t.CreatedAt,
		UpdatedAt:     t.UpdatedAt,
	}
}

func TaskFromEntities(tasks []*entity.AgentTask) []*TaskResponse {
	res := make([]*TaskResponse, len(tasks))
	for i, t := range tasks {
		res[i] = TaskFromEntity(t)
	}
	return res
}

// ─── Session ─────────────────────────────────────────────────────────────────
type SessionResponse struct {
	ID        uint       `json:"id"`
	AgentID   uint       `json:"agent_id"`
	TaskID    *uint      `json:"task_id,omitempty"`
	Status    string     `json:"status"`
	UserID    uint       `json:"user_id"`
	TokenUsed int64      `json:"token_used"`
	Duration  int64      `json:"duration"`
	Messages  int        `json:"messages"`
	StartedAt time.Time  `json:"started_at"`
	EndedAt   *time.Time `json:"ended_at,omitempty"`
	CreatedAt time.Time  `json:"created_at"`
}

func SessionFromEntity(s *entity.AgentSession) *SessionResponse {
	return &SessionResponse{
		ID:        s.ID,
		AgentID:   s.AgentID,
		TaskID:    s.TaskID,
		Status:    s.Status,
		UserID:    s.UserID,
		TokenUsed: s.TokenUsed,
		Duration:  s.Duration,
		Messages:  s.Messages,
		StartedAt: s.StartedAt,
		EndedAt:   s.EndedAt,
		CreatedAt: s.CreatedAt,
	}
}

func SessionFromEntities(sessions []*entity.AgentSession) []*SessionResponse {
	res := make([]*SessionResponse, len(sessions))
	for i, s := range sessions {
		res[i] = SessionFromEntity(s)
	}
	return res
}

// ─── Capability ──────────────────────────────────────────────────────────────
type CapabilityResponse struct {
	ID          uint      `json:"id"`
	AgentID     uint      `json:"agent_id"`
	Capability  string    `json:"capability"`
	Enabled     bool      `json:"enabled"`
	Config      string    `json:"config,omitempty"`
	MaxUsage    int       `json:"max_usage"`
	UsageCount  int       `json:"usage_count"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

func CapabilityFromEntity(c *entity.AgentCapability) *CapabilityResponse {
	return &CapabilityResponse{
		ID:         c.ID,
		AgentID:    c.AgentID,
		Capability: c.Capability,
		Enabled:    c.Enabled,
		Config:     c.Config,
		MaxUsage:   c.MaxUsage,
		UsageCount: c.UsageCount,
		CreatedAt:  c.CreatedAt,
		UpdatedAt:  c.UpdatedAt,
	}
}

func CapabilityFromEntities(caps []*entity.AgentCapability) []*CapabilityResponse {
	res := make([]*CapabilityResponse, len(caps))
	for i, c := range caps {
		res[i] = CapabilityFromEntity(c)
	}
	return res
}

// ─── Log ─────────────────────────────────────────────────────────────────────
type LogResponse struct {
	ID        uint      `json:"id"`
	AgentID   uint      `json:"agent_id"`
	TaskID    *uint     `json:"task_id,omitempty"`
	SessionID *uint     `json:"session_id,omitempty"`
	Level     string    `json:"level"`
	Action    string    `json:"action"`
	Message   string    `json:"message"`
	Detail    string    `json:"detail,omitempty"`
	Duration  int64     `json:"duration"`
	CreatedAt time.Time `json:"created_at"`
}

type LogListResponse struct {
	Logs  []*LogResponse `json:"logs"`
	Total int64          `json:"total"`
	Page  int            `json:"page"`
	Limit int            `json:"limit"`
}

func LogFromEntity(l *entity.AgentLog) *LogResponse {
	return &LogResponse{
		ID:        l.ID,
		AgentID:   l.AgentID,
		TaskID:    l.TaskID,
		SessionID: l.SessionID,
		Level:     l.Level,
		Action:    l.Action,
		Message:   l.Message,
		Detail:    l.Detail,
		Duration:  l.Duration,
		CreatedAt: l.CreatedAt,
	}
}

func LogFromEntities(logs []*entity.AgentLog) []*LogResponse {
	res := make([]*LogResponse, len(logs))
	for i, l := range logs {
		res[i] = LogFromEntity(l)
	}
	return res
}

// ─── Template ────────────────────────────────────────────────────────────────
type TemplateResponse struct {
	ID          uint      `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description,omitempty"`
	Category    string    `json:"category,omitempty"`
	Priority    int       `json:"priority"`
	MaxAttempts int       `json:"max_attempts"`
	Prompt      string    `json:"prompt"`
	Variables   string    `json:"variables,omitempty"`
	Tags        string    `json:"tags,omitempty"`
	IsPublic    bool      `json:"is_public"`
	CreatedBy   uint      `json:"created_by"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

func TemplateFromEntity(t *entity.AgentTemplate) *TemplateResponse {
	return &TemplateResponse{
		ID:          t.ID,
		Name:        t.Name,
		Description: t.Description,
		Category:    t.Category,
		Priority:    t.Priority,
		MaxAttempts: t.MaxAttempts,
		Prompt:      t.Prompt,
		Variables:   t.Variables,
		Tags:        t.Tags,
		IsPublic:    t.IsPublic,
		CreatedBy:   t.CreatedBy,
		CreatedAt:   t.CreatedAt,
		UpdatedAt:   t.UpdatedAt,
	}
}

func TemplateFromEntities(templates []*entity.AgentTemplate) []*TemplateResponse {
	res := make([]*TemplateResponse, len(templates))
	for i, t := range templates {
		res[i] = TemplateFromEntity(t)
	}
	return res
}

// ─── Schedule ────────────────────────────────────────────────────────────────
type ScheduleResponse struct {
	ID          uint       `json:"id"`
	AgentID     uint       `json:"agent_id"`
	TaskID      *uint      `json:"task_id,omitempty"`
	TemplateID  *uint      `json:"template_id,omitempty"`
	Name        string     `json:"name"`
	Description string     `json:"description,omitempty"`
	CronExpr    string     `json:"cron_expr"`
	Timezone    string     `json:"timezone"`
	Status      string     `json:"status"`
	MaxRuns     int        `json:"max_runs"`
	RunCount    int        `json:"run_count"`
	LastRunAt   *time.Time `json:"last_run_at,omitempty"`
	NextRunAt   *time.Time `json:"next_run_at,omitempty"`
	LastError   string     `json:"last_error,omitempty"`
	CreatedBy   uint       `json:"created_by"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
}

func ScheduleFromEntity(s *entity.AgentSchedule) *ScheduleResponse {
	return &ScheduleResponse{
		ID:          s.ID,
		AgentID:     s.AgentID,
		TaskID:      s.TaskID,
		TemplateID:  s.TemplateID,
		Name:        s.Name,
		Description: s.Description,
		CronExpr:    s.CronExpr,
		Timezone:    s.Timezone,
		Status:      s.Status,
		MaxRuns:     s.MaxRuns,
		RunCount:    s.RunCount,
		LastRunAt:   s.LastRunAt,
		NextRunAt:   s.NextRunAt,
		LastError:   s.LastError,
		CreatedBy:   s.CreatedBy,
		CreatedAt:   s.CreatedAt,
		UpdatedAt:   s.UpdatedAt,
	}
}

func ScheduleFromEntities(schedules []*entity.AgentSchedule) []*ScheduleResponse {
	res := make([]*ScheduleResponse, len(schedules))
	for i, s := range schedules {
		res[i] = ScheduleFromEntity(s)
	}
	return res
}

// ─── Stats ───────────────────────────────────────────────────────────────────
type AgentStatsResponse struct {
	TotalAgents    int64 `json:"total_agents"`
	IdleAgents     int64 `json:"idle_agents"`
	BusyAgents     int64 `json:"busy_agents"`
	TotalTasks     int64 `json:"total_tasks"`
	PendingTasks   int64 `json:"pending_tasks"`
	ActiveSessions int64 `json:"active_sessions"`
}
