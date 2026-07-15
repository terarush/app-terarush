package entity

import "time"

type AgentCapability struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	AgentID     uint      `gorm:"not null;index;uniqueIndex:idx_agent_cap" json:"agent_id"`
	Capability  string    `gorm:"type:varchar(100);not null;uniqueIndex:idx_agent_cap" json:"capability"` // read, write, execute, search, analyze, code, etc.
	Enabled     bool      `gorm:"default:true" json:"enabled"`
	Config      string    `gorm:"type:jsonb" json:"config"`    // per-capability config
	MaxUsage    int       `gorm:"default:0" json:"max_usage"`  // 0=unlimited
	UsageCount  int       `gorm:"default:0" json:"usage_count"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

func (AgentCapability) TableName() string { return "agent_capabilities" }

func NewCapability(agentID uint, capability string) *AgentCapability {
	now := time.Now()
	return &AgentCapability{
		AgentID:    agentID,
		Capability: capability,
		Enabled:    true,
		CreatedAt:  now,
		UpdatedAt:  now,
	}
}

// Available capabilities
const (
	CapRead      = "read"
	CapWrite     = "write"
	CapExecute   = "execute"
	CapSearch    = "search"
	CapAnalyze   = "analyze"
	CapCode      = "code"
	CapIntegrate = "integrate"
	CapMonitor   = "monitor"
	CapSchedule  = "schedule"
	CapNotify    = "notify"
)
