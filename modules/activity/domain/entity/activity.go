package entity

import "time"

type Activity struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	UserID    uint      `gorm:"not null;index" json:"user_id"`
	Type      string    `gorm:"type:varchar(50);not null;index" json:"type"` // create_blog, publish_blog, comment, favorite, follow, update_profile, etc.
	Verb      string    `gorm:"type:varchar(50);not null" json:"verb"`       // created, published, commented, favorited, followed, updated
	ActorID   uint      `gorm:"index" json:"actor_id"`                       // who performed the action
	TargetType string   `gorm:"type:varchar(50)" json:"target_type"`          // blog, comment, user
	TargetID   *uint    `json:"target_id"`
	ParentType string   `gorm:"type:varchar(50)" json:"parent_type"`
	ParentID   *uint    `json:"parent_id"`
	Metadata   string   `gorm:"type:jsonb" json:"metadata"`
	Visibility string   `gorm:"type:varchar(20);default:'public';index" json:"visibility"` // public, followers, private
	CreatedAt  time.Time `gorm:"index" json:"created_at"`
}

func (Activity) TableName() string { return "activities" }

func NewActivity(userID uint, activityType, verb string) *Activity {
	return &Activity{
		UserID:     userID,
		Type:       activityType,
		Verb:       verb,
		ActorID:    userID,
		Visibility: "public",
		CreatedAt:  time.Now(),
	}
}

type ActivityFeed struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	UserID    uint      `gorm:"not null;index" json:"user_id"`          // feed owner
	ActorID   uint      `gorm:"index" json:"actor_id"`                  // who did it
	Type      string    `gorm:"type:varchar(50);index" json:"type"`
	Verb      string    `gorm:"type:varchar(50)" json:"verb"`
	TargetType string   `gorm:"type:varchar(50)" json:"target_type"`
	TargetID   *uint    `json:"target_id"`
	Metadata   string   `gorm:"type:jsonb" json:"metadata"`
	IsRead     bool     `gorm:"default:false" json:"is_read"`
	CreatedAt  time.Time `gorm:"index" json:"created_at"`
}

func (ActivityFeed) TableName() string { return "activity_feeds" }

type Follow struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	FollowerID  uint      `gorm:"not null;index;uniqueIndex:idx_follow" json:"follower_id"`
	FollowingID uint      `gorm:"not null;index;uniqueIndex:idx_follow" json:"following_id"`
	CreatedAt   time.Time `json:"created_at"`
}

func (Follow) TableName() string { return "follows" }

type EngagementMetric struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	UserID    uint      `gorm:"uniqueIndex;not null" json:"user_id"`
	Posts     int64     `gorm:"default:0" json:"posts"`
	Comments  int64     `gorm:"default:0" json:"comments"`
	Favorites int64     `gorm:"default:0" json:"favorites"`
	Followers int64     `gorm:"default:0" json:"followers"`
	Following int64     `gorm:"default:0" json:"following"`
	Likes     int64     `gorm:"default:0" json:"likes"`
	Views     int64     `gorm:"default:0" json:"views"`
	Streak    int       `gorm:"default:0" json:"streak"`
	Level     int       `gorm:"default:1" json:"level"`
	XP        int64     `gorm:"default:0" json:"xp"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (EngagementMetric) TableName() string { return "engagement_metrics" }

// Activity types
const (
	ActCreateBlog    = "create_blog"
	ActPublishBlog   = "publish_blog"
	ActComment       = "comment"
	ActFavorite      = "favorite"
	ActFollow        = "follow"
	ActUpdateProfile = "update_profile"
)
