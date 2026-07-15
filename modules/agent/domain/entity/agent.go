package entity

import "time"

type Agent struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	Name        string    `gorm:"type:varchar(255);not null;uniqueIndex" json:"name"`
	Description string    `gorm:"type:text" json:"description"`
	Model       string    `gorm:"type:varchar(100);default:'gpt-4'" json:"model"`
	SystemPrompt string   `gorm:"type:text" json:"system_prompt"`
	Status      string    `gorm:"type:varchar(20);default:'idle';index" json:"status"` // idle, busy, error, disabled
	MaxConcurrent int     `gorm:"default:1" json:"max_concurrent"`
	Timeout     int64     `gorm:"default:300" json:"timeout"`                          // seconds
	RetryLimit  int       `gorm:"default:3" json:"retry_limit"`
	Metadata    string    `gorm:"type:jsonb" json:"metadata"`
	CreatedBy   uint      `gorm:"index" json:"created_by"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

func (Agent) TableName() string { return "agents" }

func NewAgent(name, description, model, systemPrompt, createdBy string) *Agent {
	now := time.Now()
	return &Agent{
		Name:         name,
		Description:  description,
		Model:        model,
		SystemPrompt: systemPrompt,
		Status:       "idle",
		MaxConcurrent: 1,
		Timeout:      300,
		RetryLimit:   3,
		CreatedAt:    now,
		UpdatedAt:    now,
	}
}
