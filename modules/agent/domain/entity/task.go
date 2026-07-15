package entity

import "time"

type AgentTask struct {
	ID            uint       `gorm:"primaryKey" json:"id"`
	Title         string     `gorm:"type:varchar(500);not null" json:"title"`
	Description   string     `gorm:"type:text" json:"description"`
	Status        string     `gorm:"type:varchar(20);default:'pending';index" json:"status"` // pending, assigned, in_progress, completed, failed, cancelled, paused
	Priority      int        `gorm:"default:0;index" json:"priority"`                        // 0=normal, 1=high, 2=critical, -1=low
	Category      string     `gorm:"type:varchar(100);index" json:"category"`
	Tags          string     `gorm:"type:varchar(500)" json:"tags"`                          // comma-separated
	AssignedTo    *uint      `gorm:"index" json:"assigned_to"`                               // FK -> Agent
	AssignedBy    *uint      `json:"assigned_by"`                                            // FK -> User
	ParentTaskID  *uint      `gorm:"index" json:"parent_task_id"`                            // sub-task
	Input         string     `gorm:"type:jsonb" json:"input"`
	Output        string     `gorm:"type:jsonb" json:"output"`
	ResultSummary string     `gorm:"type:text" json:"result_summary"`
	ErrorMessage  string     `gorm:"type:text" json:"error_message"`
	Attempts      int        `gorm:"default:0" json:"attempts"`
	MaxAttempts   int        `gorm:"default:3" json:"max_attempts"`
	Progress      float64    `gorm:"default:0" json:"progress"` // 0-100
	Dependencies  string     `gorm:"type:varchar(500)" json:"dependencies"` // comma-separated task IDs
	Deadline      *time.Time `json:"deadline"`
	StartedAt     *time.Time `json:"started_at"`
	CompletedAt   *time.Time `json:"completed_at"`
	CreatedBy     uint       `gorm:"index" json:"created_by"`
	CreatedAt     time.Time  `json:"created_at"`
	UpdatedAt     time.Time  `json:"updated_at"`
}

func (AgentTask) TableName() string { return "agent_tasks" }

func NewTask(title, description, category string, priority int, createdBy uint) *AgentTask {
	now := time.Now()
	return &AgentTask{
		Title:       title,
		Description: description,
		Status:      "pending",
		Priority:    priority,
		Category:    category,
		MaxAttempts: 3,
		CreatedBy:   createdBy,
		CreatedAt:   now,
		UpdatedAt:   now,
	}
}
