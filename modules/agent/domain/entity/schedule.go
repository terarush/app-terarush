package entity

import "time"

type AgentSchedule struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	AgentID     uint      `gorm:"not null;index" json:"agent_id"`
	TaskID      *uint     `gorm:"index" json:"task_id"`
	TemplateID  *uint     `gorm:"index" json:"template_id"`
	Name        string    `gorm:"type:varchar(255);not null" json:"name"`
	Description string    `gorm:"type:text" json:"description"`
	CronExpr    string    `gorm:"type:varchar(100);not null" json:"cron_expr"`     // cron expression
	Timezone    string    `gorm:"type:varchar(50);default:'UTC'" json:"timezone"`
	Status      string    `gorm:"type:varchar(20);default:'active';index" json:"status"` // active, paused, completed, error
	MaxRuns     int       `gorm:"default:0" json:"max_runs"`                             // 0=unlimited
	RunCount    int       `gorm:"default:0" json:"run_count"`
	InputConfig string    `gorm:"type:jsonb" json:"input_config"` // static input for each run
	LastRunAt   *time.Time `json:"last_run_at"`
	NextRunAt   *time.Time `json:"next_run_at"`
	LastError   string    `gorm:"type:text" json:"last_error"`
	CreatedBy   uint      `gorm:"index" json:"created_by"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

func (AgentSchedule) TableName() string { return "agent_schedules" }

func NewSchedule(agentID uint, name, cronExpr string, createdBy uint) *AgentSchedule {
	now := time.Now()
	return &AgentSchedule{
		AgentID:   agentID,
		Name:      name,
		CronExpr:  cronExpr,
		Timezone:  "UTC",
		Status:    "active",
		CreatedBy: createdBy,
		CreatedAt: now,
		UpdatedAt: now,
	}
}
