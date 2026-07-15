package entity

import "time"

type AnalyticsEvent struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	EventType   string    `gorm:"type:varchar(50);not null;index" json:"event_type"` // page_view, blog_view, blog_create, comment, search, auth, click
	UserID      *uint     `gorm:"index" json:"user_id"`
	SessionID   string    `gorm:"type:varchar(100);index" json:"session_id"`
	IP          string    `gorm:"type:varchar(45)" json:"ip"`
	UserAgent   string    `gorm:"type:varchar(500)" json:"user_agent"`
	Referrer    string    `gorm:"type:varchar(1000)" json:"referrer"`
	Path        string    `gorm:"type:varchar(500);index" json:"path"`
	Method      string    `gorm:"type:varchar(10)" json:"method"`
	StatusCode  int       `json:"status_code"`
	Duration    int64     `json:"duration"`   // ms
	Metadata    string    `gorm:"type:jsonb" json:"metadata"`
	ResourceType string   `gorm:"type:varchar(50)" json:"resource_type"` // blog, user, comment
	ResourceID  *uint     `gorm:"index" json:"resource_id"`
	Country     string    `gorm:"type:varchar(5)" json:"country"`
	Device      string    `gorm:"type:varchar(20)" json:"device"` // desktop, tablet, mobile
	Browser     string    `gorm:"type:varchar(50)" json:"browser"`
	OS          string    `gorm:"type:varchar(50)" json:"os"`
	CreatedAt   time.Time `gorm:"index" json:"created_at"`
}

func (AnalyticsEvent) TableName() string { return "analytics_events" }

type AnalyticsAggregate struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Date      string    `gorm:"type:date;not null;index" json:"date"`
	Metric    string    `gorm:"type:varchar(50);not null;index" json:"metric"` // page_views, unique_visitors, blog_views, signups, comments
	Value     int64     `gorm:"default:0" json:"value"`
	Dimension string    `gorm:"type:varchar(100)" json:"dimension"` // blog_id, user_id, path, country, device
	DimensionValue string `gorm:"type:varchar(255)" json:"dimension_value"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (AnalyticsAggregate) TableName() string { return "analytics_aggregates" }

type BlogStats struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	BlogID      uint      `gorm:"uniqueIndex;not null" json:"blog_id"`
	ViewCount   int64     `gorm:"default:0" json:"view_count"`
	UniqueViews int64     `gorm:"default:0" json:"unique_views"`
	ReadTimeAvg float64   `gorm:"default:0" json:"read_time_avg"` // seconds
	ShareCount  int64     `gorm:"default:0" json:"share_count"`
	CommentCount int64    `gorm:"default:0" json:"comment_count"`
	FavoriteCount int64   `gorm:"default:0" json:"favorite_count"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

func (BlogStats) TableName() string { return "blog_stats" }

type UserAnalytics struct {
	ID            uint      `gorm:"primaryKey" json:"id"`
	UserID        uint      `gorm:"uniqueIndex;not null" json:"user_id"`
	TotalViews    int64     `gorm:"default:0" json:"total_views"`
	TotalBlogs    int64     `gorm:"default:0" json:"total_blogs"`
	TotalComments int64     `gorm:"default:0" json:"total_comments"`
	TotalFavorites int64    `gorm:"default:0" json:"total_favorites"`
	LastActiveAt  *time.Time `json:"last_active_at"`
	StreakDays    int       `gorm:"default:0" json:"streak_days"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

func (UserAnalytics) TableName() string { return "user_analytics" }

// Event types
const (
	EvPageView  = "page_view"
	EvBlogView  = "blog_view"
	EvBlogCreate = "blog_create"
	EvSearch    = "search"
	EvAuth      = "auth"
)
