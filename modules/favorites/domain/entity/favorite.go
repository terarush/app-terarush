package entity

import "time"

type Favorite struct {
	ID        string `gorm:"primaryKey" json:"id"`
	UserID    string `json:"user_id"`
	PostID    string `json:"post_id"`
	CreatedAt int64  `json:"created_at"`
}

func (*Favorite) TableName() string {
	return "favorites"
}

func NewFavorite(userID, postID string) *Favorite {
	time := time.Now().Unix()
	return &Favorite{
		UserID:    userID,
		PostID:    postID,
		CreatedAt: time,
	}
}
