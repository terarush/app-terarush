package entity

import "time"

type AgentLog struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	AgentID   uint      `gorm:"not null;index" json:"agent_id"`
	TaskID    *uint     `gorm:"index" json:"task_id"`
	SessionID *uint     `gorm:"index" json:"session_id"`
	Level     string    `gorm:"type:varchar(20);default:'info';index" json:"level"` // debug, info, warn, error, fatal
	Action    string    `gorm:"type:varchar(100);index" json:"action"`              // task_started, task_completed, error, warning, state_change
	Message   string    `gorm:"type:text" json:"message"`
	Detail    string    `gorm:"type:jsonb" json:"detail"` // structured log payload
	Duration  int64     `json:"duration"`                  // ms
	Metadata  string    `gorm:"type:jsonb" json:"metadata"`
	CreatedAt time.Time `gorm:"index" json:"created_at"`
}

func (AgentLog) TableName() string { return "agent_logs" }

func NewLog(agentID uint, level, action, message string) *AgentLog {
	return &AgentLog{
		AgentID:   agentID,
		Level:     level,
		Action:    action,
		Message:   message,
		CreatedAt: time.Now(),
	}
}
