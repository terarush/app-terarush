package entity

import "time"

type Comment struct {
	ID        string `gorm:"primaryKey" json:"id"`
	Content   string `gorm:"type:text" json:"content"`
	UserID    uint   `gorm:"index" json:"user_id"`
	PostID    uint   `gorm:"index" json:"post_id"`
	CreatedAt int64  `json:"created_at"`
	UpdatedAt int64  `json:"updated_at"`
}

func (*Comment) TableName() string {
	return "comments"
}

func NewComment(content string, userID, postID uint) *Comment {
	time := time.Now().Unix()
	return &Comment{
		Content:   content,
		UserID:    userID,
		PostID:    postID,
		CreatedAt: time,
		UpdatedAt: time,
	}
}
