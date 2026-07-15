package response

import (
	"time"
)

type AnalyticsEventResponse struct {
	ID        uint      `json:"id"`
	EventType string    `json:"event_type"`
	UserID    *uint     `json:"user_id,omitempty"`
	Path      string    `json:"path,omitempty"`
	Duration  int64     `json:"duration"`
	Resource  string    `json:"resource,omitempty"`
	CreatedAt time.Time `json:"created_at"`
}

type BlogStatsResponse struct {
	BlogID    uint   `json:"blog_id"`
	ViewCount int64  `json:"view_count"`
	FavCount  int64  `json:"favorite_count"`
	CmtCount  int64  `json:"comment_count"`
}

type DashboardStatsResponse struct {
	RecentViews int64               `json:"recent_views"`
	TopBlogs    []*BlogStatsResponse `json:"top_blogs"`
}
