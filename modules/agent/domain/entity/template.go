package entity

import "time"

type AgentTemplate struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	Name        string    `gorm:"type:varchar(255);not null;uniqueIndex" json:"name"`
	Description string    `gorm:"type:text" json:"description"`
	Category    string    `gorm:"type:varchar(100);index" json:"category"`
	Priority    int       `gorm:"default:0" json:"priority"`
	MaxAttempts int       `gorm:"default:3" json:"max_attempts"`
	InputSchema string    `gorm:"type:jsonb" json:"input_schema"`   // expected input structure
	Prompt      string    `gorm:"type:text;not null" json:"prompt"` // template prompt
	Variables   string    `gorm:"type:jsonb" json:"variables"`      // template variables
	Tags        string    `gorm:"type:varchar(500)" json:"tags"`
	IsPublic    bool      `gorm:"default:false" json:"is_public"`
	CreatedBy   uint      `gorm:"index" json:"created_by"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

func (AgentTemplate) TableName() string { return "agent_templates" }

func NewTemplate(name, description, category, prompt string, createdBy uint) *AgentTemplate {
	now := time.Now()
	return &AgentTemplate{
		Name:        name,
		Description: description,
		Category:    category,
		MaxAttempts: 3,
		Prompt:      prompt,
		CreatedBy:   createdBy,
		CreatedAt:   now,
		UpdatedAt:   now,
	}
}
