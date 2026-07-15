package response

import (
	"time"
)

type ActivityResponse struct {
	ID         uint      `json:"id"`
	Type       string    `json:"type"`
	Verb       string    `json:"verb"`
	ActorID    uint      `json:"actor_id"`
	TargetType string    `json:"target_type,omitempty"`
	TargetID   *uint     `json:"target_id,omitempty"`
	CreatedAt  time.Time `json:"created_at"`
}

type FeedEntryResponse struct {
	ID        uint      `json:"id"`
	ActorID   uint      `json:"actor_id"`
	Type      string    `json:"type"`
	Verb      string    `json:"verb"`
	IsRead    bool      `json:"is_read"`
	CreatedAt time.Time `json:"created_at"`
}

type FollowResponse struct {
	UserID   uint      `json:"user_id"`
	Since    time.Time `json:"since"`
}

type MetricResponse struct {
	Posts     int64 `json:"posts"`
	Followers int64 `json:"followers"`
	Following int64 `json:"following"`
	Views     int64 `json:"views"`
	XP        int64 `json:"xp"`
	Level     int   `json:"level"`
}
