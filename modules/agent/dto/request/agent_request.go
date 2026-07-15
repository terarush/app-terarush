package request

// ─── Agent ───────────────────────────────────────────────────────────────────
type CreateAgentRequest struct {
	Name         string `json:"name" validate:"required,max=255"`
	Description  string `json:"description" validate:"omitempty,max=2000"`
	Model        string `json:"model" validate:"omitempty,max=100"`
	SystemPrompt string `json:"system_prompt" validate:"omitempty,max=10000"`
	MaxConcurrent int   `json:"max_concurrent" validate:"omitempty,min=1,max=10"`
	Timeout      int64  `json:"timeout" validate:"omitempty,min=1,max=86400"`
	RetryLimit   int    `json:"retry_limit" validate:"omitempty,min=0,max=10"`
}

type UpdateAgentRequest struct {
	Name         string `json:"name" validate:"omitempty,max=255"`
	Description  string `json:"description" validate:"omitempty,max=2000"`
	Model        string `json:"model" validate:"omitempty,max=100"`
	SystemPrompt string `json:"system_prompt" validate:"omitempty,max=10000"`
	Status       string `json:"status" validate:"omitempty,oneof=idle busy error disabled"`
	MaxConcurrent int   `json:"max_concurrent" validate:"omitempty,min=1,max=10"`
	Timeout      int64  `json:"timeout" validate:"omitempty,min=1,max=86400"`
	RetryLimit   int    `json:"retry_limit" validate:"omitempty,min=0,max=10"`
}

type SetAgentStatusRequest struct {
	Status string `json:"status" validate:"required,oneof=idle busy error disabled"`
}

// ─── Task ────────────────────────────────────────────────────────────────────
type CreateTaskRequest struct {
	Title        string `json:"title" validate:"required,max=500"`
	Description  string `json:"description" validate:"omitempty,max=10000"`
	Priority     int    `json:"priority" validate:"omitempty,oneof=-1 0 1 2"`
	Category     string `json:"category" validate:"omitempty,max=100"`
	Tags         string `json:"tags" validate:"omitempty,max=500"`
	Input        string `json:"input" validate:"omitempty"`
	MaxAttempts  int    `json:"max_attempts" validate:"omitempty,min=1,max=10"`
	ParentTaskID *uint  `json:"parent_task_id" validate:"omitempty"`
	Deadline     string `json:"deadline" validate:"omitempty"` // RFC3339
}

type UpdateTaskRequest struct {
	Title        string  `json:"title" validate:"omitempty,max=500"`
	Description  string  `json:"description" validate:"omitempty,max=10000"`
	Priority     int     `json:"priority" validate:"omitempty,oneof=-1 0 1 2"`
	Category     string  `json:"category" validate:"omitempty,max=100"`
	Tags         string  `json:"tags" validate:"omitempty,max=500"`
	Status       string  `json:"status" validate:"omitempty,oneof=pending assigned in_progress completed failed cancelled paused"`
	Progress     float64 `json:"progress" validate:"omitempty,min=0,max=100"`
	Input        string  `json:"input" validate:"omitempty"`
	MaxAttempts  int     `json:"max_attempts" validate:"omitempty,min=1,max=10"`
}

type AssignTaskRequest struct {
	AgentID uint `json:"agent_id" validate:"required"`
}

type CompleteTaskRequest struct {
	Output        string `json:"output" validate:"omitempty"`
	ResultSummary string `json:"result_summary" validate:"omitempty,max=5000"`
}

type FailTaskRequest struct {
	ErrorMessage string `json:"error_message" validate:"required"`
}

// ─── Session ─────────────────────────────────────────────────────────────────
type CreateSessionRequest struct {
	AgentID uint  `json:"agent_id" validate:"required"`
	TaskID  *uint `json:"task_id" validate:"omitempty"`
}

// ─── Capability ──────────────────────────────────────────────────────────────
type AddCapabilityRequest struct {
	AgentID    uint   `json:"agent_id" validate:"required"`
	Capability string `json:"capability" validate:"required,oneof=read write execute search analyze code integrate monitor schedule notify"`
	Config     string `json:"config" validate:"omitempty"`
}

type ToggleCapabilityRequest struct {
	Enabled bool `json:"enabled" validate:"required"`
}

// ─── Template ────────────────────────────────────────────────────────────────
type CreateTemplateRequest struct {
	Name        string `json:"name" validate:"required,max=255"`
	Description string `json:"description" validate:"omitempty,max=2000"`
	Category    string `json:"category" validate:"omitempty,max=100"`
	Priority    int    `json:"priority" validate:"omitempty,oneof=-1 0 1 2"`
	MaxAttempts int    `json:"max_attempts" validate:"omitempty,min=1,max=10"`
	InputSchema string `json:"input_schema" validate:"omitempty"`
	Prompt      string `json:"prompt" validate:"required"`
	Variables   string `json:"variables" validate:"omitempty"`
	Tags        string `json:"tags" validate:"omitempty,max=500"`
	IsPublic    bool   `json:"is_public"`
}

type CreateTaskFromTemplateRequest struct {
	TemplateID uint   `json:"template_id" validate:"required"`
	Input      string `json:"input" validate:"omitempty"`
}

// ─── Schedule ────────────────────────────────────────────────────────────────
type CreateScheduleRequest struct {
	AgentID     uint   `json:"agent_id" validate:"required"`
	TemplateID  *uint  `json:"template_id" validate:"omitempty"`
	Name        string `json:"name" validate:"required,max=255"`
	Description string `json:"description" validate:"omitempty,max=2000"`
	CronExpr    string `json:"cron_expr" validate:"required,max=100"`
	Timezone    string `json:"timezone" validate:"omitempty,max=50"`
	MaxRuns     int    `json:"max_runs" validate:"omitempty,min=0"`
	InputConfig string `json:"input_config" validate:"omitempty"`
}
