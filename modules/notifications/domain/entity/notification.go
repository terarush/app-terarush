package entity

import "time"

type Notification struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	UserID    uint      `gorm:"not null;index" json:"user_id"`
	Type      string    `gorm:"type:varchar(50);not null;index" json:"type"` // blog_published, comment_reply, favorite, mention, system, digest
	Title     string    `gorm:"type:varchar(500);not null" json:"title"`
	Body      string    `gorm:"type:text" json:"body"`
	Link      string    `gorm:"type:varchar(500)" json:"link"`
	Icon      string    `gorm:"type:varchar(100)" json:"icon"`
	ReferenceType string `gorm:"type:varchar(50)" json:"reference_type"` // blog, comment, user
	ReferenceID   *uint  `json:"reference_id"`
	IsRead    bool      `gorm:"default:false;index" json:"is_read"`
	IsArchived bool     `gorm:"default:false" json:"is_archived"`
	Priority  string    `gorm:"type:varchar(20);default:'normal'" json:"priority"` // low, normal, high, urgent
	Channel   string    `gorm:"type:varchar(50);default:'in_app'" json:"channel"` // in_app, email, push
	SentAt    *time.Time `json:"sent_at"`
	ReadAt    *time.Time `json:"read_at"`
	CreatedAt time.Time  `gorm:"index" json:"created_at"`
}

func (Notification) TableName() string { return "notifications" }

func NewNotification(userID uint, notifType, title, body string) *Notification {
	now := time.Now()
	return &Notification{
		UserID:    userID,
		Type:      notifType,
		Title:     title,
		Body:      body,
		Priority:  "normal",
		Channel:   "in_app",
		SentAt:    &now,
		CreatedAt: now,
	}
}

type NotificationTemplate struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	Name        string    `gorm:"type:varchar(255);uniqueIndex;not null" json:"name"`
	Description string    `gorm:"type:text" json:"description"`
	TitleTmpl   string    `gorm:"type:varchar(500);not null" json:"title_tmpl"`
	BodyTmpl    string    `gorm:"type:text;not null" json:"body_tmpl"`
	DefaultIcon string    `gorm:"type:varchar(100)" json:"default_icon"`
	Channels    string    `gorm:"type:varchar(200);default:'in_app'" json:"channels"` // comma-separated
	Priority    string    `gorm:"type:varchar(20);default:'normal'" json:"priority"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

func (NotificationTemplate) TableName() string { return "notification_templates" }

type NotificationPreference struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	UserID    uint      `gorm:"uniqueIndex;not null" json:"user_id"`
	EmailDigest string `gorm:"type:varchar(20);default:'daily'" json:"email_digest"` // none, instant, daily, weekly
	Types     string    `gorm:"type:jsonb" json:"types"` // {"blog_published": true, "comment_reply": true, ...}
	Channels  string    `gorm:"type:varchar(200)" json:"channels"` // comma-separated enabled channels
	QuietHours string   `gorm:"type:jsonb" json:"quiet_hours"` // {"start":"22:00","end":"08:00","timezone":"UTC"}
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (NotificationPreference) TableName() string { return "notification_preferences" }

// Notification types
const (
	NotifBlogPublished = "blog_published"
	NotifCommentReply  = "comment_reply"
	NotifFavorite      = "favorite"
	NotifMention       = "mention"
	NotifFollow        = "follow"
	NotifSystem        = "system"
	NotifModeration    = "moderation"
	NotifDigest        = "digest"
)
