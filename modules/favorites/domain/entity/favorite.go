package entity

import "time"

type Favorite struct {
	ID        uint  `gorm:"primaryKey" json:"id"`
	UserID    uint  `json:"user_id"`
	PostID    uint  `json:"post_id"`
	CreatedAt int64 `json:"created_at"`
}

func (*Favorite) TableName() string {
	return "favorites"
}

func NewFavorite(userID, postID uint) *Favorite {
	time := time.Now().Unix()
	return &Favorite{
		UserID:    userID,
		PostID:    postID,
		CreatedAt: time,
	}
}
