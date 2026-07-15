package entity

import "time"

type Report struct {
	ID            uint      `gorm:"primaryKey" json:"id"`
	ReporterID    uint      `gorm:"not null;index" json:"reporter_id"`
	ResourceType  string    `gorm:"type:varchar(50);not null;index" json:"resource_type"` // blog, comment, user
	ResourceID    uint      `gorm:"not null;index" json:"resource_id"`
	Reason        string    `gorm:"type:varchar(100);not null" json:"reason"` // spam, harassment, violence, copyright, illegal, other
	Description   string    `gorm:"type:text" json:"description"`
	Status        string    `gorm:"type:varchar(20);default:'pending';index" json:"status"` // pending, reviewing, resolved, dismissed
	Severity      string    `gorm:"type:varchar(20);default:'medium'" json:"severity"`      // low, medium, high, critical
	AssignedTo    *uint     `gorm:"index" json:"assigned_to"` // moderator user ID
	ReviewNotes   string    `gorm:"type:text" json:"review_notes"`
	ReviewedBy    *uint     `json:"reviewed_by"`
	ActionTaken   string    `gorm:"type:varchar(50)" json:"action_taken"` // warned, content_removed, user_banned, user_suspended, dismissed
	ActionDetail  string    `gorm:"type:text" json:"action_detail"`
	ReviewedAt    *time.Time `json:"reviewed_at"`
	CreatedAt     time.Time  `gorm:"index" json:"created_at"`
	UpdatedAt     time.Time  `json:"updated_at"`
}

func (Report) TableName() string { return "moderation_reports" }

type ModerationAction struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	ModeratorID uint      `gorm:"not null;index" json:"moderator_id"`
	UserID      uint      `gorm:"not null;index" json:"user_id"`
	Action      string    `gorm:"type:varchar(50);not null" json:"action"` // warn, suspend, ban, unban, mute, unmute
	Reason      string    `gorm:"type:varchar(500);not null" json:"reason"`
	Duration    string    `gorm:"type:varchar(50)" json:"duration"` // permanent, 1d, 3d, 7d, 30d
	ExpiresAt   *time.Time `json:"expires_at"`
	AutoTrigger bool      `gorm:"default:false" json:"auto_trigger"` // system auto-triggered
	CreatedAt   time.Time  `json:"created_at"`
}

func (ModerationAction) TableName() string { return "moderation_actions" }

type BannedUser struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	UserID      uint      `gorm:"uniqueIndex;not null" json:"user_id"`
	Reason      string    `gorm:"type:varchar(500);not null" json:"reason"`
	BannedBy    uint      `json:"banned_by"`
	IsPermanent bool      `gorm:"default:true" json:"is_permanent"`
	ExpiresAt   *time.Time `json:"expires_at"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
}

func (BannedUser) TableName() string { return "banned_users" }

type ContentFilter struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Type      string    `gorm:"type:varchar(50);not null;index" json:"type"` // keyword, regex, domain, ip
	Pattern   string    `gorm:"type:varchar(500);not null" json:"pattern"`
	Action    string    `gorm:"type:varchar(20);default:'flag'" json:"action"` // flag, block, review, warn
	Severity  string    `gorm:"type:varchar(20);default:'medium'" json:"severity"`
	CreatedBy uint      `json:"created_by"`
	IsActive  bool      `gorm:"default:true" json:"is_active"`
	HitCount  int       `gorm:"default:0" json:"hit_count"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (ContentFilter) TableName() string { return "content_filters" }

// Report reasons
const (
	ReasonSpam       = "spam"
	ReasonHarassment = "harassment"
	ReasonViolence   = "violence"
	ReasonCopyright  = "copyright"
	ReasonIllegal    = "illegal"
	ReasonOther      = "other"
)
