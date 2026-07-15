package entity

import "time"

type Subscription struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	UserID      uint      `gorm:"not null;index;uniqueIndex:idx_user_resource" json:"user_id"`
	ResourceType string   `gorm:"type:varchar(50);not null;uniqueIndex:idx_user_resource" json:"resource_type"` // author, blog, topic, newsletter
	ResourceID  uint      `gorm:"not null;uniqueIndex:idx_user_resource" json:"resource_id"`
	Status      string    `gorm:"type:varchar(20);default:'active';index" json:"status"` // active, paused, cancelled
	NotifyTypes string    `gorm:"type:varchar(200)" json:"notify_types"` // comma-separated: new_blog, comment, digest
	Frequency   string    `gorm:"type:varchar(20);default:'instant'" json:"frequency"` // instant, daily, weekly
	LastNotifiedAt *time.Time `json:"last_notified_at"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
}

func (Subscription) TableName() string { return "subscriptions" }

func NewSubscription(userID uint, resourceType string, resourceID uint) *Subscription {
	return &Subscription{
		UserID:       userID,
		ResourceType: resourceType,
		ResourceID:   resourceID,
		Status:       "active",
		Frequency:    "instant",
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}
}

type Newsletter struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	Title       string    `gorm:"type:varchar(255);not null" json:"title"`
	Content     string    `gorm:"type:text;not null" json:"content"`
	Summary     string    `gorm:"type:text" json:"summary"`
	Status      string    `gorm:"type:varchar(20);default:'draft';index" json:"status"` // draft, queued, sending, sent, failed
	ScheduledFor *time.Time `json:"scheduled_for"`
	SentCount   int       `gorm:"default:0" json:"sent_count"`
	OpenCount   int       `gorm:"default:0" json:"open_count"`
	ClickCount  int       `gorm:"default:0" json:"click_count"`
	BounceCount int       `gorm:"default:0" json:"bounce_count"`
	CreatedBy   uint      `json:"created_by"`
	SentAt      *time.Time `json:"sent_at"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
}

func (Newsletter) TableName() string { return "newsletters" }

type NewsletterSubscriber struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	UserID    *uint     `gorm:"index" json:"user_id"`
	Email     string    `gorm:"type:varchar(255);not null;uniqueIndex" json:"email"`
	Name      string    `gorm:"type:varchar(255)" json:"name"`
	Status    string    `gorm:"type:varchar(20);default:'active';index" json:"status"` // active, unsubscribed, bounced
	Token     string    `gorm:"type:varchar(100);uniqueIndex" json:"token"`
	SubscribedAt time.Time `json:"subscribed_at"`
	UnsubscribedAt *time.Time `json:"unsubscribed_at"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (NewsletterSubscriber) TableName() string { return "newsletter_subscribers" }

// Resource types
const (
	SubAuthor    = "author"
	SubBlog      = "blog"
	SubTopic     = "topic"
	SubNewsletter = "newsletter"
)
