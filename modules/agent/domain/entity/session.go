package entity

import "time"

type AgentSession struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	AgentID   uint      `gorm:"not null;index" json:"agent_id"`
	TaskID    *uint     `gorm:"index" json:"task_id"`
	Status    string    `gorm:"type:varchar(20);default:'active';index" json:"status"` // active, paused, completed, expired, terminated
	UserID    uint      `gorm:"index" json:"user_id"`
	Context   string    `gorm:"type:jsonb" json:"context"`   // session context/memory snapshot
	TokenUsed int64     `gorm:"default:0" json:"token_used"`
	Duration  int64     `gorm:"default:0" json:"duration"`   // seconds
	Messages  int       `gorm:"default:0" json:"messages"`
	Metadata  string    `gorm:"type:jsonb" json:"metadata"`
	ExpiresAt *time.Time `json:"expires_at"`
	StartedAt time.Time  `json:"started_at"`
	EndedAt   *time.Time `json:"ended_at"`
	CreatedAt time.Time  `json:"created_at"`
}

func (AgentSession) TableName() string { return "agent_sessions" }

func NewSession(agentID uint, userID uint, taskID *uint) *AgentSession {
	now := time.Now()
	return &AgentSession{
		AgentID:    agentID,
		TaskID:     taskID,
		Status:     "active",
		UserID:     userID,
		StartedAt:  now,
		CreatedAt:  now,
	}
}
