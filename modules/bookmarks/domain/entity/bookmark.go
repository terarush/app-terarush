package entity

import "time"

type Bookmark struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	UserID      uint      `gorm:"not null;index;uniqueIndex:idx_user_resource" json:"user_id"`
	ResourceType string   `gorm:"type:varchar(50);not null;uniqueIndex:idx_user_resource" json:"resource_type"` // blog
	ResourceID  uint      `gorm:"not null;uniqueIndex:idx_user_resource" json:"resource_id"`
	CollectionID *uint    `gorm:"index" json:"collection_id"`
	Note        string    `gorm:"type:varchar(1000)" json:"note"`
	Tags        string    `gorm:"type:varchar(500)" json:"tags"` // comma-separated
	Progress    float64   `gorm:"default:0" json:"progress"`    // read progress 0-100
	Priority    string    `gorm:"type:varchar(20);default:'normal'" json:"priority"` // low, normal, high
	Status      string    `gorm:"type:varchar(20);default:'saved'" json:"status"`    // saved, reading, completed, archived
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

func (Bookmark) TableName() string { return "bookmarks" }

func NewBookmark(userID uint, resourceType string, resourceID uint) *Bookmark {
	return &Bookmark{
		UserID:       userID,
		ResourceType: resourceType,
		ResourceID:   resourceID,
		Priority:     "normal",
		Status:       "saved",
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}
}

type Collection struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	UserID      uint      `gorm:"not null;index" json:"user_id"`
	Name        string    `gorm:"type:varchar(255);not null" json:"name"`
	Description string    `gorm:"type:text" json:"description"`
	Icon        string    `gorm:"type:varchar(50)" json:"icon"`
	Color       string    `gorm:"type:varchar(7)" json:"color"` // hex color
	IsPublic    bool      `gorm:"default:false" json:"is_public"`
	SortOrder   int       `gorm:"default:0" json:"sort_order"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

func (Collection) TableName() string { return "collections" }

func NewCollection(userID uint, name string) *Collection {
	return &Collection{
		UserID:    userID,
		Name:      name,
		IsPublic:  false,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
}

type ReadingGoal struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	UserID    uint      `gorm:"uniqueIndex;not null" json:"user_id"`
	Year      int       `gorm:"not null" json:"year"`
	Goal      int       `gorm:"not null" json:"goal"` // target count
	Progress  int       `gorm:"default:0" json:"progress"`
	StartedAt time.Time `json:"started_at"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (ReadingGoal) TableName() string { return "reading_goals" }

func NewReadingGoal(userID uint, year, goal int) *ReadingGoal {
	return &ReadingGoal{
		UserID:    userID,
		Year:      year,
		Goal:      goal,
		StartedAt: time.Now(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
}
